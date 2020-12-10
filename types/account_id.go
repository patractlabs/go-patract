package types

import (
	"github.com/centrifuge/go-substrate-rpc-client/scale"
)

type option struct {
	hasValue bool
}

// IsNone returns true if the value is missing
func (o option) IsNone() bool {
	return !o.hasValue
}

// IsNone returns true if a value is present
func (o option) IsSome() bool {
	return o.hasValue
}

// OptionAccountID is a structure that can store a AccountID or a missing value
type OptionAccountID struct {
	option
	value AccountID
}

// NewOptionAccountID creates an OptionAccountID with a value
func NewOptionAccountID(value AccountID) OptionAccountID {
	return OptionAccountID{option{true}, value}
}

// NewOptionAccountIDEmpty creates an OptionAccountID without a value
func NewOptionAccountIDEmpty() OptionAccountID {
	return OptionAccountID{option: option{false}}
}

func (o OptionAccountID) Encode(encoder scale.Encoder) error {
	return encoder.EncodeOption(o.hasValue, o.value)
}

func (o *OptionAccountID) Decode(decoder scale.Decoder) error {
	return decoder.DecodeOption(&o.hasValue, &o.value)
}

// SetSome sets a value
func (o *OptionAccountID) SetSome(value AccountID) {
	o.hasValue = true
	o.value = value
}

// SetNone removes a value and marks it as missing
func (o *OptionAccountID) SetNone() {
	o.hasValue = false
	o.value = AccountID{}
}

// Unwrap returns a flag that indicates whether a value is present and the stored value
func (o OptionAccountID) Unwrap() (ok bool, value AccountID) {
	return o.hasValue, o.value
}

// OptionValue is a structure that can store a Value or a missing value
type OptionValue struct {
	option
	value interface{}
}

// NewOptionValue creates an OptionValue with a value
func NewOptionValue(value interface{}) OptionValue {
	return OptionValue{option{true}, value}
}

// NewOptionValueEmpty creates an OptionValue without a value
func NewOptionValueEmpty() OptionValue {
	return OptionValue{option: option{false}}
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
