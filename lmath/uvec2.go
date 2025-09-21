package lmath

import (
	"github.com/oriolus-software/script-go/internal/msgpack"
)

type UVec2 struct {
	X uint
	Y uint
}

func (v UVec2) MarshalMsgpack(w *msgpack.Writer) error {
	if err := w.WriteArrayHeader(2); err != nil {
		return err
	}

	w.Encode(v.X)
	w.Encode(v.Y)

	return nil
}
