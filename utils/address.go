package utils

import (
	"github.com/patractlabs/go-patract/types"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
)

var (
	DecodeAccountIDFromSS58 = types.DecodeAccountIDFromSS58
	EncodeAccountIDToSS58   = types.EncodeAccountIDToSS58
)

// NewAccountIDFromSS58 TODO: to go-sdk
func NewAccountIDFromSS58(address string) (types.AccountID, error) {
	return DecodeAccountIDFromSS58(address)
}

// MustAccountIDFromSS58 must account for ss58 if not panic
func MustAccountIDFromSS58(address string) types.AccountID {
	res, err := NewAccountIDFromSS58(address)
	if err != nil {
		panic(err)
	}
	return res
}

// MustDecodeAccountIDFromSS58 if error panic
func MustDecodeAccountIDFromSS58(address string) types.AccountID {
	res, err := DecodeAccountIDFromSS58(address)
	if err != nil {
		panic(err)
	}

	return res
}

func Hash256(buf []byte) []byte {
	hash, err := blake2b.New256([]byte{})
	if err != nil {
		panic(errors.Wrap(err, "invalid blake2b"))
	}

	_, err = hash.Write(buf)
	if err != nil {
		panic(errors.Wrap(err, "invalid blake2b write"))
	}

	h := hash.Sum(nil)
	return h
}

func Sum256(buf []byte) [32]byte {
	return blake2b.Sum256(buf)
}
