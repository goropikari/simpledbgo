package metadata_test

import (
	"fmt"
	"testing"

	"github.com/goropikari/simpledbgo/index/hash"
	"github.com/goropikari/simpledbgo/metadata"
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestMetadataManager(t *testing.T) {
	const (
		blockSize = 400
		numBuf    = 8
	)

	cr := fake.NewTransactionCreater(blockSize, numBuf)
	defer cr.Finish()

	txn := cr.NewTxn()

	fac := hash.NewIndexFactory()

	metaMgr, err := metadata.CreateManager(fac, txn)
	require.NoError(t, err)

	sch := domain.NewSchema()
	sch.AddInt32Field("A")
	sch.AddStringField("B", 9)

	// table metadata
	metaMgr.CreateTable("MyTable", sch, txn)
	layout, err := metaMgr.GetTableLayout("MyTable", txn)
	require.NoError(t, err)
	sch2 := layout.Schema()

	types := make([]domain.FieldType, 0)
	for _, fld := range sch2.Fields() {
		types = append(types, sch.Type(fld))
	}
	require.Equal(t, []domain.FieldType{domain.FInt32, domain.FString}, types)
	require.Equal(t, int64(21), layout.SlotSize())

	// Statistics Metadata
	tbl, err := domain.NewTable(txn, "MyTable", layout)
	require.NoError(t, err)
	for i := 0; i < 50; i++ {
		err = tbl.AdvanceNextInsertSlotID()
		require.NoError(t, err)

		err = tbl.SetInt32("A", int32(i))
		require.NoError(t, err)

		err = tbl.SetString("B", fmt.Sprintf("rec%v", i))
		require.NoError(t, err)
	}

	si, err := metaMgr.GetStatInfo("MyTable", layout, txn)
	require.NoError(t, err)

	require.Equal(t, 3, si.EstNumBlocks())
	require.Equal(t, 50, si.EstNumRecord())
	require.Equal(t, 1+50/3, si.EstDistinctVals("A"))
	require.Equal(t, 1+50/3, si.EstDistinctVals("B"))

	// view manager
	viewDef := domain.NewViewDef("SELECT B FROM MyTable WHERE A = 1")
	err = metaMgr.CreateView("viewA", viewDef, txn)
	require.NoError(t, err)
	gotViewDef, err := metaMgr.GetViewDef("viewA", txn)
	require.NoError(t, err)
	require.Equal(t, viewDef, gotViewDef)

	// index metadata
	err = metaMgr.CreateIndex("indexA", "MyTable", "A", txn)
	require.NoError(t, err)
	err = metaMgr.CreateIndex("indexB", "MyTable", "B", txn)
	require.NoError(t, err)
	idxMap, err := metaMgr.GetIndexInfo("MyTable", txn)

	ii, found := idxMap["A"]
	require.True(t, found)
	require.Equal(t, 0, ii.EstBlockAccessed())
	require.Equal(t, 2, ii.EstNumRecord())
	require.Equal(t, 1, ii.EstDistinctVals("A"))
	require.Equal(t, 17, ii.EstDistinctVals("B"))

	ii2, found := idxMap["B"]
	require.True(t, found)
	require.Equal(t, 0, ii2.EstBlockAccessed())
	require.Equal(t, 2, ii2.EstNumRecord())
	require.Equal(t, 17, ii2.EstDistinctVals("A"))
	require.Equal(t, 1, ii2.EstDistinctVals("B"))

	err = txn.Commit()
	require.NoError(t, err)
}
