package msgpack

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
)

type Reader struct {
	input  []byte
	offset int
}

type Unmarshaler interface {
	UnmarshalMsgpack(r *Reader) error
}

func NewReader(input []byte) *Reader {
	return &Reader{input: input, offset: 0}
}

var readBuf = make([]byte, 8)

func (r *Reader) readByte() (byte, error) {
	b := r.input[r.offset]
	r.offset++
	return b, nil
}

func (r *Reader) readBytes(n int) ([]byte, error) {
	end := r.offset + n
	if end > len(r.input) {
		return nil, fmt.Errorf("read out of bounds")
	}
	buf := r.input[r.offset : r.offset+n]
	r.offset += n
	return buf, nil
}

func (r *Reader) readUint8() (uint8, error) {
	b, err := r.readByte()
	return uint8(b), err
}

func (r *Reader) readUint16() (uint16, error) {
	buf, err := r.readBytes(2)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(buf), nil
}

func (r *Reader) readUint32() (uint32, error) {
	buf, err := r.readBytes(4)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(buf), nil
}

func (r *Reader) readUint64() (uint64, error) {
	buf, err := r.readBytes(8)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(buf), nil
}

func (r *Reader) readInt8() (int8, error) {
	b, err := r.readByte()
	return int8(b), err
}

func (r *Reader) readInt16() (int16, error) {
	v, err := r.readUint16()
	return int16(v), err
}

func (r *Reader) readInt32() (int32, error) {
	v, err := r.readUint32()
	return int32(v), err
}

func (r *Reader) readInt64() (int64, error) {
	v, err := r.readUint64()
	return int64(v), err
}

func (r *Reader) readFloat32() (float32, error) {
	v, err := r.readUint32()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(v), nil
}

func (r *Reader) readFloat64() (float64, error) {
	v, err := r.readUint64()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(v), nil
}

// ReadNil reads a nil value
func (r *Reader) ReadNil() error {
	b, err := r.readByte()
	if err != nil {
		return err
	}
	if b != Nil {
		return fmt.Errorf("expected nil (0x%02x) but got 0x%02x", Nil, b)
	}
	return nil
}

// ReadBool reads a boolean value
func (r *Reader) ReadBool() (bool, error) {
	b, err := r.readByte()
	if err != nil {
		return false, err
	}
	switch b {
	case True:
		return true, nil
	case False:
		return false, nil
	default:
		return false, fmt.Errorf("expected bool but got 0x%02x", b)
	}
}

// ReadInt reads an integer value
func (r *Reader) ReadInt() (int64, error) {
	b, err := r.readByte()
	if err != nil {
		return 0, err
	}

	// Positive fixint (0x00-0x7f)
	if b <= PositiveFixintMax {
		return int64(b), nil
	}

	// Negative fixint (0xe0-0xff)
	if b >= NegativeFixintMin {
		return int64(int8(b)), nil
	}

	switch b {
	case Int8:
		v, err := r.readInt8()
		return int64(v), err
	case Int16:
		v, err := r.readInt16()
		return int64(v), err
	case Int32:
		v, err := r.readInt32()
		return int64(v), err
	case Int64:
		return r.readInt64()
	case Uint8:
		v, err := r.readUint8()
		return int64(v), err
	case Uint16:
		v, err := r.readUint16()
		return int64(v), err
	case Uint32:
		v, err := r.readUint32()
		return int64(v), err
	case Uint64:
		v, err := r.readUint64()
		if v > math.MaxInt64 {
			return 0, fmt.Errorf("uint64 value %d too large for int64", v)
		}
		return int64(v), err
	default:
		return 0, fmt.Errorf("expected int but got 0x%02x", b)
	}
}

// ReadUint reads an unsigned integer value
func (r *Reader) ReadUint() (uint64, error) {
	b, err := r.readByte()
	if err != nil {
		return 0, err
	}

	// Positive fixint (0x00-0x7f)
	if b <= PositiveFixintMax {
		return uint64(b), nil
	}

	switch b {
	case Uint8:
		v, err := r.readUint8()
		return uint64(v), err
	case Uint16:
		v, err := r.readUint16()
		return uint64(v), err
	case Uint32:
		v, err := r.readUint32()
		return uint64(v), err
	case Uint64:
		return r.readUint64()
	case Int8:
		v, err := r.readInt8()
		if v < 0 {
			return 0, fmt.Errorf("negative int8 %d cannot be converted to uint", v)
		}
		return uint64(v), err
	case Int16:
		v, err := r.readInt16()
		if v < 0 {
			return 0, fmt.Errorf("negative int16 %d cannot be converted to uint", v)
		}
		return uint64(v), err
	case Int32:
		v, err := r.readInt32()
		if v < 0 {
			return 0, fmt.Errorf("negative int32 %d cannot be converted to uint", v)
		}
		return uint64(v), err
	case Int64:
		v, err := r.readInt64()
		if v < 0 {
			return 0, fmt.Errorf("negative int64 %d cannot be converted to uint", v)
		}
		return uint64(v), err
	default:
		return 0, fmt.Errorf("expected uint but got 0x%02x", b)
	}
}

