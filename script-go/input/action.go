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

var actions = make(map[string]registerAction)

func RegisterAction(id, defaultKey string) {
	actions[id] = registerAction{
		Id:         id,
		DefaultKey: defaultKey,
	}
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
func getState(actionId uint64) uint64

func State(actionId string) ActionState {
	var state ActionState
	ffi.DeserializeInto(getState(ffi.Serialize(actionId).ToPacked()), &state)
	return state
}

//export register_actions
func register_actions() {
	for _, action := range actions {
		register_action(ffi.Serialize(action).ToPacked())
	}
}
