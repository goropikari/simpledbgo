package domain_test

import (
	"fmt"
	"testing"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/goropikari/simpledbgo/tx"
	"github.com/stretchr/testify/require"
)

func TestTableScan(t *testing.T) {
	const (
		blockSize = 100
		numBuf    = 2
	)

	dbPath := fake.RandString()
	factory := fake.NewNonDirectBufferManagerFactory(dbPath, blockSize, numBuf)
	fileMgr, logMgr, bufMgr := factory.Create()
	defer factory.Finish()

	cfg := tx.LockTableConfig{LockTimeoutMillisecond: 1000}
	lt := tx.NewLockTable(cfg)

	gen := tx.NewNumberGenerator()

	t.Run("test scanning table", func(t *testing.T) {
		// transaction 1
		// insert 50 records
		txn, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, lt, gen)
		require.NoError(t, err)

		sch := domain.NewSchema()
		sch.AddInt32Field("A")
		sch.AddStringField("B", 9)
		layout := domain.NewLayout(sch)

		table, err := domain.NewTableScan(txn, "T.tbl", layout)
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
		// delete even id records
		txn2, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, lt, gen)
		require.NoError(t, err)

		sch2 := domain.NewSchema()
		sch2.AddInt32Field("A")
		sch2.AddStringField("B", 9)
		layout2 := domain.NewLayout(sch2)

		table2, err := domain.NewTableScan(txn2, "T.tbl", layout2)
		require.NoError(t, err)

		err = table2.MoveToFirst()
		require.NoError(t, err)
		actual2 := make([]string, 0)
		for table2.HasNext() {
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
		require.NoError(t, table2.Err())
		table2.Close()
		txn2.Commit()

		expected2 := make([]string, 0)
		for i := 1; i <= 50; i++ {
			expected2 = append(expected2, fmt.Sprintf("%v rec%v", i, i))
		}
		require.Equal(t, expected2, actual2)

		// transaction 3
		// get 25 records
		txn3, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, lt, gen)
		require.NoError(t, err)

		sch3 := domain.NewSchema()
		sch3.AddInt32Field("A")
		sch3.AddStringField("B", 9)
		layout3 := domain.NewLayout(sch3)

		table3, err := domain.NewTableScan(txn3, "T.tbl", layout3)
		require.NoError(t, err)

		err = table3.MoveToFirst()
		require.NoError(t, err)
		actual3 := make([]string, 0)
		for table3.HasNext() {
			a, err := table3.GetInt32("A")
			require.NoError(t, err)
			b, err := table3.GetString("B")
			require.NoError(t, err)
			actual3 = append(actual3, fmt.Sprintf("%v %v", a, b))
		}
		require.NoError(t, table3.Err())
		table3.Close()
		txn3.Commit()

		expected3 := make([]string, 0)
		for i := 1; i <= 50; i += 2 {
			expected3 = append(expected3, fmt.Sprintf("%v rec%v", i, i))
		}
		require.Equal(t, expected3, actual3)
	})
}

func TestTableScan2(t *testing.T) {
	const (
		blockSize = 100
		numBuf    = 2
	)

	dbPath := fake.RandString()
	factory := fake.NewNonDirectBufferManagerFactory(dbPath, blockSize, numBuf)
	fileMgr, logMgr, bufMgr := factory.Create()
	defer factory.Finish()

	cfg := tx.LockTableConfig{LockTimeoutMillisecond: 1000}
	lt := tx.NewLockTable(cfg)

	gen := tx.NewNumberGenerator()

	t.Run("test scannig table", func(t *testing.T) {
		// transaction 1
		// insert 50 records
		txn, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, lt, gen)
		require.NoError(t, err)

		sch := domain.NewSchema()
		sch.AddInt32Field("A")
		sch.AddStringField("B", 9)
		layout := domain.NewLayout(sch)

		table, err := domain.NewTableScan(txn, "T.tbl", layout)
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
		// delete last 25 records
		txn2, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, lt, gen)
		require.NoError(t, err)

		sch2 := domain.NewSchema()
		sch2.AddInt32Field("A")
		sch2.AddStringField("B", 9)
		layout2 := domain.NewLayout(sch2)

		table2, err := domain.NewTableScan(txn2, "T.tbl", layout2)
		require.NoError(t, err)

		err = table2.MoveToFirst()
		require.NoError(t, err)
		for table2.HasNext() {
			a, err := table2.GetInt32("A")
			require.NoError(t, err)
			if a > 25 {
				table2.Delete()
			}
		}
		require.NoError(t, table2.Err())
		table2.Close()
		txn2.Commit()

		// transaction 3
		// insert a record
		// gets 26 records
		txn3, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, lt, gen)
		require.NoError(t, err)

		sch3 := domain.NewSchema()
		sch3.AddInt32Field("A")
		sch3.AddStringField("B", 9)
		layout3 := domain.NewLayout(sch3)

		table3, err := domain.NewTableScan(txn3, "T.tbl", layout3)
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
		for table3.HasNext() {
			a, err := table3.GetInt32("A")
			require.NoError(t, err)
			b, err := table3.GetString("B")
			require.NoError(t, err)
			actual3 = append(actual3, fmt.Sprintf("%v %v", a, b))
		}
		require.NoError(t, table3.Err())
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

