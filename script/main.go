package main

import (
	"example-script/ffi"
	"unsafe"
)

//export add
func Add(a, b int) int {
	return a + b
}

//export addWithTime
func AddWithTime(a, b int) int64 {
	result := a + b
	timestamp := ffi.GetCurrentTime()

	return int64(result) + timestamp
}

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

func main() {}
