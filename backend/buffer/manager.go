package buffer

import (
	"errors"
	"sync"
	"time"

	"github.com/goropikari/simpledb_go/backend/domain"
)

const maxTimeoutSecond = 10

// // ErrFailedPin is an error type that means failed to pin block.
// ErrFailedPin = errors.New("failed to pin block")
//
// // ErrTimeoutExceeded is an error type that means timeout exceeded.
// ErrTimeoutExceeded = errors.New("timeout exceeded")
//
// // ErrInvalidArgs is an error that means given args is invalid.
// ErrInvalidArgs = errors.New("arguments is invalid")

// ErrInvalidNumberOfBuffer is an error that means number of buffer must be positive.
var ErrInvalidNumberOfBuffer = errors.New("number of buffer must be positive")

// Config is a configure of buffer manager.
type Config struct {
	NumberBuffer int
}

// Manager is model of buffer manager.
type Manager struct {
	cond               *sync.Cond
	bufferPool         []*domain.Buffer
	numAvailableBuffer int
	timeout            time.Duration
}

// NewManager is a constructor of Manager.
func NewManager(fileMgr domain.FileManager, logMgr domain.LogManager, pageFactory *domain.PageFactory, config Config) (*Manager, error) {
	numBuffer := config.NumberBuffer
	if numBuffer <= 0 {
		return nil, ErrInvalidNumberOfBuffer
	}

	bufferPool := make([]*domain.Buffer, numBuffer)

	for i := 0; i < numBuffer; i++ {
		buf, err := domain.NewBuffer(fileMgr, logMgr, pageFactory)
		if err != nil {
			return nil, err
		}

		bufferPool[i] = buf
	}

	return &Manager{
		cond:               sync.NewCond(&sync.Mutex{}),
		bufferPool:         bufferPool,
		numAvailableBuffer: numBuffer,
		timeout:            time.Second * maxTimeoutSecond,
	}, nil
}
