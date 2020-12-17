package api

import (
	"context"
	"sync"

	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/pkg/errors"
)

type scanner struct {
	wg                sync.WaitGroup
	cli               *Client
	logger            log.Logger
	eventChann        chan *types.EventRecords
	mutex             sync.RWMutex
	latestBlockHeight uint64
}

type EventHandler func(logger log.Logger, evt *types.EventRecords) error

func NewScanner(logger log.Logger, url string) *scanner {
	cli, err := NewClient(logger, url)
	if err != nil {
		panic(err)
	}

	return &scanner{
		cli:        cli,
		logger:     logger,
		eventChann: make(chan *types.EventRecords, 4096),
	}
}

func (s *scanner) LastestBlockHeight() uint64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.latestBlockHeight
}

func (s *scanner) setToHeight(height uint64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.latestBlockHeight = height
}

func (s *scanner) Scan(ctx context.Context, fromHeight uint64, h EventHandler) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		for {
			evts, ok := <-s.eventChann
			if !ok {
				// closed scanner, closed by last
				s.logger.Info("stop handler gorountinue")
				return
			}

			if err := h(s.logger, evts); err != nil {
				s.logger.Error("handler error", "err", err.Error())
			}
		}
	}()

	s.wg.Add(1)
	go func() {
		defer func() {
			s.wg.Done()

			// stop the scanner handler
			close(s.eventChann)
		}()

		if err := s.scanBlocksImp(ctx, fromHeight); err != nil {
			s.logger.Error("scan block error", "err", err)
		}
	}()

}

func (s *scanner) scanBlocksImp(ctx context.Context, fromHeight uint64) error {
	currentBlockHeight := fromHeight
	if currentBlockHeight < 1 {
		currentBlockHeight = 1
	}

	last, err := s.cli.API().RPC.Chain.GetBlockLatest()
	if err != nil {
		return errors.Wrapf(err, "get latest block err")
	}

	lastHeight := uint64(last.Block.Header.Number)

	if lastHeight <= currentBlockHeight {
		// has scan all
		return nil
	}

	currToBlockHeight := lastHeight
	s.setToHeight(currToBlockHeight)

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("scanner stoped")
			return nil
		default:
			curr := currentBlockHeight
			if curr > currToBlockHeight {
				// has to the last
				return s.scanBlocksImp(ctx, curr)
			}

			blockHash, err := s.cli.API().RPC.Chain.GetBlockHash(curr)
			if err != nil {
				return errors.Wrapf(err, "query block %d", curr)
			}

			evts, err := s.cli.QueryEventByBlockHash(blockHash)
			if err != nil {
				return errors.Wrapf(err, "query events by block %d", curr)
			}

			// to handler loop
			s.eventChann <- evts

			currentBlockHeight = curr + 1
		}
	}
}

func (s *scanner) Wait() {
	s.wg.Wait()
}
