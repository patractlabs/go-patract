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

func TestScanEventByBlock(t *testing.T) {
	test.ByCanvasEnv(t, func(logger log.Logger, env test.Env) {
		time.Sleep(10 * time.Second)

		scanner := api.NewScanner(logger, env.URL())

		scanner.Scan(context.Background(), 1, func(l log.Logger, height uint64, evt *types.EventRecords) error {
			l.Info("scan block event %d", height)

			for _, e := range evt.System_ExtrinsicSuccess {
				l.Info("System_ExtrinsicSuccess:: (phase=%#v) : %v, %v, %v", e.Phase, e.DispatchInfo, e.Topics)
			}

			for _, e := range evt.Balances_Transfer {
				l.Info("Balances_Transfer:: (phase=%#v) : %v, %v, %v", e.Phase, e.From, e.To, e.Value)
			}

			for _, e := range evt.Contracts_ContractExecution {
				l.Info("Contracts_ContractExecution:: (phase=%#v) : %v, %v, %v", e.Phase, e.Account, e.Data, e.Topics)
			}

			return nil
		})

		scanner.Wait()
	})
}
