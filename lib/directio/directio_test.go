package directio_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/goropikari/simpledb_go/lib/directio"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestDirectBuffer(t *testing.T) {
	t.Run("test direct io", func(t *testing.T) {
		// Since tmpfs doesn't support O_DIRECT, dummy data is created at current directory
		// https://github.com/ncw/directio/issues/9
		dir, err := os.MkdirTemp(".", "directio")
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		filename := fake.RandString(10)
		flag := os.O_RDWR | os.O_CREATE
		f, err := directio.OpenFile(filepath.Join(dir, filename), flag, os.ModePerm)
		require.NoError(t, err)

		buf, err := directio.AlignedBlock(directio.BlockSize)
		require.NoError(t, err)

		s := fake.RandString(directio.BlockSize)
		copy(buf, []byte(s))

		_, err = f.Write(buf)
		require.NoError(t, err)

		_, err = f.Seek(0, io.SeekStart)
		require.NoError(t, err)

		readbytes, err := directio.AlignedBlock(directio.BlockSize)
		require.NoError(t, err)

		_, err = f.Read(readbytes)
		require.NoError(t, err)
		require.Equal(t, string(readbytes), s)
	})
}
