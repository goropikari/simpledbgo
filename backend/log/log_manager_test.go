package log_test

import (
	goos "os"
	"testing"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/backend/file"
	"github.com/goropikari/simpledb_go/backend/log"
	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/goropikari/simpledb_go/os"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestManager_Flush(t *testing.T) {
	t.Run("flush", func(t *testing.T) {
		const size = 20
		blockSize, _ := domain.NewBlockSize(size)
		bsf := bytes.NewByteSliceCreater()
		pageFactory := domain.NewPageFactory(bsf, blockSize)

		// initialize file manager
		dbPath := "."
		explorer := os.NewNormalExplorer(dbPath)
		fileConfig := file.ManagerConfig{BlockSize: size}
		fileMgr, _ := file.NewManager(explorer, bsf, fileConfig)

		// initialize log manager
		logfile := "logfile_" + fake.RandString()
		defer goos.Remove(logfile)
		logFileName, err := domain.NewFileName(logfile)
		require.NoError(t, err)

		logConfig := log.ManagerConfig{LogFileName: logfile}
		logBlock, logPage, err := log.PrepareManager(fileMgr, pageFactory, logFileName)
		require.NoError(t, err)

		logMgr, err := log.NewManager(fileMgr, logBlock, logPage, logConfig)
		require.NoError(t, err)

		err = logPage.SetInt32(0, 15)
		require.NoError(t, err)

		err = logMgr.Flush()
		require.NoError(t, err)
	})
}

func TestManager_FlushLSN(t *testing.T) {
	var tests = []struct {
		name   string
		latest int32
		saved  int32
		lsn    int32
	}{
		{
			name:   "flush by lsn",
			latest: 2,
			saved:  1,
			lsn:    2,
		},
		{
			name:   "no flush",
			latest: 3,
			saved:  2,
			lsn:    1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const size = 10
			blockSize, _ := domain.NewBlockSize(size)
			bsf := bytes.NewByteSliceCreater()
			pageFactory := domain.NewPageFactory(bsf, blockSize)

			// initialize file manager
			dbPath := "."
			explorer := os.NewNormalExplorer(dbPath)
			fileConfig := file.ManagerConfig{BlockSize: size}
			fileMgr, _ := file.NewManager(explorer, bsf, fileConfig)

			// initialize log manager
			logfile := "logfile_" + fake.RandString()
			defer goos.Remove(logfile)
			logFileName, err := domain.NewFileName(logfile)
			require.NoError(t, err)

			logConfig := log.ManagerConfig{LogFileName: logfile}
			logBlock, logPage, err := log.PrepareManager(fileMgr, pageFactory, logFileName)
			require.NoError(t, err)

			logMgr, err := log.NewManager(fileMgr, logBlock, logPage, logConfig)
			require.NoError(t, err)

			logMgr.SetLastSavedLSN(tt.saved)
			logMgr.SetLatestLSN(tt.latest)

			err = logPage.SetInt32(4, 1)
			require.NoError(t, err)

			err = logMgr.FlushLSN(tt.lsn)
			require.NoError(t, err)
		})
	}
}

