package bytes_test

import (
	"io"
	"testing"

	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/stretchr/testify/require"
)

func TestBuffer_Read(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		buf := bytes.NewBufferBytes([]byte("hello"))
		actual := make([]byte, 3)
		n, err := buf.Read(actual)

		require.NoError(t, err)
		require.Equal(t, 3, n)
		require.Equal(t, []byte("hel"), actual)
	})

	t.Run("reach EOF", func(t *testing.T) {
		buf := bytes.NewBufferBytes([]byte("hello"))
		actual := make([]byte, 6)
		n, err := buf.Read(actual)

		require.NoError(t, err)
		require.Equal(t, 5, n)
		require.Equal(t, append([]byte("hello"), 0), actual)

		n2, err := buf.Read(actual)
		require.Error(t, err)
		require.Equal(t, 0, n2)
	})
}

func TestBuffer_Write(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		buf := bytes.NewBuffer(5)
		n, err := buf.Write([]byte("hello"))
		require.NoError(t, err)
		require.Equal(t, 5, n)
	})

	t.Run("overflow", func(t *testing.T) {
		buf := bytes.NewBuffer(5)
		n, err := buf.Write([]byte("hello"))
		require.NoError(t, err)
		require.Equal(t, 5, n)

		n2, err := buf.Write([]byte("hello"))
		require.Error(t, err)
		require.Equal(t, 0, n2)
	})
}

func TestBuffer_Seek(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		buf := bytes.NewBuffer(5)
		n, err := buf.Seek(2, io.SeekStart)
		require.NoError(t, err)
		require.Equal(t, int64(2), n)
	})

	t.Run("invalid request", func(t *testing.T) {
		buf := bytes.NewBuffer(5)
		// Unsupported whence
		n, err := buf.Seek(2, io.SeekCurrent)
		require.Error(t, err)
		require.Equal(t, int64(0), n)

		// out of range
		n2, err := buf.Seek(10, io.SeekStart)
		require.Error(t, err)
		require.Equal(t, int64(0), n2)
	})
}

func TestBuffer_GetData(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		data := []byte("hello")
		buf := bytes.NewBufferBytes(data)
		require.Equal(t, data, buf.GetData())
	})
}

func TestBuffer_Reset(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		data := []byte("hello")
		buf := bytes.NewBufferBytes(data)
		require.Equal(t, data, buf.GetData())

		buf.Reset()
		require.Equal(t, make([]byte, 5), buf.GetData())
	})
}

func TestBuffer_GetInt32(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		data := []byte{0, 0, 0, 9}
		buf := bytes.NewBufferBytes(data)
		n, err := buf.GetInt32(0)
		require.NoError(t, err)
		require.Equal(t, int32(9), n)
	})

	t.Run("invalid request", func(t *testing.T) {
		data := []byte{0, 0, 9}
		buf := bytes.NewBufferBytes(data)
		n, err := buf.GetInt32(0)
		require.Error(t, err)
		require.Equal(t, int32(0), n)
	})
}

func TestBuffer_SetInt32(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		buf := bytes.NewBuffer(4)
		err := buf.SetInt32(0, 9)
		require.NoError(t, err)
		require.Equal(t, []byte{0, 0, 0, 9}, buf.GetData())
	})

	t.Run("invalid request", func(t *testing.T) {
		buf := bytes.NewBuffer(3)
		err := buf.SetInt32(0, 9)
		require.Error(t, err)
	})
}

func TestBuffer_GetUint32(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		data := []byte{0, 0, 0, 9}
		buf := bytes.NewBufferBytes(data)
		n, err := buf.GetUint32(0)
		require.NoError(t, err)
		require.Equal(t, uint32(9), n)
	})

	t.Run("invalid request", func(t *testing.T) {
		data := []byte{0, 0, 9}
		buf := bytes.NewBufferBytes(data)
		n, err := buf.GetUint32(0)
		require.Error(t, err)
		require.Equal(t, uint32(0), n)
	})
}

func TestBuffer_SetUint32(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		buf := bytes.NewBuffer(4)
		err := buf.SetUint32(0, 9)
		require.NoError(t, err)
		require.Equal(t, []byte{0, 0, 0, 9}, buf.GetData())
	})

	t.Run("invalid request", func(t *testing.T) {
		buf := bytes.NewBuffer(3)
		err := buf.SetUint32(0, 9)
		require.Error(t, err)
	})
}

func TestBuffer_GetString(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		str := "hello"
		data := append([]byte{0, 0, 0, 5}, []byte(str)...)
		buf := bytes.NewBufferBytes(data)
		s, err := buf.GetString(0)

		require.NoError(t, err)
		require.Equal(t, str, s)
	})

	t.Run("invalid request", func(t *testing.T) {
		str := "hello"
		data := append([]byte{0, 0, 0, 5}, []byte(str)...)
		buf := bytes.NewBufferBytes(data)

		s, err := buf.GetString(-1)
		require.Error(t, err)
		require.Equal(t, "", s)

		s2, err := buf.GetString(8)
		require.Error(t, err)
		require.Equal(t, "", s2)

		s3, err := buf.GetString(1)
		require.Error(t, err)
		require.Equal(t, "", s3)
	})
}

func TestBuffer_SetString(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		str := "hello"
		buf := bytes.NewBuffer(9)
		err := buf.SetString(0, str)
		require.NoError(t, err)
		require.Equal(t, append([]byte{0, 0, 0, 5}, []byte("hello")...), buf.GetData())
	})

	t.Run("invalid request", func(t *testing.T) {
		str := "hello"
		buf := bytes.NewBuffer(9)
		err := buf.SetString(-1, str)
		require.Error(t, err)

		err = buf.SetString(1, str)
		require.Error(t, err)
	})
}

func TestBuffer_GetBytes(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		str := "hello"
		data := append([]byte{0, 0, 0, 5}, []byte(str)...)
		buf := bytes.NewBufferBytes(data)
		b, err := buf.GetBytes(0)

		require.NoError(t, err)
		require.Equal(t, []byte(str), b)
	})

	t.Run("invalid request", func(t *testing.T) {
		str := "hello"
		data := append([]byte{0, 0, 0, 5}, []byte(str)...)
		buf := bytes.NewBufferBytes(data)

		b, err := buf.GetBytes(-1)
		require.Error(t, err)
		require.Equal(t, []byte(nil), b)

		b2, err := buf.GetBytes(8)
		require.Error(t, err)
		require.Equal(t, []byte(nil), b2)

		b3, err := buf.GetBytes(1)
		require.Error(t, err)
		require.Equal(t, []byte(nil), b3)
	})
}

func TestBuffer_SetBytes(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		b := []byte("hello")
		buf := bytes.NewBuffer(9)
		err := buf.SetBytes(0, b)
		require.NoError(t, err)
		require.Equal(t, append([]byte{0, 0, 0, 5}, b...), buf.GetData())
	})

	t.Run("invalid request", func(t *testing.T) {
		b := []byte("hello")
		buf := bytes.NewBuffer(9)
		err := buf.SetBytes(-1, b)
		require.Error(t, err)

		err = buf.SetBytes(1, b)
		require.Error(t, err)
	})
}
