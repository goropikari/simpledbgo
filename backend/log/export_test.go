package log

import "github.com/goropikari/simpledb_go/backend/domain"

func (mgr *Manager) CurrentBlock() *domain.Block {
	return mgr.currentBlock
}

func (mgr *Manager) SetLatestLSN(x int32) {
	mgr.latestLSN = x
}

func (mgr *Manager) SetLastSavedLSN(x int32) {
	mgr.lastSavedLSN = x
}
