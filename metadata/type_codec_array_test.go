package metadata_test

import (
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

func TestDefArrayEncodeAndDecode(t *testing.T) {
	raw := loadMetaData4Test(`
{
    "types": [
        {
            "def": {
                "primitive": "u8"
            }
        },
        {
            "def": {
                "primitive": "u128"
            }
        },
        {
            "def": {
                "primitive": "bool"
            }
        },
        {
            "def": {
                "composite": {
                    "fields": [
                        {
                            "name": "i1",
                            "type": 2
                        },
                        {
                            "name": "i2",
                            "type": 2
                        },
                        {
                            "name": "b1",
                            "type": 3
                        }
                    ]
                }
            },
            "path": [
                "tester"
            ]
        },
        {
            "def": {
                "array": {
                    "len": 8,
                    "type": 1
                }
            }
        },
        {
            "def": {
                "array": {
                    "len": 8,
                    "type": 2
                }
            }
        },
        {
            "def": {
                "array": {
                    "len": 8,
                    "type": 4
                }
            }
        }
    ]
}
	`)

	typeDefs := make([]metadata.DefCodec, 0, 16)
	for _, ty := range raw.Types {
		typeDefs = append(typeDefs, metadata.NewTypeDef(&ty))
	}

	logger := log.NewLogger()

	v1 := [8]byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7}
	v1t := [8]byte{}
	testTypeEncodeAndDecode(t, logger, typeDefs, 4, v1, &v1t)
	require.Equal(t, v1, v1t, "encode and decode []byte should match")

	v2 := [8]types.U128{
		types.NewU128(*big.NewInt(1)),
		types.NewU128(*big.NewInt(2)),
		types.NewU128(*big.NewInt(3)),
		types.NewU128(*big.NewInt(4)),
		types.NewU128(*big.NewInt(5)),
		types.NewU128(*big.NewInt(6)),
		types.NewU128(*big.NewInt(7)),
		types.NewU128(*big.NewInt(8)),
	}
	v2t := [8]types.U128{}
	testTypeEncodeAndDecode(t, logger, typeDefs, 5, v2, &v2t)
	require.Equal(t, v2, v2t, "encode and decode []types.U128 should match")

	val := testCompos{
		I1: types.NewU128(*big.NewInt(1000000000000000000)),
		I2: types.NewU128(*big.NewInt(1000000000000)),
		B1: types.NewBool(true),
	}

	v3 := [8]testCompos{
		val,
		val,
		val,
		val,
		val,
		val,
		val,
		val,
	}
	v3t := [8]testCompos{}
	testTypeEncodeAndDecode(t, logger, typeDefs, 6, v3, &v3t)
	require.Equal(t, v3, v3t, "encode and decode []testCompos should match")
}
