package alloc

import (
	"fmt"
	"unsafe"
)

// Use a small static ring buffer hostArena in linear memory. This avoids GC
// interaction and keeps pointers stable within a single host call.
var hostArena *Arena = NewArena(8 * 1024) // 8 KiB static arena
var arenaOffset int

//export allocate
func Allocate(size int) unsafe.Pointer {
	return hostArena.Allocate(size)
}

//export deallocate
func Deallocate(ptr unsafe.Pointer) {
	// no-op: ring buffer does not free individual blocks
}

type Arena struct {
	data   []byte
	len    int
	offset int
}

func NewArena(size int) *Arena {
	arena := make([]byte, size)
	return &Arena{
		data:   arena,
		len:    size,
		offset: 0,
	}
}

func (a *Arena) Allocate(size int) unsafe.Pointer {
	if size <= 0 || size > a.len {
		panic(fmt.Sprintf("invalid allocation size: %d", size))
	}
	// 8-byte alignment
	const align = 8
	if rem := arenaOffset % align; rem != 0 {
		arenaOffset += align - rem
	}
	// Check if we need to wrap around AFTER alignment
	if arenaOffset+size > a.len {
		arenaOffset = 0
		// No need to re-align since we're starting at 0, which is already aligned
	}

	ptr := unsafe.Pointer(&a.data[a.offset])
	arenaOffset += size
	return ptr
}

func (a *Arena) AllocateSlice(size int) []byte {
	ptr := a.Allocate(size)
	return unsafe.Slice((*byte)(ptr), size)
}
