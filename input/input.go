package input

import (
	"github.com/oriolus-software/script-go/internal/ffi"
	"github.com/oriolus-software/script-go/lmath"
)

func MouseDelta() lmath.Vec2 {
	return ffi.Deserialize[lmath.Vec2](mouse_delta())
}

//go:wasm-module input
//export mouse_delta
func mouse_delta() uint64
