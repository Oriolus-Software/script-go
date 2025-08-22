package msgpack

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
)

type Writer struct {
	w io.Writer
}

type Marshaler interface {
	MarshalMsgpack(w *Writer) error
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func (w *Writer) write(data []byte) error {
	_, err := w.w.Write(data)
	return err
}

var writeTemp = make([]byte, 8)

func (w *Writer) writeByte(b byte) error {
	writeTemp[0] = b
	return w.write(writeTemp[:1])
}

func (w *Writer) writeUint8(v uint8) error {
	return w.writeByte(v)
}

func (w *Writer) writeUint16(v uint16) error {
	binary.BigEndian.PutUint16(writeTemp[:2], v)
	return w.write(writeTemp[:2])
}

func (w *Writer) writeUint32(v uint32) error {
	binary.BigEndian.PutUint32(writeTemp[:4], v)
	return w.write(writeTemp[:4])
}

func (w *Writer) writeUint64(v uint64) error {
	binary.BigEndian.PutUint64(writeTemp[:8], v)
	return w.write(writeTemp[:8])
}

func (w *Writer) writeInt8(v int8) error {
	return w.writeByte(byte(v))
}

func (w *Writer) writeInt16(v int16) error {
	return w.writeUint16(uint16(v))
}

func (w *Writer) writeInt32(v int32) error {
	return w.writeUint32(uint32(v))
}

func (w *Writer) writeInt64(v int64) error {
	return w.writeUint64(uint64(v))
}

func (w *Writer) writeFloat32(v float32) error {
	return w.writeUint32(math.Float32bits(v))
}

func (w *Writer) writeFloat64(v float64) error {
	return w.writeUint64(math.Float64bits(v))
}

// WriteNil writes a nil value
func (w *Writer) WriteNil() error {
	return w.writeByte(Nil)
}

// WriteBool writes a boolean value
func (w *Writer) WriteBool(b bool) error {
	if b {
		return w.writeByte(True)
	}
	return w.writeByte(False)
}

// WriteInt writes an integer value
func (w *Writer) WriteInt(v int64) error {
	if v >= 0 {
		return w.WriteUint(uint64(v))
	}

	if v >= -32 {
		return w.writeByte(byte(v))
	}

	if v >= math.MinInt8 {
		if err := w.writeByte(Int8); err != nil {
			return err
		}
		return w.writeInt8(int8(v))
	}

	if v >= math.MinInt16 {
		if err := w.writeByte(Int16); err != nil {
			return err
		}
		return w.writeInt16(int16(v))
	}

	if v >= math.MinInt32 {
		if err := w.writeByte(Int32); err != nil {
			return err
		}
		return w.writeInt32(int32(v))
	}

	if err := w.writeByte(Int64); err != nil {
		return err
	}
	return w.writeInt64(v)
}

// WriteUint writes an unsigned integer value
func (w *Writer) WriteUint(v uint64) error {
	if v <= PositiveFixintMax {
		return w.writeByte(byte(v))
	}

	if v <= math.MaxUint8 {
		if err := w.writeByte(Uint8); err != nil {
			return err
		}
		return w.writeUint8(uint8(v))
	}

	if v <= math.MaxUint16 {
		if err := w.writeByte(Uint16); err != nil {
			return err
		}
		return w.writeUint16(uint16(v))
	}

	if v <= math.MaxUint32 {
		if err := w.writeByte(Uint32); err != nil {
			return err
		}
		return w.writeUint32(uint32(v))
	}

	if err := w.writeByte(Uint64); err != nil {
		return err
	}
	return w.writeUint64(v)
}

// WriteFloat32 writes a float32 value
func (w *Writer) WriteFloat32(v float32) error {
	if err := w.writeByte(Float32); err != nil {
		return err
	}
	return w.writeFloat32(v)
}

// WriteFloat64 writes a float64 value
func (w *Writer) WriteFloat64(v float64) error {
	if err := w.writeByte(Float64); err != nil {
		return err
	}
	return w.writeFloat64(v)
}

// WriteString writes a string value
func (w *Writer) WriteString(s string) error {
	data := []byte(s)
	return w.writeStringBytes(data)
}

func (w *Writer) writeStringBytes(data []byte) error {
	l := len(data)

	if l <= FixstrMax {
		if err := w.writeByte(FixstrMask | byte(l)); err != nil {
			return err
		}
	} else if l <= math.MaxUint8 {
		if err := w.writeByte(Str8); err != nil {
			return err
		}
		if err := w.writeUint8(uint8(l)); err != nil {
			return err
		}
	} else if l <= math.MaxUint16 {
		if err := w.writeByte(Str16); err != nil {
			return err
		}
		if err := w.writeUint16(uint16(l)); err != nil {
			return err
		}
	} else {
		if err := w.writeByte(Str32); err != nil {
			return err
		}
		if err := w.writeUint32(uint32(l)); err != nil {
			return err
		}
	}

	return w.write(data)
}

// WriteBinary writes binary data
func (w *Writer) WriteBinary(data []byte) error {
	l := len(data)

	if l <= math.MaxUint8 {
		if err := w.writeByte(Bin8); err != nil {
			return err
		}
		if err := w.writeUint8(uint8(l)); err != nil {
			return err
		}
	} else if l <= math.MaxUint16 {
		if err := w.writeByte(Bin16); err != nil {
			return err
		}
		if err := w.writeUint16(uint16(l)); err != nil {
			return err
		}
	} else {
		if err := w.writeByte(Bin32); err != nil {
			return err
		}
		if err := w.writeUint32(uint32(l)); err != nil {
			return err
		}
	}

	return w.write(data)
}

// WriteArrayHeader writes an array header with the given length
func (w *Writer) WriteArrayHeader(l int) error {
	if l <= FixarrayMax {
		return w.writeByte(FixarrayMask | byte(l))
	} else if l <= math.MaxUint16 {
		if err := w.writeByte(Array16); err != nil {
			return err
		}
		return w.writeUint16(uint16(l))
	} else {
		if err := w.writeByte(Array32); err != nil {
			return err
		}
		return w.writeUint32(uint32(l))
	}
}

// WriteMapHeader writes a map header with the given length
func (w *Writer) WriteMapHeader(l int) error {
	if l <= FixmapMax {
		return w.writeByte(FixmapMask | byte(l))
	} else if l <= math.MaxUint16 {
		if err := w.writeByte(Map16); err != nil {
			return err
		}
		return w.writeUint16(uint16(l))
	} else {
		if err := w.writeByte(Map32); err != nil {
			return err
		}
		return w.writeUint32(uint32(l))
	}
}

// Encode encodes any value to msgpack format
func (w *Writer) Encode(v any) error {
	if v == nil {
		return w.WriteNil()
	}

	rv := reflect.ValueOf(v)
	return w.encodeValue(rv)
}

func (w *Writer) WriteMap(m map[string]any) error {
	if err := w.WriteMapHeader(len(m)); err != nil {
		return err
	}
	for k, v := range m {
		if err := w.encodeValue(reflect.ValueOf(k)); err != nil {
			return err
		}
		if err := w.encodeValue(reflect.ValueOf(v)); err != nil {
			return err
		}
	}
	return nil
}

var structMeta = make(map[reflect.Type][]StructField)

type StructField struct {
	Exported bool
	Name     string
}

func getStructMeta(rt reflect.Type) []StructField {
	if meta, ok := structMeta[rt]; ok {
		return meta
	}

	meta := make([]StructField, 0)

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldName := field.Name
		if tag := field.Tag.Get("msgpack"); tag != "" && tag != "-" {
			fieldName = tag
		}

		meta = append(meta, StructField{
			Exported: field.IsExported(),
			Name:     fieldName,
		})
	}

	structMeta[rt] = meta
	return meta
}

