package buffer

import "github.com/goropikari/simpledb_go/backend/core"

func (buf *Buffer) GetPage() *core.Page {
	return buf.page
}

func (buf *Buffer) Pin() {
	buf.pin()
}

func (buf *Buffer) Unpin() {
	buf.unpin()
}

func (buf *Buffer) GetPins() int {
	return buf.pins
}

func (buf *Buffer) GetTxNum() int32 {
	return buf.txnum
}

func (buf *Buffer) GetLSN() int32 {
	return buf.lsn
}

func (buf *Buffer) SetModified(txnum, lsn int32) {
	buf.setModified(txnum, lsn)
}

func (buf *Buffer) AssignToBlock(block *core.Block) error {
	return buf.assignToBlock(block)
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
