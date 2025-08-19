package msgpack_test

import (
	"bytes"
	"testing"

	"github.com/oriolus-software/script-go/internal/msgpack"
)

func TestRoundTripNil(t *testing.T) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	err := w.WriteNil()
	if err != nil {
		t.Fatal(err)
	}

	r := msgpack.NewReader(buf)
	err = r.ReadNil()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRoundTripBool(t *testing.T) {
	tests := []bool{true, false}

	for _, test := range tests {
		buf := &bytes.Buffer{}
		w := msgpack.NewWriter(buf)

		err := w.WriteBool(test)
		if err != nil {
			t.Fatal(err)
		}

		r := msgpack.NewReader(buf)
		result, err := r.ReadBool()
		if err != nil {
			t.Fatal(err)
		}

		if result != test {
			t.Fatalf("expected %v, got %v", test, result)
		}
	}
}

func TestRoundTripInt(t *testing.T) {
	tests := []int64{
		0, 1, -1, 127, -128,
		128, -129, 32767, -32768,
		32768, -32769, 2147483647, -2147483648,
		2147483648, -2147483649, 9223372036854775807, -9223372036854775808,
	}

	for _, test := range tests {
		buf := &bytes.Buffer{}
		w := msgpack.NewWriter(buf)

		err := w.WriteInt(test)
		if err != nil {
			t.Fatal(err)
		}

		r := msgpack.NewReader(buf)
		result, err := r.ReadInt()
		if err != nil {
			t.Fatal(err)
		}

		if result != test {
			t.Fatalf("expected %v, got %v", test, result)
		}
	}
}

func TestRoundTripUint(t *testing.T) {
	tests := []uint64{
		0, 1, 127, 128, 255, 256, 65535, 65536, 4294967295, 4294967296, 18446744073709551615,
	}

	for _, test := range tests {
		buf := &bytes.Buffer{}
		w := msgpack.NewWriter(buf)

		err := w.WriteUint(test)
		if err != nil {
			t.Fatal(err)
		}

		r := msgpack.NewReader(buf)
		result, err := r.ReadUint()
		if err != nil {
			t.Fatal(err)
		}

		if result != test {
			t.Fatalf("expected %v, got %v", test, result)
		}
	}
}

func TestRoundTripFloat32(t *testing.T) {
	tests := []float32{0.0, 1.0, -1.0, 3.14159, -3.14159, 1.23456e10, -1.23456e-10}

	for _, test := range tests {
		buf := &bytes.Buffer{}
		w := msgpack.NewWriter(buf)

		err := w.WriteFloat32(test)
		if err != nil {
			t.Fatal(err)
		}

		r := msgpack.NewReader(buf)
		result, err := r.ReadFloat32()
		if err != nil {
			t.Fatal(err)
		}

		if result != test {
			t.Fatalf("expected %v, got %v", test, result)
		}
	}
}

func TestRoundTripFloat64(t *testing.T) {
	tests := []float64{0.0, 1.0, -1.0, 3.141592653589793, -3.141592653589793, 1.234567890123456e100, -1.234567890123456e-100}

	for _, test := range tests {
		buf := &bytes.Buffer{}
		w := msgpack.NewWriter(buf)

		err := w.WriteFloat64(test)
		if err != nil {
			t.Fatal(err)
		}

		r := msgpack.NewReader(buf)
		result, err := r.ReadFloat64()
		if err != nil {
			t.Fatal(err)
		}

		if result != test {
			t.Fatalf("expected %v, got %v", test, result)
		}
	}
}

func TestRoundTripString(t *testing.T) {
	tests := []string{
		"", "hello", "world",
		"a", // 1 byte fixstr
		"hello world this is a longer string that exceeds fixstr", // longer string
		"こんにちは世界",                                                 // unicode
	}

	for _, test := range tests {
		buf := &bytes.Buffer{}
		w := msgpack.NewWriter(buf)

		err := w.WriteString(test)
		if err != nil {
			t.Fatal(err)
		}

		r := msgpack.NewReader(buf)
		result, err := r.ReadString()
		if err != nil {
			t.Fatal(err)
		}

		if result != test {
			t.Fatalf("expected %q, got %q", test, result)
		}
	}
}

