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
	buf, _ := directio.AlignedBlock(directio.BlockSize)
	bb, err := bytes.NewDirectBufferBytes(buf)
	require.NoError(t, err)
	page := file.NewPage(bb)

	n := int32(10000)
	x1 := rand.Int31n(n)
	page.SetInt32(0, x1)
	x2 := rand.Int31n(n)
	page.SetInt32(4, x2)
	page.SetInt32(8, -x2)
	y1, err := page.GetInt32(0)
	require.NoError(t, err)
	y2, err := page.GetInt32(4)
	require.NoError(t, err)
	y3, err := page.GetInt32(8)
	require.NoError(t, err)
	require.Equal(t, x1, y1)
	require.Equal(t, x2, y2)
	require.Equal(t, -x2, y3)
}

func TestPage_SetGetUInt32(t *testing.T) {
	buf, _ := directio.AlignedBlock(directio.BlockSize)
	bb, err := bytes.NewDirectBufferBytes(buf)
	require.NoError(t, err)
	page := file.NewPage(bb)

	x1 := rand.Uint32()
	page.SetUInt32(0, x1)
	x2 := rand.Uint32()
	page.SetUInt32(4, x2)
	y1, err := page.GetUInt32(0)
	require.NoError(t, err)
	y2, err := page.GetUInt32(4)
	require.NoError(t, err)
	require.Equal(t, x1, y1)
	require.Equal(t, x2, y2)
}

func TestPage_SetGetString(t *testing.T) {
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
}
