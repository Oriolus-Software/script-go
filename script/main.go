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

//export allocate
func Allocate(size int) unsafe.Pointer {
}

func main() {}
