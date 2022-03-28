package file_test

import (
	"io/ioutil"
	goos "os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goropikari/simpledb_go/backend/core"
	"github.com/goropikari/simpledb_go/backend/file"
	"github.com/goropikari/simpledb_go/infra"
	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/goropikari/simpledb_go/lib/directio"
	"github.com/goropikari/simpledb_go/lib/os"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestFileManager_Config(t *testing.T) {
	t.Run("test file config", func(t *testing.T) {
		config := infra.Config{}
		config.SetDefaults()

		dbDir := "simpledb"
		expected := infra.NewConfig(dbDir, directio.BlockSize, "logfile")

		require.Equal(t, expected, config)
	})
}

func TestFileManager_LastBlock(t *testing.T) {
	t.Run("test LastBlock", func(t *testing.T) {
		db := "/tmp"
		blocksize := 100
		config := infra.NewConfig(db, blocksize, "logfile")

		filename := fake.RandString()
		f, _ := goos.OpenFile(db+"/"+filename, goos.O_RDWR|goos.O_CREATE, goos.ModePerm)
		f.Write(make([]byte, blocksize))
		f.Close()
		defer goos.Remove(db + "/" + filename)

		exp := os.NewNormalExplorer()

		mgr := file.NewManager(exp, config)
		actual, err := mgr.LastBlock(core.FileName(filename))
		blk := core.NewBlock(core.FileName(filename), core.BlockNumber(0))

		require.NoError(t, err)
		require.Equal(t, blk, actual)

		n, err := mgr.FileSize(core.FileName(filename))
		require.NoError(t, err)
		require.Equal(t, int64(blocksize), n)
	})
}

func TestFileManager_PreparePage(t *testing.T) {
	t.Run("test PreparePage", func(t *testing.T) {
		db := "/tmp"
		blocksize := 100
		config := infra.NewConfig(db, blocksize, "logfile")

		filename := fake.RandString()
		f, _ := goos.OpenFile(db+"/"+filename, goos.O_RDWR|goos.O_CREATE, goos.ModePerm)
		f.Write(make([]byte, blocksize))
		f.Close()
		defer goos.Remove(db + "/" + filename)

		exp := os.NewNormalExplorer()

		mgr := file.NewManager(exp, config)

		_, err := mgr.PreparePage()
		require.NoError(t, err)
	})

	t.Run("test PreparePage direct io", func(t *testing.T) {
		db := "."
		blocksize := directio.BlockSize
		config := infra.NewConfig(db, blocksize, "logfile")

		filename := fake.RandString()
		f, _ := goos.OpenFile(db+"/"+filename, goos.O_RDWR|goos.O_CREATE, goos.ModePerm)
		f.Write(make([]byte, blocksize))
		f.Close()
		defer goos.Remove(db + "/" + filename)

		exp := os.NewDirectIOExplorer()

		mgr := file.NewManager(exp, config)

		_, err := mgr.PreparePage()
		require.NoError(t, err)
	})
}

func TestManager(t *testing.T) {
	t.Run("test file manager", func(t *testing.T) {
		// Since tmpfs doesn't support O_DIRECT, dummy data is created at current directory
		// https://github.com/ncw/directio/issues/9
		dir, _ := goos.MkdirTemp(".", "manager-")
		f, _ := goos.CreateTemp(dir, "")
		filename, err := core.NewFileName(filepath.Base(f.Name()))
		require.NoError(t, err)
		err = f.Close()
		require.NoError(t, err)
		err = goos.MkdirAll(dir, goos.ModePerm)
		require.NoError(t, err)
		defer goos.RemoveAll(dir)

		config := infra.NewConfig(dir, directio.BlockSize, "logfile")
		exp := os.NewNormalExplorer()
		fileMgr := file.NewManager(exp, config)

		testFilePath := filepath.Join(dir, string(filename))

		// block size
		require.Equal(t, directio.BlockSize, fileMgr.GetBlockSize())

		// create new file
		_, err = fileMgr.AppendBlock(filename)
		require.NoError(t, err)

		f, err = goos.OpenFile(testFilePath, goos.O_RDONLY, goos.ModePerm)
		require.NoError(t, err)
		info, err := f.Stat()
		require.NoError(t, err)

		require.Equal(t, info.Size(), int64(directio.BlockSize))

		// append block
		_, err = fileMgr.AppendBlock(filename)
		require.NoError(t, err)
		info, err = f.Stat()
		require.NoError(t, err)
		require.Equal(t, info.Size(), int64(directio.BlockSize*2))

		// write page to block
		buf, err := directio.AlignedBlock(directio.BlockSize)
		require.NoError(t, err)
		copy(buf, []byte(strings.Repeat("A", directio.BlockSize)))
		bb, err := bytes.NewDirectBufferBytes(buf)
		require.NoError(t, err)
		page := core.NewPage(bb)
		block := core.NewBlock(filename, 0)
		err = fileMgr.CopyPageToBlock(page, block)
		require.NoError(t, err)

		buf, err = directio.AlignedBlock(directio.BlockSize)
		require.NoError(t, err)
		copy(buf, []byte(strings.Repeat("B", directio.BlockSize)))
		bb, err = bytes.NewDirectBufferBytes(buf)
		require.NoError(t, err)
		page = core.NewPage(bb)
		block = core.NewBlock(filename, 1)
		err = fileMgr.CopyPageToBlock(page, block)
		require.NoError(t, err)
		err = fileMgr.CloseFile(filename)
		require.NoError(t, err)

		content, err := ioutil.ReadFile(testFilePath)
		require.NoError(t, err)
		require.Equal(t, string(content), strings.Repeat("A", directio.BlockSize)+strings.Repeat("B", directio.BlockSize))

		// write block to page
		buf, _ = directio.AlignedBlock(directio.BlockSize)
		bb, err = bytes.NewDirectBufferBytes(buf)
		require.NoError(t, err)
		page = core.NewPage(bb)
		block = core.NewBlock(filename, 0)
		err = fileMgr.CopyBlockToPage(block, page)
		require.NoError(t, err)
		require.Equal(t, strings.Repeat("A", directio.BlockSize), string(page.GetBufferBytes()))

		buf, err = directio.AlignedBlock(directio.BlockSize)
		require.NoError(t, err)
		bb, err = bytes.NewDirectBufferBytes(buf)
		require.NoError(t, err)
		page = core.NewPage(bb)
		block = core.NewBlock(filename, 1)
		err = fileMgr.CopyBlockToPage(block, page)
		require.NoError(t, err)
		require.Equal(t, strings.Repeat("B", directio.BlockSize), string(page.GetBufferBytes()))
	})
}
