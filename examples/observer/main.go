package main

import (
	"context"
	"flag"
	"io/ioutil"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/patractlabs/go-patract/contracts/erc20"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/patractlabs/go-patract/observer"
	"github.com/patractlabs/go-patract/utils"
	"github.com/patractlabs/go-patract/utils/log"
)

var (
	flagURL              = flag.String("url", "ws://localhost:9944", "url to chain node")
	flagContractMetadata = flag.String("metadata", "", "metadata for contract")
	flagCodeHash         = flag.String("hash", "", "codeHash for contract")
)

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}

	logger := log.NewLogger()
	defer func() {
		logger.Flush()
	}()

	o := observer.New(logger, *flagURL)
	ctx, cancelFunc := context.WithCancel(context.Background())

	metaBz, err := ioutil.ReadFile(*flagContractMetadata)
	if err != nil {
		logger.Error("read metadata err", "err", err, "path", *flagContractMetadata)
		return
	}

	hash, err := types.HexDecodeString(*flagCodeHash)
	if err != nil {
		logger.Error("hex decode code hash err", "err", err, "path", *flagCodeHash)
		return
	}

	o = o.WithFromHeight(0).WithMetaDataBytes(types.NewHash(hash), metaBz)

	metaData := o.MetaData(types.NewHash(hash))

	h := observer.NewEvtHandler()
	h = h.WithContractExecution(func(l log.Logger, height uint64, evt types.EventContractsContractExecution) {
		data := evt.Data

		l.Debug("handler contract execution", "height", height)

		typ := metadata.GetEvtTypeIdx(data)
		switch typ {
		case 0:
			var transfer erc20.EventTransfer
			err := metaData.V1.Spec.Events.DecodeEvt(metaData.NewCtxForDecode(data).WithLogger(l), &transfer)
			if err != nil {
				logger.Error("evt decode transfer error", "err", err, "height", height)
			}
			logger.Info("transfer event", "evt", transfer)
		case 1:
			var approve erc20.EventApproval
			err := metaData.V1.Spec.Events.DecodeEvt(metaData.NewCtxForDecode(data).WithLogger(l), &approve)
			if err != nil {
				logger.Error("evt decode approve error", "err", err, "height", height)
			}
			logger.Info("approve event", "evt", approve)
		}
	})

	if err := o.WatchEvent(ctx, h); err != nil {
		logger.Error("watch event error", "err", err)
		return
	}

	utils.HoldToClose(func() {
		cancelFunc()
		o.Wait()

		logger.Info("observer stop")
		logger.Flush()
	})
}
