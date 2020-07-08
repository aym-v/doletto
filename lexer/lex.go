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
	r         io.RuneReader // input reader
	peekRunes []rune        // peek runes queue
	buf       bytes.Buffer  // input buffer to hold current lexeme
}

// New creates a new Scanner.
func New(r *io.RuneReader) *Scanner {
	return &Scanner{
		r: *r,
	}
}

// nextRune reads the next rune from the input.
func (l *Scanner) nextRune() rune {
	r, _, err := l.r.ReadRune()
	if err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr)
		}
		r = -1 // EOF rune
	}
	return r
}

// read consumes the peekRunes queue then calls nextRune.
func (l *Scanner) read() rune {
	if len(l.peekRunes) > 0 {
		r := l.peekRunes[0]
		l.peekRunes = l.peekRunes[1:]
		return r
	}
	return l.nextRune()
}

// peek returns but does not consume the next n rune in the input.
func (l *Scanner) peek(n int) rune {
	if len(l.peekRunes) >= n {
		return l.peekRunes[n-1]
	}

	p := l.nextRune()
	l.peekRunes = append(l.peekRunes, p)

	return p
}

// resetPeek resets the peekRunes queue and calls mkToken
func (l *Scanner) mkPeekTok(typ Type, text string) *Token {
	l.peekRunes = nil
	return mkToken(typ, text)
}

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
		case isNumber(r):
			// return l.number(r)

		case r == '=':
			// '=' or '=>' or '==' or '==='
			switch l.peek(1) {
			case '=':
				if l.peek(2) == '=' {
					return l.mkPeekTok(tokEqualsEqualsEquals, "===")
				}
				return l.mkPeekTok(tokEqualsEquals, "==")
			case '>':
				return l.mkPeekTok(tokEqualsGreaterThan, "=>")
			}
			return l.mkPeekTok(tokEquals, "=")

		case r == '+':
			// '+' or '+=' or '++'
			switch l.peek(1) {
			case '=':
				return l.mkPeekTok(tokPlusEquals, "+=")
			case '+':
				return l.mkPeekTok(tokPlusPlus, "++")
			}
			return l.mkPeekTok(tokPlus, "+")

		case r == '-':
			// '-' or '-=' or '--'
			switch l.peek(1) {
			case '=':
				return l.mkPeekTok(tokMinusEquals, "-=")
			case '-':
				return l.mkPeekTok(tokMinusMinus, "--")
			}
			return l.mkPeekTok(tokMinus, "-")

		case r == '*':
			// '*' or '*=' or '**' or '**='
			switch l.peek(1) {
			case '=':
				return l.mkPeekTok(tokAsteriskEquals, "*=")
			case '*':
				if l.peek(2) == '=' {
					return l.mkPeekTok(tokAsteriskAsteriskEquals, "**=")
				}
				return l.mkPeekTok(tokAsteriskAsterisk, "**")
			}
			return l.mkPeekTok(tokAsterisk, "*")

		case r == '/':
			// '/' or '/=' or '//' or '/* ... */'
			switch l.peek(1) {
			case '=':
				return l.mkPeekTok(tokSlashEquals, "/=")
			case '/':
				// Single line comment
			case '*':
				// Multi line comment
			}
			return l.mkPeekTok(tokSlash, "/")

		case r == '>':
			// '>' or '>>' or '>>>' or '>=' or '>>=' or '>>>='
			switch l.peek(1) {
			case '>':
				switch l.peek(2) {
				case '>':
					if l.peek(3) == '=' {
						return l.mkPeekTok(tokGreaterThanGreaterThanGreaterThanEquals, ">>>=")
					}
					return l.mkPeekTok(tokGreaterThanGreaterThanGreaterThan, ">>>")
				case '=':
					return l.mkPeekTok(tokGreaterThanGreaterThanEquals, ">>=")
				}
				return l.mkPeekTok(tokGreaterThanGreaterThan, ">>")
			case '=':
				return l.mkPeekTok(tokGreaterThanEquals, ">=")
			}
			return l.mkPeekTok(tokGreaterThan, ">")

		case r == '<':
			// '<' or '<<' or '<=' or '<<='
			switch l.peek(1) {
			case '<':
				if l.peek(2) == '=' {
					return l.mkPeekTok(tokLessThanLessThanEquals, "<<=")
				}
				return l.mkPeekTok(tokLessThanLessThan, "<<")
			case '=':
				return l.mkPeekTok(tokLessThanEquals, "<=")
			}
			return l.mkPeekTok(tokLessThan, "<")

		case r == '!':
			// '!' or '!=' or '!=='
			if l.peek(1) == '=' {
				if l.peek(2) == '=' {
					return l.mkPeekTok(tokExclamationEqualsEquals, "!==")
				}
				return l.mkPeekTok(tokExclamationEquals, "!=")
			}
			return l.mkPeekTok(tokExclamation, "!")

		case r == '^':
			// '^' or '^='
			if l.peek(1) == '=' {
				return l.mkPeekTok(tokCaretEquals, "^=")
			}
			return l.mkPeekTok(tokCaret, "^")

		case r == '|':
			// '|' or '|=' or '||' or '||='
			switch l.peek(1) {
			case '=':
				return l.mkPeekTok(tokBarEquals, "|=")
			case '|':
				if l.peek(2) == '=' {
					return l.mkPeekTok(tokBarBarEquals, "||=")
				}
				return l.mkPeekTok(tokBarBar, "||")
			}
			return l.mkPeekTok(tokBar, "|")
		}
	}
}

// accum appends the current rune to the buffer until
// the valid function returns false
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

// alphanum creates a keyword or identifier token using the buffer.
func (l *Scanner) alphanum(typ Type, r rune) *Token {
	l.accum(r, isAlphanum)
	return mkToken(typ, l.buf.String())
}

// alphanum creates a numeric literal token using the buffer.
// func (l *Scanner) number(r rune) *Token {

// }

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphanum(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isAlphaNumeric reports whether r is a numeric literal.
func isNumber(r rune) bool {
	return '0' <= r && r <= '9'
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
