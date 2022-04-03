package buffer

import (
	"errors"
	"sync"
	"time"

	"github.com/goropikari/simpledb_go/backend/domain"
)

const item = 1

var (
	// ErrTimeoutExceeded is an error type that means timeout exceeded.
	ErrTimeoutExceeded = errors.New("timeout exceeded")

	// ErrInvalidNumberOfBuffer is an error that means number of buffer must be positive.
	ErrInvalidNumberOfBuffer = errors.New("number of buffer must be positive")
)

// Config is a configure of buffer manager.
type Config struct {
	NumberBuffer       int
	TimeoutMillisecond int
}

// Manager is model of buffer manager.
type Manager struct {
	mu                 sync.Mutex
	ch                 chan int
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

	ch := make(chan int, config.NumberBuffer)
	for i := 0; i < config.NumberBuffer; i++ {
		ch <- item
	}

	return &Manager{
		mu:                 sync.Mutex{},
		ch:                 ch,
		bufferPool:         bufferPool,
		numAvailableBuffer: numBuffer,
		timeout:            time.Millisecond * time.Duration(config.TimeoutMillisecond),
	}, nil
}

// Available returns the number of available buffers.
func (mgr *Manager) Available() int {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	return mgr.numAvailableBuffer
}

// FlushAll flushes buffer with specified transaction number.
func (mgr *Manager) FlushAll(txnum domain.TransactionNumber) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	for _, buf := range mgr.bufferPool {
		if buf.TxNumber() == txnum {
			err := buf.Flush()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Unpin unpins buffer.
func (mgr *Manager) Unpin(buf *domain.Buffer) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	buf.Unpin()
	if !buf.IsPinned() {
		mgr.numAvailableBuffer++
		mgr.ch <- item
	}
}

// Pin pins buffer.
func (mgr *Manager) Pin(block *domain.Block) (*domain.Buffer, error) {
	mgr.mu.Lock()

	buf, err := mgr.tryToPin(block, naiveSearchUnpinnedBuffer)
	if err != nil {
		mgr.mu.Unlock()

		return nil, err
	}

	if buf == nil {
		mgr.mu.Unlock()
		select {
		case <-mgr.ch:
			mgr.mu.Lock()
			mgr.ch <- item

			buf, err = mgr.tryToPin(block, naiveSearchUnpinnedBuffer)
			if err != nil {
				mgr.mu.Unlock()

				return nil, err
			}
		case <-time.After(mgr.timeout):
			return nil, ErrTimeoutExceeded
		}
	}

	mgr.mu.Unlock()

	return buf, nil
}

// tryToPin tries to pin the block to a buffer.
func (mgr *Manager) tryToPin(block *domain.Block, chooseUnpinnedBuffer func([]*domain.Buffer) *domain.Buffer) (*domain.Buffer, error) {
	buf := mgr.findExistingBuffer(block)
	if buf == nil {
		buf = chooseUnpinnedBuffer(mgr.bufferPool)
		if buf == nil {
			return buf, nil
		}
		if err := buf.AssignToBlock(block); err != nil {
			return nil, err
		}
	}

	if !buf.IsPinned() {
		mgr.numAvailableBuffer--
		<-mgr.ch
	}

	buf.Pin()

	return buf, nil
}

// findExistingBuffer returns the buffer whose block is same as given block.
// If there is no such buffer, returns nil.
func (mgr *Manager) findExistingBuffer(block *domain.Block) *domain.Buffer {
	for _, buf := range mgr.bufferPool {
		if block.Equal(buf.Block()) {
			return buf
		}
	}

	return nil
}

// naiveSearchUnpinnedBuffer chooses unpinned buffer by linear search.
func naiveSearchUnpinnedBuffer(bufferPool []*domain.Buffer) *domain.Buffer {
	for _, buf := range bufferPool {
		if !buf.IsPinned() {
			return buf
		}
	}

	return nil
}
