package btree

import (
	"fmt"
	"math"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/errors"
)

// Index is btree index.
type Index struct {
	txn         domain.Transaction
	dirLayout   *domain.Layout
	leafLayout  *domain.Layout // domain.createIdxLayout で作られた Layout が入ってくる。
	leafTblName domain.FileName
	leaf        *LeafNode
	rootBlk     domain.Block
	err         error
}

// NewIndex constructs a index.
func NewIndex(txn domain.Transaction, idxName domain.IndexName, leafLayout *domain.Layout) (*Index, error) {
	// leaf node
	leafTblName, err := domain.NewFileName(idxName.String() + "leaf")
	if err != nil {
		return nil, errors.Err(err, "NewFileName")
	}

	leafTblSize, err := txn.BlockLength(leafTblName)
	if err != nil {
		return nil, errors.Err(err, "BlockLength")
	}
	if leafTblSize == 0 {
		blk, err := txn.ExtendFile(leafTblName)
		if err != nil {
			return nil, errors.Err(err, "ExtendFile")
		}
		leafNode, err := NewPage(txn, blk, leafLayout)
		if err != nil {
			return nil, errors.Err(err, "NewPage")
		}
		if err := leafNode.format(blk, noOverflowLeafNode); err != nil {
			return nil, errors.Err(err, "format")
		}
	}

	// dir node
	dirSch := domain.NewSchema()
	dirSch.Add(domain.FldBlock, leafLayout.Schema())
	dirSch.Add(domain.FldDataVal, leafLayout.Schema())
	dirTblName, err := domain.NewFileName(idxName.String() + "dir")
	if err != nil {
		return nil, errors.Err(err, "NewFileName")
	}
	dirLayout := domain.NewLayout(dirSch)
	rootBlk := domain.NewBlock(dirTblName, 0)

	dirTblSize, err := txn.BlockLength(dirTblName)
	if err != nil {
		return nil, errors.Err(err, "BlockLength")
	}
	if dirTblSize == 0 {
		// create new root block
		blk, err := txn.ExtendFile(dirTblName)
		if err != nil {
			return nil, errors.Err(err, "ExtendFile")
		}
		dirPage, err := NewDirPage(txn, blk, dirLayout)
		if err != nil {
			return nil, errors.Err(err, "NewDirNode")
		}
		if err := dirPage.format(blk, zeroLevelDirNodeFlag); err != nil {
			return nil, errors.Err(err, "format")
		}

		// insert initial directory entry (sentinel)
		fldType := dirSch.Type(domain.FldDataVal)
		var minVal domain.Constant
		switch fldType {
		case domain.Int32FieldType:
			minVal = domain.NewConstant(fldType, int32(math.MinInt32))
		case domain.StringFieldType:
			minVal = domain.NewConstant(fldType, "")
		case domain.UnknownFieldType:
			panic(fmt.Errorf("not supported FieldType %v", fldType))
		default:
			panic(fmt.Errorf("not supported FieldType %v", fldType))
		}
		initSlotID := domain.NewSlotID(0)
		initBlkNum, err := domain.NewBlockNumber(0)
		if err != nil {
			return nil, errors.Err(err, "NewBlockNumber")
		}
		if err := dirPage.insertDir(initSlotID, minVal, initBlkNum); err != nil {
			return nil, err
		}
		dirPage.close()
	}

	return &Index{
		txn:         txn,
		dirLayout:   dirLayout,
		leafLayout:  leafLayout,
		leafTblName: leafTblName,
		leaf:        nil,
		rootBlk:     rootBlk,
	}, nil
}

// BeforeFirst moves current slot where that is before the lowerbound of given search key.
// search key を lowerbound とする index の1つ手前の index を返す。
func (idx *Index) BeforeFirst(searchKey domain.Constant) error {
	idx.Close()
	root, err := NewDirNode(idx.txn, idx.rootBlk, idx.dirLayout)
	if err != nil {
		return errors.Err(err, "NewDirNode")
	}

	// root が level 0 以外だと search している間に root.contests は root node とは全く別物が入っている
	// root.search 後は root.contents に flag が 0 の dir が入っている
	blkNum, err := root.search(searchKey)
	if err != nil {
		return errors.Err(err, "search")
	}
	root.close()
	leafBlk := domain.NewBlock(idx.leafTblName, blkNum)

	leaf, err := NewLeafNode(idx.txn, leafBlk, idx.leafLayout, searchKey)
	if err != nil {
		return errors.Err(err, "NewLeafNode")
	}
	idx.leaf = leaf

	return nil
}

// HasNext checks the existence of given search key.
func (idx *Index) HasNext() bool {
	found := idx.leaf.hasNext()
	if idx.leaf.Err() != nil {
		idx.err = idx.leaf.Err()
	}

	return found
}

// GetDataRecordID gets record id from current slot.
func (idx *Index) GetDataRecordID() (domain.RecordID, error) {
	return idx.leaf.getDataRecordID()
}

// Insert inserts an index record.
func (idx *Index) Insert(dataval domain.Constant, datarid domain.RecordID) error {
	if err := idx.BeforeFirst(dataval); err != nil { // leaf を dataval を持つ leaf node に更新
		return errors.Err(err, "beforeFirst")
	}

	e, splitted, err := idx.leaf.insert(datarid)
	if err != nil {
		return errors.Err(err, "insert")
	}
	idx.Close()
	if !splitted { // split が起こらなかったとき
		return nil
	}

	root, err := NewDirNode(idx.txn, idx.rootBlk, idx.dirLayout)
	if err != nil {
		return errors.Err(err, "NewDirNode")
	}

	e2, splitted, err := root.insert(e)
	if err != nil {
		return errors.Err(err, "insert")
	}
	if splitted { // directory node で split が起こったとき
		if err := root.makeNewRoot(e2); err != nil {
			return errors.Err(err, "makeNewRoot")
		}
	}
	root.close()

	return nil
}

// Delete deletes given index record.
func (idx *Index) Delete(dataval domain.Constant, datarid domain.RecordID) error {
	if err := idx.BeforeFirst(dataval); err != nil {
		return errors.Err(err, "BeforeFirst")
	}
	if err := idx.leaf.delete(datarid); err != nil {
		return errors.Err(err, "delete")
	}
	idx.Close()

	return nil
}

// Close closes the index.
func (idx *Index) Close() {
	if idx.leaf != nil {
		idx.leaf.close()
	}
}

// Err returns the stored error.
func (idx *Index) Err() error {
	return idx.err
}
