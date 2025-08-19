package assets

import "github.com/oriolus-software/script-go/internal/ffi"

type ContentId struct {
	UserId int `msgpack:"user_id"`
	SubId  int `msgpack:"sub_id"`
}

//go:wasm-module assets
//export preload
func preload(contentId uint64)

func Preload(contentId ContentId) {
	preload(ffi.Serialize(contentId).ToPacked())
}