func TestProductScan(t *testing.T) {
	const (
		blockSize = 100
		numBuf    = 2
	)

	dbPath := fake.RandString()
	factory := fake.NewNonDirectBufferManagerFactory(dbPath, blockSize, numBuf)
	fileMgr, logMgr, bufMgr := factory.Create()
	defer factory.Finish()

	cfg := tx.LockTableConfig{LockTimeoutMillisecond: 1000}
	lt := tx.NewLockTable(cfg)

	gen := tx.NewNumberGenerator()

	t.Run("test ProductScan", func(t *testing.T) {
		txn, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, lt, gen)
		require.NoError(t, err)

		sch1 := domain.NewSchema()
		sch1.AddInt32Field("A")
		sch1.AddStringField("B", 9)
		layout1 := domain.NewLayout(sch1)
		table1, err := domain.NewTableScan(txn, "T1.tbl", layout1)
		require.NoError(t, err)

		sch2 := domain.NewSchema()
		sch2.AddInt32Field("C")
		sch2.AddStringField("D", 9)
		layout2 := domain.NewLayout(sch2)
		table2, err := domain.NewTableScan(txn, "T2.tbl", layout2)
		require.NoError(t, err)

		n := int32(3)
		for i := int32(1); i <= n; i++ {
			// insert records into table1
			err := table1.AdvanceNextInsertSlotID()
			require.NoError(t, err)

			err = table1.SetInt32("A", i)
			require.NoError(t, err)

			err = table1.SetString("B", fmt.Sprintf("rec%v", i))
			require.NoError(t, err)

			// insert records into table2
			err = table2.AdvanceNextInsertSlotID()
			require.NoError(t, err)

			err = table2.SetInt32("C", i)
			require.NoError(t, err)

			err = table2.SetString("D", fmt.Sprintf("rec%v", i))
			require.NoError(t, err)
		}

		table, err := domain.NewProductScan(table1, table2)
		require.NoError(t, err)

		actual := make([]string, 0)
		for table.HasNext() {
			a, err := table.GetInt32("A")
			require.NoError(t, err)
			b, err := table.GetString("B")
			require.NoError(t, err)
			c, err := table.GetInt32("C")
			require.NoError(t, err)
			d, err := table.GetString("D")
			require.NoError(t, err)

			actual = append(actual, fmt.Sprintf("%v %v %v %v", a, b, c, d))
		}
		require.NoError(t, table.Err())

		table.Close()
		txn.Commit()

		expected := make([]string, 0)
		for i := int32(1); i <= n; i++ {
			for j := int32(1); j <= n; j++ {
				expected = append(expected, fmt.Sprintf("%v rec%v %v rec%v", i, i, j, j))
			}
		}

		require.Equal(t, expected, actual)
	})
}

func TestSelectScan(t *testing.T) {
	const (
		blockSize = 100
		numBuf    = 2
	)

	dbPath := fake.RandString()
	factory := fake.NewNonDirectBufferManagerFactory(dbPath, blockSize, numBuf)
	fileMgr, logMgr, bufMgr := factory.Create()
	defer factory.Finish()

	cfg := tx.LockTableConfig{LockTimeoutMillisecond: 1000}
	lt := tx.NewLockTable(cfg)

	gen := tx.NewNumberGenerator()

	t.Run("test SelectScan", func(t *testing.T) {
		txn, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, lt, gen)
		require.NoError(t, err)

		sch1 := domain.NewSchema()
		sch1.AddInt32Field("A")
		sch1.AddStringField("B", 9)
		layout1 := domain.NewLayout(sch1)
		table1, err := domain.NewTableScan(txn, "T1.tbl", layout1)
		require.NoError(t, err)

		sch2 := domain.NewSchema()
		sch2.AddInt32Field("C")
		sch2.AddStringField("D", 9)
		layout2 := domain.NewLayout(sch2)
		table2, err := domain.NewTableScan(txn, "T2.tbl", layout2)
		require.NoError(t, err)

		n := int32(3)
		for i := int32(1); i <= n; i++ {
			// insert records into table1
			err := table1.AdvanceNextInsertSlotID()
			require.NoError(t, err)

			err = table1.SetInt32("A", i)
			require.NoError(t, err)

			err = table1.SetString("B", fmt.Sprintf("rec%v", i))
			require.NoError(t, err)

			// insert records into table2
			err = table2.AdvanceNextInsertSlotID()
			require.NoError(t, err)

			x := int32(100 + i%2)
			err = table2.SetInt32("C", x)
			require.NoError(t, err)

			err = table2.SetString("D", fmt.Sprintf("rec%v", x))
			require.NoError(t, err)
		}

		productTable, err := domain.NewProductScan(table1, table2)
		require.NoError(t, err)
		pred := domain.NewPredicate([]domain.Term{
			domain.NewTerm(
				domain.NewConstExpression(domain.NewConstant(domain.StringFieldType, "rec2")),
				domain.NewFieldNameExpression("B"),
			),
			domain.NewTerm(
				domain.NewFieldNameExpression("C"),
				domain.NewConstExpression(domain.NewConstant(domain.Int32FieldType, int32(100))),
			),
		})
		table := domain.NewSelectScan(productTable, pred)

		actual := make([]string, 0)
		for table.HasNext() {
			a, err := table.GetInt32("A")
			require.NoError(t, err)
			b, err := table.GetString("B")
			require.NoError(t, err)
			c, err := table.GetInt32("C")
			require.NoError(t, err)
			d, err := table.GetString("D")
			require.NoError(t, err)

			actual = append(actual, fmt.Sprintf("%v %v %v %v", a, b, c, d))
		}
		require.NoError(t, table.Err())

		table.Close()
		txn.Commit()

		expected := []string{"2 rec2 100 rec100"}

		require.Equal(t, expected, actual)
	})
}

