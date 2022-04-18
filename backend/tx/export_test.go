package tx

import (
	"github.com/goropikari/simpledbgo/backend/domain"
	"github.com/goropikari/simpledbgo/lib/list"
)

func (list *BufferList) PinnedBlocks() list.List[domain.Block] {
	return list.pinnedBlocks
}

func (tx *Transaction) Number() domain.TransactionNumber {
	return tx.number
}
