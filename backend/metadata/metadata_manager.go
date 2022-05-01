package metadata

import "github.com/goropikari/simpledbgo/backend/domain"

// Manager manages metadata.
type Manager struct {
	tblMgr  *TableManager
	viewMgr *ViewManager
	statMgr *StatManager
	idxMgr  *IndexManager
}

// CreateManager creates metadata manager with initializing tables related to metadata.
func CreateManager(txn domain.Transaction) (*Manager, error) {
	tblMgr, err := CreateTableManager(txn)
	if err != nil {
		return nil, err
	}

	viewMgr, err := CreateViewManager(tblMgr, txn)
	if err != nil {
		return nil, err
	}

	statMgr, err := NewStatManager(tblMgr, txn)
	if err != nil {
		return nil, err
	}

	idxMgr, err := CreateIndexManager(tblMgr, statMgr, txn)
	if err != nil {
		return nil, err
	}

	return &Manager{
		tblMgr:  tblMgr,
		viewMgr: viewMgr,
		statMgr: statMgr,
		idxMgr:  idxMgr,
	}, nil
}

// NewManager constructs metadata manager.
func NewManager(txn domain.Transaction) (*Manager, error) {
	tblMgr := NewTableManager()
	viewMgr := NewViewManager(tblMgr)
	statMgr, err := NewStatManager(tblMgr, txn)
	if err != nil {
		return nil, err
	}

	idxMgr, err := NewIndexManager(tblMgr, statMgr, txn)
	if err != nil {
		return nil, err
	}

	return &Manager{
		tblMgr:  tblMgr,
		viewMgr: viewMgr,
		statMgr: statMgr,
		idxMgr:  idxMgr,
	}, nil
}

// CreateTable creates a table.
func (mgr *Manager) CreateTable(tblName domain.TableName, sch *domain.Schema, txn domain.Transaction) error {
	return mgr.tblMgr.CreateTable(tblName, sch, txn)
}

// GetTableLayout returns given table layout.
func (mgr *Manager) GetTableLayout(tblName domain.TableName, txn domain.Transaction) (*domain.Layout, error) {
	return mgr.tblMgr.GetTableLayout(tblName, txn)
}

// CreateView creates a view.
func (mgr *Manager) CreateView(viewName ViewName, viewDef ViewDef, txn domain.Transaction) error {
	return mgr.viewMgr.CreateView(viewName, viewDef, txn)
}

// GetViewDef returns given view definition.
func (mgr *Manager) GetViewDef(viewName ViewName, txn domain.Transaction) (ViewDef, error) {
	return mgr.viewMgr.GetViewDef(viewName, txn)
}

// CreateIndex creates an index.
func (mgr *Manager) CreateIndex(idxName domain.IndexName, tblName domain.TableName, fldName domain.FieldName, txn domain.Transaction) error {
	return mgr.idxMgr.CreateIndex(idxName, tblName, fldName, txn)
}

// GetIndexInfo returns the index information of given table.
func (mgr *Manager) GetIndexInfo(tblName domain.TableName, txn domain.Transaction) (map[domain.FieldName]*IndexInfo, error) {
	return mgr.idxMgr.GetIndexInfo(tblName, txn)
}

// GetStatInfo returns the statistical information of given table.
func (mgr *Manager) GetStatInfo(tblName domain.TableName, layout *domain.Layout, txn domain.Transaction) (StatInfo, error) {
	return mgr.statMgr.GetStatInfo(tblName, layout, txn)
}
