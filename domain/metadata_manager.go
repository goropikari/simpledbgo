package domain

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// MetadataManager is an interface of MetadataManager.
type MetadataManager interface {
	CreateTable(tblName TableName, sch *Schema, txn Transaction) error
	GetTableLayout(tblName TableName, txn Transaction) (*Layout, error)
	CreateView(viewName ViewName, viewDef ViewDef, txn Transaction) error
	GetViewDef(viewName ViewName, txn Transaction) (ViewDef, error)
	CreateIndex(idxName IndexName, tblName TableName, fldName FieldName, txn Transaction) error
	GetIndexInfo(tblName TableName, txn Transaction) (map[FieldName]*IndexInfo, error)
	GetStatInfo(tblName TableName, layout *Layout, txn Transaction) (StatInfo, error)
}
