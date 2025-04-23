package seventveventapi

import "log/slog"

type Logger interface {
	Info(input string, fields ...any)
	Error(input string, fields ...any)
	Debug(input string, fields ...any)
	Warn(input string, fields ...any)

	WithComponent(name string) Logger
}

type SlogLogger struct {
	slog *slog.Logger
}

// Ensure SlogLogger implements your Logger interface
var _ Logger = (*SlogLogger)(nil)

func NewSlogLogger(base *slog.Logger) *SlogLogger {
	return &SlogLogger{slog: base}
}

func (l *SlogLogger) Info(msg string, fields ...any) {
	l.slog.Info(msg, fields...)
}

func (l *SlogLogger) Error(msg string, fields ...any) {
	l.slog.Error(msg, fields...)
}

func (l *SlogLogger) Debug(msg string, fields ...any) {
	l.slog.Debug(msg, fields...)
}

func (l *SlogLogger) Warn(msg string, fields ...any) {
	l.slog.Warn(msg, fields...)
}

func (l *SlogLogger) WithComponent(name string) Logger {
	return &SlogLogger{
		slog: l.slog.With("component", name),
	}
}
