package buffer_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/goropikari/simpledb_go/backend/buffer"
	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/goropikari/simpledb_go/testing/mock"
	"github.com/stretchr/testify/require"
)

func TestBufferMgr_NewManager(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		const size = 5
		const numbuf = 3

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileMgr := mock.NewMockFileManager(ctrl)
		logMgr := mock.NewMockLogManager(ctrl)

		bsf := mock.NewMockByteSliceFactory(ctrl)
		bsf.EXPECT().Create(gomock.Any()).Return(make([]byte, size), nil).AnyTimes()
		factory := domain.NewPageFactory(bsf, domain.BlockSize(size))

		config := buffer.Config{NumberBuffer: numbuf}
		_, err := buffer.NewManager(fileMgr, logMgr, factory, config)
		require.NoError(t, err)
	})

	t.Run("valid request", func(t *testing.T) {
		const size = 5
		const numbuf = 0

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileMgr := mock.NewMockFileManager(ctrl)
		logMgr := mock.NewMockLogManager(ctrl)

		bsf := mock.NewMockByteSliceFactory(ctrl)
		bsf.EXPECT().Create(gomock.Any()).Return(make([]byte, size), nil).AnyTimes()
		factory := domain.NewPageFactory(bsf, domain.BlockSize(size))

		config := buffer.Config{NumberBuffer: numbuf}
		_, err := buffer.NewManager(fileMgr, logMgr, factory, config)
		require.Error(t, err)
	})
}

func TestBufferMgr_Available(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		const size = 5
		const numbuf = 3

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileMgr := mock.NewMockFileManager(ctrl)
		logMgr := mock.NewMockLogManager(ctrl)

		bsf := mock.NewMockByteSliceFactory(ctrl)
		bsf.EXPECT().Create(gomock.Any()).Return(make([]byte, size), nil).AnyTimes()
		factory := domain.NewPageFactory(bsf, domain.BlockSize(size))

		config := buffer.Config{NumberBuffer: numbuf}
		bufMgr, err := buffer.NewManager(fileMgr, logMgr, factory, config)
		require.NoError(t, err)

		require.Equal(t, numbuf, bufMgr.Available())
	})
}

