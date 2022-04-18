package log

import "github.com/goropikari/simpledbgo/backend/domain"

func (mgr *Manager) CurrentBlock() domain.Block {
	return mgr.currentBlock
}

func (mgr *Manager) SetLatestLSN(x domain.LSN) {
	mgr.latestLSN = x
}

func (mgr *Manager) SetLastSavedLSN(x domain.LSN) {
	mgr.lastSavedLSN = x
}
