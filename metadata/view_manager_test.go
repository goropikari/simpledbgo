package metadata_test

import (
	"testing"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/metadata"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestViewManager(t *testing.T) {
	const (
		blockSize = 4096
		numBuf    = 3
		view      = "foo_view"
	)
	query := domain.NewViewDef("SELECT B FROM MyTable where A = 1")

	t.Run("test ViewManager", func(t *testing.T) {
		cr := fake.NewTransactionCreater(blockSize, numBuf)
		defer cr.Finish()
		txn := cr.NewTxn()

		tblMgr, err := metadata.CreateTableManager(txn)
		require.NoError(t, err)

		viewMgr, err := metadata.CreateViewManager(tblMgr, txn)
		require.NoError(t, err)

		err = viewMgr.CreateView(view, query, txn)
		require.NoError(t, err)

		actual, err := viewMgr.GetViewDef(view, txn)
		require.NoError(t, err)
		require.Equal(t, query, actual)
	})
}
