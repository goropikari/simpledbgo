package domain

// Field is type for field.
type Field string

// NewField constructs a Field.
func NewField(fld string) Field {
	return Field(fld)
}

// TableName is type of table name.
type TableName string

// NewTableName constructs a TableName.
func NewTableName(name string) TableName {
	return TableName(name)
}
