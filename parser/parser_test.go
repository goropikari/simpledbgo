package parser_test

import (
	"testing"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/lexer"
	"github.com/goropikari/simpledbgo/parser"
	"github.com/stretchr/testify/require"
)

func TestParser_Query(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []lexer.Token
		expected *domain.QueryData
	}{
		{
			name: "parse select",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "select"),
				lexer.NewToken(lexer.TStar, "*"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TIdentifier, "foo_bar"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "fizz_baz"),
				lexer.NewToken(lexer.TKeyword, "where"),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TInt32, int32(123)),
				lexer.NewToken(lexer.TKeyword, "and"),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TString, "Mike's dog"),
			},
			expected: domain.NewQueryData(
				[]domain.FieldName{"*", "id", "name"},
				[]domain.TableName{"foo_bar", "fizz_baz"},
				domain.NewPredicate([]domain.Term{
					domain.NewTerm(
						domain.NewFieldNameExpression("id"),
						domain.NewConstExpression(domain.NewConstant(domain.Int32FieldType, int32(123))),
					),
					domain.NewTerm(
						domain.NewFieldNameExpression("name"),
						domain.NewConstExpression(domain.NewConstant(domain.StringFieldType, "Mike's dog")),
					),
				}),
			),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			p := parser.NewParser(tt.tokens)
			q, err := p.Query()

			require.NoError(t, err)
			require.Equal(t, tt.expected, q)
		})
	}
}

func TestParser_Error(t *testing.T) {
	tests := []struct {
		name   string
		tokens []lexer.Token
	}{
		{
			name: "missing select",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "where"),
			},
		},
		{
			name: "error at select list",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "select"),
				lexer.NewToken(lexer.TStar, "*"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
			},
		},
		{
			name: "missing from",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "select"),
				lexer.NewToken(lexer.TStar, "*"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
			},
		},
		{
			name: "error at table list",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "select"),
				lexer.NewToken(lexer.TStar, "*"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TComma, ","),
			},
		},
		{
			name: "error at predicate",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "select"),
				lexer.NewToken(lexer.TStar, "*"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TIdentifier, "foo_bar"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "fizz_baz"),
				lexer.NewToken(lexer.TKeyword, "where"),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TKeyword, "and"),
			},
		},
		{
			name: "missing lhs",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "select"),
				lexer.NewToken(lexer.TStar, "*"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TIdentifier, "foo_bar"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "fizz_baz"),
				lexer.NewToken(lexer.TKeyword, "where"),
				lexer.NewToken(lexer.TEqual, "="),
			},
		},
		{
			name: "missing =",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "select"),
				lexer.NewToken(lexer.TStar, "*"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TIdentifier, "foo_bar"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "fizz_baz"),
				lexer.NewToken(lexer.TKeyword, "where"),
				lexer.NewToken(lexer.TIdentifier, "id"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			p := parser.NewParser(tt.tokens)
			_, err := p.Query()

			require.Error(t, err)
		})
	}
}

func TestParser_ExecCmd_Insert(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []lexer.Token
		expected *domain.InsertData
	}{
		{
			name: "parse insert",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "insert"),
				lexer.NewToken(lexer.TKeyword, "into"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TRParen, ")"),
				lexer.NewToken(lexer.TKeyword, "values"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TInt32, int32(123)),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TString, "mike"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
			expected: domain.NewInsertData(
				domain.TableName("foo"),
				[]domain.FieldName{"id", "name"},
				[]domain.Constant{
					domain.NewConstant(domain.Int32FieldType, int32(123)),
					domain.NewConstant(domain.StringFieldType, "mike"),
				},
			),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			p := parser.NewParser(tt.tokens)
			cmd, err := p.ExecCmd()

			require.NoError(t, err)

			got, ok := cmd.(*domain.InsertData)
			require.True(t, ok)

			require.Equal(t, tt.expected, got)
		})
	}
}

