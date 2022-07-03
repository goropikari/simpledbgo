package parser

import (
	"errors"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/lexer"
)

// ErrParse is a parse error.
var ErrParse = errors.New("parse error")

// Parser is a model of parser.
type Parser struct {
	tokens []lexer.Token
	pos    int
	len    int
}

// NewParser constructs a Parser.
func NewParser(tokens []lexer.Token) *Parser {
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

	err = parser.eatToken(lexer.TEqual)
	if err != nil {
		return domain.Term{}, err
	}

	rhs, err := parser.expression()
	if err != nil {
		return domain.Term{}, err
	}

	return domain.NewTerm(lhs, rhs), nil
}

// Query connstructs a query parse tree.
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

	pred := &domain.Predicate{}
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

func (parser *Parser) selectList() ([]domain.FieldName, error) {
	fields := make([]domain.FieldName, 0)
	fld, err := parser.field()
	if err != nil {
		return nil, err
	}

	fields = append(fields, fld)

	for parser.match(lexer.TComma) {
		err = parser.eatToken(lexer.TComma)
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

	for parser.match(lexer.TComma) {
		err = parser.eatToken(lexer.TComma)
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

// ExecCmd parses execution command.
func (parser *Parser) ExecCmd() (domain.ExecData, error) {
	switch {
	case parser.matchKeyword("insert"):
		return parser.insertCmd()
	case parser.matchKeyword("delete"):
		return parser.deleteCmd()
	case parser.matchKeyword("update"):
		return parser.modifyCmd()
	case parser.matchKeyword("create"):
		return parser.createCmd()
	default:
		return nil, ErrParse
	}
}

func (parser *Parser) insertCmd() (domain.ExecData, error) {
	err := parser.eatKeyword("insert")
	if err != nil {
		return nil, err
	}

	err = parser.eatKeyword("into")
	if err != nil {
		return nil, err
	}

	tblNameStr, err := parser.eatIdentifier()
	if err != nil {
		return nil, err
	}
	tblName, err := domain.NewTableName(tblNameStr)
	if err != nil {
		return nil, err
	}

	err = parser.eatToken(lexer.TLParen)
	if err != nil {
		return nil, err
	}

	flds, err := parser.fieldList()
	if err != nil {
		return nil, err
	}

	err = parser.eatToken(lexer.TRParen)
	if err != nil {
		return nil, err
	}

	err = parser.eatKeyword("values")
	if err != nil {
		return nil, err
	}

	err = parser.eatToken(lexer.TLParen)
	if err != nil {
		return nil, err
	}

	vals, err := parser.constList()
	if err != nil {
		return nil, err
	}

	err = parser.eatToken(lexer.TRParen)
	if err != nil {
		return nil, err
	}

	return domain.NewInsertData(tblName, flds, vals), nil
}

func (parser *Parser) deleteCmd() (domain.ExecData, error) {
	err := parser.eatKeyword("delete")
	if err != nil {
		return nil, err
	}

	err = parser.eatKeyword("from")
	if err != nil {
		return nil, err
	}

	tblStr, err := parser.eatIdentifier()
	if err != nil {
		return nil, err
	}

	tblName, err := domain.NewTableName(tblStr)
	if err != nil {
		return nil, err
	}

	pred := &domain.Predicate{}
	if parser.matchKeyword("where") {
		err := parser.eatKeyword("where")
		if err != nil {
			return nil, err
		}

		pred, err = parser.predicate()
		if err != nil {
			return nil, err
		}
	}

	return domain.NewDeleteData(tblName, pred), nil
}

func (parser *Parser) modifyCmd() (domain.ExecData, error) {
	err := parser.eatKeyword("update")
	if err != nil {
		return nil, err
	}

	tblStr, err := parser.eatIdentifier()
	if err != nil {
		return nil, err
	}

	tblName, err := domain.NewTableName(tblStr)
	if err != nil {
		return nil, err
	}

	err = parser.eatKeyword("set")
	if err != nil {
		return nil, err
	}

	fld, err := parser.field()
	if err != nil {
		return nil, err
	}

	err = parser.eatToken(lexer.TEqual)
	if err != nil {
		return nil, err
	}

	expr, err := parser.expression()
	if err != nil {
		return nil, err
	}

	pred := &domain.Predicate{}
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

	return domain.NewModifyData(tblName, fld, expr, pred), nil
}

func (parser *Parser) createCmd() (domain.ExecData, error) {
	err := parser.eatKeyword("create")
	if err != nil {
		return nil, err
	}

	switch {
	case parser.matchKeyword("table"):
		return parser.createTable()
	case parser.matchKeyword("view"):
		return parser.createView()
	case parser.matchKeyword("index"):
		return parser.createIndex()
	default:
		return nil, ErrParse
	}
}

func (parser *Parser) createTable() (domain.ExecData, error) {
	err := parser.eatKeyword("table")
	if err != nil {
		return nil, err
	}

	tblStr, err := parser.eatIdentifier()
	if err != nil {
		return nil, err
	}
	tblName, err := domain.NewTableName(tblStr)
	if err != nil {
		return nil, err
	}

	err = parser.eatToken(lexer.TLParen)
	if err != nil {
		return nil, err
	}

	sch, err := parser.fieldDefs()
	if err != nil {
		return nil, err
	}

	err = parser.eatToken(lexer.TRParen)
	if err != nil {
		return nil, err
	}

	return domain.NewCreateTableData(tblName, sch), nil
}

func (parser *Parser) fieldDefs() (*domain.Schema, error) {
	sch, err := parser.fieldDef()
	if err != nil {
		return nil, err
	}

	for parser.match(lexer.TComma) {
		err = parser.eatToken(lexer.TComma)
		if err != nil {
			return nil, err
		}

		sch2, err := parser.fieldDef()
		if err != nil {
			return nil, err
		}

		sch.AddAllFields(sch2)
	}

	return sch, nil
}

func (parser *Parser) fieldDef() (*domain.Schema, error) {
	fld, err := parser.field()
	if err != nil {
		return nil, err
	}

	return parser.fieldType(fld)
}

func (parser *Parser) fieldType(fld domain.FieldName) (*domain.Schema, error) {
	sch := domain.NewSchema()

	switch {
	case parser.matchKeyword("int"):
		err := parser.eatKeyword("int")
		if err != nil {
			return nil, err
		}
		sch.AddInt32Field(fld)
	case parser.matchKeyword("varchar"):
		err := parser.eatKeyword("varchar")
		if err != nil {
			return nil, err
		}

		err = parser.eatToken(lexer.TLParen)
		if err != nil {
			return nil, err
		}

		num, err := parser.eatInt32()
		if err != nil {
			return nil, err
		}

		err = parser.eatToken(lexer.TRParen)
		if err != nil {
			return nil, err
		}

		sch.AddStringField(fld, int(num))
	default:
		return nil, ErrParse
	}

	return sch, nil
}

func (parser *Parser) createView() (domain.ExecData, error) {
	err := parser.eatKeyword("view")
	if err != nil {
		return nil, err
	}

	viewNameStr, err := parser.eatIdentifier()
	if err != nil {
		return nil, err
	}

	viewName, err := domain.NewViewName(viewNameStr)
	if err != nil {
		return nil, err
	}

	err = parser.eatKeyword("as")
	if err != nil {
		return nil, err
	}

	qd, err := parser.Query()
	if err != nil {
		return nil, err
	}

	return domain.NewCreateViewData(viewName, qd), nil
}

func (parser *Parser) createIndex() (domain.ExecData, error) {
	err := parser.eatKeyword("index")
	if err != nil {
		return nil, err
	}

	idxNameStr, err := parser.eatIdentifier()
	if err != nil {
		return nil, err
	}
	idxName, err := domain.NewIndexName(idxNameStr)
	if err != nil {
		return nil, err
	}

	err = parser.eatKeyword("on")
	if err != nil {
		return nil, err
	}

	tblNameStr, err := parser.eatIdentifier()
	if err != nil {
		return nil, err
	}
	tblName, err := domain.NewTableName(tblNameStr)
	if err != nil {
		return nil, err
	}

	err = parser.eatToken(lexer.TLParen)
	if err != nil {
		return nil, err
	}

	fldNameStr, err := parser.eatIdentifier()
	if err != nil {
		return nil, err
	}
	fldName, err := domain.NewFieldName(fldNameStr)
	if err != nil {
		return nil, err
	}

	err = parser.eatToken(lexer.TRParen)
	if err != nil {
		return nil, err
	}

	return domain.NewCreateIndexData(idxName, tblName, fldName), nil
}

func (parser *Parser) fieldList() ([]domain.FieldName, error) {
	fields := make([]domain.FieldName, 0)
	fld, err := parser.field()
	if err != nil {
		return nil, err
	}

	fields = append(fields, fld)

	for parser.match(lexer.TComma) {
		err = parser.eatToken(lexer.TComma)
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

func (parser *Parser) constList() ([]domain.Constant, error) {
	consts := make([]domain.Constant, 0)
	cons, err := parser.constant()
	if err != nil {
		return nil, err
	}

	consts = append(consts, cons)

	for parser.match(lexer.TComma) {
		err = parser.eatToken(lexer.TComma)
		if err != nil {
			return nil, err
		}
		cons, err = parser.constant()
		if err != nil {
			return nil, err
		}
		consts = append(consts, cons)
	}

	return consts, nil
}

func (parser *Parser) predicate() (*domain.Predicate, error) {
	terms := make([]domain.Term, 0)
	term, err := parser.term()
	if err != nil {
		return &domain.Predicate{}, err
	}

	terms = append(terms, term)

	for parser.matchKeyword("and") {
		err = parser.eatKeyword("and")
		if err != nil {
			return &domain.Predicate{}, err
		}

		term, err := parser.term()
		if err != nil {
			return &domain.Predicate{}, err
		}
		terms = append(terms, term)
	}

	return domain.NewPredicate(terms), nil
}

func (parser *Parser) field() (domain.FieldName, error) {
	if parser.match(lexer.TStar) {
		err := parser.eatToken(lexer.TStar)
		if err != nil {
			return "", err
		}

		return domain.NewFieldName("*")
	}

	id, err := parser.eatIdentifier()
	if err != nil {
		return "", err
	}

	return domain.NewFieldName(id)
}

func (parser *Parser) table() (domain.TableName, error) {
	id, err := parser.eatIdentifier()
	if err != nil {
		return "", err
	}

	return domain.NewTableName(id)
}

func (parser *Parser) expression() (domain.Expression, error) {
	switch {
	case parser.match(lexer.TIdentifier):
		id, err := parser.eatIdentifier()
		if err != nil {
			return domain.Expression{}, err
		}

		fldName, err := domain.NewFieldName(id)
		if err != nil {
			return domain.Expression{}, err
		}

		return domain.NewFieldNameExpression(fldName), nil
	case parser.match(lexer.TString) || parser.match(lexer.TInt32):
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
	case parser.match(lexer.TString):
		str, err := parser.eatString()
		if err != nil {
			return domain.Constant{}, err
		}

		return domain.NewConstant(domain.FString, str), nil
	case parser.match(lexer.TInt32):
		num, err := parser.eatInt32()
		if err != nil {
			return domain.Constant{}, err
		}

		return domain.NewConstant(domain.FInt32, num), nil
	default:
		return domain.Constant{}, ErrParse
	}
}

func (parser *Parser) match(typ lexer.TokenType) bool {
	if parser.pos >= parser.len {
		return false
	}

	return parser.tokens[parser.pos].Type() == typ
}

func (parser *Parser) matchKeyword(kw string) bool {
	return parser.match(lexer.TKeyword) && parser.tokens[parser.pos].Value() == kw
}

func (parser *Parser) eatKeyword(kw string) error {
	if !parser.match(lexer.TKeyword) {
		return ErrParse
	}
	if !(parser.tokens[parser.pos].Value() == kw) {
		return ErrParse
	}

	parser.pos++

	return nil
}

func (parser *Parser) eatIdentifier() (string, error) {
	if !parser.match(lexer.TIdentifier) {
		return "", ErrParse
	}

	id, _ := parser.tokens[parser.pos].Value().(string)
	parser.pos++

	return id, nil
}

func (parser *Parser) eatInt32() (int32, error) {
	if !parser.match(lexer.TInt32) {
		return 0, ErrParse
	}

	num, _ := parser.tokens[parser.pos].Value().(int32)
	parser.pos++

	return num, nil
}

func (parser *Parser) eatString() (string, error) {
	if !parser.match(lexer.TString) {
		return "", ErrParse
	}

	str, _ := parser.tokens[parser.pos].Value().(string)
	parser.pos++

	return str, nil
}

func (parser *Parser) eatToken(typ lexer.TokenType) error {
	if !parser.match(typ) {
		return ErrParse
	}

	parser.pos++

	return nil
}
