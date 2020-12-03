package metadata

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Raw metadata for contracts in json
type Raw struct {
	MetadataVersion string `json:"metadataVersion"`
	Source          struct {
		Hash     string `json:"hash"`
		Language string `json:"language"`
		Compiler string `json:"compiler"`
	} `json:"source"`
	Contract struct {
		Name    string   `json:"name"`
		Version string   `json:"version"`
		Authors []string `json:"authors"`
	} `json:"contract"`
	Spec struct {
		Constructors []ConstructorRaw `json:"constructors"`
		Docs         []string         `json:"docs"`
		Events       []EventRaw       `json:"events"`
		Messages     []MessageRaw     `json:"messages"`
	} `json:"spec"`
	Types []rawTypeDef `json:"types"`
}

// TypeIndex type index to def params type
type TypeIndex struct {
	DisplayName []string `json:"displayName"`
	Type        int      `json:"type"`
}

// EventRaw raw json data for event
type EventRaw struct {
	Args []struct {
		Docs    []string  `json:"docs"`
		Indexed bool      `json:"indexed"`
		Name    string    `json:"name"`
		Type    TypeIndex `json:"type"`
	} `json:"args"`
	Docs []string `json:"docs"`
	Name string   `json:"name"`
}

// MessageRaw raw data for message
type MessageRaw struct {
	Args []struct {
		Name string    `json:"name"`
		Type TypeIndex `json:"type"`
	} `json:"args"`
	Docs       []string  `json:"docs"`
	Mutates    bool      `json:"mutates"`
	Name       []string  `json:"name"`
	Payable    bool      `json:"payable"`
	ReturnType TypeIndex `json:"returnType"`
	Selector   string    `json:"selector"`
}

// ConstructorRaw raw data for constructor
type ConstructorRaw struct {
	Args []struct {
		Name string    `json:"name"`
		Type TypeIndex `json:"type"`
	} `json:"args"`
	Docs     []string `json:"docs"`
	Name     []string `json:"name"`
	Selector string   `json:"selector"`
}

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
