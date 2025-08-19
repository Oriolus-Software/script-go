package input

import (
	"github.com/oriolus-software/script-go/internal/ffi"
	"github.com/oriolus-software/script-go/mathf"
)

//go:wasm-module input
//export mouse_delta
func mouse_delta() uint64

func MouseDelta() mathf.Vec2 {
	return ffi.Deserialize[mathf.Vec2](mouse_delta())
}
