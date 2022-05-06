package domain

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// ByteSliceFactory is a factory of byte slice.
type ByteSliceFactory interface {
	Create(int) ([]byte, error)
}
