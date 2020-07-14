package scan

import (
	"io"
	"strings"
	"testing"

	"github.com/valaymerick/doletto/test"
)

func TestPeek(t *testing.T) {
	in := `
	( ) { } [ ] , : ; @ ~ === == => = + += ++ - -= -- * *= ** **= / /= > >> >>> >= >>= >>>=
	! != !== < << <= <<= ^ ^= | |= || ||= & &= && &&= % %= ? ?. ?? ??= ... .
	`

	out := []struct {
		expTyp  Type
		expText string
	}{
		{tokOpenParen, "("},
		{tokCloseParen, ")"},
		{tokOpenBrace, "{"},
		{tokCloseBrace, "}"},
		{tokOpenBracket, "["},
		{tokCloseBracket, "]"},
		{tokComma, ","},
		{tokColon, ":"},
		{tokSemicolon, ";"},
		{tokAt, "@"},
		{tokTilde, "~"},

		// '=' or '=>' or '==' or '==='
		{tokEqualsEqualsEquals, "==="},
		{tokEqualsEquals, "=="},
		{tokEqualsGreaterThan, "=>"},
		{tokEquals, "="},

		// '+' or '+=' or '++'
		{tokPlus, "+"},
		{tokPlusEquals, "+="},
		{tokPlusPlus, "++"},
		{tokMinus, "-"},
		{tokMinusEquals, "-="},
		{tokMinusMinus, "--"},

		// '*' or '*=' or '**' or '**='
		{tokAsterisk, "*"},
		{tokAsteriskEquals, "*="},
		{tokAsteriskAsterisk, "**"},
		{tokAsteriskAsteriskEquals, "**="},

		// '/' or '/='
		{tokSlash, "/"},
		{tokSlashEquals, "/="},

		// '>' or '>>' or '>>>' or '>=' or '>>=' or '>>>='
		{tokGreaterThan, ">"},
		{tokGreaterThanGreaterThan, ">>"},
		{tokGreaterThanGreaterThanGreaterThan, ">>>"},
		{tokGreaterThanEquals, ">="},
		{tokGreaterThanGreaterThanEquals, ">>="},
		{tokGreaterThanGreaterThanGreaterThanEquals, ">>>="},

		// '!' or '!=' or '!=='
		{tokExclamation, "!"},
		{tokExclamationEquals, "!="},
		{tokExclamationEqualsEquals, "!=="},

		// '<' or '<<' or '<=' or '<<='
		{tokLessThan, "<"},
		{tokLessThanLessThan, "<<"},
		{tokLessThanEquals, "<="},
		{tokLessThanLessThanEquals, "<<="},

		// '^' or '^='
		{tokCaret, "^"},
		{tokCaretEquals, "^="},

		// '|' or '|=' or '||' or '||='
		{tokBar, "|"},
		{tokBarEquals, "|="},
		{tokBarBar, "||"},
		{tokBarBarEquals, "||="},

		// '&' or '&=' or '&&' or '&&='
		{tokAmpersand, "&"},
		{tokAmpersandEquals, "&="},
		{tokAmpersandAmpersand, "&&"},
		{tokAmpersandAmpersandEquals, "&&="},

		// '%' or '%='
		{tokPercent, "%"},
		{tokPercentEquals, "%="},

		// '?' or '?.' or '??' or '??='
		{tokQuestion, "?"},
		{tokQuestionDot, "?."},
		{tokQuestionQuestion, "??"},
		{tokQuestionQuestionEquals, "??="},

		// '...' or '.'
		{tokDotDotDot, "..."},
		{tokDot, "."},
	}

	r := io.RuneReader(strings.NewReader(in))

	l := New(&r)

	for _, c := range out {
		tok := l.next()

		test.AssertEqual(t, tok.typ, c.expTyp)
		test.AssertEqual(t, tok.text, c.expText)
	}
}

// expectNumber tests whether a numeric literal token
// holds the right value.
func expectNumber(t *testing.T, in string, expected float64) {
	t.Run(in, func(t *testing.T) {
		r := io.RuneReader(strings.NewReader(in))

		l := New(&r)
		out := l.next()

		test.AssertEqual(t, tokNumericLiteral, out.typ)
		test.AssertEqual(t, out.num, expected)
	})
}

