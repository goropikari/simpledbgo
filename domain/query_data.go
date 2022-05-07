package domain

// QueryData is node of query.
type QueryData struct {
	fields []FieldName
	tables []TableName
	pred   *Predicate
}

// NewQueryData constructs a QueryData.
func NewQueryData(fields []FieldName, tables []TableName, pred *Predicate) *QueryData {
	return &QueryData{
		fields: fields,
		tables: tables,
		pred:   pred,
	}
}
