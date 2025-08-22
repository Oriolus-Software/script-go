package message

import "github.com/oriolus-software/script-go/internal/ffi"

var Handler func(RawMessage)

//export late_tick
func late_tick() {
	messages := ffi.Deserialize[[]RawMessage](take())
	for _, message := range messages {
		Handler(message)
	}
}