// ReadFloat32 reads a float32 value
func (r *Reader) ReadFloat32() (float32, error) {
	b, err := r.readByte()
	if err != nil {
		return 0, err
	}
	if b != Float32 {
		return 0, fmt.Errorf("expected float32 (0x%02x) but got 0x%02x", Float32, b)
	}
	return r.readFloat32()
}

// ReadFloat64 reads a float64 value
func (r *Reader) ReadFloat64() (float64, error) {
	b, err := r.readByte()
	if err != nil {
		return 0, err
	}
	if b != Float64 {
		return 0, fmt.Errorf("expected float64 (0x%02x) but got 0x%02x", Float64, b)
	}
	return r.readFloat64()
}

// ReadString reads a string value
func (r *Reader) ReadString() (string, error) {
	data, err := r.readStringBytes()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (r *Reader) readStringBytes() ([]byte, error) {
	b, err := r.readByte()
	if err != nil {
		return nil, err
	}

	var length uint32

	// fixstr (0xa0-0xbf)
	if (b & FixstrMask) == FixstrMask {
		length = uint32(b & FixstrMax)
	} else {
		switch b {
		case Str8:
			l, err := r.readUint8()
			if err != nil {
				return nil, err
			}
			length = uint32(l)
		case Str16:
			l, err := r.readUint16()
			if err != nil {
				return nil, err
			}
			length = uint32(l)
		case Str32:
			l, err := r.readUint32()
			if err != nil {
				return nil, err
			}
			length = l
		default:
			return nil, fmt.Errorf("expected string but got 0x%02x", b)
		}
	}

	return r.readBytes(int(length))
}

// ReadBinary reads binary data
func (r *Reader) ReadBinary() ([]byte, error) {
	b, err := r.readByte()
	if err != nil {
		return nil, err
	}

	var length uint32

	switch b {
	case Bin8:
		l, err := r.readUint8()
		if err != nil {
			return nil, err
		}
		length = uint32(l)
	case Bin16:
		l, err := r.readUint16()
		if err != nil {
			return nil, err
		}
		length = uint32(l)
	case Bin32:
		l, err := r.readUint32()
		if err != nil {
			return nil, err
		}
		length = l
	default:
		return nil, fmt.Errorf("expected binary but got 0x%02x", b)
	}

	final := make([]byte, length)
	copy(final, r.input[r.offset:r.offset+int(length)])
	r.offset += int(length)

	return final, nil
}

// ReadArrayHeader reads an array header and returns the number of elements
func (r *Reader) ReadArrayHeader() (int, error) {
	b, err := r.readByte()
	if err != nil {
		return 0, err
	}

	// fixarray (0x90-0x9f)
	if b >= FixarrayMask && b <= FixarrayEnd {
		return int(b & FixarrayMax), nil
	}

	switch b {
	case Array16:
		length, err := r.readUint16()
		return int(length), err
	case Array32:
		length, err := r.readUint32()
		return int(length), err
	default:
		return 0, fmt.Errorf("expected array but got 0x%02x", b)
	}
}

// ReadMapHeader reads a map header and returns the number of key-value pairs
func (r *Reader) ReadMapHeader() (int, error) {
	b, err := r.readByte()
	if err != nil {
		return 0, err
	}

	// fixmap (0x80-0x8f)
	if b >= FixmapMask && b <= FixmapEnd {
		return int(b & FixmapMax), nil
	}

	switch b {
	case Map16:
		length, err := r.readUint16()
		return int(length), err
	case Map32:
		length, err := r.readUint32()
		return int(length), err
	default:
		return 0, fmt.Errorf("expected map but got 0x%02x", b)
	}
}

// ReadArray reads an entire array into a slice
func (r *Reader) ReadArray() ([]any, error) {
	length, err := r.ReadArrayHeader()
	if err != nil {
		return nil, err
	}

	arr := make([]any, length)
	for i := 0; i < length; i++ {
		val, err := r.ReadValue()
		if err != nil {
			return nil, err
		}
		arr[i] = val
	}
	return arr, nil
}

// ReadTypedSlice reads an array of a specific type
func ReadTypedSlice[T any](r *Reader, s *[]T) error {
	length, err := r.ReadArrayHeader()
	if err != nil {
		return err
	}

	for i := 0; i < length; i++ {
		val, err := r.ReadValue()
		if err != nil {
			return err
		}

		typedVal, ok := val.(T)
		if !ok {
			return fmt.Errorf("expected %T but got %T", reflect.TypeOf(s).Elem(), reflect.TypeOf(val))
		}

		*s = append(*s, typedVal)
	}
	return nil
}

// ReadMap reads an entire map
func (r *Reader) ReadMap() (map[string]any, error) {
	length, err := r.ReadMapHeader()
	if err != nil {
		return nil, err
	}

	m := make(map[string]any, length)
	for i := 0; i < length; i++ {
		// Read key
		key, err := r.ReadValue()
		if err != nil {
			return nil, err
		}
		keyStr, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("map key must be string, got %T", key)
		}

		// Read value
		val, err := r.ReadValue()
		if err != nil {
			return nil, err
		}

		m[keyStr] = val
	}
	return m, nil
}

