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

	message.Handler = func(in message.RawMessage) {
		log.Info(fmt.Sprintf("received message: %v", in.Value))
	}

	vars.SetI64("on_tick_i64", int64(time.TicksAlive()))
	vars.SetF64("on_tick_f64", float64(time.TicksAlive()))
	vars.SetString("on_tick_string", fmt.Sprintf("hello %d", time.TicksAlive()))

	vars.SetString("mouse_delta", fmt.Sprintf("%v", input.MouseDelta()))

	vars.SetString("throttle_state", fmt.Sprintf("%v", input.State("doing_things")))

	vars.SetString("game_time", fmt.Sprintf("%v", time.GetGameTime()))

	message.Send(TestMessage{Value: "test"}, message.Myself{})
}

type TestMessage struct {
	Value string
}

func (m TestMessage) Meta() message.Meta {
	return message.Meta{
		Namespace:  "test",
		Identifier: "test",
	}
}

//export register_actions
func RegisterActions() {
	log.Info("registering actions")
	input.RegisterAction("doing_things", "key_g")
}
