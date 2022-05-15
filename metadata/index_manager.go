package metadata

import (
	"github.com/goropikari/simpledbgo/domain"
)

// IndexManager is an index manager.
type IndexManager struct {
	layout     *domain.Layout
	tblMgr     *TableManager
	statMgr    *StatManager
	idxFactory domain.IndexFactory
}

// NewIndexManager constructs an index manager.
func NewIndexManager(factory domain.IndexFactory, tblMgr *TableManager, statMgr *StatManager, txn domain.Transaction) (*IndexManager, error) {
	cat, err := tblMgr.GetTableLayout(fldIndexCatalog, txn)
	if err != nil {
		return nil, err
	}

	return &IndexManager{
		layout:     cat,
		tblMgr:     tblMgr,
		statMgr:    statMgr,
		idxFactory: factory,
	}, nil
}

// CreateIndexManager creates new index manager.
func CreateIndexManager(factory domain.IndexFactory, tblMgr *TableManager, statMgr *StatManager, txn domain.Transaction) (*IndexManager, error) {
	sch := domain.NewSchema()
	sch.AddStringField(fldIndexName, domain.MaxIndexNameLength)
	sch.AddStringField(fldTableName, domain.MaxTableNameLength)
	sch.AddStringField(fldFieldName, domain.MaxFieldNameLength)
	if err := tblMgr.CreateTable(fldIndexCatalog, sch, txn); err != nil {
		return nil, err
	}

	layout, err := tblMgr.GetTableLayout(fldIndexCatalog, txn)
	if err != nil {
		return nil, err
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
	tbl, err := domain.NewTable(txn, fldIndexCatalog, idxMgr.layout)
	if err != nil {
		return err
	}
	defer tbl.Close()

	if err := tbl.AdvanceNextInsertSlotID(); err != nil {
		return err
	}

	err = tbl.SetString(fldIndexName, idxName.String())
	if err != nil {
		return err
	}

	err = tbl.SetString(fldTableName, tblName.String())
	if err != nil {
		return err
	}

	err = tbl.SetString(fldFieldName, fldName.String())
	if err != nil {
		return err
	}

	return nil
}

// GetIndexInfo returns all index information of given table.
func (idxMgr *IndexManager) GetIndexInfo(tblName domain.TableName, txn domain.Transaction) (map[domain.FieldName]*domain.IndexInfo, error) {
	infos := make(map[domain.FieldName]*domain.IndexInfo)
	tbl, err := domain.NewTable(txn, fldIndexCatalog, idxMgr.layout)
	if err != nil {
		return nil, err
	}
	defer tbl.Close()

	for tbl.HasNext() {
		storedTblName, err := tbl.GetString(fldTableName)
		if err != nil {
			return nil, err
		}
		if storedTblName != tblName.String() {
			continue
		}

		idxNameStr, err := tbl.GetString(fldIndexName)
		if err != nil {
			return nil, err
		}
		idxName, err := domain.NewIndexName(idxNameStr)
		if err != nil {
			return nil, err
		}

		fldNameStr, err := tbl.GetString(fldFieldName)
		if err != nil {
			return nil, err
		}
		fldName, err := domain.NewFieldName(fldNameStr)
		if err != nil {
			return nil, err
		}

		tblLayout, err := idxMgr.tblMgr.GetTableLayout(tblName, txn)
		if err != nil {
			return nil, err
		}

		tblsi, err := idxMgr.statMgr.GetStatInfo(tblName, tblLayout, txn)
		if err != nil {
			return nil, err
		}

		idxInfo := domain.NewIndexInfo(idxMgr.idxFactory, idxName, fldName, tblLayout.Schema(), txn, tblsi)

		infos[fldName] = idxInfo
	}
	if tbl.Err() != nil {
		return nil, tbl.Err()
	}

	return infos, nil
}
