package main

import (
	"context"
	"io/ioutil"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/patractlabs/go-patract/contracts/erc20"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/patractlabs/go-patract/observer"
	"github.com/patractlabs/go-patract/utils"
	"github.com/patractlabs/go-patract/utils/log"
)

func observerEvts(logger log.Logger, db *erc20DB) (*observer.ContractObserver, context.CancelFunc, error) {
	o := observer.New(logger, *flagNodeURL)
	ctx, cancelFunc := context.WithCancel(context.Background())

	metaBz, err := ioutil.ReadFile(*flagContractMetadata)
	if err != nil {
		logger.Error("read metadata err", "err", err, "path", *flagContractMetadata)
		return o, cancelFunc, err
	}

	hash, err := utils.DecodeAccountIDFromSS58(*flagCodeHash)
	if err != nil {
		logger.Error("hex decode code hash err", "err", err, "path", *flagCodeHash)
		return o, cancelFunc, err
	}

	if err := db.init(*flagCodeHash); err != nil {
		return o, cancelFunc, err
	}

	o = o.WithFromHeight(0).WithMetaDataBytes(types.NewHash(hash[:]), metaBz)
	metaData := o.MetaData(types.NewHash(hash[:]))

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

			if err := db.OnEventTransfer(&transfer); err != nil {
				logger.Error("evt process transfer error", "err", err, "height", height)
				return
			}
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
		return o, cancelFunc, err
	}

	return o, cancelFunc, nil
}
