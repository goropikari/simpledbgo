package bytes

import (
	"errors"
	"io"
)

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/tests/mock/mock_${GOFILE} -package=mock

var (
	// OutOfRangeError is an error type that refer out of range.
	OutOfRangeError = errors.New("reference out of range of buffer")

	// UnsupportedWhenceError is an error type that given whence is unsupported.
	UnsupportedWhenceError = errors.New("unsupported whence")
)

// ByteBuffer is an interface that implement io.ReadWriteSeeker.
type ByteBuffer interface {
	io.ReadWriteSeeker
	GetBytes() []byte
}

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
func (bb *Buffer) Read(p []byte) (n int, err error) {
	if bb.off < 0 || bb.off >= int64(bb.capacity) {
		return 0, io.EOF
	}

	cnt := copy(p, bb.buf[bb.off:])

	bb.off += int64(cnt)
	if int(bb.off) == bb.capacity {
		return cnt, io.EOF
	}

	return cnt, nil
}

// Write writes given bytes to writer.
func (bb *Buffer) Write(p []byte) (n int, err error) {
	if bb.off < 0 || bb.off >= int64(bb.capacity) {
		return 0, io.EOF
	}

	cnt := copy(bb.buf[bb.off:], p)

	bb.off += int64(cnt)
	if bb.off == int64(bb.capacity) {
		return cnt, io.EOF
	}

	return cnt, nil
}

// Seek seeks position.
func (bb *Buffer) Seek(offset int64, whence int) (int64, error) {
	off := int64(0)

	switch whence {
	case io.SeekStart:
		off = offset
	case io.SeekCurrent:
		off += offset
	default:
		return 0, UnsupportedWhenceError
	}
	if off < 0 || off > int64(bb.capacity) {
		return 0, OutOfRangeError
	}
	bb.off = off

	return off, nil
}

// GetBytes returns buffer.
func (bb *Buffer) GetBytes() []byte {
	return bb.buf
}
