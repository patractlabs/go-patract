package erc721

import "github.com/patractlabs/go-patract/types"

func (a *API) BalanceOf(ctx Context, owner AccountID) (TokenId, error) {
	ownerParam := struct {
		Address AccountID
	}{
		Address: owner,
	}

	var res TokenId

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"balance_of"},
		ownerParam,
	)

	return res, err
}

func (a *API) OwnerOf(ctx Context, id TokenId) (types.OptionAccountID, error) {
	var res types.OptionAccountID

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"owner_of"},
		id,
	)

	return res, err
}
