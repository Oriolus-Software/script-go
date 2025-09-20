package rand

// U64 returns a random 64-bit unsigned integer.
func U64(min, max uint64) uint64 {
	return u64(min, max)
}

// F64 returns a random 64-bit floating point number between 0 and 1.
func F64() float64 {
	return f64()
}

// RandomSeed seeds the random number generator with a random seed.
func RandomSeed() {
	randomSeed()
}

// Seed seeds the random number generator with a given seed.
func Seed(s uint64) {
	seed(s)
}

//go:wasm-module rand
//export u64
func u64(min, max uint64) uint64

//go:wasm-module rand
//export f64
func f64() float64

//go:wasm-module rand
//export random_seed
func randomSeed()

//go:wasm-module rand
//export seed
func seed(seed uint64)
