package trait_erc20

import (
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
)

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

func (a *API) Transfer(ctx Context, to AccountID, amt U128) (Hash, error) {
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
		[]string{"transfer"},
		toParam, valueParam,
	)
}

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

func (a *API) BalanceOf(ctx Context, owner AccountID) (U128, error) {
	req := struct {
		Address types.AccountID
	}{
		Address: owner,
	}

	var res types.U128

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"balance_of"},
		req,
	)

	return res, err
}

func (a *API) TotalSupply(ctx Context) (U128, error) {
	var res types.U128

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"total_supply"},
	)

	return res, err
}
