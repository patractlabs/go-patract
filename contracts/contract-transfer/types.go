package contract_transfer

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
	Balance   = types.Balance
	Hash      = types.Hash
)

// API for contract-transfer
type API struct {
	*rpc.Contract

	ContractAccountID types.AccountID
}

// New creates a new API for contract-transfer
func New(a *rpc.Contract, contractAccountID AccountID) *API {
	return &API{
		Contract:          a,
		ContractAccountID: contractAccountID,
	}
}
