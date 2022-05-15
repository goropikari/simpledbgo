package domain

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// UpdateExecutor is an interface of UpdateExecutor.
type UpdateExecutor interface {
	ExecuteInsert(*InsertData, Transaction) (int, error)
	ExecuteDelete(*DeleteData, Transaction) (int, error)
	ExecuteModify(*ModifyData, Transaction) (int, error)
	ExecuteCreateTable(*CreateTableData, Transaction) (int, error)
	ExecuteCreateView(*CreateViewData, Transaction) (int, error)
	ExecuteCreateIndex(*CreateIndexData, Transaction) (int, error)
}