// WriteStruct writes a struct as a msgpack map
func (w *Writer) WriteStruct(v any) error {
	if m, ok := v.(Marshaler); ok {
		return m.MarshalMsgpack(w)
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return w.WriteNil()
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return errors.New("not a struct")
	}

	rt := rv.Type()
	fields := getStructMeta(rt)

	if err := w.WriteMapHeader(len(fields)); err != nil {
		return err
	}

	for i, field := range fields {
		if !field.Exported {
			continue
		}

		fieldValue := rv.Field(i)

		if err := w.WriteString(field.Name); err != nil {
			return err
		}
		if err := w.encodeValue(fieldValue); err != nil {
			return err
		}
	}

	return nil
}

func (w *Writer) encodeValue(rv reflect.Value) error {
	marshaler, ok := rv.Interface().(Marshaler)
	if ok {
		return marshaler.MarshalMsgpack(w)
	}

	switch rv.Kind() {
	case reflect.Bool:
		return w.WriteBool(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return w.WriteInt(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return w.WriteUint(rv.Uint())
	case reflect.Float32:
		return w.WriteFloat32(float32(rv.Float()))
	case reflect.Float64:
		return w.WriteFloat64(rv.Float())
	case reflect.String:
		return w.WriteString(rv.String())
	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return w.WriteBinary(rv.Bytes())
		}
		return w.encodeSlice(rv)
	case reflect.Array:
		return w.encodeArray(rv)
	case reflect.Map:
		return w.encodeMap(rv)
	case reflect.Ptr:
		if rv.IsNil() {
			return w.WriteNil()
		}
		return w.encodeValue(rv.Elem())
	case reflect.Interface:
		if rv.IsNil() {
			return w.WriteNil()
		}
		return w.encodeValue(rv.Elem())
	case reflect.Struct:
		return w.WriteStruct(rv.Interface())
	default:
		return fmt.Errorf("unsupported type: %s", rv.Type())
	}
}

func (w *Writer) encodeSlice(rv reflect.Value) error {
	l := rv.Len()
	if err := w.WriteArrayHeader(l); err != nil {
		return err
	}

	for i := 0; i < l; i++ {
		if err := w.encodeValue(rv.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) encodeArray(rv reflect.Value) error {
	l := rv.Len()
	if err := w.WriteArrayHeader(l); err != nil {
		return err
	}

	for i := 0; i < l; i++ {
		if err := w.encodeValue(rv.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) encodeMap(rv reflect.Value) error {
	keys := rv.MapKeys()
	l := len(keys)
	if err := w.WriteMapHeader(l); err != nil {
		return err
	}

	for _, key := range keys {
		if err := w.encodeValue(key); err != nil {
			return err
		}
		if err := w.encodeValue(rv.MapIndex(key)); err != nil {
			return err
		}
	}
	return nil
}

// Serialize is a convenience function that encodes any value to msgpack bytes
func Serialize(val any) ([]byte, error) {
	writer := &bytes.Buffer{}
	w := NewWriter(writer)

	if err := w.Encode(val); err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}
