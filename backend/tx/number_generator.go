package tx

import (
	"sync"

	"github.com/goropikari/simpledbgo/backend/domain"
)

// NumberGenerator is generator of tx number.
type NumberGenerator struct {
	mu    *sync.Mutex
	txnum domain.TransactionNumber
}

// NewNumberGenerator constructs a NumberGenerator.
func NewNumberGenerator() *NumberGenerator {
	return &NumberGenerator{
		mu:    &sync.Mutex{},
		txnum: 0,
	}
}

// Generate generates new tx number.
func (gen *NumberGenerator) Generate() domain.TransactionNumber {
	gen.mu.Lock()
	defer gen.mu.Unlock()

	gen.txnum++

	return gen.txnum
}
