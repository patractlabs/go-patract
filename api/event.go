package api

import (
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/pkg/errors"
)

// QueryEventByBlockHash get event by block hash
func (c *Client) QueryEventByBlockHash(hash types.Hash) (*types.EventRecords, error) {
	meta, err := c.API().RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, errors.Wrap(err, "get metadata latest")
	}

	key, err := types.CreateStorageKey(meta, "System", "Events", nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "create storage key")
	}

	raw, err := c.API().RPC.State.GetStorageRaw(key, hash)

	res := &types.EventRecords{}
	err = types.EventRecordsRaw(*raw).DecodeEventRecords(meta, res)
	if err != nil {
		return nil, errors.Wrap(err, "decode event records")
	}

	return res, nil
}

// QueryEventByBlockNumber get event by block number
func (c *Client) QueryEventByBlockNumber(blockNumber uint64) (*types.EventRecords, error) {
	hash, err := c.API().RPC.Chain.GetBlockHash(blockNumber)
	if err != nil {
		return nil, errors.Wrap(err, "get block hash")
	}

	return c.QueryEventByBlockHash(hash)
}
