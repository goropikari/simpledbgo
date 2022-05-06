package domain

// Predicate is node of predicate.
type Predicate struct {
	terms []Term
}

// NewPredicate constructs a predicate.
func NewPredicate() Predicate {
	return Predicate{}
}

// Add adds a term.
func (pred *Predicate) Add(term Term) {
	pred.terms = append(pred.terms, term)
}
