package contract_terminate_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/signature"
	"github.com/patractlabs/go-patract/contracts/contract-terminate"
	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

func TestFlip(t *testing.T) {
	test.ByNodeEnv(t, func(logger log.Logger, env test.Env) {
		require := require.New(t)
		contractAccountID := initContractTerminate(t, logger, env, signature.TestKeyringPairAlice)
		rpcAPI, err := rpc.NewContractAPI(env.URL())
		require.Nil(err)

		metaBz, err := ioutil.ReadFile(contractTerminateMetaPath)
		require.Nil(err)
		rpcAPI.WithMetaData(metaBz)

		rpcAPI.WithLogger(logger)

		contractAPI := contract_terminate.New(rpcAPI, contractAccountID)

		// transfer alice to bob
		ctx := rpc.NewCtx(context.Background()).WithFrom(signature.TestKeyringPairAlice)

		_, err = contractAPI.TerminateMe(ctx)

		require.Nil(err)
	})
}
