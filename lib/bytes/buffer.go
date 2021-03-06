package bytes

import (
	"encoding/binary"
	"io"

	"github.com/goropikari/simpledbgo/errors"

	"github.com/goropikari/simpledbgo/common"
)

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

var endianness = binary.BigEndian

const (
	byteFieldByteLength   = common.Uint32Length
	stringFieldByteLength = common.Uint32Length
)

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
	pos      int64
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
		pos:      0,
	}
}

// Size returns a buffer size.
func (buf *Buffer) Size() int64 {
	return buf.capacity
}

// Read reads bytes from Reader.
func (buf *Buffer) Read(p []byte) (n int, err error) {
	if buf.pos >= buf.capacity {
		return 0, io.EOF
	}

	cnt := copy(p, buf.data[buf.pos:])

	buf.pos += int64(cnt)

	return cnt, nil
}

// Write writes given bytes to writer.
func (buf *Buffer) Write(p []byte) (n int, err error) {
	if buf.pos+int64(len(p)) > buf.capacity {
		return 0, ErrNotEnoughSpace
	}

	cnt := copy(buf.data[buf.pos:], p)
	buf.pos += int64(cnt)

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

	buf.pos = offset

	return offset, nil
}

// GetData returns buffer.
func (buf *Buffer) GetData() []byte {
	return buf.data
}

// Reset resets buffer.
func (buf *Buffer) Reset() {
	buf.pos = 0
	for i := 0; i < int(buf.capacity); i++ {
		buf.data[i] = 0
	}
}

// GetInt32 returns int32 from buffer.
func (buf *Buffer) GetInt32(offset int64) (int32, error) {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return 0, errors.Err(err, "GetInt32")
	}

	var ret int32
	if !buf.hasSpace(ret) {
		return 0, ErrInvalidOffset
	}
	if err := binary.Read(buf, endianness, &ret); err != nil {
		return 0, errors.Err(err, "GetInt32")
	}

	return ret, nil
}

// SetInt32 returns int32 from buffer.
// --------------------
// |  int32 (4 bytes) |
// --------------------.
func (buf *Buffer) SetInt32(offset int64, x int32) error {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return errors.Err(err, "Seek")
	}

	if !buf.hasSpace(x) {
		return ErrInvalidOffset
	}

	if err := binary.Write(buf, endianness, x); err != nil {
		return errors.Err(err, "Write")
	}

	return nil
}

// getUint32 returns uint32 from buffer.
func (buf *Buffer) getUint32(offset int64) (uint32, error) {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return 0, errors.Err(err, "Seek")
	}

	var ret uint32
	if !buf.hasSpace(ret) {
		return 0, ErrInvalidOffset
	}
	if err := binary.Read(buf, endianness, &ret); err != nil {
		return 0, errors.Err(err, "Read")
	}

	return ret, nil
}

// setUint32 returns uint32 from buffer.
// --------------------
// | uint32 (4 bytes) |
// --------------------.
func (buf *Buffer) setUint32(offset int64, x uint32) error {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return errors.Err(err, "Seek")
	}

	if !buf.hasSpace(x) {
		return ErrInvalidOffset
	}

	err := binary.Write(buf, endianness, x)
	if err != nil {
		return errors.Err(err, "Write")
	}

	return nil
}

// GetString returns string from buffer.
func (buf *Buffer) GetString(offset int64) (string, error) {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return "", errors.Err(err, "Seek")
	}

	length, err := buf.getUint32(offset)
	if err != nil {
		return "", errors.Err(err, "getUint32")
	}

	bytes := make([]byte, length)

	_, err = buf.Read(bytes)
	if err != nil {
		return "", errors.Err(err, "Read")
	}

	return string(bytes), nil
}

// SetString returns string from buffer.
// ----------------------------------------
// | string length (uint32)| body (string)|
// ----------------------------------------.
func (buf *Buffer) SetString(offset int64, str string) error {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return errors.Err(err, "Seek")
	}

	if !buf.hasSpace(str) {
		return ErrInvalidOffset
	}

	if err := buf.setUint32(offset, uint32(len(str))); err != nil {
		return errors.Err(err, "setUint32")
	}

	if _, err := buf.Write([]byte(str)); err != nil {
		return errors.Err(err, "Write")
	}

	return nil
}

// GetBytes returns bytes from page.
func (buf *Buffer) GetBytes(offset int64) ([]byte, error) {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return nil, errors.Err(err, "Seek")
	}

	length, err := buf.getUint32(offset)
	if err != nil {
		return nil, errors.Err(err, "getUint32")
	}
	if length == 0 {
		return []byte{}, nil
	}

	bytes := make([]byte, length)

	_, err = buf.Read(bytes)
	if err != nil {
		return nil, errors.Err(err, "Read")
	}

	return bytes, nil
}

// SetBytes writes bytes to page.
// ---------------------------------------
// | bytes length (uint32) | body (bytes)|
// ---------------------------------------.
func (buf *Buffer) SetBytes(offset int64, p []byte) error {
	if _, err := buf.Seek(offset, io.SeekStart); err != nil {
		return errors.Err(err, "Seek")
	}

	if !buf.hasSpace(p) {
		return ErrInvalidOffset
	}

	if err := buf.setUint32(offset, uint32(len(p))); err != nil {
		return errors.Err(err, "setUint32")
	}

	if _, err := buf.Write(p); err != nil && !errors.Is(err, io.EOF) {
		return errors.Err(err, "Write")
	}

	return nil
}

func (buf *Buffer) hasSpace(x any) bool {
	return buf.NeededByteLength(x)+buf.pos <= buf.capacity
}

// NeededByteLength returns needed byte length for given type.
func (buf *Buffer) NeededByteLength(x any) int64 {
	return NeededByteLength(x)
}

// NeededByteLength returns needed byte length for given type.
func NeededByteLength(x any) int64 {
	switch v := x.(type) {
	case int32, uint32:
		return common.Int32Length
	case int64:
		return common.Int64Length
	case []byte:
		return byteFieldByteLength + int64(len(v))
	case string:
		return stringFieldByteLength + int64(len(v))
	default:
		panic(errors.New("invalid data type"))
	}
}
