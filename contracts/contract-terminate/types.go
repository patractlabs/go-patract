package contract_terminate

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
)

// API for contract-terminate
type API struct {
	*rpc.Contract

	ContractAccountID types.AccountID
}

// New creates a new API for contract-terminate
func New(a *rpc.Contract, contractAccountID AccountID) *API {
	return &API{
		Contract:          a,
		ContractAccountID: contractAccountID,
	}
}
