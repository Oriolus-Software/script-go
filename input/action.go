package input

import (
	"github.com/oriolus-software/script-go/internal/ffi"
	"github.com/oriolus-software/script-go/lmath"
)

const (
	KindNone = iota
	KindJustPressed
	KindPressed
	KindJustReleased
)

type registerAction struct {
	Id         string `msgpack:"id"`
	DefaultKey string `msgpack:"default_key"`
}

var actions = make(map[string]registerAction, 0)

func RegisterAction(id, defaultKey string) {
	actions[id] = registerAction{
		Id:         id,
		DefaultKey: defaultKey,
	}
}

type ActionState struct {
	Kind         int         `msgpack:"kind"`
	CockpitIndex *int        `msgpack:"cockpit_index"`
	UV           *lmath.Vec2 `msgpack:"uv"`
}

func (a ActionState) IsPressed() bool {
	return a.Kind == KindPressed
}

func (a ActionState) IsJustPressed() bool {
	return a.Kind == KindJustPressed
}

func (a ActionState) IsJustReleased() bool {
	return a.Kind == KindJustReleased
}

func (a ActionState) IsNone() bool {
	return a.Kind == KindNone
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

	actions = make(map[string]registerAction, 0)
}

//go:wasm-module action
//export register
func register_action(action uint64)
