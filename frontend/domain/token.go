package domain

// TokenType is a type of Token.
type TokenType uint

const (
	// UnknowToken is unknown token type.
	UnknowToken TokenType = iota

	// Keyword is keyword token type.
	Keyword

	// Int32 is int32 token type.
	Int32

	// String is string token type.
	String

	// Identifier is identifier token type.
	Identifier

	// Star is star token type.
	Star

	// Equal is equal token type.
	Equal

	// Comma is comma token type.
	Comma
)

// Token is model of token.
type Token struct {
	typ   TokenType
	value any
}

// NewToken constructs a token.
func NewToken(typ TokenType, value any) Token {
	return Token{
		typ:   typ,
		value: value,
	}
}
