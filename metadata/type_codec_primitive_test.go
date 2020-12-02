package metadata_test

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/centrifuge/go-substrate-rpc-client/types"
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
	encoder := scale.NewEncoder(bz)

	v := big.NewInt(0).Mul(big.NewInt(1000000000000000000), big.NewInt(1000000000000000000))
	err := def.Encode(encoder, types.NewU128(*v))
	assert.Nil(t, err)

	toData := types.MustHexDecodeString("0x00000000109f4bb31507c97bce97c000")

	assert.Equal(t, bz.Bytes(), toData)
}
