package metadata

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

const (
	OptionNoneName = "None"
	OptionSomeName = "Some"
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

func (d *defVariant) IsOptional(v interface{}) (OptionValue, int, bool) {
	if len(d.Variants) == 2 &&
		d.Variants[0].Name == OptionNoneName &&
		len(d.Variants[0].Fields) == 0 &&
		d.Variants[1].Name == OptionSomeName &&
		len(d.Variants[1].Fields) == 1 {
		ov, ok := v.(OptionValue)
		if ok {
			return ov, d.Variants[1].Fields[0].Typ, true
		}
	}

	return nil, 0, false
}

func (d *defVariant) Encode(ctx CodecContext, value interface{}) error {
	if o, ok := value.(VariantsI); ok {
		return d.encodeVariant(ctx, o)
	}

	if ov, typIdx, ok := d.IsOptional(value); ok {
		return d.encodeOption(ctx, typIdx, ov)
	}

	return d.encodeCommonStruct(ctx, value)
}

func (d *defVariant) EncodeJSON(ctx CodecContext, value json.RawMessage) error {
	return nil
}

func (d *defVariant) Decode(ctx CodecContext, value interface{}) error {
	if v, typIdx, ok := d.IsOptional(value); ok {
		return d.decodeOption(ctx, typIdx, v)
	}

	return d.decodeCommonStruct(ctx, value)
}

func (d *defVariant) encodeFields(ctx CodecContext, index int, value interface{}) error {
	err := ctx.encoder.Encode(byte(index))
	if err != nil {
		return errors.Wrap(err, "ctx encoder Encode err")
	}

	if value != nil {
		if index >= len(d.Variants) {
			return errors.Errorf("err index for %d, all %d", index, len(d.Variants))
		}

		v := reflect.ValueOf(value).Elem().Interface()

		fields := d.Variants[index]
		if len(fields.Fields) == 0 {
			return errors.Errorf("err not nil val to zero fields %d", index)
		}

		// now just one
		field := fields.Fields[0]

		codec := ctx.GetDefCodecByIndex(field.Typ)
		if err := codec.Encode(ctx, v); err != nil {
			return errors.Wrapf(err, "encode fields %d with i: %d, t: %d",
				index, 0, field.Typ)
		}
	}

	return nil
}

const (
	OptionValueNilPrefix  byte = 0x00
	OptionValueSomePrefix byte = 0x01
)

type OptionValue interface {
	IsNone() bool
	SetHasValue(h bool)
}

func (d *defVariant) encodeOption(ctx CodecContext, _ int, value OptionValue) error {
	if value.IsNone() {
		return ctx.encoder.PushByte(OptionValueNilPrefix)
	}

	err := ctx.encoder.PushByte(OptionValueSomePrefix)
	if err != nil {
		return errors.Wrap(err, "push byte 0x01 err")
	}

	t := reflect.ValueOf(value)
	numFields := t.NumField()
	for i := 0; i < numFields; i++ {
		vi := t.Field(i)
		vt := t.Type().Field(i)

		tv, ok := vt.Tag.Lookup("scale")
		if ok && tv == "-" {
			continue
		}

		if tv == OptionSomeName {
			return d.encodeFields(ctx, 1, vi.Interface())
		}
	}

	return errors.Errorf("no found field some")
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

func (d *defVariant) decodeOption(ctx CodecContext, typIdx int, value OptionValue) error {
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

	typ, err := ctx.decoder.ReadOneByte()
	if err != nil {
		return errors.Wrapf(err, "decode type error")
	}

	if typ == OptionValueNilPrefix {
		value.SetHasValue(false)
		return nil
	}

	value.SetHasValue(true)

	// find field to value
	for i := 0; i < target.NumField(); i++ {
		ft := target.Type().Field(i)
		tv, ok := ft.Tag.Lookup("scale")
		if ok && tv == "-" {
			continue
		}

		if tv == OptionSomeName {
			ctx.logger.Debug("target option", "field", tv, "typIdx", typIdx, "v", target.Field(i))

			def := ctx.GetDefCodecByIndex(typIdx)
			fi := target.Field(i).Addr().Interface()

			if err := def.Decode(ctx, fi); err != nil {
				return errors.Wrapf(err, "decode composite field %d %d", typIdx, i)
			}

			return nil
		}
	}

	return errors.Errorf("no Some field find for optional")
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

	typ, err := ctx.decoder.ReadOneByte()
	if err != nil {
		return errors.Wrapf(err, "decode type error")
	}

	ctx.logger.Debug("decodeCommonStruct target", "typ", typ, "target", target)

	typIdx := int(typ)

	if typIdx >= len(d.Variants) {
		return errors.Errorf("typeIdx no exited %d", typIdx)
	}

	variant := d.Variants[typIdx]
	name := d.Variants[typIdx].Name

	if len(variant.Fields) == 0 {
		ctx.logger.Debug("decodeCommonStruct nil target", "typ", typ, "target", target)

		return nil
	}

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
