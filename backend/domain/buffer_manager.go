package domain

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// BufferManager is an interface of buffer manager.
type BufferManager interface {
	Available() int
	FlushAll(txnum TransactionNumber) error
	Unpin(buf *Buffer)
	Pin(*Block) (*Buffer, error)
}
