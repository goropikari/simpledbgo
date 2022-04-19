package record

import (
	"log"

	"github.com/goropikari/simpledbgo/backend/domain"
	"github.com/goropikari/simpledbgo/meta"
	"github.com/pkg/errors"
)

// Table is a model of database table.
type Table struct {
	txn           domain.Transaction
	layout        *Layout
	rp            *Page
	filename      domain.FileName
	currentSlotID SlotID
}

// NewTable constructs a Table.
func NewTable(txn domain.Transaction, filename domain.FileName, layout *Layout) (*Table, error) {
	ts := &Table{
		txn:           txn,
		layout:        layout,
		filename:      filename,
		rp:            nil,
		currentSlotID: -1,
	}

	blkLen, err := txn.BlockLength(filename)
	if err != nil {
		return nil, err
	}

	if blkLen == 0 {
		err := ts.moveToNewBlock()
		if err != nil {
			return nil, err
		}
	} else {
		err := ts.moveToBlock(0)
		if err != nil {
			return nil, err
		}
	}

	return ts, nil
}

// Close closes the table.
func (ts *Table) Close() {
	if ts.rp != nil {
		ts.txn.Unpin(ts.rp.Block())
	}
}

// GetInt32 gets int32 from the table.
func (ts *Table) GetInt32(fldname FieldName) (int32, error) {
	return ts.rp.GetInt32(ts.currentSlotID, fldname)
}

// GetString gets string from the table.
func (ts *Table) GetString(fldname FieldName) (string, error) {
	return ts.rp.GetString(ts.currentSlotID, fldname)
}

// GetVal gets value from the table.
func (ts *Table) GetVal(fldname FieldName) (meta.Constant, error) {
	typ := ts.layout.schema.typ(fldname)
	switch typ {
	case Integer:
		val, err := ts.GetInt32(fldname)
		if err != nil {
			return meta.Constant{}, err
		}

		return meta.Constant{I32val: val}, nil
	case String:
		val, err := ts.GetString(fldname)
		if err != nil {
			return meta.Constant{}, err
		}

		return meta.Constant{Sval: val}, nil
	case Unknown:
		log.Fatal(errors.New("unexpected field type"))
	}

	return meta.Constant{}, errors.New("GetVal error")
}

// SetInt32 sets int32 to the table.
func (ts *Table) SetInt32(fldname FieldName, val int32) error {
	return ts.rp.SetInt32(ts.currentSlotID, fldname, val)
}

// SetString sets string to the table.
func (ts *Table) SetString(fldname FieldName, val string) error {
	return ts.rp.SetString(ts.currentSlotID, fldname, val)
}

// SetVal sets value to the table.
func (ts *Table) SetVal(fldname FieldName, val meta.Constant) error {
	typ := ts.layout.schema.typ(fldname)
	switch typ {
	case Integer:
		err := ts.SetInt32(fldname, val.I32val)
		if err != nil {
			return err
		}
	case String:
		err := ts.SetString(fldname, val.Sval)
		if err != nil {
			return err
		}
	case Unknown:
		log.Fatal(errors.New("unexpected field type"))
	}

	return errors.New("failed to SetVal")
}

// AdvanceNextInsertSlotID  advances current slot id to next to unused slot id.
// If there is no unused record, append file block.
func (ts *Table) AdvanceNextInsertSlotID() error {
	slotID, err := ts.rp.InsertAfter(ts.currentSlotID)
	if err != nil {
		return err
	}
	ts.currentSlotID = slotID

	for ts.currentSlotID < 0 {
		last, err := ts.isAtLastBlock()
		if err != nil {
			return err
		}
		if last {
			err = ts.moveToNewBlock()
			if err != nil {
				return err
			}
		} else {
			blk := ts.rp.Block()
			blkNum := blk.Number()
			err := ts.moveToBlock(blkNum + 1)
			if err != nil {
				return err
			}
		}
		slotID, err := ts.rp.InsertAfter(ts.currentSlotID)
		if err != nil {
			return err
		}
		ts.currentSlotID = slotID
	}

	return nil
}

// Delete deletes the current slot logically.
func (ts *Table) Delete() error {
	return ts.rp.Delete(ts.currentSlotID)
}

// // moveToRecordID moves to the record id.
// func (ts *Table) moveToRecordID(rid RecordID) error {
// 	ts.Close()
// 	blk := domain.NewBlock(ts.filename, ts.txn.BlockSize(), rid.BlockNumber())
// 	page, err := NewPage(ts.txn, blk, ts.layout)
// 	if err != nil {
// 		return err
// 	}
// 	ts.rp = page
// 	ts.currentSlotID = rid.SlotID()

// 	return nil
// }

// // RecordID is a identifier of record.
// func (ts *Table) RecordID() RecordID {
// 	blk := ts.rp.Block()

// 	return NewRecordID(blk.Number(), ts.currentSlotID)
// }

// HasField checks the existence of the field.
func (ts *Table) HasField(fldname FieldName) bool {
	return ts.layout.schema.HasField(fldname)
}

// HasNextUsedSlot checks the existence of next used slot.
func (ts *Table) HasNextUsedSlot() (bool, error) {
	currentSlotID, err := ts.rp.NextAfter(ts.currentSlotID)
	if err != nil {
		return false, err
	}
	ts.currentSlotID = currentSlotID

	for ts.currentSlotID < 0 {
		last, err := ts.isAtLastBlock()
		if err != nil {
			return false, err
		}
		if last {
			return false, nil
		}

		blk := ts.rp.Block()
		if err := ts.moveToBlock(blk.Number() + 1); err != nil {
			return false, err
		}

		slotID, err := ts.rp.NextAfter(ts.currentSlotID)
		if err != nil {
			return false, err
		}
		ts.currentSlotID = slotID
	}

	return true, nil
}

// isAtLastBlock checks whether the current block is last block or not.
func (ts *Table) isAtLastBlock() (bool, error) {
	blk := ts.rp.Block()
	size, err := ts.txn.BlockLength(ts.filename)
	if err != nil {
		return false, err
	}

	blkNum, err := domain.NewBlockNumber(size - 1)
	if err != nil {
		return false, err
	}

	return blk.Number() == blkNum, nil
}

func (ts *Table) moveToNewBlock() error {
	ts.Close()
	blk, err := ts.txn.ExtendFile(ts.filename)
	if err != nil {
		return err
	}

	page, err := NewPage(ts.txn, blk, ts.layout)
	if err != nil {
		return err
	}

	ts.rp = page

	err = ts.rp.Format()
	if err != nil {
		return err
	}

	ts.currentSlotID = -1

	return nil
}

// MoveToFirst move to the first block of the table.
func (ts *Table) MoveToFirst() error {
	return ts.moveToBlock(0)
}

func (ts *Table) moveToBlock(blkNum domain.BlockNumber) error {
	ts.Close()
	blk := domain.NewBlock(ts.filename, ts.txn.BlockSize(), blkNum)
	page, err := NewPage(ts.txn, blk, ts.layout)
	if err != nil {
		return err
	}

	ts.rp = page
	ts.currentSlotID = -1

	return nil
}
