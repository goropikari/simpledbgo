package domain

import (
	"errors"
)

var (
	// ErrNegativeBlockNumber means given block number is non negative.
	ErrNegativeBlockNumber = errors.New("block number must be non negative")

	// ErrNonPositiveBlockSize means given block size must be positive.
	ErrNonPositiveBlockSize = errors.New("block size must be positive")
)

// BlockNumber is value object of block number.
type BlockNumber int32

// NewBlockNumber is a constructor of BlockNumber.
func NewBlockNumber(n int32) (BlockNumber, error) {
	if n < 0 {
		return 0, ErrNegativeBlockNumber
	}

	return BlockNumber(n), nil
}

// BlockSize is value object of block size.
type BlockSize int32

// NewBlockSize is a constructor of BlockSize.
func NewBlockSize(n int32) (BlockSize, error) {
	if n <= 0 {
		return 0, ErrNonPositiveBlockSize
	}

	return BlockSize(n), nil
}

// Block is a model of block.
type Block struct {
	filename FileName
	number   BlockNumber
}

// NewBlock is a constructor of Block.
func NewBlock(filename FileName, number BlockNumber) Block {
	return Block{
		filename: filename,
		number:   number,
	}
}

// NewDummyBlock constructs a dummy Block.
func NewDummyBlock(filename FileName) Block {
	return Block{
		filename: filename,
		number:   0,
	}
}

// NewZeroBlock is constructor of zero value Block.
func NewZeroBlock() Block {
	return Block{}
}

// Equal compares equality of two blocks.
func (b Block) Equal(other Block) bool {
	return b == other
}

// FileName returns corresponding file name.
func (b Block) FileName() FileName {
	return b.filename
}

// Number returns block number.
func (b Block) Number() BlockNumber {
	return b.number
}
