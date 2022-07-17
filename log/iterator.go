package log

import (
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/errors"
)

// Iterator is iterator of log.
type Iterator struct {
	fileMgr    domain.FileManager
	block      domain.Block
	page       *Page
	currentPos int32
}

// NewIterator is a constructor of Iterator.
func NewIterator(fileMgr domain.FileManager, block domain.Block, page *Page) (*Iterator, error) {
	err := fileMgr.CopyBlockToPage(block, page.getDomainPage())
	if err != nil {
		return nil, errors.Err(err, "CopyBlockToPage")
	}

	currentPos, err := page.getBoundaryOffset()
	if err != nil {
		return nil, errors.Err(err, "getBoundaryOffset")
	}

	return &Iterator{
		fileMgr:    fileMgr,
		block:      block,
		page:       page,
		currentPos: currentPos,
	}, nil
}

// HasNext checks whether iterator has next items or not.
func (iter *Iterator) HasNext() bool {
	return iter.currentPos < int32(iter.page.Size()) || iter.block.Number() > 0
}

// Next returns a next item.
func (iter *Iterator) Next() ([]byte, error) {
	if iter.currentPos == int32(iter.fileMgr.BlockSize()) {
		blk := domain.NewBlock(iter.block.FileName(), iter.block.Number()-1)
		err := iter.moveToBlock(blk)
		if err != nil {
			return nil, errors.Err(err, "moveToBlock")
		}
		iter.block = blk
	}

	record, err := iter.page.getRecord(iter.currentPos)
	if err != nil {
		return nil, errors.Err(err, "getRecord")
	}

	iter.currentPos += int32(iter.page.neededByteLength(record))

	return record, nil
}

func (iter *Iterator) moveToBlock(block domain.Block) error {
	err := iter.fileMgr.CopyBlockToPage(block, iter.page.getDomainPage())
	if err != nil {
		return errors.Err(err, "CopyBlockToPage")
	}

	boundary, err := iter.page.getBoundaryOffset()
	if err != nil {
		return errors.Err(err, "getBoundaryOffset")
	}

	iter.currentPos = boundary

	return nil
}

// Err returns iterator's error.
func (iter *Iterator) Err() error {
	return nil
}
