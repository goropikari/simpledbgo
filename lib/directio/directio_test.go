package directio_test

import (
	"os"
	"testing"

	"github.com/goropikari/simpledb_go/lib/directio"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestDirectIO_AlignedBlock(t *testing.T) {
	t.Run("test direct io: AlignedBlock", func(t *testing.T) {
		_, err := directio.AlignedBlock(directio.BlockSize)
		require.NoError(t, err)

		_, err = directio.AlignedBlock(10)
		require.Error(t, err)
	})
}

func TestDirectIO_IsAligned(t *testing.T) {
	t.Run("test direct io: IsAligned", func(t *testing.T) {
		b, err := directio.AlignedBlock(directio.BlockSize)
		require.NoError(t, err)
		require.Equal(t, true, directio.IsAligned(b))

		require.Equal(t, false, directio.IsAligned(make([]byte, 10)))
	})
}

func TestDirectIO_OpenFile(t *testing.T) {
	t.Run("test direct io: IsAligned", func(t *testing.T) {
		// Since tmpfs doesn't support O_DIRECT, dummy data is created at current directory
		// https://github.com/ncw/directio/issues/9
		filename := fake.RandString()
		_, err := directio.OpenFile(filename, os.O_CREATE, os.ModePerm)
		defer os.Remove(filename)
		require.NoError(t, err)
	})
}
