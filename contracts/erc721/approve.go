package erc721

import (
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
)

func (a *API) Approve(ctx Context, spender AccountID, id TokenId) (Hash, error) {
	spenderParam := struct {
		Address AccountID
	}{
		Address: spender,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"approve"},
		spenderParam, id,
	)
}

func (a *API) GetApproved(ctx Context, id TokenId) (types.OptionAccountID, error) {
	var res types.OptionAccountID

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"get_approved"},
		id,
	)
	return res, err
}

func (a *API) IsApprovedForAll(ctx Context, owner, operator AccountID) (Bool, error) {
	ownerParam := struct {
		Address AccountID
	}{
		Address: owner,
	}

	operatorParam := struct {
		Address AccountID
	}{
		Address: operator,
	}

	var res Bool

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"is_approved_for_all"},
		ownerParam, operatorParam,
	)

	return res, err
}

func (a *API) SetApprovalForAll(ctx Context, to AccountID, approved Bool) (Hash, error) {
	toParam := struct {
		Address AccountID
	}{
		Address: to,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"set_approval_for_all"},
		toParam, approved,
	)
}
