package incrementer

import (
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
)

func (a *API) Inc(ctx Context, by I32) (Hash, error) {
	byParam := struct {
		Value I32
	}{
		Value: by,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"inc"},
		byParam,
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
