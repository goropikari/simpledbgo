package btree_test

import (
	"testing"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/index/btree"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/stretchr/testify/require"
)

const (
	idxName   = "idx_id"
	blockSize = 400
	numBuf    = 10
)

func TestIndex(t *testing.T) {
	var err error

	cr := fake.NewTransactionCreater(blockSize, numBuf)
	defer cr.Finish()
	txn := cr.NewTxn()

	sch := domain.NewSchema()
	sch.AddInt32Field("id")
	sch.AddInt32Field("block")
	sch.AddInt32Field("dataval")
	layout := domain.NewLayout(sch)
	idx, err := btree.NewIndex(txn, idxName, layout)
	require.NoError(t, err)

	skey1 := domain.NewConstant(domain.Int32FieldType, int32(123))
	blkNum1 := domain.BlockNumber(2)
	slotID1 := domain.SlotID(10)
	rid1 := domain.NewRecordID(blkNum1, slotID1)
	err = idx.Insert(skey1, rid1)
	require.NoError(t, err)

	skey2 := domain.NewConstant(domain.Int32FieldType, int32(-456))
	blkNum2 := domain.BlockNumber(1)
	slotID2 := domain.SlotID(123)
	rid2 := domain.NewRecordID(blkNum2, slotID2)
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
	require.NoError(t, idx.Err())

	idx.Close()
	err = txn.Commit()
	require.NoError(t, err)

	// cal := btree.NewSearchCostCalculator()
	// cost := cal.Calculate(202, fake.RandInt())
	// require.Equal(t, 2, cost)

	err = idx.Err()
	require.NoError(t, err)
}

func TestIndex_Overflow(t *testing.T) {
	var err error

	cr := fake.NewTransactionCreater(blockSize, numBuf)
	defer cr.Finish()
	txn := cr.NewTxn()

	sch := domain.NewSchema()
	sch.AddInt32Field("id")
	sch.AddInt32Field("block")
	sch.AddInt32Field("dataval")
	layout := domain.NewLayout(sch)
	idx, err := btree.NewIndex(txn, idxName, layout)
	require.NoError(t, err)

	nitems := 200

	skey := domain.NewConstant(domain.Int32FieldType, fake.RandInt32())
	rids := make([]domain.RecordID, 0)
	for i := 0; i < nitems; i++ {
		blkNum := domain.BlockNumber(int32(i))
		slotID := domain.SlotID(int32(i))
		rid := domain.NewRecordID(blkNum, slotID)
		err = idx.Insert(skey, rid)
		require.NoError(t, err)
		rids = append(rids, rid)
	}

	err = idx.BeforeFirst(skey)
	require.NoError(t, err)
	for i := 0; i < nitems; i++ {
		found := idx.HasNext()
		require.True(t, found)
	}
	found := idx.HasNext()
	require.False(t, found)

	for i := 0; i < nitems; i++ {
		err = idx.Delete(skey, rids[i])
		require.NoError(t, err)
	}
	err = idx.Delete(skey, rids[0])
	require.NoError(t, err)

	err = idx.BeforeFirst(skey)
	require.NoError(t, err)
	found4 := idx.HasNext()
	require.False(t, found4)
	require.NoError(t, idx.Err())

	idx.Close()
	err = txn.Commit()
	require.NoError(t, err)
}

func TestIndex_SplitLeaf(t *testing.T) {
	var err error

	cr := fake.NewTransactionCreater(blockSize, numBuf)
	defer cr.Finish()
	txn := cr.NewTxn()

	sch := domain.NewSchema()
	sch.AddInt32Field("id")
	sch.AddInt32Field("block")
	sch.AddInt32Field("dataval")
	layout := domain.NewLayout(sch)
	idx, err := btree.NewIndex(txn, idxName, layout)
	require.NoError(t, err)

	nitems := 40

	skeys := make([]domain.Constant, 0)
	for i := 0; i < nitems; i++ {
		skey := domain.NewConstant(domain.Int32FieldType, int32(i))
		skeys = append(skeys, skey)
		blkNum := domain.BlockNumber(int32(i))
		slotID := domain.SlotID(int32(i))
		rid := domain.NewRecordID(blkNum, slotID)
		err = idx.Insert(skey, rid)
		require.NoError(t, err)
	}

	for i := 0; i < nitems; i++ {
		err = idx.BeforeFirst(skeys[i])
		require.NoError(t, err)
		found := idx.HasNext()
		require.True(t, found)
		found = idx.HasNext()
		require.False(t, found)
	}

	idx.Close()
	err = txn.Commit()
	require.NoError(t, err)
}

