package domain

import (
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
	schema.AddField(fldname, FInt32, 0)
}

// AddStringField adds an string field with maximum length is length.
// length is maximum length of string. It is not actual length of the value.
func (schema *Schema) AddStringField(fldname FieldName, length int) {
	schema.AddField(fldname, FString, length)
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

	return FUnknown
}

// Length returns field byte length.
func (schema *Schema) Length(fldname FieldName) int {
	if v, found := schema.info[fldname]; found {
		return v.length
	}

	return -1
}
