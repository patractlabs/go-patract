package log

import (
	"github.com/patractlabs/go-patract/types"
)

// LoggerCodec a codec to encode args to log
type LoggerCodec struct {
	ss58Codec *types.SS58Codec
}

// NewLoggerCodec creates a new LoggerCodec
func NewLoggerCodec() *LoggerCodec {
	return &LoggerCodec{
		ss58Codec: types.NewSS58Codec(types.DefaultSS58Prefix),
	}
}

// SetSS58Codec set the ss58 codec to logger codec
func (l *LoggerCodec) SetSS58Codec(ss58Codec *types.SS58Codec) {
	l.ss58Codec = ss58Codec
}

func (l *LoggerCodec) TryEncodeArg(arg interface{}) (string, bool) {
	if accountID, ok := arg.(types.AccountID); ok {
		if l.ss58Codec == nil {
			return "", false
		}

		str, err := l.ss58Codec.EncodeAccountID(accountID)
		if err != nil {
			return "", false
		}
		return str, true
	}

	if accountID, ok := arg.(*types.AccountID); ok {
		if l.ss58Codec == nil {
			return "", false
		}

		str, err := l.ss58Codec.EncodeAccountID(*accountID)
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
