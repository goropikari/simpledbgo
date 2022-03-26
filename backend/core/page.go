package core

import (
	"errors"
	"fmt"
	"io"

	"github.com/goropikari/simpledb_go/lib/bytes"
)

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
	readLen, err := page.bb.GetInt32(offset)
	if errors.Is(err, io.EOF) {
		return readLen, fmt.Errorf("%w", err)
	} else if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	return readLen, nil
}

// SetInt32 returns int32 from buffer.
// --------------------
// |  int32 (4 bytes) |
// --------------------.
func (page *Page) SetInt32(offset int64, x int32) error {
	if err := page.bb.SetInt32(offset, x); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// GetUint32 returns uint32 from buffer.
func (page *Page) GetUint32(offset int64) (uint32, error) {
	n, err := page.bb.GetUint32(offset)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	return n, nil
}

// SetUint32 returns uint32 from buffer.
// --------------------
// | uint32 (4 bytes) |
// --------------------.
func (page *Page) SetUint32(offset int64, x uint32) error {
	if err := page.bb.SetUint32(offset, x); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// GetBytes returns bytes from page.
func (page *Page) GetBytes(offset int64) ([]byte, error) {
	bytes, err := page.bb.GetBytes(offset)
	if errors.Is(err, io.EOF) {
		return bytes, fmt.Errorf("%w", err)
	} else if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return bytes, nil
}

// SetBytes writes bytes to page.
// ---------------------------------------
// | bytes length (uint32) | body (bytes)|
// ---------------------------------------.
func (page *Page) SetBytes(offset int64, p []byte) error {
	if err := page.bb.SetBytes(offset, p); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// GetString returns string from buffer.
func (page *Page) GetString(offset int64) (string, error) {
	s, err := page.bb.GetString(offset)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return s, nil
}

// SetString returns string from buffer.
// ----------------------------------------
// | string length (uint32)| body (string)|
// ----------------------------------------.
func (page *Page) SetString(offset int64, s string) error {
	if err := page.bb.SetString(offset, s); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// Write writes bytes to page.
func (page *Page) Write(p []byte) (int, error) {
	n, err := page.bb.Write(p)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	return n, nil
}

// GetBufferBytes returns page buffer.
func (page *Page) GetBufferBytes() []byte {
	return page.bb.GetBufferBytes()
}

func (page *Page) Reset() {
	page.bb.Reset()
}

func (page *Page) Seek(offset int64, whence int) (int64, error) {
	offset, err := page.bb.Seek(offset, whence)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	return offset, nil
}
