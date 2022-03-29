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

		config := infra.NewConfig(dir, 400, "logfile")
		infra.InitServer(config)

		fileMgr, err := file.NewManager(exp, config)
		require.NoError(t, err)

		require.NoError(t, err)
		logMgr, err := log.NewManager(fileMgr, config)
		require.NoError(t, err)

		createRecord(logMgr, 1, 35)
		actual, err := iteratorRecords(logMgr)
		require.NoError(t, err)
		expected := expectedRecords(1, 20)
		require.Equal(t, expected, actual)

		createRecord(logMgr, 36, 70)
		actual, err = iteratorRecords(logMgr)
		require.NoError(t, err)

		expected = expectedRecords(1, 58)
		require.Equal(t, expected, actual)

		err = logMgr.FlushByLSN(65)
		require.NoError(t, err)
		actual, err = iteratorRecords(logMgr)
		require.NoError(t, err)

		expected = expectedRecords(1, 70)
		require.Equal(t, expected, actual)

		err = exp.RemoveAll(dir)
		require.NoError(t, err)
	})
}

func iteratorRecords(logMgr service.LogManager) ([]string, error) {
	strs := make([]string, 0, 400)
	it, _ := logMgr.Iterator()
	for it.HasNext() {
		v, err := it.Next()
		if err != nil {
			return nil, err
		}

		bb := bytes.NewBufferBytes(v)
		page := core.NewPage(bb)
		s, err := page.GetString(0)
		if err != nil {
			return nil, err
		}

		x, err := page.GetUint32(int64(len(s) + 4))
		if err != nil {
			return nil, err
		}

		strs = append(strs, fmt.Sprintf("%v %v", s, x))
	}

	return strs, nil
}

func createRecord(logMgr service.LogManager, start, end int) error {
	for i := start; i <= end; i++ {
		record, err := createLogRecord(fmt.Sprintf("record%d", i), uint32(i+100))
		if err != nil {
			return err
		}

		err = logMgr.AppendRecord(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func createLogRecord(s string, n uint32) ([]byte, error) {
	bufferSize := len(s) + core.Uint32Length*2
	bb := bytes.NewBuffer(bufferSize)
	page := core.NewPage(bb)
	err := page.SetString(0, s)
	if err != nil {
		return nil, err
	}

	err = page.SetUint32(int64(len(s)+core.Uint32Length), n)
	if err != nil {
		return nil, err
	}

	return page.GetBufferBytes(), nil
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