func TestRoundTripBinary(t *testing.T) {
	tests := [][]byte{
		{},
		{1, 2, 3},
		{0, 255, 128},
		make([]byte, 300), // larger binary data
	}

	for i, test := range tests {
		buf := &bytes.Buffer{}
		w := msgpack.NewWriter(buf)

		err := w.WriteBinary(test)
		if err != nil {
			t.Fatal(err)
		}

		r := msgpack.NewReader(buf)
		result, err := r.ReadBinary()
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(result, test) {
			t.Fatalf("test %d: expected %v, got %v", i, test, result)
		}
	}
}

func TestRoundTripArray(t *testing.T) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	// Write array: [1, "hello", true]
	err := w.WriteArrayHeader(3)
	if err != nil {
		t.Fatal(err)
	}

	err = w.WriteInt(1)
	if err != nil {
		t.Fatal(err)
	}

	err = w.WriteString("hello")
	if err != nil {
		t.Fatal(err)
	}

	err = w.WriteBool(true)
	if err != nil {
		t.Fatal(err)
	}

	r := msgpack.NewReader(buf)
	result, err := r.ReadArray()
	if err != nil {
		t.Fatal(err)
	}

	expected := []any{int64(1), "hello", true}
	if len(result) != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), len(result))
	}

	for i, exp := range expected {
		if result[i] != exp {
			t.Fatalf("index %d: expected %v, got %v", i, exp, result[i])
		}
	}
}

func TestRoundTripMap(t *testing.T) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	// Write map: {"name": "John", "age": 30}
	err := w.WriteMapHeader(2)
	if err != nil {
		t.Fatal(err)
	}

	err = w.WriteString("name")
	if err != nil {
		t.Fatal(err)
	}
	err = w.WriteString("John")
	if err != nil {
		t.Fatal(err)
	}

	err = w.WriteString("age")
	if err != nil {
		t.Fatal(err)
	}
	err = w.WriteInt(30)
	if err != nil {
		t.Fatal(err)
	}

	r := msgpack.NewReader(buf)
	result, err := r.ReadMap()
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]any{"name": "John", "age": int64(30)}
	if len(result) != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), len(result))
	}

	for k, exp := range expected {
		if result[k] != exp {
			t.Fatalf("key %s: expected %v, got %v", k, exp, result[k])
		}
	}
}

func TestDecodeIntoStruct(t *testing.T) {
	type TestStruct struct {
		Name    string `msgpack:"name"`
		Age     int    `msgpack:"age"`
		IsAdult bool
	}

	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	original := TestStruct{
		Name:    "John",
		Age:     30,
		IsAdult: true,
	}

	err := w.WriteStruct(original)
	if err != nil {
		t.Fatal(err)
	}

	r := msgpack.NewReader(buf)
	var result TestStruct
	err = r.Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	if result.Name != original.Name {
		t.Fatalf("Name: expected %v, got %v", original.Name, result.Name)
	}
	if result.Age != original.Age {
		t.Fatalf("Age: expected %v, got %v", original.Age, result.Age)
	}
	if result.IsAdult != original.IsAdult {
		t.Fatalf("IsAdult: expected %v, got %v", original.IsAdult, result.IsAdult)
	}
}

func TestDecodeIntoSlice(t *testing.T) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	original := []int{1, 2, 3, 4, 5}

	err := w.WriteArrayHeader(len(original))
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range original {
		err = w.WriteInt(int64(v))
		if err != nil {
			t.Fatal(err)
		}
	}

	r := msgpack.NewReader(buf)
	var result []int
	err = r.Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(original) {
		t.Fatalf("expected length %d, got %d", len(original), len(result))
	}

	for i, exp := range original {
		if result[i] != exp {
			t.Fatalf("index %d: expected %v, got %v", i, exp, result[i])
		}
	}
}

