package os_test

import (
	"io"
	goos "os"
	"testing"

	"github.com/goropikari/simpledb_go/lib/os"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestFile(t *testing.T) {
	t.Run("test file", func(t *testing.T) {
		filename := fake.RandString(10)
		f, err := os.OpenFile(filename, false)
		defer goos.RemoveAll(filename)

		require.NoError(t, err)

		n, err := f.Write(make([]byte, 100))
		require.NoError(t, err)
		require.Equal(t, 100, n)

		size, err := f.Size()
		require.NoError(t, err)
		require.Equal(t, int64(100), size)

		_, err = f.Seek(10, io.SeekStart)
		require.NoError(t, err)

		_, err = f.Write([]byte{1, 2, 3, 4})
		require.NoError(t, err)

		_, err = f.Seek(10, io.SeekStart)
		require.NoError(t, err)

		buf := make([]byte, 4)
		n, err = f.Read(buf)
		require.NoError(t, err)
		require.Equal(t, []byte{1, 2, 3, 4}, buf)

		err = f.Close()
		require.NoError(t, err)
	})
}
