package domain

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// Indexer is an interface of index.
type Indexer interface {
	BeforeFirst(searchKey Constant) error
	HasNext() bool
	GetDataRecordID() (RecordID, error)
	Insert(Constant, RecordID) error
	Delete(Constant, RecordID) error
	Close()
}

// IndexName is a value object of index name.
type IndexName string

// NewIndexName constructs IndexName.
func NewIndexName(name string) (IndexName, error) {
	if len(name) > MaxFieldNameLength {
		return "", ErrExceedMaxFieldNameLength
	}

	return IndexName(name), nil
}

// String stringfy name.
func (name IndexName) String() string {
	return string(name)
}
