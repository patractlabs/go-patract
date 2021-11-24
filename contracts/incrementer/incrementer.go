package incrementer

import (
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
)

func (a *API) Inc(ctx Context, by I32) (Hash, error) {
	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"inc"},
		by,
	)
}

func (a *API) Get(ctx Context) (I32, error) {
	var res I32

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"get"},
	)
	return res, err
}
