package btree

import (
	"fmt"

	"github.com/goropikari/simpledbgo/common"
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/errors"
)

/*
	BTreePage
*/
const (
	pageFlagByteLength  = common.Int32Length
	numRecordByteLength = common.Int32Length
	pageFlagOffset      = 0
	numRecordOffset     = pageFlagByteLength
)

type pageFlag = int32

func newPageFlag(flag int32) pageFlag {
	return flag
}

// Page is page for btree node.
// page structure
// ------------------------------------------------------------------------------
// | page flag (int32) | # of records (int32) | record 1 | ... | record n | ... |
// ------------------------------------------------------------------------------
// 本では Slot といっているが domain.RecordPage と違って usage flag がない素の record を使っている
// 使わない record はそもそもアクセスしない想定になっているから usage flag をもたせていないのだと思われる.
type Page struct {
	txn     domain.Transaction
	currBlk domain.Block
	layout  *domain.Layout
}

// NewPage constructs a Page.
func NewPage(txn domain.Transaction, currblk domain.Block, layout *domain.Layout) (*Page, error) {
	if err := txn.Pin(currblk); err != nil {
		return nil, errors.Err(err, "Pin")
	}

	return &Page{
		txn:     txn,
		currBlk: currblk,
		layout:  layout,
	}, nil
}

func (page *Page) findSlotBefore(searchKey domain.Constant) (domain.SlotID, error) {
	slotID := domain.NewSlotID(0)
	for {
		lastID, err := page.getLastSlotID()
		if err != nil {
			return 0, errors.Err(err, "getLastSlotID")
		}

		val, err := page.getDataVal(slotID)
		if err != nil {
			return 0, errors.Err(err, "getDataVal")
		}

		if slotID <= lastID && val.Less(searchKey) {
			slotID++
		} else {
			break
		}
	}

	return slotID - 1, nil
}

func (page *Page) close() {
	if (page.currBlk != domain.Block{}) {
		page.txn.Unpin(page.currBlk)
	}
	page.currBlk = domain.Block{}
}

func (page *Page) isFull() (bool, error) {
	lastID, err := page.getLastSlotID()
	if err != nil {
		return false, errors.Err(err, "getLastSlotID")
	}

	return page.slotEndPos(lastID+1) >= int64(page.txn.BlockSize()), nil
}

func (page *Page) split(splitPos domain.SlotID, flag pageFlag) (domain.Block, error) {
	newBlk, err := page.appendNew(flag)
	if err != nil {
		return domain.Block{}, errors.Err(err, "appendNew")
	}

	newPage, err := NewPage(page.txn, newBlk, page.layout)
	if err != nil {
		return domain.Block{}, errors.Err(err, "NewPage")
	}

	if err := page.transferRecords(splitPos, newPage); err != nil {
		return domain.Block{}, errors.Err(err, "transferRecords")
	}

	if err := newPage.setFlag(flag); err != nil {
		return domain.Block{}, errors.Err(err, "setFlag")
	}

	newPage.close()

	return newBlk, nil
}

func (page *Page) getDataVal(slotID domain.SlotID) (domain.Constant, error) {
	return page.getVal(slotID, domain.FldDataVal)
}

func (page *Page) getFlag() (pageFlag, error) {
	flag, err := page.txn.GetInt32(page.currBlk, pageFlagOffset)
	if err != nil {
		return 0, errors.Err(err, "GetInt32")
	}

	return newPageFlag(flag), nil
}

func (page *Page) setFlag(flag pageFlag) error {
	if err := page.txn.SetInt32(page.currBlk, pageFlagOffset, flag, true); err != nil {
		return errors.Err(err, "SetInt32")
	}

	return nil
}

func (page *Page) appendNew(flag pageFlag) (domain.Block, error) {
	blk, err := page.txn.ExtendFile(page.currBlk.FileName())
	if err != nil {
		return domain.Block{}, errors.Err(err, "ExtendFile")
	}

	if err := page.txn.Pin(blk); err != nil {
		return domain.Block{}, errors.Err(err, "Pin")
	}

	if err := page.format(blk, flag); err != nil {
		return domain.Block{}, errors.Err(err, "format")
	}

	return blk, nil
}

