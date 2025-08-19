package message

import "github.com/oriolus-software/script-go/internal/ffi"

var Handler func(Message)

//export late_tick
func late_tick() {
	messages := ffi.Deserialize[[]Message](take())
	for _, message := range messages {
		Handler(message)
	}
}
