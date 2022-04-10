package list_test

import (
	"testing"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/lib/list"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		list := list.NewList[*domain.Block]()

		blk := fake.Block()
		blk2 := fake.Block()
		list.Add(blk)
		list.Add(blk)

		require.Equal(t, 2, list.Length())
		require.Equal(t, []*domain.Block{blk, blk}, list.Data())

		ok := list.Contains(blk)
		require.Equal(t, true, ok)

		ok = list.Contains(blk2)
		require.Equal(t, false, ok)

		list.Remove(blk)
		ok = list.Contains(blk)
		require.Equal(t, true, ok)

		list.Remove(blk)
		ok = list.Contains(blk)
		require.Equal(t, false, ok)
	})
}
