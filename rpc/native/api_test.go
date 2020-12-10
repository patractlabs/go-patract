package native_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/signature"
	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/rpc/native"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/test/contracts"
	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

var (
	bob     = types.MustAddressFromHexAccount("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	tester  = types.MustAddressFromHexAccount("")
	authKey = signature.TestKeyringPairAlice
)

func TestPutCode(t *testing.T) {
	require := require.New(t)

	codeBytes, err := ioutil.ReadFile("../../test/contracts/ink/erc20.wasm")
	require.Nil(err)

	test.ByCanvasEnv(t, func(logger log.Logger, env test.Env) {
		cli, err := api.NewClient(logger, env.URL())
		require.Nil(err)

		contractCli := native.NewContractAPI(cli)

		_, err = contractCli.PutCode(api.NewCtx(context.Background()).WithFrom(authKey), codeBytes)
		require.Nil(err)

		// check code
		var codeBz []byte
		if err := cli.GetStorageLatest(&codeBz,
			"Contracts", "PristineCode",
			[]byte(contracts.CodeHashERC20[:]), nil); err != nil {
			require.Nil(err)
		}

		require.True(bytes.Equal(codeBytes, codeBz), "code should be equal")
	})
}

func TestDeployAndCallERC20(t *testing.T) {
	require := require.New(t)

	codeBytes, err := ioutil.ReadFile("../../test/contracts/ink/erc20.wasm")
	require.Nil(err)

	test.ByCanvasEnv(t, func(logger log.Logger, env test.Env) {
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

		instantiate4ERC20 := types.MustHexDecodeString("0xd183512b00000000109f4bb31507c97bce97c000")

		// Instantiate
		_, err = contractCli.Instantiate(ctx,
			types.NewCompactBalance(1000000000000000),
			types.NewCompactGas(test.DefaultGas),
			contracts.CodeHashERC20,
			instantiate4ERC20,
		)
		require.Nil(err)
	})
}