func TestParser_ExecCmd_Insert_Error(t *testing.T) {
	tests := []struct {
		name   string
		tokens []lexer.Token
	}{
		{
			name: "missing into",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "insert"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
			},
		},
		{
			name: "missing table name",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "insert"),
				lexer.NewToken(lexer.TKeyword, "into"),
				lexer.NewToken(lexer.TLParen, "("),
			},
		},
		{
			name: "missing first left paren",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "insert"),
				lexer.NewToken(lexer.TKeyword, "into"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TIdentifier, "id"),
			},
		},
		{
			name: "missing field",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "insert"),
				lexer.NewToken(lexer.TKeyword, "into"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TComma, ","),
			},
		},
		{
			name: "missing comma",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "insert"),
				lexer.NewToken(lexer.TKeyword, "into"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TIdentifier, "name"),
			},
		},
		{
			name: "missing first right paren",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "insert"),
				lexer.NewToken(lexer.TKeyword, "into"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "values"),
			},
		},
		{
			name: "missing values",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "insert"),
				lexer.NewToken(lexer.TKeyword, "into"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TRParen, ")"),
				lexer.NewToken(lexer.TLParen, "("),
			},
		},
		{
			name: "missing second left paren",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "insert"),
				lexer.NewToken(lexer.TKeyword, "into"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TRParen, ")"),
				lexer.NewToken(lexer.TKeyword, "values"),
				lexer.NewToken(lexer.TInt32, int32(123)),
			},
		},
		{
			name: "missing value",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "insert"),
				lexer.NewToken(lexer.TKeyword, "into"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TRParen, ")"),
				lexer.NewToken(lexer.TKeyword, "values"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TComma, ","),
			},
		},
		{
			name: "missing second rigth paren",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "insert"),
				lexer.NewToken(lexer.TKeyword, "into"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TRParen, ")"),
				lexer.NewToken(lexer.TKeyword, "values"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TInt32, int32(123)),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TString, "mike"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			p := parser.NewParser(tt.tokens)
			_, err := p.ExecCmd()

			require.Error(t, err)
		})
	}
}

func TestParser_ExecCmd_Delete(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []lexer.Token
		expected *domain.DeleteData
	}{
		{
			name: "parse delete",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "delete"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
			},
			expected: domain.NewDeleteData(domain.TableName("foo"), &domain.Predicate{}),
		},
		{
			name: "parse delete with predicate",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "delete"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TKeyword, "where"),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TInt32, int32(123)),
			},
			expected: domain.NewDeleteData(
				domain.TableName("foo"),
				domain.NewPredicate([]domain.Term{
					domain.NewTerm(
						domain.NewFieldNameExpression("id"),
						domain.NewConstExpression(
							domain.NewConstant(domain.Int32FieldType, int32(123)),
						),
					),
				}),
			),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			p := parser.NewParser(tt.tokens)
			cmd, err := p.ExecCmd()

			require.NoError(t, err)

			got, ok := cmd.(*domain.DeleteData)
			require.True(t, ok)

			require.Equal(t, tt.expected, got)
		})
	}
}

func TestParser_ExecCmd_Delete_Error(t *testing.T) {
	tests := []struct {
		name   string
		tokens []lexer.Token
	}{
		{
			name: "missing delete",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "from"),
			},
		},
		{
			name: "missing from",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "delete"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
			},
		},
		{
			name: "missing table name",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "delete"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TKeyword, "where"),
			},
		},
		{
			name: "missing predicate",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "delete"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TKeyword, "where"),
				lexer.NewToken(lexer.TEqual, "="),
			},
		},
		{
			name: "missing predicate",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "delete"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TKeyword, "where"),
				lexer.NewToken(lexer.TEqual, "="),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			p := parser.NewParser(tt.tokens)
			_, err := p.ExecCmd()

			require.Error(t, err)
		})
	}
}

func TestParser_ExecCmd_Modify(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []lexer.Token
		expected *domain.ModifyData
	}{
		{
			name: "parse update cmd",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "update"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TKeyword, "set"),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TString, "mike"),
			},
			expected: domain.NewModifyData(
				domain.TableName("foo"),
				domain.FieldName("name"),
				domain.NewConstExpression(
					domain.NewConstant(
						domain.StringFieldType,
						"mike",
					),
				),
				&domain.Predicate{},
			),
		},
		{
			name: "parse update cmd with predicate",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "update"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TKeyword, "set"),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TString, "mike"),
				lexer.NewToken(lexer.TKeyword, "where"),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TString, "neko"),
			},
			expected: domain.NewModifyData(
				domain.TableName("foo"),
				domain.FieldName("name"),
				domain.NewConstExpression(
					domain.NewConstant(
						domain.StringFieldType,
						"mike",
					),
				),
				domain.NewPredicate(
					[]domain.Term{
						domain.NewTerm(
							domain.NewFieldNameExpression(domain.FieldName("name")),
							domain.NewConstExpression(
								domain.NewConstant(
									domain.StringFieldType,
									"neko",
								),
							),
						),
					},
				),
			),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			p := parser.NewParser(tt.tokens)
			cmd, err := p.ExecCmd()

			require.NoError(t, err)

			got, ok := cmd.(*domain.ModifyData)
			require.True(t, ok)

			require.Equal(t, tt.expected, got)
		})
	}
}

