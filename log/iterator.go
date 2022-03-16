package log

import (
	"io"

	"github.com/goropikari/simpledb_go/core"
	"github.com/goropikari/simpledb_go/file"
)

// Iterator is iterator of log manager.
type Iterator struct {
	fileMgr               *file.Manager
	block                 *file.Block
	page                  *file.Page
	currentRecordPosition uint32
	boundary              uint32
}

func iterator(fileMgr *file.Manager, block *file.Block) (<-chan []byte, error) {
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

func newIterator(fileMgr *file.Manager, block *file.Block) (*Iterator, error) {
	if err := validateArgs(fileMgr, block); err != nil {
		return nil, err
	}

	page, err := fileMgr.PreparePage()
	if err != nil {
		return nil, err
	}

	if err := fileMgr.CopyBlockToPage(block, page); err != nil {
		return nil, err
	}

	boundary, err := page.GetUInt32(0)
	if err != nil && err != io.EOF {
		return nil, err
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
	blockSize, _ := logIt.fileMgr.GetBlockSize()

	return int(logIt.currentRecordPosition) < blockSize || logIt.block.GetBlockNumber() > 0
}

func (logIt *Iterator) next() []byte {
	blockSize, _ := logIt.fileMgr.GetBlockSize()

	if logIt.currentRecordPosition == uint32(blockSize) {
		block := file.NewBlock(logIt.block.GetFileName(), logIt.block.GetBlockNumber()-1)
		logIt.moveToBlock(block)
	}

	record, _ := logIt.page.GetBytes(int64(logIt.currentRecordPosition))
	logIt.currentRecordPosition += uint32(core.UInt32Length + len(record))

	return record
}

func (logIt *Iterator) moveToBlock(block *file.Block) {
	logIt.fileMgr.CopyBlockToPage(block, logIt.page)
	boundary, _ := logIt.page.GetUInt32(0)
	logIt.currentRecordPosition = boundary
	logIt.block = block
}

func validateArgs(fileMgr *file.Manager, block *file.Block) error {
	if fileMgr == nil {
		return core.NilReceiverError
	}
	if block == nil {
		return core.NilReceiverError
	}

	return nil
}
