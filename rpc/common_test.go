package rpc_test

import (
	"context"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/signature"
	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/rpc/native"
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

	totalSupply = *big.NewInt(0).Mul(
		big.NewInt(1000000000000000000),
		big.NewInt(1000000000000000000))
)

const (
	erc20WasmPath = "../test/contracts/ink/erc20.wasm"
	erc20MetaPath = "../test/contracts/ink/erc20.json"
)

func initERC20(t *testing.T, logger log.Logger, env test.Env) {
	require := require.New(t)

	codeBytes, err := ioutil.ReadFile(erc20WasmPath)
	require.Nil(err)

	cli, err := api.NewClient(logger, env.URL())
	require.Nil(err)

	contractCli := native.NewContractAPI(cli)
	ctx := api.NewCtx(context.Background()).WithFrom(authKey)

	_, err = contractCli.PutCode(ctx, codeBytes)
	require.Nil(err)

	// check code
	var codeBz []byte
	if err := cli.GetStorageLatest(&codeBz,
		"Contracts", "PristineCode",
		[]byte(contracts.CodeHashERC20[:]), nil); err != nil {
		require.Nil(err)
	}
}
