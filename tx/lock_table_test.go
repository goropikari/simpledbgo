package tx_test

import (
	"sync"
	"testing"
	"time"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/goropikari/simpledbgo/tx"
	"github.com/stretchr/testify/require"
)

func TestLockTable_Lock(t *testing.T) {
	t.Run("W", func(t *testing.T) {
		cfg := tx.LockTableConfig{LockTimeoutMillisecond: 10000}
		lt := tx.NewLockTable(cfg)

		blk := domain.NewBlock(domain.FileName(fake.RandString()), domain.BlockNumber(0))

		tryLock := make([]string, 0)
		actualLock := make([]string, 0)
		wg := &sync.WaitGroup{}
		wg.Add(1)

		go func() {
			tryLock = append(tryLock, "write")
			err := lt.SLock(blk)
			require.NoError(t, err)
			err = lt.XLock(blk)
			require.NoError(t, err)
			actualLock = append(actualLock, "write")
			lt.Unlock(blk)
			wg.Done()
		}()

		wg.Wait()

		require.Equal(t, []string{"write"}, tryLock, "try lock not equal")
		require.Equal(t, []string{"write"}, actualLock, "actual lock not equal")
	})

	t.Run("RR", func(t *testing.T) {
		cfg := tx.LockTableConfig{LockTimeoutMillisecond: 1000}
		lt := tx.NewLockTable(cfg)
		blk := domain.NewBlock(domain.FileName(fake.RandString()), domain.BlockNumber(0))

		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			err := lt.SLock(blk)
			require.NoError(t, err)
			time.Sleep(150 * time.Millisecond)
			lt.Unlock(blk)
			wg.Done()
		}()

		go func() {
			time.Sleep(50 * time.Millisecond)
			err := lt.SLock(blk)
			require.NoError(t, err)
			lt.Unlock(blk)
			wg.Done()
		}()

		wg.Wait()
	})

	t.Run("RW", func(t *testing.T) {
		cfg := tx.LockTableConfig{LockTimeoutMillisecond: 10000}
		lt := tx.NewLockTable(cfg)

		blk := domain.NewBlock(domain.FileName(fake.RandString()), domain.BlockNumber(0))

		tryLock := make([]string, 0)
		actualLock := make([]string, 0)
		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			tryLock = append(tryLock, "read1")
			err := lt.SLock(blk)
			require.NoError(t, err)
			actualLock = append(actualLock, "read1")
			time.Sleep(300 * time.Millisecond)
			lt.Unlock(blk)
			wg.Done()
		}()

		go func() {
			time.Sleep(100 * time.Millisecond)
			tryLock = append(tryLock, "write")
			var err error
			err = lt.SLock(blk)
			require.NoError(t, err)
			err = lt.XLock(blk)
			require.NoError(t, err)
			actualLock = append(actualLock, "write")
			lt.Unlock(blk)
			wg.Done()
		}()

		wg.Wait()

		require.Equal(t, []string{"read1", "write"}, tryLock, "try lock not equal")
		require.Equal(t, []string{"read1", "write"}, actualLock, "actual lock not equal")
	})

	t.Run("RW timeout", func(t *testing.T) {
		cfg := tx.LockTableConfig{LockTimeoutMillisecond: 10}
		lt := tx.NewLockTable(cfg)
		blk := domain.NewBlock(domain.FileName(fake.RandString()), domain.BlockNumber(0))

		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			err := lt.SLock(blk)
			require.NoError(t, err)
			time.Sleep(100 * time.Millisecond)
			lt.Unlock(blk)
			wg.Done()
		}()

		go func() {
			time.Sleep(30 * time.Millisecond)
			err := lt.SLock(blk)
			require.NoError(t, err)
			err = lt.XLock(blk)
			require.Error(t, err)
			wg.Done()
		}()

		wg.Wait()
	})

	t.Run("WR timeout", func(t *testing.T) {
		cfg := tx.LockTableConfig{LockTimeoutMillisecond: 50}
		lt := tx.NewLockTable(cfg)
		blk := domain.NewBlock(domain.FileName(fake.RandString()), domain.BlockNumber(0))

		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			var err error
			err = lt.SLock(blk)
			require.NoError(t, err)
			err = lt.XLock(blk)
			require.NoError(t, err)
			time.Sleep(100 * time.Millisecond)
			lt.Unlock(blk)
			wg.Done()
		}()

		go func() {
			time.Sleep(10 * time.Millisecond)
			err := lt.SLock(blk)
			require.Error(t, err)
			wg.Done()
		}()

		wg.Wait()
	})
}
