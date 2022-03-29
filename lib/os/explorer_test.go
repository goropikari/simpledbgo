package os_test

import (
	goos "os"
	"testing"

	"github.com/goropikari/simpledb_go/lib/os"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestDirectIOExplorer(t *testing.T) {
	t.Run("test DirectIOExplorer", func(t *testing.T) {
		filename := fake.RandString()

		exp := os.NewDirectIOExplorer()
		defer goos.RemoveAll(filename)

		f, err := exp.OpenFile(filename)
		require.NoError(t, err)

		err = f.Close()
		require.NoError(t, err)
	})
}
