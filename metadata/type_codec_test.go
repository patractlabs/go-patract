package metadata_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

type testCompos struct {
	I1 types.U128 `scale:"i1"`
	I2 types.U128 `scale:"i2"`
	B1 types.Bool `scale:"b1"`
}

func loadMetaDataTest(str string) metadata.Raw {
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

func testTypeEncodeAndDecode(
	t *testing.T,
	logger log.Logger,
	typeDefs []metadata.DefCodec,
	typeIdx int,
	val interface{},
	res interface{}) {
	require := require.New(t)
	bz := bytes.NewBuffer(make([]byte, 0, 1024))
	encoder := scale.NewEncoder(bz)
	ctx := metadata.NewCtxForEncoder(typeDefs, encoder).WithLogger(logger)

	err := typeDefs[typeIdx].Encode(ctx, val)
	require.Nil(err, "encode")

	bytesEncode := bz.Bytes()

	decoder := scale.NewDecoder(bytes.NewReader(bytesEncode))
	ctx = metadata.NewCtxForDecoder(typeDefs, decoder).WithLogger(logger)

	err = typeDefs[typeIdx].Decode(ctx, res)
	require.Nil(err, "decode")
}
