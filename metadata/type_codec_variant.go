package metadata

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

type VariantsI interface {
	ToFieldValue() (string, interface{})
}

type VariantsFromJSON struct {
	Type string          `json:"typ"`
	Val  json.RawMessage `json:"val"`
}

type variantType struct {
	Name   string `json:"name"`
	Fields []struct {
		Typ int `json:"type"`
	}
	Discriminant int `json:"discriminant"`
}

type defVariant struct {
	Variants      []variantType `json:"variants"`
	noneFieldsIdx int           // none fields index + 1
}

func newDefVariant(raw json.RawMessage) *defVariant {
	res := &defVariant{
		noneFieldsIdx: -1,
	}

	if err := json.Unmarshal(raw, &res); err != nil {
		panic(err)
	}

	for idx, v := range res.Variants {
		if len(v.Fields) == 0 {
			res.noneFieldsIdx = idx
		}
	}

	return res
}

func (d *defVariant) findVarIdxByName(name string) int {
	for idx, v := range d.Variants {
		if v.Name == name {
			return idx
		}
	}
	return -1
}

func (d *defVariant) Encode(ctx CodecContext, value interface{}) error {
	if o, ok := value.(VariantsI); ok {
		return d.encodeVariant(ctx, o)
	}

	return d.encodeCommonStruct(ctx, value)
}

func (d *defVariant) EncodeJSON(ctx CodecContext, value json.RawMessage) error {
	return nil
}

func (d *defVariant) Decode(ctx CodecContext, value interface{}) error {
	return d.decodeCommonStruct(ctx, value)
}

func (d *defVariant) encodeFields(ctx CodecContext, index int, value interface{}) error {
	ctx.encoder.Encode(byte(index))
	if value != nil {
		if index >= len(d.Variants) {
			return errors.Errorf("err index for %d, all %d", index, len(d.Variants))
		}

		v := reflect.ValueOf(value).Elem().Interface()

		fields := d.Variants[index]
		if len(fields.Fields) == 0 {
			return errors.Errorf("err not nil val to zero fields %d", index)
		}

		for idx, f := range fields.Fields {
			codec := ctx.GetDefCodecByIndex(f.Typ)
			if err := codec.Encode(ctx, v); err != nil {
				return errors.Wrapf(err, "encode fields %d with i: %d, t: %d",
					index, idx, f.Typ)
			}

			// now just one
			return nil
		}
	}
	return nil
}

/*
encodeCommonStruct A Variants common type


type VariantsType struct {
	V1 *Type1 `scale:"type1"`
	V2 *Type2 `scale:"type2"`
}

if V1 is nil, V2 is nil and variants has a none opt which fields is empty, then is a null type,
if V1 is not nil, then variants will be v1 with its idx in header
if V1 is nil, V2 is not nil and variants will be v2 with its idx in header
*/
func (d *defVariant) encodeCommonStruct(ctx CodecContext, value interface{}) error {
	t := reflect.ValueOf(value)
	if t.Kind() == reflect.Ptr && t.IsNil() && d.noneFieldsIdx >= 0 {
		return errors.Wrapf(d.encodeFields(ctx, d.noneFieldsIdx, nil),
			"encode none %d fields err", d.noneFieldsIdx)
	}

	numFields := t.NumField()
	for i := 0; i < numFields; i++ {
		vi := t.Field(i)
		vt := t.Type().Field(i)

		tv, ok := vt.Tag.Lookup("scale")
		if ok && tv == "-" {
			continue
		}

		index := d.findVarIdxByName(tv)
		if index < 0 {
			continue
		}

		if !vi.IsNil() {
			return errors.Wrapf(d.encodeFields(ctx, index, vi.Interface()),
				"encode %d fields err", index)
		}
	}

	// none
	if d.noneFieldsIdx >= 0 {
		return errors.Wrapf(d.encodeFields(ctx, d.noneFieldsIdx, nil),
			"encode %d fields err", d.noneFieldsIdx)
	}

	return errors.Errorf("no fields has value and variants not allow none")
}

func (d *defVariant) decodeCommonStruct(ctx CodecContext, value interface{}) error {
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

	var typ byte
	if err := ctx.decoder.Decode(&typ); err != nil {
		return errors.Wrapf(err, "decode type error")
	}

	ctx.logger.Debug("target", "typ", typ)

	typIdx := int(typ)

	if typIdx >= len(d.Variants) {
		return errors.Errorf("typeIdx no exited %d", typIdx)
	}

	variant := d.Variants[typIdx]
	name := d.Variants[typIdx].Name

	// find field to value
	for i := 0; i < target.NumField(); i++ {
		ft := target.Type().Field(i)
		tv, ok := ft.Tag.Lookup("scale")
		if ok && tv == "-" {
			continue
		}

		if tv == name {
			ctx.logger.Debug("target", "field", tv, "v", target.Field(i))

			for _, f := range variant.Fields {
				def := ctx.GetDefCodecByIndex(f.Typ)
				fi := target.Field(i).Interface()

				if err := def.Decode(ctx, fi); err != nil {
					return errors.Wrapf(err, "decode composite field %d %d", f.Typ, i)
				}
				break
			}
		}
	}

	return nil
}

func (d *defVariant) encodeVariant(ctx CodecContext, value VariantsI) error {
	name, field := value.ToFieldValue()
	index := d.findVarIdxByName(name)
	if index < 0 {
		return errors.Errorf("no find variant field name %s", name)
	}
	return errors.Wrapf(d.encodeFields(ctx, index, field),
		"encode a variant interface %d", index)
}