func TestManager_AppendRecord(t *testing.T) {
	t.Run("append record", func(t *testing.T) {
		const size = 15
		blockSize, _ := domain.NewBlockSize(size)
		bsf := bytes.NewByteSliceCreater()
		pageFactory := domain.NewPageFactory(bsf, blockSize)

		// initialize file manager
		dbPath := "."
		explorer := os.NewNormalExplorer(dbPath)
		fileConfig := file.ManagerConfig{BlockSize: size}
		fileMgr, _ := file.NewManager(explorer, bsf, fileConfig)

		// initialize log manager
		logfile := "logfile_" + fake.RandString()
		defer goos.Remove(logfile)
		logFileName, err := domain.NewFileName(logfile)
		require.NoError(t, err)

		logConfig := log.ManagerConfig{LogFileName: logfile}
		logBlock, logPage, err := log.PrepareManager(fileMgr, pageFactory, logFileName)
		require.NoError(t, err)

		logMgr, err := log.NewManager(fileMgr, logBlock, logPage, logConfig)
		require.NoError(t, err)

		_, err = logMgr.AppendRecord([]byte("hello"))
		require.NoError(t, err)
		_, err = logMgr.AppendRecord([]byte("world"))
		require.NoError(t, err)
		err = logMgr.Flush()
		require.NoError(t, err)
	})

	t.Run("append record error: too long record", func(t *testing.T) {
		const size = 10
		blockSize, _ := domain.NewBlockSize(size)
		bsf := bytes.NewByteSliceCreater()
		pageFactory := domain.NewPageFactory(bsf, blockSize)

		// initialize file manager
		dbPath := "."
		explorer := os.NewNormalExplorer(dbPath)
		fileConfig := file.ManagerConfig{BlockSize: size}
		fileMgr, _ := file.NewManager(explorer, bsf, fileConfig)

		// initialize log manager
		logfile := "logfile_" + fake.RandString()
		defer goos.Remove(logfile)
		logFileName, err := domain.NewFileName(logfile)
		require.NoError(t, err)

		logConfig := log.ManagerConfig{LogFileName: logfile}
		logBlock, logPage, err := log.PrepareManager(fileMgr, pageFactory, logFileName)
		require.NoError(t, err)

		logMgr, err := log.NewManager(fileMgr, logBlock, logPage, logConfig)
		require.NoError(t, err)

		_, err = logMgr.AppendRecord([]byte("hello"))
		require.Error(t, err)
	})
}

func TestManager_AppendNewBlock(t *testing.T) {
	t.Run("prepare from empty file", func(t *testing.T) {
		const size = 20
		blockSize, _ := domain.NewBlockSize(size)
		bsf := bytes.NewByteSliceCreater()
		pageFactory := domain.NewPageFactory(bsf, blockSize)

		// initialize file manager
		dbPath := "."
		explorer := os.NewNormalExplorer(dbPath)
		fileConfig := file.ManagerConfig{BlockSize: size}
		fileMgr, _ := file.NewManager(explorer, bsf, fileConfig)

		// initialize log manager
		logfile := "logfile_" + fake.RandString()
		defer goos.Remove(logfile)
		logFileName, err := domain.NewFileName(logfile)
		require.NoError(t, err)

		logConfig := log.ManagerConfig{LogFileName: logfile}
		logBlock, logPage, err := log.PrepareManager(fileMgr, pageFactory, logFileName)
		require.NoError(t, err)

		logMgr, err := log.NewManager(fileMgr, logBlock, logPage, logConfig)
		require.NoError(t, err)

		blk0, err := logMgr.AppendNewBlock()
		require.NoError(t, err)
		expected0 := domain.NewBlock(logFileName, blockSize, domain.BlockNumber(1))
		require.Equal(t, expected0, blk0)

		blk1, err := logMgr.AppendNewBlock()
		require.NoError(t, err)
		expected1 := domain.NewBlock(logFileName, blockSize, domain.BlockNumber(2))
		require.Equal(t, expected1, blk1)
	})

	t.Run("prepare from exsting file", func(t *testing.T) {
		const size = 20
		blockSize, _ := domain.NewBlockSize(size)
		bsf := bytes.NewByteSliceCreater()
		pageFactory := domain.NewPageFactory(bsf, blockSize)

		// initialize file manager
		dbPath := "."
		explorer := os.NewNormalExplorer(dbPath)
		fileConfig := file.ManagerConfig{BlockSize: size}
		fileMgr, err := file.NewManager(explorer, bsf, fileConfig)
		require.NoError(t, err)

		// initialize log manager
		logfile := "logfile_" + fake.RandString()
		defer goos.Remove(logfile)
		logFileName, err := domain.NewFileName(logfile)
		require.NoError(t, err)

		_, err = fileMgr.ExtendFile(logFileName)
		require.NoError(t, err)

		_, err = fileMgr.ExtendFile(logFileName)
		require.NoError(t, err)

		blk, _, err := log.PrepareManager(fileMgr, pageFactory, logFileName)
		require.NoError(t, err)

		expected := domain.NewBlock(logFileName, blockSize, domain.BlockNumber(1))
		require.Equal(t, expected, blk)
	})
}
