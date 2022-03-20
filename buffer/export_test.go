package buffer

import "github.com/goropikari/simpledb_go/file"

type Buffer = buffer

func (buf *buffer) GetPage() *file.Page {
	if buf == nil {
		return nil
	}

	return buf.page
}

func (buf *buffer) SetModified(txnum, lsn int) {
	buf.setModified(txnum, lsn)
}
func (mgr *Manager) Pin(block *file.Block) (*buffer, error) {
	return mgr.pin(block)
}

func (mgr *Manager) Unpin(buf *buffer) error {
	return mgr.unpin(buf)
}

func (mgr *Manager) Available() int {
	return mgr.available()
}
