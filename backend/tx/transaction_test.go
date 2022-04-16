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
			rec, err := tx.ParseRecord(data)
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
			rec, err := tx.ParseRecord(data)
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
			rec, err := tx.ParseRecord(data)
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
			rec, err := tx.ParseRecord(data)
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

func TestTransaction_Rollback(t *testing.T) {
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
		txn1, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)

		blk := *domain.NewBlock(
			domain.FileName("table_"+fake.RandString()),
			domain.BlockSize(blockSize),
			domain.BlockNumber(0),
		)

		offset := int64(10)
		val := int32(100)
		writeLog := true
		err = txn1.Pin(blk)
		require.NoError(t, err)
		err = txn1.SetInt32(blk, offset, val, writeLog)
		require.NoError(t, err)
		err = txn1.SetString(blk, offset+4, "foo", writeLog)
		require.NoError(t, err)
		txn1.Commit()

		txn2, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)
		err = txn2.Pin(blk)
		require.NoError(t, err)
		err = txn2.SetInt32(blk, offset, val+1, writeLog)
		require.NoError(t, err)
		err = txn2.SetString(blk, offset+4, "bar", writeLog)
		require.NoError(t, err)
		err = txn2.SetInt32(blk, offset, val+2, writeLog)
		require.NoError(t, err)
		v, err := txn2.GetInt32(blk, offset)
		require.NoError(t, err)
		require.Equal(t, val+2, v)
		vs, err := txn2.GetString(blk, offset+4)
		require.NoError(t, err)
		require.Equal(t, "bar", vs)
		txn2.Rollback()

		txn3, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)
		err = txn3.Pin(blk)
		require.NoError(t, err)
		v2, err := txn3.GetInt32(blk, offset)
		require.NoError(t, err)
		require.Equal(t, val, v2)
		vs2, err := txn3.GetString(blk, offset+4)
		require.NoError(t, err)
		require.Equal(t, "foo", vs2)
		txn3.Commit()

		it, err := logMgr.Iterator()
		require.NoError(t, err)
		records := make([]logrecord.LogRecorder, 0)
		for it.HasNext() {
			data, err := it.Next()
			require.NoError(t, err)
			rec, err := tx.ParseRecord(data)
			require.NoError(t, err)
			records = append(records, rec)
		}

		expected := []logrecord.LogRecorder{
			// third transaction
			&logrecord.CommitRecord{TxNum: domain.TransactionNumber(3)},
			&logrecord.StartRecord{TxNum: domain.TransactionNumber(3)},
			// second transaction
			&logrecord.RollbackRecord{TxNum: domain.TransactionNumber(2)},
			&logrecord.SetInt32Record{
				FileName:    blk.FileName(),
				TxNum:       txn2.Number(),
				BlockNumber: blk.Number(),
				Offset:      offset,
				Val:         val + 1,
			},
			&logrecord.SetStringRecord{
				FileName:    blk.FileName(),
				TxNum:       txn2.Number(),
				BlockNumber: blk.Number(),
				Offset:      offset + 4,
				Val:         "foo",
			},
			&logrecord.SetInt32Record{
				FileName:    blk.FileName(),
				TxNum:       txn2.Number(),
				BlockNumber: blk.Number(),
				Offset:      offset,
				Val:         val,
			},
			&logrecord.StartRecord{TxNum: domain.TransactionNumber(2)},
			// first transaction
			&logrecord.CommitRecord{TxNum: domain.TransactionNumber(1)},
			&logrecord.SetStringRecord{
				FileName:    blk.FileName(),
				TxNum:       txn1.Number(),
				BlockNumber: blk.Number(),
				Offset:      offset + 4,
				Val:         "",
			},
			&logrecord.SetInt32Record{
				FileName:    blk.FileName(),
				TxNum:       txn1.Number(),
				BlockNumber: blk.Number(),
				Offset:      offset,
				Val:         0,
			},
			&logrecord.StartRecord{TxNum: domain.TransactionNumber(1)},
		}
		require.Equal(t, expected, records)
	})
}

