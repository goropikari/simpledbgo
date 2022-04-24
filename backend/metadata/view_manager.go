package metadata

import (
	"github.com/goropikari/simpledbgo/backend/domain"
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
		return nil, err
	}

	return viewMgr, nil
}

// CreateView defines a view.
func (viewMgr *ViewManager) CreateView(vName ViewName, vDef ViewDef, txn domain.Transaction) error {
	layout, err := viewMgr.tblMgr.GetTableLayout(fldViewCatalog, txn)
	if err != nil {
		return err
	}

	tbl, err := domain.NewTable(txn, fldViewCatalog, layout)
	if err != nil {
		return err
	}
	if err := tbl.AdvanceNextInsertSlotID(); err != nil {
		return err
	}
	if err := tbl.SetString(fldViewName, vName); err != nil {
		return err
	}
	if err := tbl.SetString(fldViewDef, vDef); err != nil {
		return err
	}
	tbl.Close()

	return nil
}

// GetViewDef gets the definition of view.
func (viewMgr *ViewManager) GetViewDef(viewName ViewName, txn domain.Transaction) (ViewDef, error) {
	layout, err := viewMgr.tblMgr.GetTableLayout(fldViewCatalog, txn)
	if err != nil {
		return "", err
	}

	tbl, err := domain.NewTable(txn, fldViewCatalog, layout)
	if err != nil {
		return "", err
	}

	def := ""
	for {
		found, err := tbl.HasNextUsedSlot()
		if err != nil {
			return "", err
		}
		if !found {
			break
		}

		view, err := tbl.GetString(fldViewName)
		if err != nil {
			return "", err
		}

		if view == viewName {
			def, err = tbl.GetString(fldViewDef)
			if err != nil {
				return "", err
			}

			break
		}
	}

	tbl.Close()

	return def, nil
}
