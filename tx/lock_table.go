package tx

import (
	"errors"
	"sync"
	"time"

	"github.com/goropikari/simpledbgo/domain"
)

// ErrTransactionTimeoutExceeded is an error that means exceeding timeout.
var ErrTransactionTimeoutExceeded = errors.New("transaction timeout exceeded")

type result struct {
	err error
}

// LockTable manages locked Block which used by transaction.
type LockTable struct {
	mu                 *sync.Mutex
	cond               *sync.Cond
	locks              map[domain.Block]int
	timeoutMillisecond time.Duration
}

type LockTableConfig struct {
	LockTimeoutMillisecond int
}

func NewLockTableConfig() LockTableConfig {
	const timeout = 10000

	return LockTableConfig{LockTimeoutMillisecond: timeout}
}

// NewLockTable constructs LockTable.
func NewLockTable(cfg LockTableConfig) *LockTable {
	mu := &sync.Mutex{}
	cond := sync.NewCond(mu)

	return &LockTable{
		mu:                 mu,
		cond:               cond,
		locks:              make(map[domain.Block]int),
		timeoutMillisecond: time.Duration(cfg.LockTimeoutMillisecond) * time.Millisecond,
	}
}

func (lt *LockTable) SLock(blk domain.Block) error {
	done := make(chan *result)

	go lt.slock(done, blk)
	res := <-done

	return res.err
}

func (lt *LockTable) slock(done chan *result, blk domain.Block) {
	now := time.Now()
	defer close(done)
	for !lt.mu.TryLock() && time.Since(now) < lt.timeoutMillisecond {
	}
	if time.Since(now) >= lt.timeoutMillisecond {
		lt.mu.Unlock()
		done <- &result{err: ErrTransactionTimeoutExceeded}

		return
	}

	go func() {
		time.Sleep(lt.timeoutMillisecond - time.Since(now))
		lt.cond.Broadcast()
	}()

	for lt.hasXLock(blk) && time.Since(now) < lt.timeoutMillisecond {
		lt.cond.Wait()
	}
	if lt.hasXLock(blk) {
		lt.mu.Unlock()
		done <- &result{err: ErrTransactionTimeoutExceeded}

		return
	}
	val := lt.getLockVal(blk)
	lt.locks[blk] = val + 1
	done <- &result{}
	lt.mu.Unlock()
}

func (lt *LockTable) XLock(blk domain.Block) error {
	done := make(chan *result)

	go lt.xlock(done, blk)
	res := <-done

	return res.err
}

func (lt *LockTable) xlock(done chan *result, blk domain.Block) {
	now := time.Now()
	defer close(done)
	for !lt.mu.TryLock() && time.Since(now) < lt.timeoutMillisecond {
	}
	if time.Since(now) >= lt.timeoutMillisecond {
		lt.mu.Unlock()
		done <- &result{err: ErrTransactionTimeoutExceeded}

		return
	}

	go func() {
		time.Sleep(lt.timeoutMillisecond - time.Since(now))
		lt.cond.Broadcast()
	}()

	for lt.hasOtherSlocks(blk) && time.Since(now) < lt.timeoutMillisecond {
		lt.cond.Wait()
	}
	if lt.hasOtherSlocks(blk) {
		lt.mu.Unlock()
		done <- &result{err: ErrTransactionTimeoutExceeded}

		return
	}
	lt.locks[blk] = -1
	done <- &result{}
	lt.mu.Unlock()
}

func (lt *LockTable) Unlock(blk domain.Block) {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	val := lt.getLockVal(blk)
	if val > 1 {
		lt.locks[blk] = val - 1
		if val-1 == 1 {
			// lt.locks[blk] = 1 のとき、単に shared lock が取られているだけの場合もあれば、xlock 用の slock が取られている場合の2通りがある。
			lt.cond.Broadcast()
		}
	} else {
		delete(lt.locks, blk)
		lt.cond.Broadcast()
	}
}

func (lt *LockTable) hasXLock(blk domain.Block) bool {
	return lt.getLockVal(blk) < 0
}

func (lt *LockTable) hasOtherSlocks(blk domain.Block) bool {
	return lt.getLockVal(blk) > 1
}

func (lt *LockTable) getLockVal(blk domain.Block) int {
	n, ok := lt.locks[blk]
	if ok {
		return n
	}

	return 0
}

func (lt *LockTable) GetLockVal(blk domain.Block) int {
	return lt.getLockVal(blk)
}
