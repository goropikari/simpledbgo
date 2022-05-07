package domain

// Predicate is node of predicate.
type Predicate struct {
	terms []Term
}

// NewPredicate constructs a predicate.
func NewPredicate(terms []Term) *Predicate {
	return &Predicate{
		terms: terms,
	}
}
