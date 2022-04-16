package record

import (
	"log"

	"github.com/goropikari/simpledbgo/meta"
	"github.com/pkg/errors"
)

// Layout is model of layout.
type Layout struct {
	schema   *Schema
	offsets  map[FieldName]int64
	slotsize int
}

// NewLayout constructs Layout.
func NewLayout(schema *Schema) *Layout {
	pos := int64(meta.Int32Length) // flag for used/unused
	offsets := make(map[FieldName]int64)
	for _, fld := range schema.fields {
		offsets[fld] = pos

		switch schema.typ(fld) {
		case Integer:
			pos += meta.Int32Length
		case String:
			pos += meta.Int32Length + int64(schema.length(fld))
		case Unknown:
			log.Fatal(errors.New("Invalid field type"))
		}
	}

	return &Layout{
		schema:   schema,
		offsets:  offsets,
		slotsize: int(pos),
	}
}
