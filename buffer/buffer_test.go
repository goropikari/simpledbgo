package buffer_test

import (
	"os"
	"testing"

	"github.com/goropikari/simpledb_go/buffer"
	"github.com/goropikari/simpledb_go/core"
	"github.com/goropikari/simpledb_go/file"
	"github.com/goropikari/simpledb_go/log"
	"github.com/stretchr/testify/require"
)

func TestBuffer(t *testing.T) {
	blockSize := 400
	isDirectIO := false
	dir := "test" + core.RandomString()
	logFile := core.FileName("logfile" + core.RandomString())
	fileName := core.FileName("testfile" + core.RandomString())
	defer os.RemoveAll(dir)
	fileConfig, err := file.NewConfig(dir, blockSize, isDirectIO)
	require.NoError(t, err)

	fm, err := file.NewManager(fileConfig)
	require.NoError(t, err)

	logConfig := log.NewConfig(logFile)
	lm, err := log.NewManager(fm, logConfig)
	require.NoError(t, err)

	numBeffer := 3
	bm, err := buffer.NewManager(fm, lm, numBeffer)
	require.NoError(t, err)

	t.Run("test buffer", func(t *testing.T) {
		block := file.NewBlock(fileName, 1)
		buf1, err := bm.Pin(block)
		require.NoError(t, err)

		p := buf1.GetPage()
		n, err := p.GetInt32(80)
		require.NoError(t, err)
		err = p.SetInt32(80, n+1)
		require.NoError(t, err)
		buf1.SetModified(1, 0)
		require.NoError(t, err)
		err = bm.Unpin(buf1)
		require.NoError(t, err)

		// One of these pins will flush buff1 to disk:
		buf2, err := bm.Pin(file.NewBlock(fileName, 2))
		require.NoError(t, err)
		_, err = bm.Pin(file.NewBlock(fileName, 3))
		require.NoError(t, err)
		_, err = bm.Pin(file.NewBlock(fileName, 4))
		require.NoError(t, err)

		err = bm.Unpin(buf2)
		require.NoError(t, err)

		buf1, err = bm.Pin(file.NewBlock(fileName, 1))
		require.NoError(t, err)
		p1 := buf1.GetPage()
		// This modification won't get written to disk.
		err = p1.SetInt32(80, 9999)
		require.NoError(t, err)
		buf1.SetModified(1, 0)
		require.NoError(t, err)
	})
}
