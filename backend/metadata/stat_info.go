package metadata

import "github.com/goropikari/simpledbgo/domain"

// StatInfo is a model of statistical information.
type StatInfo struct {
	numBlocks int
	numRecs   int
}

// NewStatInfo construts a StatInfo.
func NewStatInfo(numBlocks int, numRecs int) StatInfo {
	return StatInfo{
		numBlocks: numBlocks,
		numRecs:   numRecs,
	}
}

// EstNumBlocks returns estimated the number of blocks.
func (info StatInfo) EstNumBlocks() int {
	return info.numBlocks
}

// EstNumRecord returns estimated the number of records.
func (info StatInfo) EstNumRecord() int {
	return info.numRecs
}

// EstDistinctVals returns estimated the number of distinct values.
func (info StatInfo) EstDistinctVals(fldname domain.FieldName) int {
	roughEst := 3

	return 1 + (info.numRecs / roughEst)
}
