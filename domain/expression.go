package domain

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

// String stringfy expr.
func (expr Expression) String() string {
	if expr.value.IsZero() {
		return expr.field.String()
	}

	return expr.value.AsString()
}
