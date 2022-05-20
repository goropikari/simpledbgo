package domain

import (
	"github.com/pkg/errors"
)

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// Buffer is a buffer of database.
type Buffer struct {
	fileMgr FileManager
	logMgr  LogManager
	page    *Page
	block   Block
	pins    int
	txnum   TransactionNumber
	lsn     LSN
}

// NewBuffer creates a buffer.
func NewBuffer(fileMgr FileManager, logMgr LogManager, pageFactory *PageFactory) (*Buffer, error) {
	page, err := pageFactory.Create()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a page.")
	}

	return &Buffer{
		fileMgr: fileMgr,
		logMgr:  logMgr,
		page:    page,
		block:   Block{},
		pins:    0,
		txnum:   DummyTransactionNumber,
		lsn:     DummyLSN,
	}, nil
}

// Block returns buffer's block.
func (buf *Buffer) Block() Block {
	return buf.block
}

// SetModifiedTxNumber modifing tx number and lsn.
func (buf *Buffer) SetModifiedTxNumber(txnum TransactionNumber, lsn LSN) {
	buf.txnum = txnum
	if lsn >= 0 {
		buf.lsn = lsn
	}
}

// IsPinned checks whether the buffer is pinned or not.
func (buf *Buffer) IsPinned() bool {
	return buf.pins > 0
}

// TxNumber returns the transaction number which modifies the buffer.
func (buf *Buffer) TxNumber() TransactionNumber {
	return buf.txnum
}

// AssignToBlock assigns block to the buffer.
func (buf *Buffer) AssignToBlock(block Block) error {
	err := buf.Flush()
	if err != nil {
		return errors.Wrap(err, "failed to flush the buffer")
	}

	buf.block = block
	err = buf.fileMgr.CopyBlockToPage(block, buf.page)
	if err != nil {
		return errors.Wrap(err, "failed to flush block")
	}

	buf.pins = 0

	return nil
}

// Flush flushes the buffer content.
func (buf *Buffer) Flush() error {
	if buf.txnum >= 0 {
		err := buf.logMgr.FlushLSN(buf.lsn)
		if err != nil {
			return errors.Wrap(err, "failed to flush lsn block")
		}

		err = buf.fileMgr.CopyPageToBlock(buf.page, buf.block)
		if err != nil {
			return errors.Wrap(err, "failed to copy page to block")
		}

		buf.txnum = DummyTransactionNumber
	}

	return nil
}

// Pin increments the number of pin of the buffer.
func (buf *Buffer) Pin() {
	buf.pins++
}

// Unpin decrements the number of pin of the buffer.
func (buf *Buffer) Unpin() {
	buf.pins--
}

// Page returns buffer's page.
func (buf *Buffer) Page() *Page {
	return buf.page
}
