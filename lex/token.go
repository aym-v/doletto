package lex

// Type represents an ECMAScript token type
type Type int

// Definition of tokens constants
const (
	tokEndOfFile Type = iota
	tokSyntaxError

	tokHashbang

	tokNoSubstitutionTemplateLiteral
	tokNumericLiteral
	tokStringLiteral

	// Punctuation
	tokAmpersand
	tokAmpersandAmpersand
	tokAsterisk
	tokAsteriskAsterisk
	tokAt
	tokBar
	tokBarBar
	tokCaret
	tokCloseBrace
	tokCloseBracket
	tokCloseParen
	tokColon
	tokComma
	tokDot
	tokDotDotDot
	tokEqualsEquals
	tokEqualsEqualsEquals
	tokEqualsGreaterThan
	tokExclamation
	tokExclamationEquals
	tokExclamationEqualsEquals
	tokGreaterThan
	tokGreaterThanEquals
	tokGreaterThanGreaterThan
	tokGreaterThanGreaterThanGreaterThan
	tokLessThan
	tokLessThanEquals
	tokLessThanLessThan
	tokMinus
	tokMinusMinus
	tokOpenBrace
	tokOpenBracket
	tokOpenParen
	tokPercent
	tokPlus
	tokPlusPlus
	tokQuestion
	tokQuestionDot
	tokQuestionQuestion
	tokSemicolon
	tokSlash
	tokTilde

	// Assignments
	tokAmpersandAmpersandEquals
	tokAmpersandEquals
	tokAsteriskAsteriskEquals
	tokAsteriskEquals
	tokBarBarEquals
	tokBarEquals
	tokCaretEquals
	tokEquals
	tokGreaterThanGreaterThanEquals
	tokGreaterThanGreaterThanGreaterThanEquals
	tokLessThanLessThanEquals
	tokMinusEquals
	tokPercentEquals
	tokPlusEquals
	tokQuestionQuestionEquals
	tokSlashEquals

	// Class-private fields and methods
	tokPrivateIdentifier

	// Identifiers
	tokIdentifier
	tokEscapedKeyword

	// Reserved words
	tokBreak
	tokCase
	tokCatch
	tokClass
	tokConst
	tokContinue
	tokDebugger
	tokDefault
	tokDelete
	tokDo
	tokElse
	tokEnum
	tokExport
	tokExtends
	tokFalse
	tokFinally
	tokFor
	tokFunction
	tokIf
	tokImport
	tokIn
	tokInstanceof
	tokNew
	tokNull
	tokReturn
	tokSuper
	tokSwitch
	tokThis
	tokThrow
	tokTrue
	tokTry
	tokTypeof
	tokVar
	tokVoid
	tokWhile
	tokWith

	// Strict mode reserved words
	tokImplements
	tokInterface
	tokLet
	tokPackage
	tokPrivate
	tokProtected
	tokPublic
	tokStatic
	tokYield
)

var keywords = map[string]Type{
	// Reserved words
	"break":      tokBreak,
	"case":       tokCase,
	"catch":      tokCatch,
	"class":      tokClass,
	"const":      tokConst,
	"continue":   tokContinue,
	"debugger":   tokDebugger,
	"default":    tokDefault,
	"delete":     tokDelete,
	"do":         tokDo,
	"else":       tokElse,
	"enum":       tokEnum,
	"export":     tokExport,
	"extends":    tokExtends,
	"false":      tokFalse,
	"finally":    tokFinally,
	"for":        tokFor,
	"function":   tokFunction,
	"if":         tokIf,
	"import":     tokImport,
	"in":         tokIn,
	"instanceof": tokInstanceof,
	"new":        tokNew,
	"null":       tokNull,
	"return":     tokReturn,
	"super":      tokSuper,
	"switch":     tokSwitch,
	"this":       tokThis,
	"throw":      tokThrow,
	"true":       tokTrue,
	"try":        tokTry,
	"typeof":     tokTypeof,
	"var":        tokVar,
	"void":       tokVoid,
	"while":      tokWhile,
	"with":       tokWith,

	// Strict mode reserved words
	"implements": tokImplements,
	"interface":  tokInterface,
	"let":        tokLet,
	"package":    tokPackage,
	"private":    tokPrivate,
	"protected":  tokProtected,
	"public":     tokPublic,
	"static":     tokStatic,
	"yield":      tokYield,
}

// Token represents an Ecmascript token
type Token struct {
	typ  Type
	text string  // empty for numbers
	num  float64 // nil for non-numbers
}

func mkToken(typ Type, text string) *Token {
	t, ok := keywords[text]

	if !ok {
		t = typ
	}

	return &Token{t, text, 0}
}

func mkNumericLiteral(num float64) *Token {
	return &Token{tokNumericLiteral, "", num}
}
