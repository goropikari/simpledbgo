package domain

import (
	"errors"
)

// ErrNegativeBlockNumber means given block number is non negative.
var ErrNegativeBlockNumber = errors.New("block number must be non negative")

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
		return 0, errors.New("block size must be positive")
	}

	return BlockSize(n), nil
}

// Block is a model of block.
type Block struct {
	filename FileName
	size     BlockSize
	number   BlockNumber
	offset   int64
}

// NewBlock is a constructor of Block.
func NewBlock(filename FileName, size BlockSize, number BlockNumber) *Block {
	offset := int64(size) * int64(number)

	return &Block{
		filename: filename,
		size:     size,
		number:   number,
		offset:   offset,
	}
}

// Equal compares equality of two blocks.
func (b *Block) Equal(other *Block) bool {
	return b.filename == other.filename && b.number == other.number
}

// FileName returns corresponding file name.
func (b *Block) FileName() FileName {
	return b.filename
}

// Size returns block size.
func (b *Block) Size() BlockSize {
	return b.size
}

// Offset returns file's corresponding offset.
func (b *Block) Offset() int64 {
	return b.offset
}
