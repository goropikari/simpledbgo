package domain_test

import (
	"io"
	goos "os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/log"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/goropikari/simpledbgo/testing/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestBuffer_NewBuffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("valid request", func(t *testing.T) {
		fileMgr := mock.NewMockFileManager(ctrl)
		fileMgr.EXPECT().CreatePage().Return(&domain.Page{}, nil).AnyTimes()
		logMgr := mock.NewMockLogManager(ctrl)
		_, err := domain.NewBuffer(fileMgr, logMgr)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		fileMgr := mock.NewMockFileManager(ctrl)
		fileMgr.EXPECT().CreatePage().Return(nil, errors.New("error")).AnyTimes()
		logMgr := mock.NewMockLogManager(ctrl)
		_, err := domain.NewBuffer(fileMgr, logMgr)
		require.Error(t, err)
	})
}

func TestBuffer_Block(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileMgr := mock.NewMockFileManager(ctrl)
	fileMgr.EXPECT().CreatePage().Return(&domain.Page{}, nil).AnyTimes()
	logMgr := mock.NewMockLogManager(ctrl)

	t.Run("valid request", func(t *testing.T) {
		buf, err := domain.NewBuffer(fileMgr, logMgr)
		require.NoError(t, err)
		require.Equal(t, domain.Block{}, buf.Block())
	})
}

func TestBuffer_SetModifiedTxNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileMgr := mock.NewMockFileManager(ctrl)
	fileMgr.EXPECT().CreatePage().Return(&domain.Page{}, nil).AnyTimes()
	logMgr := mock.NewMockLogManager(ctrl)

	t.Run("valid request", func(t *testing.T) {
		buf, err := domain.NewBuffer(fileMgr, logMgr)
		require.NoError(t, err)

		txnum := domain.TransactionNumber(fake.RandInt32())
		lsn := domain.LSN(fake.RandInt32())
		buf.SetModifiedTxNumber(txnum, lsn)

		require.Equal(t, txnum, buf.TxNumber())
		require.Equal(t, lsn, buf.LSN())
	})
}

func TestBuffer_PinUnpin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileMgr := mock.NewMockFileManager(ctrl)
	fileMgr.EXPECT().CreatePage().Return(&domain.Page{}, nil).AnyTimes()
	logMgr := mock.NewMockLogManager(ctrl)

	t.Run("valid request", func(t *testing.T) {
		buf, err := domain.NewBuffer(fileMgr, logMgr)
		require.NoError(t, err)

		require.Equal(t, false, buf.IsPinned())

		buf.Pin()
		require.Equal(t, true, buf.IsPinned())

		buf.Unpin()
		require.Equal(t, false, buf.IsPinned())
	})
}

func TestBuffer_AssignToBlock(t *testing.T) {
	const size = 10

	t.Run("valid request", func(t *testing.T) {
		dbPath := fake.RandString()
		blockSize := int32(size)

		factory := fake.NewNonDirectLogManagerFactory(dbPath, blockSize)
		defer factory.Finish()
		fileMgr, logMgr := factory.Create()

		buf, err := domain.NewBuffer(fileMgr, logMgr)
		require.NoError(t, err)

		fileName := fake.RandString()
		f, _ := goos.OpenFile(filepath.Join(dbPath, fileName), goos.O_RDWR|goos.O_CREATE, goos.ModePerm)
		f.Write(make([]byte, blockSize))
		f.Seek(0, io.SeekStart)
		f.Write([]byte("hello"))
		f.Close()

		block := domain.NewBlock(domain.FileName(fileName), domain.BlockNumber(0))
		buf.AssignToBlock(block)
		expected := make([]byte, blockSize)
		copy(expected, []byte("hello"))
		require.Equal(t, expected, buf.Page().GetData())
	})

	t.Run("valid request: flush lsn", func(t *testing.T) {
		dbPath := fake.RandString()
		blockSize := int32(size)
		factory := fake.NewNonDirectFileManagerFactory(dbPath, blockSize)
		defer factory.Finish()
		fileMgr := factory.Create()

		logFileName := fake.RandString()
		logConfig := log.ManagerConfig{LogFileName: logFileName}
		logMgr, err := log.NewManager(fileMgr, logConfig)
		require.NoError(t, err)

		fileName := fake.RandString()
		block := domain.NewBlock(domain.FileName(fileName), domain.BlockNumber(0))

		buf, err := domain.NewBuffer(fileMgr, logMgr)
		require.NoError(t, err)

		f, _ := goos.OpenFile(filepath.Join(dbPath, fileName), goos.O_RDWR|goos.O_CREATE, goos.ModePerm)
		f.Write(make([]byte, blockSize))
		f.Seek(0, io.SeekStart)
		f.Write([]byte("hello"))
		f.Close()

		buf.AssignToBlock(block)
		buf.SetModifiedTxNumber(1, 1)
		buf.AssignToBlock(block)
		expected := make([]byte, blockSize)
		copy(expected, []byte("hello"))
		require.Equal(t, expected, buf.Page().GetData())
	})
}
