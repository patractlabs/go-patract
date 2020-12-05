package rpc

import (
	"bytes"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/pkg/errors"
)

func (c *Contract) getConstructorsData(name []string, args ...interface{}) ([]byte, error) {
	constructor, err := c.metaData.Raw.GetConstructor(name)
	if err != nil {
		return nil, err
	}

	if len(constructor.Args) != len(args) {
		return nil, errors.Errorf(
			"constructor args count error, expected %d, got %d",
			len(constructor.Args), len(args))
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	_, err = buf.Write(constructor.SelectorData)
	if err != nil {
		return nil, errors.Wrap(err, "write selector data")
	}

	bz := bytes.NewBuffer(make([]byte, 0, 1024))
	encoder := scale.NewEncoder(bz)
	ctx := metadata.NewCtxForEncoder(c.metaData.Codecs, encoder).WithLogger(c.logger)

	for i := 0; i < len(args); i++ {
		cdc, err := c.metaData.GetCodecByTypeIdx(constructor.Args[i].Type)
		if err != nil {
			return nil, errors.Wrapf(err, "get codec args %d", i)
		}

		if err := cdc.Encode(ctx, args[i]); err != nil {
			return nil, errors.Wrapf(err, "encode args %d", i)
		}
	}

	_, err = buf.Write(bz.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "write encode data")
	}

	return buf.Bytes(), nil
}
