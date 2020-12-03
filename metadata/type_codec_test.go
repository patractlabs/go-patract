package metadata_test

import (
	"encoding/json"
	"io/ioutil"

	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/patractlabs/go-patract/metadata"
)

type testCompos struct {
	I1 types.U128 `scale:"i1"`
	I2 types.U128 `scale:"i2"`
	B1 types.Bool `scale:"b1"`
}

func loadMetaData4Test(str string) metadata.Raw {
	res := metadata.Raw{}

	if err := json.Unmarshal([]byte(str), &res); err != nil {
		panic(err)
	}

	return res
}

func loadMetaDataFromFile(path string) metadata.Raw {
	res := metadata.Raw{}

	bz, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(bz, &res); err != nil {
		panic(err)
	}

	return res
}
