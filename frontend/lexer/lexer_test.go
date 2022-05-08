package lexer_test

import (
	"testing"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/frontend/lexer"
	"github.com/stretchr/testify/require"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		tokens []domain.Token
	}{
		{
			name:  "select query",
			query: "SELECT *, id, name FROM foo_bar WHERE id = 123 and name = 'Mike\\'s dog'",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "select"),
				domain.NewToken(domain.TStar, "*"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TKeyword, "from"),
				domain.NewToken(domain.TIdentifier, "foo_bar"),
				domain.NewToken(domain.TKeyword, "where"),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TEqual, "="),
				domain.NewToken(domain.TInt32, int32(123)),
				domain.NewToken(domain.TKeyword, "and"),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TEqual, "="),
				domain.NewToken(domain.TString, "Mike's dog"),
			},
		},
		{
			name:  "select query with combination upper/lower",
			query: "SeLEcT * from foo",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "select"),
				domain.NewToken(domain.TStar, "*"),
				domain.NewToken(domain.TKeyword, "from"),
				domain.NewToken(domain.TIdentifier, "foo"),
			},
		},
		{
			name:  "insert command",
			query: "INSERT INTO foo (id,name, address) VALUES (123, 'mike','tokyo')",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "insert"),
				domain.NewToken(domain.TKeyword, "into"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TLParen, "("),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "address"),
				domain.NewToken(domain.TRParen, ")"),
				domain.NewToken(domain.TKeyword, "values"),
				domain.NewToken(domain.TLParen, "("),
				domain.NewToken(domain.TInt32, int32(123)),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TString, "mike"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TString, "tokyo"),
				domain.NewToken(domain.TRParen, ")"),
			},
		},
		{
			name:  "delete command",
			query: "DELETE FROM foo",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "delete"),
				domain.NewToken(domain.TKeyword, "from"),
				domain.NewToken(domain.TIdentifier, "foo"),
			},
		},
		{
			name:  "delete command with predicate",
			query: "DELETE FROM foo where id = 123",
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "delete"),
				domain.NewToken(domain.TKeyword, "from"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TKeyword, "where"),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TEqual, "="),
				domain.NewToken(domain.TInt32, int32(123)),
			},
		},
		{
			name: "create table command",
			query: `CREATE TABLE foo (
				id int,
				name varchar(255)
			)`,
			tokens: []domain.Token{
				domain.NewToken(domain.TKeyword, "create"),
				domain.NewToken(domain.TKeyword, "table"),
				domain.NewToken(domain.TIdentifier, "foo"),
				domain.NewToken(domain.TLParen, "("),
				domain.NewToken(domain.TIdentifier, "id"),
				domain.NewToken(domain.TKeyword, "int"),
				domain.NewToken(domain.TComma, ","),
				domain.NewToken(domain.TIdentifier, "name"),
				domain.NewToken(domain.TKeyword, "varchar"),
				domain.NewToken(domain.TLParen, "("),
				domain.NewToken(domain.TInt32, int32(255)),
				domain.NewToken(domain.TRParen, ")"),
				domain.NewToken(domain.TRParen, ")"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			lexer := lexer.NewLexer(tt.query)
			tokens, err := lexer.ScanTokens()

			require.NoError(t, err)
			require.Equal(t, tt.tokens, tokens)
		})
	}
}

func TestLexer_Error(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{
			name:  "invalid identifier",
			query: "123abc",
		},
		{
			name:  "invalid character",
			query: "あ",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			lexer := lexer.NewLexer(tt.query)
			_, err := lexer.ScanTokens()

			require.Error(t, err)
		})
	}
}
