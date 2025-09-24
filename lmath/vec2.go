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
	l, err := r.ReadArrayHeader()
	if err != nil {
		return err
	}

	if l != 2 {
		return fmt.Errorf("expected 2 elements, got %d", l)
	}

	v.X, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	v.Y, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	return nil
}
