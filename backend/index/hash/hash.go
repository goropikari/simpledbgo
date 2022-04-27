package hash

import (
	"fmt"

	"github.com/goropikari/simpledbgo/backend/domain"
	"github.com/goropikari/simpledbgo/meta"
)

const (
	numBuckets = 100
	fldDataVal = "dataval"
	fldBlock   = "block"
	fldID      = "id"
)

// Index is a model of index.
type Index struct {
	txn       domain.Transaction
	idxName   domain.IndexName
	layout    *domain.Layout
	searchKey meta.Constant
	tbl       *domain.Table
}

// NewIndex constructs an Index.
func NewIndex(txn domain.Transaction, idxName domain.IndexName, layout *domain.Layout) *Index {
	return &Index{
		txn:     txn,
		idxName: idxName,
		layout:  layout,
	}
}

// BeforeFirst ...
// searchKey を持っている table をセットする.
func (idx *Index) BeforeFirst(searchKey meta.Constant) error {
	idx.Close()
	idx.searchKey = searchKey
	bucket := searchKey.HashCode() % numBuckets
	// hash 値なので searchKey の値が違っても同じ index file に書き込まれることはある
	tblName, err := domain.NewTableName(fmt.Sprintf("%v%v", idx.idxName, bucket))
	if err != nil {
		return err
	}

	tbl, err := domain.NewTable(idx.txn, tblName, idx.layout)
	if err != nil {
		return err
	}
	idx.tbl = tbl

	return nil
}

// HasNext checks whether tbl has a record having the searchKey.
func (idx *Index) HasNext() (bool, error) {
	for {
		found, err := idx.tbl.HasNextUsedSlot()
		if err != nil {
			return false, err
		}
		if !found {
			break
		}

		v, err := idx.tbl.GetVal(fldDataVal)
		if err != nil {
			return false, err
		}

		if v.Equal(idx.searchKey) {
			return true, nil
		}
	}

	return false, nil
}

// GetDataRecordID gets record id associated with searchKey from index file.
func (idx *Index) GetDataRecordID() (domain.RecordID, error) {
	num, err := idx.tbl.GetInt32(fldBlock)
	if err != nil {
		return domain.NewZeroRecordID(), err
	}

	blkNum, err := domain.NewBlockNumber(num)
	if err != nil {
		return domain.NewZeroRecordID(), err
	}

	id, err := idx.tbl.GetInt32(fldID)
	if err != nil {
		return domain.NewZeroRecordID(), err
	}

	slotID := domain.NewSlotID(id)

	return domain.NewRecordID(blkNum, slotID), nil
}

// Insert inserts search key with record id into the index file.
func (idx *Index) Insert(searchKey meta.Constant, rid domain.RecordID) error {
	if err := idx.BeforeFirst(searchKey); err != nil {
		return err
	}

	if err := idx.tbl.SetInt32(fldBlock, int32(rid.BlockNumber())); err != nil {
		return err
	}
	if err := idx.tbl.SetInt32(fldID, int32(rid.SlotID())); err != nil {
		return err
	}
	if err := idx.tbl.SetVal(fldDataVal, searchKey); err != nil {
		return err
	}

	return nil
}

// Delete given record from index file.
func (idx *Index) Delete(searchKey meta.Constant, rid domain.RecordID) error {
	if err := idx.BeforeFirst(searchKey); err != nil {
		return err
	}

	for {
		found, err := idx.tbl.HasNextUsedSlot()
		if err != nil {
			return err
		}
		if !found {
			break
		}

		getRID, err := idx.GetDataRecordID()
		if err != nil {
			return err
		}

		if rid.Equal(getRID) {
			if err := idx.tbl.Delete(); err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}

// Close closes table.
func (idx *Index) Close() {
	if idx.tbl == nil {
		return
	}

	idx.tbl.Close()
}

// SearchCost calculates search cost.
// rpb: record per block
func (idx *Index) SearchCost(numBlks int, rpb int) int {
	return numBlks / numBuckets
}
