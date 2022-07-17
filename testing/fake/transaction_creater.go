package fake

import (
	"log"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/tx"
)

type TransactionCreater struct {
	factory *NonDirectBufferManagerFactory
	FileMgr domain.FileManager
	LogMgr  domain.LogManager
	BufMgr  domain.BufferPoolManager
	LockTbl *tx.LockTable
	Gen     *tx.NumberGenerator
}

func NewTransactionCreater(blockSize int32, numBuf int) *TransactionCreater {
	dbPath := RandString()
	factory := NewNonDirectBufferManagerFactory(dbPath, blockSize, numBuf)
	fileMgr, logMgr, bufMgr := factory.Create()

	cfg := tx.LockTableConfig{LockTimeoutMillisecond: 1000}
	lt := tx.NewLockTable(cfg)
	gen := tx.NewNumberGenerator()

	return &TransactionCreater{
		factory: factory,
		FileMgr: fileMgr,
		LogMgr:  logMgr,
		BufMgr:  bufMgr,
		LockTbl: lt,
		Gen:     gen,
	}
}

func (cr *TransactionCreater) NewTxn() domain.Transaction {
	txn, err := tx.NewTransaction(cr.FileMgr, cr.LogMgr, cr.BufMgr, cr.LockTbl, cr.Gen)
	if err != nil {
		log.Fatal(err)
	}

	return txn
}

func (cr *TransactionCreater) Finish() {
	cr.factory.Finish()
}
