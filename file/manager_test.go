package file_test

import (
	"io"
	goos "os"
	"testing"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/file"
	"github.com/goropikari/simpledbgo/lib/bytes"
	"github.com/goropikari/simpledbgo/lib/directio"
	"github.com/goropikari/simpledbgo/os"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestManager_NewManager(t *testing.T) {
	t.Run("test file manager", func(t *testing.T) {
		path := "file_" + fake.RandString()
		defer goos.RemoveAll(path)
		config := file.ManagerConfig{
			DBPath:    path,
			DirectIO:  false,
			BlockSize: fake.RandInt32(),
		}

		_, err := file.NewManager(config)
		require.NoError(t, err)
	})
}

func TestManager_NewManager_Error(t *testing.T) {
	t.Run("test file manager: non positive block size", func(t *testing.T) {
		path := "file_" + fake.RandString()
		defer goos.RemoveAll(path)
		config := file.ManagerConfig{
			DBPath:    path,
			DirectIO:  false,
			BlockSize: 0,
		}

		_, err := file.NewManager(config)
		require.Error(t, err)
	})
}

func TestManager_CopyBlockToPage(t *testing.T) {
	// fixture
	blocksize := directio.BlockSize
	bsf := bytes.NewDirectByteSliceCreater()
	buf, err := bsf.Create(blocksize)
	require.NoError(t, err)
	bb := bytes.NewBufferBytes(buf)

	dbpath := "file_" + fake.RandString()
	exp := os.NewDirectIOExplorer(dbpath)

	filename := fake.RandString()
	f, err := exp.OpenFile(domain.FileName(filename))
	require.NoError(t, err)
	defer goos.RemoveAll(dbpath)

	buf[0] = 65
	f.Write(buf)
	f.Seek(0)
	f.Write([]byte("hello"))
	f.Seek(0)

	t.Run("test CopyBlockToPage", func(t *testing.T) {

		blk := domain.NewBlock(domain.FileName(filename), domain.BlockNumber(0))
		page := domain.NewPage(bb)

		config := file.ManagerConfig{
			DBPath:    dbpath,
			BlockSize: int32(blocksize),
			DirectIO:  true,
		}
		mgr, err := file.NewManager(config)
		require.NoError(t, err)

		err = mgr.CopyBlockToPage(blk, page)
		require.NoError(t, err)
		require.Equal(t, buf, page.GetData())
	})
}

// func TestManager_CopyBlockToPage_Error(t *testing.T) {
// 	const size = 10

// 	t.Run("test CopyBlockToPage", func(t *testing.T) {
// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish()

// 		exp := mock.NewMockExplorer(ctrl)
// 		exp.EXPECT().OpenFile(gomock.Any()).Return(nil, errors.New("error"))

// 		bsf := mock.NewMockByteSliceFactory(ctrl)

// 		config := file.ManagerConfig{BlockSize: size}
// 		mgr, err := file.NewManager(exp, bsf, config)
// 		require.NoError(t, err)

// 		blk := fake.Block()
// 		bb := mock.NewMockByteBuffer(ctrl)
// 		bb.EXPECT().Reset()
// 		page := domain.NewPage(bb)

// 		err = mgr.CopyBlockToPage(blk, page)
// 		require.Error(t, err)
// 	})
// }

func TestManager_CopyPageToBlock(t *testing.T) {
	// fixture
	blocksize := directio.BlockSize
	bsf := bytes.NewDirectByteSliceCreater()
	buf, err := bsf.Create(blocksize)
	require.NoError(t, err)
	bb := bytes.NewBufferBytes(buf)

	page := domain.NewPage(bb)
	page.SetString(0, "hoge")
	page.Seek(0, io.SeekStart)

	dbpath := "file_" + fake.RandString()
	exp := os.NewDirectIOExplorer(dbpath)

	filename := fake.RandString()
	f, err := exp.OpenFile(domain.FileName(filename))
	require.NoError(t, err)
	defer goos.RemoveAll(dbpath)

	f.Write(page.GetData())
	f.Seek(0)

	t.Run("test CopyPageToBlock", func(t *testing.T) {
		blk := domain.NewBlock(domain.FileName(filename), domain.BlockNumber(0))
		config := file.ManagerConfig{
			DBPath:    dbpath,
			BlockSize: int32(blocksize),
			DirectIO:  true,
		}
		mgr, err := file.NewManager(config)
		require.NoError(t, err)

		err = mgr.CopyPageToBlock(page, blk)
		require.NoError(t, err)
		f.Close()

		// check written data
		f2, _ := goos.OpenFile(dbpath+"/"+filename, goos.O_RDWR, goos.ModePerm)
		b := make([]byte, 8)
		f2.Read(b)
		require.Equal(t, append([]byte{0, 0, 0, 4}, []byte("hoge")...), b)
	})
}

func TestManager_ExtendFile(t *testing.T) {
	t.Run("test extend file", func(t *testing.T) {
		dbpath := "file_" + fake.RandString()
		blocksize := directio.BlockSize
		config := file.ManagerConfig{
			DBPath:    dbpath,
			BlockSize: int32(blocksize),
			DirectIO:  true,
		}

		mgr, err := file.NewManager(config)
		require.NoError(t, err)
		defer goos.RemoveAll(dbpath)

		filename := fake.FileName()

		// first extend
		blk, err := mgr.ExtendFile(filename)
		require.NoError(t, err)
		require.Equal(t, domain.BlockNumber(0), blk.Number())

		// second extend
		blk, err = mgr.ExtendFile(filename)
		require.NoError(t, err)
		require.Equal(t, domain.BlockNumber(1), blk.Number())

		f, err := mgr.OpenFile(filename)
		require.NoError(t, err)

		n, err := f.Size()
		require.NoError(t, err)
		require.Equal(t, int64(blocksize)*2, n)
	})
}

func TestManager_BlockLength(t *testing.T) {
	t.Run("test extend file", func(t *testing.T) {
		blocksize := directio.BlockSize

		dbpath := "file_" + fake.RandString()
		config := file.ManagerConfig{
			DBPath:    dbpath,
			BlockSize: int32(blocksize),
			DirectIO:  true,
		}

		mgr, err := file.NewManager(config)
		require.NoError(t, err)

		filename, err := domain.NewFileName(fake.RandString())
		require.NoError(t, err)
		defer goos.RemoveAll(dbpath)

		// empty file
		nb, err := mgr.BlockLength(filename)
		require.NoError(t, err)
		require.Equal(t, int32(0), nb)

		// extend file
		_, err = mgr.ExtendFile(filename)
		require.NoError(t, err)
		nb, err = mgr.BlockLength(filename)
		require.NoError(t, err)
		require.Equal(t, int32(1), nb)

		// extend file: second
		_, err = mgr.ExtendFile(filename)
		require.NoError(t, err)
		nb, err = mgr.BlockLength(filename)
		require.NoError(t, err)
		require.Equal(t, int32(2), nb)
	})
}

func TestManager_BlockSize(t *testing.T) {
	t.Run("test extend file", func(t *testing.T) {
		blocksize := directio.BlockSize

		dbpath := "file_" + fake.RandString()
		config := file.ManagerConfig{
			DBPath:    dbpath,
			BlockSize: int32(blocksize),
			DirectIO:  true,
		}

		mgr, err := file.NewManager(config)
		require.NoError(t, err)
		defer goos.RemoveAll(dbpath)

		require.Equal(t, domain.BlockSize(blocksize), mgr.BlockSize())
	})
}
