package metadata

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/centrifuge/go-substrate-rpc-client/types"
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

	encoder *scale.Encoder
	decoder *scale.Decoder
	typs    []DefCodec
}

// NewCtxForEncoder new ctx for encoder
func NewCtxForEncoder(typs []DefCodec, encoder *scale.Encoder) CodecContext {
	return CodecContext{
		logger:  log.NewNopLogger(),
		typs:    typs,
		encoder: encoder,
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

// GetDefCodecByIndex get def codec by index
func (c CodecContext) GetDefCodecByIndex(i int) DefCodec {
	return c.typs[i-1]
}

// DefCodec interface to def codec
type DefCodec interface {
	Encode(ctx CodecContext, value interface{}) error
	Decode(ctx CodecContext, target interface{}) error
}

var _, _, _, _, _ DefCodec = &defPrimitive{},
	&defComposite{},
	&defArray{},
	&defVariant{},
	&defTuple{}

type defPrimitive struct {
	typ reflect.Kind
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

	if tk != d.typ {
		return errors.Errorf("type not equal, expect %v, got %v", d.typ, tk)
	}

	return ctx.encoder.Encode(value)
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

	if tk != d.typ {
		return errors.Errorf("type not equal, expect %v, got %v", d.typ, tk)
	}

	return ctx.decoder.Decode(target)
}

func getKindFromTypeString(typ string) reflect.Kind {
	switch typ {
	case "bool":
		return reflect.TypeOf(types.NewBool(false)).Kind()
	case "u8":
		return reflect.TypeOf(types.NewU8(0)).Kind()
	case "u16":
		return reflect.TypeOf(types.NewU16(0)).Kind()
	case "u32":
		return reflect.TypeOf(types.NewU32(0)).Kind()
	case "u64":
		return reflect.TypeOf(types.NewU64(0)).Kind()
	case "u128":
		return reflect.TypeOf(types.NewU128(big.Int{})).Kind()
	case "i8":
		return reflect.TypeOf(types.NewI8(0)).Kind()
	case "i16":
		return reflect.TypeOf(types.NewI16(0)).Kind()
	case "i32":
		return reflect.TypeOf(types.NewI32(0)).Kind()
	case "i64":
		return reflect.TypeOf(types.NewI64(0)).Kind()
	case "i128":
		return reflect.TypeOf(types.NewI128(big.Int{})).Kind()
	default:
		panic(errors.Errorf("unknown type by %s", typ))
	}
}

type compositeField struct {
	Name      string `json:"name"`
	TypeIndex int    `json:"type"`
}

type defComposite struct {
	Fields []compositeField `json:"fields"`
}

func newDefComposite(raw json.RawMessage) *defComposite {
	res := &defComposite{}

	if err := json.Unmarshal(raw, res); err != nil {
		panic(err)
	}

	return res
}

func (d *defComposite) Encode(ctx CodecContext, value interface{}) error {
	target := reflect.ValueOf(value)

	for idx, field := range d.Fields {
		ctx.logger.Debug("defComposite encode", "idx", idx, "field", field)

		// find field to value
		for i := 0; i < target.NumField(); i++ {
			ft := target.Type().Field(i)
			tv, ok := ft.Tag.Lookup("scale")
			if ok && tv == "-" {
				continue
			}

			if tv == field.Name {
				ctx.logger.Debug("target", "field", tv, "v", target.Field(i))

				def := ctx.GetDefCodecByIndex(field.TypeIndex)
				if err := def.Encode(ctx, target.Field(i).Interface()); err != nil {
					return errors.Wrapf(err,
						"encode composite field %s %d", field.Name, i)
				}
				break
			}
		}
	}

	return nil
}

func (d *defComposite) Decode(ctx CodecContext, value interface{}) error {
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

	for idx, field := range d.Fields {
		ctx.logger.Debug("defComposite decode", "idx", idx, "field", field)

		// find field to value
		for i := 0; i < target.NumField(); i++ {
			ft := target.Type().Field(i)
			tv, ok := ft.Tag.Lookup("scale")
			if ok && tv == "-" {
				continue
			}

			if tv == field.Name {
				ctx.logger.Debug("target", "field", tv, "v", target.Field(i))
				def := ctx.GetDefCodecByIndex(field.TypeIndex)

				fi := target.Field(i).Addr().Interface()

				if err := def.Decode(ctx, fi); err != nil {
					return errors.Wrapf(err,
						"decode composite field %s %d", field.Name, i)
				}
				break
			}
		}
	}
	return nil
}

type defArray struct {
}

func newDefArray(raw json.RawMessage) *defArray {
	res := &defArray{}
	return res
}

func (d *defArray) Encode(ctx CodecContext, value interface{}) error {
	return nil
}

func (d *defArray) Decode(ctx CodecContext, value interface{}) error {
	return nil
}

type defTuple struct {
}

func newDefTuple(raw json.RawMessage) *defTuple {
	res := &defTuple{}
	return res
}

func (d *defTuple) Encode(ctx CodecContext, value interface{}) error {
	return nil
}

func (d *defTuple) Decode(ctx CodecContext, value interface{}) error {
	return nil
}

type defVariant struct {
}

func newDefVariant(raw json.RawMessage) *defVariant {
	res := &defVariant{}
	return res
}

func (d *defVariant) Encode(ctx CodecContext, value interface{}) error {
	return nil
}

func (d *defVariant) Decode(ctx CodecContext, value interface{}) error {
	return nil
}
