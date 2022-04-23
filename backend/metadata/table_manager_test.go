package metadata_test

import (
	"testing"

	"github.com/goropikari/simpledbgo/backend/domain"
	"github.com/goropikari/simpledbgo/backend/metadata"
	"github.com/goropikari/simpledbgo/backend/record"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestTableManager(t *testing.T) {
	const (
		blockSize = 4096
		numBuf    = 3
	)

	t.Run("test TableManager", func(t *testing.T) {
		cr := fake.NewTransactionCreater(blockSize, numBuf)
		defer cr.Finish()
		txn := cr.NewTxn()

		tblMgr, err := metadata.CreateTableManager(txn)
		require.NoError(t, err)

		sch := record.NewSchema()
		sch.AddInt32Field("A")
		sch.AddStringField("B", 9)

		tblName := domain.FileName(fake.RandString())
		err = tblMgr.CreateTable(tblName, sch, txn)
		require.NoError(t, err)

		layout, err := tblMgr.GetTableLayout(tblName, txn)
		require.NoError(t, err)

		expected := record.NewLayout(sch)
		require.Equal(t, expected, layout)
	})
}
