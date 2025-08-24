package env

func IsRC() bool {
	return isRC()
}

//go:wasm-module env
//export is_rc
func isRC() bool
