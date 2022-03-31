package domain_test

import (
	"errors"
	"io"
	goos "os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/goropikari/simpledb_go/lib/directio"
	"github.com/goropikari/simpledb_go/os"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/goropikari/simpledb_go/testing/mock"
	"github.com/stretchr/testify/require"
)

func TestFileManager_NewFileManager(t *testing.T) {
	t.Run("test file manager", func(t *testing.T) {
		config := domain.FileManagerConfig{BlockSize: fake.RandInt32()}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		exp := mock.NewMockExplorer(ctrl)
		bsf := mock.NewMockByteSliceFactory(ctrl)

		_, err := domain.NewFileManager(exp, bsf, config)
		require.NoError(t, err)
	})
}

func TestFileManager_NewFileManager_Error(t *testing.T) {
	t.Run("test file manager: non positive block size", func(t *testing.T) {
		config := domain.FileManagerConfig{BlockSize: 0}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		exp := mock.NewMockExplorer(ctrl)
		bsf := mock.NewMockByteSliceFactory(ctrl)

		_, err := domain.NewFileManager(exp, bsf, config)
		require.Error(t, err)
	})
}

func TestFileManager_CopyBlockToPage(t *testing.T) {
	t.Run("test CopyBlockToPage", func(t *testing.T) {
		blocksize := directio.BlockSize
		bsf := bytes.NewDirectSliceCreater()
		buf, err := bsf.Create(blocksize)
		require.NoError(t, err)
		bb := bytes.NewBufferBytes(buf)

		dbpath := "."
		exp := os.NewDirectIOExplorer(dbpath)

		filename := fake.RandString()
		file, err := exp.OpenFile(domain.FileName(filename))
		require.NoError(t, err)
		defer file.Remove()

		buf[0] = 65
		file.Write(buf)
		file.Seek(0)
		file.Write([]byte("hello"))
		file.Seek(0)

		blk := domain.NewBlock(file.Name(), domain.BlockSize(blocksize), domain.BlockNumber(0))
		page := domain.NewPage(bb)

		config := domain.FileManagerConfig{BlockSize: int32(blocksize)}
		mgr, err := domain.NewFileManager(exp, bsf, config)
		require.NoError(t, err)

		err = mgr.CopyBlockToPage(blk, page)
		require.NoError(t, err)
		require.Equal(t, buf, page.GetData())
	})
}

func TestFileManager_CopyBlockToPage_Error(t *testing.T) {
	t.Run("test CopyBlockToPage", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		exp := mock.NewMockExplorer(ctrl)
		exp.EXPECT().OpenFile(gomock.Any()).Return(nil, errors.New("error"))

		bsf := mock.NewMockByteSliceFactory(ctrl)

		config := domain.FileManagerConfig{BlockSize: 10}
		mgr, err := domain.NewFileManager(exp, bsf, config)
		require.NoError(t, err)

		blk := fake.Block()
		bb := mock.NewMockByteBuffer(ctrl)
		bb.EXPECT().Reset()
		page := domain.NewPage(bb)

		err = mgr.CopyBlockToPage(blk, page)
		require.Error(t, err)
	})
}

func TestFileManager_CopyPageToBlock(t *testing.T) {
	t.Run("test CopyPageToBlock", func(t *testing.T) {
		blocksize := directio.BlockSize
		bsf := bytes.NewDirectSliceCreater()
		buf, err := bsf.Create(blocksize)
		require.NoError(t, err)
		bb := bytes.NewBufferBytes(buf)

		page := domain.NewPage(bb)
		page.SetString(0, "hoge")
		page.Seek(0, io.SeekStart)

		dbpath := "."
		exp := os.NewDirectIOExplorer(dbpath)

		filename := fake.RandString()
		file, err := exp.OpenFile(domain.FileName(filename))
		require.NoError(t, err)
		defer file.Remove()

		file.Write(buf)
		file.Seek(0)

		blk := domain.NewBlock(file.Name(), domain.BlockSize(blocksize), domain.BlockNumber(0))

		config := domain.FileManagerConfig{BlockSize: int32(blocksize)}
		mgr, err := domain.NewFileManager(exp, bsf, config)
		require.NoError(t, err)

		err = mgr.CopyPageToBlock(page, blk)
		require.NoError(t, err)
		file.Close()

		f, _ := goos.OpenFile(filename, goos.O_RDWR, goos.ModePerm)
		b := make([]byte, 8)
		f.Read(b)
		require.Equal(t, append([]byte{0, 0, 0, 4}, []byte("hoge")...), b)
	})
}

func TestFileManager_ExtendFile(t *testing.T) {
	t.Run("test extend file", func(t *testing.T) {
		blocksize := directio.BlockSize

		dbpath := "."
		exp := os.NewDirectIOExplorer(dbpath)
		bsc := bytes.NewDirectSliceCreater()
		config := domain.FileManagerConfig{
			BlockSize: int32(blocksize),
		}

		mgr, err := domain.NewFileManager(exp, bsc, config)
		require.NoError(t, err)

		filename := fake.FileName()
		defer goos.Remove(string(filename))

		// first extend
		_, err = mgr.ExtendFile(filename)
		require.NoError(t, err)

		// second extend
		blk, err := mgr.ExtendFile(filename)
		require.NoError(t, err)

		expected := domain.NewBlock(filename, domain.BlockSize(blocksize), domain.BlockNumber(1))
		require.Equal(t, expected, blk)
	})
}
