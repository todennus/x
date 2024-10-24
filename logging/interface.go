package logging

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelCritical
)

type Logger interface {
	With(a ...any) Logger
	Log(level Level, msg string, a ...any)
	Debug(msg string, a ...any)
	Info(msg string, a ...any)
	Warn(msg string, a ...any)
	Critical(msg string, a ...any)
}
