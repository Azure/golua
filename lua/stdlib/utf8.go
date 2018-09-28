package stdlib

import (
	//"unicode/utf8"

    "github.com/Azure/golua/lua"
)

//
// Lua Standard Library -- utf8
//

// OpenUTF8 opens the Lua standard utf8 library. This library provides basic support for UTF-8 encoding.
// It provides all its functions inside the table utf8. This library does not provide any support for
// Unicode other than the handling of the encoding. Any operation that needs the meaning of a character,
// such as character classification, is outside its scope.
//
// Unless stated otherwise, all functions that expect a byte position as a parameter assume that the given
// position is either the start of a byte sequence or one plus the length of the subject string. As in the
// string library, negative indices count from the end of the string.
//
// See https://www.lua.org/manual/5.3/manual.html#6.5
func OpenUTF8(state *lua.State) int {
	// Create 'utf8' table.
	var utf8Funcs = map[string]lua.Func{
		"char":        lua.Func(utf8Char),
		"codepoint":   lua.Func(utf8CodePoint),
		"codes":       lua.Func(utf8Codes),
		"len":         lua.Func(utf8Len),
		"offset":      lua.Func(utf8Offset),
	}
	state.NewTableSize(0, len(utf8Funcs))
	state.SetFuncs(utf8Funcs, 0)

	// pattern to match a single UTF-8 character.
	const pattern = `[\0-\x7F\xC2-\xF4][\x80-\xBF]*`

	// The pattern (a string, not a function) "[\0-\x7F\xC2-\xF4][\x80-\xBF]*" (see §6.4.1),
	// which matches exactly one UTF-8 byte sequence, assuming that the subject is a valid
	// UTF-8 string.
	state.Push(pattern)
	state.SetField(-2, "charpattern")

	// Return 'utf8' table.
	return 1
}

// utf8.char (···)
//
// Receives zero or more integers, converts each one to its corresponding UTF-8 byte sequence and
// returns a string with the concatenation of all these sequences.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-utf8.char
func utf8Char(state *lua.State) int {
	unimplemented("utf8.char")
	return 0
}

// utf8.codes (s)
//
// Returns values so that the construction ```for p, c in utf8.codes(s) do body end```
// will iterate over all characters in string s, with p being the position (in bytes)
// and c the code point of each character.
//
// It raises an error if it meets any invalid byte sequence.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-utf8.codepoint
func utf8CodePoint(state *lua.State) int {
	unimplemented("utf8.codepoint")
	return 0
}

// utf8.codepoint (s [, i [, j]])
//
// Returns the codepoints (as integers) from all characters in s that start between byte position i and j (both included).
// The default for i is 1 and for j is i. It raises an error if it meets any invalid byte sequence.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-utf8.codes
func utf8Codes(state *lua.State) int {
	unimplemented("utf8.codes")
	return 0
}

// utf8.len (s [, i [, j]])
//
// Returns the number of UTF-8 characters in string s that start between positions i and j (both inclusive).
// The default for i is 1 and for j is -1. If it finds any invalid byte sequence, returns a false value plus
// the position of the first invalid byte.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-utf8.len
func utf8Len(state *lua.State) int {
	unimplemented("utf8.len")
	return 0
}

// utf8.offset (s, n [, i])
//
// Returns the position (in bytes) where the encoding of the n-th character of s (counting from position i) starts.
// A negative n gets characters before position i. The default for i is 1 when n is non-negative and #s + 1 otherwise,
// so that utf8.offset(s, -n) gets the offset of the n-th character from the end of the string. If the specified character
// is neither in the subject nor right after its end, the function returns nil. As a special case, when n is 0 the function
// returns the start of the encoding of the character that contains the i-th byte of s.
//
// This function assumes that s is a valid UTF-8 string.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-utf8.offset
func utf8Offset(state *lua.State) int {
	unimplemented("utf8.offset")
	return 0
}