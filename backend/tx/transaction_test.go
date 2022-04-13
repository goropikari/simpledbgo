package tx_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/backend/tx"
	"github.com/goropikari/simpledb_go/backend/tx/logrecord"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/assert"
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

		logFileName := logMgr.LogFileName()
		f, _ := os.OpenFile(filepath.Join(dbPath, string(logFileName)), os.O_RDONLY, os.ModePerm)
		data, _ := io.ReadAll(f)
		expected := []byte{
			0, 0, 0, 8, 0, 0, 0, 0,
			0, 0, 0, 8, 0, 0, 0, byte(tx.Start),
			0, 0, 0, 1,
		}

		assert.Equal(t, expected, data)
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

		transaction.Commit()
		transaction.Commit()

		txnum := 2

		logFileName := logMgr.LogFileName()
		f, _ := os.OpenFile(filepath.Join(dbPath, string(logFileName)), os.O_RDONLY, os.ModePerm)
		data, _ := io.ReadAll(f)
		expected := []byte{
			0, 0, 0, 8, 0, 0, 0, 0,
			0, 0, 0, 8, 0, 0, 0, byte(tx.Start),
			0, 0, 0, byte(txnum),
			0, 0, 0, 8, 0, 0, 0, 0,
			0, 0, 0, 8, 0, 0, 0, byte(tx.Commit),
			0, 0, 0, byte(txnum),
			0, 0, 0, 8, 0, 0, 0, 0,
			0, 0, 0, 8, 0, 0, 0, byte(tx.Commit),
			0, 0, 0, byte(txnum),
		}

		assert.Equal(t, expected, data)
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

		logFileName := logMgr.LogFileName()
		f, _ := os.OpenFile(filepath.Join(dbPath, string(logFileName)), os.O_RDONLY, os.ModePerm)
		data, _ := io.ReadAll(f)
		rec := &logrecord.SetInt32Record{}
		err = rec.Unmarshal(data[30:54])
		require.NoError(t, err)

		expected := logrecord.SetInt32Record{
			FileName:    blk.FileName(),
			TxNum:       1,
			BlockNumber: 0,
			Offset:      offset,
			Val:         val,
		}
		require.Equal(t, expected, *rec)
	})
}
