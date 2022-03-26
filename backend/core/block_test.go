package core_test

import (
	"testing"

	"github.com/goropikari/simpledb_go/backend/core"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	t.Run("test block", func(t *testing.T) {
		filename, _ := core.NewFileName("hoge")
		blknum, _ := core.NewBlockNumber(0)
		block := core.NewBlock(filename, blknum)

		require.Equal(t, filename, block.GetFileName())
		require.Equal(t, blknum, block.GetBlockNumber())
		require.Equal(t, true, block.Equal(block))
		require.Equal(t, "[file hoge, block 0]", block.String())
	})
}
