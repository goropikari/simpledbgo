package domain_test

import (
	"os"
	"testing"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/testing/fake"
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

	t.Run("stringfy", func(t *testing.T) {
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
		n, err = file.Read(buf2)
		require.NoError(t, err)
		require.Equal(t, append([]byte("hello"), 0), buf2)

		n, err = file.Read(buf2)
		require.Error(t, err)
		require.Equal(t, 0, n)

		file.Close()
	})
}
