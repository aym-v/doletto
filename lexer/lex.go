package lexer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"unicode"
)

// Scanner holds the state of the scanner.
type Scanner struct {
	r   io.RuneReader // input reader
	buf bytes.Buffer  // input buffer to hold current lexeme
}

// New creates a new Scanner.
func New(r *io.RuneReader) *Scanner {
	return &Scanner{
		r: *r,
	}
}

// read reads the next rune from the input.
func (l *Scanner) read() rune {
	r, _, err := l.r.ReadRune()
	if err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr)
		}
		r = -1 // EOF rune
	}
	return r
}

// peek returns but does not consume the next rune in the input.
// func (l *Scanner) peek(n int) rune {

// }

// next returns the next token.
func (l *Scanner) next() *Token {
	for {
		r := l.read()
		switch {
		case isSpace(r):
		case r == '(':
			return mkToken(tokOpenParen, "(")
		case r == ')':
			return mkToken(tokCloseParen, ")")
		case r == '{':
			return mkToken(tokOpenBrace, "{")
		case r == '}':
			return mkToken(tokCloseBrace, "}")
		case r == '[':
			return mkToken(tokOpenBracket, "[")
		case r == ']':
			return mkToken(tokCloseBracket, "]")
		case r == ',':
			return mkToken(tokComma, ",")
		case r == ':':
			return mkToken(tokColon, ":")
		case r == ';':
			return mkToken(tokSemicolon, ";")
		case r == '@':
			return mkToken(tokAt, "@")
		case r == '~':
			return mkToken(tokTilde, "~")
		case isAlphanum(r):
			return l.alphanum(tokIdentifier, r)
		}
	}
}

func (l *Scanner) accum(r rune, valid func(rune) bool) {
	l.buf.Reset()
	for {
		l.buf.WriteRune(r)
		r = l.read()
		if r == -1 {
			return
		}
		if !valid(r) {
			return
		}
	}
}

func (l *Scanner) alphanum(typ Type, r rune) *Token {
	l.accum(r, isAlphanum)
	return mkToken(typ, l.buf.String())
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphanum(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isSpace checks whether r is a space as defined
// in the Unicode standard or the ECMAScript specification.
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
