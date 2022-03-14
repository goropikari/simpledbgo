package file_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goropikari/simpledb_go/bytes"
	"github.com/goropikari/simpledb_go/core"
	"github.com/goropikari/simpledb_go/directio"
	"github.com/goropikari/simpledb_go/file"
	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {
	// Since tmpfs doesn't support O_DIRECT, dummy data is created at current directory
	// https://github.com/ncw/directio/issues/9
	dir, _ := os.MkdirTemp(".", "manager-")
	f, _ := os.CreateTemp(dir, "")
	filename := core.FileName(filepath.Base(f.Name()))
	err := f.Close()
	require.NoError(t, err)
	err = os.MkdirAll(dir, os.ModePerm)
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	config, err := file.NewConfig(dir, directio.BlockSize, true)
	if err != nil {
		log.Fatal(err)
	}
	fileMgr, err := file.NewManager(config)
	if err != nil {
		log.Fatal(err)
	}

	testFilePath := filepath.Join(dir, string(filename))

	// create new file
	if _, err := fileMgr.AppendBlock(filename); err != nil {
		log.Fatal(err)
	}
	f, _ = os.OpenFile(testFilePath, os.O_RDONLY, os.ModePerm)
	info, _ := f.Stat()
	require.Equal(t, info.Size(), int64(directio.BlockSize))

	// append block
	if _, err := fileMgr.AppendBlock(filename); err != nil {
		log.Fatal(err)
	}
	info, _ = f.Stat()
	require.Equal(t, info.Size(), int64(directio.BlockSize*2))

	// write page to block
	buf, _ := directio.AlignedBlock(directio.BlockSize)
	copy(buf, []byte(strings.Repeat("A", directio.BlockSize)))
	bb, err := bytes.NewDirectBufferBytes(buf)
	require.NoError(t, err)
	page := file.NewPage(bb)
	block := file.NewBlock(filename, 0)
	err = fileMgr.CopyPageToBlock(page, block)
	require.NoError(t, err)

	buf, _ = directio.AlignedBlock(directio.BlockSize)
	copy(buf, []byte(strings.Repeat("B", directio.BlockSize)))
	bb, err = bytes.NewDirectBufferBytes(buf)
	require.NoError(t, err)
	page = file.NewPage(bb)
	block = file.NewBlock(filename, 1)
	err = fileMgr.CopyPageToBlock(page, block)
	require.NoError(t, err)
	err = fileMgr.CloseFile(filename)
	require.NoError(t, err)

	content, _ := ioutil.ReadFile(testFilePath)
	require.Equal(t, string(content), strings.Repeat("A", directio.BlockSize)+strings.Repeat("B", directio.BlockSize))

	// write block to page
	buf, _ = directio.AlignedBlock(directio.BlockSize)
	bb, err = bytes.NewDirectBufferBytes(buf)
	require.NoError(t, err)
	page = file.NewPage(bb)
	block = file.NewBlock(filename, 0)
	err = fileMgr.CopyBlockToPage(block, page)
	require.NoError(t, err)
	require.Equal(t, strings.Repeat("A", directio.BlockSize), string(page.GetFullBytes()))

	buf, _ = directio.AlignedBlock(directio.BlockSize)
	bb, err = bytes.NewDirectBufferBytes(buf)
	require.NoError(t, err)
	page = file.NewPage(bb)
	block = file.NewBlock(filename, 1)
	err = fileMgr.CopyBlockToPage(block, page)
	require.NoError(t, err)
	require.Equal(t, strings.Repeat("B", directio.BlockSize), string(page.GetFullBytes()))
}
