package tx

import (
	"errors"
	"sync"
	"time"

	"github.com/goropikari/simpledbgo/domain"
)

// ErrTransactionTimeoutExceeded is an error that means exceeding timeout.
var ErrTransactionTimeoutExceeded = errors.New("transaction timeout exceeded")

// Config is configuration of LockTable.
type Config struct {
	lockTimeoutMillisecond int
}

// NewConfig constructs Config for LockTable.
func NewConfig(timeoutMilliSec int) Config {
	return Config{lockTimeoutMillisecond: timeoutMilliSec}
}

// LockTable manages locked Block which used by transaction.
type LockTable struct {
	locks              *sync.Map
	timeoutMillisecond time.Duration
}

// NewLockTable constructs LockTable.
func NewLockTable(config Config) *LockTable {
	return &LockTable{
		locks:              &sync.Map{},
		timeoutMillisecond: time.Duration(int64(config.lockTimeoutMillisecond)) * time.Millisecond,
	}
}

type result struct {
	err error
}

// SLock aquires shared lock on the blk.
func (lt *LockTable) SLock(blk domain.Block) error {
	done := make(chan *result)
	defer close(done)

	go lt.slock(done, blk)

	select {
	case result := <-done:
		if result.err != nil {
			return result.err
		}

		return nil
	case <-time.After(lt.timeoutMillisecond):
		return ErrTransactionTimeoutExceeded
	}
}

func (lt *LockTable) slock(done chan *result, blk domain.Block) {
	now := time.Now()
	lock, _ := lt.locks.LoadOrStore(blk, &sync.RWMutex{})
	lock.(*sync.RWMutex).RLock()

	if time.Since(now) > lt.timeoutMillisecond {
		lock.(*sync.RWMutex).RUnlock()

		return
	}

	done <- &result{err: nil}
}

// SUnlock releases shared lock on the blk.
func (lt *LockTable) SUnlock(blk domain.Block) {
	lock, loaded := lt.locks.Load(blk)
	if loaded {
		lock.(*sync.RWMutex).RUnlock()
	}
}

// XLock aquires exclusive lock on the blk.
func (lt *LockTable) XLock(blk domain.Block) error {
	done := make(chan *result)
	defer close(done)

	go lt.xlock(done, blk)

	select {
	case result := <-done:
		if result.err != nil {
			return result.err
		}

		return nil
	case <-time.After(lt.timeoutMillisecond):
		return ErrTransactionTimeoutExceeded
	}
}

func (lt *LockTable) xlock(done chan *result, blk domain.Block) {
	now := time.Now()
	lock, _ := lt.locks.LoadOrStore(blk, &sync.RWMutex{})
	lock.(*sync.RWMutex).Lock()

	if time.Since(now) > lt.timeoutMillisecond {
		lock.(*sync.RWMutex).Unlock()

		return
	}

	done <- &result{err: nil}
}

// XUnlock releases exclusive lock on the blk.
func (lt *LockTable) XUnlock(blk domain.Block) {
	lock, loaded := lt.locks.Load(blk)
	if loaded {
		lock.(*sync.RWMutex).Unlock()
	}
}
