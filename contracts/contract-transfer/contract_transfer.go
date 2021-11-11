package contract_transfer

/**
合约转让
*/

import (
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
)

// GiveMe 获取指定的value值，但是需要初始的gas费用。
func (a *API) GiveMe(ctx Context, value Balance) (Hash, error) {
	valueParam := struct {
		Value Balance
	}{
		Value: value,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"give_me"},
		valueParam,
	)
}

//
//func (a *API) WasItTen(ctx Context) {
//	var res interface{}
//
//	a.CallToRead(ctx,
//		&res,
//		a.ContractAccountID,
//		[]string{"was_it_ten"},
//	)
//}
