package lexer

import (
	"bytes"
	"fmt"
	"io"
	"os"
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
		switch r {
		case '(':
			return mkToken(TOpenParen, "(")
		case ')':
			return mkToken(TCloseParen, ")")
		case '{':
			return mkToken(TOpenBrace, "{")
		case '}':
			return mkToken(TCloseBrace, "}")
		case '[':
			return mkToken(TOpenBracket, "[")
		case ']':
			return mkToken(TCloseBracket, "]")
		case ',':
			return mkToken(TComma, ",")
		case ':':
			return mkToken(TColon, ":")
		case ';':
			return mkToken(TSemicolon, ";")
		case '@':
			return mkToken(TAt, "@")
		case '~':
			return mkToken(TTilde, "~")
		}
	}
}
