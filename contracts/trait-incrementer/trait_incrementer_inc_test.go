package trait_incrementer_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/signature"
	"github.com/patractlabs/go-patract/contracts/trait-incrementer"
	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

func TestIncrementerInc(t *testing.T) {
	test.ByNodeEnv(t, func(logger log.Logger, env test.Env) {
		require := require.New(t)
		contractAccountID := initTraitIncrementer(t, logger, env, signature.TestKeyringPairAlice)

		rpcAPI, err := rpc.NewContractAPI(env.URL())
		require.Nil(err)

		metaBz, err := ioutil.ReadFile(traitIncrementerMetaPath)
		require.Nil(err)
		rpcAPI.WithMetaData(metaBz)

		rpcAPI.WithLogger(logger)

		traitIncrementerAPI := trait_incrementer.New(rpcAPI, contractAccountID)

		ctx := rpc.NewCtx(context.Background()).WithFrom(signature.TestKeyringPairAlice)
		_, err = traitIncrementerAPI.IncBy(ctx, addValue)
		require.Nil(err)

		value, err := traitIncrementerAPI.Get(ctx)
		require.Nil(err)
		require.Equalf(value, targetValue, "The value after Inc must be targetValue.")
	})
}
