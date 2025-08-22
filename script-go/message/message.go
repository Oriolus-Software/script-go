package message

import (
	"github.com/oriolus-software/script-go/internal/ffi"
)

type Meta struct {
	Namespace  string `msgpack:"namespace"`
	Identifier string `msgpack:"identifier"`
	Bus        string `msgpack:"bus"`
}

type Message struct {
	Meta   Meta          `msgpack:"meta"`
	Source MessageSource `msgpack:"source"`
	Value  any           `msgpack:"value"`
}

type MessageSource struct {
}

type Type interface {
	Meta() Meta
}

type Target interface {
	ToMessageTarget() any
}

type Myself struct{}

func (m Myself) ToMessageTarget() any {
	return "Myself"
}

type Parent struct{}

func (p Parent) ToMessageTarget() any {
	return "Parent"
}

type ChildByIndex uint32

func (c ChildByIndex) ToMessageTarget() any {
	return struct {
		ChildByIndex uint32
	}{
		ChildByIndex: uint32(c),
	}
}

type Cockpit uint8

func (c Cockpit) ToMessageTarget() any {
	return struct {
		Cockpit uint8
	}{
		Cockpit: uint8(c),
	}
}

type Broadcast struct {
	AcrossCouplings bool
	IncludeSelf     bool
}

func (b Broadcast) ToMessageTarget() any {
	type payload struct {
		AcrossCouplings bool `msgpack:"across_couplings"`
		IncludeSelf     bool `msgpack:"include_self"`
	}

	return struct {
		Broadcast payload
	}{
		Broadcast: payload{
			AcrossCouplings: b.AcrossCouplings,
			IncludeSelf:     b.IncludeSelf,
		},
	}
}

type AcrossCoupling struct {
	// Coupling is the name of the coupling to send the message across.
	// Either "Front" or "Back".
	Coupling string
	Cascade  bool
}

func Send(message Message, targets ...Target) {
	tgts := make([]any, len(targets))
	for i, target := range targets {
		tgts[i] = target.ToMessageTarget()
	}

	m := ffi.Serialize(message)
	t := ffi.Serialize(tgts)
	send(t.ToPacked(), m.ToPacked())
}

//go:wasm-module messages
//export take
func take() uint64

//go:wasm-module messages
//export send
func send(targets, message uint64)
