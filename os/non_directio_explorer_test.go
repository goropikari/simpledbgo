package os_test

import (
	goos "os"
	"testing"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/os"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestNonDirectIOExplorer(t *testing.T) {
	t.Run("test normal explorer", func(t *testing.T) {
		exp := os.NewNonDirectIOExplorer(".")

		filename := fake.RandString()
		defer goos.Remove(filename)

		_, err := exp.OpenFile(domain.FileName(filename))
		require.NoError(t, err)
	})
}
