package directio

import (
	"fmt"
	"os"

	"github.com/ncw/directio"
)

var InvalidBlockSize = fmt.Errorf("block size must be multiple of %d", BlockSize)

var BlockSize = directio.BlockSize

func AlignedBlock(n int) ([]byte, error) {
	if n%BlockSize != 0 {
		return nil, InvalidBlockSize
	}

	return directio.AlignedBlock(n), nil
}

func IsAligned(p []byte) bool {
	return directio.IsAligned(p)
}

func OpenFile(name string, flag int, perm os.FileMode) (file *os.File, err error) {
	return directio.OpenFile(name, flag, perm)
}
