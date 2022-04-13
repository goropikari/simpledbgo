package tx

import (
	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/lib/list"
)

func (list *BufferList) PinnedBlocks() list.List[*domain.Block] {
	return list.pinnedBlocks
}

func (tx *Transaction) Number() domain.TransactionNumber {
	return tx.number
}
