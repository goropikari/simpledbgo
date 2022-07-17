package tx_test

import (
	"testing"

	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/goropikari/simpledbgo/tx"
	"github.com/stretchr/testify/require"
)

func TestConcurrencyManager(t *testing.T) {
	t.Run("lock", func(t *testing.T) {
		blk1 := fake.Block()
		blk2 := fake.Block()
		blk3 := fake.Block()

		cfg := tx.LockTableConfig{LockTimeoutMillisecond: 100}
		lt := tx.NewLockTable(cfg)
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
