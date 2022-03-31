package domain_test

import (
	goos "os"
	"testing"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/goropikari/simpledb_go/os"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestLogManager_AppendNewBlock(t *testing.T) {
	t.Run("test append new block", func(t *testing.T) {
		blockSize, _ := domain.NewBlockSize(4096)
		bsf := bytes.NewDirectSliceCreater()
		pageFactory := domain.NewPageFactory(bsf, blockSize)

		// initialize file manager
		dbPath := "."
		explorer := os.NewDirectIOExplorer(dbPath)
		fileConfig := domain.FileManagerConfig{BlockSize: 4096}
		fileMgr, _ := domain.NewFileManager(explorer, bsf, fileConfig)

		// initialize log manager
		filename := "logfile" + fake.RandString()
		defer goos.Remove(filename)
		logConfig := domain.LogManagerConfig{LogFileName: filename}
		logPage, _ := pageFactory.Create()
		logFileName, _ := domain.NewFileName(logConfig.LogFileName)
		logBlock, _ := fileMgr.LastBlock(logFileName)

		logMgr, err := domain.NewLogManager(fileMgr, logBlock, logPage, logConfig)
		require.NoError(t, err)

		blk, err := logMgr.AppendNewBlock()
		require.NoError(t, err)
		expected := domain.NewBlock(logFileName, blockSize, domain.BlockNumber(0))
		require.Equal(t, expected, blk)

		blk2, err := logMgr.AppendNewBlock()
		require.NoError(t, err)
		expected2 := domain.NewBlock(logFileName, blockSize, domain.BlockNumber(1))
		require.Equal(t, expected2, blk2)
	})
}
