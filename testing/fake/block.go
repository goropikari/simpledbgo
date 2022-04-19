package fake

import "github.com/goropikari/simpledbgo/backend/domain"

func FileName() domain.FileName {
	return domain.FileName(RandString())
}

func BlockNumber() domain.BlockNumber {
	return domain.BlockNumber(RandInt32())
}

func BlockSize() domain.BlockSize {
	return domain.BlockSize(RandInt32())
}

func Block() domain.Block {
	return domain.NewBlock(FileName(), BlockNumber())
}
