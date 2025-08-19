package msgpack_test

import (
	"bytes"
	"testing"

	"github.com/oriolus-software/script-go/internal/msgpack"
	orig "github.com/vmihailenco/msgpack/v5"
)

// Tests

func TestWriteNil(t *testing.T) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	err := w.WriteNil()
	if err != nil {
		t.Fatal(err)
	}

	// Verify the written bytes
	expected := []byte{0xc0}
	if !bytes.Equal(buf.Bytes(), expected) {
		t.Fatalf("expected %v, got %v", expected, buf.Bytes())
	}
}

func TestWriteBool(t *testing.T) {
	tests := []struct {
		value    bool
		expected []byte
	}{
		{true, []byte{0xc3}},
		{false, []byte{0xc2}},
	}

	for _, test := range tests {
		buf := &bytes.Buffer{}
		w := msgpack.NewWriter(buf)

		err := w.WriteBool(test.value)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(buf.Bytes(), test.expected) {
			t.Fatalf("value %v: expected %v, got %v", test.value, test.expected, buf.Bytes())
		}
	}
}

func TestWriteInt(t *testing.T) {
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

		// Verify we can read it back correctly (round-trip test)
		r := msgpack.NewReader(buf)
		result, err := r.ReadInt()
		if err != nil {
			t.Fatal(err)
		}

		if result != test {
			t.Fatalf("value %d: round-trip failed, got %d", test, result)
		}
	}
}

func TestWriteUint(t *testing.T) {
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

		// Verify we can read it back correctly (round-trip test)
		r := msgpack.NewReader(buf)
		result, err := r.ReadUint()
		if err != nil {
			t.Fatal(err)
		}

		if result != test {
			t.Fatalf("value %d: round-trip failed, got %d", test, result)
		}
	}
}

func TestWriteFloat32(t *testing.T) {
	tests := []float32{0.0, 1.0, -1.0, 3.14159, -3.14159, 1.23456e10, -1.23456e-10}

	for _, test := range tests {
		buf := &bytes.Buffer{}
		w := msgpack.NewWriter(buf)

		err := w.WriteFloat32(test)
		if err != nil {
			t.Fatal(err)
		}

		// Verify compatibility with reference implementation
		ref, err := orig.Marshal(test)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(buf.Bytes(), ref) {
			t.Fatalf("value %v: expected %v, got %v", test, ref, buf.Bytes())
		}
	}
}

func TestWriteFloat64(t *testing.T) {
	tests := []float64{0.0, 1.0, -1.0, 3.141592653589793, -3.141592653589793, 1.234567890123456e100, -1.234567890123456e-100}

	for _, test := range tests {
		buf := &bytes.Buffer{}
		w := msgpack.NewWriter(buf)

		err := w.WriteFloat64(test)
		if err != nil {
			t.Fatal(err)
		}

		// Verify compatibility with reference implementation
		ref, err := orig.Marshal(test)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(buf.Bytes(), ref) {
			t.Fatalf("value %v: expected %v, got %v", test, ref, buf.Bytes())
		}
	}
}

func TestWriteString(t *testing.T) {
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

		// Verify compatibility with reference implementation
		ref, err := orig.Marshal(test)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(buf.Bytes(), ref) {
			t.Fatalf("value %q: expected %v, got %v", test, ref, buf.Bytes())
		}
	}
}

func TestWriteBinary(t *testing.T) {
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

		// For binary data, we need to verify it can be read back correctly
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

func TestWriteArray(t *testing.T) {
	tests := [][]any{
		{},
		{1, 2, 3},
		{1, "two", 3.0, []any{4, 5, 6}},
		{1, "two", map[string]any{"three": 4}},
	}

	for i := range 100 {
		arr := make([]any, 0, i)
		for j := 0; j < i; j++ {
			arr = append(arr, j)
		}
		tests = append(tests, arr)
	}

	for i, test := range tests {
		self, err := msgpack.Serialize(test)
		if err != nil {
			t.Fatal(err)
		}

		orig, err := orig.Marshal(test)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(self, orig) {
			t.Fatalf("test %d: expected %v, got %v", i, orig, self)
		}
	}
}

func TestWriteMap(t *testing.T) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	m := map[string]any{
		"foo": int8(1),
		"bar": "baz",
	}

	err := w.WriteMap(m)
	if err != nil {
		t.Fatal(err)
	}

	// Verify compatibility with reference implementation
	var des map[string]any
	err = orig.Unmarshal(buf.Bytes(), &des)
	if err != nil {
		t.Fatal(err)
	}

	if !structuralEquals(des, m) {
		t.Fatalf("deserializing map not equal: %v != %v", des, m)
	}
}

func TestWriteStruct(t *testing.T) {
	type TestStruct struct {
		Name    string `msgpack:"name"`
		Age     int    `msgpack:"age"`
		IsAdult bool
		private string // this should be ignored
	}

	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	obj := TestStruct{
		Name:    "John",
		Age:     30,
		IsAdult: true,
		private: "secret",
	}

	err := w.WriteStruct(obj)
	if err != nil {
		t.Fatal(err)
	}

	// Deserialize and check the result
	var des map[string]any
	err = orig.Unmarshal(buf.Bytes(), &des)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]any{
		"name":    "John",
		"age":     int8(30), // msgpack deserializes small ints as int8
		"IsAdult": true,
	}

	if !structuralEquals(des, expected) {
		t.Fatalf("deserializing struct not equal: %v != %v", des, expected)
	}

	// Verify private field is not included
	if _, exists := des["private"]; exists {
		t.Fatal("private field should not be serialized")
	}
}

type testWriterMarshalerStruct struct {
	Foo string `msgpack:"foo"`
	Bar int    `msgpack:"bar"`
}

