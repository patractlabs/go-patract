package metadata

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

// rawTypeDef type definition
type rawTypeDef struct {
	Id   int `json:"id"`
	Type struct {
		Def    map[string]json.RawMessage `json:"def"`
		Params []struct {
			Name string `json:"name"`
			Type int    `json:"type"`
		} `json:"params"`
		Path []string `json:"path"`
	} `json:"type"`
}

// TypeDef type definition
type TypeDef struct {
	def    DefCodec
	Params []struct {
		Name string `json:"name"`
		Type int    `json:"type"`
	}
	Path []string
}

func NewTypeDef(raw *rawTypeDef) *TypeDef {
	if len(raw.Type.Def) != 1 {
		panic(errors.Errorf("type def raw error by not key %v", raw.Type.Def))
	}

	res := &TypeDef{
		Params: raw.Type.Params,
		Path:   raw.Type.Path,
	}

	for k, jsonRaw := range raw.Type.Def {
		switch k {
		case "primitive":
			res.def = newDefPrimitive(jsonRaw)
		case "composite":
			res.def = newDefComposite(jsonRaw, raw.Type.Path)
		case "array":
			res.def = newDefArray(jsonRaw)
		case "variant":
			res.def = newDefVariant(jsonRaw)
		case "tuple":
			res.def = newDefTuple(jsonRaw)
		default:
			panic(errors.Errorf("type def raw error by key %v not expect", k))
		}
	}

	return res
}

func (t *TypeDef) Encode(ctx CodecContext, v interface{}) error {
	//fmt.Println("--------------------------------------- 1 here is the")
	//fmt.Println("--------------------------------------- 1 here is the")
	//fmt.Println("--------------------------------------- 1 here is the")
	//fmt.Println("--------------------------------------- 1 here is the")
	//fmt.Println()
	return t.def.Encode(ctx, v)
}

func (t *TypeDef) EncodeJSON(ctx CodecContext, v json.RawMessage) error {
	return t.def.EncodeJSON(ctx, v)
}

func (t *TypeDef) Decode(ctx CodecContext, v interface{}) error {
	ctx.logger.Debug("decode type def", "v", reflect.TypeOf(v).Elem().Name())
	return t.def.Decode(ctx, v)
}
