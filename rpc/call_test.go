package rpc_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/test/contracts"
	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

func TestCallERC20(t *testing.T) {
	test.ByCanvasEnv(t, func(logger log.Logger, env test.Env) {
		require := require.New(t)

		if !env.IsUseExtToTest() {
			initERC20(t, logger, env)
		}

		metaBz, err := ioutil.ReadFile(erc20MetaPath)
		require.Nil(err)

		api, err := rpc.NewContractAPI(env.URL())
		require.Nil(err)

		api.WithLogger(logger)
		err = api.WithMetaData(metaBz)
		require.Nil(err)

		ctx := rpc.NewCtx(context.Background()).WithFrom(authKey)

		var initSupply uint64 = 100000000000000

		var contractAccount types.AccountID

		if !env.IsUseExtToTest() {
			// Instantiate
			_, contractAccount, err = api.Instantiate(ctx,
				types.NewCompactBalance(initSupply),
				types.NewCompactGas(test.DefaultGas),
				contracts.CodeHashERC20,
				types.NewU128(totalSupply),
			)
			require.Nil(err)

			t.Logf("constract %s", contractAccount)
		} else {
			contractAccount = utils.MustDecodeAccountIDFromSS58("5HKinTRKW9THEJxbQb22Nfyq9FPWNVZ9DQ2GEQ4Vg1LqTPuk")
		}

		req := struct {
			Address types.AccountID
		}{
			Address: utils.MustAccountIDFromSS58(authKey.Address),
		}
		var res types.U128

		// Instantiate
		err = api.CallToRead(ctx,
			&res,
			contractAccount,
			[]string{"balance_of"},
			req,
		)
		require.Nil(err)
		t.Logf("res %v", res)

		// transfer
		to := struct {
			Address types.AccountID
		}{
			Address: bob,
		}

		value := struct {
			Value types.U128
		}{
			Value: types.NewBalanceByU64(1),
		}

		hash, err := api.CallToExec(ctx,
			contractAccount,
			types.NewCompactBalance(0),
			types.NewCompactGas(test.DefaultGas),
			[]string{"transfer"},
			to, value,
		)
		require.Nil(err)
		t.Logf("transfer hash %v", hash)

		{
			req := struct {
				Address types.AccountID
			}{
				Address: bob,
			}
			var res types.U128

			err = api.CallToRead(ctx,
				&res,
				contractAccount,
				[]string{"balance_of"},
				req,
			)
			require.Nil(err)
			t.Logf("res %v", res)

			require.Equal(res.Int.Int64(), value.Value.Int.Int64())
		}
	})
}
