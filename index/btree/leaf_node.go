package btree

import (
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/errors"
)

// LeafPage is page for leaf node.
type LeafPage struct {
	*Page
}

// NewLeafPage constructs a LeafPage.
func NewLeafPage(txn domain.Transaction, blk domain.Block, layout *domain.Layout) (*LeafPage, error) {
	page, err := NewPage(txn, blk, layout)
	if err != nil {
		return nil, errors.Err(err, "NewPage")
	}

	return &LeafPage{Page: page}, nil
}

func (page *LeafPage) getFirstDataVal() (domain.Constant, error) {
	return page.getDataVal(domain.NewSlotID(0))
}

func (page *LeafPage) getOverflowBlockNumber() (domain.BlockNumber, error) {
	nextBlkNumInt32, err := page.getFlag()
	if err != nil {
		return 0, errors.Err(err, "getFlag")
	}

	return domain.NewBlockNumber(nextBlkNumInt32)
}

func (page *LeafPage) setOverflowBlockNumber(blkNum domain.BlockNumber) error {
	return page.setFlag(pageFlag(blkNum))
}

func (page *LeafPage) insertLeaf(slotID domain.SlotID, val domain.Constant, rid domain.RecordID) error {
	if err := page.insert(slotID); err != nil {
		return errors.Err(err, "insert")
	}

	if err := page.setVal(slotID, domain.FldDataVal, val); err != nil {
		return errors.Err(err, "setVal")
	}
	if err := page.setInt32(slotID, domain.FldBlock, rid.BlockNumber().ToInt32()); err != nil {
		return errors.Err(err, "setInt32")
	}
	if err := page.setInt32(slotID, domain.FldID, rid.SlotID().ToInt32()); err != nil {
		return errors.Err(err, "setInt32")
	}

	return nil
}

func (page *LeafPage) split(splitPos domain.SlotID, blkNum domain.BlockNumber) (domain.Block, error) {
	return page.Page.split(splitPos, pageFlag(blkNum))
}

const (
	noOverflowLeafNode = -1
)

// LeafNode is btree leaf node.
type LeafNode struct {
	txn         domain.Transaction
	layout      *domain.Layout
	searchKey   domain.Constant
	contents    *LeafPage
	currentSlot domain.SlotID
	fileName    domain.FileName
	err         error
}

// NewLeafNode constructs a LeafNode.
func NewLeafNode(txn domain.Transaction, blk domain.Block, layout *domain.Layout, searchKey domain.Constant) (*LeafNode, error) {
	contents, err := NewLeafPage(txn, blk, layout)
	if err != nil {
		return nil, errors.Err(err, "NewLeafPage")
	}

	currentSlot, err := contents.findSlotBefore(searchKey)
	if err != nil {
		return nil, errors.Err(err, "findSlotBefore")
	}

	return &LeafNode{
		txn:         txn,
		layout:      layout,
		searchKey:   searchKey,
		contents:    contents,
		currentSlot: currentSlot,
		fileName:    blk.FileName(),
		err:         nil,
	}, nil
}

func (node *LeafNode) close() {
	node.contents.close()
}

func (node *LeafNode) getDataRecordID() (domain.RecordID, error) {
	return node.contents.getDataRecordID(node.currentSlot)
}

func (node *LeafNode) hasNext() bool {
	node.currentSlot++
	lastID, err := node.contents.getLastSlotID()
	if err != nil {
		node.err = errors.Err(err, "getLastSlotID")

		return false
	}
	if node.currentSlot > lastID {
		return node.tryOverflow()
	}

	val, err := node.contents.getDataVal(node.currentSlot)
	if err != nil {
		node.err = errors.Err(err, "getDataVal")

		return false
	}

	if val.Equal(node.searchKey) {
		return true
	}

	return node.tryOverflow()
}

// Err returns the stored error.
func (node *LeafNode) Err() error {
	return node.err
}

