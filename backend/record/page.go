package record

import (
	"log"

	"github.com/goropikari/simpledbgo/backend/domain"
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

// SlotID is identifier of slot.
type SlotID = int64

// Page is a model of page.
type Page struct {
	txn    domain.Transaction
	blk    domain.Block
	layout *Layout
}

// NewPage constructs a page.
func NewPage(txn domain.Transaction, blk domain.Block, layout *Layout) (*Page, error) {
	if err := txn.Pin(blk); err != nil {
		return nil, err
	}

	return &Page{
		txn:    txn,
		blk:    blk,
		layout: layout,
	}, nil
}

// GetInt32 gets int32 from the block.
func (page *Page) GetInt32(slot SlotID, fldname FieldName) (int32, error) {
	offset := page.offset(slot) + page.layout.Offset(fldname)

	return page.txn.GetInt32(page.blk, offset)
}

// SetInt32 sets int32 to the block.
func (page *Page) SetInt32(slot SlotID, fldname FieldName, val int32) error {
	offset := page.offset(slot) + page.layout.Offset(fldname)

	return page.txn.SetInt32(page.blk, offset, val, true)
}

// GetString gets string from the block.
func (page *Page) GetString(slot SlotID, fldname FieldName) (string, error) {
	offset := page.offset(slot) + page.layout.Offset(fldname)

	return page.txn.GetString(page.blk, offset)
}

// SetString sets the string from the block.
func (page *Page) SetString(slot SlotID, fldname FieldName, val string) error {
	offset := page.offset(slot) + page.layout.Offset(fldname)

	return page.txn.SetString(page.blk, offset, val, true)
}

// Delete deletes the slot.
func (page *Page) Delete(slot SlotID) error {
	return page.setFlag(slot, Empty)
}

// Format formats blk.
func (page *Page) Format() error {
	slot := int64(0)
	for page.isValidSlot(slot) {
		if err := page.txn.SetInt32(page.blk, page.offset(slot), Empty, false); err != nil {
			return err
		}

		sch := page.layout.schema
		for _, fldname := range sch.fields {
			typ := sch.Type(fldname)
			fldpos := page.offset(slot) + page.layout.Offset(fldname)
			switch typ {
			case Int32:
				if err := page.txn.SetInt32(page.blk, fldpos, 0, false); err != nil {
					return err
				}
			case String:
				if err := page.txn.SetString(page.blk, fldpos, "", false); err != nil {
					return err
				}
			case Unknown:
				log.Fatal(errors.New("unexpected record type"))
			}
		}
		slot++
	}

	return nil
}

// NextAfter returns the slot id with Used flag after slot.
func (page *Page) NextAfter(slot SlotID) (SlotID, error) {
	return page.searchAfter(slot, Used)
}

// InsertAfter searches the slot id after slot with Empty flag, set Used flag and returns its id.
func (page *Page) InsertAfter(slot SlotID) (SlotID, error) {
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
func (page *Page) searchAfter(slot SlotID, flag SlotCondition) (SlotID, error) {
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

func (page *Page) isValidSlot(slot SlotID) bool {
	return page.offset(slot+1) <= int64(page.txn.BlockSize())
}

func (page *Page) setFlag(slot SlotID, flag SlotCondition) error {
	return page.txn.SetInt32(page.blk, page.offset(slot), flag, true)
}

func (page *Page) offset(slot int64) int64 {
	return slot * page.layout.slotsize
}

// Block returns page's block.
func (page *Page) Block() domain.Block {
	return page.blk
}
