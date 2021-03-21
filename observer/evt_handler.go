package observer

import (
	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
	"github.com/patractlabs/go-patract/utils/log"
)

type (
	InstantiateHandler       func(l log.Logger, h uint64, evt types.EventContractsInstantiated)
	EvictateHandler          func(l log.Logger, h uint64, evt types.EventContractsEvicted)
	RestoredHandler          func(l log.Logger, h uint64, evt types.EventContractsRestored)
	CodeStoredHandler        func(l log.Logger, h uint64, evt types.EventContractsCodeStored)
	ScheduleUpdateHandler    func(l log.Logger, h uint64, evt types.EventContractsScheduleUpdated)
	ContractExecutionHandler func(l log.Logger, h uint64, evt types.EventContractsContractExecution)
)

type EvtHandler struct {
	instantiate       []InstantiateHandler
	evictate          []EvictateHandler
	restored          []RestoredHandler
	codeStored        []CodeStoredHandler
	scheduleUpdate    []ScheduleUpdateHandler
	contractExecution []ContractExecutionHandler
}

func NewEvtHandler() *EvtHandler {
	return &EvtHandler{}
}

func (e *EvtHandler) WithInstantiate(h InstantiateHandler) *EvtHandler {
	e.instantiate = append(e.instantiate, h)
	return e
}

func (e *EvtHandler) WithEvictate(h EvictateHandler) *EvtHandler {
	e.evictate = append(e.evictate, h)
	return e
}

func (e *EvtHandler) WithRestored(h RestoredHandler) *EvtHandler {
	e.restored = append(e.restored, h)
	return e
}

func (e *EvtHandler) WithCodeStored(h CodeStoredHandler) *EvtHandler {
	e.codeStored = append(e.codeStored, h)
	return e
}

func (e *EvtHandler) WithScheduleUpdate(h ScheduleUpdateHandler) *EvtHandler {
	e.scheduleUpdate = append(e.scheduleUpdate, h)
	return e
}

func (e *EvtHandler) WithContractExecution(h ContractExecutionHandler) *EvtHandler {
	e.contractExecution = append(e.contractExecution, h)
	return e
}

func (e *EvtHandler) handler(l log.Logger, height uint64, evts *types.EventRecords) {
	for _, evt := range evts.Contracts_Instantiated {
		evtc := evt
		for _, h := range e.instantiate {
			h(l, height, evtc)
		}
	}

	for _, evt := range evts.Contracts_Evicted {
		evtc := evt
		for _, h := range e.evictate {
			h(l, height, evtc)
		}
	}

	for _, evt := range evts.Contracts_Restored {
		evtc := evt
		for _, h := range e.restored {
			h(l, height, evtc)
		}
	}

	for _, evt := range evts.Contracts_CodeStored {
		evtc := evt
		for _, h := range e.codeStored {
			h(l, height, evtc)
		}
	}

	for _, evt := range evts.Contracts_ScheduleUpdated {
		evtc := evt
		for _, h := range e.scheduleUpdate {
			h(l, height, evtc)
		}
	}

	for _, evt := range evts.Contracts_ContractExecution {
		evtc := evt
		for _, h := range e.contractExecution {
			h(l, height, evtc)
		}
	}
}
