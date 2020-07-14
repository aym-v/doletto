package scan

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"unicode"
)

// Scanner holds the state of the scanner.
type Scanner struct {
	r         io.RuneReader // input reader
	peekRunes []rune        // peek runes queue
	num       float64       // number buffer
	buf       bytes.Buffer  // input buffer to hold current scaneme
}

// New creates a new Scanner.
func New(r *io.RuneReader) *Scanner {
	return &Scanner{
		r: *r,
	}
}

// nextRune reads the next rune from the input.
func (s *Scanner) nextRune() rune {
	r, _, err := s.r.ReadRune()
	if err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr)
		}
		r = -1 // EOF rune
	}
	return r
}

// read consumes the peekRunes queue then calls nextRune.
func (s *Scanner) read() rune {
	if len(s.peekRunes) > 0 {
		r := s.peekRunes[0]
		s.peekRunes = s.peekRunes[1:]
		return r
	}
	return s.nextRune()
}

// peek returns but does not consume the next n rune in the input.
func (s *Scanner) peek(n int) rune {
	if len(s.peekRunes) >= n {
		return s.peekRunes[n-1]
	}

	p := s.nextRune()
	s.peekRunes = append(s.peekRunes, p)

	return p
}

// tok resets the peekRunes queue and calls mkToken
func (s *Scanner) tok(typ Type, text string) *Token {
	s.peekRunes = nil
	return mkToken(typ, text)
}

// next returns the next token.
func (s *Scanner) next() *Token {
	for {
		r := s.read()
		switch {
		case r == '@':
			return mkToken(tokAt, "@")
		case isSpace(r):
		case isIdentifierStart(r):
			return s.alphanum(tokIdentifier, r)
		case isDigit(r):
			return s.number(r)
		case isPunctuator(r):
			return s.punctuator(r)
		}
	}
}

// punctuator returns the next punctuator token
func (s *Scanner) punctuator(r rune) *Token {
	switch r {
	case '(':
		return mkToken(tokOpenParen, "(")
	case ')':
		return mkToken(tokCloseParen, ")")
	case '{':
		return mkToken(tokOpenBrace, "{")
	case '}':
		return mkToken(tokCloseBrace, "}")
	case '[':
		return mkToken(tokOpenBracket, "[")
	case ']':
		return mkToken(tokCloseBracket, "]")
	case ',':
		return mkToken(tokComma, ",")
	case ':':
		return mkToken(tokColon, ":")
	case ';':
		return mkToken(tokSemicolon, ";")
	case '~':
		return mkToken(tokTilde, "~")

	case '=':
		// '=' or '=>' or '==' or '==='
		switch s.peek(1) {
		case '=':
			if s.peek(2) == '=' {
				return s.tok(tokEqualsEqualsEquals, "===")
			}
			return s.tok(tokEqualsEquals, "==")
		case '>':
			return s.tok(tokEqualsGreaterThan, "=>")
		}
		return s.tok(tokEquals, "=")

	case '+':
		// '+' or '+=' or '++'
		switch s.peek(1) {
		case '=':
			return s.tok(tokPlusEquals, "+=")
		case '+':
			return s.tok(tokPlusPlus, "++")
		}
		return s.tok(tokPlus, "+")

	case '-':
		// '-' or '-=' or '--'
		switch s.peek(1) {
		case '=':
			return s.tok(tokMinusEquals, "-=")
		case '-':
			return s.tok(tokMinusMinus, "--")
		}
		return s.tok(tokMinus, "-")

	case '*':
		// '*' or '*=' or '**' or '**='
		switch s.peek(1) {
		case '=':
			return s.tok(tokAsteriskEquals, "*=")
		case '*':
			if s.peek(2) == '=' {
				return s.tok(tokAsteriskAsteriskEquals, "**=")
			}
			return s.tok(tokAsteriskAsterisk, "**")
		}
		return s.tok(tokAsterisk, "*")

	case '/':
		// '/' or '/=' or '//' or '/* ... */'
		switch s.peek(1) {
		case '=':
			return s.tok(tokSlashEquals, "/=")
		case '/':
			// Single line comment
		case '*':
			// Multi line comment
		}
		return s.tok(tokSlash, "/")

	case '>':
		// '>' or '>>' or '>>>' or '>=' or '>>=' or '>>>='
		switch s.peek(1) {
		case '>':
			switch s.peek(2) {
			case '>':
				if s.peek(3) == '=' {
					return s.tok(tokGreaterThanGreaterThanGreaterThanEquals, ">>>=")
				}
				return s.tok(tokGreaterThanGreaterThanGreaterThan, ">>>")
			case '=':
				return s.tok(tokGreaterThanGreaterThanEquals, ">>=")
			}
			return s.tok(tokGreaterThanGreaterThan, ">>")
		case '=':
			return s.tok(tokGreaterThanEquals, ">=")
		}
		return s.tok(tokGreaterThan, ">")

	case '<':
		// '<' or '<<' or '<=' or '<<='
		switch s.peek(1) {
		case '<':
			if s.peek(2) == '=' {
				return s.tok(tokLessThanLessThanEquals, "<<=")
			}
			return s.tok(tokLessThanLessThan, "<<")
		case '=':
			return s.tok(tokLessThanEquals, "<=")
		}
		return s.tok(tokLessThan, "<")

	case '!':
		// '!' or '!=' or '!=='
		if s.peek(1) == '=' {
			if s.peek(2) == '=' {
				return s.tok(tokExclamationEqualsEquals, "!==")
			}
			return s.tok(tokExclamationEquals, "!=")
		}
		return s.tok(tokExclamation, "!")

	case '^':
		// '^' or '^='
		if s.peek(1) == '=' {
			return s.tok(tokCaretEquals, "^=")
		}
		return s.tok(tokCaret, "^")

	case '|':
		// '|' or '|=' or '||' or '||='
		switch s.peek(1) {
		case '=':
			return s.tok(tokBarEquals, "|=")
		case '|':
			if s.peek(2) == '=' {
				return s.tok(tokBarBarEquals, "||=")
			}
			return s.tok(tokBarBar, "||")
		}
		return s.tok(tokBar, "|")

	case '&':
		// '&' or '&=' or '&&' or '&&='
		switch s.peek(1) {
		case '=':
			return s.tok(tokAmpersandEquals, "&=")
		case '&':
			if s.peek(2) == '=' {
				return s.tok(tokAmpersandAmpersandEquals, "&&=")
			}
			return s.tok(tokAmpersandAmpersand, "&&")
		}
		return s.tok(tokAmpersand, "&")

	case '%':
		// '%' or '%='
		if s.peek(1) == '=' {
			return s.tok(tokPercentEquals, "%=")
		}
		return s.tok(tokPercent, "%")

	case '?':
		// '?' or '?.' or '??' or '??='
		switch s.peek(1) {
		case '?':
			if s.peek(2) == '=' {
				return s.tok(tokQuestionQuestionEquals, "??=")
			}
			return s.tok(tokQuestionQuestion, "??")
		case '.':
			// differentiate optional chaining punctuators (?.id) from conditional operators (? :)
			if !unicode.IsDigit(s.peek(2)) {
				return s.tok(tokQuestionDot, "?.")
			}
		}
		return s.tok(tokQuestion, "?")
	default:
		return s.tok(tokSyntaxError, "")

	}
}

