package tx

import (
	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/lib/bytes"
)

// RecordType is type of log record.
type RecordType = int32

const (
	// Commit is commit record type.
	Commit RecordType = iota
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
	buf := bytes.NewBuffer(bytes.Int32Length * 2)
	if err := buf.SetInt32(0, Commit); err != nil {
		return domain.DummyLSN, nil
	}

	if err := buf.SetInt32(bytes.Int32Length, int32(tx.number)); err != nil {
		return domain.DummyLSN, nil
	}

	record := buf.GetData()

	return tx.logMgr.AppendRecord(record)
}
