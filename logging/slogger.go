package logging

import (
	"context"
	"log/slog"
	"os"
)

var _ Logger = (*SLogger)(nil)

type SLogger struct {
	core *slog.Logger
}

func levelToSlogLevel(level Level) slog.Level {
	switch level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelCritical:
		return slog.LevelError
	default:
		panic("invalid level")
	}
}

func NewSLogger(level Level) *SLogger {
	option := &slog.HandlerOptions{
		Level: levelToSlogLevel(level),
	}

	handler := slog.NewTextHandler(os.Stdout, option)

	core := slog.New(handler)
	if !core.Enabled(context.Background(), levelToSlogLevel(level)) {
		panic("Failed to enable level")
	}

	return &SLogger{core: core}
}

func (l *SLogger) With(a ...any) Logger {
	return &SLogger{
		core: l.core.With(a...),
	}
}

func (l *SLogger) Log(level Level, msg string, a ...any) {
	l.core.Log(context.Background(), levelToSlogLevel(level), msg, a...)
}

func (l *SLogger) Debug(msg string, a ...any) {
	l.Log(LevelDebug, msg, a...)
}

func (l *SLogger) Info(msg string, a ...any) {
	l.Log(LevelInfo, msg, a...)
}

func (l *SLogger) Warn(msg string, a ...any) {
	l.Log(LevelWarn, msg, a...)
}

func (l *SLogger) Critical(msg string, a ...any) {
	l.Log(LevelCritical, msg, a...)
}