func TestTransaction_Recover(t *testing.T) {
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

		filename := "table_" + fake.RandString()
		blk := *domain.NewBlock(
			domain.FileName(filename),
			domain.BlockSize(blockSize),
			domain.BlockNumber(0),
		)
		blk2 := *domain.NewBlock(
			domain.FileName(filename),
			domain.BlockSize(blockSize),
			domain.BlockNumber(1),
		)

		offset := int64(10)
		val := int32(100)
		writeLog := true

		// commit
		txn1, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)
		err = txn1.Pin(blk)
		require.NoError(t, err)
		err = txn1.SetInt32(blk, offset, val, writeLog)
		require.NoError(t, err)
		err = txn1.SetString(blk, offset+4, "foo", writeLog)
		require.NoError(t, err)
		txn1.Commit()

		// uncommit
		txn2, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)
		err = txn2.Pin(blk2)
		require.NoError(t, err)
		err = txn2.SetInt32(blk2, offset, val+1, writeLog)
		require.NoError(t, err)
		err = txn2.SetString(blk2, offset+4, "bar", writeLog)
		require.NoError(t, err)
		err = txn2.SetInt32(blk2, offset, val+2, writeLog)
		require.NoError(t, err)
		v, err := txn2.GetInt32(blk2, offset)
		require.NoError(t, err)
		require.Equal(t, val+2, v)
		vs, err := txn2.GetString(blk2, offset+4)
		require.NoError(t, err)
		require.Equal(t, "bar", vs)
		txn2.Unpin(blk2)

		// recover
		txn3, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)
		err = txn3.Recover()
		require.NoError(t, err)

		// uncommit
		txn4, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)
		err = txn4.Pin(blk)
		require.NoError(t, err)
		vs4, err := txn4.GetString(blk, offset+4)
		require.NoError(t, err)
		require.Equal(t, "foo", vs4)
		v4, err := txn4.GetInt32(blk, offset)
		require.NoError(t, err)
		require.Equal(t, val, v4)

		// uncommit
		txn5, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)
		err = txn5.Pin(blk2)
		require.NoError(t, err)
		err = txn5.SetInt32(blk2, offset, val+1, writeLog)
		require.NoError(t, err)
		err = txn5.SetString(blk2, offset+4, "baz", writeLog)
		require.NoError(t, err)

		// recover
		txn6, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)
		err = txn6.Recover()
		require.NoError(t, err)

		it, err := logMgr.Iterator()
		require.NoError(t, err)
		records := make([]logrecord.LogRecorder, 0)
		for it.HasNext() {
			data, err := it.Next()
			require.NoError(t, err)
			rec, err := tx.ParseRecord(data)
			require.NoError(t, err)
			records = append(records, rec)
		}

		expected := []logrecord.LogRecorder{
			&logrecord.CheckpointRecord{},
			&logrecord.StartRecord{TxNum: domain.TransactionNumber(6)},
			// fifth transaction
			&logrecord.SetStringRecord{
				FileName:    blk2.FileName(),
				TxNum:       txn5.Number(),
				BlockNumber: blk2.Number(),
				Offset:      offset + 4,
				Val:         "",
			},
			&logrecord.SetInt32Record{
				FileName:    blk2.FileName(),
				TxNum:       txn5.Number(),
				BlockNumber: blk2.Number(),
				Offset:      offset,
				Val:         0,
			},
			&logrecord.StartRecord{TxNum: domain.TransactionNumber(5)},
			// forth transaction
			&logrecord.StartRecord{TxNum: domain.TransactionNumber(4)},
			// third transaction
			&logrecord.CheckpointRecord{},
			&logrecord.StartRecord{TxNum: domain.TransactionNumber(3)},
			&logrecord.SetInt32Record{
				FileName:    blk2.FileName(),
				TxNum:       txn2.Number(),
				BlockNumber: blk2.Number(),
				Offset:      offset,
				Val:         val + 1,
			},
			&logrecord.SetStringRecord{
				FileName:    blk2.FileName(),
				TxNum:       txn2.Number(),
				BlockNumber: blk2.Number(),
				Offset:      offset + 4,
				Val:         "",
			},
			&logrecord.SetInt32Record{
				FileName:    blk2.FileName(),
				TxNum:       txn2.Number(),
				BlockNumber: blk2.Number(),
				Offset:      offset,
				Val:         0,
			},
			&logrecord.StartRecord{TxNum: domain.TransactionNumber(2)},
			// first transaction
			&logrecord.CommitRecord{TxNum: domain.TransactionNumber(1)},
			&logrecord.SetStringRecord{
				FileName:    blk.FileName(),
				TxNum:       txn1.Number(),
				BlockNumber: blk.Number(),
				Offset:      offset + 4,
				Val:         "",
			},
			&logrecord.SetInt32Record{
				FileName:    blk.FileName(),
				TxNum:       txn1.Number(),
				BlockNumber: blk.Number(),
				Offset:      offset,
				Val:         0,
			},
			&logrecord.StartRecord{TxNum: domain.TransactionNumber(1)},
		}
		require.Equal(t, expected, records)
	})
}

func TestTransaction_Size(t *testing.T) {
	t.Run("test size", func(t *testing.T) {
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
		txn, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)

		err = logMgr.Flush()
		require.NoError(t, err)

		logfile := logMgr.LogFileName()
		size, err := txn.Size(logfile)
		require.NoError(t, err)
		require.Equal(t, int32(1), size)
	})
}

func TestTransaction_ExtendFIle(t *testing.T) {
	t.Run("test extend file", func(t *testing.T) {
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
		txn, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)

		err = logMgr.Flush()
		require.NoError(t, err)

		logfile := logMgr.LogFileName()
		size, err := txn.Size(logfile)
		require.NoError(t, err)
		require.Equal(t, int32(1), size)

		_, err = txn.ExtendFile(logfile)
		require.NoError(t, err)

		size2, err := txn.Size(logfile)
		require.NoError(t, err)
		require.Equal(t, int32(2), size2)
	})
}

func TestTransaction_Available(t *testing.T) {
	t.Run("test available", func(t *testing.T) {
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
		txn, err := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)
		require.NoError(t, err)

		nbuf := txn.Available()
		require.Equal(t, numBuf, nbuf)

		blk := *domain.NewBlock(domain.FileName(fake.RandString()), blockSize, domain.BlockNumber(1))
		err = txn.Pin(blk)
		require.NoError(t, err)

		nbuf2 := txn.Available()
		require.Equal(t, numBuf-1, nbuf2)

		txn.Unpin(blk)
		nbuf3 := txn.Available()
		require.Equal(t, numBuf, nbuf3)
	})
}
