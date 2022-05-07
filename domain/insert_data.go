package domain

// InsertData is parse tree of insert command.
type InsertData struct {
	tableName TableName
	fields    []FieldName
	values    []Constant
}

// NewInsertData constructs insert parse tree.
func NewInsertData(tblName TableName, fields []FieldName, vals []Constant) *InsertData {
	return &InsertData{
		tableName: tblName,
		fields:    fields,
		values:    vals,
	}
}
