package logrecord_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/backend/tx/logrecord"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/goropikari/simpledbgo/testing/mock"
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

	t.Run("start record misc", func(t *testing.T) {
		n := fake.RandInt32()
		rec := &logrecord.StartRecord{
			TxNum: domain.TransactionNumber(n),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		visitor := mock.NewMockTxVisitor(ctrl)

		require.Equal(t, logrecord.Start, rec.Operator())
		require.Equal(t, domain.TransactionNumber(n), rec.TxNumber())
		require.NoError(t, rec.Undo(visitor))
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

	t.Run("commit record misc", func(t *testing.T) {
		n := fake.RandInt32()
		rec := &logrecord.CommitRecord{
			TxNum: domain.TransactionNumber(n),
		}

		require.Equal(t, logrecord.Commit, rec.Operator())
		require.Equal(t, domain.TransactionNumber(n), rec.TxNumber())
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

	t.Run("rollback record misc", func(t *testing.T) {
		n := fake.RandInt32()
		rec := &logrecord.RollbackRecord{
			TxNum: domain.TransactionNumber(n),
		}

		require.Equal(t, logrecord.Rollback, rec.Operator())
		require.Equal(t, domain.TransactionNumber(n), rec.TxNumber())
	})
}

func TestCheckpointRecord(t *testing.T) {
	t.Run("marshal/unmarshal", func(t *testing.T) {
		rec := &logrecord.CheckpointRecord{}

		bytes, err := rec.Marshal()
		require.NoError(t, err)

		rec2 := &logrecord.CheckpointRecord{}
		rec2.Unmarshal(bytes)

		require.Equal(t, *rec, *rec2)
	})

	t.Run("checkpoint record misc", func(t *testing.T) {
		rec := &logrecord.CheckpointRecord{}

		require.Equal(t, logrecord.Checkpoint, rec.Operator())
		require.Equal(t, domain.DummyTransactionNumber, rec.TxNumber())
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

	t.Run("set int32 record misc", func(t *testing.T) {
		rec := &logrecord.SetInt32Record{
			FileName:    "hoge",
			TxNum:       123,
			BlockNumber: 456,
			Offset:      789,
			Val:         111,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		visitor := mock.NewMockTxVisitor(ctrl)
		visitor.EXPECT().UndoSetInt32(gomock.Any()).Return(nil)

		require.Equal(t, logrecord.SetInt32, rec.Operator())
		require.Equal(t, domain.TransactionNumber(123), rec.TxNumber())
		require.NoError(t, rec.Undo(visitor))
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

	t.Run("set string record misc", func(t *testing.T) {
		rec := &logrecord.SetStringRecord{
			FileName:    "hoge",
			TxNum:       123,
			BlockNumber: 456,
			Offset:      789,
			Val:         "foo",
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		visitor := mock.NewMockTxVisitor(ctrl)
		visitor.EXPECT().UndoSetString(gomock.Any()).Return(nil)

		require.Equal(t, logrecord.SetString, rec.Operator())
		require.Equal(t, domain.TransactionNumber(123), rec.TxNumber())
		require.NoError(t, rec.Undo(visitor))
	})
}

func TestSetStringRecord_Error(t *testing.T) {
	t.Run("marshal/unmarshal", func(t *testing.T) {
		rec := &logrecord.SetStringRecord{}
		err := rec.Unmarshal([]byte{1, 2, 3})
		require.Error(t, err)
	})
}
