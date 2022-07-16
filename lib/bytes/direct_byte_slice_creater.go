package bytes

import (
	"github.com/goropikari/simpledbgo/errors"
	"github.com/goropikari/simpledbgo/lib/directio"
)

// DirectByteSliceCreater is a slice creature with direct io.
type DirectByteSliceCreater struct{}

// NewDirectByteSliceCreater is a constructor of DirectByteSliceCreater.
func NewDirectByteSliceCreater() *DirectByteSliceCreater {
	return &DirectByteSliceCreater{}
}

// Create creates a byte slice.
func (s *DirectByteSliceCreater) Create(n int) ([]byte, error) {
	b, err := directio.AlignedBlock(n)
	if err != nil {
		return nil, errors.Err(err, "AlignedBlock")
	}

	return b, nil
}
