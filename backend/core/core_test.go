package core_test

import (
	"testing"

	"github.com/goropikari/simpledb_go/backend/core"
	"github.com/stretchr/testify/require"
)

func TestFileName(t *testing.T) {
	t.Run("test file name", func(t *testing.T) {
		name, err := core.NewFileName("hoge")
		require.NoError(t, err)
		require.Equal(t, core.FileName("hoge"), name)
	})
}

func TestFileName_Error(t *testing.T) {
	t.Run("test file name", func(t *testing.T) {
		_, err := core.NewFileName("")
		require.Error(t, err)
	})
}

func TestBlockNumber(t *testing.T) {
	t.Run("test block number", func(t *testing.T) {
		num, err := core.NewBlockNumber(10)
		require.NoError(t, err)
		require.Equal(t, core.BlockNumber(10), num)
	})
}

func TestBlockNumber_Error(t *testing.T) {
	t.Run("test block number", func(t *testing.T) {
		_, err := core.NewBlockNumber(-1)
		require.Error(t, err)
	})
}

func TestBlockSize(t *testing.T) {
	t.Run("test block size", func(t *testing.T) {
		size, err := core.NewBlockSize(10)
		require.NoError(t, err)
		require.Equal(t, core.BlockSize(10), size)
	})
}

func TestBlockSize_Error(t *testing.T) {
	t.Run("test block size", func(t *testing.T) {
		_, err := core.NewBlockSize(-1)
		require.Error(t, err)
	})
}
