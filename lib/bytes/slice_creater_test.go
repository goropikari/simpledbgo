package bytes_test

import (
	"testing"

	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/stretchr/testify/require"
)

func TestSliceCreater(t *testing.T) {
	t.Run("test slice creater", func(t *testing.T) {
		sc := bytes.NewSliceCreater()
		b, err := sc.Create(10)
		require.NoError(t, err)
		require.Equal(t, 10, len(b))
	})
}
