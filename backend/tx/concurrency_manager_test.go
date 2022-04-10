package tx_test

import (
	"testing"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/backend/tx"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestConcurrencyManager(t *testing.T) {
	t.Run("lock", func(t *testing.T) {
		config := tx.NewConfig(100)
		lt := tx.NewLockTable(config)

		blkSize := domain.BlockSize(100)
		blkNum := domain.BlockNumber(0)
		blk1 := *domain.NewBlock(domain.FileName(fake.RandString()), blkSize, blkNum)
		blk2 := *domain.NewBlock(domain.FileName(fake.RandString()), blkSize, blkNum)
		blk3 := *domain.NewBlock(domain.FileName(fake.RandString()), blkSize, blkNum)

		mgr := tx.NewConcurrencyManager(lt)

		err := mgr.SLock(blk1)
		require.NoError(t, err)
		err = mgr.SLock(blk2)
		require.NoError(t, err)
		err = mgr.XLock(blk3)
		require.NoError(t, err)

		mgr.Release()
	})
}
