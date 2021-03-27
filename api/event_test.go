package api_test

import (
	"context"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

func TestQueryEventByBlock(t *testing.T) {
	require := require.New(t)

	test.ByCanvasEnv(t, func(logger log.Logger, env test.Env) {
		cli, err := api.NewClient(logger, env.URL())
		require.Nil(err)

		hash, err := cli.SubmitAndWaitExtrinsic(
			api.NewCtx(context.Background()).WithFrom(authKey),
			"Balances.transfer", bob, types.NewUCompactFromUInt(1000000000000000))
		require.Nil(err)

		events, err := cli.QueryEventByBlockHash(hash)
		require.Nil(err)

		t.Logf("evt %v", events)
	})

}
