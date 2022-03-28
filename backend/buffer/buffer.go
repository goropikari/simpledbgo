package buffer

import (
	"fmt"

	"github.com/goropikari/simpledb_go/backend/core"
	"github.com/goropikari/simpledb_go/backend/service"
)

// Buffer is a buffer of database.
type Buffer struct {
	fileMgr service.FileManager
	logMgr  service.LogManager
	page    *core.Page
	block   *core.Block
	pins    int
	txnum   int32
	lsn     int32
}

// NewBuffer creates a buffer.
func NewBuffer(fileMgr service.FileManager, logMgr service.LogManager) (*Buffer, error) {
	page, err := fileMgr.PreparePage()
	if err != nil {
		return nil, err
	}

	return &Buffer{
		fileMgr: fileMgr,
		logMgr:  logMgr,
		page:    page,
		block:   nil,
		pins:    0,
		txnum:   -1,
		lsn:     -1,
	}, nil
}

// GetBlock returns buffer's block.
func (buf *Buffer) GetBlock() *core.Block {
	return buf.block
}

// setModified modifing tx number and lsn.
func (buf *Buffer) setModified(txnum, lsn int32) {
	buf.txnum = txnum
	if lsn >= 0 {
		buf.lsn = lsn
	}
}

// isPinned checks whether the buffer is pinned or not.
func (buf *Buffer) isPinned() bool {
	return buf.pins > 0
}

// modifyingTx returns the transaction number which modifies the buffer.
func (buf *Buffer) modifyingTx() int32 {
	return buf.txnum
}

// assignToBlock assigns block to the buffer.
func (buf *Buffer) assignToBlock(block *core.Block) error {
	if err := buf.flush(); err != nil {
		return err
	}

	buf.block = block
	if err := buf.fileMgr.CopyBlockToPage(block, buf.page); err != nil {
		return fmt.Errorf("%w", err)
	}

	buf.pins = 0

	return nil
}

// flush flushes the buffer content.
func (buf *Buffer) flush() error {
	if buf.txnum >= 0 {
		if err := buf.logMgr.FlushByLSN(buf.txnum); err != nil {
			return fmt.Errorf("%w", err)
		}

		if err := buf.fileMgr.CopyPageToBlock(buf.page, buf.block); err != nil {
			return fmt.Errorf("%w", err)
		}

		buf.txnum = -1
	}

	return nil
}

// pin increments the number of pin of the buffer.
func (buf *Buffer) pin() {
	buf.pins++
}

// pin decrements the number of pin of the buffer.
func (buf *Buffer) unpin() {
	buf.pins--
}
