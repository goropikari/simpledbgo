package buffer_test

import (
	goos "os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/goropikari/simpledb_go/backend/buffer"
	"github.com/goropikari/simpledb_go/backend/core"
	"github.com/goropikari/simpledb_go/backend/file"
	"github.com/goropikari/simpledb_go/backend/log"
	"github.com/goropikari/simpledb_go/infra"
	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/goropikari/simpledb_go/lib/os"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/goropikari/simpledb_go/testing/mock"
	"github.com/stretchr/testify/require"
)

func TestBufferManager_mock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fm := mock.NewMockFileManager(ctrl)
	lm := mock.NewMockLogManager(ctrl)

	bb := bytes.NewBuffer(100)
	page := core.NewPage(bb)
	fm.EXPECT().PreparePage().Return(page, nil).AnyTimes()
	fm.EXPECT().CopyBlockToPage(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	fm.EXPECT().CopyPageToBlock(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	lm.EXPECT().FlushByLSN(gomock.Any()).Return(nil).AnyTimes()

	numBeffer := 3
	bm := buffer.NewManager(fm, lm, numBeffer)

	t.Run("test buffer", func(t *testing.T) {
		fileName, _ := core.NewFileName("hoge")
		block := core.NewBlock(fileName, 1)
		buf1, err := bm.Pin(block)
		require.NoError(t, err)
		require.Equal(t, numBeffer-1, bm.Available())

		err = bm.Unpin(buf1)
		require.NoError(t, err)
		require.Equal(t, numBeffer, bm.Available())

		// One of these pins will flush buff1 to disk:
		buf2, err := bm.Pin(core.NewBlock(fileName, 2))
		require.NoError(t, err)
		require.Equal(t, numBeffer-1, bm.Available())

		_, err = bm.Pin(core.NewBlock(fileName, 3))
		require.NoError(t, err)
		require.Equal(t, numBeffer-2, bm.Available())

		_, err = bm.Pin(core.NewBlock(fileName, 4))
		require.NoError(t, err)
		require.Equal(t, numBeffer-3, bm.Available())

		err = bm.Unpin(buf2)
		require.NoError(t, err)
		require.Equal(t, numBeffer-2, bm.Available())

		buf1, err = bm.Pin(core.NewBlock(fileName, 1))
		require.NoError(t, err)
		require.Equal(t, numBeffer-3, bm.Available())
	})
}

func TestBufferManager_FlushAll(t *testing.T) {
	t.Run("test FlushAll", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fm := mock.NewMockFileManager(ctrl)
		lm := mock.NewMockLogManager(ctrl)

		bb := bytes.NewBuffer(100)
		page := core.NewPage(bb)

		fm.EXPECT().PreparePage().Return(page, nil).AnyTimes()
		fm.EXPECT().CopyBlockToPage(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		fm.EXPECT().CopyPageToBlock(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

		lm.EXPECT().FlushByLSN(gomock.Any()).Return(nil).AnyTimes()

		numBuffer := 3
		bm := buffer.NewManager(fm, lm, numBuffer)

		for i := 0; i < numBuffer; i++ {
			buf, err := bm.Pin(fake.Block())
			require.NoError(t, err)
			buf.SetModified(1, 0)
		}

		err := bm.FlushAll(1)
		require.NoError(t, err)
	})
}

func TestBufferManager_pin(t *testing.T) {
	exp := os.NewNormalExplorer()

	blockSize := 400
	dir := "test" + fake.RandString(10)
	logFile := core.FileName("logfile" + fake.RandString(10))
	fileName := core.FileName("testfile" + fake.RandString(10))
	defer exp.RemoveAll(dir)
	fileConfig := infra.NewConfig(dir, blockSize, "logfile")

	fm := file.NewManager(exp, fileConfig)

	logConfig := log.NewConfig(logFile)
	lm := log.NewManager(fm, logConfig)

	numBeffer := 3
	bm := buffer.NewManager(fm, lm, numBeffer)

	t.Run("test buffer", func(t *testing.T) {
		block := core.NewBlock(fileName, 1)
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
		buf2, err := bm.Pin(core.NewBlock(fileName, 2))
		require.NoError(t, err)
		_, err = bm.Pin(core.NewBlock(fileName, 3))
		require.NoError(t, err)
		_, err = bm.Pin(core.NewBlock(fileName, 4))
		require.NoError(t, err)

		err = bm.Unpin(buf2)
		require.NoError(t, err)

		buf1, err = bm.Pin(core.NewBlock(fileName, 1))
		require.NoError(t, err)
		p1 := buf1.GetPage()
		// This modification won't get written to disk.
		err = p1.SetInt32(80, 9999)
		require.NoError(t, err)
		buf1.SetModified(1, 0)
		require.NoError(t, err)
	})
}

func TestBufferManager(t *testing.T) {
	blockSize := 400
	dir := "test" + fake.RandString(10)
	logFile := core.FileName("logfile" + fake.RandString(10))
	fileName := core.FileName("testfile" + fake.RandString(10))
	defer goos.RemoveAll(dir)
	fileConfig := infra.NewConfig(dir, blockSize, "logfile")

	exp := os.NewNormalExplorer()
	fm := file.NewManager(exp, fileConfig)

	logConfig := log.NewConfig(logFile)
	lm := log.NewManager(fm, logConfig)

	numBeffer := 3
	bm := buffer.NewManager(fm, lm, numBeffer)

	t.Run("test buffer manager", func(t *testing.T) {
		var err error

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
