package domain

// ModifyData is parse tree of modify data.
type ModifyData struct {
	tblName TableName
	fldName FieldName
	expr    Expression
	pred    *Predicate
}

// NewModifyData constructs a parse tree of modify data.
func NewModifyData(tblName TableName, fldName FieldName, expr Expression, pred *Predicate) *ModifyData {
	return &ModifyData{
		tblName: tblName,
		fldName: fldName,
		expr:    expr,
		pred:    pred,
	}
}
