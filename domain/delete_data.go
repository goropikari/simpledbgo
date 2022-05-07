package domain

// DeleteData is a parse tree of delete command.
type DeleteData struct {
	tableName TableName
	pred      *Predicate
}

// NewDeleteData constructs a DeleteData.
func NewDeleteData(name TableName, pred *Predicate) *DeleteData {
	return &DeleteData{
		tableName: name,
		pred:      pred,
	}
}
