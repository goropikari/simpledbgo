package log_test

import (
	"fmt"
	golog "log"
	"testing"

	"github.com/goropikari/simpledbgo/common"
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/lib/bytes"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestIterator(t *testing.T) {
	const size = 400

	t.Run("valid request", func(t *testing.T) {
		dbPath := "log_" + fake.RandString()
		logMgrFactory := fake.NewNonDirectLogManagerFactory(dbPath, size)
		defer logMgrFactory.Finish()
		_, logMgr := logMgrFactory.Create()

		createRecord(logMgr, 1, 35)

		actual := actualRecord(logMgr)
		expected := expectedRecords(1, 35)
		require.Equal(t, expected, actual)

		createRecord(logMgr, 36, 70)

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

// record structure
// -----------------------------------------------------------------
// | record length (uint32) | byte expr of string (len(s)) | int32 |
// -----------------------------------------------------------------
func createLogRecord(s string, n int32) []byte {
	stringOffset := int64(0)
	intOffset := bytes.NeededByteLength(s)
	needed := int(intOffset + bytes.NeededByteLength(n))
	bb := bytes.NewBuffer(needed)

	err := bb.SetString(stringOffset, s)
	if err != nil {
		golog.Fatal(err)
	}
	err = bb.SetInt32(intOffset, n)
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

		n, err := bb.GetInt32(int64(len(s) + common.Int32Length))
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
