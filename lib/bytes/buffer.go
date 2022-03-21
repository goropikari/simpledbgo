package bytes

import (
	"errors"
	"io"
)

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/tests/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// ByteBuffer is an interface that implement io.ReadWriteSeeker.
type ByteBuffer interface {
	io.ReadWriteSeeker
	GetBytes() []byte
	Reset()
}

var (
	// ErrOutOfRange is an error type that refer out of range.
	ErrOutOfRange = errors.New("reference out of range of buffer")

	// ErrUnsupportedWhence is an error type that given whence is unsupported.
	ErrUnsupportedWhence = errors.New("unsupported whence")
)

// Buffer is a buffer.
type Buffer struct {
	capacity int
	buf      []byte
	off      int64
}

// NewBuffer is a constructor of Buffer.
func NewBuffer(n int) (*Buffer, error) {
	return NewBufferBytes(make([]byte, n)), nil
}

// NewBufferBytes is a constructor of Buffer by byte slice.
func NewBufferBytes(buf []byte) *Buffer {
	return &Buffer{
		capacity: len(buf),
		buf:      buf,
		off:      0,
	}
}

// Read reads bytes from Reader.
func (buf *Buffer) Read(p []byte) (n int, err error) {
	if buf.off < 0 || buf.off >= int64(buf.capacity) {
		return 0, io.EOF
	}

	cnt := copy(p, buf.buf[buf.off:])

	buf.off += int64(cnt)
	if int(buf.off) == buf.capacity {
		return cnt, io.EOF
	}

	return cnt, nil
}

// Write writes given bytes to writer.
func (buf *Buffer) Write(p []byte) (n int, err error) {
	if buf.off < 0 || buf.off >= int64(buf.capacity) {
		return 0, io.EOF
	}

	cnt := copy(buf.buf[buf.off:], p)

	buf.off += int64(cnt)
	if buf.off == int64(buf.capacity) {
		return cnt, io.EOF
	}

	return cnt, nil
}

// Seek seeks position.
func (buf *Buffer) Seek(offset int64, whence int) (int64, error) {
	off := int64(0)

	switch whence {
	case io.SeekStart:
		off = offset
	case io.SeekCurrent:
		off += offset
	default:
		return 0, ErrUnsupportedWhence
	}
	if off < 0 || off > int64(buf.capacity) {
		return 0, ErrOutOfRange
	}
	buf.off = off

	return off, nil
}

// GetBytes returns buffer.
func (buf *Buffer) GetBytes() []byte {
	return buf.buf
}

// Reset resets buffer.
func (buf *Buffer) Reset() {
	buf.off = 0
	for i := 0; i < buf.capacity; i++ {
		buf.buf[i] = 0
	}
}
