package metadata

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// rawTypeDef type definition
type rawTypeDef struct {
	Def    map[string]json.RawMessage `json:"def"`
	Params []int                      `json:"params"`
	Path   []string                   `json:"path"`
}

// TypeDef type definition
type TypeDef struct {
	def    DefCodec
	Params []int
	Path   []string
}

func NewTypeDef(raw *rawTypeDef) *TypeDef {
	if len(raw.Def) != 1 {
		panic(errors.Errorf("type def raw error by not key %v", raw.Def))
	}

	res := &TypeDef{
		Params: raw.Params,
		Path:   raw.Path,
	}

	for k, jsonRaw := range raw.Def {
		switch k {
		case "primitive":
			res.def = newDefPrimitive(jsonRaw)
		case "composite":
			res.def = newDefComposite(jsonRaw)
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
	return t.def.Encode(ctx, v)
}

func (t *TypeDef) Decode(ctx CodecContext, v interface{}) error {
	return t.def.Decode(ctx, v)
}
