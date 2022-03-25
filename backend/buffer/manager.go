package buffer

import (
	"errors"
	"sync"
	"time"

	"github.com/goropikari/simpledb_go/backend/core"
	"github.com/goropikari/simpledb_go/backend/service"
)

const maxTimeoutSecond = 10

var (
	// ErrFailedPin is an error type that means failed to pin block.
	ErrFailedPin = errors.New("failed to pin block")

	// ErrTimeout is an error type that means timeout exceeded.
	ErrTimeoutExceeded = errors.New("timeout exceeded")

	// ErrArgsInvalid is an error that means given args is invalid.
	ErrInvalidArgs = errors.New("arguments is invalid")
)

// Manager is model of buffer manager.
type Manager struct {
	cond               *sync.Cond
	bufferPool         []*Buffer
	numAvailableBuffer int
	timeout            time.Duration
}

// NewManager is a constructor of Manager.
func NewManager(fileMgr service.FileManager, logMgr service.LogManager, numBuffer int) (*Manager, error) {
	if numBuffer <= 0 {
		return nil, ErrInvalidArgs
	}

	bufferPool := make([]*Buffer, 0, numBuffer)

	for i := 0; i < numBuffer; i++ {
		buf, err := NewBuffer(fileMgr, logMgr)
		if err != nil {
			return nil, err
		}

		bufferPool = append(bufferPool, buf)
	}

	cond := sync.NewCond(&sync.Mutex{})

	return &Manager{
		cond:               cond,
		bufferPool:         bufferPool,
		numAvailableBuffer: numBuffer,
		timeout:            time.Second * maxTimeoutSecond,
	}, nil
}

// available returns the number of unpinned buffer.
func (mgr *Manager) available() int {
	mgr.cond.L.Lock()
	defer mgr.cond.L.Unlock()

	return mgr.numAvailableBuffer
}

// FlushAll flushes all record about given txnum.
func (mgr *Manager) FlushAll(txnum int) error {
	mgr.cond.L.Lock()
	defer mgr.cond.L.Unlock()

	for _, buf := range mgr.bufferPool {
		n := buf.modifyingTx()
		if n == txnum {
			if err := buf.flush(); err != nil {
				return err
			}
		}
	}

	return nil
}

// unpin unpins the buffer.
func (mgr *Manager) unpin(buf *Buffer) error {
	mgr.cond.L.Lock()
	defer mgr.cond.L.Unlock()

	buf.unpin()

	if pinned := buf.isPinned(); !pinned {
		mgr.numAvailableBuffer++
		mgr.cond.Broadcast()
	}

	return nil
}

// pin pins the block and return pinned buffer.
func (mgr *Manager) pin(block *core.Block) (*Buffer, error) {
	mgr.cond.L.Lock()
	defer mgr.cond.L.Unlock()

	startTime := time.Now()

	buf, err := mgr.tryToPin(block, chooseUnpinnedBuffer)
	if err != nil && !errors.Is(err, ErrFailedPin) {
		return nil, err
	}

	for buf == nil && mgr.waitingTooLong(startTime) {
		mgr.cond.Wait()
		buf, err = mgr.tryToPin(block, chooseUnpinnedBuffer)

		if err != nil && !errors.Is(err, ErrFailedPin) {
			return nil, err
		}
	}

	if buf == nil {
		return nil, ErrTimeoutExceeded
	}

	return buf, nil
}

// tryToPin tries to pin the block to a buffer.
func (mgr *Manager) tryToPin(block *core.Block, chooseUnpinnedBuffer func([]*Buffer) *Buffer) (*Buffer, error) {
	buf := mgr.findExistingBuffer(block)
	if buf == nil {
		buf = chooseUnpinnedBuffer(mgr.bufferPool)
		if buf == nil {
			return nil, ErrFailedPin
		}

		if err := buf.assignToBlock(block); err != nil {
			return nil, err
		}
	}

	if !buf.isPinned() {
		mgr.numAvailableBuffer--
	}

	buf.pin()

	return buf, nil
}

// findExistingBuffer returns the buffer whose block is same as given block.
// If there is no such buffer, returns nil.
func (mgr *Manager) findExistingBuffer(block *core.Block) *Buffer {
	for _, buf := range mgr.bufferPool {
		other := buf.GetBlock()
		if other != nil && other.Equal(block) {
			return buf
		}
	}

	return nil
}

// chooseUnpinnedBuffer chooses unpinned buffer.
func chooseUnpinnedBuffer(bufferPool []*Buffer) *Buffer {
	for _, buf := range bufferPool {
		if !buf.isPinned() {
			return buf
		}
	}

	return nil
}

// waitingTooLong checks whether if wait time is too long or not.
func (mgr *Manager) waitingTooLong(start time.Time) bool {
	return time.Since(start) > mgr.timeout
}
