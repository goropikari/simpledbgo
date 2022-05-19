package domain

import (
	"strings"

	"github.com/goropikari/simpledbgo/common"
	"github.com/goropikari/simpledbgo/math"
)

// Expression is node of expression.
type Expression struct {
	value Constant
	field FieldName
}

// NewConstExpression constructs a const expression.
func NewConstExpression(c Constant) Expression {
	return Expression{value: c}
}

// NewFieldNameExpression constructs a field expression.
func NewFieldNameExpression(fld FieldName) Expression {
	return Expression{field: fld}
}

// Evaluate evaluates scanner.
func (expr Expression) Evaluate(s Scanner) (Constant, error) {
	if expr.value.IsZero() {
		return s.GetVal(expr.field)
	}

	return expr.value, nil
}

// IsFieldName checks whether expr is field name or not.
func (expr Expression) IsFieldName() bool {
	return !expr.field.IsZero()
}

// AsFieldName returns expr as FieldName.
func (expr Expression) AsFieldName() FieldName {
	return expr.field
}

// AsConstant returns expr as Constant.
func (expr Expression) AsConstant() Constant {
	return expr.value
}

// String stringfies expr.
func (expr Expression) String() string {
	if expr.value.IsZero() {
		return expr.field.String()
	}

	return expr.value.String()
}

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

	return common.MaxInt
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

// String stringfies the term.
func (term Term) String() string {
	return term.lhs.String() + "=" + term.rhs.String()
}

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

// IsSatisfied checks whether a term is satisfied or not.
func (pred *Predicate) IsSatisfied(s Scanner) bool {
	for _, term := range pred.terms {
		if !term.IsSatisfied(s) {
			return false
		}
	}

	return true
}

// ReductionFactor ...
func (pred *Predicate) ReductionFactor(p Planner) int {
	factor := 1
	for _, term := range pred.terms {
		factor *= term.ReductionFactor(p)
	}

	return factor
}

// EquatesWithConstant ...
func (pred *Predicate) EquatesWithConstant(fldName FieldName) Constant {
	for _, term := range pred.terms {
		c := term.EquatesWithConstant(fldName)
		if !c.IsZero() {
			return c
		}
	}

	return Constant{}
}

// EquatesWithField ...
func (pred *Predicate) EquatesWithField(fldName FieldName) FieldName {
	for _, term := range pred.terms {
		if s := term.EquatesWithField(fldName); !s.IsZero() {
			return s
		}
	}

	return ""
}

// String stringfies predicate.
func (pred *Predicate) String() string {
	if len(pred.terms) == 0 {
		return ""
	}

	conds := make([]string, 0, len(pred.terms))
	for _, term := range pred.terms {
		conds = append(conds, term.String())
	}

	return strings.Join(conds, " and ")
}
