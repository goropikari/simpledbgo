package fake

import (
	golog "log"
	goos "os"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/backend/log"
	"github.com/goropikari/simpledb_go/lib/bytes"
)

type NonDirectLogManagerFactory struct {
	dbPath  string
	fileMgr domain.FileManager
	logMgr  domain.LogManager
}

func NewNonDirectLogManagerFactory(dbPath string, blockSize int32) *NonDirectLogManagerFactory {
	blkSize, err := domain.NewBlockSize(blockSize)
	if err != nil {
		golog.Fatal(err)
	}

	bsf := bytes.NewByteSliceCreater()
	pageFactory := domain.NewPageFactory(bsf, blkSize)

	fileMgrFactory := NewNonDirectFileManagerFactory(dbPath, blockSize)
	fileMgr := fileMgrFactory.Create()

	logConfig := log.ManagerConfig{LogFileName: RandString()}
	logMgr, err := log.NewManager(fileMgr, pageFactory, logConfig)
	if err != nil {
		golog.Fatal(err)
	}

	return &NonDirectLogManagerFactory{
		dbPath:  dbPath,
		fileMgr: fileMgr,
		logMgr:  logMgr,
	}
}

func (factory *NonDirectLogManagerFactory) Create() domain.LogManager {
	return factory.logMgr
}

func (factory *NonDirectLogManagerFactory) Finish() {
	goos.RemoveAll(factory.dbPath)
}