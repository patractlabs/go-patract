package observer_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/patractlabs/go-patract/observer"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/test/contracts"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

func TestWatch(t *testing.T) {
	require := require.New(t)

	test.ByExternCanvasEnv(t, func(logger log.Logger, env test.Env) {
		o := observer.New(logger, env.URL())
		ctx, _ := context.WithCancel(context.Background())

		metaBz, err := ioutil.ReadFile(erc20MetaPath)
		require.Nil(err)

		o = o.WithFromHeight(0).WithMetaDataBytes(contracts.CodeHashERC20, metaBz)

		o.WatchEvent(ctx)
	})
}
