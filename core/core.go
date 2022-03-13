package core

import "errors"

var (
	// InvalidFileNameFormatError is an error that means file name is invalid format.
	InvalidFileNameFormatError = errors.New("invalid filename format")

	// NonPositiveBlockNumberError is an error that means BlockNumber is non positive.
	NonPositiveBlockNumberError = errors.New("block number must be positive")

	// NilReceiverError is an error that means receiver is nil.
	NilReceiverError = errors.New("receiver is nil")
)

// FileName is a type for filename.
type FileName string

// NewFileName is a constructor of FileName.
func NewFileName(name string) (FileName, error) {
	if name == "" {
		return "", InvalidFileNameFormatError
	}

	return FileName(name), nil
}

// BlockNumber is a type for block number.
type BlockNumber int

// NewBlockNumber is a constructor of BlockNumber.
func NewBlockNumber(bn int) (BlockNumber, error) {
	if bn <= 0 {
		return 0, NonPositiveBlockNumberError
	}

	return BlockNumber(bn), nil
}

// HashCode calculates given string hash value.
func HashCode(s string) int {
	result := 1
	base := 139
	mod := 1000000009
	b := 1

	for c := range []byte(s) {
		result *= c * b
		b = (b * base) % mod
		result %= mod
	}

	return result
}
