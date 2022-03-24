package buffer

import "github.com/goropikari/simpledb_go/backend/core"

func (buf *Buffer) GetPage() *core.Page {
	return buf.page
}

func (buf *Buffer) SetModified(txnum, lsn int) {
	buf.setModified(txnum, lsn)
}

func (mgr *Manager) Pin(block *core.Block) (*Buffer, error) {
	return mgr.pin(block)
}

func (mgr *Manager) Unpin(buf *Buffer) error {
	return mgr.unpin(buf)
}

func (mgr *Manager) Available() int {
	return mgr.available()
}
