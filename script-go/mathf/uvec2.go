package mathf

import "github.com/oriolus-software/script-go/internal/msgpack"

type UVec2 struct {
	X uint
	Y uint
}

func (v *UVec2) UnmarshalMsgpack(r *msgpack.Reader) error {
	var arr []uint
	err := msgpack.ReadTypedSlice(r, &arr)
	if err != nil {
		return err
	}
	v.X = arr[0]
	v.Y = arr[1]
	return nil
}
