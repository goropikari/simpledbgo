package lexer

// TokenType is a type of Token.
type TokenType uint

const (
	// UnknowToken is unknown token type.
	UnknowToken TokenType = iota

	// TKeyword is keyword token type.
	TKeyword

	// TInt32 is int32 token type.
	TInt32

	// TString is string token type.
	TString

	// TIdentifier is identifier token type.
	TIdentifier

	// TStar is star token type.
	TStar

	// TEqual is equal token type.
	TEqual

	// TComma is comma token type.
	TComma

	// TLParen is left parentheses token type.
	TLParen

	// TRParen is right parentheses token type.
	TRParen
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

// Type returns token type.
func (tok Token) Type() TokenType {
	return tok.typ
}

// Value returns token value.
func (tok Token) Value() any {
	return tok.value
}
