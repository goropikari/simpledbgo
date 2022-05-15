package tx_test

import (
	"testing"

	"github.com/goropikari/simpledbgo/tx"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestConcurrencyManager(t *testing.T) {
	t.Run("lock", func(t *testing.T) {
		config := tx.NewConfig(100)
		lt := tx.NewLockTable(config)

		blk1 := fake.Block()
		blk2 := fake.Block()
		blk3 := fake.Block()

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
