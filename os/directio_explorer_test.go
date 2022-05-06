package os_test

import (
	goos "os"
	"testing"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/os"
	"github.com/goropikari/simpledbgo/testing/fake"
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
