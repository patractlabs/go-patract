package main

import (
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils/log"
)

func StartRestServer(logger log.Logger, db *erc20DB) {
	go func() {
		router := gin.Default()

		logger.Info("Starting rest server")

		router.GET("/get/:address", func(ctx *gin.Context) {
			address := ctx.Param("address")

			logger.Info("get address", "add", address)

			db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket(db.BucketKey)

				res := b.Get([]byte(address))

				logger.Debug("res", "data", string(res))

				if res != nil {
					ctx.String(http.StatusOK, string(res))
					return nil
				}

				return nil
			})

			ctx.String(http.StatusOK, "")
		})

		router.Run(*flagRestURL)
	}()
}

type Router struct {
	*gin.Engine

	cli             *rpc.Contract
	runtimeMetadata *types.Metadata
	metadatas       map[string]*metadata.Data
}

func NewRouter(router *gin.Engine) Router {
	return Router{
		Engine:    router,
		metadatas: make(map[string]*metadata.Data, 16),
	}
}
