package contract_terminate_test

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
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

const (
	contractTerminateWasmPath = "../../test/contracts/ink/contract_terminate.wasm"
	contractTerminateMetaPath = "../../test/contracts/ink/contract_terminate.json"
)

var (
	instantiateSalt = []byte("ysncz3nbjjzoc7s07of3malp9d")
)

func initContractTerminate(t *testing.T, logger log.Logger, env test.Env, authKey signature.KeyringPair) types.AccountID {
	require := require.New(t)

	codeBytes, err := ioutil.ReadFile(contractTerminateWasmPath)
	require.Nil(err)

	cApi, err := rpc.NewContractAPI(env.URL())
	require.Nil(err)

	metaBz, err := ioutil.ReadFile(contractTerminateMetaPath)
	require.Nil(err)
	cApi.WithMetaData(metaBz)

	ctx := api.NewCtx(context.Background()).WithFrom(authKey)

	var endowment uint64 = 1000000000000

	// Instantiate
	_, contractAccount, err := cApi.InstantiateWithCode(ctx, logger,
		types.NewCompactBalance(endowment),
		types.NewCompactGas(test.DefaultGas),
		contracts.CodeHashContractTerminate,
		codeBytes,
		instantiateSalt,
		[]string{"new"},
	)
	require.Nil(err)

	// check code
	var codeBz []byte
	if err := cApi.Native().Cli.GetStorageLatest(&codeBz,
		"Contracts", "PristineCode",
		contracts.CodeHashContractTerminate[:], nil); err != nil {
		require.Nil(err)
	}

	t.Logf("constract %s", types.HexEncodeToString(contractAccount[:]))
	return contractAccount
}
