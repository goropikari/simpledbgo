package domain

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// Planner is an interface of planner.
type Planner interface {
	Open() (Scanner, error)
	// EstNumBlocks estimates the number of block access.
	EstNumBlocks() int
	// EstNumRecord estimates the number of record access.
	EstNumRecord() int
	// EstDistinctVals estimates the number of distinct value at given fldName.
	EstDistinctVals(FieldName) int
	Schema() *Schema
}

// QueryPlanner is an interface of query planner.
type QueryPlanner interface {
	CreatePlan(*QueryData, Transaction) (Planner, error)
}

// UpdateExecutor is an interface of UpdateExecutor.
type UpdateExecutor interface {
	ExecuteInsert(*InsertData, Transaction) (int, error)
	ExecuteDelete(*DeleteData, Transaction) (int, error)
	ExecuteModify(*ModifyData, Transaction) (int, error)
	ExecuteCreateTable(*CreateTableData, Transaction) (int, error)
	ExecuteCreateView(*CreateViewData, Transaction) (int, error)
	ExecuteCreateIndex(*CreateIndexData, Transaction) (int, error)
}
