package log

import "github.com/oriolus-software/script-go/internal/ffi"

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

func Write(level int, message string) {
	m := ffi.Serialize(message)
	write(level, m.ToPacked())
}

func Debug(message string) {
	Write(LevelDebug, message)
}

func Info(message string) {
	Write(LevelInfo, message)
}

func Warn(message string) {
	Write(LevelWarn, message)
}

func Error(message string) {
	Write(LevelError, message)
}

//go:wasm-module log
//export write
func write(level int, message uint64)