func (node *LeafNode) delete(datarid domain.RecordID) error {
	for node.hasNext() {
		rid, err := node.getDataRecordID()
		if err != nil {
			return errors.Err(err, "getDataRecordID")
		}
		if rid == datarid {
			if err := node.contents.delete(node.currentSlot); err != nil {
				return errors.Err(err, "delete")
			}

			return nil
		}
	}
	if node.Err() != nil {
		return errors.Err(node.Err(), "hasNext")
	}

	return nil
}

func (node *LeafNode) insert(datarid domain.RecordID) (ent DirEntry, splitted bool, err error) {
	// BUG: record 数 0 でも slotID = 0 の値を取ってきてしまっている。 record 数が 0 でなかったらという条件を加えるべき
	firstVal, err := node.getFirstDataVal()
	if err != nil {
		return DirEntry{}, false, errors.Err(err, "getFirstDataVal")
	}

	isOver, err := node.isOverflow()
	if err != nil {
		return DirEntry{}, false, errors.Err(err, "isOverflow")
	}

	// isOver: overflow block が存在する
	// node.searchKey.Less(firstVal): 挿入する record が slot id = 0 に入っている値よりも小さいとき
	// overflow している場合は slot_id = 0 の値を overflow しているとしたいからここだけ特別に処理している
	if isOver && node.searchKey.Less(firstVal) {
		initSlotID := domain.NewSlotID(0)
		nextBlkNum, err := node.contents.getOverflowBlockNumber()
		if err != nil {
			return DirEntry{}, false, errors.Err(err, "getOverflowBlockNumber")
		}

		newBlk, err := node.contents.split(initSlotID, nextBlkNum)
		if err != nil {
			return DirEntry{}, false, errors.Err(err, "split")
		}

		node.currentSlot = 0
		if err := node.contents.setOverflowBlockNumber(noOverflowLeafNode); err != nil {
			return DirEntry{}, false, errors.Err(err, "setOverflowBlockNumber")
		}
		if err := node.contents.insertLeaf(node.currentSlot, node.searchKey, datarid); err != nil {
			return DirEntry{}, false, errors.Err(err, "insertLeaf")
		}

		return NewDirEntry(firstVal, newBlk.Number()), true, nil
	}

	// insert method が呼ばれる前には findSlotBefore が呼ばれる前提だからこの時点での currentslot には searchkey <= dataval(x) となる最小の x の **一つ手前** の slot id が入っている。
	// そのような x のことを x' とすると x'-1 が入っている
	// ここで dataval(x) は slot_id = x の dataval のこと
	// 気持ち的には contents.lower_bound(searchKey) - 1
	node.currentSlot++
	if err := node.contents.insertLeaf(node.currentSlot, node.searchKey, datarid); err != nil {
		return DirEntry{}, false, errors.Err(err, "insertLeaf")
	}

	// isFull は insertLeaf の前にやらなくてよいのか？常に一度は insert できる程度には space を確保しておけということか？
	full, err := node.contents.isFull()
	if err != nil {
		return DirEntry{}, false, errors.Err(err, "isFull")
	}
	if !full {
		return DirEntry{}, false, nil
	}

	// page is full, so split it
	firstSlotID := domain.NewSlotID(0)
	lastSlotID, err := node.contents.getLastSlotID()
	if err != nil {
		return DirEntry{}, false, errors.Err(err, "getLastSlotID")
	}
	firstKey, err := node.contents.getDataVal(firstSlotID)
	if err != nil {
		return DirEntry{}, false, errors.Err(err, "getDataVal")
	}
	lastKey, err := node.contents.getDataVal(lastSlotID)
	if err != nil {
		return DirEntry{}, false, errors.Err(err, "getDataVal")
	}

	// make overflow block
	if firstKey.Equal(lastKey) {
		nextBlkNum, err := node.contents.getOverflowBlockNumber()
		if err != nil {
			return DirEntry{}, false, errors.Err(err, "getOverflowBlockNumber")
		}

		// すでに overflow していることがあるから flag が -1 ではないこともある
		newBlk, err := node.contents.split(1, nextBlkNum)
		if err != nil {
			return DirEntry{}, false, errors.Err(err, "split")
		}

		// 新しく作った overflow block の block number をもとの block の flag として入れている
		if err := node.contents.setOverflowBlockNumber(newBlk.Number()); err != nil {
			return DirEntry{}, false, errors.Err(err, "setOverflowBlockNumber")
		}

		// overflow block のときは directory node は増えない
		return DirEntry{}, false, nil
	}

	// overflow でなかったときは半分で分ける
	splitPos := (lastSlotID + 1) / 2
	splitKey, err := node.contents.getDataVal(splitPos)
	if err != nil {
		return DirEntry{}, false, errors.Err(err, "getDataVal")
	}

	if splitKey.Equal(firstKey) {
		// firstKey, lastKey が同じケースは then 節でチェックしてあるから、この部分では firstKey != lastKey が保証されている.
		// そのため、splitPos が範囲外になることはない
		for {
			currVal, err := node.contents.getDataVal(splitPos)
			if err != nil {
				return DirEntry{}, false, errors.Err(err, "getDataVal")
			}
			if currVal.Equal(splitKey) {
				splitPos++
			} else {
				break
			}
		}

		splitKey, err = node.contents.getDataVal(splitPos)
		if err != nil {
			return DirEntry{}, false, errors.Err(err, "getDataVal")
		}
	} else {
		// firstKey != splitKey だから splitpos が範囲外(-1) にいくことはないことは保証されている
		// 同じ値の record は同じ node に収めるために searchKey と同じ値である最小の slot id を探索する
		for {
			prevVal, err := node.contents.getDataVal(splitPos - 1)
			if err != nil {
				return DirEntry{}, false, errors.Err(err, "getDataVal")
			}
			if prevVal.Equal(splitKey) {
				splitPos--
			} else {
				break
			}
		}
	}

	newBlk, err := node.contents.split(splitPos, noOverflowLeafNode)
	if err != nil {
		return DirEntry{}, false, err
	}

	return NewDirEntry(splitKey, newBlk.Number()), true, nil
}

