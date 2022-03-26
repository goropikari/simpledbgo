package bytes_test

import (
	"io"
	"testing"

	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/goropikari/simpledb_go/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestBuffer(t *testing.T) {
	t.Run("test buffer", func(t *testing.T) {
		buf := bytes.NewBufferBytes([]byte("hello world"))

		x := 5
		actual := make([]byte, x)
		n, err := buf.Read(actual)
		require.NoError(t, err)
		require.Equal(t, x, n)
		require.Equal(t, int64(x), buf.GetOff())
		require.Equal(t, []byte("hello"), actual)

		x = 4
		actual = make([]byte, x)
		n, err = buf.Read(actual)
		require.NoError(t, err)
		require.Equal(t, x, n)
		require.Equal(t, int64(9), buf.GetOff())
		require.Equal(t, []byte(" wor"), actual)

		_, err = buf.Seek(1, io.SeekStart)
		require.NoError(t, err)
		actual = make([]byte, 5)
		n, err = buf.Read(actual)
		require.NoError(t, err)
		require.Equal(t, 5, n)
		require.Equal(t, int64(6), buf.GetOff())
		require.Equal(t, []byte("ello "), actual)

		_, err = buf.Seek(1, io.SeekStart)
		require.NoError(t, err)
		_, err = buf.Write([]byte("1234"))
		require.NoError(t, err)
		_, err = buf.Seek(0, io.SeekStart)
		require.NoError(t, err)
		actual = make([]byte, 5)
		n, err = buf.Read(actual)
		require.NoError(t, err)
		require.Equal(t, 5, n)
		require.Equal(t, int64(5), buf.GetOff())
		require.Equal(t, []byte("h1234"), actual)

		_, err = buf.Seek(0, io.SeekStart)
		require.NoError(t, err)
		n, err = buf.Write([]byte("abcde"))
		require.NoError(t, err)
		require.Equal(t, 5, n)
		require.Equal(t, int64(5), buf.GetOff())

		n, err = buf.Read(actual)
		require.NoError(t, err)
		require.Equal(t, 5, n)
		require.Equal(t, int64(10), buf.GetOff())
		require.Equal(t, []byte(" worl"), actual)

		n, err = buf.Read(actual)
		require.EqualError(t, err, "EOF")
		require.Equal(t, 1, n)
		require.Equal(t, int64(11), buf.GetOff())
		require.Equal(t, []byte("dworl"), actual)
	})
}

func TestBuffer_Read(t *testing.T) {
	tests := []struct {
		name     string
		buf      []byte
		length   int
		expected []byte
	}{
		{
			name:     "valid request",
			buf:      []byte("hello world"),
			length:   5,
			expected: []byte("hello"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferBytes(tt.buf)
			actual := make([]byte, tt.length)

			n, err := buf.Read(actual)

			require.NoError(t, err, "eof")
			require.Equal(t, tt.length, n)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestBuffer_Read_Error(t *testing.T) {
	tests := []struct {
		name      string
		buf       []byte
		sliceSize int
		retSize   int
		errMsg    string
		expected  []byte
	}{
		{
			name:      "valid request: reach EOF",
			buf:       []byte("hello"),
			sliceSize: 10,
			retSize:   5,
			errMsg:    "EOF",
			expected:  []byte{104, 101, 108, 108, 111, 0, 0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferBytes(tt.buf)
			actual := make([]byte, tt.sliceSize)

			n, err := buf.Read(actual)

			require.EqualError(t, err, tt.errMsg)
			require.Equal(t, tt.retSize, n)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestBuffer_Write(t *testing.T) {
	tests := []struct {
		name      string
		sliceSize int
		retSize   int
		given     []byte
		expected  []byte
	}{
		{
			name:      "write short bytes",
			sliceSize: 10,
			retSize:   5,
			given:     []byte("hello"),
			expected:  []byte{104, 101, 108, 108, 111, 0, 0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferBytes(make([]byte, tt.sliceSize))
			n, err := buf.Write(tt.given)
			actual := buf.GetBuf()

			require.NoError(t, err)
			require.Equal(t, tt.retSize, n)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestBuffer_Write_Error(t *testing.T) {
	tests := []struct {
		name      string
		sliceSize int
		retSize   int
		errMsg    string
		given     []byte
		expected  []byte
	}{
		{
			name:      "write overflow",
			sliceSize: 3,
			retSize:   3,
			errMsg:    "EOF",
			given:     []byte("hello"),
			expected:  []byte("hel"),
		},
		{
			name:      "write full size",
			sliceSize: 10,
			retSize:   10,
			errMsg:    "EOF",
			given:     []byte("hellohello"),
			expected:  []byte("hellohello"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferBytes(make([]byte, tt.sliceSize))
			n, err := buf.Write(tt.given)

			actual := buf.GetBuf()

			require.EqualError(t, err, tt.errMsg)
			require.Equal(t, tt.retSize, n)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestBuffer_Seek(t *testing.T) {
	tests := []struct {
		name      string
		seek      int64
		sliceSize int
		whence    int
		expected  int64
	}{
		{
			name:      "seek 3",
			seek:      3,
			sliceSize: 5,
			whence:    io.SeekStart,
			expected:  3,
		},
		{
			name:      "seek 3",
			seek:      3,
			sliceSize: 5,
			whence:    io.SeekCurrent,
			expected:  3,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferBytes(make([]byte, tt.sliceSize))

			n, err := buf.Seek(tt.seek, tt.whence)

			require.NoError(t, err)
			require.Equal(t, tt.expected, n)
			require.Equal(t, tt.expected, buf.GetOff())
		})
	}
}

func TestBuffer_Seek_Error(t *testing.T) {
	tests := []struct {
		name      string
		seek      int64
		sliceSize int
		whence    int
		errMsg    string
		expected  int64
	}{
		{
			name:      "out of range",
			seek:      10,
			sliceSize: 5,
			whence:    io.SeekStart,
			errMsg:    "reference out of range of buffer",
			expected:  0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferBytes(make([]byte, tt.sliceSize))

			n, err := buf.Seek(tt.seek, tt.whence)

			require.EqualError(t, err, tt.errMsg)
			require.Equal(t, tt.expected, n)
			require.Equal(t, tt.expected, buf.GetOff())
		})
	}
}

func TestBuffer_GetInt32(t *testing.T) {
	tests := []struct {
		name     string
		buf      []byte
		expected int32
	}{
		{
			name:     "valid request",
			buf:      []byte{0, 0, 0, 3},
			expected: int32(3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBufferBytes(tt.buf)

			n, err := bb.GetInt32(0)

			require.NoError(t, err)
			require.Equal(t, tt.expected, n)
		})
	}
}

func TestBuffer_GetInt32_Error(t *testing.T) {
	tests := []struct {
		name   string
		buf    []byte
		errMsg string
	}{
		{
			name:   "EOF",
			buf:    []byte{0, 0, 3},
			errMsg: "EOF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBufferBytes(tt.buf)
			_, err := bb.GetInt32(0)

			require.Error(t, err, tt.errMsg)

		})
	}
}

func TestBuffer_SetInt32(t *testing.T) {
	tests := []struct {
		name     string
		buf      []byte
		setint   int32
		expected []byte
	}{
		{
			name:     "valid request",
			buf:      []byte{0, 0, 0, 0},
			setint:   int32(3),
			expected: []byte{0, 0, 0, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBufferBytes(tt.buf)

			err := bb.SetInt32(0, tt.setint)

			require.NoError(t, err)
			require.Equal(t, tt.expected, bb.GetBufferBytes())
		})
	}
}

func TestBuffer_SetInt32_Error(t *testing.T) {
	tests := []struct {
		name   string
		buf    []byte
		setint int32
		errMsg string
	}{
		{
			name:   "not enough space",
			buf:    []byte{0, 0, 0},
			setint: 3,
			errMsg: "invalid offset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBufferBytes(tt.buf)
			err := bb.SetInt32(0, tt.setint)

			require.EqualError(t, err, tt.errMsg)
		})
	}
}

func TestSetGetInt32(t *testing.T) {
	t.Run("test set/get int32", func(t *testing.T) {
		for i := 0; i < 10000; i++ {
			buf := make([]byte, 4)
			bb := bytes.NewBufferBytes(buf)

			x := fake.RandInt32()

			err := bb.SetInt32(0, x)
			require.NoError(t, err)

			n, err := bb.GetInt32(0)
			require.NoError(t, err)
			require.Equal(t, x, n)
		}

		for i := 0; i < 10000; i++ {
			buf := make([]byte, 4)
			bb := bytes.NewBufferBytes(buf)

			x := -fake.RandInt32()

			err := bb.SetInt32(0, x)
			require.NoError(t, err)

			n, err := bb.GetInt32(0)
			require.NoError(t, err)
			require.Equal(t, x, n)
		}
	})
}

func TestBuffer_GetUint32(t *testing.T) {
	tests := []struct {
		name     string
		buf      []byte
		expected uint32
	}{
		{
			name:     "valid request",
			buf:      []byte{0, 0, 0, 3},
			expected: uint32(3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBufferBytes(tt.buf)

			n, err := bb.GetUint32(0)

			require.NoError(t, err)
			require.Equal(t, tt.expected, n)
		})
	}
}

func TestBuffer_GetUint32_Error(t *testing.T) {
	tests := []struct {
		name   string
		buf    []byte
		errMsg string
	}{
		{
			name:   "not enough space",
			buf:    []byte{0, 0, 3},
			errMsg: "invalid offset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBufferBytes(tt.buf)
			_, err := bb.GetUint32(0)

			require.EqualError(t, err, tt.errMsg)
		})
	}
}
func TestBuffer_SetUint32(t *testing.T) {
	tests := []struct {
		name     string
		buf      []byte
		setuint  uint32
		expected []byte
	}{
		{
			name:     "valid request",
			buf:      []byte{0, 0, 0, 0},
			setuint:  3,
			expected: []byte{0, 0, 0, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBufferBytes(tt.buf)

			err := bb.SetUint32(0, tt.setuint)

			require.NoError(t, err)
			require.Equal(t, tt.expected, bb.GetBufferBytes())
		})
	}
}

func TestBuffer_SetUint32_Error(t *testing.T) {
	tests := []struct {
		name    string
		buf     []byte
		setuint uint32
		errMsg  string
	}{
		{
			name:    "invalid request",
			buf:     []byte{0, 0, 0},
			setuint: 3,
			errMsg:  "invalid offset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBufferBytes(tt.buf)

			err := bb.SetUint32(0, tt.setuint)

			require.EqualError(t, err, tt.errMsg)
		})
	}
}

func TestSetGetUint32(t *testing.T) {
	t.Run("test set/get uint32", func(t *testing.T) {
		for i := 0; i < 10000; i++ {
			buf := make([]byte, 4)
			bb := bytes.NewBufferBytes(buf)

			x := fake.RandUint32()

			err := bb.SetUint32(0, x)
			require.NoError(t, err)

			n, err := bb.GetUint32(0)
			require.NoError(t, err)
			require.Equal(t, x, n)
		}
	})
}

func TestBuffer_GetString(t *testing.T) {
	tests := []struct {
		name     string
		buf      []byte
		expected string
	}{
		{
			name:     "valid request",
			buf:      []byte{0, 0, 0, 4, 65, 66, 67, 68, 0, 0},
			expected: "ABCD",
		},
		{
			name:     "valid request: reach EOF",
			buf:      []byte{0, 0, 0, 4, 65, 66, 67, 68},
			expected: "ABCD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBufferBytes(tt.buf)

			s, err := bb.GetString(0)

			require.NoError(t, err)
			require.Equal(t, tt.expected, s)
		})
	}
}
