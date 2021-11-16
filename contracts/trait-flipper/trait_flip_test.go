package trait_flipper_test

import (
	"context"
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/patractlabs/go-patract/contracts/trait-flipper"
	"io/ioutil"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/signature"
	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

func TestTraitFlip(t *testing.T) {
	test.ByNodeEnv(t, func(logger log.Logger, env test.Env) {
		require := require.New(t)
		contractAccountID := initTraitFlipper(t, logger, env, signature.TestKeyringPairAlice)

		rpcAPI, err := rpc.NewContractAPI(env.URL())
		require.Nil(err)

		metaBz, err := ioutil.ReadFile(flipperMetaPath)
		require.Nil(err)
		rpcAPI.WithMetaData(metaBz)

		rpcAPI.WithLogger(logger)

		traitFlipAPI := trait_flipper.New(rpcAPI, contractAccountID)

		// transfer alice to bob
		ctx := rpc.NewCtx(context.Background()).WithFrom(signature.TestKeyringPairAlice)

		resGet, err := traitFlipAPI.Get(ctx)
		require.Nil(err)
		require.Equalf(resGet, types.NewBool(false), "flipper should be false")

		_, err = traitFlipAPI.Flip(ctx)
		require.Nil(err)

		resGet, err = traitFlipAPI.Get(ctx)

		require.Nil(err)
		require.Equalf(resGet, types.NewBool(true), "flipper should be true")
	})
}
