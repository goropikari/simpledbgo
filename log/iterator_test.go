package log_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/goropikari/simpledb_go/bytes"
	"github.com/goropikari/simpledb_go/core"
	"github.com/goropikari/simpledb_go/file"
	"github.com/goropikari/simpledb_go/log"
	"github.com/stretchr/testify/require"
)

func TestLogIterator(t *testing.T) {
	t.Run("test iterator", func(t *testing.T) {
		dir := "test_db_dir"
		filename := "log_iterator"
		config, err := file.NewConfig(dir, 400, false)
		require.NoError(t, err)
		fileMgr, err := file.NewManager(config)
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

		err = os.RemoveAll(dir)
		require.NoError(t, err)
	})
}

func iteratorRecords(logMgr *log.Manager) []string {
	strs := make([]string, 0, 400)
	ch, _ := logMgr.Iterator()
	for v := range ch {
		bb := bytes.NewBufferBytes(v)
		page := file.NewPage(bb)
		s, _ := page.GetString(0)
		x, _ := page.GetUint32(int64(len(s) + 4))
		strs = append(strs, fmt.Sprintf("%v %v", s, x))
	}

	return strs
}

func createRecord(logMgr *log.Manager, start, end int) {
	for i := start; i <= end; i++ {
		record := createLogRecord(fmt.Sprintf("record%d", i), uint32(i+100))
		logMgr.AppendRecord(record)
	}
}

func createLogRecord(s string, n uint32) []byte {
	bufferSize, _ := core.NewBlockSize(len(s) + core.Uint32Length*2)
	bb, _ := bytes.NewBuffer(bufferSize)
	page := file.NewPage(bb)
	page.SetString(0, s)
	page.SetUint32(int64(len(s)+core.Uint32Length), n)

	return page.GetFullBytes()
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
