package core

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
)

var (
	// ErrInvalidFileNameFormat is an error that means file name is invalid format.
	ErrInvalidFileNameFormat = errors.New("invalid filename format")

	// NonPositiveBlockNumberError is an error that means BlockNumber is non positive.
	ErrNonNegativeBlockNumber = errors.New("block number must be non negative")

	// ErrBlockSize is an error that means BlockSize in non positive.
	ErrBlockSize = errors.New("block size must be positive")
)

// Uint32Length is byte length of int32.
const Uint32Length = 4

// FileName is a type for filename.
type FileName string

// NewFileName is a constructor of FileName.
func NewFileName(name string) (FileName, error) {
	if name == "" {
		return "", ErrInvalidFileNameFormat
	}

	return FileName(name), nil
}

// BlockNumber is a type for block number.
type BlockNumber int

// NewBlockNumber is a constructor of BlockNumber.
func NewBlockNumber(bn int) (BlockNumber, error) {
	if bn < 0 {
		return 0, ErrNonNegativeBlockNumber
	}

	return BlockNumber(bn), nil
}

// BlockSize is a type for block size.
type BlockSize int

// NewBlockSize is a constructor of BlockSize.
func NewBlockSize(x int) (BlockSize, error) {
	if x <= 0 {
		return 0, ErrBlockSize
	}

	return BlockSize(x), nil
}

// HashCode calculates given string hash value.
func HashCode(str string) int {
	result := 1
	base := 139
	mod := 1000000009
	b := 1

	for c := range []byte(str) {
		result *= c * b
		b = (b * base) % mod
		result %= mod
	}

	return result
}

// RandomString returns random string.
func RandomString() string {
	return fmt.Sprintf("%v", rand.Uint32())
}

// FileSize returns file size.
func FileSize(f *os.File) (int64, error) {
	info, err := f.Stat()
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	return info.Size(), nil
}
