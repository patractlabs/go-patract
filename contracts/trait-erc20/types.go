package trait_erc20

import (
	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/types"
)

type (
	Context = api.Context
)

type (
	AccountID    = types.AccountID
	Balance      = types.Balance
	U128         = types.U128
	Hash         = types.Hash
	U32          = types.U32
	VecU8        = []types.U8
	VecBalance   = []types.Balance
	VecAccountID = []types.AccountID
	Bool         = types.Bool
)

// API for trait-erc20
type API struct {
	*rpc.Contract

	ContractAccountID types.AccountID
}

// New creates a new API for trait-erc20
func New(a *rpc.Contract, contractAccountID AccountID) *API {
	return &API{
		Contract:          a,
		ContractAccountID: contractAccountID,
	}
}
