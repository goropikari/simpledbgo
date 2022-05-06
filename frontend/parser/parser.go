package parser

import (
	"errors"

	"github.com/goropikari/simpledbgo/frontend/domain"
)

var (
	ErrParse = errors.New("parse error")
)

type Parser struct {
	tokens []domain.Token
	pos    int
	len    int
}

func NewParser(tokens []domain.Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
		len:    len(tokens),
	}
}

func (parser *Parser) term() (domain.Term, error) {
	lhs, err := parser.expression()
	if err != nil {
		return domain.Term{}, err
	}

	err = parser.eatToken(domain.TEqual)
	if err != nil {
		return domain.Term{}, err
	}

	rhs, err := parser.expression()
	if err != nil {
		return domain.Term{}, err
	}

	return domain.NewTerm(lhs, rhs), nil
}

func (parser *Parser) Query() (*domain.QueryData, error) {
	err := parser.eatKeyword("select")
	if err != nil {
		return nil, err
	}

	fields, err := parser.selectList()
	if err != nil {
		return nil, err
	}

	err = parser.eatKeyword("from")
	if err != nil {
		return nil, err
	}

	tables, err := parser.tableList()
	if err != nil {
		return nil, err
	}

	var pred domain.Predicate
	if parser.matchKeyword("where") {
		err = parser.eatKeyword("where")
		if err != nil {
			return nil, err
		}

		pred, err = parser.predicate()
		if err != nil {
			return nil, err
		}
	}

	return domain.NewQueryData(fields, tables, pred), nil
}

func (parser *Parser) selectList() ([]domain.Field, error) {
	fields := make([]domain.Field, 0)
	fld, err := parser.field()
	if err != nil {
		return nil, err
	}

	fields = append(fields, fld)

	for parser.match(domain.TComma) {
		err = parser.eatToken(domain.TComma)
		if err != nil {
			return nil, err
		}
		fld, err = parser.field()
		if err != nil {
			return nil, err
		}
		fields = append(fields, fld)
	}

	return fields, nil
}

func (parser *Parser) tableList() ([]domain.TableName, error) {
	tables := make([]domain.TableName, 0)

	tbl, err := parser.table()
	if err != nil {
		return nil, err
	}

	tables = append(tables, tbl)

	for parser.match(domain.TComma) {
		err = parser.eatToken(domain.TComma)
		if err != nil {
			return nil, err
		}

		tbl, err = parser.table()
		if err != nil {
			return nil, err
		}

		tables = append(tables, tbl)
	}

	return tables, nil
}

func (parser *Parser) predicate() (domain.Predicate, error) {
	pred := domain.NewPredicate()
	term, err := parser.term()
	if err != nil {
		return domain.Predicate{}, err
	}
	pred.Add(term)

	for parser.matchKeyword("and") {
		err = parser.eatKeyword("and")
		if err != nil {
			return domain.Predicate{}, err
		}

		term, err := parser.term()
		if err != nil {
			return domain.Predicate{}, err
		}
		pred.Add(term)
	}

	return pred, nil
}

func (parser *Parser) field() (domain.Field, error) {
	if parser.match(domain.TStar) {
		err := parser.eatToken(domain.TStar)
		if err != nil {
			return "", err
		}

		return domain.NewField("*"), nil
	}

	id, err := parser.eatIdentifier()
	if err != nil {
		return "", err
	}

	return domain.NewField(id), nil
}

func (parser *Parser) table() (domain.TableName, error) {
	id, err := parser.eatIdentifier()
	if err != nil {
		return "", err
	}

	return domain.NewTableName(id), nil
}

func (parser *Parser) expression() (domain.Expression, error) {
	switch {
	case parser.match(domain.TIdentifier):
		id, err := parser.eatIdentifier()
		if err != nil {
			return domain.Expression{}, err
		}

		return domain.NewFieldExpression(domain.NewField(id)), nil
	case parser.match(domain.TString) || parser.match(domain.TInt32):
		c, err := parser.constant()
		if err != nil {
			return domain.Expression{}, err
		}

		return domain.NewConstExpression(c), nil
	default:
		return domain.Expression{}, ErrParse
	}
}

func (parser *Parser) constant() (domain.Constant, error) {
	switch {
	case parser.match(domain.TString):
		str, err := parser.eatString()
		if err != nil {
			return domain.Constant{}, ErrParse
		}

		return domain.NewConstant(domain.VString, str), nil
	case parser.match(domain.TInt32):
		num, err := parser.eatInt32()
		if err != nil {
			return domain.Constant{}, ErrParse
		}

		return domain.NewConstant(domain.VInt32, num), nil
	default:
		return domain.Constant{}, ErrParse
	}
}

func (parser *Parser) match(typ domain.TokenType) bool {
	if parser.pos >= parser.len {
		return false
	}

	return parser.tokens[parser.pos].Type() == typ
}

func (parser *Parser) matchKeyword(kw string) bool {
	return parser.match(domain.TKeyword) && parser.tokens[parser.pos].Value() == kw
}

func (parser *Parser) eatKeyword(kw string) error {
	if !parser.match(domain.TKeyword) {
		return ErrParse
	}
	if !(parser.tokens[parser.pos].Value() == kw) {
		return ErrParse
	}

	parser.pos++

	return nil
}

func (parser *Parser) eatIdentifier() (string, error) {
	if !parser.match(domain.TIdentifier) {
		return "", ErrParse
	}

	id := parser.tokens[parser.pos].Value().(string)
	parser.pos++

	return id, nil
}

func (parser *Parser) eatInt32() (int32, error) {
	if !parser.match(domain.TInt32) {
		return 0, ErrParse
	}

	num := parser.tokens[parser.pos].Value().(int32)
	parser.pos++

	return num, nil
}

func (parser *Parser) eatString() (string, error) {
	if !parser.match(domain.TString) {
		return "", ErrParse
	}

	str := parser.tokens[parser.pos].Value().(string)
	parser.pos++

	return str, nil
}

func (parser *Parser) eatToken(typ domain.TokenType) error {
	if !parser.match(typ) {
		return ErrParse
	}

	parser.pos++

	return nil
}
