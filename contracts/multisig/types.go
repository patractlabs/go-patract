package multisig

import (
	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/types"
)

type (
	Context = api.Context
)

type (
	AccountID     = types.AccountID
	Balance       = types.Balance
	TransactionId = types.U32
)

// API for erc20
type API struct {
	*rpc.Contract

	ContractAccountID types.AccountID
}

// New creates a new API for erc20
func New(a *rpc.Contract, contractAccountID AccountID) *API {
	return &API{
		Contract:          a,
		ContractAccountID: contractAccountID,
	}
}
