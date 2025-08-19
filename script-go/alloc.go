package scriptgo

import "unsafe"

// Use a small static ring buffer arena in linear memory. This avoids GC
// interaction and keeps pointers stable within a single host call.
var arena [8 * 1024]byte // 8 KiB static arena
var arenaOffset int

//export allocate
func Allocate(size int) unsafe.Pointer {
	if size <= 0 || size > len(arena) {
		panic("invalid allocation size")
	}
	// 8-byte alignment
	const align = 8
	if rem := arenaOffset % align; rem != 0 {
		arenaOffset += align - rem
	}
	// Check if we need to wrap around AFTER alignment
	if arenaOffset+size > len(arena) {
		arenaOffset = 0
		// No need to re-align since we're starting at 0, which is already aligned
	}

	ptr := unsafe.Pointer(&arena[arenaOffset])
	arenaOffset += size
	return ptr
}

//export deallocate
func Deallocate(ptr unsafe.Pointer) {
	// no-op: ring buffer does not free individual blocks
}
