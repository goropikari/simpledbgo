package metadata

import (
	"github.com/goropikari/simpledbgo/backend/domain"
	"github.com/goropikari/simpledbgo/backend/record"
)

// TableManager is manager of table.
type TableManager struct {
	tblCatalogLayout *record.Layout
	fldCatalogLayout *record.Layout
}

// NewTableManager constructs TableManager.
func NewTableManager() *TableManager {
	tblCatalogSchema := record.NewSchema()
	tblCatalogSchema.AddStringField(domain.FieldName(fldTableName), maxTableNameLength)
	tblCatalogSchema.AddInt32Field(domain.FieldName(fldSlotSize))
	tblCatalogLayout := record.NewLayout(tblCatalogSchema)

	fldCatalogSchema := record.NewSchema()
	fldCatalogSchema.AddStringField(fldTableName, maxTableNameLength)
	fldCatalogSchema.AddStringField(fldFieldName, maxFieldNameLength)
	fldCatalogSchema.AddInt32Field(fldType)
	fldCatalogSchema.AddInt32Field(fldLength)
	fldCatalogSchema.AddInt32Field(fldOffset)
	fldCatalogLayout := record.NewLayout(fldCatalogSchema)

	return &TableManager{
		tblCatalogLayout: tblCatalogLayout,
		fldCatalogLayout: fldCatalogLayout,
	}
}

// CreateTableManager creates table manager and table catalog.
func CreateTableManager(txn domain.Transaction) (*TableManager, error) {
	tblMgr := NewTableManager()
	if err := tblMgr.CreateTable(tableCatalog, tblMgr.tblCatalogLayout.Schema(), txn); err != nil {
		return nil, err
	}
	if err := tblMgr.CreateTable(fieldCatalog, tblMgr.fldCatalogLayout.Schema(), txn); err != nil {
		return nil, err
	}

	return tblMgr, nil
}

// // TableCatalogLayout returns layout of table catalog.
// func (tblMgr *TableManager) TableCatalogLayout() *record.Layout {
// 	return tblMgr.tblCatalogLayout
// }

// // FieldCatalogLayout returns layout of field catalog.
// func (tblMgr *TableManager) FieldCatalogLayout() *record.Layout {
// 	return tblMgr.fldCatalogLayout
// }

// CreateTable create a table.
func (tblMgr *TableManager) CreateTable(tblName domain.FileName, sch *record.Schema, txn domain.Transaction) error {
	layout := record.NewLayout(sch)

	// register table
	tcat, err := record.NewTable(txn, tableCatalog, tblMgr.tblCatalogLayout)
	if err != nil {
		return err
	}
	if err := tcat.AdvanceNextInsertSlotID(); err != nil {
		return err
	}

	if err := tcat.SetString(fldTableName, tblName.String()); err != nil {
		return err
	}
	if err := tcat.SetInt32(fldSlotSize, int32(layout.SlotSize())); err != nil {
		return err
	}
	tcat.Close()

	// register fields
	fcat, err := record.NewTable(txn, fieldCatalog, tblMgr.fldCatalogLayout)
	if err != nil {
		return err
	}

	for _, fld := range sch.Fields() {
		if err := fcat.AdvanceNextInsertSlotID(); err != nil {
			return err
		}
		if err := fcat.SetString(fldTableName, tblName.String()); err != nil {
			return err
		}
		if err := fcat.SetString(fldFieldName, string(fld)); err != nil {
			return err
		}
		if err := fcat.SetInt32(fldType, int32(sch.Type(fld))); err != nil {
			return err
		}
		if err := fcat.SetInt32(fldLength, int32(sch.Length(fld))); err != nil {
			return err
		}
		if err := fcat.SetInt32(fldOffset, int32(layout.Offset(fld))); err != nil {
			return err
		}
	}
	fcat.Close()

	return nil
}

// GetTableLayout returns the layout of given table name.
func (tblMgr *TableManager) GetTableLayout(tblName domain.FileName, txn domain.Transaction) (*record.Layout, error) {
	slotsize, err := tblMgr.tableSlotSize(tblName, txn)
	if err != nil {
		return nil, err
	}

	sch, offsets, err := tblMgr.tableSchema(tblName, txn)
	if err != nil {
		return nil, err
	}

	return record.NewLayoutWithFields(sch, offsets, int64(slotsize)), nil
}

func (tblMgr *TableManager) tableSlotSize(tblName domain.FileName, txn domain.Transaction) (int32, error) {
	slotsize := int32(-1)
	tcat, err := record.NewTable(txn, tableCatalog, tblMgr.tblCatalogLayout)
	if err != nil {
		return -1, err
	}
	defer tcat.Close()

	for {
		found, err := tcat.HasNextUsedSlot()
		if err != nil {
			return -1, err
		}
		if !found {
			break
		}

		v, err := tcat.GetString(fldTableName)
		if err != nil {
			return -1, err
		}
		if v == tblName.String() {
			slotsize, err = tcat.GetInt32(fldSlotSize)
			if err != nil {
				return -1, err
			}

			break
		}
	}

	return slotsize, nil
}

func (tblMgr *TableManager) tableSchema(tblName domain.FileName, txn domain.Transaction) (*record.Schema, map[domain.FieldName]int64, error) {
	sch := record.NewSchema()
	offsets := make(map[domain.FieldName]int64)
	fcat, err := record.NewTable(txn, fieldCatalog, tblMgr.fldCatalogLayout)
	if err != nil {
		return nil, nil, err
	}
	defer fcat.Close()

	for {
		found, err := fcat.HasNextUsedSlot()
		if err != nil {
			return nil, nil, err
		}
		if !found {
			break
		}

		v, err := fcat.GetString(fldTableName)
		if err != nil {
			return nil, nil, err
		}
		if v == tblName.String() {
			fldNameStr, err := fcat.GetString(fldFieldName)
			if err != nil {
				return nil, nil, err
			}
			typ, err := fcat.GetInt32(fldType)
			if err != nil {
				return nil, nil, err
			}
			length, err := fcat.GetInt32(fldLength)
			if err != nil {
				return nil, nil, err
			}
			offset, err := fcat.GetInt32(fldOffset)
			if err != nil {
				return nil, nil, err
			}

			fldName, err := domain.NewFieldName(fldNameStr)
			if err != nil {
				return nil, nil, err
			}
			offsets[fldName] = int64(offset)

			fldType := record.FieldType(typ)
			sch.AddField(fldName, fldType, int(length))
		}
	}

	return sch, offsets, nil
}
