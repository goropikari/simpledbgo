package bytes

import (
	"errors"

	"github.com/goropikari/simpledb_go/directio"
)

// NewDirectBuffer is a constructor of DirectBuffer.
func NewDirectBuffer(n int) (*Buffer, error) {
	buf, err := directio.AlignedBlock(n)
	if err != nil {
		return nil, err
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
		return nil, errors.New("buffer must satisfy O_DIRECT constraints")
	}

	return NewBufferBytes(buf), nil
}
