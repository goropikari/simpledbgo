package file_test

import (
	"math/rand"
	"testing"

	"github.com/goropikari/simpledb_go/bytes"
	"github.com/goropikari/simpledb_go/directio"
	"github.com/goropikari/simpledb_go/file"
	"github.com/stretchr/testify/require"
)

func TestPage_SetGetInt32(t *testing.T) {
	t.Run("test SetGetInt32", func(t *testing.T) {
		buf, _ := directio.AlignedBlock(directio.BlockSize)
		bb, err := bytes.NewDirectBufferBytes(buf)
		require.NoError(t, err)
		page := file.NewPage(bb)

		n := int32(10000)
		x1 := rand.Int31n(n)
		err = page.SetInt32(0, x1)
		require.NoError(t, err)
		x2 := rand.Int31n(n)
		err = page.SetInt32(4, x2)
		require.NoError(t, err)
		err = page.SetInt32(8, -x2)
		require.NoError(t, err)
		y1, err := page.GetInt32(0)
		require.NoError(t, err)
		y2, err := page.GetInt32(4)
		require.NoError(t, err)
		y3, err := page.GetInt32(8)
		require.NoError(t, err)
		require.Equal(t, x1, y1)
		require.Equal(t, x2, y2)
		require.Equal(t, -x2, y3)
	})
}

func TestPage_SetGetUint32(t *testing.T) {
	t.Run("test SetGetUint32", func(t *testing.T) {
		buf, _ := directio.AlignedBlock(directio.BlockSize)
		bb, err := bytes.NewDirectBufferBytes(buf)
		require.NoError(t, err)
		page := file.NewPage(bb)

		x1 := rand.Uint32()
		err = page.SetUint32(0, x1)
		require.NoError(t, err)
		x2 := rand.Uint32()
		err = page.SetUint32(4, x2)
		require.NoError(t, err)
		y1, err := page.GetUint32(0)
		require.NoError(t, err)
		y2, err := page.GetUint32(4)
		require.NoError(t, err)
		require.Equal(t, x1, y1)
		require.Equal(t, x2, y2)
	})
}

func TestPage_SetGetBytes(t *testing.T) {
	t.Run("test SetGetUint32", func(t *testing.T) {
		buf, _ := directio.AlignedBlock(directio.BlockSize)
		bb, err := bytes.NewDirectBufferBytes(buf)
		require.NoError(t, err)
		page := file.NewPage(bb)

		x1 := []byte("hello")
		err = page.SetBytes(0, x1)
		require.NoError(t, err)
		x2 := []byte("world")
		err = page.SetBytes(9, x2)
		require.NoError(t, err)
		y1, err := page.GetBytes(0)
		require.NoError(t, err)
		y2, err := page.GetBytes(9)
		require.NoError(t, err)
		require.Equal(t, x1, y1)
		require.Equal(t, x2, y2)
	})
}

func TestPage_SetGetString(t *testing.T) {
	t.Run("test SetGetUint32", func(t *testing.T) {
		buf, _ := directio.AlignedBlock(directio.BlockSize)
		bb, err := bytes.NewDirectBufferBytes(buf)
		require.NoError(t, err)
		page := file.NewPage(bb)

		x1 := "foo"
		x2 := "bar"
		err = page.SetString(0, x1)
		require.NoError(t, err)
		err = page.SetString(7, x2)
		require.NoError(t, err)
		y1, err := page.GetString(0)
		require.NoError(t, err)
		y2, err := page.GetString(7)
		require.NoError(t, err)
		require.Equal(t, x1, y1)
		require.Equal(t, x2, y2)
	})
}
