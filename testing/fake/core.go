package fake

import "github.com/goropikari/simpledb_go/backend/core"

// Block returns fake block.
func Block() *core.Block {
	filename, err := core.NewFileName(RandString())
	if err != nil {
		panic(err)
	}

	blkNum, err := core.NewBlockNumber(RandInt())
	if err != nil {
		panic(err)
	}

	return core.NewBlock(filename, blkNum)
}
