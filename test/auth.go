package test

import (
	"github.com/centrifuge/go-substrate-rpc-client/signature"
	"github.com/patractlabs/go-patract/utils"
)

var (
	AliceAccountID = utils.MustDecodeAccountIDFromSS58(signature.TestKeyringPairAlice.Address)
)
