package main

import (
	"fmt"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v3"
	"github.com/centrifuge/go-substrate-rpc-client/v3/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/pkg/errors"
)

func pushTransfer(api *gsrpc.SubstrateAPI, authKey signature.KeyringPair, to types.Address, amount uint64) error {
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return errors.Wrap(err, "get metadata lastest error")
	}

	c, err := types.NewCall(meta, "Balances.transfer", to, types.NewUCompactFromUInt(amount))
	if err != nil {
		return errors.Wrap(err, "new call error")
	}

	return pushExtrinsic(api, authKey, meta, c)
}

func pushExtrinsic(api *gsrpc.SubstrateAPI, authKey signature.KeyringPair, meta *types.Metadata, c types.Call) error {
	// Create the extrinsic
	ext := types.NewExtrinsic(c)

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return errors.Wrap(err, "get block hash error")
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return errors.Wrap(err, "get runtime version lastest error")
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", authKey.PublicKey, nil)
	if err != nil {
		return errors.Wrap(err, "create storage key error")
	}

	var accountInfo types.AccountInfo
	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return errors.Wrapf(err, "create storage key error by %s", authKey.Address)
	} else if !ok {
		return errors.Errorf("no accountInfo found by %s", authKey.Address)
	}

	nonce := uint32(accountInfo.Nonce)

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction using Alice's default account
	err = ext.Sign(authKey, o)
	if err != nil {
		return errors.Wrap(err, "sign extrinsic error")
	}

	// Send the extrinsic
	sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return errors.Wrap(err, "submit error")
	}

	defer sub.Unsubscribe()

	for {
		status := <-sub.Chan()
		fmt.Printf("Transaction status: %#v\n", status)

		if status.IsInBlock {
			fmt.Printf("Completed at block hash: %#x\n", status.AsInBlock)
			return nil
		}
	}
}
