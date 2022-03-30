package os_test

import (
	goos "os"
	"testing"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/os"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestDirectIOExplorer(t *testing.T) {
	t.Run("test direct io explorer", func(t *testing.T) {
		filename := domain.FileName(fake.RandString())
		defer goos.Remove(string(filename))

		exp := os.NewDirectIOExplorer(".")
		_, err := exp.OpenFile(filename)
		require.NoError(t, err)

	})
}
