package api_test

import (
	"context"
	"testing"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/signature"
	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/assert"
)

var (
	bob     = types.MustAddressFromHexAccount("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	authKey = signature.TestKeyringPairAlice
)

func TestSubmitAndWaitExtrinsic(t *testing.T) {
	assert := assert.New(t)

	test.ByCanvasEnv(t, func(logger log.Logger, env test.Env) {
		cli, err := api.NewClient(logger, env.URL())
		assert.Nil(err)

		_, err = cli.SubmitAndWaitExtrinsic(
			api.NewCtx(context.Background()).WithFrom(authKey),
			"Balances.transfer", bob, types.NewUCompactFromUInt(1000000000000000))
		assert.Nil(err)
	})
}

func TestSubmitAndWaitExtrinsicCancel(t *testing.T) {
	assert := assert.New(t)

	test.ByCanvasEnv(t, func(logger log.Logger, env test.Env) {
		cli, err := api.NewClient(logger, env.URL())
		assert.Nil(err)

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		_, err = cli.SubmitAndWaitExtrinsic(
			api.NewCtx(ctx).WithFrom(authKey),
			"Balances.transfer", bob, types.NewCompactBalance(1000000000000000))
		assert.Equal(err, context.Canceled)
	})
}
