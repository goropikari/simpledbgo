package domain

import (
	"fmt"
	"strings"
)

// ExecData is parse tree of execute command.
type ExecData interface{}

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

// Fields returns field list.
func (data *QueryData) Fields() []FieldName {
	return data.fields
}

// Tables returns table list.
func (data *QueryData) Tables() []TableName {
	return data.tables
}

// Predicate returns predicate.
func (data *QueryData) Predicate() *Predicate {
	return data.pred
}

// String stringfies data.
func (data *QueryData) String() string {
	fields := make([]string, 0, len(data.fields))
	for _, f := range data.fields {
		fields = append(fields, f.String())
	}

	tables := make([]string, 0, len(data.tables))
	for _, t := range data.tables {
		tables = append(tables, t.String())
	}

	query := fmt.Sprintf("select %v from %v ", strings.Join(fields, ","), strings.Join(tables, ","))

	if pred := data.pred.String(); pred != "" {
		return query + " where " + pred
	}

	return query
}

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

// TableName returns table name.
func (data *InsertData) TableName() TableName {
	return data.tableName
}

// Fields returns field list.
func (data *InsertData) Fields() []FieldName {
	return data.fields
}

// Values returns insertion values.
func (data *InsertData) Values() []Constant {
	return data.values
}

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

// TableName returns table name.
func (data *DeleteData) TableName() TableName {
	return data.tableName
}

// Predicate returns predicate.
func (data *DeleteData) Predicate() *Predicate {
	return data.pred
}

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

// TableName returns table name.
func (data *ModifyData) TableName() TableName {
	return data.tblName
}

// FieldName returns field name.
func (data *ModifyData) FieldName() FieldName {
	return data.fldName
}

// Expression returns an expression.
func (data *ModifyData) Expression() Expression {
	return data.expr
}

// Predicate returns a predicate.
func (data *ModifyData) Predicate() *Predicate {
	return data.pred
}

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

// TableName returns a table name.
func (data *CreateTableData) TableName() TableName {
	return data.tblName
}

// Schema returns a schema.
func (data *CreateTableData) Schema() *Schema {
	return data.sch
}

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

// ViewName returns a view name.
func (data *CreateViewData) ViewName() ViewName {
	return data.viewName
}

// ViewDef returns the definition of view.
func (data *CreateViewData) ViewDef() ViewDef {
	return NewViewDef(data.queryData.String())
}

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

// IndexName returns index name.
func (data *CreateIndexData) IndexName() IndexName {
	return data.idxName
}

// TableName returns a table name.
func (data *CreateIndexData) TableName() TableName {
	return data.tblName
}

// FieldName returns a field name.
func (data *CreateIndexData) FieldName() FieldName {
	return data.fldName
}