func TestParser_ExecCmd_Modify_Error(t *testing.T) {
	tests := []struct {
		name   string
		tokens []lexer.Token
	}{
		{
			name: "missing update",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TIdentifier, "foo"),
			},
		},
		{
			name: "missing table name",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "update"),
				lexer.NewToken(lexer.TKeyword, "set"),
			},
		},
		{
			name: "missing set",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "update"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TIdentifier, "name"),
			},
		},
		{
			name: "missing field",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "update"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TKeyword, "set"),
				lexer.NewToken(lexer.TEqual, "="),
			},
		},
		{
			name: "missing equal",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "update"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TKeyword, "set"),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TString, "mike"),
			},
		},
		{
			name: "missing value",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "update"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TKeyword, "set"),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TEqual, "="),
			},
		},
		{
			name: "missing predicate field",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "update"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TKeyword, "set"),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TString, "mike"),
				lexer.NewToken(lexer.TKeyword, "where"),
				lexer.NewToken(lexer.TIdentifier, "name"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			p := parser.NewParser(tt.tokens)
			_, err := p.ExecCmd()

			require.Error(t, err)
		})
	}
}

func TestParser_ExecCmd_CreateTable(t *testing.T) {
	sch := domain.NewSchema()
	sch.AddInt32Field("id")
	sch.AddStringField("name", 255)

	tests := []struct {
		name     string
		tokens   []lexer.Token
		expected *domain.CreateTableData
	}{
		{
			name: "parse create table",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "table"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TKeyword, "int"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "varchar"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TInt32, int32(255)),
				lexer.NewToken(lexer.TRParen, ")"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
			expected: domain.NewCreateTableData(
				domain.TableName("foo"),
				sch,
			),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			p := parser.NewParser(tt.tokens)
			cmd, err := p.ExecCmd()

			require.NoError(t, err)

			got, ok := cmd.(*domain.CreateTableData)
			require.True(t, ok)

			require.Equal(t, tt.expected, got)
		})
	}
}

