package btree

import "math"

// SearchCostCalculator calculates search cost.
type SearchCostCalculator struct{}

// NewSearchCostCalculator constructs a SearchCostCalculator.
func NewSearchCostCalculator() *SearchCostCalculator {
	return &SearchCostCalculator{}
}

// Calculate calculates search cost.
// rpb: record per block
func (cal *SearchCostCalculator) Calculate(numBlocks, rpb int) int {
	return 1 + int(math.Round(math.Log(float64(numBlocks))/math.Log(float64(rpb))))
}
