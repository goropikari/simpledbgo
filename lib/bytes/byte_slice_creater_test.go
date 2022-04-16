package bytes_test

import (
	"testing"

	"github.com/goropikari/simpledbgo/lib/bytes"
	"github.com/stretchr/testify/require"
)

func TestByteSliceCreater(t *testing.T) {
	t.Run("test slice creater", func(t *testing.T) {
		sc := bytes.NewByteSliceCreater()
		b, err := sc.Create(10)
		require.NoError(t, err)
		require.Equal(t, 10, len(b))
	})
}
