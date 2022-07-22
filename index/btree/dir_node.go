package btree

import (
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/errors"
)

// DirEntry is entry of directory node.
type DirEntry struct {
	dataval domain.Constant
	blkNum  domain.BlockNumber
}

// NewDirEntry constructs a DirEntry.
func NewDirEntry(key domain.Constant, blkNum domain.BlockNumber) DirEntry {
	return DirEntry{
		dataval: key,
		blkNum:  blkNum,
	}
}

func (de DirEntry) getDataVal() domain.Constant {
	return de.dataval
}

func (de DirEntry) getBlockNumber() domain.BlockNumber {
	return de.blkNum
}

// DirNodeLevel represents directory node level.
type DirNodeLevel = int32

const (
	zeroLevelDirNodeFlag DirNodeLevel = 0
)

// DirPage is page for directory node.
type DirPage struct {
	*Page
}

// NewDirPage constructs a directory page.
func NewDirPage(txn domain.Transaction, blk domain.Block, layout *domain.Layout) (*DirPage, error) {
	page, err := NewPage(txn, blk, layout)
	if err != nil {
		return nil, errors.Err(err, "NewPage")
	}

	return &DirPage{Page: page}, nil
}

func (p *DirPage) close() {
	p.Page.close()
}

func (p *DirPage) getLevel() (DirNodeLevel, error) {
	return p.getFlag()
}

func (p *DirPage) setLevel(l DirNodeLevel) error {
	return p.setFlag(l)
}

func (p *DirPage) insertDir(slotID domain.SlotID, val domain.Constant, blkNum domain.BlockNumber) error {
	if err := p.insert(slotID); err != nil {
		return errors.Err(err, "insert")
	}
	if err := p.setVal(slotID, domain.FldDataVal, val); err != nil {
		return errors.Err(err, "setVal")
	}
	if err := p.setInt32(slotID, domain.FldBlock, int32(blkNum)); err != nil {
		return errors.Err(err, "setInt32")
	}

	return nil
}

func (p *DirPage) format(blk domain.Block, level DirNodeLevel) error {
	return p.Page.format(blk, level)
}

func (p *DirPage) getChildBlockNumber(slotID domain.SlotID) (domain.BlockNumber, error) {
	num, err := p.getInt32(slotID, domain.FldBlock)
	if err != nil {
		return 0, errors.Err(err, "getInt32")
	}

	return domain.NewBlockNumber(num)
}

// DirNode is directory node.
type DirNode struct {
	txn      domain.Transaction
	layout   *domain.Layout
	contents *DirPage
	fileName domain.FileName
}

// NewDirNode constructs a directory node.
func NewDirNode(txn domain.Transaction, blk domain.Block, dirLayout *domain.Layout) (*DirNode, error) {
	page, err := NewDirPage(txn, blk, dirLayout)
	if err != nil {
		return nil, errors.Err(err, "NewDirPage")
	}

	return &DirNode{
		txn:      txn,
		layout:   dirLayout,
		contents: page,
		fileName: blk.FileName(),
	}, nil
}

func (node *DirNode) close() {
	node.contents.close()
}

// search searches the leaf node block number which contains given key.
// During searching, update contents that is parent of leaf node.
// TODO: rename method.
func (node *DirNode) search(key domain.Constant) (domain.BlockNumber, error) {
	childBlk, err := node.findChildBlock(key)
	if err != nil {
		return 0, errors.Err(err, "findChildBlock")
	}

	for {
		level, err := node.contents.getLevel()
		if err != nil {
			return 0, errors.Err(err, "getLevel")
		}
		if level <= zeroLevelDirNodeFlag {
			break
		}

		node.close()
		node.contents, err = NewDirPage(node.txn, childBlk, node.layout)
		if err != nil {
			return 0, errors.Err(err, "NewDirPage")
		}
		childBlk, err = node.findChildBlock(key)
		if err != nil {
			return 0, errors.Err(err, "findChildBlock")
		}
	}

	return childBlk.Number(), nil
}

