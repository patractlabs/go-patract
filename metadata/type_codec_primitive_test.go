package metadata_test

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v2/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/stretchr/testify/assert"
)

func TestPrimitiveEncode(t *testing.T) {
	raw := loadMetaData4Test(`
  {
    "types": [
        {
            "def": {
                "primitive": "u128"
            }
        }
    ]
  }
	`)

	def := metadata.NewTypeDef(&raw.Types[0])

	bz := bytes.NewBuffer(make([]byte, 0, 64))

	// check encode
	v := big.NewInt(0).Mul(big.NewInt(1000000000000000000), big.NewInt(1000000000000000000))
	toData := types.MustHexDecodeString("0x00000000109f4bb31507c97bce97c000")

	encoder := scale.NewEncoder(bz)
	ctx := metadata.NewCtxForEncoder(nil, encoder)

	err := def.Encode(ctx, types.NewU128(*v))
	assert.Nil(t, err)
	assert.Equal(t, bz.Bytes(), toData)

	decoder := scale.NewDecoder(bytes.NewReader(toData))
	i128 := types.NewU128(*big.NewInt(0))
	ctx = metadata.NewCtxForDecoder(nil, decoder)

	err = def.Decode(ctx, &i128)
	assert.Nil(t, err)
	assert.Equal(t, i128, types.NewU128(*v))

}
