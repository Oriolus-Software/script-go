package lmath

import (
	"fmt"

	"github.com/oriolus-software/script-go/internal/msgpack"
)

type UVec2 struct {
	X uint
	Y uint
}

func (v *UVec2) UnmarshalMsgpack(r *msgpack.Reader) error {
	length, err := r.ReadArrayHeader()
	if err != nil {
		return err
	}

	if length != 2 {
		return fmt.Errorf("expected 2 elements, got %d", length)
	}

	x, err := r.ReadValue()
	if err != nil {
		return err
	}
	v.X = x.(uint)
	y, err := r.ReadValue()
	if err != nil {
		return err
	}
	v.Y = y.(uint)

	return nil
}

func (v UVec2) MarshalMsgpack(w *msgpack.Writer) error {
	if err := w.WriteArrayHeader(2); err != nil {
		return err
	}

	if err := w.WriteUint(uint64(v.X)); err != nil {
		return err
	}

	if err := w.WriteUint(uint64(v.Y)); err != nil {
		return err
	}

	return nil
}
