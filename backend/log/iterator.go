package log

import (
	"github.com/goropikari/simpledbgo/backend/domain"
	"github.com/goropikari/simpledbgo/meta"
)

// Iterator is iterator of log.
type Iterator struct {
	fileMgr    domain.FileManager
	block      *domain.Block
	page       *domain.Page
	currentPos int32
	// boundary   int32
}

// NewIterator is a constructor of Iterator.
func NewIterator(fileMgr domain.FileManager, block *domain.Block, page *domain.Page) (*Iterator, error) {
	err := fileMgr.CopyBlockToPage(block, page)
	if err != nil {
		return nil, err
	}

	currentPos, err := page.GetInt32(0)
	if err != nil {
		return nil, err
	}

	return &Iterator{
		fileMgr:    fileMgr,
		block:      block,
		page:       page,
		currentPos: currentPos,
		// boundary:   0,
	}, nil
}

// HasNext checks whether iterator has next items or not.
func (iter *Iterator) HasNext() bool {
	return iter.currentPos < int32(iter.block.Size()) || int32(iter.block.Number()) > 0
}

// Next returns a next item.
func (iter *Iterator) Next() ([]byte, error) {
	if iter.currentPos == int32(iter.block.Size()) {
		blk := domain.NewBlock(iter.block.FileName(), iter.block.Size(), iter.block.Number()-1)
		err := iter.moveToBlock(blk)
		if err != nil {
			return nil, err
		}
		iter.block = blk
	}

	record, err := iter.page.GetBytes(int64(iter.currentPos))
	if err != nil {
		return nil, err
	}

	iter.currentPos += meta.Int32Length + int32(len(record))

	return record, nil
}

func (iter *Iterator) moveToBlock(block *domain.Block) error {
	err := iter.fileMgr.CopyBlockToPage(block, iter.page)
	if err != nil {
		return err
	}

	boundary, err := iter.page.GetInt32(0)
	if err != nil {
		return err
	}

	iter.currentPos = boundary

	return nil
}
