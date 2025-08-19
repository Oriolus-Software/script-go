package msgpack_test

import (
	"bytes"
	"testing"

	"github.com/oriolus-software/script-go/internal/msgpack"
	orig "github.com/vmihailenco/msgpack/v5"
)

func FuzzInteger(f *testing.F) {
	f.Add(1)
	f.Add(1)
	f.Add(2)
	f.Add(3)
	f.Add(4)
	f.Add(5)
	f.Add(6)

	buf := bytes.NewBuffer(nil)

	f.Fuzz(func(t *testing.T, i int) {
		buf.Reset()

		self := msgpack.NewWriter(buf)
		self.WriteInt(int64(i))

		ref, err := orig.Marshal(i)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(buf.Bytes(), ref) {
			t.Fatalf("serializing %d not equal: %v != %v", i, buf.Bytes(), ref)
		}
	})
}

func BenchmarkInteger(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	self := msgpack.NewWriter(buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := self.WriteInt(int64(i)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBool(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	self := msgpack.NewWriter(buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := self.WriteBool(i%2 == 0); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkString(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	self := msgpack.NewWriter(buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := self.WriteString("my_variable_name"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFloat(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	self := msgpack.NewWriter(buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := self.WriteFloat64(float64(i)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMap(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	self := msgpack.NewWriter(buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := self.WriteMap(map[string]any{"foo": "bar"}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStruct(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	self := msgpack.NewWriter(buf)

	obj := struct {
		Foo string
		Bar int
	}{
		Foo: "foo",
		Bar: 1,
	}

	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := self.WriteStruct(obj); err != nil {
			b.Fatal(err)
		}
	}
}

type testMarshalerStruct struct {
	Foo string `msgpack:"foo"`
	Bar int    `msgpack:"bar"`
}

func (t *testMarshalerStruct) MarshalMsgpack(w *msgpack.Writer) error {
	w.WriteMapHeader(2)
	w.WriteString("foo")
	w.WriteString(t.Foo)
	w.WriteString("bar")
	w.WriteInt(int64(t.Bar))
	return nil
}

func BenchmarkHandwrittenMarshaler(b *testing.B) {
	obj := &testMarshalerStruct{
		Foo: "foo",
		Bar: 1,
	}

	buf := bytes.NewBuffer(nil)
	self := msgpack.NewWriter(buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := self.WriteStruct(obj); err != nil {
			b.Fatal(err)
		}
	}
}

func TestHandwrittenMarshaler(t *testing.T) {
	obj := &testMarshalerStruct{
		Foo: "foo",
		Bar: 1,
	}

	buf := bytes.NewBuffer(nil)
	self := msgpack.NewWriter(buf)
	if err := self.WriteStruct(obj); err != nil {
		t.Fatal(err)
	}

	var des testMarshalerStruct
	err := orig.Unmarshal(buf.Bytes(), &des)
	if err != nil {
		t.Fatal(err)
	}

	if des.Foo != obj.Foo || des.Bar != obj.Bar {
		t.Fatalf("deserializing marshaler not equal: des=%+v != obj=%+v", des, obj)
	}
}

func FuzzString(f *testing.F) {
	f.Add("hello")
	f.Add("world")

	buf := bytes.NewBuffer(nil)

	f.Fuzz(func(t *testing.T, s string) {
		buf.Reset()

		self := msgpack.NewWriter(buf)
		self.WriteString(s)

		ref, err := orig.Marshal(s)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(buf.Bytes(), ref) {
			t.Fatalf("serializing %q not equal: %v != %v", s, buf.Bytes(), ref)
		}
	})
}

func TestMap(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	m := map[string]any{
		"foo": int8(1),
		"bar": "baz",
	}

	self := msgpack.NewWriter(buf)
	self.WriteMap(m)

	var des map[string]any
	err := orig.Unmarshal(buf.Bytes(), &des)
	if err != nil {
		t.Fatal(err)
	}

	if !structuralEquals(des, m) {
		t.Fatalf("deserializing map not equal: %v != %v", des, m)
	}
}

func TestStruct(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	type TestStruct struct {
		Name    string `msgpack:"name"`
		Age     int    `msgpack:"age"`
		IsAdult bool
		private string // this should be ignored
	}

	obj := TestStruct{
		Name:    "John",
		Age:     30,
		IsAdult: true,
		private: "secret",
	}

	self := msgpack.NewWriter(buf)
	err := self.WriteStruct(obj)
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

func structuralEquals(a, b any) bool {
	switch a := a.(type) {
	case map[string]any:
		b, ok := b.(map[string]any)
		if !ok {
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
