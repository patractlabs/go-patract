package rpc_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/test/contracts"
	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

func TestDeployAndCallERC20(t *testing.T) {
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

		// Instantiate
		hash, contractAccount, err := api.Instantiate(ctx,
			types.NewCompactBalance(10000000000000000),
			types.NewCompactGas(test.DefaultGas),
			contracts.CodeHashERC20,
			types.NewU128(totalSupply),
		)

		require.Nil(err)
		t.Logf("instantiate hash %v %v", hash, contractAccount)
	})
}
