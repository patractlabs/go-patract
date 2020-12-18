package api_test

import (
	"context"
	"testing"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/utils/log"
)

func TestWatchEventByBlock(t *testing.T) {
	test.ByCanvasEnv(t, func(logger log.Logger, env test.Env) {
		time.Sleep(10 * time.Second)

		w := api.NewWatcher(logger, env.URL())
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			time.Sleep(10 * time.Second)
			cancel()
		}()

		w.Watch(ctx, 1, func(l log.Logger, height uint64, evt *types.EventRecords) error {
			l.Info("scan block event", "height", height)

			for _, e := range evt.System_ExtrinsicSuccess {
				l.Info("System_ExtrinsicSuccess", "phase", e.Phase, "dispatchInfo", e.DispatchInfo, "topics", e.Topics)
			}

			for _, e := range evt.Balances_Transfer {
				l.Info("Balances_Transfer", "phase", e.Phase, "from", e.From, "to", e.To, "value", e.Value)
			}

			for _, e := range evt.Contracts_ContractExecution {
				l.Info("Contracts_ContractExecution", "phase", e.Phase, "account", e.Account, "data", e.Data, "topics", e.Topics)
			}

			return nil
		})

		w.Wait()
	})
}
