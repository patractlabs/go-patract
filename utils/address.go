package utils

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/patractlabs/go-patract/types"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
)

const (
	base58VersionPrefix = 0x2a // 42
	addressLength       = 32 + 1 + 2
)

var (
	ss58Prefix = []byte("SS58PRE")
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

// DecodeAccountIDFromSS58 encode address SS58 to AccountID
func DecodeAccountIDFromSS58(address string) (types.AccountID, error) {
	a := base58.Decode(address)

	if len(a) == 0 {
		return types.AccountID{}, errors.New("no address bytes encode")
	}

	if a[0] == base58VersionPrefix {
		if len(a) == addressLength {
			addr := a[:addressLength-2]

			hash, err := blake2b.New512([]byte{})
			if err != nil {
				return types.AccountID{}, errors.Wrap(err, "invalid blake2b")
			}

			buf := make([]byte, 0, len(ss58Prefix)+len(addr)+1)
			buf = append(buf, ss58Prefix...)
			buf = append(buf, addr...)

			_, err = hash.Write(buf)
			if err != nil {
				return types.AccountID{}, errors.Wrap(err, "invalid blake2b write")
			}

			h := hash.Sum(nil)

			if (a[addressLength-2] == h[0]) && (a[addressLength-1] == h[1]) {
				return types.NewAccountID(a[1:]), nil
			}

			return types.AccountID{},
				errors.Errorf("invalid checksum %x%x, expected %x%x",
					a[addressLength-2], a[addressLength-1],
					h[0], h[1])
		}

		return types.AccountID{}, errors.New("invalid length")
	}

	return types.AccountID{}, errors.New("invalid version")
}

// MustDecodeAccountIDFromSS58 if error panic
func MustDecodeAccountIDFromSS58(address string) types.AccountID {
	res, err := DecodeAccountIDFromSS58(address)
	if err != nil {
		panic(err)
	}

	return res
}

// EncodeAccountIDToSS58 encode accountID to ss58 format
func EncodeAccountIDToSS58(account types.AccountID) (string, error) {
	bz := make([]byte, 0, len(account)+1)
	bz = append(bz, base58VersionPrefix)
	bz = append(bz, account[:]...)

	buf := make([]byte, 0, len(ss58Prefix)+len(bz)+1)
	buf = append(buf, ss58Prefix...)
	buf = append(buf, bz...)

	hash, err := blake2b.New512([]byte{})
	if err != nil {
		return "", errors.Wrap(err, "invalid blake2b")
	}

	_, err = hash.Write(buf)
	if err != nil {
		return "", errors.Wrap(err, "invalid blake2b write")
	}

	h := hash.Sum(nil)

	complete := make([]byte, 0, addressLength+1)
	complete = append(complete, bz...)
	complete = append(complete, h[0], h[1])

	return base58.Encode(complete), nil
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
