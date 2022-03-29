package bytes

import (
	"errors"
	"fmt"

	"github.com/goropikari/simpledb_go/lib/directio"
)

// ErrInvalidDirectIO is an error that means invalid direct io.
var ErrInvalidDirectIO = errors.New("invalid direct io")

// NewDirectBuffer is a constructor of DirectBuffer.
func NewDirectBuffer(n int) (*Buffer, error) {
	buf, err := directio.AlignedBlock(n)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	buffer, err := NewDirectBufferBytes(buf)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

// NewDirectBufferBytes is a constructor of DirectBuffer by byte slice.
func NewDirectBufferBytes(buf []byte) (*Buffer, error) {
	if !directio.IsAligned(buf) {
		return nil, ErrInvalidDirectIO
	}

	return NewBufferBytes(buf), nil
}
