package types

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/centrifuge/go-substrate-rpc-client/v2/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
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

// DecodeAccountIDFromSS58 encode address SS58 to AccountID
func DecodeAccountIDFromSS58(address string) (AccountID, error) {
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

			buf := make([]byte, 0, len(ss58Prefix)+len(addr)+1)
			buf = append(buf, ss58Prefix...)
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

// EncodeAccountIDToSS58 encode accountID to ss58 format
func EncodeAccountIDToSS58(account AccountID) (string, error) {
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

type Option struct {
	hasValue bool
}

// IsNone returns true if the value is missing
func (o Option) IsNone() bool {
	return !o.hasValue
}

// IsNone returns true if a value is present
func (o Option) IsSome() bool {
	return o.hasValue
}

func (o *Option) SetHasValue(h bool) {
	o.hasValue = h
}

// OptionAccountID is a structure that can store a AccountID or a missing value
type OptionAccountID struct {
	Option
	Value AccountID `scale:"Some"`
}

// NewOptionAccountID creates an OptionAccountID with a value
func NewOptionAccountID(value AccountID) OptionAccountID {
	return OptionAccountID{Option{true}, value}
}

// NewOptionAccountIDEmpty creates an OptionAccountID without a value
func NewOptionAccountIDEmpty() OptionAccountID {
	return OptionAccountID{Option: Option{false}}
}

func (o OptionAccountID) Encode(encoder scale.Encoder) error {
	return encoder.EncodeOption(o.hasValue, o.Value)
}

func (o *OptionAccountID) Decode(decoder scale.Decoder) error {
	return decoder.DecodeOption(&o.hasValue, &o.Value)
}

// SetSome sets a value
func (o *OptionAccountID) SetSome(value AccountID) {
	o.hasValue = true
	o.Value = value
}

// SetNone removes a value and marks it as missing
func (o *OptionAccountID) SetNone() {
	o.hasValue = false
	o.Value = AccountID{}
}

// Unwrap returns a flag that indicates whether a value is present and the stored value
func (o OptionAccountID) Unwrap() (ok bool, value AccountID) {
	return o.hasValue, o.Value
}

func (o OptionAccountID) String() string {
	if o.hasValue {
		str, _ := EncodeAccountIDToSS58(o.Value)
		return str
	}

	return ""
}

// OptionValue is a structure that can store a Value or a missing value
type OptionValue struct {
	Option
	value interface{}
}

// NewOptionValue creates an OptionValue with a value
func NewOptionValue(value interface{}) OptionValue {
	return OptionValue{Option{true}, value}
}

// NewOptionValueEmpty creates an OptionValue without a value
func NewOptionValueEmpty() OptionValue {
	return OptionValue{Option: Option{false}}
}

func (o OptionValue) Encode(encoder scale.Encoder) error {
	return encoder.EncodeOption(o.hasValue, o.value)
}

func (o *OptionValue) Decode(decoder scale.Decoder) error {
	return decoder.DecodeOption(&o.hasValue, &o.value)
}

// SetSome sets a value
func (o *OptionValue) SetSome(value interface{}) {
	o.hasValue = true
	o.value = value
}

// SetNone removes a value and marks it as missing
func (o *OptionValue) SetNone() {
	o.hasValue = false
	o.value = nil
}

// Unwrap returns a flag that indicates whether a value is present and the stored value
func (o OptionValue) Unwrap() (ok bool, value interface{}) {
	return o.hasValue, o.value
}
