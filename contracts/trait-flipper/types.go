package trait_flipper

import (
	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/types"
)

type (
	Context = api.Context
)

type (
	AccountID = types.AccountID
	Hash      = types.Hash
	Bool      = types.Bool
)

// API for trait-flipper
type API struct {
	*rpc.Contract

	ContractAccountID types.AccountID
}

// New creates a new API for trait-flipper
func New(a *rpc.Contract, contractAccountID AccountID) *API {
	return &API{
		Contract:          a,
		ContractAccountID: contractAccountID,
	}
}
