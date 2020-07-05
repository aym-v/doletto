package lexer

import (
	"log"
	"math/big"
)

// T represents a generic token
type T int

const (
	TEndOfFile T = iota
	TSyntaxError

	THashbang

	TNoSubstitutionTemplateLiteral
	TNumericLiteral
	TStringLiteral

	// Punctuation
	TAmpersand
	TAmpersandAmpersand
	TAsterisk
	TAsteriskAsterisk
	TAt
	TBar
	TBarBar
	TCaret
	TCloseBrace
	TCloseBracket
	TCloseParen
	TColon
	TComma
	TDot
	TDotDotDot
	TEqualsEquals
	TEqualsEqualsEquals
	TEqualsGreaterThan
	TExclamation
	TExclamationEquals
	TExclamationEqualsEquals
	TGreaterThan
	TGreaterThanEquals
	TGreaterThanGreaterThan
	TGreaterThanGreaterThanGreaterThan
	TLessThan
	TLessThanEquals
	TLessThanLessThan
	TMinus
	TMinusMinus
	TOpenBrace
	TOpenBracket
	TOpenParen
	TPercent
	TPlus
	TPlusPlus
	TQuestion
	TQuestionDot
	TQuestionQuestion
	TSemicolon
	TSlash
	TTilde

	// Assignments
	TAmpersandAmpersandEquals
	TAmpersandEquals
	TAsteriskAsteriskEquals
	TAsteriskEquals
	TBarBarEquals
	TBarEquals
	TCaretEquals
	TEquals
	TGreaterThanGreaterThanEquals
	TGreaterThanGreaterThanGreaterThanEquals
	TLessThanLessThanEquals
	TMinusEquals
	TPercentEquals
	TPlusEquals
	TQuestionQuestionEquals
	TSlashEquals

	// Class-private fields and methods
	TPrivateIdentifier

	// Identifiers
	TIdentifier
	TEscapedKeyword

	// Reserved words
	TBreak
	TCase
	TCatch
	TClass
	TConst
	TContinue
	TDebugger
	TDefault
	TDelete
	TDo
	TElse
	TEnum
	TExport
	TExtends
	TFalse
	TFinally
	TFor
	TFunction
	TIf
	TImport
	TIn
	TInstanceof
	TNew
	TNull
	TReturn
	TSuper
	TSwitch
	TThis
	TThrow
	TTrue
	TTry
	TTypeof
	TVar
	TVoid
	TWhile
	TWith

	// Strict mode reserved words
	TImplements
	TInterface
	TLet
	TPackage
	TPrivate
	TProtected
	TPublic
	TStatic
	TYield
)

var keywords = map[string]T{
	// Reserved words
	"break":      TBreak,
	"case":       TCase,
	"catch":      TCatch,
	"class":      TClass,
	"const":      TConst,
	"continue":   TContinue,
	"debugger":   TDebugger,
	"default":    TDefault,
	"delete":     TDelete,
	"do":         TDo,
	"else":       TElse,
	"enum":       TEnum,
	"export":     TExport,
	"extends":    TExtends,
	"false":      TFalse,
	"finally":    TFinally,
	"for":        TFor,
	"function":   TFunction,
	"if":         TIf,
	"import":     TImport,
	"in":         TIn,
	"instanceof": TInstanceof,
	"new":        TNew,
	"null":       TNull,
	"return":     TReturn,
	"super":      TSuper,
	"switch":     TSwitch,
	"this":       TThis,
	"throw":      TThrow,
	"true":       TTrue,
	"try":        TTry,
	"typeof":     TTypeof,
	"var":        TVar,
	"void":       TVoid,
	"while":      TWhile,
	"with":       TWith,

	// Strict mode reserved words
	"implements": TImplements,
	"interface":  TInterface,
	"let":        TLet,
	"package":    TPackage,
	"private":    TPrivate,
	"protected":  TProtected,
	"public":     TPublic,
	"static":     TStatic,
	"yield":      TYield,
}

// Token represents an Ecmascript token
type Token struct {
	typ  T
	text string   // empty for numbers
	num  *big.Int // nil for non-numbers
}

func mkToken(typ T, text string) *Token {
	if typ == TNumericLiteral {
		var z big.Int
		num, ok := z.SetString(text, 0)
		if !ok {
			log.Fatalf("bad number syntax: %s", text)
		}
		return &Token{TNumericLiteral, "", num}
	}

	t, ok := keywords[text]

	if !ok {
		t = typ
	}

	return &Token{t, text, nil}
}
