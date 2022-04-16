package domain_test

import (
	"testing"

	"github.com/goropikari/simpledbgo/backend/domain"
	"github.com/goropikari/simpledbgo/lib/bytes"
	"github.com/stretchr/testify/require"
)

func TestPageFactory_Create(t *testing.T) {
	bsf := bytes.NewDirectByteSliceCreater()

	t.Run("test page factory", func(t *testing.T) {
		blockSize, err := domain.NewBlockSize(4096)
		require.NoError(t, err)

		factory := domain.NewPageFactory(bsf, blockSize)
		_, err = factory.Create()
		require.NoError(t, err)
	})

	t.Run("invalid request: test page factory", func(t *testing.T) {
		blockSize, err := domain.NewBlockSize(100)
		require.NoError(t, err)

		factory := domain.NewPageFactory(bsf, blockSize)
		_, err = factory.Create()
		require.Error(t, err)
	})
}
