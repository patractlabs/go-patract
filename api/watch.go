package api

import (
	"context"
	"sync"

	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
	"github.com/patractlabs/go-patract/utils/log"
)

type Watcher struct {
	wg         sync.WaitGroup
	eventChann chan evtMsgInChann
	stat       int
	mutex      sync.RWMutex
	scanner    *scanner
	cli        *Client
	logger     log.Logger
}

func NewWatcher(logger log.Logger, url string) *Watcher {
	scanner := NewScanner(logger, url)

	return &Watcher{
		scanner:    scanner,
		cli:        scanner.Cli(),
		logger:     logger,
		eventChann: make(chan evtMsgInChann, 4096),
	}
}

func (s *Watcher) Cli() *Client {
	return s.cli
}

func (w *Watcher) nextStatStep() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.stat++
}

func (w *Watcher) Status() int {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	return w.stat
}

func (w *Watcher) Wait() {
	w.logger.Debug("watcher start wait stopped")

	w.wg.Wait()

	w.logger.Debug("watcher stopped")
}

func (w *Watcher) Watch(ctx context.Context, fromHeight uint64, h EventHandler) error {
	w.logger.Debug("start scanner first", "from", fromHeight)

	// init the client

	scannerCtx, cancelScanner := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()

		// first stop scanner
		cancelScanner()
		if w.scanner != nil {
			w.scanner.Wait()
		}

		// second, stop watcher

		w.stop()
	}()

	// start handler
	w.logger.Debug("start handler")

	// call all handler h in one gorountinue
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()

		for {
			evts, ok := <-w.eventChann
			if !ok {
				// closed scanner, closed by last
				w.logger.Info("stop handler gorountinue")
				return
			}

			if err := h(w.logger, evts.height, evts.records); err != nil {
				w.logger.Error("handler error", "err", err.Error())
			}
		}
	}()

	// from init to sync
	w.nextStatStep()

	// first scanner all old blocks
	w.scanner.Scan(scannerCtx, fromHeight,
		func(logger log.Logger, height uint64, records *types.EventRecords) error {
			w.eventChann <- evtMsgInChann{
				height:  height,
				records: records,
			}
			return nil
		})

	// start watch blocks from no scaned
	currentHasHandled := w.scanner.LastestBlockHeight()
	w.logger.Debug("current block has get", "height", currentHasHandled)

	// from sync to watch
	w.nextStatStep()

	// on watch
	w.cli.WatchEvents(ctx, func(logger log.Logger, height uint64, records *types.EventRecords) error {
		logger.Debug("handler block events", "height", height)

		w.eventChann <- evtMsgInChann{
			height:  height,
			records: records,
		}
		return nil
	})

	return nil
}

func (w *Watcher) stop() {
	w.logger.Debug("watcher start stop")
	close(w.eventChann)
}
