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

		initERC20(t, logger, env)

		metaBz, err := ioutil.ReadFile(erc20MetaPath)
		require.Nil(err)

		api, err := rpc.NewContractAPI(env.URL())
		require.Nil(err)

		api.WithLogger(logger)
		err = api.WithMetaData(metaBz)
		require.Nil(err)

		ctx := rpc.NewCtx(context.Background()).WithFrom(authKey)

		var initSupply uint64 = 100000000000000

		// Instantiate
		_, contractAccount, err := api.Instantiate(ctx,
			types.NewCompactBalance(initSupply),
			types.NewCompactGas(test.DefaultGas),
			contracts.CodeHashERC20,
			types.NewU128(totalSupply),
		)
		require.Nil(err)

		t.Logf("constract %s", contractAccount)

		req := struct {
			Address types.AccountID
		}{
			Address: utils.MustAccountIDFromSS58(authKey.Address),
		}
		var res types.U128

		// Instantiate
		err = api.Call(ctx,
			&res,
			contractAccount,
			[]string{"balance_of"},
			req,
		)

		require.Nil(err)
		t.Logf("res %v", res)

		// transfer
	})
}
