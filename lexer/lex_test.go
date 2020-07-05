package lexer

import (
	"bufio"
	"strings"
	"testing"
)

func TestToken(t *testing.T) {
	sample := "({[]})"

	tests := []struct {
		expTyp  T
		expText string
	}{
		{TOpenParen, "("},
		{TOpenBrace, "{"},
		{TOpenBracket, "["},
		{TCloseBracket, "]"},
		{TCloseBrace, "}"},
		{TCloseParen, ")"},
	}

	rd := bufio.NewReader(strings.NewReader(sample))

	l := NewLexer(rd)

	for _, c := range tests {
		tok := l.Next()

		if tok.typ != c.expTyp {
			t.Fatalf("token type is wrong. expected=%q, got=%q", c.expTyp, tok.typ)
		}

		if tok.text != c.expText {
			t.Fatalf("token text is wrong. expected=%q, got=%q", c.expText, tok.text)
		}
	}
}
