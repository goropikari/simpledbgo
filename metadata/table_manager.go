package metadata

import (
	"fmt"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/errors"
)

// TableManager is manager of table.
type TableManager struct {
	tblCatalogLayout *domain.Layout
	fldCatalogLayout *domain.Layout
}

// NewTableManager constructs TableManager.
func NewTableManager() *TableManager {
	tblCatalogSchema := domain.NewSchema()
	tblCatalogSchema.AddStringField(domain.FieldName(fldTableName), domain.MaxTableNameLength)
	tblCatalogSchema.AddInt32Field(domain.FieldName(fldSlotSize))
	tblCatalogLayout := domain.NewLayout(tblCatalogSchema)

	fldCatalogSchema := domain.NewSchema()
	fldCatalogSchema.AddStringField(fldTableName, domain.MaxTableNameLength)
	fldCatalogSchema.AddStringField(fldFieldName, domain.MaxFieldNameLength)
	fldCatalogSchema.AddInt32Field(fldType)
	fldCatalogSchema.AddInt32Field(fldLength)
	fldCatalogSchema.AddInt32Field(fldOffset)
	fldCatalogLayout := domain.NewLayout(fldCatalogSchema)

	return &TableManager{
		tblCatalogLayout: tblCatalogLayout,
		fldCatalogLayout: fldCatalogLayout,
	}
}

// CreateTableManager creates table manager and table catalog.
func CreateTableManager(txn domain.Transaction) (*TableManager, error) {
	tblMgr := NewTableManager()
	if err := tblMgr.CreateTable(tableCatalog, tblMgr.tblCatalogLayout.Schema(), txn); err != nil {
		return nil, errors.Err(err, "CreateTable")
	}
	if err := tblMgr.CreateTable(fieldCatalog, tblMgr.fldCatalogLayout.Schema(), txn); err != nil {
		return nil, errors.Err(err, "CreateTable")
	}

	return tblMgr, nil
}

// // TableCatalogLayout returns layout of table catalog.
// func (tblMgr *TableManager) TableCatalogLayout() *domain.Layout {
// 	return tblMgr.tblCatalogLayout
// }

// // FieldCatalogLayout returns layout of field catalog.
// func (tblMgr *TableManager) FieldCatalogLayout() *domain.Layout {
// 	return tblMgr.fldCatalogLayout
// }

// CreateTable create a table.
func (tblMgr *TableManager) CreateTable(tblName domain.TableName, sch *domain.Schema, txn domain.Transaction) error {
	layout := domain.NewLayout(sch)

	// register table
	tcat, err := domain.NewTableScan(txn, tableCatalog, tblMgr.tblCatalogLayout)
	if err != nil {
		return errors.Err(err, "NewTableScan")
	}
	if err := tcat.AdvanceNextInsertSlotID(); err != nil {
		return errors.Err(err, "AdvanceNextInsertSlotID")
	}

	if err := tcat.SetString(fldTableName, tblName.String()); err != nil {
		return errors.Err(err, "SetString")
	}
	if err := tcat.SetInt32(fldSlotSize, int32(layout.SlotSize())); err != nil {
		return errors.Err(err, "SetInt32")
	}
	tcat.Close()

	// register fields
	fcat, err := domain.NewTableScan(txn, fieldCatalog, tblMgr.fldCatalogLayout)
	if err != nil {
		return errors.Err(err, "NewTableScan")
	}

	for _, fld := range sch.Fields() {
		if err := fcat.AdvanceNextInsertSlotID(); err != nil {
			return errors.Err(err, "AdvanceNextInsertSlotID")
		}
		if err := fcat.SetString(fldTableName, tblName.String()); err != nil {
			return errors.Err(err, "SetString")
		}
		if err := fcat.SetString(fldFieldName, string(fld)); err != nil {
			return errors.Err(err, "SetString")
		}
		if err := fcat.SetInt32(fldType, int32(sch.Type(fld))); err != nil {
			return errors.Err(err, "SetInt32")
		}
		if err := fcat.SetInt32(fldLength, int32(sch.Length(fld))); err != nil {
			return errors.Err(err, "SetInt32")
		}
		if err := fcat.SetInt32(fldOffset, int32(layout.Offset(fld))); err != nil {
			return errors.Err(err, "SetInt32")
		}
	}
	fcat.Close()

	return nil
}