func (node *DirNode) makeNewRoot(e DirEntry) error {
	firstVal, err := node.contents.getDataVal(0)
	if err != nil {
		return errors.Err(err, "getFirstVal")
	}

	level, err := node.contents.getLevel()
	if err != nil {
		return errors.Err(err, "getLevel")
	}

	newBlk, err := node.contents.split(0, level)
	if err != nil {
		return errors.Err(err, "split")
	}
	oldRoot := NewDirEntry(firstVal, newBlk.Number())
	_, _, err = node.insertEntry(oldRoot)
	if err != nil {
		return errors.Err(err, "insertEntry")
	}
	_, _, err = node.insertEntry(e)
	if err != nil {
		return errors.Err(err, "insertEntry")
	}
	if err := node.contents.setLevel(level + 1); err != nil {
		return errors.Err(err, "setLevel")
	}

	return nil
}

func (node *DirNode) insert(e DirEntry) (newEnt DirEntry, makeNew bool, err error) {
	level, err := node.contents.getLevel()
	if err != nil {
		return DirEntry{}, false, errors.Err(err, "getLevel")
	}
	if level == 0 {
		return node.insertEntry(e)
	}

	childBlk, err := node.findChildBlock(e.getDataVal())
	if err != nil {
		return DirEntry{}, false, errors.Err(err, "findChildBlock")
	}
	child, err := NewDirNode(node.txn, childBlk, node.layout)
	if err != nil {
		return DirEntry{}, false, errors.Err(err, "NewDirNode")
	}
	myentry, makeNew, err := child.insert(e)
	if err != nil {
		return DirEntry{}, false, errors.Err(err, "insert")
	}
	if makeNew {
		return node.insertEntry(myentry)
	}

	return DirEntry{}, false, nil
}

func (node *DirNode) insertEntry(e DirEntry) (newDir DirEntry, makeNew bool, err error) {
	newSlot, err := node.contents.findSlotBefore(e.getDataVal())
	if err != nil {
		return DirEntry{}, false, errors.Err(err, "findSlotBefore")
	}
	newSlot++
	if err = node.contents.insertDir(newSlot, e.getDataVal(), e.getBlockNumber()); err != nil {
		return DirEntry{}, false, errors.Err(err, "insertDir")
	}

	full, err := node.contents.isFull()
	if err != nil {
		return DirEntry{}, false, errors.Err(err, "isFull")
	}
	if !full {
		return DirEntry{}, false, nil
	}

	splitPos, err := node.contents.getLastSlotID()
	if err != nil {
		return DirEntry{}, false, errors.Err(err, "getLastSlotID")
	}
	splitPos /= 2
	splitVal, err := node.contents.getDataVal(splitPos)
	if err != nil {
		return DirEntry{}, false, errors.Err(err, "getDataVal")
	}
	level, err := node.contents.getLevel()
	if err != nil {
		return DirEntry{}, false, errors.Err(err, "getLevel")
	}
	newBlk, err := node.contents.split(splitPos, level)
	if err != nil {
		return DirEntry{}, false, errors.Err(err, "split")
	}

	return NewDirEntry(splitVal, newBlk.Number()), true, nil
}

func (node *DirNode) findChildBlock(searchKey domain.Constant) (domain.Block, error) {
	slotID, err := node.contents.findSlotBefore(searchKey)
	if err != nil {
		return domain.Block{}, errors.Err(err, "findSlotBefore")
	}

	nextVal, err := node.contents.getDataVal(slotID + 1)
	if err != nil {
		return domain.Block{}, errors.Err(err, "getDataVal")
	}
	if searchKey.Equal(nextVal) {
		slotID++
	}

	blkNum, err := node.contents.getChildBlockNumber(slotID)
	if err != nil {
		return domain.Block{}, errors.Err(err, "getChildBlockNumber")
	}

	return domain.NewBlock(node.fileName, blkNum), nil
}
