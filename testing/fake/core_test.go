package fake_test

import (
	"testing"

	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestBlock(t *testing.T) {
	t.Run("test Block", func(t *testing.T) {
		block := fake.Block()

		require.Regexp(t, "[a-zA-Z]", string(block.GetFileName()))
		require.GreaterOrEqual(t, uint32(block.GetBlockNumber()), uint32(0))
	})
}
