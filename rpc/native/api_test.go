package native_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/signature"
	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/rpc/native"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/test/contracts"
	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

var authKey = signature.TestKeyringPairAlice

func TestDeployAndCallERC20(t *testing.T) {
	require := require.New(t)

	codeBytes, err := ioutil.ReadFile("../../test/contracts/ink/erc20.wasm")
	require.Nil(err)

	test.ByNodeEnv(t, func(logger log.Logger, env test.Env) {
		cli, err := api.NewClient(logger, env.URL())
		require.Nil(err)

		contractCli := native.NewContractAPI(cli)
		ctx := api.NewCtx(context.Background()).WithFrom(authKey)

		instantiate4ERC20 := types.MustHexDecodeString("0xd183512b00000000109f4bb31507c97bce97c000")

		// Instantiate
		_, err = contractCli.InstantiateWithCode(ctx,
			types.NewCompactBalance(1000000000000000),
			types.NewCompactGas(test.DefaultGas),
			codeBytes,
			instantiate4ERC20,
			nil,
		)
		require.Nil(err)

		// check code
		var codeBz []byte
		if err := cli.GetStorageLatest(&codeBz,
			"Contracts", "PristineCode",
			contracts.CodeHashERC20[:], nil); err != nil {
			require.Nil(err)
		}
	})
}
