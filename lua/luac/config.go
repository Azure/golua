package luac

import "github.com/Azure/golua/lua/code"

// Maximum depth for nested Go calls and syntactical nested non-terminals
// in a program.
//
// Value must fit in a byte.
const maxCalls = 200

// Maximum number of local variables per function.
//
// Must be smaller than 250, due to the bytecode format.
const maxVars = 200

// maximum number of registers in a Lua function (must fit in 8-bits).
const maxRegs = 255

// option for multiple returns in 'lua_pcall' and 'lua_call'
const multRet = -1

// MaxUp is the maximum number of upvalues in a closure (both Go and Lua).
//
// Value must fit in a VM register.
const maxUp = 255

// maxID gives the maximum size for the description of the source
// of a function in debug information.
const maxID = 60

// maximum integer size
const maxInt = int(^uint(0)>>1)

// noJump marks the end of a patch list. It is an invalid value both as an absolute
// address, and as a list link (would link an element to itself).
const noJump = -1

// noReg represents an invalid register that fits in 8 bits.
const noReg = code.MaxArgA

// fieldsPerFlush is the number of list items to accumulate in a table before emitting a SETLIST instruction.
const fieldsPerFlush = 50

// A short literal string can be delimited by matching single or double quotes, and can contain the following
// C-like escape sequences:
//
//		'\a' (bell),
//		'\b' (backspace),
//		'\f' (form feed),
//		'\n' (newline),
//		 '\r' (carriage return),
//		'\t' (horizontal tab),
//		'\v' (vertical tab),
//		'\\' (backslash),
//		'\"' (quotation mark [double quote]),
//		'\'' (apostrophe [single quote]).
//
// A backslash followed by a line break results in a newline in the string.
//
// The escape sequence '\z' skips the following span of white-space characters, including line breaks;
// it is particularly useful to break and indent a long literal string into multiple lines without adding
// the newlines and spaces into the string contents. A short literal string cannot contain unescaped line
// breaks nor escapes not forming a valid escape sequence.
//
// We can specify any byte in a short literal string by its numeric value (including embedded zeros).
// This can be done with the escape sequence \xXX, where XX is a sequence of exactly two hexadecimal digits,
// or with the escape sequence \ddd, where ddd is a sequence of up to three decimal digits. (Note that if a
// decimal escape sequence is to be followed by a digit, it must be expressed using exactly three digits.)
// 
// The UTF-8 encoding of a Unicode character can be inserted in a literal string with the escape sequence
// \u{XXX} (note the mandatory enclosing brackets), where XXX is a sequence of one or more hexadecimal digits
// representing the character code point.
//
// recognized escape sequences
var escapes = map[rune]rune{ 
	'a':  '\a',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'v':  '\v',
	'\\': '\\',
	'"':  '"', 
	'\'': '\'',
}

// priority for unary operators
const unaryPriority = 12

var priority = []struct{
	lhs int // left priority for each binary operator
	rhs int // right priority
}{
	{10, 10}, // '+'
	{10, 10}, // '-'
	{11, 11}, // '*'
	{11, 11}, // '%'
	{14, 13}, // '^' (right associative)
	{11, 11}, // '/'
	{11, 11}, // '//'
	{6, 6},   // '&'
	{4, 4},   // '|'
	{5, 5},   // '~'
	{7, 7},   // '<<'
	{7, 7},   // '>>'
	{9, 8},   // '..' (right associative)
	{3, 3},   // '=='
	{3, 3},   // '<'
	{3, 3},   // '<='
	{3, 3},   // ~=
	{3, 3},   // >
	{3, 3},   // >=
	{2, 2},   // and
	{1, 1},   // or
}

func binaryOp(op rune) code.Op {
	switch op {
		case tConcat:
			return code.OpConcat
		case tDivI:
			return code.OpDivI
		case tShl:
			return code.OpShl
		case tShr:
			return code.OpShr
		case tAnd:
			return code.OpAnd
		case tNe:
			return code.OpNe
		case tEq:
			return code.OpEq
		case tLe:
			return code.OpLe
		case tGe:
			return code.OpGe
		case tOr:
			return code.OpOr
		case '+':
			return code.OpAdd
		case '-':
			return code.OpSub
		case '*':
			return code.OpMul
		case '%':
			return code.OpMod
		case '^':
			return code.OpPow
		case '/':
			return code.OpDivF
		case '&':
			return code.OpBand
		case '|':
			return code.OpBor
		case '~':
			return code.OpBxor
		case '<':
			return code.OpLt
		case '>':
			return code.OpGt
	}
	return code.OpNone
}

func unaryOp(op rune) code.Op {
	switch op {
		case tNot:
			return code.OpNot
		case '-':
			return code.OpMinus
		case '~':
			return code.OpBnot
		case '#':
			return code.OpLen
	}
	return code.OpNone
}