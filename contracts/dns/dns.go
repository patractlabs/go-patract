package dns

import (
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
)

func (a *API) Register(ctx Context, name Hash) (Hash, error) {
	nameParam := struct {
		Name Hash
	}{
		Name: name,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"register"},
		nameParam,
	)
}

func (a *API) SetAddress(ctx Context, name Hash, newAddress AccountID) (Hash, error) {
	nameParam := struct {
		Name Hash
	}{
		Name: name,
	}

	NewAddressParam := struct {
		NewAddress AccountID
	}{
		NewAddress: newAddress,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"set_address"},
		nameParam, NewAddressParam,
	)
}

func (a *API) Transfer(ctx Context, name Hash, to AccountID) (Hash, error) {
	nameParam := struct {
		Name Hash
	}{
		Name: name,
	}

	toParam := struct {
		Address AccountID
	}{
		Address: to,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"transfer"},
		nameParam, toParam,
	)
}

func (a *API) GetAddress(ctx Context, name Hash) (AccountID, error) {
	nameParam := struct {
		Name Hash
	}{
		Name: name,
	}

	var res AccountID

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"get_address"},
		nameParam,
	)
	return res, err
}
