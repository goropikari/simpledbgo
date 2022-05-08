package domain

// CreateViewData is parse tree of create view command.
type CreateViewData struct {
	viewName  ViewName
	queryData *QueryData
}

// NewCreateViewData constructs create view parse tree.
func NewCreateViewData(name ViewName, qd *QueryData) *CreateViewData {
	return &CreateViewData{
		viewName:  name,
		queryData: qd,
	}
}
