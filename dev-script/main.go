package main

import (
	"fmt"

	"github.com/oriolus-software/script-go/input"
	"github.com/oriolus-software/script-go/log"
	"github.com/oriolus-software/script-go/message"
	"github.com/oriolus-software/script-go/time"
	"github.com/oriolus-software/script-go/vars"
)

//export init
func ScriptInit() {
	vars.SetI64("on_init_i64", 1)
	vars.SetBool("on_init_bool_true", true)
	vars.SetBool("on_init_bool_false", false)
	vars.SetString("on_init_string", "hello")
}

//export tick
func Tick() {
	// for range 4000 {
	// 	vars.SetI64("on_tick_i64", int64(time.TicksAlive()))
	// }

	// log.Info(fmt.Sprintf("hello world on tick %d", time.TicksAlive()))

	message.Handler = func(in message.Message) {
		log.Info(fmt.Sprintf("received message: %v", in))
	}

	vars.SetI64("on_tick_i64", int64(time.TicksAlive()))
	vars.SetF64("on_tick_f64", float64(time.TicksAlive()))
	vars.SetString("on_tick_string", fmt.Sprintf("hello %d", time.TicksAlive()))

	vars.SetString("mouse_delta", fmt.Sprintf("%v", input.MouseDelta()))

	vars.SetString("move_forward_state", fmt.Sprintf("%v", input.State("move_forward")))
	vars.SetString("move_forward_state", fmt.Sprintf("%v", input.State("key_s")))

	vars.SetString("game_time", fmt.Sprintf("%v", time.GetGameTime()))
}

//export register_actions
func RegisterActions() {
	input.RegisterAction("move_forward", "key_w")
}
