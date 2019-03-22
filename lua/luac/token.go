package luac

import "fmt"

const reserved = 257

const (
	// terminal symbols denoted by reserved words
	tAnd   = reserved + iota
	tBreak
	tDo
	tElse
	tElseIf
	tEnd
	tFalse
	tFor
	tFunction
	tGoto
	tIf
	tIn
	tLocal
	tNil
	tNot
	tOr
	tRepeat
	tReturn
	tThen
	tTrue
	tUntil
	tWhile
	// other terminal symbols
	tDivI
	tConcat
	tDots
	tEq
	tGe
	tLe
	tNe
	tShl
	tShr
	tColon2
	tEOS
	tFloat
	tInt
	tName
	tString
)

var keywords = map[string]rune{
	"and":		tAnd,
	"break":	tBreak,
	"do":		tDo,
	"else":		tElse,
	"elseif":	tElseIf,
	"end":		tEnd,
	"false":	tFalse,
	"for":		tFor,
	"function":	tFunction,
	"goto":		tGoto,
	"if":		tIf,
	"in":		tIn,
	"local":	tLocal,
	"nil":		tNil,
	"not":		tNot,
	"or":		tOr,
	"repeat":	tRepeat,
	"return":	tReturn,
	"then":		tThen,
	"true":		tTrue,
	"until": 	tUntil,
	"while":    tWhile,
}

var tokens = []string{
	"and",
	"break",
	"do",
	"else",
	"elseif",
	"end",
	"false",
	"for",
	"function",
	"goto",
	"if",
	"in",
	"local",
	"nil",
	"not",
	"or",
	"repeat",
	"return",
	"then",
	"true",
	"until",
	"while",
	"//",
	"..",
	"...",
	"==",
	">=",
	"<=",
	"~=",
	"<<",
	">>",
	"::",
	"<eof>",
	"<number>",
	"<integer>",
	"<name>",
	"<string>",
}

type token struct {
	char rune
	ival int64
	sval string
	nval float64
}

func (tok token) String() string {
	if tok.char < reserved { // single-byte symbols?
		return fmt.Sprintf("'%c'", tok.char)
	}
	str := tokens[int(tok.char) - reserved]
	if tok.char < tEOS {
		return fmt.Sprintf("'%s'", str)
	}
	return str
}

func isHexDigit(r rune) bool { return isDigit(r) || 'a' <= r && r <= 'f' || 'A' <= r && r <= 'F' }
func isAlnum(r rune) bool { return isDigit(r) || isAlpha(r) }
func isSpace(r rune) bool { return r == ' ' || r == '\t' || r == '\f' || r == '\v' }
func isAlpha(r rune) bool { return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' }
func isIdent(r rune) bool { return isAlpha(r) || r == '_' }
func isNewLine(r rune) bool { return r == '\n' || r == '\r' }
func isDigit(r rune) bool { return '0' <= r && r <= '9' }