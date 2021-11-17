package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/centrifuge/go-substrate-rpc-client/v3/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
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
	fmt.Println("1____________________________________________")

	if err := json.Unmarshal(bz, &res.Raw); err != nil {
		fmt.Println(err)
		return nil, errors.Wrap(err, "unmarshal json")
	}
	fmt.Println("2____________________________________________")

	res.Codecs = make([]DefCodec, 0, len(res.Raw.V1.Types))

	for idx := range res.Raw.V1.Types {
		res.Codecs = append(res.Codecs, NewTypeDef(&res.Raw.V1.Types[idx]))
	}

	// parse datas
	fmt.Println("3____________________________________________")
	fmt.Println("-____________________________________________")
	fmt.Println("-____________________________________________")
	fmt.Println("-____________________________________________")
	//fmt.Println(string(types.HexDecodeString(res.Raw.V1.Spec.Constructors[0].Selector)))
	fmt.Println(len(res.Raw.V1.Spec.Constructors))
	for i := 0; i < len(res.Raw.V1.Spec.Constructors); i++ {
		selectorStr := res.Raw.V1.Spec.Constructors[i].Selector

		bz, err := types.HexDecodeString(selectorStr)
		fmt.Println(bz)
		fmt.Println("------------------------------------------")
		fmt.Println("------------------------------------------")
		if err != nil {
			return nil, errors.Wrapf(err, "decode str selector from %s", selectorStr)
		}

		res.Raw.V1.Spec.Constructors[i].SelectorData = bz
	}

	for i := 0; i < len(res.Raw.V1.Spec.Messages); i++ {
		selectorStr := res.Raw.V1.Spec.Messages[i].Selector

		bz, err := types.HexDecodeString(selectorStr)
		if err != nil {
			return nil, errors.Wrapf(err, "decode str selector from %s", selectorStr)
		}

		res.Raw.V1.Spec.Messages[i].SelectorData = bz
	}

	return res, nil
}

// NewFromFile creates a new metadata data from read metadata.json
func NewFromFile(path string) (*Data, error) {
	bz, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "read file %s", path)
	}

	return New(bz)
}

func (d *Data) GetCodecByTypeIdx(i TypeIndex) (DefCodec, error) {
	if len(d.Codecs) < i.Type {
		return nil, errors.Errorf("codec idx no found to %d, all len %d",
			i.Type, len(d.Codecs))
	}

	return d.Codecs[i.Type], nil
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
