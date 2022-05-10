package domain

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// SearchCostCalculator calculate search cost.
type SearchCostCalculator interface {
	Calculate(numBlk int, rpb int) int
}
