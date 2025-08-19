package time

//go:wasm-module time
//export delta_f64
func Delta64() float64

//go:wasm-module time
//export ticks_alive
func TicksAlive() uint64

//TODO: game_time
//go:wasm-module time
//export game_time
// func gameTime() int64
