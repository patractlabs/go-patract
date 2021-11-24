package delegator

import (
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
)

func (a *API) Get(ctx Context) (I32, error) {
	var res I32

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"get"},
	)
	return res, err
}

func (a *API) Change(ctx Context, by I32) (Hash, error) {
	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"change"},
		by,
	)
}

func (a *API) Switch(ctx Context) (Hash, error) {
	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"switch"},
	)
}
