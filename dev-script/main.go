package main

import (
	"fmt"

	"github.com/oriolus-software/script-go/input"
	"github.com/oriolus-software/script-go/log"
	"github.com/oriolus-software/script-go/message"
	"github.com/oriolus-software/script-go/script"
	"github.com/oriolus-software/script-go/vars"
)

func init() {
	input.RegisterAction("doing_things", "key_g")
}

//export init
func Init() {
	// message.Handler = func(in message.RawMessage) {
	// 	log.Info(fmt.Sprintf("[%d] received message: %v", time.TicksAlive(), in.Value))
	// }

	vars.SetI64("on_init_i64", 1)
	vars.SetBool("on_init_bool_true", true)
	vars.SetBool("on_init_bool_false", false)
	vars.SetString("on_init_string", "hello")

	vars.SetString("doing_things", fmt.Sprintf("%v", input.State("doing_things")))

	log.Info("script initialized")

	script.OnTick = Tick
}

func Tick() {
	// for range 4000 {
	// 	vars.SetI64("on_tick_i64", int64(time.TicksAlive()))
	// }

	// log.Info(fmt.Sprintf("hello world on tick %d", time.TicksAlive()))

	// vars.SetI64("on_tick_i64", int64(time.TicksAlive()))
	// vars.SetF64("on_tick_f64", float64(time.TicksAlive()))
	// vars.SetString("on_tick_string", fmt.Sprintf("hello %d", time.TicksAlive()))

	// vars.SetString("mouse_delta", fmt.Sprintf("%v", input.MouseDelta()))

	// vars.SetString("throttle_state", fmt.Sprintf("%v", input.State("doing_things")))

	// vars.SetString("game_time", fmt.Sprintf("%v", time.GetGameTime()))

	for i := 0; i < 1000; i++ {
		message.Send(TestMessage(true), message.Broadcast{IncludeSelf: false})
	}
}

type TestMessage bool

func (m TestMessage) Meta() message.Meta {
	return message.Meta{
		Namespace:  "test",
		Identifier: "test",
	}
}
