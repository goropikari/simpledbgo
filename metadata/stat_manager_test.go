package metadata_test

import (
	"fmt"
	"testing"

	"github.com/goropikari/simpledbgo/common"
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/metadata"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestStatManager(t *testing.T) {
	const (
		blockSize = 4096
		numBuf    = 10
	)

	t.Run("test stat manager", func(t *testing.T) {
		cr := fake.NewTransactionCreater(blockSize, numBuf)
		defer cr.Finish()

		txn := cr.NewTxn()
		tblMgr, err := metadata.CreateTableManager(txn)
		require.NoError(t, err)
		sch := domain.NewSchema()
		sch.AddInt32Field("A")
		sch.AddStringField("B", 10)

		tblName := domain.TableName(fake.RandString())
		err = tblMgr.CreateTable(tblName, sch, txn)
		require.NoError(t, err)

		tblLayout, err := tblMgr.GetTableLayout(tblName, txn)
		require.NoError(t, err)

		tbl, err := domain.NewTable(txn, tblName, tblLayout)
		require.NoError(t, err)

		n := int((blockSize - common.Int32Length) / tblLayout.SlotSize())
		// n := 187
		for i := 1; i <= n; i++ {
			err = tbl.AdvanceNextInsertSlotID()
			require.NoError(t, err)

			err = tbl.SetInt32("A", int32(i))
			require.NoError(t, err)

			err = tbl.SetString("B", fmt.Sprintf("rec%v", i))
			require.NoError(t, err)
		}
		err = txn.Commit()
		require.NoError(t, err)

		// test stat manager
		txn2 := cr.NewTxn()
		statMgr, err := metadata.NewStatManager(tblMgr, txn2)
		require.NoError(t, err)

		si, err := statMgr.GetStatInfo(tblName, tblLayout, txn2)
		require.NoError(t, err)
		require.Equal(t, domain.NewStatInfo(1, n), si)
		err = txn2.Commit()
		require.NoError(t, err)

		// block size 2
		txn31 := cr.NewTxn()
		tbl3, err := domain.NewTable(txn31, tblName, tblLayout)
		require.NoError(t, err)

		i := n + 1
		err = tbl3.AdvanceNextInsertSlotID()
		require.NoError(t, err)
		err = tbl3.SetInt32("A", int32(i))
		require.NoError(t, err)
		err = tbl3.SetString("B", fmt.Sprintf("rec%v", i))
		require.NoError(t, err)
		err = txn31.Commit()
		require.NoError(t, err)

		txn32 := cr.NewTxn()
		si3, err := statMgr.GetStatInfo(tblName, tblLayout, txn32)
		require.NoError(t, err)
		require.Equal(t, domain.NewStatInfo(1, n), si3)
		txn32.Commit()
		require.NoError(t, err)

		// update stat info
		// GetStatInfo が 100 より呼ばれたときに StatInfo は更新される
		txn4 := cr.NewTxn()
		for i := 0; i < 98; i++ {
			_, err = statMgr.GetStatInfo(tblName, tblLayout, txn4)
			require.NoError(t, err)
		}
		si4, err := statMgr.GetStatInfo(tblName, tblLayout, txn4)
		require.NoError(t, err)
		require.Equal(t, domain.NewStatInfo(2, n+1), si4)
		txn4.Commit()
		require.NoError(t, err)
	})
}

func TestStatManager_Error(t *testing.T) {
	const (
		blockSize = 4096
		numBuf    = 10
	)

	t.Run("test stat manager", func(t *testing.T) {
		cr := fake.NewTransactionCreater(blockSize, numBuf)
		defer cr.Finish()

		txn := cr.NewTxn()
		tblMgr, err := metadata.CreateTableManager(txn)
		require.NoError(t, err)
		sch := domain.NewSchema()
		sch.AddInt32Field("A")
		sch.AddStringField("B", 10)

		tblName := domain.TableName(fake.RandString())
		err = tblMgr.CreateTable(tblName, sch, txn)
		require.NoError(t, err)

		tblLayout, err := tblMgr.GetTableLayout(tblName, txn)
		require.NoError(t, err)
		err = txn.Commit()
		require.NoError(t, err)

		// test stat manager
		txn2 := cr.NewTxn()
		statMgr, err := metadata.NewStatManager(tblMgr, txn2)
		require.NoError(t, err)

		si, err := statMgr.GetStatInfo(tblName, tblLayout, txn2)
		require.NoError(t, err)
		require.Equal(t, domain.NewStatInfo(0, 0), si)
		err = txn2.Commit()
		require.NoError(t, err)

		// create another table
		txn31 := cr.NewTxn()
		sch2 := domain.NewSchema()
		sch2.AddInt32Field("AA")
		sch2.AddStringField("BB", 10)
		tblName2 := domain.TableName(fake.RandString())
		err = tblMgr.CreateTable(tblName2, sch2, txn31)
		require.NoError(t, err)
		err = txn31.Commit()
		require.NoError(t, err)

		txn32 := cr.NewTxn()
		tblLayout2, err := tblMgr.GetTableLayout(tblName2, txn32)
		require.NoError(t, err)
		tbl3, err := domain.NewTable(txn32, tblName2, tblLayout2)
		require.NoError(t, err)
		err = tbl3.AdvanceNextInsertSlotID()
		require.NoError(t, err)
		err = tbl3.SetInt32("AA", 1)
		require.NoError(t, err)
		err = tbl3.SetString("BB", fmt.Sprintf("rec%v", 1))
		require.NoError(t, err)
		err = txn32.Commit()
		require.NoError(t, err)

		txn33 := cr.NewTxn()
		si2, err := statMgr.GetStatInfo(tblName2, tblLayout2, txn33)
		require.NoError(t, err)
		require.Equal(t, domain.NewStatInfo(1, 1), si2)
	})
}
