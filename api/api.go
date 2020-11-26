package api

import (
	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/pkg/errors"
)

// Client a high-level chain-sdk warpper to interaction to chain node
type Client struct {
	rpcAPI *gsrpc.SubstrateAPI
	logger log.Logger
}

// NewClient create a client to chain
func NewClient(logger log.Logger, url string) (*Client, error) {
	api, err := gsrpc.NewSubstrateAPI(url)
	if err != nil {
		return nil, err
	}

	return &Client{
		rpcAPI: api,
		logger: logger,
	}, nil
}

// SubmitAndWaitExtrinsic submit and wait extrinsic into chain
func (c *Client) SubmitAndWaitExtrinsic(ctx Context, call string, args ...interface{}) (string, error) {
	c.logger.Debug("submitAndWatchExtrinsic", "call", call)

	meta, err := c.rpcAPI.RPC.State.GetMetadataLatest()
	if err != nil {
		return "", errors.Wrap(err, "get metadata lastest error")
	}

	cc, err := types.NewCall(meta, call, args...)
	if err != nil {
		return "", errors.Wrap(err, "new call error")
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(cc)

	genesisHash, err := c.rpcAPI.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return "", errors.Wrap(err, "get block hash error")
	}

	rv, err := c.rpcAPI.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return "", errors.Wrap(err, "get runtime version lastest error")
	}

	authKey := ctx.From()

	key, err := types.CreateStorageKey(meta, "System", "Account", authKey.PublicKey, nil)
	if err != nil {
		return "", errors.Wrap(err, "create storage key error")
	}

	var accountInfo types.AccountInfo
	ok, err := c.rpcAPI.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return "", errors.Wrapf(err, "create storage key error by %s", authKey.Address)
	} else if !ok {
		return "", errors.Errorf("no accountInfo found by %s", authKey.Address)
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
	if err := ext.Sign(authKey, o); err != nil {
		return "", errors.Wrap(err, "sign extrinsic error")
	}

	// Send the extrinsic
	sub, err := c.rpcAPI.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return "", errors.Wrap(err, "submit error")
	}
	defer sub.Unsubscribe()

	c.logger.Debug("start watch extrinsic", "call", call)
	for {
		select {
		case status := <-sub.Chan():
			c.logger.Debug("on status", "call", call,
				"isInBlock", status.IsInBlock, "isFinalized", status.IsFinalized)

			if status.IsInBlock {
				c.logger.Debug("Completed at block", "hash", status.AsInBlock.Hex())
				// if is in block, should return
				return status.AsInBlock.Hex(), nil
			}
		case err := <-sub.Err():
			return "", errors.Wrap(err, "subscribe error")
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}
