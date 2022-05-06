package domain

// Expression is node of expression.
type Expression struct {
	constant Constant
	field    FieldName
}

// NewConstExpression constructs a const expression.
func NewConstExpression(c Constant) Expression {
	return Expression{constant: c}
}

// NewFieldNameExpression constructs a field expression.
func NewFieldNameExpression(fld FieldName) Expression {
	return Expression{field: fld}
}