// accum appends the current rune to the buffer until
// the valid function returns false
func (s *Scanner) accum(r rune, valid func(rune) bool) {
	s.buf.Reset()
	for {
		s.buf.WriteRune(r)
		r = s.read()
		if r == -1 {
			return
		}
		if !valid(r) {
			return
		}
	}
}

// alphanum creates a keyword or identifier token using the buffer.
func (s *Scanner) alphanum(typ Type, r rune) *Token {
	s.accum(r, isIdentifierContinue)
	return mkToken(typ, s.buf.String())
}

func (s *Scanner) number(r rune) *Token {
	base := 10.0
	isLegacyOctal := false
	// isInvalidLegacyOctal := false
	num := 0.0

	switch r {
	case '.':
		// '.' or '...'
		if s.peek(1) == '.' && s.peek(2) == '.' {
			return s.tok(tokDotDotDot, "...")
		}
		return s.tok(tokDot, ".")

	case '0':
		// binary, octal or hexadecimal literal
		switch s.peek(1) {
		case 'b', 'B':
			base = 2

		case 'o', 'O':
			base = 8

		case 'x', 'X':
			base = 16

		case '0', '1', '2', '3', '4', '5', '6', '7':
			base = 8
			isLegacyOctal = true
		}
	default:
		num = float64(r - '0')
	}

intLiteral:
	for {
		next := s.read()
		switch next {
		case '_':
			// numeric separator
			b := s.buf.String()

			if b[len(b)-1] == '_' {
				s.syntaxError()
			}

		case '0', '1':
			num = num*base + float64(next-'0')

		case '2', '3', '4', '5', '6', '7':
			if base == 2 {
				s.syntaxError()
			}
			num = num*base + float64(next-'0')

		case '8', '9':
			if isLegacyOctal {
				// isInvalidLegacyOctal = true
			} else if base < 10 {
				s.syntaxError()
			}
			num = num*base + float64(next-'0')
		case 'A', 'B', 'C', 'D', 'E', 'F':
			if base != 16 {
				s.syntaxError()
			}
			num = num*base + float64(next+10-'A')

		case 'a', 'b', 'c', 'd', 'e', 'f':
			if base != 16 {
				s.syntaxError()
			}
			num = num*base + float64(next+10-'a')

		default:
			// The first digit must exist
			if len(s.buf.String()) == 1 {
				s.syntaxError()
			}

			break intLiteral
		}
	}

	return mkNumericLiteral(num)
}

// func (s *Scanner) intLiteral(r rune, base float64) *Token {

// }

// alphanum creates a numeric literal token using the buffer.
// func (s *Scanner) number(r rune) *Token {

// }

// Panic represents a scanner panic
type Panic struct{}

func (s *Scanner) syntaxError() {
	log.Printf("Syntax Error")
	panic(Panic{})
}

// isAlphaNumeric reports whether r is a letter, digit, or underscore.
func isAlphanum(r rune) bool {
	return r == '_' || r == '$' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isDigit reports whether r is a digit or dot.
func isDigit(r rune) bool {
	switch r {
	case '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	default:
		return false
	}
}

// isPunctuator reports whether r is a punctuator
func isPunctuator(r rune) bool {
	switch r {
	case '{', '}', '(', ')', '[', ']', '.', ';', ',', '<', '>', '=', '!', '+', '-', '*', '%', '&', '|', '^', '~', '?', ':', '/':
		return true
	default:
		return false
	}
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
