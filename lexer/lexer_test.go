package lexer_test

import (
	"testing"

	"github.com/goropikari/simpledbgo/lexer"
	"github.com/stretchr/testify/require"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		tokens []lexer.Token
	}{
		{
			name:  "select query",
			query: "SELECT id, name FROM foo_bar WHERE id = 123 and name = 'Mike\\'s dog'",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "select"),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TIdentifier, "foo_bar"),
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
			name:  "select query with combination upper/lower",
			query: "SeLEcT hoge from foo",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "select"),
				lexer.NewToken(lexer.TIdentifier, "hoge"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
			},
		},
		{
			name:  "insert command",
			query: "INSERT INTO foo (id,name, address) VALUES (-123, 'mike','tokyo')",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "insert"),
				lexer.NewToken(lexer.TKeyword, "into"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "name"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TIdentifier, "address"),
				lexer.NewToken(lexer.TRParen, ")"),
				lexer.NewToken(lexer.TKeyword, "values"),
				lexer.NewToken(lexer.TLParen, "("),
				lexer.NewToken(lexer.TInt32, int32(-123)),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TString, "mike"),
				lexer.NewToken(lexer.TComma, ","),
				lexer.NewToken(lexer.TString, "tokyo"),
				lexer.NewToken(lexer.TRParen, ")"),
			},
		},
		{
			name:  "delete command",
			query: "DELETE FROM foo",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "delete"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
			},
		},
		{
			name:  "delete command with predicate",
			query: "DELETE FROM foo where id = 123",
			tokens: []lexer.Token{
				lexer.NewToken(lexer.TKeyword, "delete"),
				lexer.NewToken(lexer.TKeyword, "from"),
				lexer.NewToken(lexer.TIdentifier, "foo"),
				lexer.NewToken(lexer.TKeyword, "where"),
				lexer.NewToken(lexer.TIdentifier, "id"),
				lexer.NewToken(lexer.TEqual, "="),
				lexer.NewToken(lexer.TInt32, int32(123)),
			},
		},
		{
			name: "create table command",
			query: `CREATE TABLE foo (
				id int,
				name varchar(255)
			)`,
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
		},
		{
			name:  "create index command",
			query: `CREATE INDEX idx_id ON foo (id)`,
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
