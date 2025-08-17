package main_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/tetratelabs/wazero"
)

// Baseline functions - native Go equivalents of WASM functions

//go:noinline
func baselineAdd(a, b int) int {
	return a + b
}

//go:noinline
func baselineAddWithTime(a, b int) int64 {
	result := a + b
	timestamp := time.Now().UnixNano()
	return int64(result) + timestamp
}

func BenchmarkAdd(b *testing.B) {
	r := wazero.NewRuntime(context.Background())

	// Set up the host time module that the WASM module imports
	r.NewHostModuleBuilder("time").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context) int64 {
			return time.Now().UnixNano()
		}).
		Export("now").
		Instantiate(context.Background())

	wasmBytes, err := os.ReadFile("../script/main.wasm")
	if err != nil {
		panic(err)
	}

	compiled, err := r.CompileModule(context.Background(), wasmBytes)
	if err != nil {
		panic(err)
	}

	mod, err := r.InstantiateModule(context.Background(), compiled, wazero.NewModuleConfig())
	if err != nil {
		panic(err)
	}

	b.ReportAllocs()

	add := mod.ExportedFunction("add")
	ctx := context.Background()

	for b.Loop() {
		if _, err := add.Call(ctx, 1, 2); err == nil {
		} else {
			panic(err)
		}
	}
}

func BenchmarkAddWithTime(b *testing.B) {
	r := wazero.NewRuntime(context.Background())

	r.NewHostModuleBuilder("time").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context) int64 {
			return time.Now().UnixNano()
		}).
		Export("now").
		Instantiate(context.Background())

	wasmBytes, err := os.ReadFile("../script/main.wasm")
	if err != nil {
		panic(err)
	}

	compiled, err := r.CompileModule(context.Background(), wasmBytes)
	if err != nil {
		panic(err)
	}

	mod, err := r.InstantiateModule(context.Background(), compiled, wazero.NewModuleConfig())
	if err != nil {
		panic(err)
	}

	b.ReportAllocs()

	addWithTime := mod.ExportedFunction("addWithTime")
	ctx := context.Background()

	for b.Loop() {
		if _, err := addWithTime.Call(ctx, 1, 2); err == nil {
		} else {
			panic(err)
		}
	}
}

// Baseline benchmarks - native Go function calls without WASM
func BenchmarkBaselineAdd(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		_ = baselineAdd(1, 2)
	}
}

func BenchmarkBaselineAddWithTime(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		_ = baselineAddWithTime(1, 2)
	}
}

func BenchmarkAllocDealloc(b *testing.B) {
	r := wazero.NewRuntime(context.Background())

	// Set up the host time module that the WASM module imports
	r.NewHostModuleBuilder("time").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context) int64 {
			return time.Now().UnixNano()
		}).
		Export("now").
		Instantiate(context.Background())

	wasmBytes, err := os.ReadFile("../script/main.wasm")
	if err != nil {
		panic(err)
	}

	compiled, err := r.CompileModule(context.Background(), wasmBytes)
	if err != nil {
		panic(err)
	}

	mod, err := r.InstantiateModule(context.Background(), compiled, wazero.NewModuleConfig())
	if err != nil {
		panic(err)
	}

	b.ReportAllocs()

	allocate := mod.ExportedFunction("allocate")
	deallocate := mod.ExportedFunction("deallocate")
	ctx := context.Background()

	const allocSize = 128

	for b.Loop() {
		results, err := allocate.Call(ctx, uint64(allocSize))
		if err != nil {
			panic(err)
		}
		ptr := results[0]

		_, err = deallocate.Call(ctx, ptr)
		if err != nil {
			panic(err)
		}
	}

	b.ReportMetric(float64(mod.Memory().Size()), "memory")
}
