package main

import (
	"fmt"

	"github.com/oriolus-software/script-go/assets"
	"github.com/oriolus-software/script-go/input"
	"github.com/oriolus-software/script-go/lmath"
	"github.com/oriolus-software/script-go/log"
	"github.com/oriolus-software/script-go/message"
	"github.com/oriolus-software/script-go/rand"
	"github.com/oriolus-software/script-go/texture"
	"github.com/oriolus-software/script-go/time"
	"github.com/oriolus-software/script-go/vars"
	"github.com/oriolus-software/script-go/vehicle"
)

//export init
func Init() {
	// message.Handler = func(in message.RawMessage) {
	// 	log.Info(fmt.Sprintf("[%d] received message: %v", time.TicksAlive(), in.Value))
	// }

	input.RegisterAction("doing_things", "key_g")

	message.RegisterHandler(func(m message.Incoming[TestMessage]) {
		// log.Info(fmt.Sprintf("received message: %+v", m.Meta))
	})

	vars.Set("on_init_ticks_alive", time.TicksAlive())
	vars.Set("on_init_bool_true", true)
	vars.Set("on_init_bool_false", false)
	vars.Set("on_init_string", "hello")

	vars.Set("doing_things", fmt.Sprintf("%v", input.State("doing_things").Kind == input.KindJustPressed))

	vars.Set("on_init_i64", 1)
	vars.Set("on_init_bool_true", true)
	vars.Set("on_init_bool_false", false)
	vars.Set("on_init_string", "hello")

	vars.Set("doing_things_state", fmt.Sprintf("%+v", input.State("doing_things")))

	log.Info("script initialized")

	t := texture.Create(texture.CreationOptions{
		Width:   100,
		Height:  100,
		MipMaps: false,
	})

	defer t.Dispose()

	t.DrawText(assets.ContentId{
		UserId: 1,
		SubId:  1,
	}, "Hello, world!", lmath.UVec2{X: 0, Y: 0}, 1, nil, texture.AlphaMask(0.5))

	t.DrawPixels([]texture.DrawPixel{
		{
			Pos:   lmath.UVec2{X: 0, Y: 0},
			Color: texture.Color{R: 255, G: 0, B: 0, A: 255},
		},
	})

	ot := texture.Create(texture.CreationOptions{
		Width:   100,
		Height:  100,
		MipMaps: false,
	})

	ot.DrawRect(lmath.UVec2{X: 0, Y: 0}, lmath.UVec2{X: 100, Y: 100}, texture.Color{R: 255, G: 0, B: 0, A: 255})

	t.DrawScriptTexture(ot, texture.DrawTextureOptions{
		SourceRect: &lmath.Rectangle{
			Start: lmath.UVec2{X: 0, Y: 0},
			End:   lmath.UVec2{X: 100, Y: 100},
		},
	})

	rand.RandomSeed()
}

//export tick
func Tick() {
	// for range 4200 {
	// vars.SetI64("on_tick_i64", int64(time.TicksAlive()))
	// }

	// vars.Set("_random", fmt.Sprintf("%v", rand.U64(0, 100)))
	vars.Set("delta", fmt.Sprintf("%v", time.Delta64()))
	// vars.Set("ticks_alive", fmt.Sprintf("%v", time.TicksAlive()))
	// vars.Set("game_time", fmt.Sprintf("%v", time.GetGameTime()))

	b, err := vehicle.GetBogie(0)
	if err != nil {
		log.Errorf("error getting bogie: %v", err)
	}

	a, err := b.GetAxle(0)
	if err != nil {
		log.Errorf("error getting axle: %v", err)
	}

	b.SetRailBrakeForceNewton(1000)

	// log.Info(fmt.Sprintf("rail quality: %v", a.RailQuality()))
	// log.Info(fmt.Sprintf("surface type: %v", a.SurfaceType()))
	// log.Info(fmt.Sprintf("inverse radius: %v", a.InverseRadius()))
	// log.Info(fmt.Sprintf("velocity vs ground: %v", vehicle.VelocityVsGround()))
	// log.Info(fmt.Sprintf("acceleration vs ground: %v", vehicle.AccelerationVsGround()))

	a.SetBrakeForceNewton(1000)
	a.SetTractionForceNewton(100000)

	// for i := 0; i < 100; i++ {
	// message.Send(TestMessage(true), message.Myself{})
	// }

	// log.Info(fmt.Sprintf("hello world on tick %d", time.TicksAlive()))

	// vars.SetI64("on_tick_i64", int64(time.TicksAlive()))
	// vars.SetF64("on_tick_f64", float64(time.TicksAlive()))
	// vars.SetString("on_tick_string", fmt.Sprintf("hello %d", time.TicksAlive()))

	// vars.SetString("mouse_delta", fmt.Sprintf("%v", input.MouseDelta()))

	// vars.SetString("throttle_state", fmt.Sprintf("%v", input.State("doing_things")))

	// vars.SetString("game_time", fmt.Sprintf("%v", time.GetGameTime()))

	// for i := 0; i < 1000; i++ {
	// 	message.Send(TestMessage(true), message.Broadcast{IncludeSelf: false})
	// }
}

type TestMessage bool

func (m TestMessage) Meta() message.Meta {
	return message.Meta{
		Namespace:  "test",
		Identifier: "test",
	}
}
