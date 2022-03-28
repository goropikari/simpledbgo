package log_test

import (
	"fmt"
	"testing"

	"github.com/goropikari/simpledb_go/backend/core"
	"github.com/goropikari/simpledb_go/backend/file"
	"github.com/goropikari/simpledb_go/backend/log"
	"github.com/goropikari/simpledb_go/backend/service"
	"github.com/goropikari/simpledb_go/infra"
	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/goropikari/simpledb_go/lib/os"
	"github.com/stretchr/testify/require"
)

func TestLogIterator(t *testing.T) {
	t.Run("test iterator", func(t *testing.T) {
		exp := os.NewNormalExplorer()

		dir := "test_db_dir"
		defer exp.RemoveAll(dir)

		filename := "log_iterator"
		config := infra.NewConfig(dir, 400, "logfile")
		fileMgr, err := file.NewManager(exp, config)
		require.NoError(t, err)

		fileName, err := core.NewFileName(filename)
		require.NoError(t, err)
		logMgr, err := log.NewManager(fileMgr, log.NewConfig(fileName))
		require.NoError(t, err)

		createRecord(logMgr, 1, 35)
		actual := iteratorRecords(logMgr)
		expected := expectedRecords(1, 20)
		require.Equal(t, expected, actual)

		createRecord(logMgr, 36, 70)
		actual = iteratorRecords(logMgr)
		expected = expectedRecords(1, 58)
		require.Equal(t, expected, actual)

		err = logMgr.FlushByLSN(65)
		require.NoError(t, err)
		actual = iteratorRecords(logMgr)
		expected = expectedRecords(1, 70)
		require.Equal(t, expected, actual)

		err = exp.RemoveAll(dir)
		require.NoError(t, err)
	})
}

func iteratorRecords(logMgr service.LogManager) []string {
	strs := make([]string, 0, 400)
	ch, _ := logMgr.Iterator()
	for v := range ch {
		bb := bytes.NewBufferBytes(v)
		page := core.NewPage(bb)
		s, _ := page.GetString(0)
		x, _ := page.GetUint32(int64(len(s) + 4))
		strs = append(strs, fmt.Sprintf("%v %v", s, x))
	}

	return strs
}

func createRecord(logMgr service.LogManager, start, end int) {
	for i := start; i <= end; i++ {
		record := createLogRecord(fmt.Sprintf("record%d", i), uint32(i+100))
		logMgr.AppendRecord(record)
	}
}

func createLogRecord(s string, n uint32) []byte {
	bufferSize := len(s) + core.Uint32Length*2
	bb := bytes.NewBuffer(bufferSize)
	page := core.NewPage(bb)
	page.SetString(0, s)
	page.SetUint32(int64(len(s)+core.Uint32Length), n)

	return page.GetBufferBytes()
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
