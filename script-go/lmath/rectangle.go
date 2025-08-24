package lmath

type Rectangle struct {
	Start UVec2 `msgpack:"start"`
	End   UVec2 `msgpack:"end"`
}
