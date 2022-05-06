package domain_test

import (
	"testing"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/stretchr/testify/require"
)

func TestLayout(t *testing.T) {
	t.Run("constructs Layout", func(t *testing.T) {
		schema := domain.NewSchema()
		schema.AddField("hoge", domain.FInt32, 0)
		schema.AddField("piyo", domain.FString, 8)

		layout := domain.NewLayout(schema)

		mp := map[domain.FieldName]int64{
			"hoge": 4,
			"piyo": 8,
		}

		expected := domain.NewLayoutByElement(schema, mp, 20)

		require.Equal(t, expected, layout)
	})
}
