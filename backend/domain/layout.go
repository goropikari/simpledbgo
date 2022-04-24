package domain

import (
	"log"

	"github.com/goropikari/simpledbgo/meta"
	"github.com/pkg/errors"
)

// Layout is model of layout.
type Layout struct {
	schema   *Schema
	offsets  map[FieldName]int64
	slotsize int64
}

// NewLayout constructs Layout.
func NewLayout(schema *Schema) *Layout {
	pos := int64(meta.Int32Length) // flag for used/unused
	offsets := make(map[FieldName]int64)
	for _, fld := range schema.fields {
		offsets[fld] = pos

		// length in bytes
		switch schema.Type(fld) {
		case Int32:
			pos += meta.Int32Length
		case String:
			pos += meta.Int32Length + int64(schema.Length(fld))
		case Unknown:
			log.Fatal(errors.New("Invalid field type"))
		}
	}

	return &Layout{
		schema:   schema,
		offsets:  offsets,
		slotsize: pos,
	}
}

// NewLayoutWithFields constructs a Layout with fields.
func NewLayoutWithFields(sch *Schema, offsets map[FieldName]int64, slotsize int64) *Layout {
	return &Layout{
		schema:   sch,
		offsets:  offsets,
		slotsize: slotsize,
	}
}

// Schema returns schema.
func (layout *Layout) Schema() *Schema {
	return layout.schema
}

// Offset returns field offset.
func (layout *Layout) Offset(fldname FieldName) int64 {
	return layout.offsets[fldname]
}

// SlotSize returns record slot size.
func (layout *Layout) SlotSize() int64 {
	return layout.slotsize
}
