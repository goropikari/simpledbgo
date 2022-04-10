package tx_test

import (
	"sync"
	"testing"
	"time"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/backend/tx"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestLockTable_Lock(t *testing.T) {

	t.Run("RRW", func(t *testing.T) {
		config := tx.NewConfig(100)
		lt := tx.NewLockTable(config)
		blk := *domain.NewBlock(domain.FileName(fake.RandString()), domain.BlockSize(10), domain.BlockNumber(0))

		tryLock := make([]string, 0)
		actualLock := make([]string, 0)
		wg := &sync.WaitGroup{}
		wg.Add(3)

		go func() {
			tryLock = append(tryLock, "read1")
			err := lt.SLock(blk)
			require.NoError(t, err)
			actualLock = append(actualLock, "read1")
			time.Sleep(50 * time.Millisecond)
			lt.SUnlock(blk)
			wg.Done()
		}()

		go func() {
			time.Sleep(10 * time.Millisecond)
			tryLock = append(tryLock, "read2")
			err := lt.SLock(blk)
			require.NoError(t, err)
			actualLock = append(actualLock, "read2")
			lt.SUnlock(blk)
			wg.Done()
		}()

		go func() {
			time.Sleep(15 * time.Millisecond)
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
		// writer 優先だから read2 は writer の後になる
		config := tx.NewConfig(100)
		lt := tx.NewLockTable(config)
		blk := *domain.NewBlock(domain.FileName(fake.RandString()), domain.BlockSize(10), domain.BlockNumber(0))

		tryLock := make([]string, 0)
		actualLock := make([]string, 0)
		wg := &sync.WaitGroup{}
		wg.Add(3)

		go func() {
			tryLock = append(tryLock, "read1")
			err := lt.SLock(blk)
			require.NoError(t, err)
			actualLock = append(actualLock, "read1")
			time.Sleep(20 * time.Millisecond)
			lt.SUnlock(blk)
			wg.Done()
		}()

		go func() {
			time.Sleep(10 * time.Millisecond)
			tryLock = append(tryLock, "write")
			err := lt.XLock(blk)
			require.NoError(t, err)
			actualLock = append(actualLock, "write")
			lt.XUnlock(blk)
			wg.Done()
		}()

		go func() {
			time.Sleep(15 * time.Millisecond)
			tryLock = append(tryLock, "read2")
			err := lt.SLock(blk)
			require.NoError(t, err)
			actualLock = append(actualLock, "read2")
			lt.SUnlock(blk)
			wg.Done()
		}()

		wg.Wait()

		require.Equal(t, []string{"read1", "write", "read2"}, tryLock, "try lock not equal")
		require.Equal(t, []string{"read1", "write", "read2"}, actualLock, "actual lock not equal")
	})

	t.Run("RW timeout", func(t *testing.T) {
		config := tx.NewConfig(10)
		lt := tx.NewLockTable(config)
		blk := *domain.NewBlock(domain.FileName(fake.RandString()), domain.BlockSize(10), domain.BlockNumber(0))

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
			lt.XUnlock(blk)
			wg.Done()
		}()

		wg.Wait()
	})

	t.Run("WR timeout", func(t *testing.T) {
		config := tx.NewConfig(10)
		lt := tx.NewLockTable(config)
		blk := *domain.NewBlock(domain.FileName(fake.RandString()), domain.BlockSize(10), domain.BlockNumber(0))

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
			lt.SUnlock(blk)
			wg.Done()
		}()

		wg.Wait()
	})
}
