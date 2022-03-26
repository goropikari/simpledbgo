package fake_test

import (
	"testing"

	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestRandInt(t *testing.T) {
	t.Run("test RandInt", func(t *testing.T) {
		for i := 0; i < 10000; i++ {
			x := fake.RandInt()

			require.GreaterOrEqual(t, x, 0)
		}
	})
}
func TestRandInt32(t *testing.T) {
	t.Run("test RandInt32", func(t *testing.T) {
		for i := 0; i < 10000; i++ {
			x := fake.RandInt32()

			require.GreaterOrEqual(t, x, int32(0))
		}
	})
}

func TestRandInt32n(t *testing.T) {
	upper := int32(100)

	t.Run("test RandInt32n", func(t *testing.T) {
		for i := 0; i < 10000; i++ {
			x := fake.RandInt32n(upper)

			require.GreaterOrEqual(t, x, int32(0))
			require.Less(t, x, upper)
		}
	})
}

func TestRandUint32(t *testing.T) {
	t.Run("test RandUint32", func(t *testing.T) {
		for i := 0; i < 10000; i++ {
			x := fake.RandUint32()

			require.GreaterOrEqual(t, x, uint32(0))
		}
	})
}

func TestRandInt64(t *testing.T) {
	t.Run("test RandInt64", func(t *testing.T) {
		for i := 0; i < 10000; i++ {
			x := fake.RandInt64()

			require.GreaterOrEqual(t, x, int64(0))
		}
	})
}

func TestRandString(t *testing.T) {
	length := 10

	t.Run("test RandString", func(t *testing.T) {
		for i := 0; i < 10000; i++ {
			s := fake.RandString(length)

			require.Equal(t, length, len(s))
			require.Regexp(t, "[a-zA-Z]", s)
		}
	})
}
