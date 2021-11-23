package flipper_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/signature"
	"github.com/patractlabs/go-patract/contracts/flipper"
	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

func TestFlip(t *testing.T) {
	test.ByNodeEnv(t, func(logger log.Logger, env test.Env) {
		require := require.New(t)
		contractAccountID := initFlipper(t, logger, env, signature.TestKeyringPairAlice)
		rpcAPI, err := rpc.NewContractAPI(env.URL())
		require.Nil(err)

		metaBz, err := ioutil.ReadFile(flipperMetaPath)
		require.Nil(err)
		rpcAPI.WithMetaData(metaBz)

		rpcAPI.WithLogger(logger)

		flipAPI := flipper.New(rpcAPI, contractAccountID)

		ctx := rpc.NewCtx(context.Background()).WithFrom(signature.TestKeyringPairAlice)

		resGet, err := flipAPI.Get(ctx)

		require.Nil(err)
		require.Equalf(resGet, initValue, "flipper should be initValue")

		_, err = flipAPI.Flip(ctx)
		require.Nil(err)

		resGet, err = flipAPI.Get(ctx)

		require.Nil(err)
		require.Equalf(resGet, !initValue, "flipper should be !initValue")
	})
}
