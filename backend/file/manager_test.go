package file_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goropikari/simpledb_go/backend/core"
	"github.com/goropikari/simpledb_go/backend/file"
	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/goropikari/simpledb_go/lib/directio"
	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {
	t.Run("test file manager", func(t *testing.T) {
		// Since tmpfs doesn't support O_DIRECT, dummy data is created at current directory
		// https://github.com/ncw/directio/issues/9
		dir, _ := os.MkdirTemp(".", "manager-")
		f, _ := os.CreateTemp(dir, "")
		filename, err := core.NewFileName(filepath.Base(f.Name()))
		require.NoError(t, err)
		err = f.Close()
		require.NoError(t, err)
		err = os.MkdirAll(dir, os.ModePerm)
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		isDirectIO := true
		config, err := file.NewConfig(dir, directio.BlockSize, isDirectIO)
		require.NoError(t, err)
		fileMgr, err := file.NewManager(config)
		require.NoError(t, err)

		testFilePath := filepath.Join(dir, string(filename))

		// create new file
		_, err = fileMgr.AppendBlock(filename)
		require.NoError(t, err)

		f, err = os.OpenFile(testFilePath, os.O_RDONLY, os.ModePerm)
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
