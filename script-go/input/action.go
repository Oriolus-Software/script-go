package input

import (
	"github.com/oriolus-software/script-go/internal/ffi"
)

const (
	StateNone = iota
	StateJustPressed
	StatePressed
	StateJustReleased
)

type registerAction struct {
	Id         string `msgpack:"id"`
	DefaultKey string `msgpack:"default_key"`
}

func RegisterAction(id, defaultKey string) {
	register_action(ffi.Serialize(registerAction{
		Id:         id,
		DefaultKey: defaultKey,
	}).ToPacked())
}

//go:wasm-module action
//export register
func register_action(action uint64)

type ActionState struct {
	Kind         int  `msgpack:"kind"`
	CockpitIndex *int `msgpack:"cockpit_index"`
}

//go:wasm-module action
//export state
func state(actionId uint64) uint64

func State(actionId string) ActionState {
	return ffi.Deserialize[ActionState](state(ffi.Serialize(actionId).ToPacked()))
}
