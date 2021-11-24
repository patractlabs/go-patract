package rand_extension

import (
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
)

func (a *API) Update(ctx Context) (Hash, error) {
	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"update"},
	)
}

func (a *API) Get(ctx Context) (VecU8, error) {
	var res VecU8

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"get"},
	)
	return res, err
}
