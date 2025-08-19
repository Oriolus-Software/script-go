package env

//go:wasm-module env
//export is_rc
func IsRC() bool
