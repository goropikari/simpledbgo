package buffer_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/goropikari/simpledb_go/backend/buffer"
	"github.com/goropikari/simpledb_go/backend/core"
	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/goropikari/simpledb_go/testing/mock"
	"github.com/stretchr/testify/require"
)

func TestBuffer_pin_unpin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fm := mock.NewMockFileManager(ctrl)
	lm := mock.NewMockLogManager(ctrl)

	bb := bytes.NewBuffer(100)
	page := core.NewPage(bb)

	fm.EXPECT().PreparePage().Return(page, nil).AnyTimes()
	fm.EXPECT().CopyBlockToPage(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	fm.EXPECT().CopyPageToBlock(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	lm.EXPECT().FlushByLSN(gomock.Any()).Return(nil).AnyTimes()

	var tests = []struct {
		name     string
		niter    int
		expected int
	}{
		{
			name:     "valid request",
			niter:    3,
			expected: 3,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf, err := buffer.NewBuffer(fm, lm)
			require.NoError(t, err)

			for i := 0; i < tt.niter; i++ {
				buf.Pin()
			}

			require.Equal(t, tt.expected, buf.GetPins())

			for i := 0; i < tt.niter; i++ {
				buf.Unpin()
			}

			require.Equal(t, 0, buf.GetPins())
		})
	}
}

func TestBuffer_setModified(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fm := mock.NewMockFileManager(ctrl)
	lm := mock.NewMockLogManager(ctrl)

	bb := mock.NewMockByteBuffer(ctrl)
	page := core.NewPage(bb)

	fm.EXPECT().PreparePage().Return(page, nil).AnyTimes()

	var tests = []struct {
		name  string
		txnum int32
		lsn   int32
	}{
		{
			name:  "valid request",
			txnum: fake.RandInt32(),
			lsn:   fake.RandInt32(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf, err := buffer.NewBuffer(fm, lm)
			require.NoError(t, err)

			buf.SetModified(tt.txnum, tt.lsn)

			require.Equal(t, buf.GetTxNum(), tt.txnum)
			require.Equal(t, buf.GetLSN(), tt.lsn)
		})
	}
}

func TestBuffer_AssignToBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fm := mock.NewMockFileManager(ctrl)
	lm := mock.NewMockLogManager(ctrl)

	bb := mock.NewMockByteBuffer(ctrl)
	page := core.NewPage(bb)

	fm.EXPECT().PreparePage().Return(page, nil).AnyTimes()
	fm.EXPECT().CopyBlockToPage(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	var tests = []struct {
		name string
		blk  *core.Block
	}{
		{
			name: "valid request",
			blk:  fake.Block(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf, err := buffer.NewBuffer(fm, lm)
			require.NoError(t, err)

			err = buf.AssignToBlock(tt.blk)

			require.NoError(t, err)
			require.Equal(t, tt.blk, buf.GetBlock())
		})
	}
}

func TestBuffer_AssignToBlock_Error(t *testing.T) {

	var tests = []struct {
		name        string
		blk         *core.Block
		txnum       int32
		errBtoP     error
		errPtoB     error
		errFlushLSN error
		errMsg      string
	}{
		{
			name:    "error at file manager: block to page",
			blk:     fake.Block(),
			txnum:   -1,
			errBtoP: errors.New("error block to page"),
			errMsg:  "error block to page",
		},
		{
			name:    "error at file manager: page to block",
			blk:     fake.Block(),
			txnum:   1,
			errPtoB: errors.New("error page to block"),
			errMsg:  "error page to block",
		},
		{
			name:        "error at log manager",
			blk:         fake.Block(),
			txnum:       1,
			errFlushLSN: errors.New("error flush by lsn"),
			errMsg:      "error flush by lsn",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fm := mock.NewMockFileManager(ctrl)
			lm := mock.NewMockLogManager(ctrl)

			bb := mock.NewMockByteBuffer(ctrl)
			page := core.NewPage(bb)

			fm.EXPECT().PreparePage().Return(page, nil).AnyTimes()
			fm.EXPECT().CopyBlockToPage(gomock.Any(), gomock.Any()).Return(tt.errBtoP).AnyTimes()
			fm.EXPECT().CopyPageToBlock(gomock.Any(), gomock.Any()).Return(tt.errPtoB).AnyTimes()

			lm.EXPECT().FlushByLSN(tt.txnum).Return(tt.errFlushLSN).AnyTimes()

			buf, err := buffer.NewBuffer(fm, lm)
			require.NoError(t, err)

			buf.SetModified(tt.txnum, -1)

			err = buf.AssignToBlock(tt.blk)

			require.EqualError(t, err, tt.errMsg)
		})
	}
}
