package logrecord_test

import (
	"testing"

	"github.com/goropikari/simpledb_go/backend/tx/logrecord"
	"github.com/stretchr/testify/require"
)

func TestSetInt32Record(t *testing.T) {
	t.Run("marshal/unmarshal", func(t *testing.T) {
		rec := &logrecord.SetInt32Record{
			FileName:    "hoge",
			TxNum:       123,
			BlockNumber: 456,
			Offset:      789,
			Val:         111,
		}

		bytes, err := rec.Marshal()
		require.NoError(t, err)

		rec2 := &logrecord.SetInt32Record{}
		rec2.Unmarshal(bytes)

		require.Equal(t, *rec, *rec2)
	})
}

func TestSetInt32Record_Error(t *testing.T) {
	t.Run("marshal/unmarshal", func(t *testing.T) {
		rec := &logrecord.SetInt32Record{}
		err := rec.Unmarshal([]byte{1, 2, 3})
		require.Error(t, err)
	})
}
