package domain

type Predicate struct {
	terms []Term
}

func NewPredicate() Predicate {
	return Predicate{}
}

func (pred *Predicate) Add(term Term) {
	pred.terms = append(pred.terms, term)
}

func (pred *Predicate) ConjoinWith(p Predicate) {
	pred.terms = append(pred.terms, p.terms...)
}
