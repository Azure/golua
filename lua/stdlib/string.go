package stdlib

import (
	"fmt"
    "github.com/Azure/golua/lua"
)
var _ = fmt.Println

//
// Lua Standard Library -- string
//

// See https://www.lua.org/manual/5.3/manual.html#6.4
func OpenString(state *lua.State) int {
	// Create 'string' table.
	var strFuncs = map[string]lua.Func{
		"byte":     lua.Func(strByte),
		"char":     lua.Func(strChar),
		"dump":     lua.Func(strDump),
		"find":     lua.Func(strFind),
		"format":   lua.Func(strFormat),
		"gmatch":   lua.Func(strGmatch),
		"gsub":     lua.Func(strGsub),
		"len": 	    lua.Func(strLen),
		"lower":    lua.Func(strLower),
		"match":    lua.Func(strMatch),
		"pack":     lua.Func(strPack),
		"packsize": lua.Func(strPackSize),
		"rep": 		lua.Func(strRep),
		"reverse":  lua.Func(strReverse),
		"sub": 	    lua.Func(strSub),
		"unpack":   lua.Func(strUnpack),
		"upper":    lua.Func(strUpper),
	}
	state.NewTableSize(0, len(strFuncs))
	state.SetFuncs(strFuncs, 0)

	// createStrMetaTable(state)

	// Return 'string' table.
	return 1
}

// https://www.lua.org/manual/5.3/manual.html#pdf-string.byte
func strByte(state *lua.State) int {
	unimplemented("string.byte")
	return 0
}

// https://www.lua.org/manual/5.3/manual.html#pdf-string.char
func strChar(state *lua.State) int {
	unimplemented("string.char")
	return 0
}

// string.dump (function [, strip])
//
// Returns a string containing a binary representation (a binary chunk) of the given function, so that
// a later load on this string returns a copy of the function (but with new upvalues). If strip is a
// true value, the binary representation may not include all debug information about the function, to
// save space.
//
// Functions with upvalues have only their number of upvalues saved. When (re)loaded, those upvalues
// receive fresh instances containing nil. (You can use the debug library to serialize and reload the
// upvalues of a function in a way adequate to your needs.)
//
// https://www.lua.org/manual/5.3/manual.html#pdf-string.dump
func strDump(state *lua.State) int {
	state.CheckType(1, lua.FuncType)
	strip := state.ToBool(2)
	state.SetTop(1)
	state.Push(string(state.Dump(strip)))
	state.Debug(true)
	return 1
}

// https://www.lua.org/manual/5.3/manual.html#pdf-string.find
func strFind(state *lua.State) int {
	unimplemented("string.find")
	return 0
}

// strFormat returns a formatted version of its variable number of arguments following
// the description given in its first argument (which must be a string). The format
// string follows the same rules as the ISO C function sprintf. The only differences
// are that the options/modifiers *, h, L, l, n, and p are not supported and that there
// is an extra option, q.
//
// The q option formats a string between double quotes, using escape sequences when necessary
// to ensure that it can safely be read back by the Lua interpreter. For instace, the call
//
//		string.format('%q', 'a string with "quotes" and \n new line')
//
// may produce the string:
//
//		"a string \"quotes\" and \
//		new line"
//
// Options A, a, E, e, f, G, and g all expect a number as argument. Options c, d, i, o, u, X,
// and x expect an integer. When Lua is compiled with C89 compiler, options A and a (hex floats)
// do not support any modifier (flags, width, length).
//
// Option s expects a string; if its argument is not a string, it is converted to one following
// the same rules of tostring. If the option has any modifier (flags, width, length), the string
// argument should not contain embedded zeros.
//
// https://www.lua.org/manual/5.3/manual.html#pdf-string.format
func strFormat(state *lua.State) int {
	state.Push("string.format: TODO")
	return 1
}

// https://www.lua.org/manual/5.3/manual.html#pdf-string.gmatch
func strGmatch(state *lua.State) int {
	unimplemented("string.gmatch")
	return 0
}

// https://www.lua.org/manual/5.3/manual.html#pdf-string.gsub
func strGsub(state *lua.State) int {
	unimplemented("string.gsub")
	return 0
}

// https://www.lua.org/manual/5.3/manual.html#pdf-string.len
func strLen(state *lua.State) int {
	unimplemented("string.len")
	return 0
}

// https://www.lua.org/manual/5.3/manual.html#pdf-string.lower
func strLower(state *lua.State) int {
	unimplemented("string.lower")
	return 0
}

// https://www.lua.org/manual/5.3/manual.html#pdf-string.match
func strMatch(state *lua.State) int {
	unimplemented("string.match")
	return 0
}

// https://www.lua.org/manual/5.3/manual.html#pdf-string.pack
func strPack(state *lua.State) int {
	unimplemented("string.pack")
	return 0
}

// https://www.lua.org/manual/5.3/manual.html#pdf-string.packsize
func strPackSize(state *lua.State) int {
	unimplemented("string.packsize")
	return 0
}

// https://www.lua.org/manual/5.3/manual.html#pdf-string.rep
func strRep(state *lua.State) int {
	unimplemented("string.rep")
	return 0
}

// https://www.lua.org/manual/5.3/manual.html#pdf-string.reverse
func strReverse(state *lua.State) int {
	unimplemented("string.reverse")
	return 0
}

// https://www.lua.org/manual/5.3/manual.html#pdf-string.sub
func strSub(state *lua.State) int {
	unimplemented("string.sub")
	return 0
}

// https://www.lua.org/manual/5.3/manual.html#pdf-string.unpack
func strUnpack(state *lua.State) int {
	unimplemented("string.unpack")
	return 0
}

// https://www.lua.org/manual/5.3/manual.html#pdf-string.upper
func strUpper(state *lua.State) int {
	unimplemented("string.upper")
	return 0
}