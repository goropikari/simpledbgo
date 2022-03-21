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

func iterator(fileMgr service.FileManager, block *core.Block) (<-chan []byte, error) {
	if err := validateArgs(fileMgr, block); err != nil {
		return nil, err
	}

	ch := make(chan []byte)

	it, err := newIterator(fileMgr, block)
	if err != nil {
		return nil, err
	}

	go func() {
		for it.hasNext() {
			ch <- it.next()
		}
		close(ch)
	}()

	return ch, nil
}

func newIterator(fileMgr service.FileManager, block *core.Block) (*Iterator, error) {
	if err := validateArgs(fileMgr, block); err != nil {
		return nil, err
	}

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

func (logIt *Iterator) hasNext() bool {
	blockSize := logIt.fileMgr.GetBlockSize()

	return int(logIt.currentRecordPosition) < blockSize || logIt.block.GetBlockNumber() > 0
}

func (logIt *Iterator) next() []byte {
	blockSize := logIt.fileMgr.GetBlockSize()

	if logIt.currentRecordPosition == uint32(blockSize) {
		block := core.NewBlock(logIt.block.GetFileName(), logIt.block.GetBlockNumber()-1)
		logIt.moveToBlock(block)
	}

	record, _ := logIt.page.GetBytes(int64(logIt.currentRecordPosition))
	logIt.currentRecordPosition += uint32(core.Uint32Length + len(record))

	return record
}

func (logIt *Iterator) moveToBlock(block *core.Block) {
	err := logIt.fileMgr.CopyBlockToPage(block, logIt.page)
	if err != nil {
		panic(err)
	}
	boundary, _ := logIt.page.GetUint32(0)
	logIt.currentRecordPosition = boundary
	logIt.block = block
}

func validateArgs(fileMgr service.FileManager, block *core.Block) error {
	if fileMgr.IsZero() {
		return errors.New("fileMgr must not be nil")
	}
	if block == nil {
		return errors.New("block must not be nil")
	}

	return nil
}
