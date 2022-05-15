package fake

import (
	golog "log"
	goos "os"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/log"
	"github.com/goropikari/simpledbgo/lib/bytes"
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

	logConfig := log.ManagerConfig{LogFileName: "logfile_" + RandString()}
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

func (factory *NonDirectLogManagerFactory) Create() (domain.FileManager, domain.LogManager) {
	return factory.fileMgr, factory.logMgr
}

func (factory *NonDirectLogManagerFactory) Finish() {
	goos.RemoveAll(factory.dbPath)
}
