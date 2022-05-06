package domain

// Expression is node of expression.
type Expression struct {
	constant Constant
	field    Field
}

// NewConstExpression constructs a const expression.
func NewConstExpression(c Constant) Expression {
	return Expression{constant: c}
}

// NewFieldExpression constructs a field expression.
func NewFieldExpression(fld Field) Expression {
	return Expression{field: fld}
}
