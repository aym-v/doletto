package lexer

import (
	"io"
	"strings"
	"testing"
)

func TestBrackets(t *testing.T) {
	sample := "({[]})"

	tests := []struct {
		expTyp  Type
		expText string
	}{
		{tokOpenParen, "("},
		{tokOpenBrace, "{"},
		{tokOpenBracket, "["},
		{tokCloseBracket, "]"},
		{tokCloseBrace, "}"},
		{tokCloseParen, ")"},
	}

	r := io.RuneReader(strings.NewReader(sample))

	l := New(&r)

	for _, c := range tests {
		tok := l.next()

		if tok.typ != c.expTyp {
			t.Fatalf("token type is wrong. expected=%q, got=%q", c.expTyp, tok.typ)
		}

		if tok.text != c.expText {
			t.Fatalf("token text is wrong. expected=%q, got=%q", c.expText, tok.text)
		}
	}
}

func TestKeyword(t *testing.T) {
	sample := "const foo = 1"

	tests := []struct {
		expTyp  Type
		expText string
	}{
		{tokConst, "const"},
		{tokIdentifier, "foo"},
		{tokEquals, "="},
		// {tokNumericLiteral, "1"},
	}

	r := io.RuneReader(strings.NewReader(sample))

	l := New(&r)

	for _, c := range tests {
		tok := l.next()

		if tok.typ != c.expTyp {
			t.Fatalf("token type is wrong. expected=%q, got=%q", c.expTyp, tok.typ)
		}

		if tok.text != c.expText {
			t.Fatalf("token text is wrong. expected=%q, got=%q", c.expText, tok.text)
		}
	}
}

func TestPeek(t *testing.T) {
	sample := "=== == => = + += ++ - -= -- * *= ** **="

	tests := []struct {
		expTyp  Type
		expText string
	}{
		{tokEqualsEqualsEquals, "==="},
		{tokEqualsEquals, "=="},
		{tokEqualsGreaterThan, "=>"},
		{tokEquals, "="},
		{tokPlus, "+"},
		{tokPlusEquals, "+="},
		{tokPlusPlus, "++"},
		{tokMinus, "-"},
		{tokMinusEquals, "-="},
		{tokMinusMinus, "--"},
		{tokAsterisk, "*"},
		{tokAsteriskEquals, "*="},
		{tokAsteriskAsterisk, "**"},
		{tokAsteriskAsteriskEquals, "**="},
	}

	r := io.RuneReader(strings.NewReader(sample))

	l := New(&r)

	for _, c := range tests {
		tok := l.next()

		if tok.typ != c.expTyp {
			t.Fatalf("token type is wrong. expected=%q, got=%q", c.expTyp, tok.typ)
		}

		if tok.text != c.expText {
			t.Fatalf("token text is wrong. expected=%q, got=%q", c.expText, tok.text)
		}
	}
}

func TestIsSpace(t *testing.T) {
	tests := []struct {
		in  rune
		out bool
	}{
		{'\u0009', true},  // tab
		{'\u000B', true},  // vertical tab
		{'\u000C', true},  // form feed
		{'\u0020', true},  // space
		{'\u00A0', true},  // no-break space
		{'\u1680', true},  // ogham space mark
		{'\u2000', true},  // en quad
		{'\u2001', true},  // em quad
		{'\u2002', true},  // en space
		{'\u2003', true},  // em space
		{'\u2004', true},  // three-per-em space
		{'\u2005', true},  // four-per-em space
		{'\u2006', true},  // six-per-em space
		{'\u2007', true},  // figure space
		{'\u2008', true},  // punctuation space
		{'\u2009', true},  // thin space
		{'\u200A', true},  // hair space
		{'\u202F', true},  // narrow no-break space
		{'\u205F', true},  // medium mathematical space
		{'\u3000', true},  // ideographic space
		{'\uFEFF', true},  // zero width non-breaking space
		{'\u0085', false}, // next line
	}

	for _, c := range tests {
		r := isSpace(c.in)

		if r != c.out {
			t.Fatalf("isSpace output is wrong. expected=%t, got=%t", c.out, r)
		}
	}
}
