package main

import "github.com/oriolus-software/script-go/vars"

//export init
func script_init() {
	vars.SetI64("on_init", 1)
}

//export tick
func tick() {
	// vars.SetI64("on_tick", 1)
}

//export late_tick
func late_tick() {}

//export register_actions
func register_actions() {
	// return ffi.Serialize([]any{}).ToPacked()
}
