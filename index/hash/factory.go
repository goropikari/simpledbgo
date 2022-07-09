package hash

import "github.com/goropikari/simpledbgo/domain"

// IndexFactory is generator of index.
type IndexFactory struct{}

// NewIndexFactory constructs an IndexFactory.
func NewIndexFactory() *IndexFactory {
	return &IndexFactory{}
}

// Create creates an Index.
func (gen *IndexFactory) Create(txn domain.Transaction, idxName domain.IndexName, layout *domain.Layout) domain.Indexer {
	return NewIndex(txn, idxName, layout)
}

// SearchCostCalculator calculates search cost.
type SearchCostCalculator struct{}

// NewSearchCostCalculator constructs a SearchCostCalculator.
func NewSearchCostCalculator() *SearchCostCalculator {
	return &SearchCostCalculator{}
}

// Calculate calculates search cost.
// rpb: record per block
func (cal *SearchCostCalculator) Calculate(numBlocks, rpb int) int {
	return numBlocks / numBuckets
}
