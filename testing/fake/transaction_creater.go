package fake

import (
	"log"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/backend/tx"
)

type TransactionCreater struct {
	factory   *NonDirectBufferManagerFactory
	FileMgr   domain.FileManager
	LogMgr    domain.LogManager
	BufMgr    domain.BufferManager
	ConcurMgr *tx.ConcurrencyManager
	Gen       *tx.NumberGenerator
}

func NewTransactionCreater(blockSize int32, numBuf int) *TransactionCreater {
	dbPath := RandString()
	factory := NewNonDirectBufferManagerFactory(dbPath, blockSize, numBuf)
	fileMgr, logMgr, bufMgr := factory.Create()

	ltConfig := tx.NewConfig(1000)
	lt := tx.NewLockTable(ltConfig)
	concurMgr := tx.NewConcurrencyManager(lt)

	gen := tx.NewNumberGenerator()

	return &TransactionCreater{
		factory:   factory,
		FileMgr:   fileMgr,
		LogMgr:    logMgr,
		BufMgr:    bufMgr,
		ConcurMgr: concurMgr,
		Gen:       gen,
	}
}

func (cr *TransactionCreater) NewTxn() domain.Transaction {
	txn, err := tx.NewTransaction(cr.FileMgr, cr.LogMgr, cr.BufMgr, cr.ConcurMgr, cr.Gen)
	if err != nil {
		log.Fatal(err)
	}

	return txn
}

func (cr *TransactionCreater) Finish() {
	cr.factory.Finish()
}
