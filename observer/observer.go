package observer

import (
	"context"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/patractlabs/go-patract/utils/log"
)

type ContractObserver struct {
	logger      log.Logger
	url         string
	fromHeight  uint64
	contractIDs map[types.AccountID]bool
	metaDatas   map[types.Hash]*metadata.Data
	watcher     *api.Watcher
}

func New(logger log.Logger, url string) *ContractObserver {
	return &ContractObserver{
		logger:      logger,
		url:         url,
		contractIDs: make(map[types.AccountID]bool, 16),
		metaDatas:   make(map[types.Hash]*metadata.Data, 16),
	}
}

func (w *ContractObserver) WithFromHeight(h uint64) *ContractObserver {
	w.fromHeight = h
	return w
}

func (w *ContractObserver) WithContractIDs(ids ...types.AccountID) *ContractObserver {
	for _, id := range ids {
		w.contractIDs[id] = true
	}

	return w
}

func (w *ContractObserver) WithMetaDataBytes(codeHash types.Hash, bz []byte) *ContractObserver {
	metaData, err := metadata.New(bz)
	if err != nil {
		panic(err)
	}

	return w.WithMetaData(codeHash, metaData)
}

func (w *ContractObserver) WithMetaData(codeHash types.Hash, data *metadata.Data) *ContractObserver {
	w.metaDatas[codeHash] = data
	return w
}

func (w *ContractObserver) MetaData(codeHash types.Hash) *metadata.Data {
	return w.metaDatas[codeHash]
}

func (w *ContractObserver) WatchEvent(ctx context.Context, handler *EvtHandler) error {
	w.watcher = api.NewWatcher(w.logger, w.url)

	return w.watcher.Watch(ctx, w.fromHeight,
		func(l log.Logger, height uint64, evt *types.EventRecords) error {
			if len(evt.Contracts_Instantiated)+
				len(evt.Contracts_Evicted)+
				len(evt.Contracts_Restored)+
				len(evt.Contracts_CodeStored)+
				len(evt.Contracts_ScheduleUpdated)+
				len(evt.Contracts_ContractExecution) == 0 {
				return nil
			}

			w.logContractEvts(height, evt)
			handler.handler(l, height, evt)

			return nil
		})
}

func (w *ContractObserver) Wait() {
	w.watcher.Wait()
}

func (w *ContractObserver) logContractEvts(height uint64, evt *types.EventRecords) {
	w.logger.Debug("block event", "height", height)

	for _, e := range evt.Contracts_Instantiated {
		w.logger.Debug("Contracts_Instantiated",
			"phase", e.Phase, "topics", e.Topics, "owner", e.Owner, "contract", e.Contract)
	}

	for _, e := range evt.Contracts_Evicted {
		w.logger.Debug("Contracts_Evicted",
			"phase", e.Phase, "topics", e.Topics, "tombstone", e.Tombstone, "contract", e.Contract)
	}

	for _, e := range evt.Contracts_Restored {
		w.logger.Debug("Contracts_Restored", "phase", e.Phase, "topics", e.Topics, "donor", e.Donor, "codeHash", e.CodeHash)
	}

	for _, e := range evt.Contracts_CodeStored {
		w.logger.Debug("Contracts_CodeStored", "phase", e.Phase, "topics", e.Topics, "codeHash", e.CodeHash)
	}

	for _, e := range evt.Contracts_ScheduleUpdated {
		w.logger.Debug("Contracts_ScheduleUpdated", "phase", e.Phase, "topics", e.Topics, "schedule", e.Schedule)
	}

	for _, e := range evt.Contracts_ContractExecution {
		w.logger.Debug("Contracts_ContractExecution",
			"phase", e.Phase, "topics", e.Topics, "account", e.Account, "data", e.Data)
	}
}
