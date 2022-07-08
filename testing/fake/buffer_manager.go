package fake

import (
	golog "log"
	goos "os"

	"github.com/goropikari/simpledbgo/buffer"
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/lib/bytes"
	"github.com/goropikari/simpledbgo/log"
)

type NonDirectBufferManagerFactory struct {
	dbPath    string
	fileMgr   domain.FileManager
	logMgr    domain.LogManager
	bufferMgr domain.BufferPoolManager
}

func NewNonDirectBufferManagerFactory(dbPath string, blockSize int32, numBuf int) *NonDirectBufferManagerFactory {
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

	bufConfig := buffer.Config{
		NumberBuffer:       numBuf,
		TimeoutMillisecond: 10000,
	}

	bufMgr, err := buffer.NewManager(fileMgr, logMgr, pageFactory, bufConfig)
	if err != nil {
		golog.Fatal(err)
	}

	return &NonDirectBufferManagerFactory{
		dbPath:    dbPath,
		fileMgr:   fileMgr,
		logMgr:    logMgr,
		bufferMgr: bufMgr,
	}
}

func (factory *NonDirectBufferManagerFactory) Create() (domain.FileManager, domain.LogManager, domain.BufferPoolManager) {
	return factory.fileMgr, factory.logMgr, factory.bufferMgr
}

func (factory *NonDirectBufferManagerFactory) Finish() {
	goos.RemoveAll(factory.dbPath)
}
