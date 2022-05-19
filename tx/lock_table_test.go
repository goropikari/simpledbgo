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
	t.Run("RRW", func(t *testing.T) {
		config := tx.NewConfig(1000)
		lt := tx.NewLockTable(config)
		blk := domain.NewBlock(domain.FileName(fake.RandString()), domain.BlockNumber(0))

		tryLock := make([]string, 0)
		actualLock := make([]string, 0)
		wg := &sync.WaitGroup{}
		wg.Add(3)

		go func() {
			tryLock = append(tryLock, "read1")
			err := lt.SLock(blk)
			require.NoError(t, err)
			actualLock = append(actualLock, "read1")
			time.Sleep(150 * time.Millisecond)
			lt.SUnlock(blk)
			wg.Done()
		}()

		go func() {
			time.Sleep(50 * time.Millisecond)
			tryLock = append(tryLock, "read2")
			err := lt.SLock(blk)
			require.NoError(t, err)
			actualLock = append(actualLock, "read2")
			lt.SUnlock(blk)
			wg.Done()
		}()

		go func() {
			time.Sleep(100 * time.Millisecond)
			tryLock = append(tryLock, "write")
			err := lt.XLock(blk)
			require.NoError(t, err)
			actualLock = append(actualLock, "write")
			lt.XUnlock(blk)
			wg.Done()
		}()

		wg.Wait()

		require.Equal(t, []string{"read1", "read2", "write"}, tryLock, "try lock not equal")
		require.Equal(t, []string{"read1", "read2", "write"}, actualLock, "actual lock not equal")
	})

	t.Run("RWR", func(t *testing.T) {
		// RWMutex は RLock 取った後に Lock を取ると、最初の RLock が release されるまで
		// 他の RLock も取ることはできない
		// ref: https://pkg.go.dev/sync#RWMutex
		config := tx.NewConfig(100)
		lt := tx.NewLockTable(config)
		blk := domain.NewBlock(domain.FileName(fake.RandString()), domain.BlockNumber(0))

		wg := &sync.WaitGroup{}
		wg.Add(3)

		now := time.Now()
		go func() {
			err := lt.SLock(blk)
			require.NoError(t, err)
			time.Sleep(20 * time.Millisecond)
			lt.SUnlock(blk)
			wg.Done()
		}()

		go func() {
			time.Sleep(10 * time.Millisecond)
			err := lt.XLock(blk)
			time.Sleep(10 * time.Millisecond)
			require.NoError(t, err)
			lt.XUnlock(blk)
			wg.Done()
		}()

		go func() {
			time.Sleep(15 * time.Millisecond)
			err := lt.SLock(blk)
			time.Sleep(10 * time.Millisecond)
			require.NoError(t, err)
			lt.SUnlock(blk)
			wg.Done()
		}()

		wg.Wait()
		require.Greater(t, time.Since(now), 40*time.Millisecond)
	})

	t.Run("RW timeout", func(t *testing.T) {
		config := tx.NewConfig(10)
		lt := tx.NewLockTable(config)
		blk := domain.NewBlock(domain.FileName(fake.RandString()), domain.BlockNumber(0))

		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			err := lt.SLock(blk)
			require.NoError(t, err)
			time.Sleep(20 * time.Millisecond)
			lt.SUnlock(blk)
			wg.Done()
		}()

		go func() {
			time.Sleep(5 * time.Millisecond)
			err := lt.XLock(blk)
			require.Error(t, err)
			wg.Done()
		}()

		wg.Wait()
	})

	t.Run("WR timeout", func(t *testing.T) {
		config := tx.NewConfig(5)
		lt := tx.NewLockTable(config)
		blk := domain.NewBlock(domain.FileName(fake.RandString()), domain.BlockNumber(0))

		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			err := lt.XLock(blk)
			require.NoError(t, err)
			time.Sleep(20 * time.Millisecond)
			lt.XUnlock(blk)
			wg.Done()
		}()

		go func() {
			time.Sleep(5 * time.Millisecond)
			err := lt.SLock(blk)
			require.Error(t, err)
			wg.Done()
		}()

		wg.Wait()
	})
}
