package domain

import "github.com/pkg/errors"

const (
	// MaxFieldNameLength is maximum field name length.
	MaxFieldNameLength = 16

	// MaxTableNameLength is maximum table name length.
	MaxTableNameLength = 16

	// MaxViewDefLength is maximum view definition length.
	MaxViewDefLength = 100
)

// ErrExceedMaxFieldNameLength is an error that means exceeding maximum field name length.
var ErrExceedMaxFieldNameLength = errors.Errorf("exceeds maximum field name length %v", MaxFieldNameLength)

// SlotID is identifier of slot.
type SlotID int32

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