func TestDecodeIntoMap(t *testing.T) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	original := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	err := w.WriteMapHeader(len(original))
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range original {
		err = w.WriteString(k)
		if err != nil {
			t.Fatal(err)
		}
		err = w.WriteInt(int64(v))
		if err != nil {
			t.Fatal(err)
		}
	}

	r := msgpack.NewReader(buf)
	var result map[string]int
	err = r.Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(original) {
		t.Fatalf("expected length %d, got %d", len(original), len(result))
	}

	for k, exp := range original {
		if result[k] != exp {
			t.Fatalf("key %s: expected %v, got %v", k, exp, result[k])
		}
	}
}

type testUnmarshalerStruct struct {
	Foo string
	Bar int
}

func (t *testUnmarshalerStruct) UnmarshalMsgpack(r *msgpack.Reader) error {
	length, err := r.ReadMapHeader()
	if err != nil {
		return err
	}

	for i := 0; i < length; i++ {
		key, err := r.ReadString()
		if err != nil {
			return err
		}

		switch key {
		case "foo":
			t.Foo, err = r.ReadString()
			if err != nil {
				return err
			}
		case "bar":
			val, err := r.ReadInt()
			if err != nil {
				return err
			}
			t.Bar = int(val)
		default:
			// Skip unknown field
			_, err = r.ReadValue()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func TestUnmarshaler(t *testing.T) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	// Write map: {"foo": "hello", "bar": 42}
	err := w.WriteMapHeader(2)
	if err != nil {
		t.Fatal(err)
	}

	err = w.WriteString("foo")
	if err != nil {
		t.Fatal(err)
	}
	err = w.WriteString("hello")
	if err != nil {
		t.Fatal(err)
	}

	err = w.WriteString("bar")
	if err != nil {
		t.Fatal(err)
	}
	err = w.WriteInt(42)
	if err != nil {
		t.Fatal(err)
	}

	r := msgpack.NewReader(buf)
	var result testUnmarshalerStruct
	err = r.Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	if result.Foo != "hello" {
		t.Fatalf("Foo: expected %v, got %v", "hello", result.Foo)
	}
	if result.Bar != 42 {
		t.Fatalf("Bar: expected %v, got %v", 42, result.Bar)
	}
}

func TestSerializeDeserialize(t *testing.T) {
	type TestStruct struct {
		Name   string            `msgpack:"name"`
		Age    int               `msgpack:"age"`
		Values []int             `msgpack:"values"`
		Nested map[string]string `msgpack:"nested"`
	}

	original := TestStruct{
		Name:   "Alice",
		Age:    25,
		Values: []int{1, 2, 3},
		Nested: map[string]string{"key": "value"},
	}

	// Serialize
	data, err := msgpack.Serialize(original)
	if err != nil {
		t.Fatal(err)
	}

	// Deserialize
	var result TestStruct
	err = msgpack.Deserialize(data, &result)
	if err != nil {
		t.Fatal(err)
	}

	// Compare
	if result.Name != original.Name {
		t.Fatalf("Name: expected %v, got %v", original.Name, result.Name)
	}
	if result.Age != original.Age {
		t.Fatalf("Age: expected %v, got %v", original.Age, result.Age)
	}
	if len(result.Values) != len(original.Values) {
		t.Fatalf("Values length: expected %d, got %d", len(original.Values), len(result.Values))
	}
	for i, v := range original.Values {
		if result.Values[i] != v {
			t.Fatalf("Values[%d]: expected %v, got %v", i, v, result.Values[i])
		}
	}
	if len(result.Nested) != len(original.Nested) {
		t.Fatalf("Nested length: expected %d, got %d", len(original.Nested), len(result.Nested))
	}
	for k, v := range original.Nested {
		if result.Nested[k] != v {
			t.Fatalf("Nested[%s]: expected %v, got %v", k, v, result.Nested[k])
		}
	}
}

// Benchmarks

func BenchmarkReadNil(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)
	w.WriteNil()
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := msgpack.NewReader(bytes.NewReader(data))
		if err := r.ReadNil(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadBool(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)
	w.WriteBool(true)
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := msgpack.NewReader(bytes.NewReader(data))
		if _, err := r.ReadBool(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadInteger(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)
	w.WriteInt(42)
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := msgpack.NewReader(bytes.NewReader(data))
		if _, err := r.ReadInt(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadUint(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)
	w.WriteUint(42)
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := msgpack.NewReader(bytes.NewReader(data))
		if _, err := r.ReadUint(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadFloat32(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)
	w.WriteFloat32(3.14159)
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := msgpack.NewReader(bytes.NewReader(data))
		if _, err := r.ReadFloat32(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadFloat64(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)
	w.WriteFloat64(3.141592653589793)
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := msgpack.NewReader(bytes.NewReader(data))
		if _, err := r.ReadFloat64(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadString(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)
	w.WriteString("my_variable_name")
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := msgpack.NewReader(bytes.NewReader(data))
		if _, err := r.ReadString(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadBinary(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)
	testData := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	w.WriteBinary(testData)
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := msgpack.NewReader(bytes.NewReader(data))
		if _, err := r.ReadBinary(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadArray(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)
	w.WriteArrayHeader(3)
	w.WriteInt(1)
	w.WriteString("hello")
	w.WriteBool(true)
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := msgpack.NewReader(bytes.NewReader(data))
		if _, err := r.ReadArray(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadMap(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)
	w.WriteMapHeader(2)
	w.WriteString("name")
	w.WriteString("John")
	w.WriteString("age")
	w.WriteInt(30)
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := msgpack.NewReader(bytes.NewReader(data))
		if _, err := r.ReadMap(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadValue(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)
	w.WriteMapHeader(2)
	w.WriteString("name")
	w.WriteString("John")
	w.WriteString("age")
	w.WriteInt(30)
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := msgpack.NewReader(bytes.NewReader(data))
		if _, err := r.ReadValue(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeStruct(b *testing.B) {
	type TestStruct struct {
		Name string `msgpack:"name"`
		Age  int    `msgpack:"age"`
	}

	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)
	original := TestStruct{Name: "John", Age: 30}
	w.WriteStruct(original)
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := msgpack.NewReader(bytes.NewReader(data))
		var result TestStruct
		if err := r.Decode(&result); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeSlice(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)
	original := []int{1, 2, 3, 4, 5}
	w.WriteArrayHeader(len(original))
	for _, v := range original {
		w.WriteInt(int64(v))
	}
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := msgpack.NewReader(bytes.NewReader(data))
		var result []int
		if err := r.Decode(&result); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeMap(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)
	original := map[string]int{"one": 1, "two": 2, "three": 3}
	w.WriteMapHeader(len(original))
	for k, v := range original {
		w.WriteString(k)
		w.WriteInt(int64(v))
	}
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := msgpack.NewReader(bytes.NewReader(data))
		var result map[string]int
		if err := r.Decode(&result); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDeserialize(b *testing.B) {
	type TestStruct struct {
		Name   string `msgpack:"name"`
		Age    int    `msgpack:"age"`
		Values []int  `msgpack:"values"`
	}

	original := TestStruct{
		Name:   "Alice",
		Age:    25,
		Values: []int{1, 2, 3},
	}

	data, err := msgpack.Serialize(original)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result TestStruct
		if err := msgpack.Deserialize(data, &result); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshaler(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)
	w.WriteMapHeader(2)
	w.WriteString("foo")
	w.WriteString("hello")
	w.WriteString("bar")
	w.WriteInt(42)
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := msgpack.NewReader(bytes.NewReader(data))
		var result testUnmarshalerStruct
		if err := r.Decode(&result); err != nil {
			b.Fatal(err)
		}
	}
}
