package lmath

import (
	"fmt"

	"github.com/oriolus-software/script-go/internal/msgpack"
)

type Vec2 struct {
	X float32 `msgpack:"x"`
	Y float32 `msgpack:"y"`
}

func (v *Vec2) UnmarshalMsgpack(r *msgpack.Reader) error {
	// Read the next value generically to support both array [x, y]
	// and map {"x": x, "y": y} encodings.
	val, err := r.ReadValue()
	if err != nil {
		return err
	}

	if val == nil {
		*v = Vec2{}
		return nil
	}

	switch t := val.(type) {
	case []any:
		if len(t) != 2 {
			return fmt.Errorf("expected 2 elements, got %d", len(t))
		}
		x, ok := asFloat32(t[0])
		if !ok {
			return fmt.Errorf("expected float for x, got %T", t[0])
		}
		y, ok := asFloat32(t[1])
		if !ok {
			return fmt.Errorf("expected float for y, got %T", t[1])
		}
		*v = Vec2{X: x, Y: y}
		return nil
	case map[string]any:
		xv, xok := t["x"]
		yv, yok := t["y"]
		if !xok || !yok {
			return fmt.Errorf("missing x or y field in map")
		}
		x, ok := asFloat32(xv)
		if !ok {
			return fmt.Errorf("expected float for x, got %T", xv)
		}
		y, ok := asFloat32(yv)
		if !ok {
			return fmt.Errorf("expected float for y, got %T", yv)
		}
		*v = Vec2{X: x, Y: y}
		return nil
	default:
		return fmt.Errorf("expected array [x,y] or map {x,y}, got %T", val)
	}
}

func asFloat32(v any) (float32, bool) {
	switch n := v.(type) {
	case float32:
		return n, true
	case float64:
		return float32(n), true
	case int64:
		return float32(n), true
	case uint64:
		return float32(n), true
	case int:
		return float32(n), true
	case uint:
		return float32(n), true
	default:
		return 0, false
	}
}
