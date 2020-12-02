package metadata

import (
	"encoding/json"
	"math/big"
	"reflect"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/centrifuge/go-substrate-rpc-client/types"
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

type defCodec interface {
	Encode(encoder *scale.Encoder, value interface{}) error
}

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

func (d defPrimitive) Encode(encoder *scale.Encoder, value interface{}) error {
	// for primitive, just can encode base types, check types
	t := reflect.TypeOf(value)
	tk := t.Kind()

	if tk != d.typ {
		return errors.Errorf("type not equal, expect %v, got %v", d.typ, tk)
	}

	return encoder.Encode(value)
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
