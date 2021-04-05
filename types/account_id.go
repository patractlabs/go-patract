package types

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
)

const (
	base58VersionPrefix = 0x2a // 42
	addressLength       = 32 + 1 + 2
)

var (
	DefaultSS58Prefix = []byte("SS58PRE")
)

type SS58Codec struct {
	prefix []byte
}

func NewSS58Codec(prefix []byte) *SS58Codec {
	p := make([]byte, 0, len(prefix))
	copy(p, prefix)

	return &SS58Codec{
		prefix: p,
	}
}

var (
	defaultSS58Codec = NewSS58Codec(DefaultSS58Prefix)
)

// GetPrefix get ss58 prefix
func (c SS58Codec) GetPrefix() []byte {
	return c.prefix
}

// DecodeAccountID decode address SS58 to AccountID
func (c SS58Codec) DecodeAccountID(address string) (AccountID, error) {
	a := base58.Decode(address)

	if len(a) == 0 {
		return AccountID{}, errors.New("no address bytes encode")
	}

	if a[0] == base58VersionPrefix {
		if len(a) == addressLength {
			addr := a[:addressLength-2]

			hash, err := blake2b.New512([]byte{})
			if err != nil {
				return AccountID{}, errors.Wrap(err, "invalid blake2b")
			}

			buf := make([]byte, 0, len(c.prefix)+len(addr)+1)
			buf = append(buf, c.prefix...)
			buf = append(buf, addr...)

			_, err = hash.Write(buf)
			if err != nil {
				return AccountID{}, errors.Wrap(err, "invalid blake2b write")
			}

			h := hash.Sum(nil)

			if (a[addressLength-2] == h[0]) && (a[addressLength-1] == h[1]) {
				return types.NewAccountID(a[1:]), nil
			}

			return AccountID{},
				errors.Errorf("invalid checksum %x%x, expected %x%x",
					a[addressLength-2], a[addressLength-1],
					h[0], h[1])
		}

		return AccountID{}, errors.New("invalid length")
	}

	return AccountID{}, errors.New("invalid version")
}

func (c SS58Codec) EncodeAccountID(id AccountID) (string, error) {
	bz := make([]byte, 0, len(id)+1)
	bz = append(bz, base58VersionPrefix)
	bz = append(bz, id[:]...)

	buf := make([]byte, 0, len(c.prefix)+len(bz)+1)
	buf = append(buf, c.prefix...)
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

// DecodeAccountIDFromSS58 encode address SS58 to AccountID
func DecodeAccountIDFromSS58(address string) (AccountID, error) {
	return defaultSS58Codec.DecodeAccountID(address)
}

// EncodeAccountIDToSS58 encode accountID to ss58 format
func EncodeAccountIDToSS58(account AccountID) (string, error) {
	return defaultSS58Codec.EncodeAccountID(account)
}
