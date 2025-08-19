package ffi

import (
	"unsafe"

	"github.com/oriolus-software/script-go/msgpack"
)

type FfiObject []byte

func Serialize(val any) FfiObject {
	data, err := msgpack.Serialize(val)
	if err != nil {
		panic(err)
	}
	return FfiObject(data)
}

func (o FfiObject) ToPacked() uint64 {
	ptr := uintptr(unsafe.Pointer(unsafe.SliceData(o)))
	l := len(o)
	return uint64(ptr) | uint64(l)<<32
}

func FromPacked(ptr uint64) FfiObject {
	p := uintptr(ptr & 0xFFFFFFFF)
	l := uint32(ptr >> 32)
	return FfiObject(unsafe.Slice((*byte)(unsafe.Pointer(p)), int(l)))
}
