package erc721

import (
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
)

func (a *API) Approve(ctx Context, spender AccountID, id U32) (Hash, error) {
	spenderParam := struct {
		Address AccountID
	}{
		Address: spender,
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
		[]string{"approve"},
		spenderParam, idParam,
	)
}

func (a *API) GetApproved(ctx Context, id U32) (types.OptionAccountID, error) {
	idParam := struct {
		Id U32
	}{
		Id: id,
	}

	var res types.OptionAccountID

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"get_approved"},
		idParam,
	)
	return res, err
}

func (a *API) IsApprovedForAll(ctx Context, owner, operator AccountID) (bool, error) {
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

	var res bool

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"is_approved_for_all"},
		ownerParam, operatorParam,
	)

	return res, err
}

func (a *API) SetApprovalForAll(ctx Context, to AccountID, approved bool) (Hash, error) {
	toParam := struct {
		Address AccountID
	}{
		Address: to,
	}

	approvedParam := struct {
		Approved bool
	}{
		Approved: approved,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"set_approval_for_all"},
		toParam, approvedParam,
	)
}
