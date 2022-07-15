package metadata

import (
	"fmt"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/tx"
)

// Manager manages metadata.
type Manager struct {
	tblMgr  *TableManager
	viewMgr *ViewManager
	statMgr *StatManager
	idxMgr  *IndexManager
}

// NewManager constructs a metadata manager.
func NewManager(driver domain.IndexDriver, fileMgr domain.FileManager, logMgr domain.LogManager, bufMgr domain.BufferPoolManager, concurMgr domain.ConcurrencyManager, gen domain.TxNumberGenerator) (*Manager, error) {
	if fileMgr.IsInit() {
		return createManager(driver, fileMgr, logMgr, bufMgr, concurMgr, gen)
	}

	return newManager(driver, fileMgr, logMgr, bufMgr, concurMgr, gen)
}

// createManager creates metadata manager with initializing tables related to metadata.
func createManager(driver domain.IndexDriver, fileMgr domain.FileManager, logMgr domain.LogManager, bufMgr domain.BufferPoolManager, concurMgr domain.ConcurrencyManager, gen domain.TxNumberGenerator) (*Manager, error) {
	txn, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
	if err != nil {
		return nil, err
	}

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

	idxMgr, err := CreateIndexManager(driver, tblMgr, statMgr, txn)
	if err != nil {
		return nil, err
	}

	err = txn.Commit()
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

// newManager constructs metadata manager.
func newManager(driver domain.IndexDriver, fileMgr domain.FileManager, logMgr domain.LogManager, bufMgr domain.BufferPoolManager, concurMgr domain.ConcurrencyManager, gen domain.TxNumberGenerator) (*Manager, error) {
	txn, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
	if err != nil {
		return nil, err
	}

	tblMgr := NewTableManager()
	viewMgr := NewViewManager(tblMgr)
	statMgr, err := NewStatManager(tblMgr, txn)
	if err != nil {
		return nil, err
	}

	idxMgr, err := NewIndexManager(driver, tblMgr, statMgr, txn)
	if err != nil {
		return nil, err
	}

	err = txn.Commit()
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
	if mgr.tblMgr.Exists(tblName, txn) {
		return fmt.Errorf("table %v already exists", tblName)
	}

	return mgr.tblMgr.CreateTable(tblName, sch, txn)
}

// GetTableLayout returns given table layout.
func (mgr *Manager) GetTableLayout(tblName domain.TableName, txn domain.Transaction) (*domain.Layout, error) {
	return mgr.tblMgr.GetTableLayout(tblName, txn)
}

// CreateView creates a view.
func (mgr *Manager) CreateView(viewName domain.ViewName, viewDef domain.ViewDef, txn domain.Transaction) error {
	return mgr.viewMgr.CreateView(viewName, viewDef, txn)
}

// GetViewDef returns given view definition.
func (mgr *Manager) GetViewDef(viewName domain.ViewName, txn domain.Transaction) (domain.ViewDef, error) {
	return mgr.viewMgr.GetViewDef(viewName, txn)
}

// CreateIndex creates an index.
func (mgr *Manager) CreateIndex(idxName domain.IndexName, tblName domain.TableName, fldName domain.FieldName, txn domain.Transaction) error {
	return mgr.idxMgr.CreateIndex(idxName, tblName, fldName, txn)
}

// GetIndexInfo returns the index information of given table.
func (mgr *Manager) GetIndexInfo(tblName domain.TableName, txn domain.Transaction) (map[domain.FieldName]*domain.IndexInfo, error) {
	return mgr.idxMgr.GetIndexInfo(tblName, txn)
}

// GetStatInfo returns the statistical information of given table.
func (mgr *Manager) GetStatInfo(tblName domain.TableName, layout *domain.Layout, txn domain.Transaction) (domain.StatInfo, error) {
	return mgr.statMgr.GetStatInfo(tblName, layout, txn)
}
