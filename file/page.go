package file

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/goropikari/simpledb_go/bytes"
	"github.com/goropikari/simpledb_go/core"
)

// InvalidOffsetError is error type that indicates you specify invalid offset.
var InvalidOffsetError = errors.New("you may be specify invalid offset")

// Int32Length is byte length of int32
var Int32Length = 4

var order = binary.BigEndian

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
	if err := binary.Read(page.bb, order, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

// SetInt32 returns int32 from buffer.
func (page *Page) SetInt32(offset int64, x int32) error {
	if page == nil {
		return nil
	}

	if _, err := page.bb.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	if err := binary.Write(page.bb, order, x); err != nil {
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
	if err := binary.Read(page.bb, order, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

// SetUInt32 returns uint32 from buffer.
func (page *Page) SetUInt32(offset int64, x uint32) error {
	if page == nil {
		return nil
	}

	if _, err := page.bb.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	if err := binary.Write(page.bb, order, x); err != nil {
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
		return nil, err
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
func (page *Page) SetBytes(offset int64, p []byte) error {
	if err := page.SetUInt32(offset, uint32(len(p))); err != nil {
		return err
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
		return "", err
	}

	length, err := page.GetUInt32(offset)
	if err != nil {
		return "", err
	}

	b := make([]byte, length)
	if n, err := page.bb.Read(b); err != nil {
		return "", err
	} else if uint32(n) != length {
		return "", InvalidOffsetError
	}

	return string(b), nil
}

// SetString returns string from buffer.
func (page *Page) SetString(offset int64, s string) error {
	if _, err := page.bb.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	if err := page.SetUInt32(offset, uint32(len(s))); err != nil {
		return err
	}
	if _, err := page.bb.Write([]byte(s)); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// Contents returns ByteBuffer.
func (page *Page) Contents() (bytes.ByteBuffer, error) {
	if _, err := page.bb.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	return page.bb, nil
}

// Write writes bytes to page.
func (page *Page) Write(p []byte) (int, error) {
	return page.bb.Write(p)
}

// GetFullBytes returns page buffer.
func (page *Page) GetFullBytes() []byte {
	return page.bb.GetBytes()
}
