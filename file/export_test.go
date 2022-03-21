package file

import "github.com/goropikari/simpledb_go/core"

func (mgr *Manager) CloseFile(filename core.FileName) error {
	return mgr.closeFile(filename)
}
