package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
)

func LoadRuntimeMetadata(path string) *types.Metadata {
	res := types.NewMetadataV12()

	bz, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(bz, &res.AsMetadataV12); err != nil {
		panic(err)
	}

	return res
}