func TestParser_ExecCmd_CreateTable_Error(t *testing.T) {
	sch := domain.NewSchema()
	sch.AddInt32Field("id")
	sch.AddStringField("name", 255)

	tests := []struct {
		name   string
		tokens []lexer.Token
	}{
		{
			name: "missing create",
			tokens: []lexer.Token{
				// lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "table"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TKeyword, "int"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "varchar"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TInt32, int32(255)),
				lexer.NewToken(lexer.TRParen, ")"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
		},
		{
			name: "missing table",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				// lexer.NewToken(lexer.TKeyword, "table"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TKeyword, "int"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "varchar"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TInt32, int32(255)),
				lexer.NewToken(lexer.TRParen, ")"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
		},
		{
			name: "missing table name",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "table"),
				// lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TKeyword, "int"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "varchar"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TInt32, int32(255)),
				lexer.NewToken(lexer.TRParen, ")"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
		},
		{
			name: "missing first left paren",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "table"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				// lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TKeyword, "int"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "varchar"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TInt32, int32(255)),
				lexer.NewToken(lexer.TRParen, ")"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
		},
		{
			name: "missing identifier",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "table"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				// lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TKeyword, "int"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "varchar"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TInt32, int32(255)),
				lexer.NewToken(lexer.TRParen, ")"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
		},
		{
			name: "missing type",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "table"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				// lexer.NewToken(lexer.TKeyword, "int"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "varchar"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TInt32, int32(255)),
				lexer.NewToken(lexer.TRParen, ")"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
		},
		{
			name: "missing comma",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "table"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TKeyword, "int"),
				// lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "varchar"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TInt32, int32(255)),
				lexer.NewToken(lexer.TRParen, ")"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
		},
		{
			name: "missing varchar left paren",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "table"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TKeyword, "int"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "varchar"),
				// lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TInt32, int32(255)),
				lexer.NewToken(lexer.TRParen, ")"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
		},
		{
			name: "missing varchar size",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "table"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TKeyword, "int"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "varchar"),
				lexer.NewToken(lexer.TLParen, "("),
				// lexer.NewToken(lexer.TInt32, int32(255)),
				lexer.NewToken(lexer.TRParen, ")"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
		},
		{
			name: "missing varchar right paren",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "table"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TKeyword, "int"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "varchar"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TInt32, int32(255)),
				// lexer.NewToken(lexer.TRParen, ")"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			p := parser.NewParser(tt.tokens)
			_, err := p.ExecCmd()

			require.Error(t, err)
		})
	}
}

func TestParser_ExecCmd_CreateView(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []lexer.Token
		expected *domain.CreateViewData
	}{
		{
			name: "parse create table",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "view"),
				lexer.NewToken(lexer.TIdentifier, "view_foo"),
				lexer.NewToken(lexer.TKeyword, "as"),
				lexer.NewToken(lexer.TKeyword, "select"),
				lexer.NewToken(lexer.TStar, "*"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TIdentifier, "foo_bar"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "fizz_baz"),
				lexer.NewToken(lexer.TKeyword, "where"),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TInt32, int32(123)),
				lexer.NewToken(lexer.TKeyword, "and"),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TString, "Mike's dog"),
			},
			expected: domain.NewCreateViewData(
				domain.ViewName("view_foo"),
				domain.NewQueryData(
					[]domain.FieldName{"*", "id", "name"},
					[]domain.TableName{"foo_bar", "fizz_baz"},
					domain.NewPredicate([]domain.Term{
						domain.NewTerm(
							domain.NewFieldNameExpression("id"),
							domain.NewConstExpression(domain.NewConstant(domain.Int32FieldType, int32(123))),
						),
						domain.NewTerm(
							domain.NewFieldNameExpression("name"),
							domain.NewConstExpression(domain.NewConstant(domain.StringFieldType, "Mike's dog")),
						),
					}),
				),
			),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			p := parser.NewParser(tt.tokens)
			cmd, err := p.ExecCmd()

			require.NoError(t, err)

			got, ok := cmd.(*domain.CreateViewData)
			require.True(t, ok)

			require.Equal(t, tt.expected, got)
		})
	}
}

func TestParser_ExecCmd_CreateView_Error(t *testing.T) {
	tests := []struct {
		name   string
		tokens []lexer.Token
	}{
		{
			name: "missing create",
			tokens: []lexer.Token{
				// lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "view"),
				lexer.NewToken(lexer.TIdentifier, "view_foo"),
				lexer.NewToken(lexer.TKeyword, "as"),
				lexer.NewToken(lexer.TKeyword, "select"),
				lexer.NewToken(lexer.TStar, "*"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TIdentifier, "foo_bar"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "fizz_baz"),
				lexer.NewToken(lexer.TKeyword, "where"),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TInt32, int32(123)),
				lexer.NewToken(lexer.TKeyword, "and"),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TString, "Mike's dog"),
			},
		},
		{
			name: "missing view",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				// lexer.NewToken(lexer.TKeyword, "view"),
				lexer.NewToken(lexer.TIdentifier, "view_foo"),
				lexer.NewToken(lexer.TKeyword, "as"),
				lexer.NewToken(lexer.TKeyword, "select"),
				lexer.NewToken(lexer.TStar, "*"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TIdentifier, "foo_bar"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "fizz_baz"),
				lexer.NewToken(lexer.TKeyword, "where"),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TInt32, int32(123)),
				lexer.NewToken(lexer.TKeyword, "and"),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TString, "Mike's dog"),
			},
		},
		{
			name: "missing view name",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "view"),
				// lexer.NewToken(lexer.TIdentifier, "view_foo"),
				lexer.NewToken(lexer.TKeyword, "as"),
				lexer.NewToken(lexer.TKeyword, "select"),
				lexer.NewToken(lexer.TStar, "*"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TIdentifier, "foo_bar"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "fizz_baz"),
				lexer.NewToken(lexer.TKeyword, "where"),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TInt32, int32(123)),
				lexer.NewToken(lexer.TKeyword, "and"),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TString, "Mike's dog"),
			},
		},
		{
			name: "missing as",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "view"),
				lexer.NewToken(lexer.TIdentifier, "view_foo"),
				// lexer.NewToken(lexer.TKeyword, "as"),
				lexer.NewToken(lexer.TKeyword, "select"),
				lexer.NewToken(lexer.TStar, "*"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TIdentifier, "foo_bar"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "fizz_baz"),
				lexer.NewToken(lexer.TKeyword, "where"),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TInt32, int32(123)),
				lexer.NewToken(lexer.TKeyword, "and"),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TString, "Mike's dog"),
			},
		},
		{
			name: "missing select",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "view"),
				lexer.NewToken(lexer.TIdentifier, "view_foo"),
				lexer.NewToken(lexer.TKeyword, "as"),
				// lexer.NewToken(lexer.TKeyword, "select"),
				lexer.NewToken(lexer.TStar, "*"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TIdentifier, "foo_bar"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "fizz_baz"),
				lexer.NewToken(lexer.TKeyword, "where"),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TInt32, int32(123)),
				lexer.NewToken(lexer.TKeyword, "and"),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TString, "Mike's dog"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			p := parser.NewParser(tt.tokens)
			_, err := p.ExecCmd()

			require.Error(t, err)
		})
	}
}

