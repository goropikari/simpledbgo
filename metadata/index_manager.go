package metadata

import (
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/errors"
)

// IndexManager is an index manager.
type IndexManager struct {
	layout     *domain.Layout
	tblMgr     *TableManager
	statMgr    *StatManager
	idxFactory domain.IndexDriver
}

// NewIndexManager constructs an index manager.
func NewIndexManager(factory domain.IndexDriver, tblMgr *TableManager, statMgr *StatManager, txn domain.Transaction) (*IndexManager, error) {
	cat, err := tblMgr.GetTableLayout(fldIndexCatalog, txn)
	if err != nil {
		return nil, errors.Err(err, "GetTableLayout")
	}

	return &IndexManager{
		layout:     cat,
		tblMgr:     tblMgr,
		statMgr:    statMgr,
		idxFactory: factory,
	}, nil
}

// CreateIndexManager creates new index manager.
func CreateIndexManager(factory domain.IndexDriver, tblMgr *TableManager, statMgr *StatManager, txn domain.Transaction) (*IndexManager, error) {
	sch := domain.NewSchema()
	sch.AddStringField(fldIndexName, domain.MaxIndexNameLength)
	sch.AddStringField(fldTableName, domain.MaxTableNameLength)
	sch.AddStringField(fldFieldName, domain.MaxFieldNameLength)
	if err := tblMgr.CreateTable(fldIndexCatalog, sch, txn); err != nil {
		return nil, err
	}

	layout, err := tblMgr.GetTableLayout(fldIndexCatalog, txn)
	if err != nil {
		return nil, errors.Err(err, "GetTableLayout")
	}

	return &IndexManager{
		layout:     layout,
		tblMgr:     tblMgr,
		statMgr:    statMgr,
		idxFactory: factory,
	}, nil
}

// CreateIndex creates an index.
func (idxMgr *IndexManager) CreateIndex(idxName domain.IndexName, tblName domain.TableName, fldName domain.FieldName, txn domain.Transaction) error {
	tbl, err := domain.NewTableScan(txn, fldIndexCatalog, idxMgr.layout)
	if err != nil {
		return errors.Err(err, "NewTableScan")
	}
	defer tbl.Close()

	if err := tbl.AdvanceNextInsertSlotID(); err != nil {
		return errors.Err(err, "AdvanceNextInsertSlotID")
	}

	err = tbl.SetString(fldIndexName, idxName.String())
	if err != nil {
		return errors.Err(err, "SetString")
	}

	err = tbl.SetString(fldTableName, tblName.String())
	if err != nil {
		return errors.Err(err, "SetString")
	}

	err = tbl.SetString(fldFieldName, fldName.String())
	if err != nil {
		return errors.Err(err, "SetString")
	}

	return nil
}

// GetIndexInfo returns all index information of given table.
func (idxMgr *IndexManager) GetIndexInfo(tblName domain.TableName, txn domain.Transaction) (map[domain.FieldName]*domain.IndexInfo, error) {
	infos := make(map[domain.FieldName]*domain.IndexInfo)
	tbl, err := domain.NewTableScan(txn, fldIndexCatalog, idxMgr.layout)
	if err != nil {
		return nil, errors.Err(err, "SetString")
	}
	defer tbl.Close()

	for tbl.HasNext() {
		storedTblName, err := tbl.GetString(fldTableName)
		if err != nil {
			return nil, errors.Err(err, "SetString")
		}
		if storedTblName != tblName.String() {
			continue
		}

		idxNameStr, err := tbl.GetString(fldIndexName)
		if err != nil {
			return nil, errors.Err(err, "SetString")
		}
		idxName, err := domain.NewIndexName(idxNameStr)
		if err != nil {
			return nil, errors.Err(err, "SetString")
		}

		fldNameStr, err := tbl.GetString(fldFieldName)
		if err != nil {
			return nil, errors.Err(err, "SetString")
		}
		fldName, err := domain.NewFieldName(fldNameStr)
		if err != nil {
			return nil, errors.Err(err, "SetString")
		}

		tblLayout, err := idxMgr.tblMgr.GetTableLayout(tblName, txn)
		if err != nil {
			return nil, errors.Err(err, "SetString")
		}

		tblsi, err := idxMgr.statMgr.GetStatInfo(tblName, tblLayout, txn)
		if err != nil {
			return nil, errors.Err(err, "SetString")
		}

		idxInfo := domain.NewIndexInfo(idxMgr.idxFactory, idxName, fldName, tblLayout.Schema(), txn, tblsi)

		infos[fldName] = idxInfo
	}
	if tbl.Err() != nil {
		return nil, errors.Err(tbl.Err(), "HasNext")
	}

	return infos, nil
}
