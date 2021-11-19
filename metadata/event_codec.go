package metadata

import (
	"reflect"

	"github.com/pkg/errors"
)

type EventRaws []EventRaw

func GetEvtTypeIdx(data []byte) int {
	return int(data[0])
}

func (e EventRaws) DecodeEvt(ctx CodecContext, target interface{}) error {
	b, err := ctx.decoder.ReadOneByte()
	if err != nil {
		return errors.Wrap(err, "decode bytes")
	}

	bi := int(b)
	if bi >= len(e) {
		return errors.Errorf("idx no found for %d but %d", bi, len(e))
	}

	ctx.logger.Debug("type event", "b", bi)

	evtRaw := e[bi]
	for idx, f := range evtRaw.Args {
		vt := reflect.ValueOf(target).Elem()
		for i := 0; i < vt.NumField(); i++ {
			vtt := vt.Type().Field(i)

			tv, ok := vtt.Tag.Lookup("scale")
			if ok && tv == "-" {
				continue
			}

			if tv == f.Name {
				ctx.logger.Debug("decode evt", "name", f.Name, "i", i, "typ", f.Type.TypeIndex)

				td := ctx.GetDefCodecByIndex(f.Type.TypeIndex)
				di := vt.Field(i).Addr().Interface()
				if err := td.Decode(ctx, di); err != nil {
					return errors.Wrapf(err,
						"decode arg %d field %d with typ %d err", idx, i, f.Type.TypeIndex)
				}

				ctx.logger.Debug("decode", "v", di)
			}
		}
	}

	return nil
}
