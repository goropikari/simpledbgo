package fake

import (
	"log"
	goos "os"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/backend/file"
	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/goropikari/simpledb_go/os"
)

type NonDirectFileManagerFactory struct {
	dbPath string
	mgr    domain.FileManager
}

func NewNonDirectFileManagerFactory(dbPath string, blockSize int32) *NonDirectFileManagerFactory {
	bsf := bytes.NewByteSliceCreater()

	// initialize file manager
	explorer := os.NewNormalExplorer(dbPath)
	fileConfig := file.ManagerConfig{BlockSize: blockSize}
	fileMgr, err := file.NewManager(explorer, bsf, fileConfig)
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
