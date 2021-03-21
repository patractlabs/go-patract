package api

import (
	"context"

	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
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
	if err != nil {
		return nil, errors.Wrap(err, "get storage raw")
	}

	res := &types.EventRecords{}
	err = types.EventRecordsRaw(*raw).DecodeEventRecords(meta, res)
	if err != nil {
		return nil, errors.Wrap(err, "decode event records")
	}

	return res, nil
}

func (c *Client) WatchEvents(ctx context.Context, h EventHandler) error {
	meta, err := c.API().RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Subscribe to system events via storage
	key, err := types.CreateStorageKey(meta, "System", "Events", nil, nil)
	if err != nil {
		panic(err)
	}

	sub, err := c.API().RPC.State.SubscribeStorageRaw([]types.StorageKey{key})
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	// outer for loop for subscription notifications
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("watcher events stoped")
			return nil
		case set := <-sub.Chan():
			// inner loop for the changes within one of those notifications
			for _, chng := range set.Changes {
				if !types.Eq(chng.StorageKey, key) || !chng.HasStorageData {
					// skip, we are only interested in events with content
					continue
				}

				// Decode the event records
				evt := types.EventRecords{}
				err = types.EventRecordsRaw(chng.StorageData).DecodeEventRecords(meta, &evt)
				if err != nil {
					return errors.Wrapf(err, "decode event records error %s", set.Block.Hex())
				}

				blk, err := c.API().RPC.Chain.GetHeader(set.Block)
				if err != nil {
					return errors.Wrap(err, "get block header error")
				}

				if err := h(c.logger, uint64(blk.Number), &evt); err != nil {
					return errors.Wrapf(err, "handler blk %d %s", uint64(blk.Number), set.Block.Hex())
				}
			}
		}
	}
}

// QueryEventByBlockNumber get event by block number
func (c *Client) QueryEventByBlockNumber(blockNumber uint64) (*types.EventRecords, error) {
	hash, err := c.API().RPC.Chain.GetBlockHash(blockNumber)
	if err != nil {
		return nil, errors.Wrap(err, "get block hash")
	}

	return c.QueryEventByBlockHash(hash)
}
