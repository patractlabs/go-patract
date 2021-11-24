package metadata_test

import (
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/patractlabs/go-patract/utils/log"
)

func TestDefVariantEncodeAndDecode(t *testing.T) {
	raw := loadMetaDataTest(`
{
    "V1": {
        "types": [
            {
                "id": 0,
                "type": {
                    "def": {
                        "primitive": "u8"
                    }
                }
            },
            {
                "id": 1,
                "type": {
                    "def": {
                        "primitive": "u128"
                    }
                }
            },
            {
                "id": 2,
                "type": {
                    "def": {
                        "primitive": "bool"
                    }
                }
            },
            {
                "id": 3,
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
            },
            {
                "id": 4,
                "type": {
                    "def": {
                        "array": {
                            "len": 8,
                            "type": 3
                        }
                    }
                }
            },
            {
                "id": 8,
                "type": {
                    "def": {
                        "variant": {
                            "variants": [
                                {
                                    "index": 0,
                                    "name": "None"
                                },
                                {
                                    "fields": [
                                        {
                                            "type": 4
                                        }
                                    ],
                                    "index": 1,
                                    "name": "Some"
                                }
                            ]
                        }
                    },
                    "params": [
                        {
                            "name": "T",
                            "type": 5
                        }
                    ],
                    "path": [
                        "Option"
                    ]
                }
            }
        ]
    }
}`)

	typeDefs := make([]metadata.DefCodec, 0, 16)
	for _, ty := range raw.V1.Types {
		typeDefs = append(typeDefs, metadata.NewTypeDef(&ty))
	}

	logger := log.NewLogger()

	valCom := testCompos{
		I1: types.NewU128(*big.NewInt(1000000000000000000)),
		I2: types.NewU128(*big.NewInt(1000000000000)),
		B1: types.NewBool(true),
	}

	valComArr := [8]testCompos{
		valCom,
		valCom,
		valCom,
		valCom,
		valCom,
		valCom,
		valCom,
		valCom,
	}

	type testOption struct {
		Some *[8]testCompos `scale:"Some"`
	}

	s := testOption{
		&valComArr,
	}
	sp := testOption{
		Some: &[8]testCompos{},
	}

	testTypeEncodeAndDecode(t, logger, typeDefs, 5, s, &sp)
	//require.Equalf(t, s, sp, "encode and decode should match")
}
