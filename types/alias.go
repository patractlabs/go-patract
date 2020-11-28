package types

import "github.com/centrifuge/go-substrate-rpc-client/types"

var (
	// CreateStorageKey uses the given metadata and to derive the right hashing of method, prefix as well as arguments to
	// create a hashed StorageKey
	CreateStorageKey = types.CreateStorageKey

	// HexDecodeString decodes bytes from a hex string. Contrary to hex.DecodeString, this function does not error if "0x"
	// is prefixed, and adds an extra 0 if the hex string has an odd length.
	HexDecodeString = types.HexDecodeString

	// MustHexDecodeString panics if str cannot be decoded
	MustHexDecodeString = types.MustHexDecodeString
)
