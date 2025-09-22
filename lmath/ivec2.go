package lmath

import "github.com/oriolus-software/script-go/internal/msgpack"

type IVec2 struct {
	X int
	Y int
}

func (v IVec2) MarshalMsgpack(w *msgpack.Writer) error {
	if err := w.WriteArrayHeader(2); err != nil {
		return err
	}

	if err := w.WriteInt(int64(v.X)); err != nil {
		return err
	}

	if err := w.WriteInt(int64(v.Y)); err != nil {
		return err
	}

	return nil
}
