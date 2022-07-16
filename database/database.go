package database

import (
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/plan"
	"github.com/goropikari/simpledbgo/tx"
)

const (
	blockSize       = 4096
	numBuf          = 20
	timeoutMilliSec = 10000
)

type Config struct {
	DBPath          string
	BlockSize       int32
	NumBuf          int
	TimeoutMilliSec int
}

func NewConfig() Config {
	c := Config{
		BlockSize:       blockSize,
		NumBuf:          numBuf,
		TimeoutMilliSec: timeoutMilliSec,
	}

	return c
}

type DB struct {
	fmgr domain.FileManager
	lmgr domain.LogManager
	bmgr domain.BufferPoolManager
	cmgr domain.ConcurrencyManager
	gen  domain.TxNumberGenerator
	pe   *plan.Executor
}

func NewDB(
	fmgr domain.FileManager,
	lmgr domain.LogManager,
	bmgr domain.BufferPoolManager,
	cmgr domain.ConcurrencyManager,
	gen domain.TxNumberGenerator,
	pe *plan.Executor,
) *DB {
	return &DB{
		fmgr: fmgr,
		lmgr: lmgr,
		bmgr: bmgr,
		cmgr: cmgr,
		gen:  gen,
		pe:   pe,
	}
}

func (db *DB) NewTx() (domain.Transaction, error) {
	txn, err := tx.NewTransaction(db.fmgr, db.lmgr, db.bmgr, db.cmgr, db.gen)
	if err != nil {
		return nil, err
	}

	return txn, nil
}

func (db *DB) Query(txn domain.Transaction, query string) (domain.Planner, error) {
	return db.pe.CreateQueryPlan(query, txn)
}

func (db *DB) Exec(txn domain.Transaction, cmd string) (int, error) {
	return db.pe.ExecuteUpdate(cmd, txn)
}
