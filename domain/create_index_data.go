package domain

// CreateIndexData is parse tree of create index command.
type CreateIndexData struct {
	idxName IndexName
	tblName TableName
	fldName FieldName
}

// NewCreateIndexData constructs a CreateIndexData.
func NewCreateIndexData(idxName IndexName, tblName TableName, fldName FieldName) *CreateIndexData {
	return &CreateIndexData{
		idxName: idxName,
		tblName: tblName,
		fldName: fldName,
	}
}
