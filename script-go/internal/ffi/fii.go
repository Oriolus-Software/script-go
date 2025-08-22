package ffi

import (
	"encoding/binary"
	"unsafe"

	"github.com/oriolus-software/script-go/internal/alloc"
	"github.com/oriolus-software/script-go/internal/msgpack"
)

type FfiObject struct {
	ptr unsafe.Pointer
	len uint32
}

func Serialize(val any) FfiObject {
	data, err := msgpack.Marshal(val)
	if err != nil {
		panic(err)
	}

	memPtr := alloc.Allocate(len(data))
	buf := unsafe.Slice((*byte)(memPtr), len(data))
	copy(buf, data)

	return FfiObject{ptr: memPtr, len: uint32(len(data))}
}

func (o FfiObject) ToPacked() uint64 {
	// Match Rust format: [ptr: 4 bytes][len: 4 bytes] in big-endian
	var packed [8]byte
	binary.BigEndian.PutUint32(packed[:4], uint32(uintptr(o.ptr)))
	binary.BigEndian.PutUint32(packed[4:], o.len)
	return binary.BigEndian.Uint64(packed[:])
}

func FromPacked(packed uint64) []byte {
	// Match Rust format: [ptr: 4 bytes][len: 4 bytes] in big-endian
	var packedBytes [8]byte
	binary.BigEndian.PutUint64(packedBytes[:], packed)

	ptr := binary.BigEndian.Uint32(packedBytes[:4])
	len := binary.BigEndian.Uint32(packedBytes[4:])

	// Safe: ptr is a valid WASM memory address from the allocator
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), int(len))
}

func Deserialize[T any](packed uint64) T {
	var result T
	data := FromPacked(packed)
	err := msgpack.Unmarshal(data, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func DeserializeInto[T any](packed uint64, v *T) {
	data := FromPacked(packed)
	err := msgpack.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}
