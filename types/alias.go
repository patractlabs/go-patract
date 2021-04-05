package types

import "github.com/centrifuge/go-substrate-rpc-client/v2/types"

const (
	ExtrinsicVersion4 = types.ExtrinsicVersion4
)

var (
	EncodeToBytes = types.EncodeToBytes
)

type (
	U32 = types.U32
	U64 = types.U64
	// U128 is an unsigned 128-bit integer, it is represented as a big.Int in Go.
	U128 = types.U128
)

var (
	NewU32 = types.NewU32
	NewU64 = types.NewU64
	// NewU128 creates a new U128 type
	NewU128 = types.NewU128
)

type (
	// Hash is the default hash that is used across the system. It is just a thin wrapper around H256
	Hash = types.Hash
)

var (
	// NewHash creates a new Hash type
	NewHash = types.NewHash
	// NewHashFromHexString creates a new Hash type from a hex string
	NewHashFromHexString = types.NewHashFromHexString
)

// Address is a wrapper around an AccountId or an AccountIndex. It is encoded with a prefix in case of an AccountID.
// Basically the Address is encoded as `[ <prefix-byte>, ...publicKey/...bytes ]` as per spec
type Address = types.Address

var (
	// NewAddressFromHexAccountID creates an Address from the given hex string that contains an AccountID (public key)
	NewAddressFromHexAccountID = types.NewAddressFromHexAccountID

	// NewAddressFromAccountID creates an Address from the given AccountID (public key)
	NewAddressFromAccountID = types.NewAddressFromAccountID
)

// MustAddressFromHexAccount address from hex account, panic if invalid
func MustAddressFromHexAccount(str string) Address {
	res, err := NewAddressFromHexAccountID(str)
	if err != nil {
		panic(err)
	}

	return res
}

type (
	// AccountID represents a public key (an 32 byte array)
	AccountID = types.AccountID

	// Bytes represents byte slices. Bytes has a variable length, it is encoded with a scale prefix
	Bytes = types.Bytes
)

var (
	// NewAccountID creates a new AccountID type
	NewAccountID = types.NewAccountID
)

var (
	// CreateStorageKey uses the given metadata and to derive the right hashing of method, prefix as well as arguments to
	// create a hashed StorageKey
	CreateStorageKey = types.CreateStorageKey

	// HexDecodeString decodes bytes from a hex string. Contrary to hex.DecodeString, this function does not error if "0x"
	// is prefixed, and adds an extra 0 if the hex string has an odd length.
	HexDecodeString = types.HexDecodeString

	// MustHexDecodeString panics if str cannot be decoded
	MustHexDecodeString = types.MustHexDecodeString

	// HexEncodeToString HexEncode encodes bytes to a hex string.
	// Contrary to hex.EncodeToString, this function prefixes the hex string with "0x"
	HexEncodeToString = types.HexEncodeToString
)

type (
	AccountInfo      = types.AccountInfo
	SignatureOptions = types.SignatureOptions
	ExtrinsicEra     = types.ExtrinsicEra
	Metadata         = types.Metadata
	Extrinsic        = types.Extrinsic
)

var (
	NewExtrinsic = types.NewExtrinsic
	NewCall      = types.NewCall
)
