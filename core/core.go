package core

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"os"
)

var (
	// InvalidFileNameFormatError is an error that means file name is invalid format.
	InvalidFileNameFormatError = errors.New("invalid filename format")

	// NonPositiveBlockNumberError is an error that means BlockNumber is non positive.
	NonNegativeBlockNumberError = errors.New("block number must be non negative")

	// NilReceiverError is an error that means receiver is nil.
	NilReceiverError = errors.New("receiver is nil")
)

// UInt32Length is byte length of int32
var UInt32Length = 4

// Endianness is endianness of this system.
var Endianness = binary.BigEndian

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
	if bn < 0 {
		return 0, NonNegativeBlockNumberError
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

func RandomString() string {
	return fmt.Sprintf("%v", rand.Uint32())
}

func FileSize(f *os.File) (int64, error) {
	info, err := f.Stat()
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}
