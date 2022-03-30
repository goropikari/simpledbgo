package domain

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

type ByteSliceFactory interface {
	Create(int) ([]byte, error)
}
