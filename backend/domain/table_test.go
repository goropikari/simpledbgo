package domain_test

import (
	"fmt"
	"testing"

	"github.com/goropikari/simpledbgo/backend/domain"
	"github.com/goropikari/simpledbgo/backend/tx"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestTable(t *testing.T) {
	const (
		blockSize = 100
		numBuf    = 2
	)

	dbPath := fake.RandString()
	factory := fake.NewNonDirectBufferManagerFactory(dbPath, blockSize, numBuf)
	fileMgr, logMgr, bufMgr := factory.Create()
	defer factory.Finish()

	ltConfig := tx.NewConfig(1000)
	lt := tx.NewLockTable(ltConfig)
	concurMgr := tx.NewConcurrencyManager(lt)

	gen := tx.NewNumberGenerator()

	t.Run("test commit", func(t *testing.T) {
		txn, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)

		sch := domain.NewSchema()
		sch.AddInt32Field("A")
		sch.AddStringField("B", 9)
		layout := domain.NewLayout(sch)

		table, err := domain.NewTable(txn, "T.tbl", layout)
		require.NoError(t, err)
		for i := 1; i <= 50; i++ {
			err := table.AdvanceNextInsertSlotID()
			require.NoError(t, err)

			n := int32(i)
			err = table.SetInt32("A", n)
			require.NoError(t, err)

			err = table.SetString("B", fmt.Sprintf("rec%v", n))
			require.NoError(t, err)
		}
		table.Close()
		txn.Commit()

		// transaction 2
		txn2, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)

		sch2 := domain.NewSchema()
		sch2.AddInt32Field("A")
		sch2.AddStringField("B", 9)
		layout2 := domain.NewLayout(sch2)

		table2, err := domain.NewTable(txn2, "T.tbl", layout2)
		require.NoError(t, err)

		err = table2.MoveToFirst()
		require.NoError(t, err)
		actual2 := make([]string, 0)
		for {
			found, err := table2.HasNextUsedSlot()
			require.NoError(t, err)
			if !found {
				break
			}

			a, err := table2.GetInt32("A")
			require.NoError(t, err)
			b, err := table2.GetString("B")
			require.NoError(t, err)
			actual2 = append(actual2, fmt.Sprintf("%v %v", a, b))
			if a%2 == 0 {
				err = table2.Delete()
				require.NoError(t, err)
			}
		}
		table2.Close()
		txn2.Commit()

		expected2 := make([]string, 0)
		for i := 1; i <= 50; i++ {
			expected2 = append(expected2, fmt.Sprintf("%v rec%v", i, i))
		}
		require.Equal(t, expected2, actual2)

		// transaction 3
		txn3, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)

		sch3 := domain.NewSchema()
		sch3.AddInt32Field("A")
		sch3.AddStringField("B", 9)
		layout3 := domain.NewLayout(sch3)

		table3, err := domain.NewTable(txn3, "T.tbl", layout3)
		require.NoError(t, err)

		err = table3.MoveToFirst()
		require.NoError(t, err)
		actual3 := make([]string, 0)
		for {
			found, err := table3.HasNextUsedSlot()
			require.NoError(t, err)
			if !found {
				break
			}

			a, err := table3.GetInt32("A")
			require.NoError(t, err)
			b, err := table3.GetString("B")
			require.NoError(t, err)
			actual3 = append(actual3, fmt.Sprintf("%v %v", a, b))
		}
		table3.Close()
		txn3.Commit()

		expected3 := make([]string, 0)
		for i := 1; i <= 50; i += 2 {
			expected3 = append(expected3, fmt.Sprintf("%v rec%v", i, i))
		}
		require.Equal(t, expected3, actual3)
	})
}

func TestTable2(t *testing.T) {
	const (
		blockSize = 100
		numBuf    = 2
	)

	dbPath := fake.RandString()
	factory := fake.NewNonDirectBufferManagerFactory(dbPath, blockSize, numBuf)
	fileMgr, logMgr, bufMgr := factory.Create()
	defer factory.Finish()

	ltConfig := tx.NewConfig(1000)
	lt := tx.NewLockTable(ltConfig)
	concurMgr := tx.NewConcurrencyManager(lt)

	gen := tx.NewNumberGenerator()

	t.Run("test commit", func(t *testing.T) {
		txn, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)

		sch := domain.NewSchema()
		sch.AddInt32Field("A")
		sch.AddStringField("B", 9)
		layout := domain.NewLayout(sch)

		table, err := domain.NewTable(txn, "T.tbl", layout)
		require.NoError(t, err)
		for i := 1; i <= 50; i++ {
			err := table.AdvanceNextInsertSlotID()
			require.NoError(t, err)

			n := int32(i)
			err = table.SetInt32("A", n)
			require.NoError(t, err)

			err = table.SetString("B", fmt.Sprintf("rec%v", n))
			require.NoError(t, err)
		}
		table.Close()
		txn.Commit()

		// transaction 2
		txn2, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)

		sch2 := domain.NewSchema()
		sch2.AddInt32Field("A")
		sch2.AddStringField("B", 9)
		layout2 := domain.NewLayout(sch2)

		table2, err := domain.NewTable(txn2, "T.tbl", layout2)
		require.NoError(t, err)

		err = table2.MoveToFirst()
		require.NoError(t, err)
		for {
			found, err := table2.HasNextUsedSlot()
			require.NoError(t, err)
			if !found {
				break
			}

			a, err := table2.GetInt32("A")
			require.NoError(t, err)
			if a > 25 {
				table2.Delete()
			}
		}
		table2.Close()
		txn2.Commit()

		// transaction 3
		txn3, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)

		sch3 := domain.NewSchema()
		sch3.AddInt32Field("A")
		sch3.AddStringField("B", 9)
		layout3 := domain.NewLayout(sch3)

		table3, err := domain.NewTable(txn3, "T.tbl", layout3)
		require.NoError(t, err)

		err = table3.MoveToFirst()
		require.NoError(t, err)

		err = table3.AdvanceNextInsertSlotID()
		require.NoError(t, err)

		n := int32(100)
		err = table3.SetInt32("A", n)
		require.NoError(t, err)

		err = table3.SetString("B", fmt.Sprintf("rec%v", n))
		require.NoError(t, err)

		err = table3.MoveToFirst()
		require.NoError(t, err)
		actual3 := make([]string, 0)
		for {
			found, err := table3.HasNextUsedSlot()
			require.NoError(t, err)
			if !found {
				break
			}

			a, err := table3.GetInt32("A")
			require.NoError(t, err)
			b, err := table3.GetString("B")
			require.NoError(t, err)
			actual3 = append(actual3, fmt.Sprintf("%v %v", a, b))
		}
		table3.Close()
		txn3.Commit()

		expected3 := make([]string, 0)
		for i := 1; i <= 25; i++ {
			expected3 = append(expected3, fmt.Sprintf("%v rec%v", i, i))
		}
		expected3 = append(expected3, fmt.Sprintf("%v rec%v", n, n))
		require.Equal(t, expected3, actual3)
	})
}
