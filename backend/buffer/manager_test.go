package buffer_test

import (
	"os"
	"testing"

	"github.com/goropikari/simpledb_go/backend/buffer"
	"github.com/goropikari/simpledb_go/backend/core"
	"github.com/goropikari/simpledb_go/backend/file"
	"github.com/goropikari/simpledb_go/backend/log"
	"github.com/stretchr/testify/require"
)

func TestBufferManager(t *testing.T) {
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

	t.Run("test buffer manager", func(t *testing.T) {
		bufs := make([]*buffer.Buffer, 6)
		bufs[0], err = bm.Pin(core.NewBlock(fileName, 0))
		require.NoError(t, err)
		bufs[1], err = bm.Pin(core.NewBlock(fileName, 1))
		require.NoError(t, err)
		bufs[2], err = bm.Pin(core.NewBlock(fileName, 2))
		require.NoError(t, err)
		err = bm.Unpin(bufs[1])
		require.NoError(t, err)
		bufs[1] = nil
		bufs[3], err = bm.Pin(core.NewBlock(fileName, 0)) // block 0 pinned twice
		require.NoError(t, err)
		bufs[4], err = bm.Pin(core.NewBlock(fileName, 1)) // block 1 repinned
		require.NoError(t, err)

		require.Equal(t, 0, bm.Available())

		bufs[5], err = bm.Pin(core.NewBlock(fileName, 3))
		require.EqualError(t, err, "timeout exceeded")

		err = bm.Unpin(bufs[2])
		require.NoError(t, err)
		bufs[2] = nil
		require.Equal(t, 1, bm.Available())

		bufs[5], err = bm.Pin(core.NewBlock(fileName, 3))
		require.NoError(t, err)
	})
}
