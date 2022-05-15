package domain

import (
	"github.com/goropikari/simpledbgo/math"
	"github.com/goropikari/simpledbgo/meta"
)

// Term is a node of term.
type Term struct {
	lhs Expression
	rhs Expression
}

// NewTerm constructs a Term.
func NewTerm(lhs, rhs Expression) Term {
	return Term{lhs: lhs, rhs: rhs}
}

// IsSatisfied checks whether a term is satisfied or not.
func (term Term) IsSatisfied(s Scanner) bool {
	lhsVal, err := term.lhs.Evaluate(s)
	if err != nil {
		return false
	}

	rhsVal, err := term.rhs.Evaluate(s)
	if err != nil {
		return false
	}

	return lhsVal.Equal(rhsVal)
}

// ReductionFactor is reduction factor due to the predicate.
// predicate によってスキャン量がどれだけ減るかの割合.
func (term Term) ReductionFactor(p Planner) int {
	if term.lhs.IsFieldName() && term.rhs.IsFieldName() {
		lhsName := term.lhs.AsFieldName()
		rhsName := term.rhs.AsFieldName()

		return math.Max(p.EstDistinctVals(lhsName), p.EstDistinctVals(rhsName))
	}

	if term.lhs.IsFieldName() {
		return p.EstDistinctVals(term.lhs.AsFieldName())
	}

	if term.rhs.IsFieldName() {
		return p.EstDistinctVals(term.rhs.AsFieldName())
	}

	if term.lhs.AsConstant().Equal(term.rhs.AsConstant()) {
		return 1
	}

	return meta.MaxInt
}

// EquatesWithConstant ...
// F=c の形式かチェックする。ここで F は field, c は Constant。
func (term Term) EquatesWithConstant(fldName FieldName) Constant {
	lhs, rhs := term.lhs, term.rhs
	if lhs.IsFieldName() && lhs.AsFieldName() == fldName && !rhs.IsFieldName() {
		return rhs.AsConstant()
	} else if rhs.IsFieldName() && rhs.AsFieldName() == fldName && !lhs.IsFieldName() {
		return lhs.AsConstant()
	}

	return Constant{}
}

// EquatesWithField ...
// F1=F2 の形式かチェック.
func (term Term) EquatesWithField(fldName FieldName) FieldName {
	lhs, rhs := term.lhs, term.rhs
	if lhs.IsFieldName() && lhs.AsFieldName() == fldName && rhs.IsFieldName() {
		return rhs.AsFieldName()
	} else if rhs.IsFieldName() && rhs.AsFieldName() == fldName && lhs.IsFieldName() {
		return lhs.AsFieldName()
	}

	return ""
}

// String stringfy the term.
func (term Term) String() string {
	return term.lhs.String() + "=" + term.rhs.String()
}
