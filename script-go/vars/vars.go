package vars

import (
	"runtime"

	"github.com/oriolus-software/script-go/ffi"
)

//go:wasm-module var
//export get_i64
func get_i64(name uint64) int64

func GetI64(name string) int64 {
	o := ffi.Serialize(name)
	v := get_i64(o.ToPacked())
	runtime.KeepAlive(o)
	return v
}

//go:wasm-module var
//export set_i64
func set_i64(name uint64, value int64)

func SetI64(name string, value int64) {
	o := ffi.Serialize(name)
	set_i64(o.ToPacked(), value)
	runtime.KeepAlive(o)
}

//go:wasm-module var
//export get_string
func get_string(name uint64) string

func GetString(name string) string {
	o := ffi.Serialize(name)
	s := get_string(o.ToPacked())
	runtime.KeepAlive(o)
	return s
}
