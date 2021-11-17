package metadata

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils"
	"github.com/pkg/errors"
)

type compositeField struct {
	Name      string `json:"name"`
	TypeIndex int    `json:"type"`
}

type defComposite struct {
	Fields []compositeField `json:"fields"`
	Path   []string
}

func newDefComposite(raw json.RawMessage, path []string) *defComposite {
	res := &defComposite{}

	if err := json.Unmarshal(raw, res); err != nil {
		panic(err)
	}

	res.Path = path
	return res
}

func (d *defComposite) Encode(ctx CodecContext, value interface{}) error {
	target := reflect.ValueOf(value)
	//fmt.Println("========================================")
	//fmt.Println(value)
	//fmt.Println("=======================================")
	//fmt.Println(target.Kind())
	//fmt.Println("=======================================")

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
				//fmt.Println("--------------------------------------- 2 here is the")
				//fmt.Println("--------------------------------------- 2 here is the")
				//fmt.Println("--------------------------------------- 2 here is the")
				argsType := target.Field(i)
				//fmt.Println("--------------giaogiao")
				argsTypeStruct := TypeIndex{}
				_ = json.Unmarshal([]byte(argsType.String()), &argsTypeStruct)
				if len(argsTypeStruct.DisplayName) == 0 {
					break
				}
				//fmt.Println("--------------------------------------- 2 here is the")
				if err := def.Encode(ctx, argsTypeStruct); err != nil {
					return errors.Wrapf(err,
						"encode composite field %s %d", field.Name, i)
				}
				break
			}
		}
	}

	return nil
}

func (d *defComposite) encodeAccountIDJSON(ctx CodecContext, value json.RawMessage) error {
	var str string
	if err := json.Unmarshal(value, &str); err != nil {
		return errors.Wrap(err, "failed to unmarshal string")
	}

	ctx.logger.Debug("encode accountID", "v", str)

	acc, err := ctx.GetSS58Codec().DecodeAccountID(str)
	if err != nil {
		return errors.Wrap(err, "failed to decode accountID SS58")
	}

	return ctx.encoder.Encode(acc)
}

func (d *defComposite) EncodeJSON(ctx CodecContext, value json.RawMessage) error {
	if utils.IsNameEqual(d.Path, []string{"ink_env", "types", "AccountId"}) {
		return d.encodeAccountIDJSON(ctx, value)
	}

	// for just one sub elem
	if len(d.Fields) == 1 {
		field := d.Fields[0]

		def := ctx.GetDefCodecByIndex(field.TypeIndex)
		ctx.logger.Debug("target", "field", field.Name, "v", value)
		if err := def.EncodeJSON(ctx, value); err != nil {
			return errors.Wrapf(err,
				"encode composite field %s %s", field.Name, value)
		}

		return nil
	}

	jsonMap := make(map[string]json.RawMessage)
	if err := json.Unmarshal(value, &jsonMap); err != nil {
		return errors.Wrap(err, "failed to unmarshal to map")
	}

	for idx, field := range d.Fields {
		ctx.logger.Debug("defComposite encode", "idx", idx, "field", field)

		def := ctx.GetDefCodecByIndex(field.TypeIndex)

		// find field to value
		raw, ok := jsonMap[field.Name]
		if !ok {
			// nil value
			return errors.Errorf("no found field to %s", field.Name)
		}

		ctx.logger.Debug("target", "field", field.Name, "v", raw)
		if err := def.EncodeJSON(ctx, raw); err != nil {
			return errors.Wrapf(err,
				"encode composite field %s %s", field.Name, raw)
		}
	}

	return nil
}

func (d *defComposite) decodeAccountID(ctx CodecContext, value interface{}) error {
	pID, ok := value.(*types.AccountID)
	if !ok {
		return errors.Errorf("decodeAccountID value err")
	}

	return ctx.decoder.Decode(pID)
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

	if utils.IsNameEqual(d.Path, []string{"ink_env", "types", "AccountId"}) {
		return d.decodeAccountID(ctx, value)
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
