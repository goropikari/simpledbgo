package fake

import (
	"github.com/goropikari/simpledbgo/domain"
)

// Buffer is a fake domain.Buffer
func Buffer() *domain.Buffer {
	buf := &domain.Buffer{}
	buf.SetModifiedTxNumber(domain.TransactionNumber(RandInt32()), domain.LSN(RandInt32()))

	return buf
}
