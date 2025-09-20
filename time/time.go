package time

import "time"

//go:wasm-module time
//export delta_f64
func Delta64() float64

//go:wasm-module time
//export ticks_alive
func TicksAlive() uint64

// gameTime in unix microseconds
//
//go:wasm-module time
//export game_time
func gameTime() int64

type GameTime int64

func GetGameTime() time.Time {
	return time.UnixMicro(gameTime())
}
