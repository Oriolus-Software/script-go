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
	Meta    Meta          `msgpack:"meta"`
	Source  MessageSource `msgpack:"source"`
	Payload any           `msgpack:"value"`
}

func (m *RawMessage) MarshalMsgpack(w *msgpack.Writer) error {
	w.WriteMapHeader(2)

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

	if err := w.WriteString("value"); err != nil {
		return err
	}

	if err := w.Encode(m.Payload); err != nil {
		return err
	}

	return nil
}

type MessageSource struct {
	Coupling               string `msgpack:"coupling"`                // "" if not set
	ModuleSlotIndex        int    `msgpack:"module_slot_index"`       // -1 if not set
	ModuleSlotCockpitIndex int    `msgpack:"module_slot_cockpit_idx"` // -1 if not set
}

func (m *MessageSource) UnmarshalMsgpack(r *msgpack.Reader) error {
	h, err := r.ReadMapHeader()
	if err != nil {
		return err
	}

	m.ModuleSlotIndex = -1
	m.ModuleSlotCockpitIndex = -1

	for i := 0; i < h; i++ {
		key, err := r.ReadString()
		if err != nil {
			return err
		}

		switch key {
		case "coupling":
			err := r.Decode(&m.Coupling)
			if err != nil {
				return err
			}
		case "module_slot_index":
			idx, err := r.ReadInt()
			if err == nil {
				m.ModuleSlotIndex = int(idx)
			}
		case "module_slot_cockpit_index":
			idx, err := r.ReadInt()
			if err == nil {
				m.ModuleSlotCockpitIndex = int(idx)
			}
		}
	}

	return nil
}

type Message interface {
	Meta() Meta
}

func Send(message Message, targets ...Target) {
	tgts := make([]any, len(targets))
	for i, target := range targets {
		tgts[i] = target.toMessageTarget()
	}

	meta := message.Meta()
	m := ffi.Serialize(&RawMessage{
		Meta:    meta,
		Payload: message,
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