func TestBufferMgr_Pin(t *testing.T) {
	t.Run("valid request: without timeout", func(t *testing.T) {
		const size = 200
		const numbuf = 3
		dbPath := "dbpath_" + fake.RandString()

		factory := fake.NewNonDirectLogManagerFactory(dbPath, size)
		defer factory.Finish()

		fileMgr, logMgr := factory.Create()

		bsf := bytes.NewByteSliceCreater()
		pageFactory := domain.NewPageFactory(bsf, domain.BlockSize(size))

		config := buffer.Config{
			NumberBuffer:       numbuf,
			TimeoutMillisecond: 500,
		}
		bufMgr, err := buffer.NewManager(fileMgr, logMgr, pageFactory, config)
		require.NoError(t, err)

		filenames := make([]string, 0)
		for i := 0; i < 3; i++ {
			filenames = append(filenames, "file_"+fake.RandString())
		}

		go func() {
			block := domain.NewBlock(domain.FileName(filenames[0]), domain.BlockSize(size), domain.BlockNumber(0))
			buf, err := bufMgr.Pin(block)
			require.NoError(t, err)
			time.Sleep(time.Millisecond * 15)
			bufMgr.Unpin(buf)
		}()

		go func() {
			block := domain.NewBlock(domain.FileName(filenames[1]), domain.BlockSize(size), domain.BlockNumber(0))
			buf, err := bufMgr.Pin(block)
			require.NoError(t, err)
			time.Sleep(time.Millisecond * 20)
			bufMgr.Unpin(buf)
		}()

		go func() {
			block := domain.NewBlock(domain.FileName(filenames[2]), domain.BlockSize(size), domain.BlockNumber(0))
			buf, err := bufMgr.Pin(block)
			require.NoError(t, err)
			time.Sleep(time.Millisecond * 20)
			bufMgr.Unpin(buf)
		}()

		// 先の goroutine よりも後で実行するために sleep
		time.Sleep(time.Millisecond * 10)
		block := domain.NewBlock(domain.FileName("file_"+fake.RandString()), domain.BlockSize(size), domain.BlockNumber(0))
		_, err = bufMgr.Pin(block)
		require.NoError(t, err)
	})

	t.Run("valid request: with timeout", func(t *testing.T) {
		const size = 200
		const numbuf = 3
		dbPath := "dbpath_" + fake.RandString()

		factory := fake.NewNonDirectLogManagerFactory(dbPath, size)
		defer factory.Finish()

		fileMgr, logMgr := factory.Create()

		bsf := bytes.NewByteSliceCreater()
		pageFactory := domain.NewPageFactory(bsf, domain.BlockSize(size))

		config := buffer.Config{
			NumberBuffer:       numbuf,
			TimeoutMillisecond: 20,
		}
		bufMgr, err := buffer.NewManager(fileMgr, logMgr, pageFactory, config)
		require.NoError(t, err)

		filenames := make([]string, 0)
		for i := 0; i < 4; i++ {
			filenames = append(filenames, "file_"+fake.RandString())
		}

		go func() {
			block := domain.NewBlock(domain.FileName(filenames[0]), domain.BlockSize(size), domain.BlockNumber(0))
			_, err := bufMgr.Pin(block)
			require.NoError(t, err)
		}()

		go func() {
			block := domain.NewBlock(domain.FileName(filenames[1]), domain.BlockSize(size), domain.BlockNumber(0))
			_, err := bufMgr.Pin(block)
			require.NoError(t, err)
		}()

		go func() {
			block := domain.NewBlock(domain.FileName(filenames[2]), domain.BlockSize(size), domain.BlockNumber(0))
			_, err := bufMgr.Pin(block)
			require.NoError(t, err)
		}()

		// 先の goroutine よりも後で実行するために sleep
		time.Sleep(time.Millisecond * 10)
		require.Equal(t, 0, bufMgr.Available())
		block := domain.NewBlock(domain.FileName(filenames[3]), domain.BlockSize(size), domain.BlockNumber(0))
		_, err = bufMgr.Pin(block)
		require.EqualError(t, err, "timeout exceeded")
	})
}

func TestBufferMgr_FlushAll(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		const size = 200
		const numbuf = 3
		dbPath := "dbpath_" + fake.RandString()

		factory := fake.NewNonDirectLogManagerFactory(dbPath, size)
		defer factory.Finish()

		fileMgr, logMgr := factory.Create()

		bsf := bytes.NewByteSliceCreater()
		pageFactory := domain.NewPageFactory(bsf, domain.BlockSize(size))

		config := buffer.Config{
			NumberBuffer:       numbuf,
			TimeoutMillisecond: 50,
		}
		bufMgr, err := buffer.NewManager(fileMgr, logMgr, pageFactory, config)
		require.NoError(t, err)

		block := domain.NewBlock(domain.FileName("file_"+fake.RandString()), domain.BlockSize(size), domain.BlockNumber(0))
		buf, err := bufMgr.Pin(block)
		require.NoError(t, err)
		buf.SetModifiedTxNumber(1, 1)

		block2 := domain.NewBlock(domain.FileName("file_"+fake.RandString()), domain.BlockSize(size), domain.BlockNumber(0))
		buf2, err := bufMgr.Pin(block2)
		require.NoError(t, err)
		buf2.SetModifiedTxNumber(1, 2)

		require.Equal(t, buf.TxNumber(), domain.TransactionNumber(1))
		require.Equal(t, buf2.TxNumber(), domain.TransactionNumber(1))

		err = bufMgr.FlushAll(domain.TransactionNumber(1))
		require.NoError(t, err)
		require.Equal(t, buf.TxNumber(), domain.TransactionNumber(-1))
		require.Equal(t, buf2.TxNumber(), domain.TransactionNumber(-1))
	})
}