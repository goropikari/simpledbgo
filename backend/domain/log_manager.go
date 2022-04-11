package domain

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// LogManager is an interface of log manager.
type LogManager interface {
	FlushLSN(LSN) error
	Flush() error
	AppendRecord([]byte) (LSN, error)
	AppendNewBlock() (*Block, error)
	Iterator() (LogIterator, error)
}

// LogIterator is a iterator of log record.
type LogIterator interface {
	HasNext() bool
	Next() ([]byte, error)
}
