package directio

import (
	"fmt"
	"os"

	"github.com/ncw/directio"
)

// ErrInvalidBlockSize is error type indicating given block does not satisfy direct IO constraint.
var ErrInvalidBlockSize = fmt.Errorf("block size must be multiple of %d", BlockSize)

// BlockSize is block size for direct IO.
const BlockSize = directio.BlockSize

// AlignedBlock returns byte slice satisfying direct IO.
func AlignedBlock(n int) ([]byte, error) {
	if n%BlockSize != 0 {
		return nil, ErrInvalidBlockSize
	}

	return directio.AlignedBlock(n), nil
}

// IsAligned check whether given byte slice satisfy direct IO constraint or not.
func IsAligned(p []byte) bool {
	return directio.IsAligned(p)
}

// OpenFile opens file with direct IO option.
func OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	file, err := directio.OpenFile(name, flag, perm)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return file, nil
}
