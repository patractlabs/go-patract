package delegator_test

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
	delegatorWasmPath = "../../test/contracts/ink/delegator/delegator.wasm"
	delegatorMetaPath = "../../test/contracts/ink/delegator/delegator.json"

	accumulatorWasmPath = "../../test/contracts/ink/delegator/accumulator/accumulator.wasm"
	accumulatorMetaPath = "../../test/contracts/ink/delegator/accumulator/accumulator.json"
	adderWasmPath       = "../../test/contracts/ink/delegator/adder/adder.wasm"
	adderMetaPath       = "../../test/contracts/ink/delegator/adder/adder.json"
	subberWasmPath      = "../../test/contracts/ink/delegator/subber/subber.wasm"
	subberMetaPath      = "../../test/contracts/ink/delegator/subber/subber.json"
)

var (
	instantiateSalt = []byte("ysncz3nbjjzoc7s07of3malp9d")

	initValue   = types.NewI32(10)
	changeValue = types.NewI32(5)

	version = types.NewU32(0)

	accumulatorParam = struct {
		AccumulatorCodeHash types.Hash
	}{
		AccumulatorCodeHash: contracts.CodeHashAccumulator,
	}

	adderParam = struct {
		AdderCodeHash types.Hash
	}{
		AdderCodeHash: contracts.CodeHashAdder,
	}

	subberParam = struct {
		SubberCodeHash types.Hash
	}{
		SubberCodeHash: contracts.CodeHashSubber,
	}
)

type AccountParam struct {
	Inner struct {
		AccountID struct {
			Account types.AccountID
		} `scale:"account_id"`
	} `scale:"inner"`
}

func initDelegator(t *testing.T, logger log.Logger, env test.Env, authKey signature.KeyringPair) types.AccountID {
	require := require.New(t)

	codeBytes, err := ioutil.ReadFile(delegatorWasmPath)
	require.Nil(err)

	cApi, err := rpc.NewContractAPI(env.URL())
	require.Nil(err)

	metaBz, err := ioutil.ReadFile(delegatorMetaPath)
	require.Nil(err)
	cApi.WithMetaData(metaBz)

	ctx := api.NewCtx(context.Background()).WithFrom(authKey)

	var endowment uint64 = 1000000000000

	// Instantiate
	_, contractAccount, err := cApi.InstantiateWithCode(ctx, logger,
		types.NewCompactBalance(endowment),
		types.NewCompactGas(test.DefaultGas),
		contracts.CodeHashDelegator,
		codeBytes,
		instantiateSalt,
		[]string{"new"},
		initValue,
		version,
		accumulatorParam,
		adderParam,
		subberParam,
	)
	require.Nil(err)

	// check code
	var codeBz []byte
	if err := cApi.Native().Cli.GetStorageLatest(&codeBz,
		"Contracts", "PristineCode",
		contracts.CodeHashDelegator[:], nil); err != nil {
		require.Nil(err)
	}

	t.Logf("constract %s", types.HexEncodeToString(contractAccount[:]))
	return contractAccount
}

func initAccumulator(t *testing.T, logger log.Logger, env test.Env, authKey signature.KeyringPair) types.AccountID {
	require := require.New(t)

	codeBytes, err := ioutil.ReadFile(accumulatorWasmPath)
	require.Nil(err)

	cApi, err := rpc.NewContractAPI(env.URL())
	require.Nil(err)

	metaBz, err := ioutil.ReadFile(accumulatorMetaPath)
	require.Nil(err)
	cApi.WithMetaData(metaBz)

	ctx := api.NewCtx(context.Background()).WithFrom(authKey)

	var endowment uint64 = 1000000000000

	// Instantiate
	_, contractAccount, err := cApi.InstantiateWithCode(ctx, logger,
		types.NewCompactBalance(endowment),
		types.NewCompactGas(test.DefaultGas),
		contracts.CodeHashAccumulator,
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
		contracts.CodeHashAccumulator[:], nil); err != nil {
		require.Nil(err)
	}

	t.Logf("constract %s", types.HexEncodeToString(contractAccount[:]))
	return contractAccount
}

func initAdder(t *testing.T, logger log.Logger, env test.Env, authKey signature.KeyringPair) types.AccountID {
	require := require.New(t)

	codeBytes, err := ioutil.ReadFile(adderWasmPath)
	require.Nil(err)

	cApi, err := rpc.NewContractAPI(env.URL())
	require.Nil(err)

	metaBz, err := ioutil.ReadFile(adderMetaPath)
	require.Nil(err)
	cApi.WithMetaData(metaBz)

	ctx := api.NewCtx(context.Background()).WithFrom(authKey)

	var endowment uint64 = 1000000000000

	// Instantiate
	_, contractAccount, err := cApi.InstantiateWithCode(ctx, logger,
		types.NewCompactBalance(endowment),
		types.NewCompactGas(test.DefaultGas),
		contracts.CodeHashAdder,
		codeBytes,
		instantiateSalt,
		[]string{"new"},
		AccountParam{
			struct {
				AccountID struct {
					Account types.AccountID
				} `scale:"account_id"`
			}{AccountID: struct {
				Account types.AccountID
			}{
				Account: test.AliceAccountID,
			}}},
	)
	require.Nil(err)

	// check code
	var codeBz []byte
	if err := cApi.Native().Cli.GetStorageLatest(&codeBz,
		"Contracts", "PristineCode",
		contracts.CodeHashAdder[:], nil); err != nil {
		require.Nil(err)
	}

	t.Logf("constract %s", types.HexEncodeToString(contractAccount[:]))
	return contractAccount
}

func initSubber(t *testing.T, logger log.Logger, env test.Env, authKey signature.KeyringPair) types.AccountID {
	require := require.New(t)

	codeBytes, err := ioutil.ReadFile(subberWasmPath)
	require.Nil(err)

	cApi, err := rpc.NewContractAPI(env.URL())
	require.Nil(err)

	metaBz, err := ioutil.ReadFile(subberMetaPath)
	require.Nil(err)
	cApi.WithMetaData(metaBz)

	ctx := api.NewCtx(context.Background()).WithFrom(authKey)

	var endowment uint64 = 1000000000000

	// Instantiate
	_, contractAccount, err := cApi.InstantiateWithCode(ctx, logger,
		types.NewCompactBalance(endowment),
		types.NewCompactGas(test.DefaultGas),
		contracts.CodeHashSubber,
		codeBytes,
		instantiateSalt,
		[]string{"new"},
		AccountParam{
			struct {
				AccountID struct {
					Account types.AccountID
				} `scale:"account_id"`
			}{AccountID: struct {
				Account types.AccountID
			}{
				Account: test.AliceAccountID,
			}}},
	)
	require.Nil(err)

	// check code
	var codeBz []byte
	if err := cApi.Native().Cli.GetStorageLatest(&codeBz,
		"Contracts", "PristineCode",
		contracts.CodeHashSubber[:], nil); err != nil {
		require.Nil(err)
	}

	t.Logf("constract %s", types.HexEncodeToString(contractAccount[:]))
	return contractAccount
}
