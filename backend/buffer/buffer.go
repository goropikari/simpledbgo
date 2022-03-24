package buffer

import (
	"fmt"

	"github.com/goropikari/simpledb_go/backend/core"
	"github.com/goropikari/simpledb_go/backend/service"
)

type Buffer struct {
	fileMgr service.FileManager
	logMgr  service.LogManager
	page    *core.Page
	block   *core.Block
	pins    int
	txnum   int
	lsn     int
}

func NewBuffer(fileMgr service.FileManager, logMgr service.LogManager) (*Buffer, error) {
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

func (buf *Buffer) GetBlock() *core.Block {
	return buf.block
}

func (buf *Buffer) GetInt32(offset int64) (int32, error) {
	n, err := buf.page.GetInt32(offset)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	return n, nil
}

func (buf *Buffer) setModified(txnum, lsn int) {
	buf.txnum = txnum
	if lsn >= 0 {
		buf.lsn = lsn
	}
}

// GetString returns string from page.
func (buf *Buffer) GetString(offset int64) (string, error) {
	s, err := buf.page.GetString(offset)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return s, nil
}

func (buf *Buffer) isPinned() bool {
	return buf.pins > 0
}

func (buf *Buffer) modifyingTx() int {
	return buf.txnum
}

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

func (buf *Buffer) pin() {
	buf.pins++
}

func (buf *Buffer) unpin() {
	buf.pins--
}
