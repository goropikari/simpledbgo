package domain

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// IndexGenerator generates Index.
type IndexGenerator interface {
	Create(Transaction, IndexName, *Layout) Index
}