func (t *testWriterMarshalerStruct) MarshalMsgpack(w *msgpack.Writer) error {
	w.WriteMapHeader(2)
	w.WriteString("foo")
	w.WriteString(t.Foo)
	w.WriteString("bar")
	w.WriteInt(int64(t.Bar))
	return nil
}

func TestWriteMarshaler(t *testing.T) {
	obj := &testWriterMarshalerStruct{
		Foo: "foo",
		Bar: 1,
	}

	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)
	err := w.WriteStruct(obj)
	if err != nil {
		t.Fatal(err)
	}

	var des testWriterMarshalerStruct
	err = orig.Unmarshal(buf.Bytes(), &des)
	if err != nil {
		t.Fatal(err)
	}

	if des.Foo != obj.Foo || des.Bar != obj.Bar {
		t.Fatalf("deserializing marshaler not equal: des=%+v != obj=%+v", des, obj)
	}
}

func TestWriteEncode(t *testing.T) {
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

	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	err := w.Encode(original)
	if err != nil {
		t.Fatal(err)
	}

	// Verify it can be read back correctly
	r := msgpack.NewReader(buf)
	var result TestStruct
	err = r.Decode(&result)
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
}

func TestWriteSerialize(t *testing.T) {
	type TestStruct struct {
		Name string `msgpack:"name"`
		Age  int    `msgpack:"age"`
	}

	original := TestStruct{Name: "Bob", Age: 42}

	data, err := msgpack.Serialize(original)
	if err != nil {
		t.Fatal(err)
	}

	var result TestStruct
	err = msgpack.Deserialize(data, &result)
	if err != nil {
		t.Fatal(err)
	}

	if result.Name != original.Name || result.Age != original.Age {
		t.Fatalf("round-trip failed: expected %+v, got %+v", original, result)
	}
}

// Fuzz tests

func FuzzWriteInteger(f *testing.F) {
	f.Add(1)
	f.Add(2)
	f.Add(3)
	f.Add(4)
	f.Add(5)
	f.Add(6)

	buf := bytes.NewBuffer(nil)

	f.Fuzz(func(t *testing.T, i int) {
		buf.Reset()

		w := msgpack.NewWriter(buf)
		w.WriteInt(int64(i))

		ref, err := orig.Marshal(i)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(buf.Bytes(), ref) {
			t.Fatalf("serializing %d not equal: %v != %v", i, buf.Bytes(), ref)
		}
	})
}

func FuzzWriteString(f *testing.F) {
	f.Add("hello")
	f.Add("world")

	buf := bytes.NewBuffer(nil)

	f.Fuzz(func(t *testing.T, s string) {
		buf.Reset()

		w := msgpack.NewWriter(buf)
		w.WriteString(s)

		ref, err := orig.Marshal(s)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(buf.Bytes(), ref) {
			t.Fatalf("serializing %q not equal: %v != %v", s, buf.Bytes(), ref)
		}
	})
}

// Benchmarks

func BenchmarkWriteNil(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := w.WriteNil(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteBool(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := w.WriteBool(i%2 == 0); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteInteger(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := w.WriteInt(int64(i)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteUint(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := w.WriteUint(uint64(i)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteFloat32(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := w.WriteFloat32(3.14159); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteFloat64(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := w.WriteFloat64(float64(i)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteString(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := w.WriteString("my_variable_name"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteBinary(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)
	testData := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := w.WriteBinary(testData); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteMap(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := w.WriteMap(map[string]any{"foo": "bar"}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteStruct(b *testing.B) {
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	obj := struct {
		Foo string
		Bar int
	}{
		Foo: "foo",
		Bar: 1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := w.WriteStruct(obj); err != nil {
			b.Fatal(err)
		}
	}
}

type Vector3 struct {
	X float32 `msgpack:"x"`
	Y float32 `msgpack:"y"`
	Z float32 `msgpack:"z"`
}

func (t *Vector3) MarshalMsgpack(w *msgpack.Writer) error {
	w.WriteArrayHeader(3)
	w.WriteFloat32(t.X)
	w.WriteFloat32(t.Y)
	w.WriteFloat32(t.Z)
	return nil
}

func BenchmarkWriteMarshalerVector3(b *testing.B) {
	obj := &Vector3{
		X: 1.0,
		Y: 2.0,
		Z: 3.0,
	}

	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := w.WriteStruct(obj); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteMarshaler(b *testing.B) {
	obj := &testWriterMarshalerStruct{
		Foo: "foo",
		Bar: 1,
	}

	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := w.WriteStruct(obj); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteEncode(b *testing.B) {
	type TestStruct struct {
		Name string `msgpack:"name"`
		Age  int    `msgpack:"age"`
	}

	obj := TestStruct{Name: "John", Age: 30}
	buf := &bytes.Buffer{}
	w := msgpack.NewWriter(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := w.Encode(obj); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteSerialize(b *testing.B) {
	type TestStruct struct {
		Name string `msgpack:"name"`
		Age  int    `msgpack:"age"`
	}

	obj := TestStruct{Name: "John", Age: 30}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := msgpack.Serialize(obj); err != nil {
			b.Fatal(err)
		}
	}
}

func structuralEquals(a, b any) bool {
	switch a := a.(type) {
	case map[string]any:
		b, ok := b.(map[string]any)
		if !ok {
			return false
		}
		if len(a) != len(b) {
			return false
		}
		for k, v := range a {
			if !structuralEquals(v, b[k]) {
				return false
			}
		}
		return true
	case string:
		b, ok := b.(string)
		if !ok {
			return false
		}
		return a == b
	default:
		return a == b
	}
}
