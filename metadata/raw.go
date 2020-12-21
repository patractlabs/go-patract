package metadata

import (
	"github.com/patractlabs/go-patract/utils"
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
		Events       EventRaws        `json:"events"`
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

type ArgRaw struct {
	Name string    `json:"name"`
	Type TypeIndex `json:"type"`
}

// MessageRaw raw data for message
type MessageRaw struct {
	Args         []ArgRaw  `json:"args"`
	Docs         []string  `json:"docs"`
	Mutates      bool      `json:"mutates"`
	Name         []string  `json:"name"`
	Payable      bool      `json:"payable"`
	ReturnType   TypeIndex `json:"returnType"`
	Selector     string    `json:"selector"`
	SelectorData []byte    `json:"-"`
}

// ConstructorRaw raw data for constructor
type ConstructorRaw struct {
	Args         []ArgRaw `json:"args"`
	Docs         []string `json:"docs"`
	Name         []string `json:"name"`
	Selector     string   `json:"selector"`
	SelectorData []byte   `json:"-"`
}

// GetConstructor get constructor by name
func (r *Raw) GetConstructor(name []string) (ConstructorRaw, error) {
	for _, c := range r.Spec.Constructors {
		if utils.IsNameEqual(c.Name, name) {
			return c, nil
		}
	}

	return ConstructorRaw{}, errors.Errorf("no found constructor for %s", name)
}

// GetMessage get message by name
func (r *Raw) GetMessage(name []string) (MessageRaw, error) {
	for _, c := range r.Spec.Messages {
		if utils.IsNameEqual(c.Name, name) {
			return c, nil
		}
	}

	return MessageRaw{}, errors.Errorf("no found constructor for %s", name)
}