func (node *LeafNode) isOverflow() (bool, error) {
	nextBlkNum, err := node.contents.getFlag()
	if err != nil {
		return false, errors.Err(err, "getFlag")
	}

	return nextBlkNum != noOverflowLeafNode, nil
}

func (node *LeafNode) getFirstDataVal() (domain.Constant, error) {
	val, err := node.contents.getFirstDataVal()
	if err != nil {
		return domain.Constant{}, errors.Err(err, "getFirstDataVal")
	}

	return val, nil
}

func (node *LeafNode) tryOverflow() bool {
	firstKey, err := node.contents.getDataVal(domain.NewSlotID(0))
	if err != nil {
		node.err = errors.Err(err, "getDataVal")

		return false
	}

	// fix original implementation bug
	// check the existence of record in overflow block
	for {
		nextBlkNum, err := node.contents.getOverflowBlockNumber()
		if err != nil {
			node.err = errors.Err(err, "getOverflowBlockNumber")

			return false
		}
		if !node.searchKey.Equal(firstKey) || nextBlkNum == noOverflowLeafNode {
			return false
		}

		node.contents.close()
		nextBlk := domain.NewBlock(node.fileName, nextBlkNum)
		contents, err := NewLeafPage(node.txn, nextBlk, node.layout)
		if err != nil {
			node.err = errors.Err(err, "NewLeafPage")

			return false
		}
		node.contents = contents
		node.currentSlot = 0

		// block の結合処理をしていないから overflow block の繋がりは途中に空の overflow block が
		// 入っていることがある. そのため、空だから false を返却するというわけには行かない。
		// この処理までくると overflow block を探索していることは保証させるので, lastSlotID が
		// 0 以上になった場合はレコードが存在することが保証される。
		// ここで、slotID は 0 から始まり、record が存在しなかった場合 getLastSlotID は -1 を
		// 返却する。
		lastSlotID, err := node.contents.getLastSlotID()
		if err != nil {
			node.err = err

			return false
		}
		if lastSlotID >= 0 {
			return true
		}
	}
}
