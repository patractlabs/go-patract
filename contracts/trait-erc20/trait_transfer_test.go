package trait_erc20_test

import (
	"context"
	trait_erc20 "github.com/patractlabs/go-patract/contracts/trait-erc20"
	"io/ioutil"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/signature"
	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	test.ByNodeEnv(t, func(logger log.Logger, env test.Env) {
		require := require.New(t)

		contractAccountID := initTraitERC20(t, logger, env, signature.TestKeyringPairAlice)
		rpcAPI, err := rpc.NewContractAPI(env.URL())
		require.Nil(err)

		metaBz, err := ioutil.ReadFile(erc20MetaPath)
		require.Nil(err)
		rpcAPI.WithMetaData(metaBz)

		rpcAPI.WithLogger(logger)

		erc20API := trait_erc20.New(rpcAPI, contractAccountID)

		// transfer alice to bob
		ctx := rpc.NewCtx(context.Background()).WithFrom(signature.TestKeyringPairAlice)

		aliceTotal, err := erc20API.BalanceOf(ctx, test.AliceAccountID)
		require.Nil(err)
		require.Equalf(*aliceTotal.Int, totalSupply, "alice should be total supply")

		// transfer
		amt2Bob := types.NewBalanceByU64(100)
		_, err = erc20API.Transfer(ctx, bob, amt2Bob)
		require.Nil(err)

		bobBalance, err := erc20API.BalanceOf(ctx, bob)
		require.Nil(err)
		require.Equalf(bobBalance, amt2Bob, "bob Balance should be amt2Bob")

		aliceNewTotal, err := erc20API.BalanceOf(ctx, test.AliceAccountID)
		require.Nil(err)
		require.Equalf(aliceTotal.Int, aliceNewTotal.Int.Add(aliceNewTotal.Int, amt2Bob.Int),
			"alice new add transfer should be alice old")
	})
}
