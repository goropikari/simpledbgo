package domain

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
