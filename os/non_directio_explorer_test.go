package os_test

import (
	goos "os"
	"testing"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/os"
	"github.com/goropikari/simpledbgo/testing/fake"
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
