package utf8

import (
	"fmt"
	"unicode/utf8"

	"github.com/Azure/golua/lua"
)

const unicodeMax = 0x10FFFF

//
// Lua Standard Library -- utf8
//

// Open opens the Lua standard utf8 library. This library provides basic support for UTF-8 encoding.
// It provides all its functions inside the table utf8. This library does not provide any support for
// Unicode other than the handling of the encoding. Any operation that needs the meaning of a character,
// such as character classification, is outside its scope.
//
// Unless stated otherwise, all functions that expect a byte position as a parameter assume that the given
// position is either the start of a byte sequence or one plus the length of the subject string. As in the
// string library, negative indices count from the end of the string.
//
// See https://www.lua.org/manual/5.3/manual.html#6.5
func Open(state *lua.State) int {
	// Create 'utf8' table.
	var utf8Funcs = map[string]lua.Func{
		"char":      lua.Func(utf8Char),
		"codepoint": lua.Func(utf8CodePoint),
		"codes":     lua.Func(utf8Codes),
		"len":       lua.Func(utf8Len),
		"offset":    lua.Func(utf8Offset),
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
// Receives zero or more integers, converts each one to its corresponding UTF-8
// byte sequence and returns a string with the concatenation of all these sequences.
//
// utf8.char(n1, n2, n3, ...) -> char(n1) .. char(n2) ..
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-utf8.char
func utf8Char(state *lua.State) int {
	runes := make([]rune, state.Top())
	for i := 1; i <= state.Top(); i++ {
		c := state.CheckInt(i)
		if c < 0 || c > unicodeMax {

			panic(fmt.Errorf("bad argument #1 to 'char' (value out of range)"))
		}
		runes[i-1] = rune(c)
	}
	var (
		buf = make([]byte, 6)
		str string
	)
	for _, r := range runes {
		w := utf8.EncodeRune(buf, r)
		str += string(buf[:w])
	}
	state.Push(str)
	return 1
}

// utf8.codes (s)
//
// Returns values so that the construction ```for p, c in utf8.codes(s) do body end```
// will iterate over all characters in string s, with p being the position (in bytes)
// and c the code point of each character.
//
// It raises an error if it meets any invalid byte sequence.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-utf8.codes
func utf8Codes(state *lua.State) int {
	iter := lua.Func(func(state *lua.State) int {
		var (
			s = state.CheckString(1)
			n = state.ToInt(2)
		)
		if n > 0 && n < int64(len(s)) {
			n++
			for n < int64(len(s)) && isContByte(s[n]) {
				n++
			}
		}
		if n >= int64(len(s)) {
			return 0
		}
		r, _ := utf8.DecodeRuneInString(s[n:])
		if r == utf8.RuneError {
			panic(fmt.Errorf("invalid UTF-8 code"))
		}
		state.Push(n + 1)
		state.Push(int64(r))
		return 2
	})
	state.CheckString(1)
	state.Push(iter)
	state.PushIndex(1)
	state.Push(0)
	return 3
}

// utf8.codepoint (s [, i [, j]])
//
// Returns the codepoints (as integers) from all characters in s that start between byte position
// i and j (both included). The default for i is 1 and for j is i. It raises an error if it meets
// any invalid byte sequence.
//
// codepoint(s, [i, [j]]) => returns codepoints for all characters that start in the range [i,j].
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-utf8.codepoint
func utf8CodePoint(state *lua.State) int {
	var (
		s = state.CheckString(1)
		i = state.OptInt(2, 1)
		j = state.OptInt(3, i)
		n = 0
	)
	if i = int64(strPos(len(s), int(i))); i < 1 {
		panic(fmt.Errorf("bad argument #2 to 'codepoint' (out of range)"))
	}
	if j = int64(strPos(len(s), int(j))); j > int64(len(s)) {
		panic(fmt.Errorf("bad argument #3 to 'codepoint' (out of range)"))
	}
	for s = s[i-1:]; i <= j; {
		r, w := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError {
			panic(fmt.Errorf("invalid UTF-8 code"))
		}
		state.Push(int64(r))
		i += int64(w)
		s = s[w:]
		n++
	}
	return n
}

// utf8.len (s [, i [, j]])
//
// Returns the number of UTF-8 characters in string s that start between positions
// i and j (both inclusive). The default for i is 1 and for j is -1. If it finds
// any invalid byte sequence, returns a false value plus the position of the first
// invalid byte.
//
// utf8.len(s [, i [, j]]) --> number of characters that start in the range [i,j],
// or nil + current position if 's' is not well formed in that interval.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-utf8.len
func utf8Len(state *lua.State) int {
	var (
		s = state.CheckString(1)
		i = int64(strPos(len(s), int(state.OptInt(2, 1))))
		j = int64(strPos(len(s), int(state.OptInt(3, -1))))
		n = int64(0)
	)
	state.ArgCheck(1 <= i && i <= int64(len(s))+1, 2, "initial position out of string")
	state.ArgCheck(j <= int64(len(s)), 3, "final position out of string")

	// if i = int64(strPos(len(s), int(i))); i > 1 {
	// 	panic(fmt.Errorf("initial position out of string"))
	// }
	// if j = int64(strPos(len(s), int(j))); j > int64(len(s)) {
	// 	panic(fmt.Errorf("final position out of string"))
	// }
	if i <= j {
		for s = s[i-1 : j]; len(s) > 0; {
			r, w := utf8.DecodeRuneInString(s)
			if r == utf8.RuneError {
				state.Push(false)
				state.Push(i)
				return 2
			}
			i += int64(w)
			s = s[w:]
			n++
		}
	}
	state.Push(n)
	return 1
}

// utf8.offset (s, n [, i])
//
// Returns the position (in bytes) where the encoding of the n-th character
// of s (counting from position i) starts. A negative n gets characters before
// position i. The default for i is 1 when n is non-negative and #s + 1 otherwise,
// so that utf8.offset(s, -n) gets the offset of the n-th character from the end of
// the string. If the specified character is neither in the subject nor right after
// its end, the function returns nil. As a special case, when n is 0 the function
// returns the start of the encoding of the character that contains the i-th byte
// of s.
//
// This function assumes that s is a valid UTF-8 string.
//
// offset(s, n, [i]) -> index where n-th character counting from position 'i'
// starts; 0 means character at 'i'.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-utf8.offset
func utf8Offset(state *lua.State) int {
	var (
		s = state.CheckString(1)
		n = state.CheckInt(2)
		i = int64(1)
	)
	if n < 0 {
		i = int64(len(s) + 1)
	}
	i = int64(strPos(len(s), int(state.OptInt(3, i))))
	if i < 1 || i > int64(len(s))+1 {
		panic(fmt.Errorf("position out of range"))
	}
	i--
	if n == 0 {
		for i > 0 && isContByte(s[i]) {
			i--
		}
	} else {
		if i < int64(len(s)) && isContByte(s[i]) {
			panic(fmt.Errorf("initial position is a continuation byte"))
		}
		if n < 0 {
			for n < 0 && i > 0 {
				for {
					i--
					if !(i > 0 && isContByte(s[i])) {
						break
					}
				}
				n++
			}
		} else {
			n--
			for n > 0 && i < int64(len(s)) {
				for {
					i++
					if i >= int64(len(s)) || !isContByte(s[i]) {
						break
					}
				}
				n--
			}
		}
	}
	if n == 0 {
		state.Push(int64(i + 1))
	} else {
		state.Push(nil)
	}
	return 1
}

// isContByte reports whether b is a continuation byte.
func isContByte(b byte) bool { return b&0xC0 == 0x80 }

// strPos converts a relative string position: negative means back
// from end. The absolute position is returned.
func strPos(len, pos int) int {
	switch {
	case pos >= 0:
		return pos
	case -pos > len:
		return 0
	default:
		return len + pos + 1
	}
}
