package main

import (
	"fmt"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/config"
)

func main() {
	// Instantiate the API
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	hash, err := api.RPC.Chain.GetFinalizedHead()
	if err != nil {
		panic(err)
	}

	blk, err := api.RPC.Chain.GetBlock(hash)
	if err != nil {
		panic(err)
	}

	fmt.Printf("blk %v\n", blk)
	for _, h := range blk.Block.Extrinsics {
		fmt.Printf("e %v\n", h)
	}

	queryEvent()
}
