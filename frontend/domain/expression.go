package domain

type Expression struct {
	constant Constant
	field    Field
}

func NewConstExpression(c Constant) Expression {
	return Expression{constant: c}
}

func NewFieldExpression(fld Field) Expression {
	return Expression{field: fld}
}