func TestSelectScan_update_column(t *testing.T) {
	const (
		blockSize = 100
		numBuf    = 2
	)

	dbPath := fake.RandString()
	factory := fake.NewNonDirectBufferManagerFactory(dbPath, blockSize, numBuf)
	fileMgr, logMgr, bufMgr := factory.Create()
	defer factory.Finish()

	cfg := tx.LockTableConfig{LockTimeoutMillisecond: 1000}
	lt := tx.NewLockTable(cfg)

	gen := tx.NewNumberGenerator()

	t.Run("test SelectScan update column", func(t *testing.T) {
		txn, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, lt, gen)
		require.NoError(t, err)

		sch := domain.NewSchema()
		sch.AddInt32Field("A")
		sch.AddStringField("B", 9)
		layout := domain.NewLayout(sch)

		table, err := domain.NewTableScan(txn, "T.tbl", layout)
		require.NoError(t, err)
		for i := 1; i <= 50; i++ {
			err := table.AdvanceNextInsertSlotID()
			require.NoError(t, err)

			var n int32
			if i%2 == 0 {
				n = 0
			} else {
				n = int32(i)
			}
			err = table.SetInt32("A", n)
			require.NoError(t, err)

			err = table.SetString("B", fmt.Sprintf("rec%v", n))
			require.NoError(t, err)
		}
		require.NoError(t, table.Err())

		pred := domain.NewPredicate([]domain.Term{
			domain.NewTerm(
				domain.NewConstExpression(domain.NewConstant(domain.Int32FieldType, int32(0))),
				domain.NewFieldNameExpression("A"),
			),
		})
		utbl := domain.NewSelectScan(table, pred)

		err = utbl.MoveToFirst()
		require.NoError(t, err)
		for utbl.HasNext() {
			utbl.SetVal("A", domain.NewConstant(domain.Int32FieldType, int32(10000)))
		}

		err = table.MoveToFirst()
		require.NoError(t, err)
		actual := make([]string, 0)
		for table.HasNext() {
			a, err := table.GetInt32("A")
			require.NoError(t, err)
			d, err := table.GetString("B")
			require.NoError(t, err)

			actual = append(actual, fmt.Sprintf("%v %v", a, d))
		}
		require.NoError(t, table.Err())

		utbl.Close()
		txn.Commit()

		expected := make([]string, 0)
		for i := 1; i <= 50; i++ {
			if i%2 == 0 {
				expected = append(expected, fmt.Sprintf("%v rec%v", 10000, 0))
			} else {
				expected = append(expected, fmt.Sprintf("%v rec%v", i, i))
			}
		}
		require.Equal(t, expected, actual)
	})
}

