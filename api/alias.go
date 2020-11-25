package api

import "github.com/centrifuge/go-substrate-rpc-client/types"

type (
	// Address is a wrapper around an AccountId or an AccountIndex
	Address = types.Address
)

var (
	// NewAddressFromHexAccountID creates an Address from the given hex string that contains an AccountID (public key)
	NewAddressFromHexAccountID = types.NewAddressFromHexAccountID
)

// MustAddressFromHexAccount address from hex account, panic if invalid
func MustAddressFromHexAccount(str string) Address {
	res, err := NewAddressFromHexAccountID(str)
	if err != nil {
		panic(err)
	}

	return res
}
