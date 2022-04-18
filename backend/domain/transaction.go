package domain

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// TransactionNumber is transaction number.
type TransactionNumber int32

const (
	// DummyTransactionNumber is the dummy transaction number.
	DummyTransactionNumber TransactionNumber = -1
)

// Transaction is an interface of transaction.
type Transaction interface {
	Pin(Block) error
	Unpin(Block)
	// Commit() error
	// Rollback() error
	// Recover() error
	GetInt32(Block, int64) (int32, error)
	SetInt32(Block, int64, int32, bool) error
	GetString(Block, int64) (string, error)
	SetString(Block, int64, string, bool) error
	BlockLength(FileName) (int32, error)
	ExtendFile(FileName) (Block, error)
	BlockSize() BlockSize
	// Available() int
}
