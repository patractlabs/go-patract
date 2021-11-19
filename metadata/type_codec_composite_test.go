package metadata_test

import (
	"bytes"
	"math/big"

	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/assert"
)

func TestCompositeEncodeDecode(t *testing.T) {
	raw := loadMetaData4Test(`
{
	"V1": {
		"types": [
		{
			"id": 0,
			"type": {
				"def": {
					"primitive": "u128"
				}
			}
		},
		{
			"id": 1,
			"type": {
				"def": {
				"primitive": "bool"
				}
			}
		},
		{
			"id": 2,
			"type": {
				"def": {
				"composite": {
					"fields": [
						{
							"name": "i1",
							"type": 1
						},
						{
							"name": "i2",
							"type": 1
						},
						{
							"name": "b1",
							"type": 2
						}
					]
				}
				},
				"path": [
					"tester"
				]
				}
			}
		]
	}
}
`)

	typeDefs := make([]metadata.DefCodec, 0, 8)
	for _, ty := range raw.V1.Types {
		typeDefs = append(typeDefs, metadata.NewTypeDef(&ty))
	}

	// check encode
	toData := types.MustHexDecodeString(
		"0x000064a7b3b6e00d00000000000000000010a5d4e8000000000000000000000001")

	logger := log.NewLogger()

	bz := bytes.NewBuffer(make([]byte, 0, 64))
	encoder := scale.NewEncoder(bz)
	ctx := metadata.NewCtxForEncoder(typeDefs, encoder).WithLogger(logger)

	val := testCompos{
		I1: types.NewU128(*big.NewInt(1000000000000000000)),
		I2: types.NewU128(*big.NewInt(1000000000000)),
		B1: types.NewBool(true),
	}
	err := typeDefs[2].Encode(ctx, val)
	assert.Nil(t, err)

	decoder := scale.NewDecoder(bytes.NewReader(toData))
	ctx = metadata.NewCtxForDecoder(typeDefs, decoder).WithLogger(logger)

	res := testCompos{}

	err = typeDefs[2].Decode(ctx, &res)
	assert.Nil(t, err)
}
