package domain

// CreateTableData is parse tree of create table command.
type CreateTableData struct {
	tblName TableName
	sch     *Schema
}

// NewCreateTableData constructs a CreateTableData.
func NewCreateTableData(tblName TableName, sch *Schema) *CreateTableData {
	return &CreateTableData{
		tblName: tblName,
		sch:     sch,
	}
}
