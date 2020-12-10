package erc20

import (
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
)

// Transfer
func (a *API) Approve(ctx Context, spender AccountID, value U128) (Hash, error) {
	spenderParam := struct {
		Address AccountID
	}{
		Address: spender,
	}

	valueParam := struct {
		Value U128
	}{
		Value: value,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"approve"},
		spenderParam, valueParam,
	)
}

func (a *API) Allowance(ctx Context, owner, spender AccountID) (U128, error) {
	ownerParam := struct {
		Address AccountID
	}{
		Address: owner,
	}

	spenderParam := struct {
		Address AccountID
	}{
		Address: spender,
	}

	var res U128

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"approve"},
		ownerParam, spenderParam,
	)

	return res, err
}
