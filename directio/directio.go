package directio

import (
	"fmt"
	"os"

	"github.com/ncw/directio"
)

// InvalidBlockSizeError is error type indicating given block does not satisfy direct IO constraint.
var InvalidBlockSizeError = fmt.Errorf("block size must be multiple of %d", BlockSize)

// BlockSize is block size for direct IO.
var BlockSize = directio.BlockSize

// AlignedBlock returns byte slice satisfying direct IO.
func AlignedBlock(n int) ([]byte, error) {
	if n%BlockSize != 0 {
		return nil, InvalidBlockSizeError
	}

	return directio.AlignedBlock(n), nil
}

// IsAligned check whether given byte slice satisfy direct IO constraint or not.
func IsAligned(p []byte) bool {
	return directio.IsAligned(p)
}

// OpenFile opens file with direct IO option.
func OpenFile(name string, flag int, perm os.FileMode) (file *os.File, err error) {
	return directio.OpenFile(name, flag, perm)
}