// GetTableLayout returns the layout of given table name.
func (tblMgr *TableManager) GetTableLayout(tblName domain.TableName, txn domain.Transaction) (*domain.Layout, error) {
	slotsize, err := tblMgr.tableSlotSize(tblName, txn)
	if err != nil {
		return nil, errors.Err(err, "tableSlotSize")
	}

	sch, offsets, err := tblMgr.tableSchema(tblName, txn)
	if err != nil {
		return nil, errors.Err(err, "tableSchema")
	}

	return domain.NewLayoutWithFields(sch, offsets, int64(slotsize)), nil
}

// Exists checks the existence of table.
func (tblMgr *TableManager) Exists(tblName domain.TableName, txn domain.Transaction) bool {
	tcat, err := domain.NewTableScan(txn, tableCatalog, tblMgr.tblCatalogLayout)
	if err != nil {
		return false
	}
	defer tcat.Close()

	for tcat.HasNext() {
		v, err := tcat.GetString(fldTableName)
		if err != nil {
			return false
		}
		if v == tblName.String() {
			return true
		}
	}
	if err := tcat.Err(); err != nil {
		return false
	}

	return false
}

func (tblMgr *TableManager) tableSlotSize(tblName domain.TableName, txn domain.Transaction) (int32, error) {
	const NonExistSlotSize = -1

	tcat, err := domain.NewTableScan(txn, tableCatalog, tblMgr.tblCatalogLayout)
	if err != nil {
		return NonExistSlotSize, errors.Err(err, "NewTableScan")
	}
	defer tcat.Close()

	slotsize := int32(NonExistSlotSize)
	for tcat.HasNext() {
		v, err := tcat.GetString(fldTableName)
		if err != nil {
			return NonExistSlotSize, errors.Err(err, "GetString")
		}
		if v == tblName.String() {
			slotsize, err = tcat.GetInt32(fldSlotSize)
			if err != nil {
				return NonExistSlotSize, errors.Err(err, "GetInt32")
			}

			break
		}
	}
	if err := tcat.Err(); err != nil {
		return NonExistSlotSize, errors.Err(err, "HasNext")
	}

	if slotsize <= 0 {
		return NonExistSlotSize, fmt.Errorf("relation \"%v\" does not exist", tblName)
	}

	return slotsize, nil
}

func (tblMgr *TableManager) tableSchema(tblName domain.TableName, txn domain.Transaction) (*domain.Schema, map[domain.FieldName]int64, error) {
	sch := domain.NewSchema()
	offsets := make(map[domain.FieldName]int64)
	fcat, err := domain.NewTableScan(txn, fieldCatalog, tblMgr.fldCatalogLayout)
	if err != nil {
		return nil, nil, errors.Err(err, "NewTableScan")
	}
	defer fcat.Close()

	for fcat.HasNext() {
		v, err := fcat.GetString(fldTableName)
		if err != nil {
			return nil, nil, errors.Err(err, "GetString")
		}
		if v == tblName.String() {
			fldNameStr, err := fcat.GetString(fldFieldName)
			if err != nil {
				return nil, nil, errors.Err(err, "GetString")
			}
			typ, err := fcat.GetInt32(fldType)
			if err != nil {
				return nil, nil, errors.Err(err, "GetInt32")
			}
			length, err := fcat.GetInt32(fldLength)
			if err != nil {
				return nil, nil, errors.Err(err, "GetInt32")
			}
			offset, err := fcat.GetInt32(fldOffset)
			if err != nil {
				return nil, nil, errors.Err(err, "GetInt32")
			}

			fldName, err := domain.NewFieldName(fldNameStr)
			if err != nil {
				return nil, nil, errors.Err(err, "NewFieldName")
			}
			offsets[fldName] = int64(offset)

			fldType := domain.FieldType(typ)
			sch.AddField(fldName, fldType, int(length))
		}
	}
	if err := fcat.Err(); err != nil {
		return nil, nil, errors.Err(err, "HasNext")
	}

	return sch, offsets, nil
}
