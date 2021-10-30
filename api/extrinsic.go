package api

import (
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/pkg/errors"
)

func GenSignBytes(e types.Extrinsic, o types.SignatureOptions) ([]byte, error) {
	if e.Type() != types.ExtrinsicVersion4 {
		return []byte{}, errors.Errorf(
			"unsupported extrinsic version: %v (isSigned: %v, type: %v)",
			e.Version, e.IsSigned(), e.Type())
	}

	mb, err := types.EncodeToBytes(e.Method)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "encode method error")
	}

	era := o.Era
	if !o.Era.IsMortalEra {
		era = types.ExtrinsicEra{IsImmortalEra: true}
	}

	payload := types.ExtrinsicPayloadV4{
		ExtrinsicPayloadV3: types.ExtrinsicPayloadV3{
			Method:      mb,
			Era:         era,
			Nonce:       o.Nonce,
			Tip:         o.Tip,
			SpecVersion: o.SpecVersion,
			GenesisHash: o.GenesisHash,
			BlockHash:   o.BlockHash,
		},
		TransactionVersion: o.TransactionVersion,
	}

	b, err := types.EncodeToBytes(payload)
	if err != nil {
		return []byte{}, errors.Wrap(err, "encode payload")
	}

	return b, nil
}
