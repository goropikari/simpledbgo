package tx_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/goropikari/simpledbgo/testing/mock"
	"github.com/goropikari/simpledbgo/tx"
	"github.com/stretchr/testify/require"
)

func TestBufferList_Pin(t *testing.T) {
	tests := []struct {
		name string
		blk  domain.Block
		buf  *domain.Buffer
	}{
		{
			name: "valid request",
			blk:  fake.Block(),
			buf:  fake.Buffer(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			bufMgr := mock.NewMockBufferPoolManager(ctrl)
			bufMgr.EXPECT().Pin(gomock.Any()).Return(tt.buf, nil).AnyTimes()

			bl := tx.NewBufferList(bufMgr)
			err := bl.Pin(tt.blk)
			require.NoError(t, err)

			buf := bl.GetBuffer(tt.blk)
			require.Equal(t, tt.buf, buf)
		})
	}
}

func TestBufferList_Pin_Error(t *testing.T) {
	tests := []struct {
		name string
		blk  domain.Block
		err  error
	}{
		{
			name: "valid request",
			blk:  fake.Block(),
			err:  errors.New("timeout"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			bufMgr := mock.NewMockBufferPoolManager(ctrl)
			bufMgr.EXPECT().Pin(gomock.Any()).Return(nil, tt.err).AnyTimes()

			bl := tx.NewBufferList(bufMgr)
			err := bl.Pin(tt.blk)
			require.Error(t, err)
		})
	}
}

func TestBufferList_Unpin(t *testing.T) {
	tests := []struct {
		name string
		blk  domain.Block
		buf  *domain.Buffer
	}{
		{
			name: "valid request",
			blk:  fake.Block(),
			buf:  fake.Buffer(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			bufMgr := mock.NewMockBufferPoolManager(ctrl)
			bufMgr.EXPECT().Pin(gomock.Any()).Return(tt.buf, nil).AnyTimes()
			bufMgr.EXPECT().Unpin(gomock.Any()).AnyTimes()

			bl := tx.NewBufferList(bufMgr)

			var err error
			err = bl.Pin(tt.blk)
			require.NoError(t, err)
			err = bl.Pin(tt.blk)
			require.NoError(t, err)

			bl.Unpin(tt.blk)
			require.Equal(t, 1, bl.PinnedBlocks().Length())
			bl.Unpin(tt.blk)
			require.Equal(t, 0, bl.PinnedBlocks().Length())
		})
	}
}

func TestBufferList_UnpinAll(t *testing.T) {
	tests := []struct {
		name string
		blk  domain.Block
		buf  *domain.Buffer
	}{
		{
			name: "valid request",
			blk:  fake.Block(),
			buf:  fake.Buffer(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			bufMgr := mock.NewMockBufferPoolManager(ctrl)
			bufMgr.EXPECT().Pin(gomock.Any()).Return(tt.buf, nil).AnyTimes()
			bufMgr.EXPECT().Unpin(gomock.Any()).AnyTimes()

			bl := tx.NewBufferList(bufMgr)

			var err error
			err = bl.Pin(tt.blk)
			require.NoError(t, err)
			err = bl.Pin(tt.blk)
			require.NoError(t, err)

			require.Equal(t, 2, bl.PinnedBlocks().Length())

			bl.UnpinAll()
			require.Equal(t, 0, bl.PinnedBlocks().Length())
		})
	}
}
