package vars

import (
	"fmt"

	"github.com/oriolus-software/script-go/assets"
	"github.com/oriolus-software/script-go/internal/ffi"
)

//go:wasm-module var
//export get_i64
func get_i64(name uint64) int64

func GetI64(name string) int64 {
	return get_i64(ffi.Serialize(name).ToPacked())
}

//go:wasm-module var
//export set_i64
func set_i64(name uint64, value int64)

func SetI64(name string, value int64) {
	set_i64(ffi.Serialize(name).ToPacked(), value)
}

//go:wasm-module var
//export get_f64
func get_f64(name uint64) float64

func GetF64(name string) float64 {
	return get_f64(ffi.Serialize(name).ToPacked())
}

//go:wasm-module var
//export set_f64
func set_f64(name uint64, value float64)

func SetF64(name string, value float64) {
	set_f64(ffi.Serialize(name).ToPacked(), value)
}

//go:wasm-module var
//export get_bool
func get_bool(name uint64) bool

func GetBool(name string) bool {
	return get_bool(ffi.Serialize(name).ToPacked())
}

//go:wasm-module var
//export set_bool
func set_bool(name uint64, value bool)

func SetBool(name string, value bool) {
	set_bool(ffi.Serialize(name).ToPacked(), value)
}

//go:wasm-module var
//export get_string
func get_string(name uint64) string

func GetString(name string) string {
	return get_string(ffi.Serialize(name).ToPacked())
}

//go:wasm-module var
//export set_string
func set_string(name uint64, value uint64)

func SetString(name string, value string) {
	set_string(ffi.Serialize(name).ToPacked(), ffi.Serialize(value).ToPacked())
}

//go:wasm-module var
//export get_content_id
func get_content_id(name uint64) uint64

func GetContentId(name string) assets.ContentId {
	return ffi.Deserialize[assets.ContentId](get_content_id(ffi.Serialize(name).ToPacked()))
}

//go:wasm-module var
//export set_content_id
func set_content_id(name uint64, value uint64)

func SetContentId(name string, value assets.ContentId) {
	set_content_id(ffi.Serialize(name).ToPacked(), ffi.Serialize(value).ToPacked())
}

func Set(name string, value any) {
	switch value := value.(type) {
	case int:
		SetI64(name, int64(value))
	case int8:
		SetI64(name, int64(value))
	case int16:
		SetI64(name, int64(value))
	case int32:
		SetI64(name, int64(value))
	case int64:
		SetI64(name, value)
	case uint:
		SetI64(name, int64(value))
	case uint8:
		SetI64(name, int64(value))
	case uint16:
		SetI64(name, int64(value))
	case uint32:
		SetI64(name, int64(value))
	case uint64:
		SetI64(name, int64(value))
	case float32:
		SetF64(name, float64(value))
	case float64:
		SetF64(name, value)
	case bool:
		SetBool(name, value)
	case string:
		SetString(name, value)
	case assets.ContentId:
		SetContentId(name, value)
	default:
		panic(fmt.Sprintf("unsupported type: %T", value))
	}
}
