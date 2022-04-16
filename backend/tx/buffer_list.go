package tx

import (
	"github.com/goropikari/simpledbgo/backend/domain"
	ls "github.com/goropikari/simpledbgo/lib/list"
)

// BufferList is list of buffer.
type BufferList struct {
	buffers      map[domain.Block]*domain.Buffer
	pinnedBlocks ls.List[*domain.Block]
	bufMgr       domain.BufferManager
}

// NewBufferList constructs a BufferList.
func NewBufferList(bufMgr domain.BufferManager) *BufferList {
	return &BufferList{
		buffers:      make(map[domain.Block]*domain.Buffer),
		pinnedBlocks: ls.NewList[*domain.Block](),
		bufMgr:       bufMgr,
	}
}

// GetBuffer gets a buffer from list.
func (list *BufferList) GetBuffer(blk domain.Block) *domain.Buffer {
	return list.buffers[blk]
}

// Pin pins a buffer.
func (list *BufferList) Pin(blk domain.Block) error {
	buf, err := list.bufMgr.Pin(&blk)
	if err != nil {
		return err
	}

	list.buffers[blk] = buf

	list.pinnedBlocks.Add(&blk)

	return nil
}

// Unpin unpins the block from buffer list.
func (list *BufferList) Unpin(blk domain.Block) {
	buf := list.buffers[blk]

	list.bufMgr.Unpin(buf)

	list.pinnedBlocks.Remove(&blk)
	if !list.pinnedBlocks.Contains(&blk) {
		delete(list.buffers, blk)
	}
}

// UnpinAll unpins all of pinned buffer on this BufferList.
func (list *BufferList) UnpinAll() {
	for _, blk := range list.pinnedBlocks.Data() {
		buf := list.buffers[*blk]
		list.bufMgr.Unpin(buf)
	}

	list.buffers = make(map[domain.Block]*domain.Buffer)
	list.pinnedBlocks = ls.NewList[*domain.Block]()
}
