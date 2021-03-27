package observer_test

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
	"github.com/patractlabs/go-patract/contracts/erc20"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/patractlabs/go-patract/observer"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/test/contracts"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

const (
	erc20WasmPath = "../test/contracts/ink/erc20.wasm"
	erc20MetaPath = "../test/contracts/ink/erc20.json"
)

func TestWatch(t *testing.T) {
	require := require.New(t)

	test.ByCanvasEnv(t, func(logger log.Logger, env test.Env) {
		o := observer.New(logger, env.URL())
		ctx, cc := context.WithCancel(context.Background())

		metaBz, err := ioutil.ReadFile(erc20MetaPath)
		require.Nil(err)

		o = o.WithFromHeight(0).WithMetaDataBytes(contracts.CodeHashERC20, metaBz)

		metaData := o.MetaData(contracts.CodeHashERC20)

		h := observer.NewEvtHandler()
		h = h.WithContractExecution(func(l log.Logger, height uint64, evt types.EventContractsContractExecution) {
			data := evt.Data

			l.Debug("handler contract execution", "height", height)

			typ := metadata.GetEvtTypeIdx(data)
			switch typ {
			case 0:
				var transfer erc20.EventTransfer
				err := metaData.Spec.Events.DecodeEvt(metaData.NewCtxForDecode(data).WithLogger(l), &transfer)
				if err != nil {
					logger.Error("evt decode transfer error", "err", err, "height", height)
				}
				logger.Info("transfer event", "evt", transfer)
			case 1:
				var approve erc20.EventApproval
				err := metaData.Spec.Events.DecodeEvt(metaData.NewCtxForDecode(data).WithLogger(l), &approve)
				if err != nil {
					logger.Error("evt decode approve error", "err", err, "height", height)
				}
				logger.Info("approve event", "evt", approve)
			}
		})

		go func() {
			time.Sleep(2 * time.Second)
			cc()
		}()

		o.WatchEvent(ctx, h)
	})
}
