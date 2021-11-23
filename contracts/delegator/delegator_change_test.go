package delegator_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/signature"
	"github.com/patractlabs/go-patract/contracts/delegator"
	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

func TestDelegator(t *testing.T) {
	test.ByExternEnv(t, func(logger log.Logger, env test.Env) {
		require := require.New(t)
		_ = initAccumulator(t, logger, env, signature.TestKeyringPairAlice)
		fmt.Println("====================================== 1")
		fmt.Println("====================================== 1")
		fmt.Println("====================================== 1")
		_ = initAdder(t, logger, env, signature.TestKeyringPairAlice)
		fmt.Println("====================================== 2")
		fmt.Println("====================================== 2")
		fmt.Println("====================================== 2")
		_ = initSubber(t, logger, env, signature.TestKeyringPairAlice)
		fmt.Println("====================================== 3")
		fmt.Println("====================================== 3")
		fmt.Println("====================================== 3")

		contractAccountID := initDelegator(t, logger, env, signature.TestKeyringPairAlice)
		rpcAPI, err := rpc.NewContractAPI(env.URL())
		require.Nil(err)

		metaBz, err := ioutil.ReadFile(delegatorMetaPath)
		require.Nil(err)
		rpcAPI.WithMetaData(metaBz)

		rpcAPI.WithLogger(logger)

		delegatorAPI := delegator.New(rpcAPI, contractAccountID)

		ctx := rpc.NewCtx(context.Background()).WithFrom(signature.TestKeyringPairAlice)

		_, err = delegatorAPI.Change(ctx, changeValue)
		require.Nil(err)

		addValue, err := delegatorAPI.Get(ctx)
		require.Nil(err)
		require.Equalf(addValue, changeValue+initValue, "It must be the result of addition")

		_, err = delegatorAPI.Switch(ctx)
		require.Nil(err)
		_, err = delegatorAPI.Change(ctx, changeValue)
		require.Nil(err)

		subValue, err := delegatorAPI.Get(ctx)
		require.Nil(err)
		require.Equalf(subValue, initValue, "It must be the result of subtraction")
	})
}
