package message

import (
	"github.com/oriolus-software/script-go/internal/ffi"
	"github.com/oriolus-software/script-go/internal/msgpack"
)

var handlers = make(map[Meta]func(RawMessage))

//export late_tick
func late_tick() {
	messages := ffi.Deserialize[[]RawMessage](take())

	for _, message := range messages {

		handler, ok := handlers[message.Meta]
		if !ok {
			continue
		}

		handler(message)
	}
}

func RegisterHandler[T Message](handler func(Incoming[T])) {
	var proto T

	handlers[proto.Meta()] = func(message RawMessage) {
		raw, err := msgpack.Marshal(&message.Payload)
		if err != nil {
			return
		}

		var data T
		err = msgpack.Unmarshal(raw, &data)
		if err != nil {
			return
		}

		handler(Incoming[T]{
			Meta:    message.Meta,
			Payload: data,
		})
	}
}

type Incoming[T Message] struct {
	Meta    Meta
	Payload T
}
