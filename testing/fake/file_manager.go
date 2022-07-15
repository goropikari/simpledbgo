package fake

import (
	"log"
	goos "os"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/file"
)

type NonDirectFileManagerFactory struct {
	dbPath string
	mgr    domain.FileManager
}

func NewNonDirectFileManagerFactory(dbPath string, blockSize int32) *NonDirectFileManagerFactory {
	// initialize file manager
	fileConfig := file.ManagerConfig{
		DBPath:    dbPath,
		BlockSize: blockSize,
		DirectIO:  false,
	}
	fileMgr, err := file.NewManager(fileConfig)
	if err != nil {
		log.Fatal(err)
	}

	return &NonDirectFileManagerFactory{
		dbPath: dbPath,
		mgr:    fileMgr,
	}
}

func (factory *NonDirectFileManagerFactory) Create() domain.FileManager {
	return factory.mgr
}

func (factory *NonDirectFileManagerFactory) Finish() {
	goos.RemoveAll(factory.dbPath)
}
