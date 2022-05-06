package domain

type QueryData struct {
	fields []Field
	tables []TableName
	pred   Predicate
}

func NewQueryData(fields []Field, tables []TableName, pred Predicate) *QueryData {
	return &QueryData{
		fields: fields,
		tables: tables,
		pred:   pred,
	}
}
