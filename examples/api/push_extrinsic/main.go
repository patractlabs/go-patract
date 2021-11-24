package main

import (
	"fmt"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v3"
	"github.com/centrifuge/go-substrate-rpc-client/v3/config"
	"github.com/centrifuge/go-substrate-rpc-client/v3/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
)

func main() {
	// Instantiate the API
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	// Create a call, transferring 12345 units to Bob
	bob, err := types.NewAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	if err != nil {
		panic(err)
	}

	var (
		authKey        = signature.TestKeyringPairAlice
		to             = bob
		amount  uint64 = 12345
	)

	if err := pushTransfer(api, authKey, to, amount); err != nil {
		fmt.Printf("push transfer error by %s", err)
	}
}
