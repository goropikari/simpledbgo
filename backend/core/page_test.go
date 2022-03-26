package core_test

import (
	"errors"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/goropikari/simpledb_go/backend/core"
	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/goropikari/simpledb_go/testing/mock"
	"github.com/stretchr/testify/require"
)

func TestPage_SetGetInt32(t *testing.T) {
	x1 := fake.RandInt32()
	x2 := -fake.RandInt32()

	tests := []struct {
		name    string
		offset  int64
		bufsize int
		setint  int32
	}{
		{
			name:    "valid request: set positive number",
			offset:  0,
			bufsize: 4,
			setint:  x1,
		},
		{
			name:    "valid request: set negative number",
			offset:  0,
			bufsize: 4,
			setint:  x2,
		},
		{
			name:    "valid request: non zero offset",
			offset:  4,
			bufsize: 12,
			setint:  x2,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBuffer(tt.bufsize)
			page := core.NewPage(bb)

			err := page.SetInt32(tt.offset, tt.setint)
			require.NoError(t, err)

			n, err := page.GetInt32(tt.offset)
			require.NoError(t, err)
			require.Equal(t, tt.setint, n)
		})
	}
}

func TestPage_SetGetInt32_Error(t *testing.T) {
	tests := []struct {
		name    string
		offset  int64
		bufsize int
		setint  int32
		errMsg  string
	}{
		{
			name:    "invalid request: invalid offset",
			offset:  1,
			bufsize: 4,
			setint:  fake.RandInt32(),
			errMsg:  "invalid offset",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBuffer(tt.bufsize)
			page := core.NewPage(bb)

			err := page.SetInt32(tt.offset, tt.setint)
			require.EqualError(t, err, tt.errMsg)

			_, err = page.GetInt32(tt.offset)
			require.EqualError(t, err, tt.errMsg)
		})
	}
}

func TestPage_SetGetUint32(t *testing.T) {
	x1 := fake.RandUint32()

	tests := []struct {
		name    string
		offset  int64
		bufsize int
		setint  uint32
	}{
		{
			name:    "valid request: set positive number",
			offset:  0,
			bufsize: 4,
			setint:  x1,
		},
		{
			name:    "valid request: non zero offset",
			offset:  4,
			bufsize: 12,
			setint:  x1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBuffer(tt.bufsize)
			page := core.NewPage(bb)

			err := page.SetUint32(tt.offset, tt.setint)
			require.NoError(t, err)

			n, err := page.GetUint32(tt.offset)
			require.NoError(t, err)
			require.Equal(t, tt.setint, n)
		})
	}
}

func TestPage_SetGetUint32_Error(t *testing.T) {
	tests := []struct {
		name    string
		offset  int64
		bufsize int
		setint  uint32
		errMsg  string
	}{
		{
			name:    "invalid request: invalid offset",
			offset:  1,
			bufsize: 4,
			setint:  fake.RandUint32(),
			errMsg:  "invalid offset",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBuffer(tt.bufsize)
			page := core.NewPage(bb)

			err := page.SetUint32(tt.offset, tt.setint)
			require.EqualError(t, err, tt.errMsg)

			_, err = page.GetUint32(tt.offset)
			require.EqualError(t, err, tt.errMsg)
		})
	}
}

func TestPage_SetGetBytes(t *testing.T) {
	var tests = []struct {
		name     string
		bufsize  int
		setbytes []byte
		offset   int64
	}{
		{
			name:     "valid request",
			bufsize:  8,
			setbytes: []byte("hoge"),
			offset:   0,
		},
		{
			name:     "valid request: non zero offset",
			bufsize:  20,
			setbytes: []byte("hoge"),
			offset:   4,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBuffer(tt.bufsize)
			page := core.NewPage(bb)

			err := page.SetBytes(tt.offset, tt.setbytes)
			require.NoError(t, err)

			b, err := page.GetBytes(tt.offset)
			require.NoError(t, err)
			require.Equal(t, tt.setbytes, b)
		})
	}
}

func TestPage_SetBytes_Error(t *testing.T) {
	tests := []struct {
		name     string
		bufsize  int
		setbytes []byte
		offset   int64
		errMsg   string
	}{
		{
			name:     "invalid request: no enough space",
			bufsize:  8,
			setbytes: []byte("hoge"),
			offset:   1,
			errMsg:   "invalid offset",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBuffer(tt.bufsize)
			page := core.NewPage(bb)

			err := page.SetBytes(tt.offset, tt.setbytes)
			require.EqualError(t, err, tt.errMsg)
		})
	}
}

