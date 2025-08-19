package ffi

import (
	"unsafe"

	scriptgo "github.com/oriolus-software/script-go"
	"github.com/oriolus-software/script-go/internal/msgpack"
)

type FfiObject struct {
	ptr unsafe.Pointer
	len uint32
}

func Serialize(val any) FfiObject {
	data, err := msgpack.Serialize(val)
	if err != nil {
		panic(err)
	}

	memPtr := scriptgo.Allocate(len(data))
	buf := unsafe.Slice((*byte)(memPtr), len(data))
	copy(buf, data)

	return FfiObject{ptr: memPtr, len: uint32(len(data))}
}

func (o FfiObject) ToPacked() uint64 {
	return uint64(uintptr(o.ptr)) | uint64(o.len)<<32
}

func FromPacked(packed uint64) []byte {
	p := uintptr(packed & 0xFFFFFFFF)
	l := uint32(packed >> 32)
	return unsafe.Slice((*byte)(unsafe.Pointer(p)), int(l))
}
