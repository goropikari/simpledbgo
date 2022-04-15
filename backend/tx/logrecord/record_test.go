package logrecord_test

import (
	"testing"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/backend/tx/logrecord"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestStartRecord(t *testing.T) {
	t.Run("marshal/unmarshal", func(t *testing.T) {
		rec := &logrecord.StartRecord{
			TxNum: domain.TransactionNumber(fake.RandInt32()),
		}

		bytes, err := rec.Marshal()
		require.NoError(t, err)

		rec2 := &logrecord.StartRecord{}
		rec2.Unmarshal(bytes)

		require.Equal(t, *rec, *rec2)
	})
}

func TestCommitRecord(t *testing.T) {
	t.Run("marshal/unmarshal", func(t *testing.T) {
		rec := &logrecord.CommitRecord{
			TxNum: domain.TransactionNumber(fake.RandInt32()),
		}

		bytes, err := rec.Marshal()
		require.NoError(t, err)

		rec2 := &logrecord.CommitRecord{}
		rec2.Unmarshal(bytes)

		require.Equal(t, *rec, *rec2)
	})
}

func TestRollbackRecord(t *testing.T) {
	t.Run("marshal/unmarshal", func(t *testing.T) {
		rec := &logrecord.RollbackRecord{
			TxNum: domain.TransactionNumber(fake.RandInt32()),
		}

		bytes, err := rec.Marshal()
		require.NoError(t, err)

		rec2 := &logrecord.RollbackRecord{}
		rec2.Unmarshal(bytes)

		require.Equal(t, *rec, *rec2)
	})
}

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

func TestSetStringRecord(t *testing.T) {
	t.Run("marshal/unmarshal", func(t *testing.T) {
		rec := &logrecord.SetStringRecord{
			FileName:    domain.FileName(fake.RandString()),
			TxNum:       123,
			BlockNumber: 456,
			Offset:      789,
			Val:         fake.RandString(),
		}

		bytes, err := rec.Marshal()
		require.NoError(t, err)

		rec2 := &logrecord.SetStringRecord{}
		rec2.Unmarshal(bytes)

		require.Equal(t, *rec, *rec2)
	})
}

func TestSetStringRecord_Error(t *testing.T) {
	t.Run("marshal/unmarshal", func(t *testing.T) {
		rec := &logrecord.SetStringRecord{}
		err := rec.Unmarshal([]byte{1, 2, 3})
		require.Error(t, err)
	})
}
