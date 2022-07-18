package database

import (
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/errors"
	"github.com/goropikari/simpledbgo/plan"
	"github.com/goropikari/simpledbgo/tx"
)

const (
	blockSize       = 4096
	numBuf          = 20
	timeoutMilliSec = 10000
)

// Config is configuration for server.
type Config struct {
	DBPath          string
	BlockSize       int32
	NumBuf          int
	TimeoutMilliSec int
}

// NewConfig constructs a Config.
func NewConfig() Config {
	c := Config{
		BlockSize:       blockSize,
		NumBuf:          numBuf,
		TimeoutMilliSec: timeoutMilliSec,
	}

	return c
}

// DB is database.
type DB struct {
	fmgr domain.FileManager
	lmgr domain.LogManager
	bmgr domain.BufferPoolManager
	lt   *tx.LockTable
	gen  domain.TxNumberGenerator
	pe   *plan.Executor
}

// NewDB constructs a DB.
func NewDB(
	fmgr domain.FileManager,
	lmgr domain.LogManager,
	bmgr domain.BufferPoolManager,
	lt *tx.LockTable,
	gen domain.TxNumberGenerator,
	pe *plan.Executor,
) *DB {
	return &DB{
		fmgr: fmgr,
		lmgr: lmgr,
		bmgr: bmgr,
		lt:   lt,
		gen:  gen,
		pe:   pe,
	}
}

// NewTx make a new transaction.
func (db *DB) NewTx() (domain.Transaction, error) {
	txn, err := tx.NewTransaction(db.fmgr, db.lmgr, db.bmgr, db.lt, db.gen)
	if err != nil {
		return nil, errors.Err(err, "NewTransaction")
	}

	return txn, nil
}

// Query queries given sql.
func (db *DB) Query(txn domain.Transaction, query string) (domain.Planner, error) {
	return db.pe.CreateQueryPlan(query, txn)
}

// Exec executes a command.
func (db *DB) Exec(txn domain.Transaction, cmd string) (int, error) {
	return db.pe.ExecuteUpdate(cmd, txn)
}
