package domain

type Term struct {
	lhs Expression
	rhs Expression
}

func NewTerm(lhs, rhs Expression) Term {
	return Term{lhs: lhs, rhs: rhs}
}
