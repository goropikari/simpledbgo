package hash_test

import (
	"testing"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/index/hash"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestHashIndex(t *testing.T) {
	const (
		idxName   = "idx_id"
		blockSize = 4096
		numBuf    = 10
	)

	t.Run("test hash index", func(t *testing.T) {
		cr := fake.NewTransactionCreater(blockSize, numBuf)
		defer cr.Finish()
		txn := cr.NewTxn()

		var err error

		sch := domain.NewSchema()
		sch.AddInt32Field("id")
		sch.AddInt32Field("block")
		sch.AddInt32Field("dataval")

		layout := domain.NewLayout(sch)

		idx := hash.NewIndex(txn, idxName, layout)

		skey1 := domain.NewConstant(domain.Int32FieldType, fake.RandInt32())
		skey2 := domain.NewConstant(domain.Int32FieldType, fake.RandInt32())

		blkNum1 := domain.BlockNumber(2)
		slotID1 := domain.SlotID(10)
		rid1 := domain.NewRecordID(blkNum1, slotID1)

		blkNum2 := domain.BlockNumber(1)
		slotID2 := domain.SlotID(123)
		rid2 := domain.NewRecordID(blkNum2, slotID2)

		err = idx.Insert(skey1, rid1)
		require.NoError(t, err)

		err = idx.Insert(skey2, rid2)
		require.NoError(t, err)

		err = idx.BeforeFirst(skey1)
		require.NoError(t, err)
		found := idx.HasNext()
		require.True(t, found)
		gotRID1, err := idx.GetDataRecordID()
		require.NoError(t, err)
		require.Equal(t, gotRID1, rid1)

		err = idx.BeforeFirst(skey2)
		require.NoError(t, err)
		found2 := idx.HasNext()
		require.True(t, found2)
		gotRID2, err := idx.GetDataRecordID()
		require.NoError(t, err)
		require.Equal(t, gotRID2, rid2)

		found3 := idx.HasNext()
		require.False(t, found3)
		require.NoError(t, idx.Err())

		err = idx.Delete(skey1, rid1)
		require.NoError(t, err)

		err = idx.BeforeFirst(skey1)
		require.NoError(t, err)
		found4 := idx.HasNext()
		require.False(t, found4)
		_, err = idx.GetDataRecordID()
		require.Error(t, err)

		idx.Close()
		err = txn.Commit()
		require.NoError(t, err)

		fac := hash.NewIndexDriver()
		_, cal := fac.Create()

		cost := cal.Calculate(202, fake.RandInt())
		require.Equal(t, 2, cost)

		err = idx.Err()
		require.NoError(t, err)
	})
}
