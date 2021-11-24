package incrementer_test

import (
	"context"
	"io/ioutil"
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
	incrementerWasmPath = "../../test/contracts/ink/incrementer.wasm"
	incrementerMetaPath = "../../test/contracts/ink/incrementer.json"
)

var (
	bob     = utils.MustAccountIDFromSS58("5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty")
	charlie = utils.MustAccountIDFromSS58("5FLSigC9HGRKVhB9FiEo4Y3koPsNmBmLJbpXg2mp1hXcS59Y")
	dave    = utils.MustAccountIDFromSS58("5DAAnrj7VHTznn2AWBemMuyBwZWs6FNFjdyVXUeYum3PTXFy")

	instantiateSalt = []byte("ysncz3nbjjzoc7s07of3malp9d")

	initValue   = types.NewI32(10)
	addValue    = types.NewI32(5)
	targetValue = initValue + addValue
)

func initIncrementer(t *testing.T, logger log.Logger, env test.Env, authKey signature.KeyringPair) types.AccountID {
	require := require.New(t)

	codeBytes, err := ioutil.ReadFile(incrementerWasmPath)
	require.Nil(err)

	cApi, err := rpc.NewContractAPI(env.URL())
	require.Nil(err)

	metaBz, err := ioutil.ReadFile(incrementerMetaPath)
	require.Nil(err)
	cApi.WithMetaData(metaBz)

	ctx := api.NewCtx(context.Background()).WithFrom(authKey)

	var endowment uint64 = 1000000000000

	// Instantiate
	_, contractAccount, err := cApi.InstantiateWithCode(ctx, logger,
		types.NewCompactBalance(endowment),
		types.NewCompactGas(test.DefaultGas),
		contracts.CodeHashIncrementer,
		codeBytes,
		instantiateSalt,
		[]string{"new"},
		initValue,
	)
	require.Nil(err)

	// check code
	var codeBz []byte
	if err := cApi.Native().Cli.GetStorageLatest(&codeBz,
		"Contracts", "PristineCode",
		contracts.CodeHashIncrementer[:], nil); err != nil {
		require.Nil(err)
	}

	t.Logf("constract %s", types.HexEncodeToString(contractAccount[:]))
	return contractAccount
}
