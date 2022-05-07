package parser_test

import (
	"testing"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/frontend/parser"
	"github.com/stretchr/testify/require"
)

func TestParser_Query(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []domain.Token
		expected *domain.QueryData
	}{
		{
			name: "parse select",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "select"),
				domain.NewToken(domain.TStar, "*"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TKeyword, "from"),
				domain.NewToken(domain.TIdentifier, "foo_bar"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "fizz_baz"),
				domain.NewToken(domain.TKeyword, "where"),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TEqual, "="),
				domain.NewToken(domain.TInt32, int32(123)),
				domain.NewToken(domain.TKeyword, "and"),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TEqual, "="),
				domain.NewToken(domain.TString, "Mike's dog"),
			},
			expected: domain.NewQueryData(
				[]domain.FieldName{"*", "id", "name"},
				[]domain.TableName{"foo_bar", "fizz_baz"},
				domain.NewPredicate([]domain.Term{
					domain.NewTerm(
						domain.NewFieldNameExpression("id"),
						domain.NewConstExpression(domain.NewConstant(domain.VInt32, int32(123))),
					),
					domain.NewTerm(
						domain.NewFieldNameExpression("name"),
						domain.NewConstExpression(domain.NewConstant(domain.VString, "Mike's dog")),
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
		tokens []domain.Token
	}{
		{
			name: "missing select",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "where"),
			},
		},
		{
			name: "error at select list",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "select"),
				domain.NewToken(domain.TStar, "*"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TComma, ","),
			},
		},
		{
			name: "missing from",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "select"),
				domain.NewToken(domain.TStar, "*"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "name"),
			},
		},
		{
			name: "error at table list",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "select"),
				domain.NewToken(domain.TStar, "*"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TKeyword, "from"),
				domain.NewToken(domain.TComma, ","),
			},
		},
		{
			name: "error at predicate",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "select"),
				domain.NewToken(domain.TStar, "*"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TKeyword, "from"),
				domain.NewToken(domain.TIdentifier, "foo_bar"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "fizz_baz"),
				domain.NewToken(domain.TKeyword, "where"),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TEqual, "="),
				domain.NewToken(domain.TKeyword, "and"),
			},
		},
		{
			name: "missing lhs",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "select"),
				domain.NewToken(domain.TStar, "*"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TKeyword, "from"),
				domain.NewToken(domain.TIdentifier, "foo_bar"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "fizz_baz"),
				domain.NewToken(domain.TKeyword, "where"),
				domain.NewToken(domain.TEqual, "="),
			},
		},
		{
			name: "missing =",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "select"),
				domain.NewToken(domain.TStar, "*"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TKeyword, "from"),
				domain.NewToken(domain.TIdentifier, "foo_bar"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "fizz_baz"),
				domain.NewToken(domain.TKeyword, "where"),
				domain.NewToken(domain.TIdentifier, "id"),
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
		tokens   []domain.Token
		expected *domain.InsertData
	}{
		{
			name: "parse insert",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "insert"),
				domain.NewToken(domain.TKeyword, "into"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TLParen, "("),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TRParen, ")"),
				domain.NewToken(domain.TKeyword, "values"),
				domain.NewToken(domain.TLParen, "("),
				domain.NewToken(domain.TInt32, int32(123)),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TString, "mike"),
				domain.NewToken(domain.TRParen, ")"),
			},
			expected: domain.NewInsertData(
				domain.TableName("foo"),
				[]domain.FieldName{"id", "name"},
				[]domain.Constant{
					domain.NewConstant(domain.VInt32, int32(123)),
					domain.NewConstant(domain.VString, "mike"),
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
		tokens []domain.Token
	}{
		{
			name: "missing into",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "insert"),
				domain.NewToken(domain.TIdentifier, "foo"),
			},
		},
		{
			name: "missing table name",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "insert"),
				domain.NewToken(domain.TKeyword, "into"),
				domain.NewToken(domain.TLParen, "("),
			},
		},
		{
			name: "missing first left paren",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "insert"),
				domain.NewToken(domain.TKeyword, "into"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TIdentifier, "id"),
			},
		},
		{
			name: "missing field",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "insert"),
				domain.NewToken(domain.TKeyword, "into"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TLParen, "("),
				domain.NewToken(domain.TComma, ","),
			},
		},
		{
			name: "missing comma",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "insert"),
				domain.NewToken(domain.TKeyword, "into"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TLParen, "("),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TIdentifier, "name"),
			},
		},
		{
			name: "missing first right paren",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "insert"),
				domain.NewToken(domain.TKeyword, "into"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TLParen, "("),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TKeyword, "values"),
			},
		},
		{
			name: "missing values",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "insert"),
				domain.NewToken(domain.TKeyword, "into"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TLParen, "("),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TRParen, ")"),
				domain.NewToken(domain.TLParen, "("),
			},
		},
		{
			name: "missing second left paren",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "insert"),
				domain.NewToken(domain.TKeyword, "into"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TLParen, "("),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TRParen, ")"),
				domain.NewToken(domain.TKeyword, "values"),
				domain.NewToken(domain.TInt32, int32(123)),
			},
		},
		{
			name: "missing value",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "insert"),
				domain.NewToken(domain.TKeyword, "into"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TLParen, "("),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TRParen, ")"),
				domain.NewToken(domain.TKeyword, "values"),
				domain.NewToken(domain.TLParen, "("),
				domain.NewToken(domain.TComma, ","),
			},
		},
		{
			name: "missing second rigth paren",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "insert"),
				domain.NewToken(domain.TKeyword, "into"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TLParen, "("),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TRParen, ")"),
				domain.NewToken(domain.TKeyword, "values"),
				domain.NewToken(domain.TLParen, "("),
				domain.NewToken(domain.TInt32, int32(123)),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TString, "mike"),
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
		tokens   []domain.Token
		expected *domain.DeleteData
	}{
		{
			name: "parse delete",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "delete"),
				domain.NewToken(domain.TKeyword, "from"),
				domain.NewToken(domain.TIdentifier, "foo"),
			},
			expected: domain.NewDeleteData(domain.TableName("foo"), nil),
		},
		{
			name: "parse delete with predicate",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "delete"),
				domain.NewToken(domain.TKeyword, "from"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TKeyword, "where"),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TEqual, "="),
				domain.NewToken(domain.TInt32, int32(123)),
			},
			expected: domain.NewDeleteData(
				domain.TableName("foo"),
				domain.NewPredicate([]domain.Term{
					domain.NewTerm(
						domain.NewFieldNameExpression("id"),
						domain.NewConstExpression(
							domain.NewConstant(domain.VInt32, int32(123)),
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
		tokens []domain.Token
	}{
		{
			name: "missing delete",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "from"),
			},
		},
		{
			name: "missing from",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "delete"),
				domain.NewToken(domain.TIdentifier, "foo"),
			},
		},
		{
			name: "missing table name",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "delete"),
				domain.NewToken(domain.TKeyword, "from"),
				domain.NewToken(domain.TKeyword, "where"),
			},
		},
		{
			name: "missing predicate",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "delete"),
				domain.NewToken(domain.TKeyword, "from"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TKeyword, "where"),
				domain.NewToken(domain.TEqual, "="),
			},
		},
		{
			name: "missing predicate",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "delete"),
				domain.NewToken(domain.TKeyword, "from"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TKeyword, "where"),
				domain.NewToken(domain.TEqual, "="),
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
		tokens   []domain.Token
		expected *domain.ModifyData
	}{
		{
			name: "parse update cmd",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "update"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TKeyword, "set"),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TEqual, "="),
				domain.NewToken(domain.TString, "mike"),
			},
			expected: domain.NewModifyData(
				domain.TableName("foo"),
				domain.FieldName("name"),
				domain.NewConstExpression(
					domain.NewConstant(
						domain.VString,
						"mike",
					),
				),
				nil,
			),
		},
		{
			name: "parse update cmd with predicate",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "update"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TKeyword, "set"),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TEqual, "="),
				domain.NewToken(domain.TString, "mike"),
				domain.NewToken(domain.TKeyword, "where"),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TEqual, "="),
				domain.NewToken(domain.TString, "neko"),
			},
			expected: domain.NewModifyData(
				domain.TableName("foo"),
				domain.FieldName("name"),
				domain.NewConstExpression(
					domain.NewConstant(
						domain.VString,
						"mike",
					),
				),
				domain.NewPredicate(
					[]domain.Term{
						domain.NewTerm(
							domain.NewFieldNameExpression(domain.FieldName("name")),
							domain.NewConstExpression(
								domain.NewConstant(
									domain.VString,
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
		tokens []domain.Token
	}{
		{
			name: "missing update",
			tokens: []domain.Token{
				domain.NewToken(domain.TIdentifier, "foo"),
			},
		},
		{
			name: "missing table name",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "update"),
				domain.NewToken(domain.TKeyword, "set"),
			},
		},
		{
			name: "missing set",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "update"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TIdentifier, "name"),
			},
		},
		{
			name: "missing field",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "update"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TKeyword, "set"),
				domain.NewToken(domain.TEqual, "="),
			},
		},
		{
			name: "missing equal",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "update"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TKeyword, "set"),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TString, "mike"),
			},
		},
		{
			name: "missing value",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "update"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TKeyword, "set"),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TEqual, "="),
			},
		},
		{
			name: "missing predicate field",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "update"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TKeyword, "set"),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TEqual, "="),
				domain.NewToken(domain.TString, "mike"),
				domain.NewToken(domain.TKeyword, "where"),
				domain.NewToken(domain.TIdentifier, "name"),
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
