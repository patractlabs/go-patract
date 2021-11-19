package metadata

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v3/scale"

	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/pkg/errors"
)

const (
	// DefTypComposite composite type def
	DefTypComposite = "composite"

	// DefTypPrimitive primitive type def
	DefTypPrimitive = "primitive"

	// DefTypVariant variant type def
	DefTypVariant = "variant"

	// DefTypArray array type def
	DefTypArray = "array"

	// DefTypTuple tuple type def
	DefTypTuple = "tuple"
)

// CodecContext context for codec by type def
type CodecContext struct {
	logger log.Logger

	encoder   *scale.Encoder
	decoder   *scale.Decoder
	ss58Codec *types.SS58Codec
	typs      []DefCodec
}

// NewCtxForEncoder new ctx for encoder
func NewCtxForEncoder(typs []DefCodec, encoder *scale.Encoder) CodecContext {
	return CodecContext{
		logger:    log.NewNopLogger(),
		typs:      typs,
		encoder:   encoder,
		ss58Codec: types.GetDefaultSS58Codec(),
	}
}

// NewCtxForDecoder new ctx for decoder
func NewCtxForDecoder(typs []DefCodec, decoder *scale.Decoder) CodecContext {
	return CodecContext{
		logger:  log.NewNopLogger(),
		typs:    typs,
		decoder: decoder,
	}
}

// WithLogger with logger for ctx
func (c CodecContext) WithLogger(logger log.Logger) CodecContext {
	c.logger = logger
	return c
}

// WithSS58Codec with ss58Codec for ctx
func (c CodecContext) WithSS58Codec(ss58Codec *types.SS58Codec) CodecContext {
	c.ss58Codec = ss58Codec
	return c
}

// GetDefCodecByIndex get def codec by index
func (c CodecContext) GetDefCodecByIndex(i int) DefCodec {
	return c.typs[i]
}

// GetSS58Codec get ss58 codec
func (c CodecContext) GetSS58Codec() *types.SS58Codec {
	return c.ss58Codec
}

// DefCodec interface to def codec
type DefCodec interface {
	Encode(ctx CodecContext, value interface{}) error
	Decode(ctx CodecContext, target interface{}) error
	EncodeJSON(ctx CodecContext, value json.RawMessage) error
}

var _, _, _, _, _ DefCodec = &defPrimitive{},
	&defComposite{},
	&defArray{},
	&defVariant{},
	&defTuple{}

type defPrimitive struct {
	typ reflect.Type
}

func newDefPrimitive(raw json.RawMessage) *defPrimitive {
	var typ string

	if err := json.Unmarshal(raw, &typ); err != nil {
		panic(errors.Wrapf(err, "failed to unmarshal primitive %s", string(raw)))
	}

	return &defPrimitive{
		typ: getKindFromTypeString(typ),
	}
}

func (d *defPrimitive) Encode(ctx CodecContext, value interface{}) error {
	// for primitive, just can encode base types, check types
	t := reflect.TypeOf(value)
	tk := t.Kind()

	if tk != d.typ.Kind() {
		return errors.Errorf("type not equal, expect %v, got %v", d.typ, tk)
	}
	return ctx.encoder.Encode(value)
}

func (d *defPrimitive) EncodeJSON(ctx CodecContext, value json.RawMessage) error {
	var str string
	if err := json.Unmarshal(value, &str); err != nil {
		return err
	}

	str = strings.ToLower(str)

	val, err := newFromKind(d.typ, str)
	if err != nil {
		return errors.Wrapf(err, "new from type by %s to %s", str, d.typ.String())
	}

	ctx.logger.Debug("EncodeJSON", "value", string(value), "val", val)

	return d.Encode(ctx, val)
}

func (d *defPrimitive) Decode(ctx CodecContext, target interface{}) error {
	// for primitive, just can encode base types, check types
	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Ptr {
		return errors.Errorf("Target must be a pointer, but was %s", fmt.Sprint(t))
	}

	val := reflect.ValueOf(target)
	if val.IsNil() {
		return errors.New("Target is a nil pointer")
	}

	tk := val.Elem().Kind()

	if tk != d.typ.Kind() {
		return errors.Errorf("type not equal, expect %v, got %v", d.typ, tk)
	}

	return ctx.decoder.Decode(target)
}

func getKindFromTypeString(typ string) reflect.Type {
	switch typ {
	case "bool":
		return reflect.TypeOf(types.NewBool(false))
	case "u8":
		return reflect.TypeOf(types.NewU8(0))
	case "u16":
		return reflect.TypeOf(types.NewU16(0))
	case "u32":
		return reflect.TypeOf(types.NewU32(0))
	case "u64":
		return reflect.TypeOf(types.NewU64(0))
	case "u128":
		return reflect.TypeOf(types.NewU128(big.Int{}))
	case "i8":
		return reflect.TypeOf(types.NewI8(0))
	case "i16":
		return reflect.TypeOf(types.NewI16(0))
	case "i32":
		return reflect.TypeOf(types.NewI32(0))
	case "i64":
		return reflect.TypeOf(types.NewI64(0))
	case "i128":
		return reflect.TypeOf(types.NewI128(big.Int{}))
	default:
		panic(errors.Errorf("unknown type by %s", typ))
	}
}

