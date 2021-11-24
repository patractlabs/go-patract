package erc721

import (
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
)

// Transfer
func (a *API) Transfer(ctx Context, destination AccountID, id TokenId) (Hash, error) {
	destinationParam := struct {
		Address AccountID
	}{
		Address: destination,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"transfer"},
		destinationParam, id,
	)
}

// TransferFrom
func (a *API) TransferFrom(ctx Context, from, to AccountID, id TokenId) (Hash, error) {
	fromParam := struct {
		Address AccountID
	}{
		Address: from,
	}

	toParam := struct {
		Address AccountID
	}{
		Address: to,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"transfer_from"},
		fromParam, toParam, id,
	)
}

func (a *API) Mint(ctx Context, id TokenId) (Hash, error) {
	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"mint"},
		id,
	)
}

func (a *API) Burn(ctx Context, id TokenId) (Hash, error) {
	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"burn"},
		id,
	)
}
