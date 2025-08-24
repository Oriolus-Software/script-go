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
	input.RegisterAction("doing_things", "key_g")

	message.RegisterHandler(func(m message.Incoming[TestMessage]) {
		log.Infof("received message: %+v", m.Source)
	})

	vars.Set("on_init_ticks_alive", time.TicksAlive())
	vars.Set("on_init_bool_true", true)
	vars.Set("on_init_bool_false", false)
	vars.Set("on_init_string", "hello")

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

	t.DrawText(&texture.DrawTextOptions{
		Font: assets.ContentId{
			UserId: 1,
			SubId:  1,
		},
		Text:          "Hello, world!",
		TopLeft:       lmath.UVec2{X: 0, Y: 0},
		LetterSpacing: 1,
		FullColor:     nil,
		AlphaMode:     texture.AlphaMask(0.5),
	})

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
		SourceRect: lmath.Rectangle{
			Start: lmath.UVec2{X: 0, Y: 0},
			End:   lmath.UVec2{X: 100, Y: 100},
		},
	})

	rand.RandomSeed()
}

//export tick
func Tick() {
	vars.Set("_random", fmt.Sprintf("%v", rand.U64(0, 100)))
	vars.Set("delta", fmt.Sprintf("%v", time.Delta64()))
	vars.Set("ticks_alive", fmt.Sprintf("%v", time.TicksAlive()))
	vars.Set("game_time", fmt.Sprintf("%v", time.GetGameTime()))

	b, err := vehicle.GetBogie(0)
	if err != nil {
		log.Errorf("error getting bogie: %v", err)
	}

	a, err := b.GetAxle(0)
	if err != nil {
		log.Errorf("error getting axle: %v", err)
	}

	b.SetRailBrakeForceNewton(1000)

	a.SetBrakeForceNewton(1000)
	a.SetTractionForceNewton(100000)

	log.Infof("rail quality: %v", a.RailQuality())
	log.Infof("surface type: %v", a.SurfaceType())
	log.Infof("inverse radius: %v", a.InverseRadius())
	log.Infof("velocity vs ground: %v", vehicle.VelocityVsGround())
	log.Infof("acceleration vs ground: %v", vehicle.AccelerationVsGround())

	log.Infof("hello world on tick %d", time.TicksAlive())

	vars.Set("on_tick_string", fmt.Sprintf("hello %d", time.TicksAlive()))
	vars.Set("mouse_delta", fmt.Sprintf("%v", input.MouseDelta()))
	vars.Set("throttle_state", fmt.Sprintf("%v", input.State("doing_things")))
	vars.Set("game_time", fmt.Sprintf("%v", time.GetGameTime()))
	vars.Set("random", fmt.Sprintf("%v", rand.U64(0, 100)))

	if input.State("doing_things").IsJustPressed() {
		log.Info("doing things just pressed")
	}

	message.Send(TestMessage(true), message.Myself)
	message.Send(TestMessage(true), message.Parent)
	message.Send(TestMessage(true), message.ChildByIndex(0))
	message.Send(TestMessage(true), message.Cockpit(0))
	message.Send(TestMessage(true), message.Broadcast{IncludeSelf: false, AcrossCouplings: false})
	message.Send(TestMessage(true), message.AcrossCoupling{Coupling: "rear", Cascade: false})
}

type TestMessage bool

func (m TestMessage) Meta() message.Meta {
	return message.Meta{
		Namespace:  "test",
		Identifier: "test",
	}
}