func TestIndex_SplitLeaf2(t *testing.T) {
	var err error

	cr := fake.NewTransactionCreater(blockSize, numBuf)
	defer cr.Finish()
	txn := cr.NewTxn()

	sch := domain.NewSchema()
	sch.AddInt32Field("id")
	sch.AddInt32Field("block")
	sch.AddInt32Field("dataval")
	layout := domain.NewLayout(sch)
	idx, err := btree.NewIndex(txn, idxName, layout)
	require.NoError(t, err)

	nitems := 40

	skey1 := domain.NewConstant(domain.Int32FieldType, int32(123))
	for i := 0; i < nitems; i++ {
		blkNum := domain.BlockNumber(int32(i))
		slotID := domain.SlotID(int32(i))
		rid := domain.NewRecordID(blkNum, slotID)
		err = idx.Insert(skey1, rid)
		require.NoError(t, err)
	}

	skey2 := domain.NewConstant(domain.Int32FieldType, int32(111))
	blkNum := domain.BlockNumber(nitems)
	slotID := domain.SlotID(nitems)
	rid := domain.NewRecordID(blkNum, slotID)
	err = idx.Insert(skey2, rid)
	require.NoError(t, err)

	idx.Close()
	err = txn.Commit()
	require.NoError(t, err)
}

func TestIndex_SplitLeaf3(t *testing.T) {
	var err error

	cr := fake.NewTransactionCreater(blockSize, numBuf)
	defer cr.Finish()
	txn := cr.NewTxn()

	sch := domain.NewSchema()
	sch.AddInt32Field("id")
	sch.AddInt32Field("block")
	sch.AddInt32Field("dataval")
	layout := domain.NewLayout(sch)
	idx, err := btree.NewIndex(txn, idxName, layout)
	require.NoError(t, err)

	nitems := 40
	skey1 := domain.NewConstant(domain.Int32FieldType, int32(123))
	for i := 0; i < nitems; i++ {
		blkNum := domain.BlockNumber(int32(i))
		slotID := domain.SlotID(int32(i))
		rid := domain.NewRecordID(blkNum, slotID)
		err = idx.Insert(skey1, rid)
		require.NoError(t, err)
	}

	skey2 := domain.NewConstant(domain.Int32FieldType, int32(456))
	for i := nitems; i < 2*nitems; i++ {
		blkNum := domain.BlockNumber(int32(i))
		slotID := domain.SlotID(int32(i))
		rid := domain.NewRecordID(blkNum, slotID)
		err = idx.Insert(skey2, rid)
		require.NoError(t, err)
	}

	err = idx.BeforeFirst(skey1)
	require.NoError(t, err)
	for i := 0; i < nitems; i++ {
		found := idx.HasNext()
		require.True(t, found)
	}
	found := idx.HasNext()
	require.False(t, found)

	err = idx.BeforeFirst(skey2)
	require.NoError(t, err)
	for i := 0; i < nitems; i++ {
		found := idx.HasNext()
		require.True(t, found)
	}
	found = idx.HasNext()
	require.False(t, found)

	idx.Close()
	err = txn.Commit()
	require.NoError(t, err)
}

func TestIndex_SplitLeaf4(t *testing.T) {
	var err error

	blockSize := int32(400)
	numBuf := 100
	cr := fake.NewTransactionCreater(blockSize, numBuf)
	defer cr.Finish()
	txn := cr.NewTxn()

	sch := domain.NewSchema()
	sch.AddInt32Field("id")
	sch.AddInt32Field("block")
	sch.AddInt32Field("dataval")
	layout := domain.NewLayout(sch)
	idx, err := btree.NewIndex(txn, idxName, layout)
	require.NoError(t, err)

	nitems := 1000
	for i := 0; i < nitems; i++ {
		blkNum := domain.BlockNumber(int32(i))
		slotID := domain.SlotID(int32(i))
		rid := domain.NewRecordID(blkNum, slotID)
		skey := domain.NewConstant(domain.Int32FieldType, int32(i))
		err = idx.Insert(skey, rid)
		require.NoError(t, err)
	}

	for i := 0; i < nitems; i++ {
		skey := domain.NewConstant(domain.Int32FieldType, int32(i))
		err = idx.BeforeFirst(skey)
		require.NoError(t, err)
		found := idx.HasNext()
		require.True(t, found)
		found = idx.HasNext()
		require.False(t, found)
	}

	idx.Close()
	err = txn.Commit()
	require.NoError(t, err)
}
