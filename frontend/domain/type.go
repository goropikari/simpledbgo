package domain

type Field string

func NewField(fld string) Field {
	return Field(fld)
}

type TableName string

func NewTableName(name string) TableName {
	return TableName(name)
}
