package rest

import (
	"encoding/json"
	"math/big"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/types"
	"github.com/pkg/errors"
)

type ExecCommonParams struct {
	Nonce       *uint64            `json:"nonce"`
	ChainStatus *types.ChainStatus `json:"chain_status"`
	Contract    string             `json:"contract"`
	Origin      string             `json:"origin"`
	Value       string             `json:"value"`
	GasLimit    string             `json:"gas_limit"`
	Args        json.RawMessage    `json:"args"`
}

func (r *Router) messageHandler(metaDate metadata.Data, message metadata.MessageRaw) func(*gin.Context) {
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

		contractID, err := r.ss58Codec.DecodeAccountID(contractIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var (
			value *big.Int = big.NewInt(0)
			ok    bool
		)

		if params.Value != "" {
			value, ok = big.NewInt(0).SetString(params.Value, 10)
			if !ok {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "parse value param failed"})
				return
			}
		}

		gasLimit, err := strconv.ParseUint(params.GasLimit, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.Wrap(err, "parse gasLimit param failed").Error()})
			return
		}

		var nonce uint64
		if params.Nonce != nil {
			nonce = *params.Nonce
		}

		data, err := rpc.GetMessagesDataFromJSON(&metaDate, message.Name, params.Args)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.Wrap(err, "gen msg data").Error()})
			return
		}

		bz, err := api.MakeExtrinisic(
			nonce, r.runtimeMetadata,
			params.ChainStatus, "Contracts.call",
			contractID,
			types.NewCompactBalanceByInt(value),
			types.NewCompactGas(types.Gas(gasLimit)),
			data)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		raw := types.HexEncodeToString(bz)

		ctx.JSON(http.StatusOK, gin.H{"raw": raw})
	}
}

func (r *Router) messages(path string) {
	for name, metadate := range r.metadatas {
		path4Contract := path + "/" + name + "/"

		for _, message := range metadate.V1.Spec.Messages {
			path4Message := path4Contract + "exec/"
			if !message.Mutates {
				path4Message = path4Contract + "read/"
			}

			for _, n := range message.Name {
				path4Message = path4Message + n + "/"
			}

			r.POST(path4Message[:len(path4Message)-1], r.messageHandler(*metadate, message))
		}
	}
}
