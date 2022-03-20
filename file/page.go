package file

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/goropikari/simpledb_go/bytes"
	"github.com/goropikari/simpledb_go/core"
)

// InvalidOffsetError is error type that indicates you specify invalid offset.
var InvalidOffsetError = errors.New("you may be specify invalid offset")

// Page is a model of a page.
type Page struct {
	bb bytes.ByteBuffer
}

// NewPage is a constructor of Page.
func NewPage(bb bytes.ByteBuffer) *Page {
	return &Page{
		bb: bb,
	}
}

// GetInt32 returns int32 from buffer.
func (page *Page) GetInt32(offset int64) (int32, error) {
	if page == nil {
		return 0, nil
	}

	if _, err := page.bb.Seek(offset, io.SeekStart); err != nil {
		return 0, err
	}

	var ret int32
	if err := binary.Read(page.bb, core.Endianness, &ret); err != nil {
		if err == io.EOF {
			return ret, err
		}
		return 0, err
	}

	return ret, nil
}

// SetInt32 returns int32 from buffer.
// --------------------
// |  int32 (4 bytes) |
// --------------------
func (page *Page) SetInt32(offset int64, x int32) error {
	if page == nil {
		return nil
	}

	if _, err := page.bb.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	if err := binary.Write(page.bb, core.Endianness, x); err != nil {
		return err
	}

	return nil
}

// GetUInt32 returns uint32 from buffer.
func (page *Page) GetUInt32(offset int64) (uint32, error) {
	if page == nil {
		return 0, nil
	}

	if _, err := page.bb.Seek(offset, io.SeekStart); err != nil {
		return 0, err
	}

	var ret uint32
	if err := binary.Read(page.bb, core.Endianness, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

// SetUInt32 returns uint32 from buffer.
// --------------------
// | uint32 (4 bytes) |
// --------------------
func (page *Page) SetUInt32(offset int64, x uint32) error {
	if page == nil {
		return nil
	}

	if _, err := page.bb.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	if err := binary.Write(page.bb, core.Endianness, x); err != nil {
		return err
	}

	return nil
}

// GetBytes returns bytes from page.
func (page *Page) GetBytes(offset int64) ([]byte, error) {
	if page == nil {
		return nil, core.NilReceiverError
	}

	if _, err := page.bb.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}

	length, err := page.GetUInt32(offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get bytes: %w", err)
	}

	buf := make([]byte, length)
	_, err = page.bb.Read(buf)
	if err == io.EOF {
		return buf, err
	} else if err != nil {
		return nil, err
	}

	return buf, nil
}

// SetBytes writes bytes to page.
// ---------------------------------------
// | bytes length (uint32) | body (bytes)|
// ---------------------------------------
func (page *Page) SetBytes(offset int64, p []byte) error {
	if page == nil {
		return core.NilReceiverError
	}
	if err := page.SetUInt32(offset, uint32(len(p))); err != nil {
		return fmt.Errorf("failed to set bytes: %w", err)
	}

	if _, err := page.bb.Write(p); err != nil {
		return err
	}

	return nil
}

// GetString returns string from buffer.
func (page *Page) GetString(offset int64) (string, error) {
	if page == nil {
		return "", nil
	}

	if _, err := page.bb.Seek(offset, io.SeekStart); err != nil {
		return "", fmt.Errorf("failed to get string: %w", err)
	}

	length, err := page.GetUInt32(offset)
	if err != nil {
		return "", err
	}

	b := make([]byte, length)
	n, err := page.bb.Read(b)
	if err != nil && err == io.EOF {
		return "", err
	}
	if uint32(n) != length {
		return "", InvalidOffsetError
	}

	return string(b), err
}

// SetString returns string from buffer.
// ----------------------------------------
// | string length (uint32)| body (string)|
// ----------------------------------------
func (page *Page) SetString(offset int64, s string) error {
	if _, err := page.bb.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	if err := page.SetUInt32(offset, uint32(len(s))); err != nil {
		return fmt.Errorf("failed to set string: %w", err)
	}
	if _, err := page.bb.Write([]byte(s)); err != nil {
		return err
	}

	return nil
}

// Write writes bytes to page.
func (page *Page) Write(p []byte) (int, error) {
	return page.bb.Write(p)
}

// GetFullBytes returns page buffer.
func (page *Page) GetFullBytes() []byte {
	return page.bb.GetBytes()
}

func (page *Page) Reset() {
	page.bb.Reset()
}
