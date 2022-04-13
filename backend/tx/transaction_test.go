package tx_test

import (
	"testing"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/backend/tx"
	"github.com/goropikari/simpledb_go/backend/tx/logrecord"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestTransaction_Start(t *testing.T) {
	t.Run("test commit", func(t *testing.T) {
		const (
			blockSize = 20
			numBuf    = 10
		)

		dbPath := fake.RandString()
		factory := fake.NewNonDirectBufferManagerFactory(dbPath, blockSize, numBuf)
		fileMgr, logMgr, bufMgr := factory.Create()
		defer factory.Finish()

		ltConfig := tx.NewConfig(1000)
		lt := tx.NewLockTable(ltConfig)
		concurMgr := tx.NewConcurrencyManager(lt)

		gen := tx.NewNumberGenerator()
		_, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)

		err = logMgr.Flush()
		require.NoError(t, err)

		it, err := logMgr.Iterator()
		require.NoError(t, err)
		records := make([]logrecord.LogRecorder, 0)
		for it.HasNext() {
			data, err := it.Next()
			require.NoError(t, err)
			rec, err := tx.RecordParse(data)
			require.NoError(t, err)
			records = append(records, rec)
		}

		require.Equal(t, 1, len(records))
		expected := &logrecord.StartRecord{TxNum: domain.TransactionNumber(1)}
		require.Equal(t, expected, records[0])
	})
}

func TestTransaction_Commit(t *testing.T) {
	t.Run("test commit", func(t *testing.T) {
		const (
			blockSize = 20
			numBuf    = 10
		)

		dbPath := fake.RandString()
		factory := fake.NewNonDirectBufferManagerFactory(dbPath, blockSize, numBuf)
		fileMgr, logMgr, bufMgr := factory.Create()
		defer factory.Finish()

		ltConfig := tx.NewConfig(1000)
		lt := tx.NewLockTable(ltConfig)
		concurMgr := tx.NewConcurrencyManager(lt)

		gen := tx.NewNumberGenerator()
		gen.Generate()
		transaction, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)

		err = transaction.Commit()
		require.NoError(t, err)

		err = logMgr.Flush()
		require.NoError(t, err)

		it, err := logMgr.Iterator()
		require.NoError(t, err)
		records := make([]logrecord.LogRecorder, 0)
		for it.HasNext() {
			data, err := it.Next()
			require.NoError(t, err)
			rec, err := tx.RecordParse(data)
			require.NoError(t, err)
			records = append(records, rec)
		}

		require.Equal(t, 2, len(records))
		expected := []logrecord.LogRecorder{
			&logrecord.CommitRecord{TxNum: domain.TransactionNumber(2)},
			&logrecord.StartRecord{TxNum: domain.TransactionNumber(2)},
		}
		require.Equal(t, expected, records)
	})
}

func TestTransaction_GetSetInt32(t *testing.T) {
	t.Run("test commit", func(t *testing.T) {
		const (
			blockSize = 100
			numBuf    = 2
		)

		dbPath := fake.RandString()
		factory := fake.NewNonDirectBufferManagerFactory(dbPath, blockSize, numBuf)
		fileMgr, logMgr, bufMgr := factory.Create()
		defer factory.Finish()

		ltConfig := tx.NewConfig(1000)
		lt := tx.NewLockTable(ltConfig)
		concurMgr := tx.NewConcurrencyManager(lt)

		gen := tx.NewNumberGenerator()
		transaction, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)

		blk := *domain.NewBlock(
			domain.FileName("table_"+fake.RandString()),
			domain.BlockSize(blockSize),
			domain.BlockNumber(0),
		)

		offset := int64(10)
		val := int32(100)
		writeLog := true
		err = transaction.Pin(blk)
		require.NoError(t, err)
		err = transaction.SetInt32(blk, offset, val, writeLog)
		require.NoError(t, err)
		err = transaction.SetInt32(blk, offset, val+1, writeLog)
		require.NoError(t, err)
		v, err := transaction.GetInt32(blk, offset)
		require.NoError(t, err)
		require.Equal(t, val+1, v)
		transaction.Commit()

		err = logMgr.Flush()
		require.NoError(t, err)

		it, err := logMgr.Iterator()
		require.NoError(t, err)
		records := make([]logrecord.LogRecorder, 0)
		for it.HasNext() {
			data, err := it.Next()
			require.NoError(t, err)
			rec, err := tx.RecordParse(data)
			require.NoError(t, err)
			records = append(records, rec)
		}

		require.Equal(t, 4, len(records))
		expected := []logrecord.LogRecorder{
			&logrecord.CommitRecord{TxNum: domain.TransactionNumber(1)},
			&logrecord.SetInt32Record{
				FileName:    blk.FileName(),
				TxNum:       transaction.Number(),
				BlockNumber: blk.Number(),
				Offset:      offset,
				Val:         val,
			},
			&logrecord.SetInt32Record{
				FileName:    blk.FileName(),
				TxNum:       transaction.Number(),
				BlockNumber: blk.Number(),
				Offset:      offset,
				Val:         0,
			},
			&logrecord.StartRecord{TxNum: domain.TransactionNumber(1)},
		}
		require.Equal(t, expected, records)
	})
}

func TestTransaction_GetSetString(t *testing.T) {
	t.Run("test commit", func(t *testing.T) {
		const (
			blockSize = 100
			numBuf    = 2
		)

		dbPath := fake.RandString()
		factory := fake.NewNonDirectBufferManagerFactory(dbPath, blockSize, numBuf)
		fileMgr, logMgr, bufMgr := factory.Create()
		defer factory.Finish()

		ltConfig := tx.NewConfig(1000)
		lt := tx.NewLockTable(ltConfig)
		concurMgr := tx.NewConcurrencyManager(lt)

		gen := tx.NewNumberGenerator()
		transaction, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)

		blk := *domain.NewBlock(
			domain.FileName("table_"+fake.RandString()),
			domain.BlockSize(blockSize),
			domain.BlockNumber(0),
		)

		offset := int64(10)
		val1 := fake.RandString()
		val2 := fake.RandString()
		writeLog := true
		err = transaction.Pin(blk)
		require.NoError(t, err)
		err = transaction.SetString(blk, offset, val1, writeLog)
		require.NoError(t, err)
		err = transaction.SetString(blk, offset, val2, writeLog)
		require.NoError(t, err)
		v, err := transaction.GetString(blk, offset)
		require.NoError(t, err)
		require.Equal(t, val2, v)
		transaction.Commit()

		err = logMgr.Flush()
		require.NoError(t, err)

		it, err := logMgr.Iterator()
		require.NoError(t, err)
		records := make([]logrecord.LogRecorder, 0)
		for it.HasNext() {
			data, err := it.Next()
			require.NoError(t, err)
			rec, err := tx.RecordParse(data)
			require.NoError(t, err)
			records = append(records, rec)
		}

		require.Equal(t, 4, len(records))
		expected := []logrecord.LogRecorder{
			&logrecord.CommitRecord{TxNum: domain.TransactionNumber(1)},
			&logrecord.SetStringRecord{
				FileName:    blk.FileName(),
				TxNum:       transaction.Number(),
				BlockNumber: blk.Number(),
				Offset:      offset,
				Val:         val1,
			},
			&logrecord.SetStringRecord{
				FileName:    blk.FileName(),
				TxNum:       transaction.Number(),
				BlockNumber: blk.Number(),
				Offset:      offset,
				Val:         "",
			},
			&logrecord.StartRecord{TxNum: domain.TransactionNumber(1)},
		}
		require.Equal(t, expected, records)
	})
}
