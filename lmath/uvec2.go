package lmath

import "github.com/oriolus-software/script-go/internal/msgpack"

type UVec2 struct {
	X uint
	Y uint
}

func (v *UVec2) UnmarshalMsgpack(r *msgpack.Reader) error {
	var arr []int
	err := msgpack.ReadTypedSlice(r, &arr)
	if err != nil {
		return err
	}
	v.X = uint(arr[0])
	v.Y = uint(arr[1])
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