func TestParser_ExecCmd_CreateIndex(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []lexer.Token
		expected *domain.CreateIndexData
	}{
		{
			name: "parse create index",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "index"),
				lexer.NewToken(lexer.TIdentifier, "idx_id"),
				lexer.NewToken(lexer.TKeyword, "on"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
			expected: domain.NewCreateIndexData(
				domain.IndexName("idx_id"),
				domain.TableName("foo"),
				domain.FieldName("id"),
			),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			p := parser.NewParser(tt.tokens)
			cmd, err := p.ExecCmd()

			require.NoError(t, err)

			got, ok := cmd.(*domain.CreateIndexData)
			require.True(t, ok)

			require.Equal(t, tt.expected, got)
		})
	}
}

func TestParser_ExecCmd_CreateIndex_Error(t *testing.T) {
	tests := []struct {
		name   string
		tokens []lexer.Token
	}{
		{
			name: "missing create",
			tokens: []lexer.Token{
				// lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "index"),
				lexer.NewToken(lexer.TIdentifier, "idx_id"),
				lexer.NewToken(lexer.TKeyword, "on"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
		},
		{
			name: "missing index",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				// lexer.NewToken(lexer.TKeyword, "index"),
				lexer.NewToken(lexer.TIdentifier, "idx_id"),
				lexer.NewToken(lexer.TKeyword, "on"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
		},
		{
			name: "missing index name",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "index"),
				// lexer.NewToken(lexer.TIdentifier, "idx_id"),
				lexer.NewToken(lexer.TKeyword, "on"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
		},
		{
			name: "missing on",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "index"),
				lexer.NewToken(lexer.TIdentifier, "idx_id"),
				// lexer.NewToken(lexer.TKeyword, "on"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
		},
		{
			name: "missing table name",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "index"),
				lexer.NewToken(lexer.TIdentifier, "idx_id"),
				lexer.NewToken(lexer.TKeyword, "on"),
				// lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
		},
		{
			name: "missing left paren",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "index"),
				lexer.NewToken(lexer.TIdentifier, "idx_id"),
				lexer.NewToken(lexer.TKeyword, "on"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				// lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
		},
		{
			name: "missing field name",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "index"),
				lexer.NewToken(lexer.TIdentifier, "idx_id"),
				lexer.NewToken(lexer.TKeyword, "on"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				// lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
		},
		{
			name: "missing right paren",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "create"),
				lexer.NewToken(lexer.TKeyword, "index"),
				lexer.NewToken(lexer.TIdentifier, "idx_id"),
				lexer.NewToken(lexer.TKeyword, "on"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				// lexer.NewToken(lexer.TRParen, ")"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			p := parser.NewParser(tt.tokens)
			_, err := p.ExecCmd()

			require.Error(t, err)
		})
	}
}
