package lexer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"unicode"
)

// Lexer represents an Ecmascript lexer
type Lexer struct {
	rd   io.RuneReader // input reader
	ch   rune          // current Unicode char
	peek rune          // next Unicode char
	buf  bytes.Buffer
}

// NewLexer creates a new lexer that will read from the RuneReader
func NewLexer(rd io.RuneReader) *Lexer {
	return &Lexer{
		rd: rd,
	}
}

func (l *Lexer) read() rune {
	r, _, err := l.rd.ReadRune()
	if err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr)
		}
		r = -1 // EOF rune
	}
	l.ch = r
	return r
}

// Next returns the next token
func (l *Lexer) Next() *Token {
	for {
		r := l.read()
		switch {
		case isSpace(r):
		case r == '(':
			return mkToken(TOpenParen, "(")
		case r == ')':
			return mkToken(TCloseParen, ")")
		case r == '{':
			return mkToken(TOpenBrace, "{")
		case r == '}':
			return mkToken(TCloseBrace, "}")
		case r == '[':
			return mkToken(TOpenBracket, "[")
		case r == ']':
			return mkToken(TCloseBracket, "]")
		case r == ',':
			return mkToken(TComma, ",")
		case r == ':':
			return mkToken(TColon, ":")
		case r == ';':
			return mkToken(TSemicolon, ";")
		case r == '@':
			return mkToken(TAt, "@")
		case r == '~':
			return mkToken(TTilde, "~")
		}
	}
}

// isSpace checks whether r is a space as defined
// in the Unicode standard or the ECMAScript specification
func isSpace(r rune) bool {
	switch {
	case r == 0x85:
		return false
	case
		unicode.IsSpace(r),
		r == '\uFEFF': // zero width non-breaking space
		return true

	default:
		return false
	}
}
