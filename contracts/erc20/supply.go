package erc20

import "github.com/patractlabs/go-patract/types"

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
