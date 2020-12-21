package log

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/patractlabs/go-patract/utils"
	"go.uber.org/zap"
)

const (
	kvLen = 2
)

// Logger a logger interface
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
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

// NewNopLogger create a nop logger
func NewNopLogger() Logger {
	return &loggerByZap{
		l: zap.NewNop(),
	}
}

func tryProcessTypes(arg interface{}) (string, bool) {
	if accountID, ok := arg.(types.AccountID); ok {
		str, err := utils.EncodeAccountIDToSS58(accountID)
		if err != nil {
			return "", false
		}
		return str, true
	}

	if accountID, ok := arg.(*types.AccountID); ok {
		str, err := utils.EncodeAccountIDToSS58(*accountID)
		if err != nil {
			return "", false
		}
		return str, true
	}

	if bz, ok := arg.([]byte); ok {
		return types.HexEncodeToString(bz), true
	}

	if bz, ok := arg.(types.Bytes); ok {
		return types.HexEncodeToString(bz), true
	}

	return "", false
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
		str, ok := tryProcessTypes(args[i+1])
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
		str, ok := tryProcessTypes(args[i+1])
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
		str, ok := tryProcessTypes(args[i+1])
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
		str, ok := tryProcessTypes(args[i+1])
		if !ok {
			fields = append(fields, zap.Any(fmt.Sprintf("%s", args[i]), args[i+1]))
		} else {
			fields = append(fields, zap.Any(fmt.Sprintf("%s", args[i]), str))
		}
	}

	l.l.Info(msg, fields...)
}