// ReadValue reads any value and returns it as interface{}
func (r *Reader) ReadValue() (any, error) {
	// Peek at the first byte to determine the type
	b := r.input[r.offset]

	// Determine type based on the first byte
	switch {
	case b == Nil:
		err := r.ReadNil()
		return nil, err
	case b == True || b == False:
		return r.ReadBool()
	case b <= PositiveFixintMax || b >= NegativeFixintMin ||
		b == Int8 || b == Int16 || b == Int32 || b == Int64 ||
		b == Uint8 || b == Uint16 || b == Uint32 || b == Uint64:
		return r.ReadInt()
	case b == Float32:
		return r.ReadFloat32()
	case b == Float64:
		return r.ReadFloat64()
	case (b&FixstrMask) == FixstrMask || b == Str8 || b == Str16 || b == Str32:
		return r.ReadString()
	case b == Bin8 || b == Bin16 || b == Bin32:
		return r.ReadBinary()
	case (b >= FixarrayMask && b <= FixarrayEnd) || b == Array16 || b == Array32:
		return r.ReadArray()
	case (b >= FixmapMask && b <= FixmapEnd) || b == Map16 || b == Map32:
		return r.ReadMap()
	default:
		return nil, fmt.Errorf("unknown type marker: 0x%02x", b)
	}
}

// Decode decodes msgpack data into the provided value
func (r *Reader) Decode(v any) error {
	if v == nil {
		return fmt.Errorf("cannot decode into nil")
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("decode target must be a pointer")
	}

	if rv.IsNil() {
		return fmt.Errorf("decode target cannot be a nil pointer")
	}

	return r.decodeValue(rv.Elem())
}

