package tx

import "github.com/goropikari/simpledb_go/backend/domain"

// LockType is a type of lock.
type LockType int8

const (
	// Shared means shared lock.
	Shared LockType = iota

	// Exclusive means exclusive lock.
	Exclusive
)

// ConcurrencyManager is a manager of concurrency.
type ConcurrencyManager struct {
	lt    *LockTable
	locks map[domain.Block]LockType
}

// NewConcurrencyManager constructs a ConcurrencyManager.
func NewConcurrencyManager(lt *LockTable) *ConcurrencyManager {
	return &ConcurrencyManager{
		lt:    lt,
		locks: make(map[domain.Block]LockType),
	}
}

// SLock takes shared lock.
func (conMgr *ConcurrencyManager) SLock(blk domain.Block) error {
	if _, ok := conMgr.locks[blk]; !ok {
		if err := conMgr.lt.SLock(blk); err != nil {
			return err
		}
		conMgr.locks[blk] = Shared
	}

	return nil
}

// XLock takes exclusive lock.
func (conMgr *ConcurrencyManager) XLock(blk domain.Block) error {
	if _, ok := conMgr.locks[blk]; !ok {
		if err := conMgr.lt.XLock(blk); err != nil {
			return err
		}
		conMgr.locks[blk] = Exclusive
	}

	return nil
}

// Release releases all taken locks.
func (conMgr *ConcurrencyManager) Release() {
	for blk, typ := range conMgr.locks {
		switch typ {
		case Shared:
			conMgr.lt.SUnlock(blk)
		case Exclusive:
			conMgr.lt.XUnlock(blk)
		}
	}
	conMgr.locks = make(map[domain.Block]LockType)
}
