package domain

import "github.com/pkg/errors"

const (
	// MaxFieldNameLength is maximum field name length.
	MaxFieldNameLength = 16

	// MaxTableNameLength is maximum table name length.
	MaxTableNameLength = 16

	// MaxViewNameLength is maximum view name length.
	MaxViewNameLength = 16

	// MaxIndexNameLength is maximum index name length.
	MaxIndexNameLength = 16

	// MaxViewDefLength is maximum view definition length.
	MaxViewDefLength = 100
)

var (
	// ErrExceedMaxFieldNameLength is an error that means exceeding maximum field name length.
	ErrExceedMaxFieldNameLength = errors.Errorf("exceeds maximum field name length %v", MaxFieldNameLength)

	// ErrExceedMaxViewNameLength is an error that means exceeding maximum view name length.
	ErrExceedMaxViewNameLength = errors.Errorf("exceeds maximum view name length %v", MaxViewNameLength)
)

// LSN is log sequence number.
type LSN int32

// SlotID is identifier of slot.
type SlotID int32

// NewSlotID constructs a slot id.
func NewSlotID(id int32) SlotID {
	return SlotID(id)
}

// FieldName is value object of field name.
type FieldName string

// NewFieldName constructs FieldName.
func NewFieldName(name string) (FieldName, error) {
	if len(name) > MaxFieldNameLength {
		return "", ErrExceedMaxFieldNameLength
	}
	if name == "" {
		return "", errors.New("field name must not be empty")
	}

	return FieldName(name), nil
}

// String stringfy name.
func (name FieldName) String() string {
	return string(name)
}

// IsZero checks whether name is zero value or not.
func (name FieldName) IsZero() bool {
	return name == ""
}

// FieldType is a type of field.
type FieldType uint

const (
	// FUnknown is unknown field type.
	FUnknown FieldType = iota

	// FInt32 is integer field type.
	FInt32

	// FString is string field type.
	FString
)

// FieldInfo is a model of field information.
type FieldInfo struct {
	typ    FieldType
	length int
}

// TableName is value object of table name.
type TableName string

// NewTableName constructs TableName.
func NewTableName(name string) (TableName, error) {
	if len(name) > MaxTableNameLength {
		return "", ErrExceedMaxFieldNameLength
	}
	if name == "" {
		return "", errors.New("table name must not be empty")
	}

	return TableName(name), nil
}

// String stringfy name.
func (name TableName) String() string {
	return string(name)
}

// ToFileName converts type from TableName into FileName.
func (name TableName) ToFileName() FileName {
	return FileName(name)
}

// ToViewName converts type from ViewName into ViewName.
func (name TableName) ToViewName() ViewName {
	return ViewName(name)
}

// ViewName is type of view name.
type ViewName string

// NewViewName constructs a ViewName.
func NewViewName(name string) (ViewName, error) {
	if len(name) > MaxViewNameLength {
		return "", ErrExceedMaxViewNameLength
	}

	return ViewName(name), nil
}

// String stringfies ViewName.
func (name ViewName) String() string {
	return string(name)
}

// ViewDef is type of view definition.
type ViewDef string

// NewViewDef constructs a ViewDef.
func NewViewDef(def string) ViewDef {
	return ViewDef(def)
}

// String stringfies ViewDef.
func (def ViewDef) String() string {
	return string(def)
}

// IsZero checks def is zero value or not.
func (def ViewDef) IsZero() bool {
	return def == ""
}
