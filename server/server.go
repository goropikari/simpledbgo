package server

import (
	golog "log"

	"github.com/goropikari/simpledbgo/buffer"
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/file"
	"github.com/goropikari/simpledbgo/index/hash"
	"github.com/goropikari/simpledbgo/lib/bytes"
	"github.com/goropikari/simpledbgo/log"
	"github.com/goropikari/simpledbgo/metadata"
	"github.com/goropikari/simpledbgo/plan"
	"github.com/goropikari/simpledbgo/tx"
)

type DB struct {
	fmgr domain.FileManager
	lmgr domain.LogManager
	bmgr domain.BufferPoolManager
	cmgr domain.ConcurrencyManager
	gen  domain.TxNumberGenerator
	pe   *plan.Executor
}

type Config struct {
	DBPath          string
	BlockSize       int32
	NumBuf          int
	TimeoutMilliSec int
}

func NewConfig() Config {
	c := Config{
		BlockSize:       4096,
		NumBuf:          20,
		TimeoutMilliSec: 10000,
	}

	return c
}

func NewDB() *DB {
	cfg := NewConfig()

	// initialize file manager
	fileConfig := file.NewManagerConfig()
	fileMgr, err := file.NewManager(fileConfig)
	if err != nil {
		golog.Fatal(err)
	}

	blkSize, err := domain.NewBlockSize(cfg.BlockSize)
	if err != nil {
		golog.Fatal(err)
	}

	logConfig := log.ManagerConfig{LogFileName: "logfile"}
	logMgr, err := log.NewManager(fileMgr, logConfig)
	if err != nil {
		golog.Fatal(err)
	}

	bufConfig := buffer.Config{
		NumberBuffer:       cfg.NumBuf,
		TimeoutMillisecond: cfg.TimeoutMilliSec,
	}
	bsc := bytes.NewDirectByteSliceCreater()
	pageFactory := domain.NewPageFactory(bsc, blkSize)
	bufMgr, err := buffer.NewManager(fileMgr, logMgr, pageFactory, bufConfig)
	if err != nil {
		golog.Fatal(err)
	}

	ltConfig := tx.NewConfig(cfg.TimeoutMilliSec)
	lt := tx.NewLockTable(ltConfig)
	concurMgr := tx.NewConcurrencyManager(lt)

	gen := tx.NewNumberGenerator()

	txn, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
	if err != nil {
		golog.Fatal(err)
	}

	// _, err = goos.Stat(cfg.dbPath)
	// isNewDatabase := err != nil
	isNewDatabase := true
	idxDriver := domain.NewIndexDriver(hash.NewIndexFactory(), hash.NewSearchCostCalculator())
	var mmgr domain.MetadataManager
	if isNewDatabase {
		mmgr, err = metadata.CreateManager(idxDriver, txn)
		if err != nil {
			golog.Fatal(err)
		}
	} else {
		mmgr, err = metadata.NewManager(idxDriver, txn)
		if err != nil {
			golog.Fatal(err)
		}
	}

	err = txn.Commit()
	if err != nil {
		golog.Fatal(err)
	}

	qp := plan.NewBasicQueryPlanner(mmgr)
	ue := plan.NewBasicUpdatePlanner(mmgr)
	pe := plan.NewExecutor(qp, ue)

	return &DB{
		fmgr: fileMgr,
		lmgr: logMgr,
		bmgr: bufMgr,
		cmgr: concurMgr,
		gen:  gen,
		pe:   pe,
	}
}

func (db *DB) NewTx() (domain.Transaction, error) {
	return newTx(db.fmgr, db.lmgr, db.bmgr, db.cmgr, db.gen)
}

func (db *DB) Query(txn domain.Transaction, query string) (domain.Planner, error) {
	return db.pe.CreateQueryPlan(query, txn)
}

func (db *DB) Exec(txn domain.Transaction, cmd string) (int, error) {

	return db.pe.ExecuteUpdate(cmd, txn)
}

func newTx(fmgr domain.FileManager, lmgr domain.LogManager, bmgr domain.BufferPoolManager, cmgr domain.ConcurrencyManager, gen domain.TxNumberGenerator) (domain.Transaction, error) {
	txn, err := tx.NewTransaction(fmgr, lmgr, bmgr, cmgr, gen)
	if err != nil {
		return nil, err
	}

	return txn, nil
}
