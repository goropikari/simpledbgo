package domain

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// LogManager is an interface of log manager.
type LogManager interface {
	FlushLSN(int32) error
	Flush() error
	AppendRecord([]byte) (int32, error)
	AppendNewBlock() (*Block, error)
}
