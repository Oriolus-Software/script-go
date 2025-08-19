package rand

//go:wasm-module rand
//export u64
func U64() uint64

//go:wasm-module rand
//export f64
func F64() float64

//go:wasm-module rand
//export random_seed
func RandomSeed() uint64

//go:wasm-module rand
//export seed
func Seed(seed uint64)
