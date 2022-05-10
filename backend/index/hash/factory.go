package hash

import "github.com/goropikari/simpledbgo/domain"

// IndexFactory is factory of IndexGenerator and SearchCostCalculator.
type IndexFactory struct{}

// NewIndexFactory constructs an IndexFactory.
func NewIndexFactory() *IndexFactory {
	return &IndexFactory{}
}

// Create creates IndexGenerator and SearchCostCalculator.
func (fac *IndexFactory) Create() (domain.IndexGenerator, domain.SearchCostCalculator) {
	return NewIndexGenerator(), NewSearchCostCalculator()
}

// IndexGenerator is generator of index.
type IndexGenerator struct{}

// NewIndexGenerator constructs an IndexGenerator.
func NewIndexGenerator() *IndexGenerator {
	return &IndexGenerator{}
}

// Create creates an Index.
func (gen *IndexGenerator) Create(txn domain.Transaction, idxName domain.IndexName, layout *domain.Layout) domain.Index {
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
