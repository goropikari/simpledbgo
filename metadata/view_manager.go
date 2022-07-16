package metadata

import (
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/errors"
)

// ViewManager is manager of view.
type ViewManager struct {
	tblMgr *TableManager
}

// NewViewManager constructs a ViewManager.
func NewViewManager(tblMgr *TableManager) *ViewManager {
	return &ViewManager{
		tblMgr: tblMgr,
	}
}

// CreateViewManager creates a view manager and view catalog.
func CreateViewManager(tblMgr *TableManager, txn domain.Transaction) (*ViewManager, error) {
	viewMgr := NewViewManager(tblMgr)
	sch := domain.NewSchema()
	sch.AddStringField(fldViewName, domain.MaxTableNameLength)
	sch.AddStringField(fldViewDef, domain.MaxViewDefLength)
	if err := tblMgr.CreateTable(fldViewCatalog, sch, txn); err != nil {
		return nil, errors.Err(err, "CreateTable")
	}

	return viewMgr, nil
}

// CreateView defines a view.
func (viewMgr *ViewManager) CreateView(vName domain.ViewName, vDef domain.ViewDef, txn domain.Transaction) error {
	layout, err := viewMgr.tblMgr.GetTableLayout(fldViewCatalog, txn)
	if err != nil {
		return errors.Err(err, "GetTableLayout")
	}

	tbl, err := domain.NewTableScan(txn, fldViewCatalog, layout)
	if err != nil {
		return errors.Err(err, "NewTableScan")
	}
	if err := tbl.AdvanceNextInsertSlotID(); err != nil {
		return errors.Err(err, "AdvanceNextInsertSlotID")
	}
	if err := tbl.SetString(fldViewName, vName.String()); err != nil {
		return errors.Err(err, "SetString")
	}
	if err := tbl.SetString(fldViewDef, vDef.String()); err != nil {
		return errors.Err(err, "SetString")
	}
	tbl.Close()

	return nil
}

// GetViewDef gets the definition of view.
func (viewMgr *ViewManager) GetViewDef(viewName domain.ViewName, txn domain.Transaction) (domain.ViewDef, error) {
	layout, err := viewMgr.tblMgr.GetTableLayout(fldViewCatalog, txn)
	if err != nil {
		return "", errors.Err(err, "GetTableLayout")
	}

	tbl, err := domain.NewTableScan(txn, fldViewCatalog, layout)
	if err != nil {
		return "", errors.Err(err, "NewTableScan")
	}
	defer tbl.Close()

	defStr := ""
	for tbl.HasNext() {
		view, err := tbl.GetString(fldViewName)
		if err != nil {
			return "", errors.Err(err, "GetString")
		}

		if view == viewName.String() {
			defStr, err = tbl.GetString(fldViewDef)
			if err != nil {
				return "", errors.Err(err, "GetString")
			}

			break
		}
	}
	if err := tbl.Err(); err != nil {
		return "", errors.Err(err, "GetString")
	}

	return domain.NewViewDef(defStr), nil
}