func TestPage_GetBytes_Error(t *testing.T) {
	tests := []struct {
		name   string
		buf    []byte
		offset int64
		errMsg string
	}{
		{
			name:   "invalid request: no enough space",
			buf:    append([]byte{0, 0, 0, 4}, []byte("ABCD")...),
			offset: 1,
			errMsg: "invalid offset",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBufferBytes(tt.buf)
			page := core.NewPage(bb)

			_, err := page.GetBytes(tt.offset)
			require.EqualError(t, err, tt.errMsg)
		})
	}
}

func TestPage_SetGetString(t *testing.T) {
	var tests = []struct {
		name      string
		bufsize   int
		setstring string
		offset    int64
	}{
		{
			name:      "valid request",
			bufsize:   10,
			setstring: "hoge",
			offset:    0,
		},
		{
			name:      "full size",
			bufsize:   8,
			setstring: "hoge",
			offset:    0,
		},
		{
			name:      "non zero offset",
			bufsize:   12,
			setstring: "hoge",
			offset:    4,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBuffer(tt.bufsize)
			page := core.NewPage(bb)

			err := page.SetString(tt.offset, tt.setstring)
			require.NoError(t, err, "set string")

			b, err := page.GetString(tt.offset)
			require.NoError(t, err, "get string")
			require.Equal(t, tt.setstring, b)
		})
	}
}

func TestPage_GetBufferBytes(t *testing.T) {
	var tests = []struct {
		name     string
		expected []byte
	}{
		{
			name:     "valid request",
			expected: []byte("hello"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run("test GetBufferBytes", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			bb := mock.NewMockByteBuffer(ctrl)
			bb.EXPECT().GetBufferBytes().Return(tt.expected).AnyTimes()

			page := core.NewPage(bb)

			b := page.GetBufferBytes()
			require.Equal(t, b, tt.expected)
		})
	}

}

func TestPage_Reset(t *testing.T) {
	t.Run("test Reset", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bb := mock.NewMockByteBuffer(ctrl)
		bb.EXPECT().Reset().AnyTimes()

		page := core.NewPage(bb)

		page.Reset()
	})
}

func TestPage_Write(t *testing.T) {
	tests := []struct {
		name      string
		bufsize   int
		writebyte []byte
	}{
		{
			name:      "write page",
			bufsize:   10,
			writebyte: []byte("hello"),
		},
		{
			name:      "write full size",
			bufsize:   5,
			writebyte: []byte("hello"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBuffer(tt.bufsize)
			page := core.NewPage(bb)

			_, err := page.Write(tt.writebyte)
			require.NoError(t, err)
		})
	}
}

func TestPage_Write_Error(t *testing.T) {
	tests := []struct {
		name   string
		bytes  []byte
		err    error
		errMsg string
	}{
		{
			name:   "write page",
			bytes:  []byte{1, 2, 3},
			err:    errors.New("unexpected error"),
			errMsg: "unexpected error",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			bb := mock.NewMockByteBuffer(ctrl)
			bb.EXPECT().Write(gomock.Any()).Return(0, tt.err).AnyTimes()

			page := core.NewPage(bb)

			_, err := page.Write(tt.bytes)
			require.EqualError(t, err, tt.errMsg)
		})
	}
}

func TestPage_Seek(t *testing.T) {
	var tests = []struct {
		name   string
		offset int64
	}{
		{
			name:   "valid request",
			offset: fake.RandInt64(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			bb := mock.NewMockByteBuffer(ctrl)
			bb.EXPECT().Seek(gomock.Any(), gomock.Any()).Return(tt.offset, nil).AnyTimes()

			page := core.NewPage(bb)

			_, err := page.Seek(tt.offset, io.SeekStart)
			require.NoError(t, err)
		})
	}
}

func TestPage_Seek_Error(t *testing.T) {
	var tests = []struct {
		name   string
		err    error
		errMsg string
	}{
		{
			name:   "valid request",
			err:    errors.New("unexpected error"),
			errMsg: "unexpected error",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			bb := mock.NewMockByteBuffer(ctrl)
			bb.EXPECT().Seek(gomock.Any(), gomock.Any()).Return(fake.RandInt64(), tt.err).AnyTimes()

			page := core.NewPage(bb)

			_, err := page.Seek(0, io.SeekStart)
			require.EqualError(t, err, tt.errMsg)
		})
	}
}
