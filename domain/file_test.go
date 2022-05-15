package domain_test

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/lib/bytes"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/goropikari/simpledbgo/testing/mock"
	"github.com/stretchr/testify/require"
)

func TestFileName(t *testing.T) {
	t.Run("valid name", func(t *testing.T) {
		_, err := domain.NewFileName(fake.RandString())
		require.NoError(t, err)
	})

	t.Run("invalid name", func(t *testing.T) {
		_, err := domain.NewFileName("")
		require.Error(t, err)
	})

	t.Run("stringfies", func(t *testing.T) {
		name, err := domain.NewFileName("hello")
		require.NoError(t, err)
		require.Equal(t, "hello", name.String())
	})
}

func TestFile_Read(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		filename := fake.RandString()
		defer os.RemoveAll(filename)

		f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
		require.NoError(t, err)

		file := domain.NewFile(f)

		// Write
		n, err := file.Write([]byte("hello"))
		require.NoError(t, err)
		require.Equal(t, 5, n)

		// Seek
		offset, err := file.Seek(1)
		require.NoError(t, err)
		require.Equal(t, int64(1), offset)

		// Read
		buf1 := make([]byte, 2)
		n, err = file.Read(buf1)
		require.NoError(t, err)
		require.Equal(t, 2, n)
		require.Equal(t, []byte("el"), buf1)

		// Size
		size, err := file.Size()
		require.NoError(t, err)
		require.Equal(t, int64(5), size)

		// FileName
		name := file.Name()
		fileName, err := domain.NewFileName(filename)
		require.NoError(t, err)
		require.Equal(t, fileName, name)

		// Close
		err = file.Close()
		require.NoError(t, err)

		// Read: EOF
		f, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
		require.NoError(t, err)
		file = domain.NewFile(f)
		buf2 := make([]byte, 6)
		_, err = file.Read(buf2)
		require.NoError(t, err)
		require.Equal(t, append([]byte("hello"), 0), buf2)

		n, err = file.Read(buf2)
		require.Error(t, err)
		require.Equal(t, 0, n)

		file.Close()
	})
}

func TestBlock(t *testing.T) {
	t.Run("test block", func(t *testing.T) {
		_, err := domain.NewBlockNumber(fake.RandInt32())
		require.NoError(t, err)
	})

	t.Run("test block", func(t *testing.T) {
		_, err := domain.NewBlockNumber(-1)
		require.Error(t, err)
	})
}

func TestBlock_Equal(t *testing.T) {
	t.Run("test equal", func(t *testing.T) {
		blk1 := fake.Block()
		blk2 := fake.Block()

		require.Equal(t, true, blk1.Equal(blk1))
		require.Equal(t, false, blk1.Equal(blk2))
	})
}

func TestPage_NewPage(t *testing.T) {
	t.Run("test page constructor", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bb := mock.NewMockByteBuffer(ctrl)

		domain.NewPage(bb)
	})
}

func TestPageFactory_Create(t *testing.T) {
	bsf := bytes.NewDirectByteSliceCreater()

	t.Run("test page factory", func(t *testing.T) {
		blockSize, err := domain.NewBlockSize(4096)
		require.NoError(t, err)

		factory := domain.NewPageFactory(bsf, blockSize)
		_, err = factory.Create()
		require.NoError(t, err)
	})

	t.Run("invalid request: test page factory", func(t *testing.T) {
		blockSize, err := domain.NewBlockSize(100)
		require.NoError(t, err)

		factory := domain.NewPageFactory(bsf, blockSize)
		_, err = factory.Create()
		require.Error(t, err)
	})
}
