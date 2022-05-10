package domain

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// IndexFactory creates IndexGenerator and SearchCostCalculator.
type IndexFactory interface {
	Create() (IndexGenerator, SearchCostCalculator)
}
