package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/patractlabs/go-patract/utils"
)

type ExecParams struct {
	Origin   string          `json:"origin"`
	Value    string          `json:"value"`
	GasLimit string          `json:"gas_limit"`
	Args     json.RawMessage `json:"args"`
}

func messageHandler(message *metadata.MessageRaw) func(*gin.Context) {
	return func(ctx *gin.Context) {
		contract := ctx.Query("contract")

		var params ExecParams
		if err := ctx.ShouldBindJSON(&params); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		contractID, err := utils.DecodeAccountIDFromSS58(contract)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("contract: %v\n", contractID)
		fmt.Printf("params: %s\n", params)
	}
}

func (router *Router) messages(path string) {
	for name, metadate := range router.metadatas {
		path4Contract := path + "/" + name + "/"

		for _, message := range metadate.Spec.Messages {
			path4Message := path4Contract + "exec/"
			if !message.Mutates {
				path4Message = path4Contract + "read/"
			}

			for _, n := range message.Name {
				path4Message = path4Message + n + "/"
			}

			router.POST(path4Message, messageHandler(&message))
		}
	}
}