func (r *Reader) decodeValue(rv reflect.Value) error {
	// Check if the target implements Unmarshaler
	if rv.CanAddr() {
		if unmarshaler, ok := rv.Addr().Interface().(Unmarshaler); ok {
			return unmarshaler.UnmarshalMsgpack(r)
		}
	} else {
		// If not addressable, try via temporary pointer if pointer type implements Unmarshaler
		ptr := reflect.New(rv.Type())
		ptr.Elem().Set(rv)
		if unmarshaler, ok := ptr.Interface().(Unmarshaler); ok {
			if err := unmarshaler.UnmarshalMsgpack(r); err != nil {
				return err
			}
			rv.Set(ptr.Elem())
			return nil
		}
	}
	// Also allow non-pointer receiver implementations
	if unmarshaler, ok := rv.Interface().(Unmarshaler); ok {
		return unmarshaler.UnmarshalMsgpack(r)
	}

	// Peek at the type marker
	b := r.input[r.offset]

	// Handle nil case
	if b == Nil {
		if err := r.ReadNil(); err != nil {
			return err
		}
		if rv.Kind() == reflect.Ptr {
			rv.Set(reflect.Zero(rv.Type()))
		}
		return nil
	}

	switch rv.Kind() {
	case reflect.Bool:
		val, err := r.ReadBool()
		if err != nil {
			return err
		}
		rv.SetBool(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := r.ReadInt()
		if err != nil {
			return err
		}
		if rv.OverflowInt(val) {
			return fmt.Errorf("value %d overflows %s", val, rv.Type())
		}
		rv.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := r.ReadUint()
		if err != nil {
			return err
		}
		if rv.OverflowUint(val) {
			return fmt.Errorf("value %d overflows %s", val, rv.Type())
		}
		rv.SetUint(val)
	case reflect.Float32, reflect.Float64:
		// Try reading as float first
		switch b {
		case Float32:
			val, err := r.ReadFloat32()
			if err != nil {
				return err
			}
			rv.SetFloat(float64(val))
		case Float64:
			val, err := r.ReadFloat64()
			if err != nil {
				return err
			}
			rv.SetFloat(val)
		default:
			// Fallback to integer and convert
			val, err := r.ReadInt()
			if err != nil {
				return err
			}
			rv.SetFloat(float64(val))
		}
	case reflect.String:
		val, err := r.ReadString()
		if err != nil {
			return err
		}
		rv.SetString(val)
	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			// Handle []byte specially
			val, err := r.ReadBinary()
			if err != nil {
				return err
			}
			rv.SetBytes(val)
		} else {
			return r.decodeSlice(rv)
		}
	case reflect.Array:
		return r.decodeArray(rv)
	case reflect.Map:
		return r.decodeMap(rv)
	case reflect.Ptr:
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		return r.decodeValue(rv.Elem())
	case reflect.Interface:
		val, err := r.ReadValue()
		if err != nil {
			return err
		}
		rv.Set(reflect.ValueOf(val))
	case reflect.Struct:
		return r.decodeStruct(rv)
	default:
		return fmt.Errorf("unsupported type: %s", rv.Type())
	}

	return nil
}

func (r *Reader) decodeSlice(rv reflect.Value) error {
	length, err := r.ReadArrayHeader()
	if err != nil {
		return err
	}

	slice := reflect.MakeSlice(rv.Type(), length, length)

	for i := 0; i < length; i++ {
		if err := r.decodeValue(slice.Index(i)); err != nil {
			return err
		}
	}
	rv.Set(slice)
	return nil
}

func (r *Reader) decodeArray(rv reflect.Value) error {
	length, err := r.ReadArrayHeader()
	if err != nil {
		return err
	}

	arrayLen := rv.Len()
	if length != arrayLen {
		return fmt.Errorf("array length mismatch: expected %d, got %d", arrayLen, length)
	}

	for i := 0; i < length; i++ {
		if err := r.decodeValue(rv.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

func (r *Reader) decodeMap(rv reflect.Value) error {
	length, err := r.ReadMapHeader()
	if err != nil {
		return err
	}

	if rv.IsNil() {
		rv.Set(reflect.MakeMap(rv.Type()))
	}

	keyType := rv.Type().Key()
	valueType := rv.Type().Elem()

	for i := 0; i < length; i++ {
		// Read key
		key := reflect.New(keyType).Elem()
		if err := r.decodeValue(key); err != nil {
			return err
		}

		// Read value
		value := reflect.New(valueType).Elem()
		if err := r.decodeValue(value); err != nil {
			return err
		}

		rv.SetMapIndex(key, value)
	}
	return nil
}

func (r *Reader) decodeStruct(rv reflect.Value) error {
	// Read as map first
	length, err := r.ReadMapHeader()
	if err != nil {
		return err
	}

	rt := rv.Type()
	fields := make(map[string]reflect.Value)
	fieldMetas := getStructMeta(rt)

	// Build field map
	for i, field := range fieldMetas {
		if !field.Exported {
			continue
		}

		fields[field.Name] = rv.Field(i)
	}

	// Read map entries
	for i := 0; i < length; i++ {
		// Read key
		key, err := r.ReadString()
		if err != nil {
			return err
		}

		// Find corresponding field
		if field, exists := fields[key]; exists {
			if err := r.decodeValue(field); err != nil {
				return err
			}
		} else {
			// Skip unknown field
			if _, err := r.ReadValue(); err != nil {
				return err
			}
		}
	}

	return nil
}

// Unmarshal is a convenience function that deserializes msgpack bytes into a value
func Unmarshal[T any](data []byte, v *T) error {
	r := NewReader(data)
	return r.Decode(v)
}
