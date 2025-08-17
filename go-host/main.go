package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/tetratelabs/wazero"
)

func main() {
	r := wazero.NewRuntime(context.Background())

	wasmBytes, err := os.ReadFile("../script/main.wasm")
	if err != nil {
		panic(err)
	}

	r.NewHostModuleBuilder("time").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context) int64 {
			return time.Now().UnixNano()
		}).
		Export("now").
		Instantiate(context.Background())

	compiled, err := r.CompileModule(context.Background(), wasmBytes)
	if err != nil {
		panic(err)
	}

	mod, err := r.InstantiateModule(context.Background(), compiled, wazero.NewModuleConfig())
	if err != nil {
		panic(err)
	}

	const runs = 10

	for i := 0; i < runs; i++ {
		start := time.Now()
		_, err := mod.ExportedFunction("add").Call(context.Background(), 1, 2)
		elapsed := time.Since(start)

		if err != nil {
			panic(err)
		}
		fmt.Printf("add: %v\n", elapsed)
	}

	for i := 0; i < runs; i++ {
		start := time.Now()
		_, err := mod.ExportedFunction("addWithTime").Call(context.Background(), 1, 2)
		elapsed := time.Since(start)

		if err != nil {
			panic(err)
		}
		fmt.Printf("addWithTime: %v\n", elapsed)
	}

	memory := runtime.MemStats{}
	runtime.ReadMemStats(&memory)
	fmt.Printf("Alloc: %s\n", Bytes(memory.Alloc))
	fmt.Printf("TotalAlloc: %s\n", Bytes(memory.TotalAlloc))
	fmt.Printf("Sys: %s\n", Bytes(memory.Sys))
	fmt.Printf("NumGC: %d\n", memory.NumGC)
}

type Bytes uint64

func (b Bytes) String() string {
	if b < 1024 {
		return fmt.Sprintf("%d B", b)
	} else if b < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(b)/1024)
	} else if b < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(b)/(1024*1024))
	} else {
		return fmt.Sprintf("%.1f GB", float64(b)/(1024*1024*1024))
	}
}
