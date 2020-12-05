package metadata

import (
	"bytes"
	"encoding/json"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/pkg/errors"
)

// Data meta data for contracts
type Data struct {
	Raw

	Codecs []DefCodec
}

// New create metadata from metadata.json
func New(bz []byte) (*Data, error) {
	res := &Data{}

	if err := json.Unmarshal(bz, &res.Raw); err != nil {
		return nil, errors.Wrap(err, "unmarshal json")
	}

	res.Codecs = make([]DefCodec, 0, len(res.Raw.Types))
	for _, ty := range res.Raw.Types {
		res.Codecs = append(res.Codecs, NewTypeDef(&ty))
	}

	// parse datas
	for i := 0; i < len(res.Raw.Spec.Constructors); i++ {
		selectorStr := res.Raw.Spec.Constructors[i].Selector

		bz, err := types.HexDecodeString(selectorStr)
		if err != nil {
			return nil, errors.Wrapf(err, "decode str selector from %s", selectorStr)
		}

		res.Raw.Spec.Constructors[i].SelectorData = bz
	}

	for i := 0; i < len(res.Raw.Spec.Messages); i++ {
		selectorStr := res.Raw.Spec.Messages[i].Selector

		bz, err := types.HexDecodeString(selectorStr)
		if err != nil {
			return nil, errors.Wrapf(err, "decode str selector from %s", selectorStr)
		}

		res.Raw.Spec.Messages[i].SelectorData = bz
	}

	return res, nil
}

func (d *Data) GetCodecByTypeIdx(i TypeIndex) (DefCodec, error) {
	if len(d.Codecs) < i.Type {
		return nil, errors.Errorf("codec idx no found to %d, all len %d",
			i.Type, len(d.Codecs))
	}

	return d.Codecs[i.Type-1], nil
}

// NewCtxForDecode new ctx for decoder
func (d *Data) NewCtxForDecode(bz []byte) CodecContext {
	decoder := scale.NewDecoder(bytes.NewReader(bz))
	return CodecContext{
		logger:  log.NewNopLogger(),
		typs:    d.Codecs,
		decoder: decoder,
	}
}

func (d *Data) Decode(res interface{}, typ TypeIndex, bz []byte) error {
	ctx := d.NewCtxForDecode(bz)

	codec, err := d.GetCodecByTypeIdx(typ)
	if err != nil {
		return err
	}

	return codec.Decode(ctx, res)
}