func TestNumericLiteral(t *testing.T) {
	expectNumber(t, "0", 0.0)
	expectNumber(t, "000", 0.0)
	expectNumber(t, "010", 8.0)
	expectNumber(t, "123", 123.0)
	expectNumber(t, "987", 987.0)
	expectNumber(t, "0000", 0.0)
	expectNumber(t, "0123", 83.0)
	expectNumber(t, "0123.4567", 83.0)
	expectNumber(t, "0987", 987.0)
	expectNumber(t, "0987.6543", 987.6543)
	// expectNumber(t, "01289", 1289.0)
	// expectNumber(t, "01289.345", 1289.0)
	// expectNumber(t, "999999999", 999999999.0)
	// expectNumber(t, "9999999999", 9999999999.0)
	// expectNumber(t, "99999999999", 99999999999.0)
	// expectNumber(t, "123456789123456789", 123456789123456780.0)
	// expectNumber(t, "123456789123456789"+strings.Repeat("0", 128), 1.2345678912345679e+145)

	// expectNumber(t, "0b00101", 5.0)
	// expectNumber(t, "0B00101", 5.0)
	// expectNumber(t, "0b1011101011101011101011101011101011101", 100352251741.0)
	// expectNumber(t, "0B1011101011101011101011101011101011101", 100352251741.0)

	// expectNumber(t, "0o12345", 5349.0)
	// expectNumber(t, "0O12345", 5349.0)
	// expectNumber(t, "0o1234567654321", 89755965649.0)
	// expectNumber(t, "0O1234567654321", 89755965649.0)

	// expectNumber(t, "0x12345678", float64(0x12345678))
	// expectNumber(t, "0xFEDCBA987", float64(0xFEDCBA987))
	// expectNumber(t, "0x000012345678", float64(0x12345678))
	// expectNumber(t, "0x123456781234", float64(0x123456781234))

	// expectNumber(t, "123.", 123.0)
	// expectNumber(t, ".0123", 0.0123)
	// expectNumber(t, "0.0123", 0.0123)
	// expectNumber(t, "2.2250738585072014e-308", 2.2250738585072014e-308)
	// expectNumber(t, "1.7976931348623157e+308", 1.7976931348623157e+308)

	// // Underflow
	// expectNumber(t, "4.9406564584124654417656879286822e-324", 5e-324)
	// expectNumber(t, "5e-324", 5e-324)
	// expectNumber(t, "1e-325", 0.0)

	// // Overflow
	// expectNumber(t, "1.797693134862315708145274237317e+308", 1.7976931348623157e+308)
	// expectNumber(t, "1.797693134862315808e+308", math.Inf(1))
	// expectNumber(t, "1e+309", math.Inf(1))

	// // int32
	// expectNumber(t, "0x7fff_ffff", 2147483647.0)
	// expectNumber(t, "0x8000_0000", 2147483648.0)
	// expectNumber(t, "0x8000_0001", 2147483649.0)

	// // uint32
	// expectNumber(t, "0xffff_ffff", 4294967295.0)
	// expectNumber(t, "0x1_0000_0000", 4294967296.0)
	// expectNumber(t, "0x1_0000_0001", 4294967297.0)

	// // int64
	// expectNumber(t, "0x7fff_ffff_ffff_fdff", 9223372036854774784)
	// expectNumber(t, "0x8000_0000_0000_0000", 9.223372036854776e+18)
	// expectNumber(t, "0x8000_0000_0000_3000", 9.223372036854788e+18)

	// // uint64
	// expectNumber(t, "0xffff_ffff_ffff_fbff", 1.844674407370955e+19)
	// expectNumber(t, "0x1_0000_0000_0000_0000", 1.8446744073709552e+19)
	// expectNumber(t, "0x1_0000_0000_0000_1000", 1.8446744073709556e+19)

	// expectNumber(t, "1.", 1.0)
	// expectNumber(t, ".1", 0.1)
	// expectNumber(t, "1.1", 1.1)
	// expectNumber(t, "1e1", 10.0)
	// expectNumber(t, "1e+1", 10.0)
	// expectNumber(t, "1e-1", 0.1)
	// expectNumber(t, ".1e1", 1.0)
	// expectNumber(t, ".1e+1", 1.0)
	// expectNumber(t, ".1e-1", 0.01)
	// expectNumber(t, "1.e1", 10.0)
	// expectNumber(t, "1.e+1", 10.0)
	// expectNumber(t, "1.e-1", 0.1)
	// expectNumber(t, "1.1e1", 11.0)
	// expectNumber(t, "1.1e+1", 11.0)
	// expectNumber(t, "1.1e-1", 0.11)

	// expectNumber(t, "1_2_3", 123)
	// expectNumber(t, ".1_2", 0.12)
	// expectNumber(t, "1_2.3_4", 12.34)
	// expectNumber(t, "1e2_3", 1e23)
	// expectNumber(t, "1_2e3_4", 12e34)
	// expectNumber(t, "1_2.3_4e5_6", 12.34e56)
	// expectNumber(t, "0b1_0", 2)
	// expectNumber(t, "0B1_0", 2)
	// expectNumber(t, "0o1_2", 10)
	// expectNumber(t, "0O1_2", 10)
	// expectNumber(t, "0x1_2", 0x12)
	// expectNumber(t, "0X1_2", 0x12)
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

		test.AssertEqual(t, c.out, r)
	}
}
