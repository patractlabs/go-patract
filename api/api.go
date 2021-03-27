package api

import (
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v2"
	"github.com/patractlabs/go-patract/types"
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

// API returns the rpcAPI for client
func (c Client) API() *gsrpc.SubstrateAPI {
	return c.rpcAPI
}

// WithLogger set logger
func (c *Client) WithLogger(logger log.Logger) {
	c.logger = logger
}

// Call call rpc
func (c *Client) Call(result interface{}, method string, args ...interface{}) error {
	c.logger.Debug("Call", "method", method)
	return c.rpcAPI.Client.Call(result, method, args...)
}

func MakeExtrinisic(
	nonce uint64, meta *types.Metadata, cs *types.ChainStatus,
	call string, args ...interface{}) ([]byte, error) {
	cc, err := types.NewCall(meta, call, args...)
	if err != nil {
		return []byte{}, errors.Wrap(err, "new call error")
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(cc)

	blkHash, err := types.NewHashFromHexString(cs.BlockHash)
	if err != nil {
		return []byte{}, errors.Wrap(err, "block hash hex error")
	}

	genHash, err := types.NewHashFromHexString(cs.GenesisHash)
	if err != nil {
		return []byte{}, errors.Wrap(err, "genesis hash hex error")
	}

	o := types.SignatureOptions{
		BlockHash:          blkHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genHash,
		Nonce:              types.NewUCompactFromUInt(nonce),
		SpecVersion:        types.NewU32(cs.SpecVersion),
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: types.NewU32(cs.TransactionVersion),
	}

	return GenSignBytes(ext, o)
}

// SubmitAndWaitExtrinsic submit and wait extrinsic into chain
func (c *Client) SubmitAndWaitExtrinsic(ctx Context, call string, args ...interface{}) (types.Hash, error) {
	c.logger.Debug("submitAndWatchExtrinsic", "call", call)

	meta, err := c.rpcAPI.RPC.State.GetMetadataLatest()
	if err != nil {
		return types.Hash{}, errors.Wrap(err, "get metadata lastest error")
	}

	cc, err := types.NewCall(meta, call, args...)
	if err != nil {
		return types.Hash{}, errors.Wrap(err, "new call error")
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(cc)

	genesisHash, err := c.rpcAPI.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return types.Hash{}, errors.Wrap(err, "get block hash error")
	}

	rv, err := c.rpcAPI.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return types.Hash{}, errors.Wrap(err, "get runtime version lastest error")
	}

	authKey := ctx.From()

	key, err := types.CreateStorageKey(meta, "System", "Account", authKey.PublicKey, nil)
	if err != nil {
		return types.Hash{}, errors.Wrap(err, "create storage key error")
	}

	var accountInfo types.AccountInfo
	ok, err := c.rpcAPI.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return types.Hash{}, errors.Wrapf(err, "create storage key error by %s", authKey.Address)
	} else if !ok {
		return types.Hash{}, errors.Errorf("no accountInfo found by %s", authKey.Address)
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

	ctx.logger.Info("genesisHash", "raw", types.HexEncodeToString(genesisHash[:]))

	// Sign the transaction using Alice's default account
	if err := ext.Sign(authKey, o); err != nil {
		return types.Hash{}, errors.Wrap(err, "sign extrinsic error")
	}

	// Send the extrinsic
	sub, err := c.rpcAPI.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return types.Hash{}, errors.Wrap(err, "submit error")
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
				return status.AsInBlock, nil
			}
		case err := <-sub.Err():
			return types.Hash{}, errors.Wrap(err, "subscribe error")
		case <-ctx.Done():
			return types.Hash{}, ctx.Err()
		}
	}
}
