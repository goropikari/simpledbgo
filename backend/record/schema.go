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

	// Integer is integer field type.
	Integer

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

// AddIntField adds an int field.
func (schema *Schema) AddIntField(fldname FieldName) {
	schema.AddField(fldname, Integer, 0)
}

// AddStringField adds an string field with maximum length is length.
// length is maximum length of string. It is not actual length of the value.
func (schema *Schema) AddStringField(fldname FieldName, length int) {
	schema.AddField(fldname, String, length)
}

// Add adds other's field into the schema.
func (schema *Schema) Add(fldname FieldName, other *Schema) {
	typ := other.typ(fldname)

	length := other.length(fldname)

	schema.AddField(fldname, typ, length)
}

// AddAllFields adds all fields of other into the schema.
func (schema *Schema) AddAllFields(other *Schema) {
	for _, fld := range other.fields {
		schema.Add(fld, other)
	}
}

func (schema *Schema) typ(fldname FieldName) FieldType {
	if v, found := schema.info[fldname]; found {
		return v.typ
	}

	return Unknown
}

func (schema *Schema) length(fldname FieldName) int {
	if v, found := schema.info[fldname]; found {
		return v.length
	}

	return -1
}
