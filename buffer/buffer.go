package buffer

import (
	"errors"

	"github.com/goropikari/simpledb_go/file"
	"github.com/goropikari/simpledb_go/log"
)

type buffer struct {
	fileMgr *file.Manager
	logMgr  *log.Manager
	page    *file.Page
	block   *file.Block
	pins    int
	txnum   int
	lsn     int
}

func newBuffer(fileMgr *file.Manager, logMgr *log.Manager) (*buffer, error) {
	if fileMgr == nil {
		return nil, errors.New("fileMgr must not be nil")
	}
	if logMgr == nil {
		return nil, errors.New("fileMgr must not be nil")
	}

	page, err := fileMgr.PreparePage()
	if err != nil {
		return nil, err
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

func (buf *buffer) getBlock() *file.Block {
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

// func (buf *buffer) modifyingTx() (int, error) {
// 	if buf == nil {
// 		return -1, core.NilReceiverError
// 	}
//
// 	return buf.txnum, nil
// }

func (buf *buffer) assignToBlock(block *file.Block) error {
	if buf == nil {
		return nil
	}

	if err := buf.flush(); err != nil {
		return err
	}

	buf.block = block
	if err := buf.fileMgr.CopyBlockToPage(block, buf.page); err != nil {
		return err
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
			return err
		}

		if err := buf.fileMgr.CopyPageToBlock(buf.page, buf.block); err != nil {
			return err
		}
	}

	return nil
}

func (buf *buffer) pin() {
	buf.pins++
}

func (buf *buffer) unpin() {
	buf.pins--
}
