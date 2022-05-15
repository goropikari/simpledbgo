package lexer

import (
	"errors"
	"io"
	"strconv"
	"strings"
)

var keywords = []string{
	"select", "from", "where", "and",
	"insert", "into", "values", "delete", "update", "set",
	"create", "table", "int", "varchar", "view", "as", "index", "on",
}

// Lexer is a model of lexer.
type Lexer struct {
	reader   *strings.Reader
	keywords map[string]bool
}

// NewLexer constructs a lexer.
func NewLexer(query string) *Lexer {
	keywordMap := make(map[string]bool)
	for _, key := range keywords {
		keywordMap[key] = true
	}

	return &Lexer{
		reader:   strings.NewReader(query),
		keywords: keywordMap,
	}
}

// ScanTokens scans token from the query.
func (lex *Lexer) ScanTokens() ([]Token, error) {
	tokens := make([]Token, 0)

	for {
		token, err := lex.scan()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (lex *Lexer) scan() (Token, error) {
	err := lex.skipWhitespace()
	if err != nil {
		return Token{}, err
	}

	c, err := lex.readByte()
	if err != nil {
		return Token{}, err
	}

	switch c {
	case '*':
		return NewToken(TStar, "*"), nil
	case '=':
		return NewToken(TEqual, "="), nil
	case ',':
		return NewToken(TComma, ","), nil
	case '(':
		return NewToken(TLParen, "("), nil
	case ')':
		return NewToken(TRParen, ")"), nil
	}

	err = lex.unreadByte()
	if err != nil {
		return Token{}, err
	}

	switch {
	case isNumber(c) || c == '-':
		numStr, err := lex.scanInteger()
		if err != nil {
			return Token{}, err
		}

		base := 10
		precision := 32
		num, err := strconv.ParseInt(numStr, base, precision)
		if err != nil {
			return Token{}, err
		}

		return NewToken(TInt32, int32(num)), nil
	case isAlpha(c):
		id, err := lex.scanIdentifier()
		if err != nil {
			return Token{}, err
		}

		if lex.isKeyword(id) {
			return NewToken(TKeyword, id), nil
		} else {
			return NewToken(TIdentifier, id), nil
		}
	case c == '\'':
		s, err := lex.scanString()
		if err != nil {
			return Token{}, err
		}

		return NewToken(TString, s), nil
	}

	return Token{}, errors.New("error at scan")
}

func (lex *Lexer) scanInteger() (string, error) {
	b := make([]byte, 0)
	c, _ := lex.readByte()
	if c == '-' {
		b = append(b, c)
	} else {
		err := lex.unreadByte()
		if err != nil {
			return "", err
		}
	}

	for {
		c, err := lex.readByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return "", err
		}

		if isAlpha(c) {
			return "", errors.New("not number")
		}
		if !isNumber(c) {
			err := lex.unreadByte()
			if err != nil {
				return "", errors.New("not number")
			}

			break
		}

		b = append(b, c)
	}

	return string(b), nil
}

func (lex *Lexer) scanIdentifier() (string, error) {
	b := make([]byte, 0)
	for {
		c, err := lex.readByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return "", err
		}

		if !isAlphaNumeric(c) {
			err = lex.unreadByte()
			if err != nil {
				return "", err
			}

			break
		}

		b = append(b, c)
	}

	return strings.ToLower(string(b)), nil
}

func (lex *Lexer) scanString() (string, error) {
	_, err := lex.readByte()
	if err != nil {
		return "", err
	}

	b := make([]byte, 0)

	esc := false
	for {
		c, err := lex.readByte()
		if err != nil {
			return "", err
		}

		switch {
		case c == '\\':
			if esc {
				esc = false
			} else {
				esc = true

				continue
			}
		case c == '\'' && !esc:
			goto Ret
		default:
			esc = false
		}

		b = append(b, c)
	}

Ret:
	return string(b), nil
}

func (lex *Lexer) isKeyword(s string) bool {
	return lex.keywords[s]
}

func (lex *Lexer) skipWhitespace() error {
	for {
		c, err := lex.readByte()
		if err != nil {
			return err
		}

		if !isWhitespace(c) {
			err = lex.unreadByte()
			if err != nil {
				return err
			}

			break
		}
	}

	return nil
}

func (lex *Lexer) readByte() (byte, error) {
	return lex.reader.ReadByte()
}

func (lex *Lexer) unreadByte() error {
	return lex.reader.UnreadByte()
}

func isNumber(c byte) bool {
	return '0' <= c && c <= '9'
}

func isAlpha(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isNumber(c)
}

func isWhitespace(c byte) bool {
	switch c {
	case ' ', '\r', '\n', '\t':
		return true
	}

	return false
}
