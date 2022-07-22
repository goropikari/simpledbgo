package dummy

import (
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/errors"
)

type IndexFactory struct{}

func NewIndexFactory() *IndexFactory {
	return &IndexFactory{}
}

func (fty *IndexFactory) Create(txn domain.Transaction, idxName domain.IndexName, layout *domain.Layout) domain.Indexer {
	return &Index{}
}

type SearchCostCalculator struct{}

func NewSearchCostCalculator() *SearchCostCalculator {
	return &SearchCostCalculator{}
}

func (cal *SearchCostCalculator) Calculate(numBlk int, rpb int) int {
	return 0
}

type Index struct{}

func (idx *Index) BeforeFirst(searchKey domain.Constant) error {
	return errors.ErrNotImplemented
}

func (idx *Index) HasNext() bool {
	return false
}

func (idx *Index) GetDataRecordID() (domain.RecordID, error) {
	return domain.RecordID{}, errors.ErrNotImplemented
}

func (idx *Index) Insert(key domain.Constant, rid domain.RecordID) error {
	return errors.ErrNotImplemented
}

func (idx *Index) Delete(key domain.Constant, rid domain.RecordID) error {
	return errors.ErrNotImplemented
}

func (idx *Index) Close() {}

func (idx *Index) Err() error {
	return errors.ErrNotImplemented
}
