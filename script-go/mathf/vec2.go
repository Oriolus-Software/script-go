package mathf

import (
	"github.com/oriolus-software/script-go/internal/msgpack"
)

type Vec2 struct {
	X float32 `msgpack:"x"`
	Y float32 `msgpack:"y"`
}

func (v *Vec2) UnmarshalMsgpack(r *msgpack.Reader) error {
	var arr []float32
	err := msgpack.ReadTypedSlice(r, &arr)
	if err != nil {
		return err
	}
	v.X = arr[0]
	v.Y = arr[1]
	return nil
}
