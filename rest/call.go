package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/patractlabs/go-patract/utils"
)

type ChainStatus struct {
	SpecVersion        uint32 `json:"spec_version"`
	TransactionVersion uint32 `json:"tx_version"`
	BlockHash          string `json:"block_hash"`
	GenesisHash        string `json:"genesis_hash"`
}

type ExecCommonParams struct {
	Nonce       *uint64         `json:"nonce"`
	ChainStatus *ChainStatus    `json:"chain_status"`
	Contract    string          `json:"contract"`
	Origin      string          `json:"origin"`
	Value       string          `json:"value"`
	GasLimit    string          `json:"gas_limit"`
	Args        json.RawMessage `json:"args"`
}

func (r *Router) messageHandler(message *metadata.MessageRaw) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var params ExecCommonParams
		if err := ctx.ShouldBindJSON(&params); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		isOffline := ctx.Query("isOffline")
		if isOffline != "" {
			// if offline, so should gice chain status
			if params.ChainStatus == nil || params.Nonce == nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "offline sign should fill chain status and nonce"})
				return
			}
		} else {
			if r.cli == nil {
				// if online call, need connect to chain
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "online call need connect to chain to get nonce"})
				return
			}
		}

		contractIDStr := ctx.Query("contract")
		if params.Contract != "" {
			contractIDStr = params.Contract
		}

		contractID, err := utils.DecodeAccountIDFromSS58(contractIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("contract: %v\n", contractID)
		fmt.Printf("params: %s\n", params)
	}
}

func (r *Router) messages(path string) {
	for name, metadate := range r.metadatas {
		path4Contract := path + "/" + name + "/"

		for _, message := range metadate.Spec.Messages {
			path4Message := path4Contract + "exec/"
			if !message.Mutates {
				path4Message = path4Contract + "read/"
			}

			for _, n := range message.Name {
				path4Message = path4Message + n + "/"
			}

			r.POST(path4Message[:len(path4Message)-1], r.messageHandler(&message))
		}
	}
}
