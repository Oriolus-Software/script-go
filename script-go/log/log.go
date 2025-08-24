package log

import (
	"fmt"

	"github.com/oriolus-software/script-go/internal/ffi"
)

const (
	levelDebug = iota
	levelInfo
	levelWarn
	levelError
)

func Debug(message string) {
	writeTrampoline(levelDebug, message)
}

func Debugf(format string, a ...any) {
	writeTrampoline(levelDebug, fmt.Sprintf(format, a...))
}

func Info(message string) {
	writeTrampoline(levelInfo, message)
}

func Infof(format string, a ...any) {
	writeTrampoline(levelInfo, fmt.Sprintf(format, a...))
}

func Warn(message string) {
	writeTrampoline(levelWarn, message)
}

func Warnf(format string, a ...any) {
	writeTrampoline(levelWarn, fmt.Sprintf(format, a...))
}

func Error(message string) {
	writeTrampoline(levelError, message)
}

func Errorf(format string, a ...any) {
	writeTrampoline(levelError, fmt.Sprintf(format, a...))
}

func writeTrampoline(level int, message string) {
	m := ffi.Serialize(message)
	write(level, m.ToPacked())
}

//go:wasm-module log
//export write
func write(level int, message uint64)
