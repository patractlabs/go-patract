package main

import (
	"flag"

	"github.com/patractlabs/go-patract/utils"
	"github.com/patractlabs/go-patract/utils/log"
)

var (
	flagNodeURL          = flag.String("ws", "ws://localhost:9944", "url to chain node")
	flagContractMetadata = flag.String("metadata", "", "metadata for contract")
	flagCodeHash         = flag.String("hash", "", "codeHash for contract")
	flagDBPath           = flag.String("path", "./erc20.db", "path to database")
	flagRestURL          = flag.String("url", ":8899", "path to url")
)

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}

	logger := log.NewLogger()
	db := NewErc20DB(*flagDBPath)

	defer func() {
		db.Close()
		logger.Flush()
	}()

	StartRestServer(logger, db)

	o, cancelFunc, err := observerEvts(logger, db)
	if err != nil {
		panic(err)
	}

	utils.HoldToClose(func() {
		cancelFunc()
		o.Wait()

		logger.Info("observer stop")
		logger.Flush()
	})
}
