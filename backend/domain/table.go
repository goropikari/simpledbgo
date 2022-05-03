package domain

import (
	"log"

	"github.com/pkg/errors"
)

// Table is a model of database table.
type Table struct {
	txn           Transaction
	layout        *Layout
	RecordPage    *RecordPage
	tblName       TableName
	currentSlotID SlotID
	err           error
}

// NewTable constructs a Table.
func NewTable(txn Transaction, tblName TableName, layout *Layout) (*Table, error) {
	tbl := &Table{
		txn:           txn,
		layout:        layout,
		tblName:       tblName,
		RecordPage:    nil,
		currentSlotID: -1,
		err:           nil,
	}

	blkLen, err := txn.BlockLength(FileName(tblName))
	if err != nil {
		return nil, err
	}

	if blkLen == 0 {
		err := tbl.moveToNewBlock()
		if err != nil {
			return nil, err
		}
	} else {
		err := tbl.moveToBlock(0)
		if err != nil {
			return nil, err
		}
	}

	return tbl, nil
}

// Err returns iteration err.
func (tbl *Table) Err() error {
	return tbl.err
}

// Close closes the table.
func (tbl *Table) Close() {
	if tbl.RecordPage != nil {
		tbl.txn.Unpin(tbl.RecordPage.Block())
	}
}

// GetInt32 gets int32 from the table.
func (tbl *Table) GetInt32(fldname FieldName) (int32, error) {
	return tbl.RecordPage.GetInt32(tbl.currentSlotID, fldname)
}

// GetString gets string from the table.
func (tbl *Table) GetString(fldname FieldName) (string, error) {
	return tbl.RecordPage.GetString(tbl.currentSlotID, fldname)
}

// GetVal gets value from the table.
func (tbl *Table) GetVal(fldname FieldName) (Constant, error) {
	typ := tbl.layout.schema.Type(fldname)
	switch typ {
	case Int32:
		val, err := tbl.GetInt32(fldname)
		if err != nil {
			return Constant{}, err
		}

		return NewConstant(VInt32, val), nil
	case String:
		val, err := tbl.GetString(fldname)
		if err != nil {
			return Constant{}, err
		}

		return NewConstant(VString, val), nil
	case Unknown:
		log.Fatal(errors.New("unexpected field type"))
	}

	return Constant{}, errors.New("GetVal error")
}

// SetInt32 sets int32 to the table.
func (tbl *Table) SetInt32(fldname FieldName, val int32) error {
	return tbl.RecordPage.SetInt32(tbl.currentSlotID, fldname, val)
}

// SetString sets string to the table.
func (tbl *Table) SetString(fldname FieldName, val string) error {
	return tbl.RecordPage.SetString(tbl.currentSlotID, fldname, val)
}

// SetVal sets value to the table.
func (tbl *Table) SetVal(fldname FieldName, val Constant) error {
	typ := tbl.layout.schema.Type(fldname)
	switch typ {
	case Int32:
		// TODO: check val type?
		err := tbl.SetInt32(fldname, val.ToInt32())
		if err != nil {
			return err
		}
	case String:
		err := tbl.SetString(fldname, val.ToString())
		if err != nil {
			return err
		}
	case Unknown:
		log.Fatal(errors.New("unexpected field type"))
	}

	return nil
}

// AdvanceNextInsertSlotID  advances current slot id to next to unused slot id.
// If there is no unused record, append file block.
func (tbl *Table) AdvanceNextInsertSlotID() error {
	slotID, err := tbl.RecordPage.InsertAfter(tbl.currentSlotID)
	if err != nil {
		return err
	}
	tbl.currentSlotID = slotID

	for tbl.currentSlotID < 0 {
		last, err := tbl.isAtLastBlock()
		if err != nil {
			return err
		}
		if last {
			err = tbl.moveToNewBlock()
			if err != nil {
				return err
			}
		} else {
			blk := tbl.RecordPage.Block()
			blkNum := blk.Number()
			err := tbl.moveToBlock(blkNum + 1)
			if err != nil {
				return err
			}
		}
		slotID, err := tbl.RecordPage.InsertAfter(tbl.currentSlotID)
		if err != nil {
			return err
		}
		tbl.currentSlotID = slotID
	}

	return nil
}

// Delete deletes the current slot logically.
func (tbl *Table) Delete() error {
	return tbl.RecordPage.Delete(tbl.currentSlotID)
}

// // moveToRecordID moves to the record id.
// func (tbl *Table) moveToRecordID(rid RecordID) error {
// 	tbl.Close()
// 	blk := NewBlock(tbl.tblName, tbl.txn.BlockSize(), rid.BlockNumber())
// 	RecordPage, err := NewPage(tbl.txn, blk, tbl.layout)
// 	if err != nil {
// 		return err
// 	}
// 	tbl.RecordPage = RecordPage
// 	tbl.currentSlotID = rid.SlotID()

// 	return nil
// }

// RecordID is a identifier of record.
func (tbl *Table) RecordID() RecordID {
	blk := tbl.RecordPage.Block()

	return NewRecordID(blk.Number(), tbl.currentSlotID)
}

// HasField checks the existence of the field.
func (tbl *Table) HasField(fldname FieldName) bool {
	return tbl.layout.schema.HasField(fldname)
}

// HasNextUsedSlot checks the existence of next used slot.
func (tbl *Table) HasNextUsedSlot() bool {
	currentSlotID, err := tbl.RecordPage.NextAfter(tbl.currentSlotID)
	if err != nil {
		tbl.err = err

		return false
	}
	tbl.currentSlotID = currentSlotID

	for tbl.currentSlotID < 0 {
		last, err := tbl.isAtLastBlock()
		if err != nil {
			tbl.err = err

			return false
		}
		if last {
			return false
		}

		blk := tbl.RecordPage.Block()
		if err := tbl.moveToBlock(blk.Number() + 1); err != nil {
			tbl.err = err

			return false
		}

		slotID, err := tbl.RecordPage.NextAfter(tbl.currentSlotID)
		if err != nil {
			tbl.err = err

			return false
		}
		tbl.currentSlotID = slotID
	}

	return true
}

// isAtLastBlock checks whether the current block is last block or not.
func (tbl *Table) isAtLastBlock() (bool, error) {
	blk := tbl.RecordPage.Block()
	size, err := tbl.txn.BlockLength(FileName(tbl.tblName))
	if err != nil {
		return false, err
	}

	blkNum, err := NewBlockNumber(size - 1)
	if err != nil {
		return false, err
	}

	return blk.Number() == blkNum, nil
}

func (tbl *Table) moveToNewBlock() error {
	tbl.Close()
	blk, err := tbl.txn.ExtendFile(FileName(tbl.tblName))
	if err != nil {
		return err
	}

	RecordPage, err := NewRecordPage(tbl.txn, blk, tbl.layout)
	if err != nil {
		return err
	}

	tbl.RecordPage = RecordPage

	err = tbl.RecordPage.Format()
	if err != nil {
		return err
	}

	tbl.currentSlotID = -1

	return nil
}

// MoveToFirst move to the first block of the table.
func (tbl *Table) MoveToFirst() error {
	return tbl.moveToBlock(0)
}

func (tbl *Table) moveToBlock(blkNum BlockNumber) error {
	tbl.Close()
	blk := NewBlock(FileName(tbl.tblName), blkNum)
	RecordPage, err := NewRecordPage(tbl.txn, blk, tbl.layout)
	if err != nil {
		return err
	}

	tbl.RecordPage = RecordPage
	tbl.currentSlotID = -1

	return nil
}
