package scriptgo

import "unsafe"

var allocatedMemory = make(map[unsafe.Pointer][]byte)

//export allocate
func Allocate(size int) unsafe.Pointer {
	data := make([]byte, size)
	ptr := unsafe.Pointer(&data[0])

	allocatedMemory[ptr] = data

	return ptr
}

//export deallocate
func Deallocate(ptr unsafe.Pointer) {
	delete(allocatedMemory, ptr)
}
