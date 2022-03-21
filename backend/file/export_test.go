package file

import "github.com/goropikari/simpledb_go/backend/core"

func (mgr *Manager) CloseFile(filename core.FileName) error {
	return mgr.closeFile(filename)
}
