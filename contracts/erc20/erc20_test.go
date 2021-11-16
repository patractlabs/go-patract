package erc20_test

import (
	"context"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/signature"
	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/test/contracts"
	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

const (
	erc20WasmPath = "../../test/contracts/ink/erc20.wasm"
	erc20MetaPath = "../../test/contracts/ink/erc20.json"
	//erc20MetaPath = "/mnt/c/Users/lizhaoyang/Desktop/metadata.json"
)

var (
	bob     = utils.MustAccountIDFromSS58("5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty")
	charlie = utils.MustAccountIDFromSS58("5FLSigC9HGRKVhB9FiEo4Y3koPsNmBmLJbpXg2mp1hXcS59Y")
	dave    = utils.MustAccountIDFromSS58("5DAAnrj7VHTznn2AWBemMuyBwZWs6FNFjdyVXUeYum3PTXFy")

	instantiateSalt = []byte("ysncz3nbjjzoc7s07of3malp9d")

	totalSupply = *big.NewInt(0).Mul(
		big.NewInt(1000000000000000000),
		big.NewInt(1000000000000000000))
)

func initERC20(t *testing.T, logger log.Logger, env test.Env, authKey signature.KeyringPair) types.AccountID {
	require := require.New(t)

	codeBytes, err := ioutil.ReadFile(erc20WasmPath)
	require.Nil(err)

	cApi, err := rpc.NewContractAPI(env.URL())
	require.Nil(err)

	metaBz, err := ioutil.ReadFile(erc20MetaPath)
	require.Nil(err)
	cApi.WithMetaData(metaBz)

	ctx := api.NewCtx(context.Background()).WithFrom(authKey)

	var endowment uint64 = 1000000000000

	// Instantiate
	_, contractAccount, err := cApi.InstantiateWithCode(ctx, logger,
		types.NewCompactBalance(endowment),
		types.NewCompactGas(test.DefaultGas),
		contracts.CodeHashERC20,
		codeBytes,
		instantiateSalt,
		[]string{"new"},
		types.NewU128(totalSupply),
	)
	require.Nil(err)

	// check code
	var codeBz []byte
	if err := cApi.Native().Cli.GetStorageLatest(&codeBz,
		"Contracts", "PristineCode",
		contracts.CodeHashERC20[:], nil); err != nil {
		require.Nil(err)
	}

	t.Logf("constract %s", types.HexEncodeToString(contractAccount[:]))
	return contractAccount
}
