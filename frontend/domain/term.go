package domain

// Term is a node of term.
type Term struct {
	lhs Expression
	rhs Expression
}

// NewTerm constructs a Term.
func NewTerm(lhs, rhs Expression) Term {
	return Term{lhs: lhs, rhs: rhs}
}
