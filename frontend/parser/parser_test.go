package parser_test

import (
	"testing"

	"github.com/goropikari/simpledbgo/frontend/domain"
	"github.com/goropikari/simpledbgo/frontend/parser"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {

	pred := domain.NewPredicate()
	pred.Add(domain.NewTerm(
		domain.NewFieldExpression("id"),
		domain.NewConstExpression(domain.NewConstant(domain.VInt32, int32(123))),
	))
	pred.Add(domain.NewTerm(
		domain.NewFieldExpression("name"),
		domain.NewConstExpression(domain.NewConstant(domain.VString, "Mike's dog")),
	))

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
				[]domain.Field{"*", "id", "name"},
				[]domain.TableName{"foo_bar", "fizz_baz"},
				pred,
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
