package flipper

import (
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
)

func (a *API) Flip(ctx Context) (Hash, error) {
	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"flip"},
	)
}

func (a *API) Get(ctx Context) (Bool, error) {
	var res Bool

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"get"},
	)
	return res, err
}