// これ Page method に入れるべき処理なのか？
func (page *Page) format(blk domain.Block, flag pageFlag) error {
	if err := page.txn.SetInt32(blk, pageFlagOffset, flag, false); err != nil {
		return errors.Err(err, "SetInt32")
	}
	if err := page.txn.SetInt32(blk, numRecordOffset, 0, false); err != nil { // # of recs = 0
		return errors.Err(err, "SetInt32")
	}

	// TODO: 何 byte 目かという基準になっているので slot id をもとにした処理に変更したい
	recSize := page.layout.SlotSize()
	pageSize := int64(page.txn.BlockSize())
	for pos := int64(pageFlagByteLength + numRecordByteLength); pos+recSize <= pageSize; pos += recSize {
		if err := page.makeDefaultRecord(blk, pos); err != nil {
			return errors.Err(err, "makeDefaultRecord")
		}
	}

	return nil
}

// これ Page method に入れるべき処理なのか？
func (page *Page) makeDefaultRecord(blk domain.Block, pos int64) error {
	for _, fldName := range page.layout.Schema().Fields() {
		offset := page.layout.Offset(fldName)

		var err error
		switch page.layout.Schema().Type(fldName) {
		case domain.Int32FieldType:
			err = page.txn.SetInt32(blk, pos+offset, 0, false)
		case domain.StringFieldType:
			err = page.txn.SetString(blk, pos+offset, "", false)
		case domain.UnknownFieldType:
			panic(errors.New("unsupported field type"))
		default:
			panic(errors.New("unsupported field type"))
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// overflow check をここではしていない。その責務は呼び出し側に任せているっぽい
// ここでは insert するスペースを開けるだけで実際の insert はしていない
// TODO: method 名 rename する.
func (page *Page) insert(slotID domain.SlotID) error {
	lastSlotID, err := page.getLastSlotID()
	if err != nil {
		return errors.Err(err, "getLastSlotID")
	}
	for i := lastSlotID + 1; i > slotID; i-- {
		if err := page.copyRecord(i-1, i); err != nil {
			return errors.Err(err, "copyRecord")
		}
	}

	if err := page.setLastSlotID(lastSlotID + 1); err != nil {
		return errors.Err(err, "setLastSlotID")
	}

	return nil
}

func (page *Page) delete(slotID domain.SlotID) error {
	lastSlotID, err := page.getLastSlotID()
	if err != nil {
		return errors.Err(err, "delete")
	}
	for i := slotID + 1; i <= lastSlotID; i++ {
		if err := page.copyRecord(i, i-1); err != nil {
			return errors.Err(err, "copyRecord")
		}
	}

	if err := page.setLastSlotID(lastSlotID - 1); err != nil {
		return errors.Err(err, "setLastSlotID")
	}

	return nil
}

func (page *Page) copyRecord(from, to domain.SlotID) error {
	sch := page.layout.Schema()
	for _, fldName := range sch.Fields() {
		val, err := page.getVal(from, fldName)
		if err != nil {
			return errors.Err(err, "getVal")
		}
		if err := page.setVal(to, fldName, val); err != nil {
			return errors.Err(err, "setVal")
		}
	}

	return nil
}

func (page *Page) transferRecords(slotID domain.SlotID, dest *Page) error {
	destSlotID := domain.NewSlotID(0)
	var lastSlotID domain.SlotID
	var err error
	for lastSlotID, err = page.getLastSlotID(); err == nil && slotID <= lastSlotID; lastSlotID-- {
		if err := dest.insert(destSlotID); err != nil {
			return errors.Err(err, "insert")
		}
		sch := page.layout.Schema()
		for _, fldName := range sch.Fields() {
			val, err := page.getVal(slotID, fldName)
			if err != nil {
				return errors.Err(err, "getVal")
			}
			if err := dest.setVal(destSlotID, fldName, val); err != nil {
				return errors.Err(err, "setVal")
			}
		}

		// slot 番目の record を消して、それより右にあった要素を左に移動させている
		if err := page.delete(slotID); err != nil {
			return errors.Err(err, "delete")
		}
		destSlotID++
	}
	if err != nil {
		return errors.Err(err, "getLastSlotID")
	}

	return nil
}

func (page *Page) getLastSlotID() (domain.SlotID, error) {
	numRecs, err := page.txn.GetInt32(page.currBlk, numRecordOffset)
	if err != nil {
		return 0, errors.Err(err, "GetInt32")
	}

	return domain.NewSlotID(numRecs) - 1, nil
}

func (page *Page) getDataRecordID(slotID domain.SlotID) (domain.RecordID, error) {
	blkNumInt32, err := page.getInt32(slotID, domain.FldBlock)
	if err != nil {
		return domain.RecordID{}, errors.Err(err, "getInt32")
	}
	blkNum, err := domain.NewBlockNumber(blkNumInt32)
	if err != nil {
		return domain.RecordID{}, errors.Err(err, "NewBlockNumber")
	}

	recSlotInt32, err := page.getInt32(slotID, domain.FldID)
	if err != nil {
		return domain.RecordID{}, errors.Err(err, "getInt32")
	}
	recSlot := domain.NewSlotID(recSlotInt32)

	return domain.NewRecordID(blkNum, recSlot), nil
}

func (page *Page) getInt32(slotID domain.SlotID, fldName domain.FieldName) (int32, error) {
	pos := page.fldPos(slotID, fldName)

	return page.txn.GetInt32(page.currBlk, pos)
}

func (page *Page) getString(slotID domain.SlotID, fldName domain.FieldName) (string, error) {
	pos := page.fldPos(slotID, fldName)

	return page.txn.GetString(page.currBlk, pos)
}

func (page *Page) getVal(slotID domain.SlotID, fldName domain.FieldName) (domain.Constant, error) {
	typ := page.layout.Schema().Type(fldName)
	switch typ {
	case domain.Int32FieldType:
		num, err := page.getInt32(slotID, fldName)
		if err != nil {
			return domain.Constant{}, errors.Err(err, "getInt32")
		}

		return domain.NewConstant(typ, num), nil
	case domain.StringFieldType:
		str, err := page.getString(slotID, fldName)
		if err != nil {
			return domain.Constant{}, errors.Err(err, "getString")
		}

		return domain.NewConstant(typ, str), nil
	case domain.UnknownFieldType:
		return domain.Constant{}, fmt.Errorf("unsupported field type %v", typ)
	default:
		return domain.Constant{}, fmt.Errorf("unsupported field type %v", typ)
	}
}

func (page *Page) setInt32(slotID domain.SlotID, fldName domain.FieldName, val int32) error {
	pos := page.fldPos(slotID, fldName)

	return page.txn.SetInt32(page.currBlk, pos, val, true)
}

func (page *Page) setString(slotID domain.SlotID, fldName domain.FieldName, val string) error {
	pos := page.fldPos(slotID, fldName)

	return page.txn.SetString(page.currBlk, pos, val, true)
}

func (page *Page) setVal(slotID domain.SlotID, fldName domain.FieldName, val domain.Constant) error {
	typ := page.layout.Schema().Type(fldName)
	switch typ {
	case domain.Int32FieldType:
		v, err := val.AsInt32()
		if err != nil {
			return errors.Err(err, "AsInt32")
		}

		return page.setInt32(slotID, fldName, v)
	case domain.StringFieldType:
		v, err := val.AsString()
		if err != nil {
			return errors.Err(err, "AsString")
		}

		return page.setString(slotID, fldName, v)
	case domain.UnknownFieldType:
		return fmt.Errorf("unsupported field type %v", typ)
	default:
		return fmt.Errorf("unsupported field type %v", typ)
	}
}

func (page *Page) setLastSlotID(n domain.SlotID) error {
	return page.txn.SetInt32(page.currBlk, numRecordOffset, int32(n+1), true)
}

func (page *Page) fldPos(slotID domain.SlotID, fldName domain.FieldName) int64 {
	offset := page.layout.Offset(fldName)

	return page.slotPos(slotID) + offset
}

func (page *Page) slotPos(slotID domain.SlotID) int64 {
	slotSize := page.layout.SlotSize()

	return pageFlagByteLength + numRecordByteLength + (int64(slotID) * slotSize)
}

func (page *Page) slotEndPos(slotID domain.SlotID) int64 {
	slotSize := page.layout.SlotSize()

	return pageFlagByteLength + numRecordByteLength + (int64(slotID+1) * slotSize)
}
