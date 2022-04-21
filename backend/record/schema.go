package record

import "github.com/pkg/errors"

const (
	maximumFieldNameLength = 64
)

var (
	// ErrFieldNotFound is an error that means specified field is not found.
	ErrFieldNotFound = errors.New("specified field is not found")

	// ErrExceedMaximumFieldNameLength is an error that means exceeding maximum field name length.
	ErrExceedMaximumFieldNameLength = errors.Errorf("exceeds maximum field name length %v", maximumFieldNameLength)
)

// FieldName is value object of field name.
type FieldName string

// NewFieldName constructs FieldName.
func NewFieldName(name string) (FieldName, error) {
	if len(name) > maximumFieldNameLength {
		return "", ErrExceedMaximumFieldNameLength
	}

	return FieldName(name), nil
}

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
	schema.AddField(fldname, Int32, 0)
}

// AddStringField adds an string field with maximum length is length.
// length is maximum length of string. It is not actual length of the value.
func (schema *Schema) AddStringField(fldname FieldName, length int) {
	schema.AddField(fldname, String, length)
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
func (sch *Schema) Fields() []FieldName {
	return sch.fields
}

// Type returns field type.
func (schema *Schema) Type(fldname FieldName) FieldType {
	if v, found := schema.info[fldname]; found {
		return v.typ
	}

	return Unknown
}

// Length returns field byte length.
func (schema *Schema) Length(fldname FieldName) int {
	if v, found := schema.info[fldname]; found {
		return v.length
	}

	return -1
}