func newFromKind(typ reflect.Type, str string) (interface{}, error) {
	switch typ {
	case reflect.TypeOf(types.NewBool(false)):
		return types.NewBool(strings.ToLower(str) == "true"), nil

	case reflect.TypeOf(types.NewU8(0)):
		i, err := strconv.ParseUint(str, 10, 8)
		if err != nil {
			return nil, errors.Wrapf(err, "parse u8 error")
		}
		return types.NewU8(uint8(i)), nil

	case reflect.TypeOf(types.NewU16(0)):
		i, err := strconv.ParseUint(str, 10, 16)
		if err != nil {
			return nil, errors.Wrapf(err, "parse u16 error")
		}
		return types.NewU16(uint16(i)), nil

	case reflect.TypeOf(types.NewU32(0)):
		i, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, "parse u32 error")
		}
		return types.NewU32(uint32(i)), nil

	case reflect.TypeOf(types.NewU64(0)):
		i, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse u64 error")
		}
		return types.NewU64(i), nil

	case reflect.TypeOf(types.NewU128(big.Int{})):
		i := big.NewInt(0)
		i, ok := i.SetString(str, 10)
		if !ok {
			return nil, errors.Errorf("big int set string error %s", str)
		}
		return types.NewU128(*i), nil

	case reflect.TypeOf(types.NewI8(0)):
		i, err := strconv.ParseInt(str, 10, 8)
		if err != nil {
			return nil, errors.Wrapf(err, "parse i8 error")
		}
		return types.NewI8(int8(i)), nil

	case reflect.TypeOf(types.NewI16(0)):
		i, err := strconv.ParseInt(str, 10, 16)
		if err != nil {
			return nil, errors.Wrapf(err, "parse i16 error")
		}
		return types.NewI16(int16(i)), nil

	case reflect.TypeOf(types.NewI32(0)):
		i, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, "parse i32 error")
		}
		return types.NewI32(int32(i)), nil

	case reflect.TypeOf(types.NewI64(0)):
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse i64 error")
		}
		return types.NewI64(i), nil

	case reflect.TypeOf(types.NewI128(big.Int{})):
		i := big.NewInt(0)
		i, ok := i.SetString(str, 10)
		if !ok {
			return nil, errors.Errorf("big int set string error %s", str)
		}
		return types.NewI128(*i), nil

	default:
		panic(errors.Errorf("unknown type by %s", typ))
	}
}

type defArray struct {
	TypeIndex int `json:"type"`
	Len       int `json:"len"`
}

func newDefArray(raw json.RawMessage) *defArray {
	res := &defArray{}

	if err := json.Unmarshal(raw, res); err != nil {
		panic(err)
	}

	return res
}

func (d *defArray) Encode(ctx CodecContext, value interface{}) error {
	// if value not len enough, append nil value
	t := reflect.ValueOf(value)
	tk := t.Kind()

	if tk != reflect.Array && tk != reflect.Slice {
		return errors.Errorf("def array need value is array type by %v %v", tk, value)
	}

	defCodec := ctx.GetDefCodecByIndex(d.TypeIndex)

	l := t.Len()
	if l > d.Len {
		return errors.Errorf("value len larger than def by %d > %d", l, d.Len)
	}

	ctx.logger.Debug("encode array", "val", value, "tk", tk)

	for i := 0; i < l && i < d.Len; i++ {
		if err := defCodec.Encode(ctx, t.Index(i).Interface()); err != nil {
			return errors.Wrapf(err, "failed encode %d", i)
		}
	}

	tNil := reflect.New(reflect.TypeOf(value))
	for i := l; i < d.Len; i++ {
		if err := defCodec.Encode(ctx, tNil.Interface()); err != nil {
			return errors.Wrapf(err, "failed encode with nil %d", i)
		}
	}

	return nil
}

func (d *defArray) encodeByteJSON(ctx CodecContext, value json.RawMessage) error {
	ctx.logger.Debug("encode bytes", "v", value)

	bytes, err := types.HexDecodeString(string(value))
	if err != nil {
		return errors.Wrap(err, "failed to decode hex")
	}

	if d.Len != len(bytes) {
		return errors.Wrapf(err, "bytes len not equal export %d to %d", d.Len, len(bytes))
	}

	return ctx.encoder.Encode(bytes)
}

