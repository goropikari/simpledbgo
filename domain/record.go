package domain

import (
	"log"

	"github.com/pkg/errors"
)

// SlotCondition is flag of record slot condition.
type SlotCondition = int32

const (
	// Empty means the slot is empty.
	Empty SlotCondition = iota
	// Used means the slot is used.
	Used
)

// RecordPage is a model of RecordPage.
type RecordPage struct {
	txn    Transaction
	blk    Block
	layout *Layout
}

// NewRecordPage constructs a RecordPage.
func NewRecordPage(txn Transaction, blk Block, layout *Layout) (*RecordPage, error) {
	if err := txn.Pin(blk); err != nil {
		return nil, err
	}

	return &RecordPage{
		txn:    txn,
		blk:    blk,
		layout: layout,
	}, nil
}

// GetInt32 gets int32 from the block.
func (page *RecordPage) GetInt32(slot SlotID, fldname FieldName) (int32, error) {
	offset := page.offset(slot) + page.layout.Offset(fldname)

	return page.txn.GetInt32(page.blk, offset)
}

// SetInt32 sets int32 to the block.
func (page *RecordPage) SetInt32(slot SlotID, fldname FieldName, val int32) error {
	offset := page.offset(slot) + page.layout.Offset(fldname)

	return page.txn.SetInt32(page.blk, offset, val, true)
}

// GetString gets string from the block.
func (page *RecordPage) GetString(slot SlotID, fldname FieldName) (string, error) {
	offset := page.offset(slot) + page.layout.Offset(fldname)

	return page.txn.GetString(page.blk, offset)
}

// SetString sets the string from the block.
func (page *RecordPage) SetString(slot SlotID, fldname FieldName, val string) error {
	offset := page.offset(slot) + page.layout.Offset(fldname)

	return page.txn.SetString(page.blk, offset, val, true)
}

// Delete deletes the slot.
func (page *RecordPage) Delete(slot SlotID) error {
	return page.setFlag(slot, Empty)
}

// Format formats blk.
func (page *RecordPage) Format() error {
	slot := SlotID(0)
	for page.isValidSlot(slot) {
		if err := page.txn.SetInt32(page.blk, page.offset(slot), Empty, false); err != nil {
			return err
		}

		sch := page.layout.schema
		for _, fldname := range sch.fields {
			typ := sch.Type(fldname)
			fldpos := page.offset(slot) + page.layout.Offset(fldname)
			switch typ {
			case Int32FieldType:
				if err := page.txn.SetInt32(page.blk, fldpos, 0, false); err != nil {
					return err
				}
			case StringFieldType:
				if err := page.txn.SetString(page.blk, fldpos, "", false); err != nil {
					return err
				}
			case UnknownFieldType:
				log.Fatal(errors.New("unexpected record type"))
			}
		}
		slot++
	}

	return nil
}

// NextAfter returns the slot id with Used flag after slot.
func (page *RecordPage) NextAfter(slot SlotID) (SlotID, error) {
	return page.searchAfter(slot, Used)
}

// InsertAfter searches the slot id after slot with Empty flag, set Used flag and returns its id.
func (page *RecordPage) InsertAfter(slot SlotID) (SlotID, error) {
	newSlot, err := page.searchAfter(slot, Empty)
	if err != nil {
		return 0, err
	}
	if newSlot >= 0 {
		err := page.setFlag(newSlot, Used)
		if err != nil {
			return 0, err
		}
	}

	return newSlot, nil
}

// searchAfter searches slot id with flag after slot.
func (page *RecordPage) searchAfter(slot SlotID, flag SlotCondition) (SlotID, error) {
	slot++
	for page.isValidSlot(slot) {
		curFlag, err := page.txn.GetInt32(page.blk, page.offset(slot))
		if err != nil {
			return 0, err
		}
		if flag == curFlag {
			return slot, nil
		}
		slot++
	}

	return -1, nil
}

func (page *RecordPage) isValidSlot(slot SlotID) bool {
	off := page.offset(slot + 1)
	x := int64(page.txn.BlockSize())

	return off <= x
}

func (page *RecordPage) setFlag(slot SlotID, flag SlotCondition) error {
	return page.txn.SetInt32(page.blk, page.offset(slot), flag, true)
}

func (page *RecordPage) offset(slot SlotID) int64 {
	return int64(slot) * page.layout.slotsize
}

// Block returns RecordPage's block.
func (page *RecordPage) Block() Block {
	return page.blk
}

// RecordID is identifier of record.
type RecordID struct {
	blkNum BlockNumber
	slotID SlotID
}

// NewRecordID constructs RecordID.
func NewRecordID(blkNum BlockNumber, slotID SlotID) RecordID {
	return RecordID{
		blkNum: blkNum,
		slotID: slotID,
	}
}

// BlockNumber returns block number.
func (rid RecordID) BlockNumber() BlockNumber {
	return rid.blkNum
}

// SlotID returns slot id.
func (rid RecordID) SlotID() SlotID {
	return rid.slotID
}

// Equal checks equality of rid and other.
func (rid RecordID) Equal(other RecordID) bool {
	return rid.blkNum == other.blkNum && rid.slotID == other.slotID
}
