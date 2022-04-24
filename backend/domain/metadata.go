package domain

import "github.com/pkg/errors"

const (
	MaximumFieldNameLength = 64
)

var (
	// ErrExceedMaximumFieldNameLength is an error that means exceeding maximum field name length.
	ErrExceedMaximumFieldNameLength = errors.Errorf("exceeds maximum field name length %v", MaximumFieldNameLength)
)

// SlotID is identifier of slot.
type SlotID int32

// FieldName is value object of field name.
type FieldName string

// NewFieldName constructs FieldName.
func NewFieldName(name string) (FieldName, error) {
	if len(name) > MaximumFieldNameLength {
		return "", ErrExceedMaximumFieldNameLength
	}
	if name == "" {
		return "", errors.New("field name must not be empty")
	}

	return FieldName(name), nil
}
