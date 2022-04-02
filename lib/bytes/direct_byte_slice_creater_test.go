package bytes_test

import (
	"testing"

	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/stretchr/testify/require"
)

func TestDirectByteSliceCreater(t *testing.T) {
	t.Run("test direct slice creater", func(t *testing.T) {
		sc := bytes.NewDirectByteSliceCreater()
		b, err := sc.Create(4096)
		require.NoError(t, err)
		require.Equal(t, 4096, len(b))
	})

	t.Run("invalid request, test direct slice creater", func(t *testing.T) {
		sc := bytes.NewDirectByteSliceCreater()
		b, err := sc.Create(10)
		require.Error(t, err)
		require.Equal(t, 0, len(b))
	})
}
