package plan_test

import (
	"fmt"
	"testing"

	"github.com/goropikari/simpledbgo/index/hash"
	"github.com/goropikari/simpledbgo/metadata"
	"github.com/goropikari/simpledbgo/plan"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestExecutor_select_table(t *testing.T) {
	const (
		blockSize = 400
		numBuf    = 8
	)

	cr := fake.NewTransactionCreater(blockSize, numBuf)
	defer cr.Finish()

	txn := cr.NewTxn()
	fac := hash.NewIndexDriver()
	mmgr, err := metadata.CreateManager(fac, txn)
	require.NoError(t, err)
	err = txn.Commit()
	require.NoError(t, err)

	t.Run("test Executor", func(t *testing.T) {
		qp := plan.NewBasicQueryPlanner(mmgr)
		ue := plan.NewBasicUpdatePlanner(mmgr)
		pe := plan.NewExecutor(qp, ue)

		txn := cr.NewTxn()
		cmd := "create table T1(A int, B varchar(9))"
		x, err := pe.ExecuteUpdate(cmd, txn)
		require.NoError(t, err)
		require.Equal(t, 0, x)

		n := 200
		for i := 0; i < n; i++ {
			cmd := fmt.Sprintf("insert into T1(A, B) values (%v, 'rec%v')", i, i)
			pe.ExecuteUpdate(cmd, txn)
		}

		qry := "select B from T1 where A=10"
		p, err := pe.CreateQueryPlan(qry, txn)
		require.NoError(t, err)
		s, err := p.Open()
		require.NoError(t, err)

		actual := make([]string, 0)
		for s.HasNext() {
			b, err := s.GetString("b")
			require.NoError(t, err)
			actual = append(actual, b)
		}
		require.NoError(t, s.Err())

		err = txn.Commit()
		require.NoError(t, err)

		expected := []string{"rec10"}
		require.Equal(t, expected, actual)
	})
}

func TestExecutor_select_multi_table(t *testing.T) {
	const (
		blockSize = 400
		numBuf    = 10
	)

	cr := fake.NewTransactionCreater(blockSize, numBuf)
	defer cr.Finish()

	txn := cr.NewTxn()
	fac := hash.NewIndexDriver()
	mmgr, err := metadata.CreateManager(fac, txn)
	require.NoError(t, err)
	err = txn.Commit()
	require.NoError(t, err)

	t.Run("test Executor", func(t *testing.T) {
		qp := plan.NewBasicQueryPlanner(mmgr)
		ue := plan.NewBasicUpdatePlanner(mmgr)
		pe := plan.NewExecutor(qp, ue)

		txn := cr.NewTxn()
		cmd1 := "create table T1(A int, B varchar(9))"
		x, err := pe.ExecuteUpdate(cmd1, txn)
		require.NoError(t, err)
		require.Equal(t, 0, x)

		n := 50
		for i := 0; i < n; i++ {
			cmd := fmt.Sprintf("insert into T1(A, B) values (%v, 'rec%v')", i, i)
			pe.ExecuteUpdate(cmd, txn)
		}

		cmd2 := "create table T2(C int, D varchar(9))"
		y, err := pe.ExecuteUpdate(cmd2, txn)
		require.NoError(t, err)
		require.Equal(t, 0, y)

		for i := 0; i < n; i++ {
			cmd := fmt.Sprintf("insert into T2(D, C) values ('ddd%v', %v)", i, i)
			pe.ExecuteUpdate(cmd, txn)
		}

		qry := " Select D,B from T1, T2 where A=C"
		p, err := pe.CreateQueryPlan(qry, txn)
		require.NoError(t, err)
		s, err := p.Open()
		require.NoError(t, err)

		actual := make([][]string, 0)
		for s.HasNext() {
			b, err := s.GetString("b")
			require.NoError(t, err)
			c, err := s.GetString("d")
			require.NoError(t, err)
			actual = append(actual, []string{b, c})
		}
		require.NoError(t, s.Err())

		err = txn.Commit()
		require.NoError(t, err)

		expected := make([][]string, 0)
		for i := 0; i < n; i++ {
			expected = append(expected, []string{
				fmt.Sprintf("rec%v", i),
				fmt.Sprintf("ddd%v", i),
			})
		}
		require.Equal(t, n, len(actual))
		require.Equal(t, expected, actual)
	})
}

func TestExecutor_select_multi_table_better(t *testing.T) {
	const (
		blockSize = 400
		numBuf    = 10
	)

	cr := fake.NewTransactionCreater(blockSize, numBuf)
	defer cr.Finish()

	txn := cr.NewTxn()
	fac := hash.NewIndexDriver()
	mmgr, err := metadata.CreateManager(fac, txn)
	require.NoError(t, err)
	err = txn.Commit()
	require.NoError(t, err)

	t.Run("test Executor", func(t *testing.T) {
		qp := plan.NewBetterQueryPlanner(mmgr)
		ue := plan.NewBasicUpdatePlanner(mmgr)
		pe := plan.NewExecutor(qp, ue)

		txn := cr.NewTxn()
		cmd1 := "create table T1(A int, B varchar(9))"
		x, err := pe.ExecuteUpdate(cmd1, txn)
		require.NoError(t, err)
		require.Equal(t, 0, x)

		n := 50
		for i := 0; i < n; i++ {
			cmd := fmt.Sprintf("insert into T1(A, B) values (%v, 'rec%v')", i, i)
			pe.ExecuteUpdate(cmd, txn)
		}

		cmd2 := "create table T2(C int, D varchar(9))"
		y, err := pe.ExecuteUpdate(cmd2, txn)
		require.NoError(t, err)
		require.Equal(t, 0, y)

		for i := 0; i < n; i++ {
			cmd := fmt.Sprintf("insert into T2(D, C) values ('ddd%v', %v)", i, i)
			pe.ExecuteUpdate(cmd, txn)
		}

		qry := " Select D,B from T1, T2 where A=C"
		p, err := pe.CreateQueryPlan(qry, txn)
		require.NoError(t, err)
		s, err := p.Open()
		require.NoError(t, err)

		actual := make([][]string, 0)
		for s.HasNext() {
			b, err := s.GetString("b")
			require.NoError(t, err)
			c, err := s.GetString("d")
			require.NoError(t, err)
			actual = append(actual, []string{b, c})
		}
		require.NoError(t, s.Err())

		err = txn.Commit()
		require.NoError(t, err)

		expected := make([][]string, 0)
		for i := 0; i < n; i++ {
			expected = append(expected, []string{
				fmt.Sprintf("rec%v", i),
				fmt.Sprintf("ddd%v", i),
			})
		}
		require.Equal(t, n, len(actual))
		require.Equal(t, expected, actual)
	})
}

func TestExecutor_update_table_without_predicate(t *testing.T) {
	const (
		blockSize = 400
		numBuf    = 8
	)

	cr := fake.NewTransactionCreater(blockSize, numBuf)
	defer cr.Finish()

	txn := cr.NewTxn()
	fac := hash.NewIndexDriver()
	mmgr, err := metadata.CreateManager(fac, txn)
	require.NoError(t, err)
	err = txn.Commit()
	require.NoError(t, err)

	t.Run("test Executor", func(t *testing.T) {
		qp := plan.NewBasicQueryPlanner(mmgr)
		ue := plan.NewBasicUpdatePlanner(mmgr)
		pe := plan.NewExecutor(qp, ue)

		txn := cr.NewTxn()
		cmd := "create table T1(A int, B varchar(9))"
		x, err := pe.ExecuteUpdate(cmd, txn)
		require.NoError(t, err)
		require.Equal(t, 0, x)

		n := 200
		for i := 0; i < n; i++ {
			x := 1000
			if i%2 == 0 {
				x = i
			}
			cmd := fmt.Sprintf("insert into T1(A, B) values (%v, 'rec%v')", x, i)
			pe.ExecuteUpdate(cmd, txn)
		}

		qry := "update T1 set B='foo'"
		x2, err := pe.ExecuteUpdate(qry, txn)
		require.NoError(t, err)
		require.Equal(t, 200, x2)

		qry2 := "select A, B from T1"
		p, err := pe.CreateQueryPlan(qry2, txn)
		require.NoError(t, err)
		s, err := p.Open()
		require.NoError(t, err)

		actual := make([][]any, 0)
		for s.HasNext() {
			a, err := s.GetInt32("a")
			require.NoError(t, err)
			b, err := s.GetString("b")
			require.NoError(t, err)
			actual = append(actual, []any{a, b})
		}
		require.NoError(t, s.Err())

		err = txn.Commit()
		require.NoError(t, err)

		expected := make([][]any, 0)
		for i := 0; i < n; i++ {
			x := 1000
			if i%2 == 0 {
				x = i
			}
			expected = append(expected, []any{int32(x), "foo"})
		}
		require.Equal(t, n, len(actual))
		require.Equal(t, expected, actual)
	})
}

func TestExecutor_update_table_with_predicate(t *testing.T) {
	const (
		blockSize = 400
		numBuf    = 8
	)

	cr := fake.NewTransactionCreater(blockSize, numBuf)
	defer cr.Finish()

	txn := cr.NewTxn()
	fac := hash.NewIndexDriver()
	mmgr, err := metadata.CreateManager(fac, txn)
	require.NoError(t, err)
	err = txn.Commit()
	require.NoError(t, err)

	t.Run("test Executor", func(t *testing.T) {
		qp := plan.NewBasicQueryPlanner(mmgr)
		ue := plan.NewBasicUpdatePlanner(mmgr)
		pe := plan.NewExecutor(qp, ue)

		txn := cr.NewTxn()
		cmd := "create table T1(A int, B varchar(9))"
		x, err := pe.ExecuteUpdate(cmd, txn)
		require.NoError(t, err)
		require.Equal(t, 0, x)

		n := 200
		for i := 0; i < n; i++ {
			x := 1000
			if i%2 == 0 {
				x = i
			}
			cmd := fmt.Sprintf("insert into T1(A, B) values (%v, 'rec%v')", x, i)
			pe.ExecuteUpdate(cmd, txn)
		}

		qry := "update T1 set A=-1 where A=1000"
		x2, err := pe.ExecuteUpdate(qry, txn)
		require.NoError(t, err)
		require.Equal(t, 100, x2)

		qry2 := "select A, B from T1"
		p, err := pe.CreateQueryPlan(qry2, txn)
		require.NoError(t, err)
		s, err := p.Open()
		require.NoError(t, err)

		actual := make([][]any, 0)
		for s.HasNext() {
			a, err := s.GetInt32("a")
			require.NoError(t, err)
			b, err := s.GetString("b")
			require.NoError(t, err)
			actual = append(actual, []any{a, b})
		}
		require.NoError(t, s.Err())

		err = txn.Commit()
		require.NoError(t, err)

		expected := make([][]any, 0)
		for i := 0; i < n; i++ {
			x := -1
			if i%2 == 0 {
				x = i
			}
			expected = append(expected, []any{int32(x), fmt.Sprintf("rec%v", i)})
		}
		require.Equal(t, n, len(actual))
		require.Equal(t, expected, actual)
	})
}

func TestExecutor_update_table_Error(t *testing.T) {
	const (
		blockSize = 400
		numBuf    = 8
	)

	cr := fake.NewTransactionCreater(blockSize, numBuf)
	defer cr.Finish()

	txn := cr.NewTxn()
	fac := hash.NewIndexDriver()
	mmgr, err := metadata.CreateManager(fac, txn)
	require.NoError(t, err)
	err = txn.Commit()
	require.NoError(t, err)

	t.Run("exceed varchar size", func(t *testing.T) {
		qp := plan.NewBasicQueryPlanner(mmgr)
		ue := plan.NewBasicUpdatePlanner(mmgr)
		pe := plan.NewExecutor(qp, ue)

		txn := cr.NewTxn()
		cmd := "create table T1(A int, B varchar(9))"
		x, err := pe.ExecuteUpdate(cmd, txn)
		require.NoError(t, err)
		require.Equal(t, 0, x)

		n := 200
		for i := 0; i < n; i++ {
			x := 1000
			if i%2 == 0 {
				x = i
			}
			cmd := fmt.Sprintf("insert into T1(A, B) values (%v, 'rec%v')", x, i)
			pe.ExecuteUpdate(cmd, txn)
		}

		qry := "update T1 set B='foo_123_bar'"
		_, err = pe.ExecuteUpdate(qry, txn)
		require.Error(t, err)
	})
}

func TestExecutor_delete_record(t *testing.T) {
	const (
		blockSize = 400
		numBuf    = 8
	)

	cr := fake.NewTransactionCreater(blockSize, numBuf)
	defer cr.Finish()

	txn := cr.NewTxn()
	fac := hash.NewIndexDriver()
	mmgr, err := metadata.CreateManager(fac, txn)
	require.NoError(t, err)
	err = txn.Commit()
	require.NoError(t, err)

	t.Run("test Executor", func(t *testing.T) {
		qp := plan.NewBasicQueryPlanner(mmgr)
		ue := plan.NewBasicUpdatePlanner(mmgr)
		pe := plan.NewExecutor(qp, ue)

		txn := cr.NewTxn()
		cmd := "create table T1(A int, B varchar(9))"
		x, err := pe.ExecuteUpdate(cmd, txn)
		require.NoError(t, err)
		require.Equal(t, 0, x)

		n := 200
		for i := 0; i < n; i++ {
			x := 1000
			if i%2 == 0 {
				x = i
			}
			cmd := fmt.Sprintf("insert into T1(A, B) values (%v, 'rec%v')", x, i)
			pe.ExecuteUpdate(cmd, txn)
		}

		qry := "delete from T1 where A=1000"
		x2, err := pe.ExecuteUpdate(qry, txn)
		require.NoError(t, err)
		require.Equal(t, 100, x2)

		qry2 := "select A, B from T1"
		p, err := pe.CreateQueryPlan(qry2, txn)
		require.NoError(t, err)
		s, err := p.Open()
		require.NoError(t, err)

		actual := make([][]any, 0)
		for s.HasNext() {
			a, err := s.GetInt32("a")
			require.NoError(t, err)
			b, err := s.GetString("b")
			require.NoError(t, err)
			actual = append(actual, []any{a, b})
		}
		require.NoError(t, s.Err())

		err = txn.Commit()
		require.NoError(t, err)

		expected := make([][]any, 0)
		for i := 0; i < n; i += 2 {
			expected = append(expected, []any{int32(i), fmt.Sprintf("rec%v", i)})
		}
		require.Equal(t, 100, len(actual))
		require.Equal(t, expected, actual)
	})
}

func TestExecutor_create_view(t *testing.T) {
	const (
		blockSize = 400
		numBuf    = 20
	)

	cr := fake.NewTransactionCreater(blockSize, numBuf)
	defer cr.Finish()

	txn := cr.NewTxn()
	fac := hash.NewIndexDriver()
	mmgr, err := metadata.CreateManager(fac, txn)
	require.NoError(t, err)
	err = txn.Commit()
	require.NoError(t, err)

	t.Run("test Executor", func(t *testing.T) {
		qp := plan.NewBasicQueryPlanner(mmgr)
		ue := plan.NewBasicUpdatePlanner(mmgr)
		pe := plan.NewExecutor(qp, ue)

		txn := cr.NewTxn()
		cmd1 := "create table T1(A int, B varchar(9))"
		x1, err := pe.ExecuteUpdate(cmd1, txn)
		require.NoError(t, err)
		require.Equal(t, 0, x1)

		cmd2 := "create table T2(C int, D varchar(9))"
		x2, err := pe.ExecuteUpdate(cmd2, txn)
		require.NoError(t, err)
		require.Equal(t, 0, x2)

		cmd3 := "create view foo_view as select D from T2"
		x3, err := pe.ExecuteUpdate(cmd3, txn)
		require.NoError(t, err)
		require.Equal(t, 0, x3)

		n := 100
		for i := 0; i < n; i++ {
			cmd1 := fmt.Sprintf("insert into T1(A, B) values (%v, 'rec%v')", i, i)
			pe.ExecuteUpdate(cmd1, txn)
			cmd2 := fmt.Sprintf("insert into T2(C, D) values (%v, 'ddd%v')", i, i)
			pe.ExecuteUpdate(cmd2, txn)
		}

		qry2 := "select A, D from T1, foo_view"
		p, err := pe.CreateQueryPlan(qry2, txn)
		require.NoError(t, err)
		s, err := p.Open()
		require.NoError(t, err)

		actual := make([][]any, 0)
		for s.HasNext() {
			a, err := s.GetInt32("a")
			require.NoError(t, err)
			d, err := s.GetString("d")
			require.NoError(t, err)
			actual = append(actual, []any{a, d})
		}
		require.NoError(t, s.Err())

		err = txn.Commit()
		require.NoError(t, err)

		expected := make([][]any, 0)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				expected = append(expected, []any{int32(i), fmt.Sprintf("ddd%v", j)})
			}
		}
		require.Equal(t, 10000, len(actual))
		require.Equal(t, expected, actual)
	})
}
