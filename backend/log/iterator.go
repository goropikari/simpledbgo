package log

import (
	"errors"
	"fmt"
	"io"

	"github.com/goropikari/simpledb_go/backend/core"
	"github.com/goropikari/simpledb_go/backend/service"
)

// Iterator is iterator of log manager.
type Iterator struct {
	fileMgr               service.FileManager
	block                 *core.Block
	page                  *core.Page
	currentRecordPosition uint32
	boundary              uint32
}

// NewIterator is an iterator of log.
func NewIterator(fileMgr service.FileManager, block *core.Block) (*Iterator, error) {
	page, err := fileMgr.PreparePage()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err := fileMgr.CopyBlockToPage(block, page); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	boundary, err := page.GetUint32(0)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("%w", err)
	}

	return &Iterator{
		fileMgr:               fileMgr,
		block:                 block,
		page:                  page,
		currentRecordPosition: boundary,
		boundary:              boundary,
	}, nil
}

// HasNext checks whether there is another next item.
func (it *Iterator) HasNext() bool {
	blockSize := it.fileMgr.GetBlockSize()

	return int(it.currentRecordPosition) < blockSize || it.block.GetBlockNumber() > 0
}

// Next returns next item.
func (it *Iterator) Next() ([]byte, error) {
	blockSize := it.fileMgr.GetBlockSize()

	if it.currentRecordPosition == uint32(blockSize) {
		block := core.NewBlock(it.block.GetFileName(), it.block.GetBlockNumber()-1)
		err := it.moveToBlock(block)
		if err != nil {
			return nil, err
		}
	}

	record, err := it.page.GetBytes(int64(it.currentRecordPosition))
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	it.currentRecordPosition += uint32(core.Uint32Length + len(record))

	return record, nil
}

func (it *Iterator) moveToBlock(block *core.Block) error {
	err := it.fileMgr.CopyBlockToPage(block, it.page)
	if err != nil {
		return err
	}

	boundary, err := it.page.GetUint32(0)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	it.currentRecordPosition = boundary

	it.block = block

	return nil
}
