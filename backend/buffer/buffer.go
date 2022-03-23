package buffer

import (
	"fmt"

	"github.com/goropikari/simpledb_go/backend/core"
	"github.com/goropikari/simpledb_go/backend/service"
)

type buffer struct {
	fileMgr service.FileManager
	logMgr  service.LogManager
	page    *core.Page
	block   *core.Block
	pins    int
	txnum   int
	lsn     int
}

func newBuffer(fileMgr service.FileManager, logMgr service.LogManager) (*buffer, error) {
	if fileMgr.IsZero() {
		return nil, ErrInvalidArgs
	}

	if logMgr.IsZero() {
		return nil, ErrInvalidArgs
	}

	page, err := fileMgr.PreparePage()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &buffer{
		fileMgr: fileMgr,
		logMgr:  logMgr,
		page:    page,
		block:   nil,
		pins:    0,
		txnum:   -1,
		lsn:     -1,
	}, nil
}

func (buf *buffer) getBlock() *core.Block {
	if buf == nil {
		return nil
	}

	return buf.block
}

func (buf *buffer) setModified(txnum, lsn int) {
	if buf == nil {
		return
	}

	buf.txnum = txnum
	if lsn >= 0 {
		buf.lsn = lsn
	}
}

func (buf *buffer) isPinned() bool {
	if buf == nil {
		return false
	}

	return buf.pins > 0
}

// func (buf *buffer) modifyingTx() int {
// 	if buf == nil {
// 		return -1
// 	}

// 	return buf.txnum
// }

func (buf *buffer) assignToBlock(block *core.Block) error {
	if buf == nil {
		return nil
	}

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

func (buf *buffer) flush() error {
	if buf == nil {
		return nil
	}

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

func (buf *buffer) pin() {
	buf.pins++
}

func (buf *buffer) unpin() {
	buf.pins--
}
