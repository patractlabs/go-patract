package erc721

import (
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
)

// Transfer
func (a *API) Transfer(ctx Context, destination AccountID, id U32) (Hash, error) {
	destinationParam := struct {
		Address AccountID
	}{
		Address: destination,
	}

	idParam := struct {
		Id U32
	}{
		Id: id,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"transfer"},
		destinationParam, idParam,
	)
}

// TransferFrom
func (a *API) TransferFrom(ctx Context, from, to AccountID, amt U128) (Hash, error) {
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

	valueParam := struct {
		Value U128
	}{
		Value: amt,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"transfer_from"},
		fromParam, toParam, valueParam,
	)
}
