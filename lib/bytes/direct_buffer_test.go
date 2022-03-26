package bytes_test

import (
	"testing"

	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/ncw/directio"
	"github.com/stretchr/testify/require"
)

func TestDirectBuffer(t *testing.T) {
	t.Run("test DirectBuffer", func(t *testing.T) {
		_, err := bytes.NewDirectBuffer(int(directio.BlockSize))
		require.NoError(t, err)

		_, err = bytes.NewDirectBuffer(10)
		require.Error(t, err)
	})
}

func TestDirectBufferBytes(t *testing.T) {
	t.Run("test DirectBufferBytes", func(t *testing.T) {
		_, err := bytes.NewDirectBufferBytes(make([]byte, 10))
		require.Error(t, err)
	})
}
