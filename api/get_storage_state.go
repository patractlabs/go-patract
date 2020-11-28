package api

import (
	"github.com/patractlabs/go-patract/types"
	"github.com/pkg/errors"
)

var (
	// ErrGetStorageNotFound error if no storage is available
	ErrGetStorageNotFound = errors.New("storage key not found")
)

// GetStorageLatest get storage in latest state from chain
func (c *Client) GetStorageLatest(res interface{}, prefix, method string, arg []byte, arg2 []byte) error {
	meta, err := c.API().RPC.State.GetMetadataLatest()
	if err != nil {
		return errors.Wrap(err, "get metadata failed")
	}

	key, err := types.CreateStorageKey(meta, prefix, method, arg, arg2)
	if err != nil {
		return errors.Wrap(err, "create storage key failed")
	}

	ok, err := c.API().RPC.State.GetStorageLatest(key, res)
	if err != nil {
		return errors.Wrap(err, "get storage failed")
	}

	if !ok {
		return ErrGetStorageNotFound
	}

	return nil
}
