package file

import (
	"fmt"

	"github.com/goropikari/simpledb_go/core"
)

// Block is a model representing `filename`'s  `blockNumber` th block.
type Block struct {
	fileName    core.FileName
	blockNumber core.BlockNumber
}

// NewBlock is a constructor of Block.
func NewBlock(fileName core.FileName, blockNumber core.BlockNumber) *Block {
	return &Block{
		fileName:    fileName,
		blockNumber: blockNumber,
	}
}

// GetFileName is a getter of fileName.
func (b *Block) GetFileName() core.FileName {
	if b == nil {
		return ""
	}

	return b.fileName
}

// GetBlockNumber is a getter of blockNumber
func (b *Block) GetBlockNumber() core.BlockNumber {
	if b == nil {
		return 0
	}

	return b.blockNumber
}

// Equal compares equivalence of receiver and other.
func (b *Block) Equal(other *Block) bool {
	if b == nil {
		return false
	}

	return b.fileName == other.fileName && b.blockNumber == other.blockNumber
}

// String stringfy the receiver.
func (b *Block) String() string {
	if b == nil {
		return "<nil>"
	}

	return fmt.Sprintf("[file %v, block %v]", b.fileName, b.blockNumber)
}

// HashCode returns receiver's hash value.
func (b *Block) HashCode() int {
	if b == nil {
		return 0
	}

	return core.HashCode(b.String())
}
