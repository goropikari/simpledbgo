package service

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/tests/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// LogManager is an interface of log manager.
type LogManager interface {
	// FlushByLSN flushes given lsn part.
	FlushByLSN(lsn int) error

	// AppendRecord appends a log record.
	AppendRecord(record []byte) error

	// Iterator returns iterator.
	Iterator() (<-chan []byte, error)
}
