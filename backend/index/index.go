package index

import "github.com/goropikari/simpledbgo/backend/domain"

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// Index is an interface of index.
type Index interface {
	BeforeFirst(searchKey domain.Constant) error
	HasNext() bool
	GetDataRecordID() (domain.RecordID, error)
	Insert(domain.Constant, domain.RecordID) error
	Delete(domain.Constant, domain.RecordID) error
	Close()
}
