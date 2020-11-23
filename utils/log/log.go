package log

import (
	"fmt"

	"go.uber.org/zap"
)

const (
	kvLen = 2
)

// Logger a logger interface
type Logger interface {
	Debug(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// loggerByZap logger imp zap
type loggerByZap struct {
	l *zap.Logger
}

// NewLogger create logger
func NewLogger() Logger {
	// TODO: log config
	return &loggerByZap{
		l: zap.NewExample(),
	}
}

// Debug imp debug log
func (l *loggerByZap) Debug(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.l.Debug(msg)
		return
	}

	if len(args)%kvLen != 0 {
		args = append(args, "missing key-value")
	}

	fields := make([]zap.Field, 0, len(args)/kvLen)
	for i := 0; i < len(args); i += kvLen {
		fields = append(fields, zap.Any(fmt.Sprintf("%s", args[i]), args[i+1]))
	}

	l.l.Debug(msg, fields...)
}

// Error imp Error log
func (l *loggerByZap) Error(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.l.Error(msg)
		return
	}

	if len(args)%kvLen != 0 {
		args = append(args, "missing key-value")
	}

	fields := make([]zap.Field, 0, len(args)/kvLen)
	for i := 0; i < len(args); i += kvLen {
		fields = append(fields, zap.Any(fmt.Sprintf("%s", args[i]), args[i+1]))
	}

	l.l.Error(msg, fields...)
}
