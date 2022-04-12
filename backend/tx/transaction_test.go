package tx_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/goropikari/simpledb_go/backend/tx"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/assert"
)

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
		transaction := tx.NewTransaction(fileMgr, logMgr, bufMgr, concurMgr, gen)

		transaction.Commit()
		transaction.Commit()

		logFileName := logMgr.LogFileName()
		f, _ := os.OpenFile(filepath.Join(dbPath, string(logFileName)), os.O_RDONLY, os.ModePerm)
		data, _ := io.ReadAll(f)
		expected := []byte{
			0, 0, 0, 8, 0, 0, 0, 0,
			0, 0, 0, 8, 0, 0, 0, 0,
			0, 0, 0, 2,
			0, 0, 0, 8, 0, 0, 0, 0,
			0, 0, 0, 8, 0, 0, 0, 0,
			0, 0, 0, 2,
		}

		assert.Equal(t, expected, data)
	})
}
