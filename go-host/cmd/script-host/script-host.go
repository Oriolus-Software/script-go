package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/vmihailenco/msgpack/v5"
)

type ScriptHost struct {
	r wazero.Runtime
}

func NewScriptHost() *ScriptHost {
	return &ScriptHost{
		r: wazero.NewRuntime(context.Background()),
	}
}

func (h *ScriptHost) LoadModule(ctx context.Context, data []byte) (api.Module, error) {
	mod, err := h.r.CompileModule(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("compile module: %w", err)
	}

	v := h.r.NewHostModuleBuilder("var")
	if _, err := v.NewFunctionBuilder().WithFunc(func(ctx context.Context, m api.Module, a uint64, b int64) {
		// Debug: inspect packed pointer/length and raw bytes
		ptr := uint32(a & 0xFFFFFFFF)
		l := uint32(a >> 32)
		data, ok := m.Memory().Read(ptr, l)
		if !ok {
			panic("failed to read memory")
		}
		fmt.Printf("set_i64 raw: ptr=0x%08x len=%d bytes=%v memSize=%d\n", ptr, l, data, m.Memory().Size())

		var name string
		FromPacked(m, a, &name)

		if name != "on_init" {
			fmt.Printf("wrong call to var::set_i64 name=%q value=%d\n", name, b)
			os.Exit(1)
		}
	}).Export("set_i64").Instantiate(ctx); err != nil {
		return nil, fmt.Errorf("instantiate var module: %w", err)
	}

	return h.r.InstantiateModule(ctx, mod, wazero.NewModuleConfig())
}

func main() {
	scriptData, err := os.ReadFile("../dev-script/dev-script.wasm")
	if err != nil {
		panic(err)
	}

	h := NewScriptHost()
	mod, err := h.LoadModule(context.Background(), scriptData)
	if err != nil {
		panic(err)
	}

	f := mod.ExportedFunction("init")

	executions := 0

	for range 100000 {
		_, err := f.Call(context.Background())
		if err != nil {
			panic(fmt.Errorf("call init on execution %d: %w", executions, err))
		}

		executions++

		// fmt.Println(mod.Memory().Size())
	}
}

func FromPacked[T any](m api.Module, packed uint64, dst *T) {
	ptr := uint32(packed & 0xFFFFFFFF)
	l := uint32(packed >> 32)

	data, ok := m.Memory().Read(uint32(ptr), l)
	if !ok {
		panic("failed to read memory")
	}

	if err := msgpack.Unmarshal(data, dst); err != nil {
		panic(err)
	}
}
