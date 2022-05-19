package server

import (
	golog "log"
	goos "os"
	"path"

	"github.com/goropikari/simpledbgo/buffer"
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/file"
	"github.com/goropikari/simpledbgo/index/hash"
	"github.com/goropikari/simpledbgo/lib/bytes"
	"github.com/goropikari/simpledbgo/log"
	"github.com/goropikari/simpledbgo/metadata"
	"github.com/goropikari/simpledbgo/os"
	"github.com/goropikari/simpledbgo/plan"
	"github.com/goropikari/simpledbgo/tx"
)

type DB struct {
	fmgr domain.FileManager
	lmgr domain.LogManager
	bmgr domain.BufferManager
	cmgr domain.ConcurrencyManager
	gen  domain.TxNumberGenerator
	pe   *plan.Executor
}

type config struct {
	dbPath          string
	blockSize       int32
	numBuf          int
	timeoutMilliSec int
}

func newConfig() config {
	dbPath := goos.Getenv("SIMPLEDB_PATH")

	return config{
		dbPath: dbPath,
	}
}

func (c *config) setDefault() {
	if path := goos.Getenv("SIMPLEDB_PATH"); path != "" {
		c.dbPath = path
	}
	if c.dbPath == "" {
		c.dbPath = path.Join(goos.Getenv("HOME"), "simpledb")
	}
	if c.blockSize == 0 {
		c.blockSize = 4096
	}
	if c.numBuf == 0 {
		c.numBuf = 20
	}
	if c.timeoutMilliSec == 0 {
		c.timeoutMilliSec = 10000
	}
}

func NewDB() *DB {
	cfg := config{}
	cfg.setDefault()

	blkSize, err := domain.NewBlockSize(cfg.blockSize)
	if err != nil {
		golog.Fatal(err)
	}

	bsc := bytes.NewDirectByteSliceCreater()
	pageFactory := domain.NewPageFactory(bsc, blkSize)

	_, err = goos.Stat(cfg.dbPath)
	isNewDatabase := err != nil

	// initialize file manager
	explorer := os.NewDirectIOExplorer(cfg.dbPath)
	fileConfig := file.ManagerConfig{BlockSize: cfg.blockSize}
	fileMgr, err := file.NewManager(explorer, bsc, fileConfig)
	if err != nil {
		golog.Fatal(err)
	}

	logConfig := log.ManagerConfig{LogFileName: "logfile"}
	logMgr, err := log.NewManager(fileMgr, pageFactory, logConfig)
	if err != nil {
		golog.Fatal(err)
	}

	bufConfig := buffer.Config{
		NumberBuffer:       cfg.numBuf,
		TimeoutMillisecond: cfg.timeoutMilliSec,
	}
	bufMgr, err := buffer.NewManager(fileMgr, logMgr, pageFactory, bufConfig)
	if err != nil {
		golog.Fatal(err)
	}

	ltConfig := tx.NewConfig(cfg.timeoutMilliSec)
	lt := tx.NewLockTable(ltConfig)
	concurMgr := tx.NewConcurrencyManager(lt)

	gen := tx.NewNumberGenerator()

	txn, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
	if err != nil {
		golog.Fatal(err)
	}

	fac := hash.NewIndexFactory()
	var mmgr domain.MetadataManager
	if isNewDatabase {
		mmgr, err = metadata.CreateManager(fac, txn)
		if err != nil {
			golog.Fatal(err)
		}
	} else {
		mmgr, err = metadata.NewManager(fac, txn)
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

func newTx(fmgr domain.FileManager, lmgr domain.LogManager, bmgr domain.BufferManager, cmgr domain.ConcurrencyManager, gen domain.TxNumberGenerator) (domain.Transaction, error) {
	txn, err := tx.NewTransaction(fmgr, lmgr, bmgr, cmgr, gen)
	if err != nil {
		return nil, err
	}

	return txn, nil
}
