package bytes

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// ErrInvalidOffset is an error that means given offset is invalid.
var ErrInvalidOffset = errors.New("invalid offset")

var endianness = binary.BigEndian

// ByteBuffer is an interface that implement io.ReadWriteSeeker.
type ByteBuffer interface {
	io.ReadWriteSeeker
	GetBufferBytes() []byte
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

const (
	int32Length  = 4
	uint32Length = 4
)

// Buffer is a buffer.
type Buffer struct {
	capacity int64
	buf      []byte
	off      int64
}

// NewBuffer is a constructor of Buffer.
func NewBuffer(n int) *Buffer {
	return NewBufferBytes(make([]byte, n))
}

// NewBufferBytes is a constructor of Buffer by byte slice.
func NewBufferBytes(buf []byte) *Buffer {
	return &Buffer{
		capacity: int64(len(buf)),
		buf:      buf,
		off:      0,
	}
}

// Read reads bytes from Reader.
func (buf *Buffer) Read(p []byte) (n int, err error) {
	if buf.off < 0 || buf.off >= buf.capacity {
		return 0, io.EOF
	}

	cnt := copy(p, buf.buf[buf.off:])

	buf.off += int64(cnt)
	if buf.off == buf.capacity {
		return cnt, io.EOF
	}

	return cnt, nil
}

// Write writes given bytes to writer.
func (buf *Buffer) Write(p []byte) (n int, err error) {
	if buf.off < 0 || buf.off >= buf.capacity {
		return 0, io.EOF
	}

	cnt := copy(buf.buf[buf.off:], p)

	buf.off += int64(cnt)
	// if buf.off == buf.capacity {
	// 	return cnt, io.EOF
	// }

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

	if off < 0 || off > buf.capacity {
		return 0, ErrOutOfRange
	}

	buf.off = off

	return off, nil
}

// GetBufferBytes returns buffer.
func (buf *Buffer) GetBufferBytes() []byte {
	return buf.buf
}

// Reset resets buffer.
func (buf *Buffer) Reset() {
	buf.off = 0
	for i := 0; i < int(buf.capacity); i++ {
		buf.buf[i] = 0
	}
}

// GetInt32 returns int32 from buffer.
func (buf *Buffer) GetInt32(offset int64) (int32, error) {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	if !buf.hasSpace(int32Length) {
		return 0, ErrInvalidOffset
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

	if !buf.hasSpace(int32Length) {
		return ErrInvalidOffset
	}

	if err := binary.Write(buf, endianness, x); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// GetUint32 returns uint32 from buffer.
func (buf *Buffer) GetUint32(offset int64) (uint32, error) {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	if !buf.hasSpace(uint32Length) {
		return 0, ErrInvalidOffset
	}

	var ret uint32
	err := binary.Read(buf, endianness, &ret)
	if err != nil && !errors.Is(err, io.EOF) {
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

	if !buf.hasSpace(uint32Length) {
		return ErrInvalidOffset
	}

	err := binary.Write(buf, endianness, x)
	if err != nil && !errors.Is(err, io.EOF) {
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

	if !buf.hasSpace(int(length)) {
		return "", ErrInvalidOffset
	}

	bytes := make([]byte, length)

	_, err = buf.Read(bytes)
	if errors.Is(err, io.EOF) {
		return string(bytes), nil
	} else if err != nil {
		return "", err
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

	if !buf.hasSpace(uint32Length + len(str)) {
		return ErrInvalidOffset
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

	if !buf.hasSpace(int(length)) {
		return nil, ErrInvalidOffset
	}

	bytes := make([]byte, length)

	_, err = buf.Read(bytes)
	if errors.Is(err, io.EOF) {
		return bytes, nil
	} else if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return bytes, nil
}

// SetBytes writes bytes to page.
// ---------------------------------------
// | bytes length (uint32) | body (bytes)|
// ---------------------------------------.
func (buf *Buffer) SetBytes(offset int64, p []byte) error {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	if !buf.hasSpace(uint32Length + len(p)) {
		return ErrInvalidOffset
	}

	if err := buf.SetUint32(offset, uint32(len(p))); err != nil {
		return fmt.Errorf("failed to set bytes: %w", err)
	}

	if _, err := buf.Write(p); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (buf *Buffer) hasSpace(x int) bool {
	return buf.capacity-buf.off >= int64(x)
}
