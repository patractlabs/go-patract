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
	"github.com/stretchr/testify/assert"
)

var (
	bob     = types.MustAddressFromHexAccount("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	tester  = types.MustAddressFromHexAccount("")
	authKey = signature.TestKeyringPairAlice
)

func TestPutCode(t *testing.T) {
	assert := assert.New(t)

	codeBytes, err := ioutil.ReadFile("../../test/contracts/ink/erc20.wasm")
	assert.Nil(err)

	test.ByCanvasEnv(t, func(logger log.Logger, env test.Env) {
		cli, err := api.NewClient(logger, env.URL())
		assert.Nil(err)

		contractCli := native.NewContractAPI(cli)

		_, err = contractCli.PutCode(api.NewCtx(context.Background()).WithFrom(authKey), codeBytes)
		assert.Nil(err)

		// check code
		var codeBz []byte
		if err := cli.GetStorageLatest(&codeBz,
			"Contracts", "PristineCode",
			[]byte(contracts.CodeHashERC20[:]), nil); err != nil {
			assert.Nil(err)
		}

		assert.True(bytes.Equal(codeBytes, codeBz), "code should be equal")
	})
}

func TestDeployAndCallERC20(t *testing.T) {
	assert := assert.New(t)

	codeBytes, err := ioutil.ReadFile("../../test/contracts/ink/erc20.wasm")
	assert.Nil(err)

	test.ByExternCanvasEnv(t, func(logger log.Logger, env test.Env) {
		cli, err := api.NewClient(logger, env.URL())
		assert.Nil(err)

		contractCli := native.NewContractAPI(cli)
		ctx := api.NewCtx(context.Background()).WithFrom(authKey)

		_, err = contractCli.PutCode(ctx, codeBytes)
		assert.Nil(err)

		// check code
		var codeBz []byte
		if err := cli.GetStorageLatest(&codeBz,
			"Contracts", "PristineCode",
			[]byte(contracts.CodeHashERC20[:]), nil); err != nil {
			assert.Nil(err)
		}

		instantiate4ERC20 := types.MustHexDecodeString("0xd183512b00000000109f4bb31507c97bce97c000")

		// Instantiate
		_, err = contractCli.Instantiate(ctx,
			types.NewCompactBalance(10000000000000000),
			types.NewCompactGas(test.DefaultGas),
			contracts.CodeHashERC20,
			instantiate4ERC20,
		)
		assert.Nil(err)
	})
}
