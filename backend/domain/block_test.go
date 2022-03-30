package domain_test

import (
	"testing"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestBlock(t *testing.T) {
	t.Run("test block", func(t *testing.T) {
		_, err := domain.NewBlockNumber(fake.RandInt32())
		require.NoError(t, err)
	})

	t.Run("test block", func(t *testing.T) {
		_, err := domain.NewBlockNumber(-1)
		require.Error(t, err)
	})
}

func TestBlock_Equal(t *testing.T) {
	t.Run("test equal", func(t *testing.T) {
		blk1 := fake.Block()
		blk2 := fake.Block()

		require.Equal(t, true, blk1.Equal(blk1))
		require.Equal(t, false, blk1.Equal(blk2))
	})
}
