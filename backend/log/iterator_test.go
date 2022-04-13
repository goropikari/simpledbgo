package log_test

import (
	"fmt"
	golog "log"
	"testing"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/goropikari/simpledb_go/meta"
	"github.com/stretchr/testify/require"
)

func TestIterator(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		const size = 400
		dbPath := fake.RandString()
		logMgrFactory := fake.NewNonDirectLogManagerFactory(dbPath, size)
		defer logMgrFactory.Finish()
		_, logMgr := logMgrFactory.Create()

		createRecord(logMgr, 1, 35)

		actual := actualRecord(logMgr)
		expected := expectedRecords(1, 20)
		require.Equal(t, expected, actual)

		createRecord(logMgr, 36, 70)

		actual2 := actualRecord(logMgr)
		expected2 := expectedRecords(1, 58)
		require.Equal(t, expected2, actual2)

		err := logMgr.FlushLSN(65)
		require.NoError(t, err)
		actual3 := actualRecord(logMgr)
		expected3 := expectedRecords(1, 70)
		require.Equal(t, expected3, actual3)
	})
}

func createRecord(logMgr domain.LogManager, start, end int) {
	for i := start; i <= end; i++ {
		record := createLogRecord(fmt.Sprintf("record%d", i), int32(i+100))

		_, err := logMgr.AppendRecord(record)
		if err != nil {
			golog.Fatal(err)
		}
	}
}

func createLogRecord(s string, n int32) []byte {
	needed := len(s) + meta.Int32Length*2
	bb := bytes.NewBuffer(needed)

	err := bb.SetString(0, s)
	if err != nil {
		golog.Fatal(err)
	}
	err = bb.SetInt32(int64(needed-meta.Int32Length), n)
	if err != nil {
		golog.Fatal(err)
	}

	return bb.GetData()
}

func actualRecord(logMgr domain.LogManager) []string {
	strs := make([]string, 0)

	iter, err := logMgr.Iterator()
	if err != nil {
		golog.Fatal(err)
	}

	for iter.HasNext() {
		nx, err := iter.Next()
		if err != nil {
			golog.Fatal(err)
		}

		bb := bytes.NewBufferBytes(nx)
		s, err := bb.GetString(0)
		if err != nil {
			golog.Fatal(err)
		}

		n, err := bb.GetInt32(int64(len(s) + meta.Int32Length))
		if err != nil {
			golog.Fatal(err)
		}

		strs = append(strs, fmt.Sprintf("%v %v", s, n))
	}

	return strs
}

func expectedRecords(start, end int) []string {
	strs := make([]string, 0)
	for i := start; i <= end; i++ {
		strs = append(strs, fmt.Sprintf("record%v %v", i, i+100))
	}

	n := len(strs)
	for i := 0; i < n/2; i++ {
		strs[i], strs[n-i-1] = strs[n-i-1], strs[i]
	}

	return strs
}
