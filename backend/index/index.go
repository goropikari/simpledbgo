package index

import (
	"github.com/goropikari/simpledbgo/backend/domain"
	"github.com/goropikari/simpledbgo/meta"
)

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// Index is an interface of index.
type Index interface {
	BeforeFirst(searchKey meta.Constant) error
	HasNext() (bool, error)
	GetDataRecordID() (domain.RecordID, error)
	Insert(meta.Constant, domain.RecordID) error
	Delete(meta.Constant, domain.RecordID) error
	Close()
	SearchCost(numBlocks int, recordPerBlock int) int
}
