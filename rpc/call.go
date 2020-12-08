package rpc

import (
	"bytes"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/pkg/errors"
)

const (
	// DefaultGasLimitForCall just default gas
	DefaultGasLimitForCall = 5000000000000
)

func (c *Contract) encodeDataFromArgs(argsToEncode []metadata.ArgRaw, args ...interface{}) ([]byte, error) {
	if len(argsToEncode) != len(args) {
		return nil, errors.Errorf(
			"constructor args count error, expected %d, got %d",
			len(argsToEncode), len(args))
	}

	bz := bytes.NewBuffer(make([]byte, 0, 1024))
	encoder := scale.NewEncoder(bz)
	ctx := metadata.NewCtxForEncoder(c.metaData.Codecs, encoder).WithLogger(c.logger)

	for i := 0; i < len(args); i++ {
		cdc, err := c.metaData.GetCodecByTypeIdx(argsToEncode[i].Type)
		if err != nil {
			return nil, errors.Wrapf(err, "get codec args %d", i)
		}

		if err := cdc.Encode(ctx, args[i]); err != nil {
			return nil, errors.Wrapf(err, "encode args %d", i)
		}
	}

	return bz.Bytes(), nil
}

func (c *Contract) getMessagesData(name []string, args ...interface{}) ([]byte, error) {
	message, err := c.metaData.Raw.GetMessage(name)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	_, err = buf.Write(message.SelectorData)
	if err != nil {
		return nil, errors.Wrap(err, "write selector data")
	}

	bz, err := c.encodeDataFromArgs(message.Args, args...)
	if err != nil {
		return nil, errors.Wrap(err, "encode data")
	}

	_, err = buf.Write(bz)
	if err != nil {
		return nil, errors.Wrap(err, "write encode data")
	}

	return buf.Bytes(), nil
}

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

	bz, err := c.encodeDataFromArgs(constructor.Args, args...)
	if err != nil {
		return nil, errors.Wrap(err, "encode data")
	}

	_, err = buf.Write(bz)
	if err != nil {
		return nil, errors.Wrap(err, "write encode data")
	}

	return buf.Bytes(), nil
}

// Call contract call
func (c *Contract) Call(ctx api.Context, result interface{},
	contractHash string,
	call []string,
	args ...interface{}) error {
	params := struct {
		Origin    string `json:"origin"`
		Dest      string `json:"dest"`
		GasLimit  uint   `json:"gasLimit"`
		InputData string `json:"inputData"`
		Value     int    `json:"value"`
	}{
		Origin:   ctx.From().Address,
		Dest:     contractHash,
		GasLimit: DefaultGasLimitForCall,
	}

	data, err := c.getMessagesData(call, args...)
	if err != nil {
		return errors.Wrap(err, "getMessagesData")
	}

	params.InputData = types.HexEncodeToString(data)

	c.logger.Debug("contracts call", "call", call, "hash", contractHash, "params", params)

	res := struct {
		Success struct {
			Data        string `json:"data"`
			Flags       int    `json:"flags"`
			GasConsumed int    `json:"gas_consumed"`
		} `json:"success"`
	}{}

	err = c.native.Cli.Call(&res, "contracts_call", params)
	if err != nil {
		return errors.Wrap(err, "call")
	}

	c.logger.Debug("contracts call", "call", call, "res", res)

	message, err := c.metaData.Raw.GetMessage(call)
	if err != nil {
		return err
	}

	bz, err := types.HexDecodeString(res.Success.Data)
	if err != nil {
		return errors.Wrap(err, "hex from string error")
	}

	err = c.metaData.Decode(result, message.ReturnType, bz)
	return errors.Wrap(err, "decode error")
}
