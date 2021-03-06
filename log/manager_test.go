package log_test

import (
	goos "os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/file"
	"github.com/goropikari/simpledbgo/log"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestManager_NewManager(t *testing.T) {
	const size = 20

	t.Run("valid request: initialize empty file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// initialize file manager
		dbPath := "log_" + fake.RandString()
		fileConfig := file.ManagerConfig{
			DBPath:    dbPath,
			BlockSize: size,
			DirectIO:  false,
		}
		fileMgr, _ := file.NewManager(fileConfig)
		defer goos.RemoveAll(dbPath)

		// initialize log manager
		logfile := "logfile_" + fake.RandString()
		logConfig := log.ManagerConfig{LogFileName: logfile}

		_, err := log.NewManager(fileMgr, logConfig)
		require.NoError(t, err)
	})

	t.Run("valid request: initialize with 2 blocks file", func(t *testing.T) {
		// make dummy file
		dbPath := "log_" + fake.RandString()
		logfile := "logfile_" + fake.RandString()
		goos.MkdirAll(dbPath, goos.ModePerm)
		defer goos.RemoveAll(dbPath)
		f, err := goos.OpenFile(dbPath+"/"+logfile, goos.O_CREATE|goos.O_RDWR, goos.ModePerm)
		require.NoError(t, err)
		_, err = f.Write(make([]byte, size*2))
		require.NoError(t, err)
		err = f.Close()
		require.NoError(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// initialize file manager
		fileConfig := file.ManagerConfig{
			DBPath:    dbPath,
			BlockSize: size,
			DirectIO:  false,
		}
		fileMgr, _ := file.NewManager(fileConfig)

		// initialize log manager
		logConfig := log.ManagerConfig{LogFileName: logfile}
		logMgr, err := log.NewManager(fileMgr, logConfig)
		require.NoError(t, err)
		require.Equal(t, domain.BlockNumber(1), logMgr.CurrentBlock().Number())
	})
}

func TestManager_Flush(t *testing.T) {
	const size = 20

	t.Run("flush", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// initialize file manager
		dbPath := "log_" + fake.RandString()
		fileConfig := file.ManagerConfig{
			DBPath:    dbPath,
			BlockSize: size,
			DirectIO:  false,
		}
		fileMgr, _ := file.NewManager(fileConfig)

		// initialize log manager
		logfile := "logfile_" + fake.RandString()
		defer goos.RemoveAll(dbPath)

		logConfig := log.ManagerConfig{LogFileName: logfile}

		logMgr, err := log.NewManager(fileMgr, logConfig)
		require.NoError(t, err)

		err = logMgr.Flush()
		require.NoError(t, err)
	})
}

func TestManager_FlushLSN(t *testing.T) {
	const size = 20

	var tests = []struct {
		name   string
		latest domain.LSN
		saved  domain.LSN
		lsn    domain.LSN
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// initialize file manager
			dbPath := fake.RandString()
			fileMgrFactory := fake.NewNonDirectFileManagerFactory(dbPath, size)
			defer fileMgrFactory.Finish()
			fileMgr := fileMgrFactory.Create()

			// initialize log manager
			logfile := "logfile_" + fake.RandString()
			defer goos.Remove(logfile)

			logConfig := log.ManagerConfig{LogFileName: logfile}

			logMgr, err := log.NewManager(fileMgr, logConfig)
			require.NoError(t, err)

			logMgr.SetLastSavedLSN(tt.saved)
			logMgr.SetLatestLSN(tt.latest)

			err = logMgr.FlushLSN(tt.lsn)
			require.NoError(t, err)
		})
	}
}

func TestManager_AppendRecord(t *testing.T) {
	t.Run("append record", func(t *testing.T) {
		const size = 15

		dbPath := fake.RandString()
		logMgrFactory := fake.NewNonDirectLogManagerFactory(dbPath, size)
		defer logMgrFactory.Finish()
		_, logMgr := logMgrFactory.Create()

		_, err := logMgr.AppendRecord([]byte("hello"))
		require.NoError(t, err)
		_, err = logMgr.AppendRecord([]byte("world"))
		require.NoError(t, err)
		err = logMgr.Flush()
		require.NoError(t, err)
	})

	t.Run("append record error: too long record", func(t *testing.T) {
		const size = 10

		dbPath := fake.RandString()
		logMgrFactory := fake.NewNonDirectLogManagerFactory(dbPath, size)
		defer logMgrFactory.Finish()
		_, logMgr := logMgrFactory.Create()

		_, err := logMgr.AppendRecord([]byte("hello"))
		require.Error(t, err)
	})
}

func TestManager_AppendNewBlock(t *testing.T) {
	const size = 20

	t.Run("prepare from empty file", func(t *testing.T) {
		dbPath := fake.RandString()
		logMgrFactory := fake.NewNonDirectLogManagerFactory(dbPath, size)
		defer logMgrFactory.Finish()
		_, logMgr := logMgrFactory.Create()

		blk0, err := logMgr.AppendNewBlock()
		require.NoError(t, err)
		require.Equal(t, domain.BlockNumber(1), blk0.Number())

		blk1, err := logMgr.AppendNewBlock()
		require.NoError(t, err)
		require.Equal(t, domain.BlockNumber(2), blk1.Number())
	})

	t.Run("prepare from exsting file", func(t *testing.T) {
		// initialize file manager
		dbPath := fake.RandString()
		fileMgrFactory := fake.NewNonDirectFileManagerFactory(dbPath, size)
		defer fileMgrFactory.Finish()
		fileMgr := fileMgrFactory.Create()

		// initialize log manager
		logfile := "logfile_" + fake.RandString()
		defer goos.Remove(logfile)
		logFileName, err := domain.NewFileName(logfile)
		require.NoError(t, err)

		_, err = fileMgr.ExtendFile(logFileName)
		require.NoError(t, err)

		_, err = fileMgr.ExtendFile(logFileName)
		require.NoError(t, err)

		logConfig := log.ManagerConfig{LogFileName: logfile}

		logMgr, err := log.NewManager(fileMgr, logConfig)
		require.NoError(t, err)

		blk0, err := logMgr.AppendNewBlock()
		require.NoError(t, err)
		expected0 := domain.NewBlock(logFileName, domain.BlockNumber(2))
		require.Equal(t, expected0, blk0)
	})
}
