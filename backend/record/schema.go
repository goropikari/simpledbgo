package record

import (
	"github.com/goropikari/simpledbgo/backend/domain"
	"github.com/pkg/errors"
)

var (
	// ErrFieldNotFound is an error that means specified field is not found.
	ErrFieldNotFound = errors.New("specified field is not found")
)

// FieldType is a type of field.
type FieldType uint

const (
	// Unknown is unknown field type.
	Unknown FieldType = iota

	// Int32 is integer field type.
	Int32

	// String is string field type.
	String
)

// FieldInfo is a model of field information.
type FieldInfo struct {
	typ    FieldType
	length int
}

// Schema is model of table schema.
type Schema struct {
	fields []domain.FieldName
	info   map[domain.FieldName]*FieldInfo
}

// NewSchema constructs a Schema.
func NewSchema() *Schema {
	return &Schema{
		fields: make([]domain.FieldName, 0),
		info:   make(map[domain.FieldName]*FieldInfo),
	}
}

// HasField checks existence of fldname.
func (schema *Schema) HasField(fldname domain.FieldName) bool {
	_, found := schema.info[fldname]

	return found
}

// AddField adds a field in to the schema.
func (schema *Schema) AddField(fldname domain.FieldName, typ FieldType, length int) {
	schema.fields = append(schema.fields, fldname)
	schema.info[fldname] = &FieldInfo{
		typ:    typ,
		length: length,
	}
}

// AddInt32Field adds an int field.
func (schema *Schema) AddInt32Field(fldname domain.FieldName) {
	schema.AddField(fldname, Int32, 0)
}

// AddStringField adds an string field with maximum length is length.
// length is maximum length of string. It is not actual length of the value.
func (schema *Schema) AddStringField(fldname domain.FieldName, length int) {
	schema.AddField(fldname, String, length)
}

// Add adds other's field into the schema.
func (schema *Schema) Add(fldname domain.FieldName, other *Schema) {
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
func (sch *Schema) Fields() []domain.FieldName {
	return sch.fields
}

// Type returns field type.
func (schema *Schema) Type(fldname domain.FieldName) FieldType {
	if v, found := schema.info[fldname]; found {
		return v.typ
	}

	return Unknown
}

// Length returns field byte length.
func (schema *Schema) Length(fldname domain.FieldName) int {
	if v, found := schema.info[fldname]; found {
		return v.length
	}

	return -1
}
