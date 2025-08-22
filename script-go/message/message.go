package message

import (
	"github.com/oriolus-software/script-go/internal/ffi"
	"github.com/oriolus-software/script-go/internal/msgpack"
)

type Meta struct {
	Namespace  string `msgpack:"namespace"`
	Identifier string `msgpack:"identifier"`
	Bus        string `msgpack:"bus"`
}

type RawMessage struct {
	Meta   Meta          `msgpack:"meta"`
	Source MessageSource `msgpack:"source"`
	Value  any           `msgpack:"value"`
}

func (m *RawMessage) MarshalMsgpack(w *msgpack.Writer) error {
	w.WriteMapHeader(3)

	if err := w.WriteString("meta"); err != nil {
		return err
	}

	{
		header := 2
		if m.Meta.Bus != "" {
			header++
		}

		w.WriteMapHeader(header)

		if err := w.WriteString("namespace"); err != nil {
			return err
		}

		if err := w.WriteString(m.Meta.Namespace); err != nil {
			return err
		}

		if err := w.WriteString("identifier"); err != nil {
			return err
		}

		if err := w.WriteString(m.Meta.Identifier); err != nil {
			return err
		}

		if m.Meta.Bus != "" {
			if err := w.WriteString("bus"); err != nil {
				return err
			}

			if err := w.WriteString(m.Meta.Bus); err != nil {
				return err
			}
		}
	}

	if err := w.WriteString("source"); err != nil {
		return err
	}

	if err := w.WriteStruct(m.Source); err != nil {
		return err
	}

	if err := w.WriteString("value"); err != nil {
		return err
	}

	if err := w.Encode(m.Value); err != nil {
		return err
	}

	return nil
}

type MessageSource struct {
}

type Message interface {
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

	m := ffi.Serialize(&RawMessage{
		Meta:  message.Meta(),
		Value: message,
	})
	t := ffi.Serialize(tgts)
	send(t.ToPacked(), m.ToPacked())
}

//go:wasm-module messages
//export take
func take() uint64

//go:wasm-module messages
//export send
func send(targets, message uint64)