func (d *defArray) EncodeJSON(ctx CodecContext, value json.RawMessage) error {
	defCodec := ctx.GetDefCodecByIndex(d.TypeIndex)
	// for bytes
	if pCodec, ok := defCodec.(*defPrimitive); ok {
		if pCodec.typ.String() == reflect.TypeOf(types.NewU8(0)).String() {
			return d.encodeByteJSON(ctx, value)
		}
	}

	arr := make([]json.RawMessage, 0, 32)
	if err := json.Unmarshal(value, &arr); err != nil {
		return errors.Wrap(err, "json unmarshal error")
	}

	l := len(arr)
	if l != d.Len {
		return errors.Errorf("value len larger than def by %d > %d", l, d.Len)
	}

	ctx.logger.Debug("encode array", "val", value)

	for i := 0; i < l && i < d.Len; i++ {
		if err := defCodec.EncodeJSON(ctx, arr[i]); err != nil {
			return errors.Wrapf(err, "failed encode %d", i)
		}
	}

	return nil
}

func (d *defArray) Decode(ctx CodecContext, value interface{}) error {
	t0 := reflect.TypeOf(value)
	if t0.Kind() != reflect.Ptr {
		return errors.New("Target must be a pointer, but was " + fmt.Sprint(t0))
	}

	val := reflect.ValueOf(value)
	if val.IsNil() {
		return errors.New("Target is a nil pointer")
	}

	target := val.Elem()

	t := target.Type()
	if !target.CanSet() {
		return errors.Errorf("Unsettable value %v", t)
	}

	defCodec := ctx.GetDefCodecByIndex(d.TypeIndex)

	ctx.logger.Debug("decode target", "val", val, "target", target)

	targetLen := target.Len()
	for i := 0; i < targetLen; i++ {
		fi := target.Index(i).Addr().Interface()
		if err := defCodec.Decode(ctx, fi); err != nil {
			return errors.Wrapf(err, "failed decode %d", i)
		}
	}

	return nil
}

type defTuple struct {
	TypeIndexs []int
}

func newDefTuple(raw json.RawMessage) *defTuple {
	res := &defTuple{}

	if err := json.Unmarshal(raw, &res.TypeIndexs); err != nil {
		panic(err)
	}

	return res
}

func (d *defTuple) Encode(ctx CodecContext, value interface{}) error {
	if len(d.TypeIndexs) == 0 {
		return nil
	}

	// if value not len enough, append nil value
	t := reflect.ValueOf(value)

	if len(d.TypeIndexs) != t.NumField() {
		return errors.Errorf("tuple fields count not equal by %d, expect %d",
			t.NumField(), len(d.TypeIndexs))
	}

	for i, typ := range d.TypeIndexs {
		defCodec := ctx.GetDefCodecByIndex(typ)
		f := t.Field(i).Interface()
		if err := defCodec.Encode(ctx, f); err != nil {
			return errors.Wrapf(err, "failed encode %d", i)
		}
	}

	return nil
}

func (d *defTuple) EncodeJSON(ctx CodecContext, value json.RawMessage) error {
	if len(d.TypeIndexs) == 0 {
		return nil
	}

	arr := make([]json.RawMessage, 0, 32)
	if err := json.Unmarshal(value, &arr); err != nil {
		return errors.Wrap(err, "json unmarshal error")
	}

	if len(d.TypeIndexs) != len(arr) {
		return errors.Errorf("tuple fields count not equal by %d, expect %d",
			len(arr), len(d.TypeIndexs))
	}

	for i, typ := range d.TypeIndexs {
		defCodec := ctx.GetDefCodecByIndex(typ)
		f := arr[i]
		if err := defCodec.EncodeJSON(ctx, f); err != nil {
			return errors.Wrapf(err, "failed encode %d", i)
		}
	}

	return nil
}

func (d *defTuple) Decode(ctx CodecContext, value interface{}) error {
	t0 := reflect.TypeOf(value)
	if t0.Kind() != reflect.Ptr {
		return errors.New("Target must be a pointer, but was " + fmt.Sprint(t0))
	}

	val := reflect.ValueOf(value)
	if val.IsNil() {
		return errors.New("Target is a nil pointer")
	}

	target := val.Elem()

	t := target.Type()
	if !target.CanSet() {
		return errors.Errorf("Unsettable value %v", t)
	}

	if len(d.TypeIndexs) != target.NumField() {
		return errors.Errorf("tuple field count not equal by %d, expect %d",
			target.NumField(), len(d.TypeIndexs))
	}

	for i, typ := range d.TypeIndexs {
		defCodec := ctx.GetDefCodecByIndex(typ)
		fi := target.Field(i).Addr().Interface()
		if err := defCodec.Decode(ctx, fi); err != nil {
			return errors.Wrapf(err, "failed decode %d", i)
		}
	}

	return nil
}
