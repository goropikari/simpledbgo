package domain

import "strings"

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

// String stringfy predicate.
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
