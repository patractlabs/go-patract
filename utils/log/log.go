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
	Flush()

	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// loggerByZap logger imp zap
type loggerByZap struct {
	l     *zap.Logger
	codec *LoggerCodec
}

// NewLogger create logger
func NewLogger() Logger {
	// TODO: log config
	return &loggerByZap{
		l:     zap.NewExample(),
		codec: NewLoggerCodec(),
	}
}

// NewNopLogger create a nop logger
func NewNopLogger() Logger {
	return &loggerByZap{
		l: zap.NewNop(),
	}
}

// Flush sync logger
func (l *loggerByZap) Flush() {
	err := l.l.Sync()
	if err != nil {
		fmt.Printf("logger flush error %s", err.Error())
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
		str, ok := l.codec.TryEncodeArg(args[i+1])
		if !ok {
			fields = append(fields, zap.Any(fmt.Sprintf("%s", args[i]), args[i+1]))
		} else {
			fields = append(fields, zap.Any(fmt.Sprintf("%s", args[i]), str))
		}
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
		str, ok := l.codec.TryEncodeArg(args[i+1])
		if !ok {
			fields = append(fields, zap.Any(fmt.Sprintf("%s", args[i]), args[i+1]))
		} else {
			fields = append(fields, zap.Any(fmt.Sprintf("%s", args[i]), str))
		}
	}

	l.l.Error(msg, fields...)
}

// Warn imp Warn log
func (l *loggerByZap) Warn(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.l.Warn(msg)
		return
	}

	if len(args)%kvLen != 0 {
		args = append(args, "missing key-value")
	}

	fields := make([]zap.Field, 0, len(args)/kvLen)
	for i := 0; i < len(args); i += kvLen {
		str, ok := l.codec.TryEncodeArg(args[i+1])
		if !ok {
			fields = append(fields, zap.Any(fmt.Sprintf("%s", args[i]), args[i+1]))
		} else {
			fields = append(fields, zap.Any(fmt.Sprintf("%s", args[i]), str))
		}
	}

	l.l.Warn(msg, fields...)
}

// Info imp Info log
func (l *loggerByZap) Info(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.l.Info(msg)
		return
	}

	if len(args)%kvLen != 0 {
		args = append(args, "missing key-value")
	}

	fields := make([]zap.Field, 0, len(args)/kvLen)
	for i := 0; i < len(args); i += kvLen {
		str, ok := l.codec.TryEncodeArg(args[i+1])
		if !ok {
			fields = append(fields, zap.Any(fmt.Sprintf("%s", args[i]), args[i+1]))
		} else {
			fields = append(fields, zap.Any(fmt.Sprintf("%s", args[i]), str))
		}
	}

	l.l.Info(msg, fields...)
}
