package bytes

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/goropikari/simpledb_go/meta"
)

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

var endianness = binary.BigEndian

var (
	// ErrInvalidOffset is an error that means given offset is invalid.
	ErrInvalidOffset = errors.New("invalid offset")

	// ErrNotEnoughSpace is an error that means there is no enough space.
	ErrNotEnoughSpace = errors.New("not enough space")

	// ErrOutOfRange is an error type that refer out of range.
	ErrOutOfRange = errors.New("reference out of range of buffer")

	// ErrUnsupportedWhence is an error type that given whence is unsupported.
	ErrUnsupportedWhence = errors.New("unsupported whence")
)

// Buffer is a buffer.
type Buffer struct {
	capacity int64
	data     []byte
	offset   int64
}

// NewBuffer is a constructor of Buffer.
func NewBuffer(n int) *Buffer {
	return NewBufferBytes(make([]byte, n))
}

// NewBufferBytes is a constructor of Buffer by byte slice.
func NewBufferBytes(data []byte) *Buffer {
	return &Buffer{
		capacity: int64(len(data)),
		data:     data,
		offset:   0,
	}
}

// Read reads bytes from Reader.
func (buf *Buffer) Read(p []byte) (n int, err error) {
	if buf.offset >= buf.capacity {
		return 0, io.EOF
	}

	cnt := copy(p, buf.data[buf.offset:])

	buf.offset += int64(cnt)

	return cnt, nil
}

// Write writes given bytes to writer.
func (buf *Buffer) Write(p []byte) (n int, err error) {
	if buf.offset+int64(len(p)) > buf.capacity {
		return 0, ErrNotEnoughSpace
	}

	cnt := copy(buf.data[buf.offset:], p)
	buf.offset += int64(cnt)

	return cnt, nil
}

// Seek seeks position.
func (buf *Buffer) Seek(offset int64, whence int) (int64, error) {
	if whence != io.SeekStart {
		return 0, ErrUnsupportedWhence
	}

	if offset < 0 || offset >= buf.capacity {
		return 0, ErrOutOfRange
	}

	buf.offset = offset

	return offset, nil
}

// GetData returns buffer.
func (buf *Buffer) GetData() []byte {
	return buf.data
}

// Reset resets buffer.
func (buf *Buffer) Reset() {
	buf.offset = 0
	for i := 0; i < int(buf.capacity); i++ {
		buf.data[i] = 0
	}
}

// GetInt32 returns int32 from buffer.
func (buf *Buffer) GetInt32(offset int64) (int32, error) {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return 0, fmt.Errorf("failed to GetInt32: %w", err)
	}

	if !buf.hasSpace(meta.Int32Length) {
		return 0, ErrInvalidOffset
	}

	var ret int32
	if err := binary.Read(buf, endianness, &ret); err != nil {
		return 0, fmt.Errorf("failed to GetInt32: %w", err)
	}

	return ret, nil
}

// SetInt32 returns int32 from buffer.
// --------------------
// |  int32 (4 bytes) |
// --------------------.
func (buf *Buffer) SetInt32(offset int64, x int32) error {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("failed to SetInt32: %w", err)
	}

	if !buf.hasSpace(meta.Int32Length) {
		return ErrInvalidOffset
	}

	if err := binary.Write(buf, endianness, x); err != nil {
		return fmt.Errorf("failed to SetInt32: %w", err)
	}

	return nil
}

// GetUint32 returns uint32 from buffer.
func (buf *Buffer) GetUint32(offset int64) (uint32, error) {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return 0, fmt.Errorf("failed to GetUint32: %w", err)
	}

	if !buf.hasSpace(meta.Uint32Length) {
		return 0, ErrInvalidOffset
	}

	var ret uint32
	err := binary.Read(buf, endianness, &ret)
	if err != nil {
		return 0, fmt.Errorf("failed to GetUint32: %w", err)
	}

	return ret, nil
}

// SetUint32 returns uint32 from buffer.
// --------------------
// | uint32 (4 bytes) |
// --------------------.
func (buf *Buffer) SetUint32(offset int64, x uint32) error {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("failed to SetUint32: %w", err)
	}

	if !buf.hasSpace(meta.Uint32Length) {
		return ErrInvalidOffset
	}

	err := binary.Write(buf, endianness, x)
	if err != nil {
		return fmt.Errorf("failed to SetUint32: %w", err)
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
	if err != nil {
		return "", fmt.Errorf("failed to get string: %w", err)
	}

	return string(bytes), nil
}

// SetString returns string from buffer.
// ----------------------------------------
// | string length (uint32)| body (string)|
// ----------------------------------------.
func (buf *Buffer) SetString(offset int64, str string) error {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("failed to set string: %w", err)
	}

	if !buf.hasSpace(meta.Uint32Length + len(str)) {
		return ErrInvalidOffset
	}

	if err := buf.SetUint32(offset, uint32(len(str))); err != nil {
		return fmt.Errorf("failed to set string: %w", err)
	}

	if _, err := buf.Write([]byte(str)); err != nil {
		return fmt.Errorf("failed to set string: %w", err)
	}

	return nil
}

// GetBytes returns bytes from page.
func (buf *Buffer) GetBytes(offset int64) ([]byte, error) {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to get bytes: %w", err)
	}

	length, err := buf.GetUint32(offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get bytes: %w", err)
	}
	if length == 0 {
		return []byte{}, nil
	}

	if !buf.hasSpace(int(length)) {
		return nil, ErrInvalidOffset
	}

	bytes := make([]byte, length)

	_, err = buf.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get bytes: %w", err)
	}

	return bytes, nil
}

// SetBytes writes bytes to page.
// ---------------------------------------
// | bytes length (uint32) | body (bytes)|
// ---------------------------------------.
func (buf *Buffer) SetBytes(offset int64, p []byte) error {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("failed to set bytes: %w", err)
	}

	if !buf.hasSpace(meta.Uint32Length + len(p)) {
		return ErrInvalidOffset
	}

	if err := buf.SetUint32(offset, uint32(len(p))); err != nil {
		return fmt.Errorf("failed to set bytes: %w", err)
	}

	if _, err := buf.Write(p); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("failed to set bytes: %w", err)
	}

	return nil
}

func (buf *Buffer) hasSpace(x int) bool {
	return int64(x)+buf.offset <= buf.capacity
}
