package erc721

import "github.com/patractlabs/go-patract/types"

func (a *API) BalanceOf(ctx Context, owner AccountID) (U32, error) {
	ownerParam := struct {
		Address AccountID
	}{
		Address: owner,
	}

	var res U32

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"balance_of"},
		ownerParam,
	)

	return res, err
}

func (a *API) OwnerOf(ctx Context, id U32) (types.OptionAccountID, error) {
	idParam := struct {
		Id U32
	}{
		Id: id,
	}

	var res types.OptionAccountID

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"owner_of"},
		idParam,
	)

	return res, err
}
