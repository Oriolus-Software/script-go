package msgpack_test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/oriolus-software/script-go/msgpack"
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

func structuralEquals(a, b any) bool {
	fmt.Println("structuralEquals", reflect.TypeOf(a), a, reflect.TypeOf(b), b)
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
