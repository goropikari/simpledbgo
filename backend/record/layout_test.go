package record_test

import (
	"testing"

	"github.com/goropikari/simpledbgo/backend/record"
	"github.com/stretchr/testify/require"
)

func TestLayout(t *testing.T) {
	t.Run("constructs Layout", func(t *testing.T) {
		schema := record.NewSchema()
		schema.AddField("hoge", record.Integer, 0)
		schema.AddField("piyo", record.String, 8)

		layout := record.NewLayout(schema)

		mp := map[record.FieldName]int64{
			"hoge": 4,
			"piyo": 8,
		}

		expected := record.NewLayoutByElement(schema, mp, 20)

		require.Equal(t, expected, layout)
	})
}
