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
	Commit() error
	Rollback() error
	Recover() error
	GetInt32(blk Block, offset int64) (val int32, err error)
	SetInt32(blk Block, offset int64, val int32, writeLog bool) error
	GetString(blk Block, offset int64) (val string, err error)
	SetString(blk Block, offset int64, val string, writeLog bool) error
	BlockLength(FileName) (int32, error)
	ExtendFile(FileName) (Block, error)
	BlockSize() BlockSize
	// Available() int
}

type TxNumberGenerator interface {
	Generate() TransactionNumber
}
