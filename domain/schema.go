package domain

import (
	"log"

	"github.com/goropikari/simpledbgo/common"
	"github.com/pkg/errors"
)

// ErrFieldNotFound is an error that means specified field is not found.
var ErrFieldNotFound = errors.New("specified field is not found")

// Schema is model of table schema.
type Schema struct {
	fields []FieldName
	info   map[FieldName]*FieldInfo
}

// NewSchema constructs a Schema.
func NewSchema() *Schema {
	return &Schema{
		fields: make([]FieldName, 0),
		info:   make(map[FieldName]*FieldInfo),
	}
}

// HasField checks existence of fldname.
func (schema *Schema) HasField(fldname FieldName) bool {
	_, found := schema.info[fldname]

	return found
}

// AddField adds a field in to the schema.
func (schema *Schema) AddField(fldname FieldName, typ FieldType, length int) {
	schema.fields = append(schema.fields, fldname)
	schema.info[fldname] = &FieldInfo{
		typ:    typ,
		length: length,
	}
}

// AddInt32Field adds an int field.
func (schema *Schema) AddInt32Field(fldname FieldName) {
	schema.AddField(fldname, Int32FieldType, 0)
}

// AddStringField adds an string field with maximum length is length.
// length is maximum length of string. It is not actual length of the value.
func (schema *Schema) AddStringField(fldname FieldName, length int) {
	schema.AddField(fldname, StringFieldType, length)
}

// Add adds other's field into the schema.
func (schema *Schema) Add(fldname FieldName, other *Schema) {
	typ := other.Type(fldname)

	length := other.Length(fldname)

	schema.AddField(fldname, typ, length)
}

// AddAllFields adds all fields of other into the schema.
func (schema *Schema) AddAllFields(other *Schema) {
	for _, fld := range other.fields {
		schema.Add(fld, other)
	}
}

// Fields returns schema fileds.
func (schema *Schema) Fields() []FieldName {
	return schema.fields
}

// Type returns field type.
func (schema *Schema) Type(fldname FieldName) FieldType {
	if v, found := schema.info[fldname]; found {
		return v.typ
	}

	return UnknownFieldType
}

// Length returns field byte length.
func (schema *Schema) Length(fldname FieldName) int {
	if v, found := schema.info[fldname]; found {
		return v.length
	}

	return -1
}

// Layout is model of table layout.
type Layout struct {
	schema   *Schema
	offsets  map[FieldName]int64
	slotsize int64
}

// NewLayout constructs Layout.
func NewLayout(schema *Schema) *Layout {
	pos := int64(common.Int32Length) // flag for used/unused
	offsets := make(map[FieldName]int64)
	for _, fld := range schema.fields {
		offsets[fld] = pos

		// length in bytes
		switch schema.Type(fld) {
		case Int32FieldType:
			pos += common.Int32Length
		case StringFieldType:
			pos += common.Int32Length + int64(schema.Length(fld))
		case UnknownFieldType:
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

// Length returns byte size of given field name.
func (layout *Layout) Length(fldName FieldName) int {
	return layout.schema.Length(fldName)
}
