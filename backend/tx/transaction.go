package tx

import (
	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/backend/tx/logrecord"
	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/goropikari/simpledb_go/typ"
)

// RecordType is type of log record.
type RecordType = int32

const (
	// Commit is commit record type.
	Commit RecordType = iota

	// SetInt32 is set int32 record type.
	SetInt32
)

// Transaction is a model of transaction.
type Transaction struct {
	fileMgr    domain.FileManager
	logMgr     domain.LogManager
	bufferMgr  domain.BufferManager
	concurMgr  *ConcurrencyManager
	bufferList *BufferList
	number     domain.TransactionNumber
}

// NewTransaction constructs Transaction.
func NewTransaction(fileMgr domain.FileManager, logMgr domain.LogManager, bufferMgr domain.BufferManager, concurMgr *ConcurrencyManager, gen *NumberGenerator) *Transaction {
	return &Transaction{
		fileMgr:    fileMgr,
		logMgr:     logMgr,
		bufferMgr:  bufferMgr,
		concurMgr:  concurMgr,
		bufferList: NewBufferList(bufferMgr),
		number:     gen.Generate(),
	}
}

// Commit commits the transaction.
func (tx *Transaction) Commit() error {
	if err := tx.commit(); err != nil {
		return err
	}

	tx.concurMgr.Release()
	tx.bufferList.UnpinAll()

	return nil
}

func (tx *Transaction) commit() error {
	if err := tx.bufferMgr.FlushAll(tx.number); err != nil {
		return err
	}

	lsn, err := tx.writeCommitLog()
	if err != nil {
		return err
	}

	if err := tx.logMgr.FlushLSN(lsn); err != nil {
		return err
	}

	return nil
}

func (tx *Transaction) writeCommitLog() (domain.LSN, error) {
	buf := bytes.NewBuffer(typ.Int32Length * 2)
	if err := buf.SetInt32(0, Commit); err != nil {
		return domain.DummyLSN, err
	}

	if err := buf.SetInt32(typ.Int32Length, int32(tx.number)); err != nil {
		return domain.DummyLSN, err
	}

	return tx.logMgr.AppendRecord(buf.GetData())
}

// GetInt32 gets int32 from the blk at offset.
func (tx *Transaction) GetInt32(blk domain.Block, offset int64) (int32, error) {
	if err := tx.concurMgr.SLock(blk); err != nil {
		return 0, err
	}

	buf := tx.bufferList.GetBuffer(blk)
	x, err := buf.Page().GetInt32(offset)
	if err != nil {
		return 0, err
	}

	return x, nil
}

// SetInt32 sets int32 on the given block.
func (tx *Transaction) SetInt32(blk domain.Block, offset int64, val int32, writeLog bool) error {
	if err := tx.concurMgr.XLock(blk); err != nil {
		return err
	}

	buf := tx.bufferList.GetBuffer(blk)
	lsn := domain.DummyLSN
	if writeLog {
		var err error
		oldval, err := buf.Page().GetInt32(offset)
		if err != nil {
			return err
		}

		lsn, err = tx.writeSetInt32Log(*buf.Block(), offset, oldval)
		if err != nil {
			return err
		}
	}

	if err := buf.Page().SetInt32(offset, val); err != nil {
		return err
	}

	buf.SetModifiedTxNumber(tx.number, lsn)

	return nil
}

func (tx *Transaction) writeSetInt32Log(blk domain.Block, offset int64, val int32) (domain.LSN, error) {
	record := &logrecord.SetInt32Record{
		FileName:    blk.FileName(),
		TxNum:       tx.number,
		BlockNumber: blk.Number(),
		Offset:      offset,
		Val:         val,
	}

	data, err := record.Marshal()
	if err != nil {
		return domain.DummyLSN, err
	}

	buf := bytes.NewBuffer(typ.Int32Length*2 + len(data))
	if err := buf.SetInt32(0, SetInt32); err != nil {
		return domain.DummyLSN, err
	}

	if err := buf.SetBytes(typ.Int32Length, data); err != nil {
		return domain.DummyLSN, err
	}

	lsn, err := tx.logMgr.AppendRecord(buf.GetData())
	if err != nil {
		return domain.DummyLSN, err
	}

	return lsn, nil
}

// Pin pins the blk by tx.
func (tx *Transaction) Pin(blk domain.Block) error {
	if err := tx.bufferList.Pin(blk); err != nil {
		return err
	}

	return nil
}
