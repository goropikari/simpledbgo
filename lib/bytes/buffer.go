package bytes

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/tests/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// ErrGetFromBuffer is an error that means failed to get object from buffer.
var ErrGetFromBuffer = errors.New("failed to get from buffer")

var endianness = binary.BigEndian

// ByteBuffer is an interface that implement io.ReadWriteSeeker.
type ByteBuffer interface {
	io.ReadWriteSeeker
	GetFullBytes() []byte
	GetInt32(int64) (int32, error)
	SetInt32(int64, int32) error
	GetUint32(int64) (uint32, error)
	SetUint32(int64, uint32) error
	GetString(int64) (string, error)
	SetString(int64, string) error
	GetBytes(int64) ([]byte, error)
	SetBytes(int64, []byte) error
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

// GetFullBytes returns buffer.
func (buf *Buffer) GetFullBytes() []byte {
	return buf.buf
}

// Reset resets buffer.
func (buf *Buffer) Reset() {
	buf.off = 0
	for i := 0; i < buf.capacity; i++ {
		buf.buf[i] = 0
	}
}

// GetInt32 returns int32 from buffer.
func (buf *Buffer) GetInt32(offset int64) (int32, error) {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	var ret int32
	if err := binary.Read(buf, endianness, &ret); err != nil {
		if errors.Is(err, io.EOF) {
			return ret, fmt.Errorf("%w", err)
		}

		return 0, fmt.Errorf("%w", err)
	}

	return ret, nil
}

// SetInt32 returns int32 from buffer.
// --------------------
// |  int32 (4 bytes) |
// --------------------.
func (buf *Buffer) SetInt32(offset int64, x int32) error {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := binary.Write(buf, endianness, x); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// GetUint32 returns uint32 from buffer.
func (buf *Buffer) GetUint32(offset int64) (uint32, error) {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	var ret uint32
	if err := binary.Read(buf, endianness, &ret); err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	return ret, nil
}

// SetUint32 returns uint32 from buffer.
// --------------------
// | uint32 (4 bytes) |
// --------------------.
func (buf *Buffer) SetUint32(offset int64, x uint32) error {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := binary.Write(buf, endianness, x); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// GetString returns string from buffer.
func (buf *Buffer) GetString(offset int64) (string, error) {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return "", fmt.Errorf("failed to get string: %w", err)
	}

	length, err := buf.GetUint32(offset)
	if err != nil {
		return "", err
	}

	bytes := make([]byte, length)
	readLen, err := buf.Read(bytes)

	if errors.Is(err, io.EOF) {
		return "", fmt.Errorf("%w", err)
	} else if err != nil {
		return "", err
	}

	if uint32(readLen) != length {
		return "", ErrGetFromBuffer
	}

	return string(bytes), err
}

// SetString returns string from buffer.
// ----------------------------------------
// | string length (uint32)| body (string)|
// ----------------------------------------.
func (buf *Buffer) SetString(offset int64, str string) error {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := buf.SetUint32(offset, uint32(len(str))); err != nil {
		return fmt.Errorf("failed to set string: %w", err)
	}

	if _, err := buf.Write([]byte(str)); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// GetBytes returns bytes from page.
func (buf *Buffer) GetBytes(offset int64) ([]byte, error) {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	length, err := buf.GetUint32(offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get bytes: %w", err)
	}

	bytes := make([]byte, length)
	readLen, err := buf.Read(bytes)

	if !errors.Is(err, io.EOF) && err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if errors.Is(err, io.EOF) {
		return bytes, fmt.Errorf("%w", err)
	}

	if uint32(readLen) != length {
		return nil, ErrGetFromBuffer
	}

	return bytes, nil
}

// SetBytes writes bytes to page.
// ---------------------------------------
// | bytes length (uint32) | body (bytes)|
// ---------------------------------------.
func (buf *Buffer) SetBytes(offset int64, p []byte) error {
	if err := buf.SetUint32(offset, uint32(len(p))); err != nil {
		return fmt.Errorf("failed to set bytes: %w", err)
	}

	if _, err := buf.Write(p); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