func TestSelectScan_delete_record(t *testing.T) {
	const (
		blockSize = 100
		numBuf    = 2
	)

	dbPath := fake.RandString()
	factory := fake.NewNonDirectBufferManagerFactory(dbPath, blockSize, numBuf)
	fileMgr, logMgr, bufMgr := factory.Create()
	defer factory.Finish()

	cfg := tx.LockTableConfig{LockTimeoutMillisecond: 1000}
	lt := tx.NewLockTable(cfg)

	gen := tx.NewNumberGenerator()

	t.Run("test SelectScan delete record", func(t *testing.T) {
		txn, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, lt, gen)
		require.NoError(t, err)

		sch := domain.NewSchema()
		sch.AddInt32Field("A")
		sch.AddStringField("B", 9)
		layout := domain.NewLayout(sch)

		table, err := domain.NewTableScan(txn, "T.tbl", layout)
		require.NoError(t, err)
		for i := 1; i <= 50; i++ {
			err := table.AdvanceNextInsertSlotID()
			require.NoError(t, err)

			var n int32
			if i%2 == 0 {
				n = 0
			} else {
				n = int32(i)
			}
			err = table.SetInt32("A", n)
			require.NoError(t, err)

			err = table.SetString("B", fmt.Sprintf("rec%v", n))
			require.NoError(t, err)
		}
		require.NoError(t, table.Err())

		pred := domain.NewPredicate([]domain.Term{
			domain.NewTerm(
				domain.NewConstExpression(domain.NewConstant(domain.Int32FieldType, int32(0))),
				domain.NewFieldNameExpression("A"),
			),
		})
		utbl := domain.NewSelectScan(table, pred)

		err = utbl.MoveToFirst()
		require.NoError(t, err)
		for utbl.HasNext() {
			err = utbl.Delete()
			require.NoError(t, err)
		}

		err = table.MoveToFirst()
		require.NoError(t, err)
		actual := make([]string, 0)
		for table.HasNext() {
			a, err := table.GetInt32("A")
			require.NoError(t, err)
			d, err := table.GetString("B")
			require.NoError(t, err)

			actual = append(actual, fmt.Sprintf("%v %v", a, d))
		}
		require.NoError(t, table.Err())

		utbl.Close()
		txn.Commit()

		expected := make([]string, 0)
		for i := 1; i <= 50; i += 2 {
			expected = append(expected, fmt.Sprintf("%v rec%v", i, i))
		}
		require.Equal(t, expected, actual)
	})
}

func TestProjectScan(t *testing.T) {
	const (
		blockSize = 100
		numBuf    = 2
	)

	dbPath := fake.RandString()
	factory := fake.NewNonDirectBufferManagerFactory(dbPath, blockSize, numBuf)
	fileMgr, logMgr, bufMgr := factory.Create()
	defer factory.Finish()

	cfg := tx.LockTableConfig{LockTimeoutMillisecond: 1000}
	lt := tx.NewLockTable(cfg)

	gen := tx.NewNumberGenerator()

	t.Run("test ProjectScan", func(t *testing.T) {
		txn, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, lt, gen)
		require.NoError(t, err)

		sch1 := domain.NewSchema()
		sch1.AddInt32Field("A")
		sch1.AddStringField("B", 9)
		layout1 := domain.NewLayout(sch1)
		table1, err := domain.NewTableScan(txn, "T1.tbl", layout1)
		require.NoError(t, err)

		sch2 := domain.NewSchema()
		sch2.AddInt32Field("C")
		sch2.AddStringField("D", 9)
		layout2 := domain.NewLayout(sch2)
		table2, err := domain.NewTableScan(txn, "T2.tbl", layout2)
		require.NoError(t, err)

		n := int32(3)
		for i := int32(1); i <= n; i++ {
			// insert records into table1
			err := table1.AdvanceNextInsertSlotID()
			require.NoError(t, err)

			err = table1.SetInt32("A", i)
			require.NoError(t, err)

			err = table1.SetString("B", fmt.Sprintf("rec%v", i))
			require.NoError(t, err)

			// insert records into table2
			err = table2.AdvanceNextInsertSlotID()
			require.NoError(t, err)

			x := int32(100 + i%2)
			err = table2.SetInt32("C", x)
			require.NoError(t, err)

			err = table2.SetString("D", fmt.Sprintf("rec%v", x))
			require.NoError(t, err)
		}

		productTable, err := domain.NewProductScan(table1, table2)
		require.NoError(t, err)
		pred := domain.NewPredicate([]domain.Term{
			domain.NewTerm(
				domain.NewConstExpression(domain.NewConstant(domain.Int32FieldType, int32(2))),
				domain.NewFieldNameExpression("A"),
			),
			domain.NewTerm(
				domain.NewFieldNameExpression("C"),
				domain.NewConstExpression(domain.NewConstant(domain.Int32FieldType, int32(100))),
			),
		})
		selectTable := domain.NewSelectScan(productTable, pred)
		table := domain.NewProjectScan(selectTable, []domain.FieldName{"A", "D"})

		actual := make([]string, 0)
		for table.HasNext() {
			a, err := table.GetInt32("A")
			require.NoError(t, err)
			d, err := table.GetString("D")
			require.NoError(t, err)

			actual = append(actual, fmt.Sprintf("%v %v", a, d))
		}
		require.NoError(t, table.Err())

		table.Close()
		txn.Commit()

		expected := []string{"2 rec100"}

		require.Equal(t, expected, actual)
	})
}
