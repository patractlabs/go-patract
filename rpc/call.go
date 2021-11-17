package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v3/scale"
	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils/log"
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
	fmt.Println(args)
	fmt.Println(len(args))
	for i := 0; i < len(args); i++ {
		cdc, err := c.metaData.GetCodecByTypeIdx(argsToEncode[i].Type)
		if err != nil {
			return nil, errors.Wrapf(err, "get codec args %d", i)
		}
		fmt.Println("------------------------------- here")
		fmt.Println("------------------------------- here")
		fmt.Println(argsToEncode[i].Type)
		fmt.Println("------------------------------- here")
		fmt.Println("------------------------------- here")
		fmt.Println(args[i])
		//err = cdc.Encode(ctx, args[i])
		//fmt.Println(args[i])
		if err := cdc.Encode(ctx, args[i]); err != nil {
			return nil, errors.Wrapf(err, "encode args %d", i)
		}
	}

	return bz.Bytes(), nil
}

func encodeDataFromArgJSONs(
	metaData *metadata.Data,
	argsToEncode []metadata.ArgRaw,
	args ...json.RawMessage) ([]byte, error) {
	if len(argsToEncode) != len(args) {
		return nil, errors.Errorf(
			"constructor args count error, expected %d, got %d",
			len(argsToEncode), len(args))
	}

	bz := bytes.NewBuffer(make([]byte, 0, 1024))
	encoder := scale.NewEncoder(bz)
	ctx := metadata.NewCtxForEncoder(metaData.Codecs, encoder).WithLogger(log.NewLogger())

	for i := 0; i < len(args); i++ {
		cdc, err := metaData.GetCodecByTypeIdx(argsToEncode[i].Type)
		if err != nil {
			return nil, errors.Wrapf(err, "get codec args %d", i)
		}

		if err := cdc.EncodeJSON(ctx, args[i]); err != nil {
			return nil, errors.Wrapf(err, "encode args %d", i)
		}
	}

	return bz.Bytes(), nil
}

func (c *Contract) GetMessageData(name []string, args ...interface{}) ([]byte, error) {
	return c.getMessagesData(name, args...)
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

func GetMessagesDataFromJSON(metaData *metadata.Data, name []string, args json.RawMessage) ([]byte, error) {
	message, err := metaData.Raw.GetMessage(name)
	if err != nil {
		return nil, err
	}

	logger := log.NewLogger()

	logger.Debug("GetMessagesDataFromJSON", "name", name)

	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	_, err = buf.Write(message.SelectorData)
	if err != nil {
		return nil, errors.Wrap(err, "write selector data")
	}

	argsMap := make(map[string]json.RawMessage)
	if err := json.Unmarshal(args, &argsMap); err != nil {
		return nil, errors.Wrap(err, "unmarshal params error")
	}

	argArr := make([]json.RawMessage, 0, len(argsMap)+1)
	for _, arg := range message.Args {
		a, ok := argsMap[arg.Name]
		if !ok {
			return nil, errors.Errorf("no params for %s", arg.Name)
		}

		argArr = append(argArr, a)
	}

	logger.Info("GetMessagesDataFromJSON", "arr", argArr)

	bz, err := encodeDataFromArgJSONs(metaData, message.Args, argArr...)
	if err != nil {
		return nil, errors.Wrap(err, "encode data")
	}

	_, err = buf.Write(bz)
	if err != nil {
		return nil, errors.Wrap(err, "write encode data")
	}

	return buf.Bytes(), nil
}

// GenConstructorsData gen input data for constructor
func (c *Contract) GenConstructorsData(name []string, args ...interface{}) ([]byte, error) {
	return c.getConstructorsData(name, args...)
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

// CallToRead contract call to read state from chain
func (c *Contract) CallToRead(
	ctx api.Context,
	result interface{},
	contractID types.AccountID,
	call []string,
	args ...interface{}) error {
	contractIDStr, err := c.ss58Codec.EncodeAccountID(contractID)
	if err != nil {
		return errors.Wrapf(err, "encode accountid for contract %v", contractID)
	}

	params := struct {
		Origin    string `json:"origin"`
		Dest      string `json:"dest"`
		GasLimit  uint   `json:"gasLimit"`
		InputData string `json:"inputData"`
		Value     int    `json:"value"`
	}{
		Origin:   ctx.From().Address,
		Dest:     contractIDStr,
		GasLimit: DefaultGasLimitForCall,
	}

	data, err := c.getMessagesData(call, args...)
	if err != nil {
		return errors.Wrap(err, "getMessagesData")
	}

	params.InputData = types.HexEncodeToString(data)

	c.logger.Debug("contracts call", "call", call, "hash", contractIDStr, "params", params)

	res := struct {
		DebugMessage string `json:"debugMessage"`
		GasConsumed  int    `json:"gasConsumed"`
		Result       struct {
			Ok struct {
				Data  string `json:"data"`
				Flags int    `json:"flags"`
			} `json:"Ok"`
		} `json:"result"`
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

	bz, err := types.HexDecodeString(res.Result.Ok.Data)
	if err != nil {
		return errors.Wrap(err, "hex from string error")
	}

	if len(bz) == 0 {
		return errors.Errorf("no data got")
	}

	err = c.metaData.Decode(result, message.ReturnType, bz)
	return errors.Wrapf(err, "decode error %s.", res.Result.Ok.Data)
}

// CallToExec contract call to exec state from chain
func (c *Contract) CallToExec(
	ctx api.Context,
	contractID types.AccountID,
	value types.CompactBalance,
	gasLimit types.CompactGas,
	call []string, args ...interface{}) (types.Hash, error) {
	data, err := c.getMessagesData(call, args...)
	if err != nil {
		return types.Hash{}, errors.Wrap(err, "getMessagesData")
	}

	c.logger.Debug("call to exec", "data", types.HexEncodeToString(data))

	return c.native.Call(ctx, contractID, value, gasLimit, data)
}
