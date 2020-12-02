package metadata_test

import (
	"encoding/json"

	"github.com/patractlabs/go-patract/metadata"
)

func loadMetaData4Test(str string) metadata.Raw {
	res := metadata.Raw{}

	if err := json.Unmarshal([]byte(str), &res); err != nil {
		panic(err)
	}

	return res
}
