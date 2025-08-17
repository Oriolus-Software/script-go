package ffi

//go:wasm-module time
//export now
func now() int64

// GetCurrentTime returns the current time in nanoseconds from the host
func GetCurrentTime() int64 {
	return now()
}
