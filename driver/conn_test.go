package driver_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/goropikari/simpledbgo/driver"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestConn(t *testing.T) {
	dbpath := "simpledb_" + fake.RandString()
	t.Setenv("SIMPLEDB_PATH", dbpath)
	defer os.RemoveAll(dbpath)

	db, err := sql.Open("simpledb", "dsn hoge")
	require.NoError(t, err)

	_, err = db.Exec("create table T1(A int, B varchar(9))")
	require.NoError(t, err)

	_, err = db.Exec("create table T1(A int, B varchar(9))")
	require.Error(t, err)

	n := 3
	for i := 0; i < n; i++ {
		cmd := fmt.Sprintf("insert into T1(A, B) values (%v, 'rec%v')", i, i)
		db.Exec(cmd)
	}

	tx, err := db.Begin()
	require.NoError(t, err)

	rows2, err := tx.QueryContext(context.Background(), "select A, B from T1")
	require.NoError(t, err)
	acnum1 := make([]int, 0)
	acstr1 := make([]string, 0)
	for rows2.Next() {
		var a int
		var b string
		err := rows2.Scan(&a, &b)
		require.NoError(t, err)

		// fmt.Printf("hoge %v %v\n", a, b)
		acnum1 = append(acnum1, a)
		acstr1 = append(acstr1, b)
	}
	require.NoError(t, rows2.Err())
	require.NoError(t, tx.Commit())
	require.Equal(t, []int{0, 1, 2}, acnum1)
	require.Equal(t, []string{"rec0", "rec1", "rec2"}, acstr1)

	tx2, err := db.Begin()
	require.NoError(t, err)
	for i := n; i < 2*n; i++ {
		cmd := fmt.Sprintf("insert into T1(A, B) values (%v, 'rec%v')", i, i)
		tx2.Exec(cmd)
	}

	rows3, err := tx2.QueryContext(context.Background(), "select A, B from T1")
	require.NoError(t, err)
	acnum2 := make([]int, 0)
	acstr2 := make([]string, 0)
	for rows3.Next() {
		var a int
		var b string
		err := rows3.Scan(&a, &b)
		require.NoError(t, err)

		// fmt.Printf("piyo %v %v\n", a, b)
		acnum2 = append(acnum2, a)
		acstr2 = append(acstr2, b)
	}
	require.NoError(t, tx2.Rollback())
	require.Equal(t, []int{0, 1, 2, 3, 4, 5}, acnum2)
	require.Equal(t, []string{"rec0", "rec1", "rec2", "rec3", "rec4", "rec5"}, acstr2)

	rows4, err := db.QueryContext(context.Background(), "select A, B from T1")
	require.NoError(t, err)
	acnum3 := make([]int, 0)
	acstr3 := make([]string, 0)
	for rows4.Next() {
		var a int
		var b string
		err = rows4.Scan(&a, &b)
		require.NoError(t, err)

		// fmt.Printf("fuga %v %v\n", a, b)
		acnum3 = append(acnum3, a)
		acstr3 = append(acstr3, b)
	}
	require.Equal(t, []int{0, 1, 2}, acnum3)
	require.Equal(t, []string{"rec0", "rec1", "rec2"}, acstr3)
}
