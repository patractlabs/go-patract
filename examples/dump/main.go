package main

import (
	"flag"

	"github.com/patractlabs/go-patract/utils"
	"github.com/patractlabs/go-patract/utils/log"
)

var (
	flagURL              = flag.String("url", "ws://localhost:9944", "url to chain node")
	flagContractMetadata = flag.String("metadata", "", "metadata for contract")
	flagCodeHash         = flag.String("hash", "", "codeHash for contract")
	flagDBPath           = flag.String("path", "./erc20.db", "path to database")
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
