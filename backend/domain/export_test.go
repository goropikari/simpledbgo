package domain

import "os"

func (f *File) Remove() {
	os.Remove(string(f.Name()))
}

func (mgr *LogManager) CurrentBlock() *Block {
	return mgr.currentBlock
}

func (mgr *LogManager) SetLatestLSN(x int32) {
	mgr.latestLSN = x
}

func (mgr *LogManager) SetLastSavedLSN(x int32) {
	mgr.lastSavedLSN = x
}
