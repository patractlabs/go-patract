package rpc_test

import (
	"context"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v2/signature"
	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/test/contracts"
	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

var (
	bob     = utils.MustAccountIDFromSS58("5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty")
	tester  = types.MustAddressFromHexAccount("")
	authKey = signature.TestKeyringPairAlice

	instantiateSalt = []byte("79z58iw8h03yozt8jefhkr1axcysncz3nbjjzoc7s07of3malp9d4xfhduyylsuc")
	contractAddress = "5HLboJfhkqpj6qgnQoLNeFPS1JKvseZJwHdpqkZRNY57aEgh"

	totalSupply = *big.NewInt(0).Mul(
		big.NewInt(1000000000000000000),
		big.NewInt(1000000000000000000))

	erc20WasmPath = "../test/contracts/ink/erc20.wasm"
	erc20MetaPath = "../test/contracts/ink/erc20.json"
)

func TestDeployAndCallERC20(t *testing.T) {
	test.ByCanvasEnv(t, func(logger log.Logger, env test.Env) {
		require := require.New(t)

		metaBz, err := ioutil.ReadFile(erc20MetaPath)
		require.Nil(err)

		codeBytes, err := ioutil.ReadFile(erc20WasmPath)
		require.Nil(err)

		api, err := rpc.NewContractAPI(env.URL())
		require.Nil(err)

		api.WithLogger(logger)
		err = api.WithMetaData(metaBz)
		require.Nil(err)

		ctx := rpc.NewCtx(context.Background()).WithFrom(authKey)

		// Instantiate
		hash, contractAccount, err := api.InstantiateWithCode(ctx,
			logger,
			types.NewCompactBalance(10000000000000000),
			types.NewCompactGas(test.DefaultGas),
			contracts.CodeHashERC20,
			codeBytes,
			instantiateSalt,
			types.NewU128(totalSupply),
		)

		require.Nil(err)
		t.Logf("instantiate hash %v %v", hash, contractAccount)
	})
}
