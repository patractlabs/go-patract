package trait_incrementer

import (
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
)

func (a *API) IncBy(ctx Context, delta U64) (Hash, error) {
	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"inc_by"},
		delta,
	)
}

func (a *API) Inc(ctx Context) (Hash, error) {
	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"Increment", "inc"},
	)
}

func (a *API) Reset(ctx Context) (Hash, error) {
	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"Reset", "reset"},
	)
}

func (a *API) Get(ctx Context) (U64, error) {
	var res U64

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"Increment", "get"},
	)
	return res, err
}
