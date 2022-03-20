package buffer

import (
	"errors"
	"sync"
	"time"

	"github.com/goropikari/simpledb_go/file"
	"github.com/goropikari/simpledb_go/log"
)

// goroutine の数を制御
// timeout を設定する

type Manager struct {
	cond               *sync.Cond
	bufferPool         []*buffer
	numAvailableBuffer int
	maxTimeMilliSec    time.Duration
}

func NewManager(fileMgr *file.Manager, logMgr *log.Manager, numBuffer int) (*Manager, error) {
	if fileMgr == nil {
		return nil, errors.New("fileMgr must not be nil")
	}
	if logMgr == nil {
		return nil, errors.New("logMgr must not be nil")
	}

	bufferPool := make([]*buffer, 0, numBuffer)
	for i := 0; i < numBuffer; i++ {
		buf, err := newBuffer(fileMgr, logMgr)
		if err != nil {
			return nil, err
		}
		bufferPool = append(bufferPool, buf)
	}

	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	return &Manager{
		cond:               cond,
		bufferPool:         bufferPool,
		numAvailableBuffer: numBuffer,
		maxTimeMilliSec:    time.Second * 10, // 10 seconds
	}, nil
}

func (mgr *Manager) available() int {
	// if mgr == nil {
	// 	return 0, core.NilReceiverError
	// }

	mgr.cond.L.Lock()
	defer mgr.cond.L.Unlock()

	return mgr.numAvailableBuffer
}

// func (mgr *Manager) flushAll(txnum int) error {
// 	mgr.cond.L.Lock()
// 	defer mgr.cond.L.Unlock()
//
// 	for _, buf := range mgr.bufferPool {
// 		n, err := buf.modifyingTx()
// 		if err != nil {
// 			return err
// 		}
// 		if n == txnum {
// 			if err := buf.flush(); err != nil {
// 				return err
// 			}
// 		}
// 	}
//
// 	return nil
// }

func (mgr *Manager) unpin(buf *buffer) error {
	// if mgr == nil {
	// 	return core.NilReceiverError
	// }

	mgr.cond.L.Lock()
	defer mgr.cond.L.Unlock()

	buf.unpin()

	pinned := buf.isPinned()
	if !pinned {
		mgr.numAvailableBuffer++
		mgr.cond.Broadcast()
	}

	return nil
}

func (mgr *Manager) pin(block *file.Block) (*buffer, error) {
	// if mgr == nil {
	// 	return nil, core.NilReceiverError
	// }

	mgr.cond.L.Lock()
	defer mgr.cond.L.Unlock()

	startTime := time.Now()
	buf, err := mgr.tryToPin(block)
	if err != nil {
		return nil, err
	}
	for buf == nil && mgr.waitingTooLong(startTime) {
		mgr.cond.Wait()
		buf, err = mgr.tryToPin(block)
		if err != nil {
			return nil, err
		}
	}
	if buf == nil {
		return nil, errors.New("timeout")
	}

	return buf, nil
}

func (mgr *Manager) tryToPin(block *file.Block) (*buffer, error) {
	// if mgr == nil {
	// 	return nil, core.NilReceiverError
	// }

	buf := mgr.findExistingBuffer(block)
	if buf == nil {
		buf = mgr.chooseUnpinnedBuffer()
		if buf == nil {
			return nil, nil
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

func (mgr *Manager) findExistingBuffer(block *file.Block) *buffer {
	for _, buf := range mgr.bufferPool {
		other := buf.getBlock()
		if other != nil && other.Equal(block) {
			return buf
		}
	}

	return nil
}

func (mgr *Manager) chooseUnpinnedBuffer() *buffer {
	for _, buf := range mgr.bufferPool {
		if !buf.isPinned() {
			return buf
		}
	}

	return nil
}

func (mgr *Manager) waitingTooLong(start time.Time) bool {
	return time.Since(start) > mgr.maxTimeMilliSec
}
