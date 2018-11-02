package str

import (
	gostrs "strings"
	"fmt"

	"github.com/Azure/golua/pkg/strings"
	"github.com/Azure/golua/pkg/packer"
    "github.com/Azure/golua/lua"
)
var _ = fmt.Println

//
// Lua Standard Library -- string
//

// See https://www.lua.org/manual/5.3/manual.html#6.4
func Open(state *lua.State) int {
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
	createStrMetaTable(state)

	// Return 'string' table.
	return 1
}

// createStrMetaTable sets the string library as the metatable for the Lua string type.
func createStrMetaTable(state *lua.State) {
	state.NewTableSize(0, 1) // table to be metatable for strings
	state.Push("") // dummy string
	state.PushIndex(-2) // copy table
	state.SetMetaTableAt(-2) // set table as metatable for strings
	state.Pop() // pop dummy string
	state.PushIndex(-2) // get string library
	state.SetField(-2, "__index") // metatable.__index = string
	state.Pop() // pop metatable
}

// string.byte (s [, i [, j]])
//
// Returns the internal numeric codes of the characters s[i], s[i+1], ..., s[j].
// The default value for i is 1; the default value for j is i. These indices are
// corrected following the same rules of function string.sub.
//
// Numeric codes are not necessarily portable across platforms.
//
// https://www.lua.org/manual/5.3/manual.html#pdf-string.byte
func strByte(state *lua.State) int {
	// s := state.CheckString(1)
	// i := state.OptInt(2, 1)
	// j := state.OptInt(3, i)
	// bytes := strs.ByteSlice(s, int(i), int(j))
	// for _, b := range bytes {
	// 	state.Push(int64(b))
	// }
	// return len(bytes)
	fmt.Println("string.byte: TODO")
	state.Debug(true)
	return 0
}

// string.char (···)
//
// Receives zero or more integers. Returns a string with length equal to the number
// of arguments, in which each character has the internal numeric code equal to its
// corresponding argument.
// 
// Numeric codes are not necessarily portable across platforms.
//
// https://www.lua.org/manual/5.3/manual.html#pdf-string.char
func strChar(state *lua.State) int {
	bytes := make([]byte, 0, state.Top())
	for i := 1; i <= state.Top(); i++ {
		bytes = append(bytes, byte(state.CheckInt(i)))
	}
	state.Push(string(bytes))
	return 1
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
	return 1
}

// string.format (formatstring, ···)
//
// strFormat returns a formatted version of its variable number of arguments
// following the description given in its first argument (which must be a
// string).
//
// The format string follows the same rules as the ISO C function sprintf.
// The only differences are that the options/modifiers *, h, L, l, n, and p
// are not supported and that there is an extra option, q.
//
// The q option formats a string between double quotes, using escape sequences
// when necessary to ensure that it can safely be read back by the Lua interpreter.
//
// For instace, the call
//
//		string.format('%q', 'a string with "quotes" and \n new line')
//
// may produce the string:
//
//		"a string \"quotes\" and \
//		new line"
//
// Options A, a, E, e, f, G, and g all expect a number as argument.
// Options c, d, i, o, u, X, and x expect an integer.
//
// When Lua is compiled with C89 compiler, options A and a (hex floats) do not
// support any modifier (flags, width, length).
//
// Option s expects a string; if its argument is not a string, it is converted to
// one following the same rules of tostring. If the option has any modifier (flags,
// width, length), the string argument should not contain embedded zeros.
//
// https://www.lua.org/manual/5.3/manual.html#pdf-string.format
func strFormat(state *lua.State) int {
	// args := make([]interface{}, state.Top()-1)
	// for argc := state.Top()-1; argc >= 1; argc-- {
	// 	args[argc-1] = state.Pop()
	// }
	// str, err := strs.Format(state.CheckString(1), args...)
	// if err != nil {
	// 	panic(err)
	// }
	// state.Push(str)
	// state.Push(fmt.Sprintf(state.CheckString(1), args...))

	// fmts := state.CheckString(1)

	// if len(fmts) <= 1 || !strings.ContainsRune(fmts, rune('%')) {
	// 	state.Push(fmts)
	// 	return 1
	// }
	// if state.Top() < 2 {
	// 	state.Push(fmts)
	// 	return 1
	// }
	// state.Push(doFmt(state, fmts, state.Top()))
	// return 1

	// var (
	// 	str strings.Builder
	// 	arg = 2
	// )
	// for _, opt := range scanfmt(state.CheckString(1)) {
	// 	if opt[0] == '%' {
	// 		if opt == "%%" {
	// 			str.WriteByte('%')
	// 		} else {
	// 			str.WriteString(fmtarg(state, opt, arg))
	// 			arg++
	// 		}
	// 	} else {
	// 		str.WriteString(opt)
	// 	}
	// }
	// state.Push(str.String())
	// return 1

	// str := format(state, state.CheckString(1), state.Top())
	// state.Push(str)
	// return 1

	fmt.Println("string.format: TODO")
	state.Debug(true)
	return 0
}

// string.gmatch (s, pattern)
//
// Returns an iterator function that, each time it is called, returns the next captures from
// pattern (see §6.4.1) over the string s. If pattern specifies no captures, then the whole
// match is produced in each call.
//
// As an example, the following loop will iterate over all the words from string s, printing
// one per line:
//
//     s = "hello world from Lua"
//     for w in string.gmatch(s, "%a+") do
//       print(w)
//     end
//
// The next example collects all pairs key=value from the given string into a table:
//
//     t = {}
//     s = "from=world, to=Lua"
//     for k, v in string.gmatch(s, "(%w+)=(%w+)") do
//       t[k] = v
//     end
//
// For this function, a caret '^' at the start of a pattern does not work as an anchor, as this
// would prevent the iteration.
//
// https://www.lua.org/manual/5.3/manual.html#pdf-string.gmatch
func strGmatch(state *lua.State) int {
	// var (
	// 	subj = state.CheckString(1)
	// 	patt = state.CheckString(2)
	// )
	// iter, err := strs.MatchIter(subj, patt)
	// if err != nil {
	// 	panic(err)
	// }
	// state.Push(func(state *lua.State) int {
	// 	next, more := iter()
	// 	if !more {
	// 		return 0
	// 	}
	// 	state.Push(next)
	// 	return 1
	// })
	// return 1
	fmt.Println("string.gmatch")
	state.Debug(true)
	return 0
}

// string.match (s, pattern [, init])
//
// Looks for the first match of pattern (see §6.4.1) in the string s. If it finds one, then match returns
// the captures from the pattern; otherwise it returns nil. If pattern specifies no captures, then the whole
// match is returned. A third, optional numeric argument init specifies where to start the search; its default
// value is 1 and can be negative.
//
// https://www.lua.org/manual/5.3/manual.html#pdf-string.match
func strMatch(state *lua.State) int {
	// var (
	// 	subj = state.CheckString(1)
	// 	patt = state.CheckString(2)
	// 	from = int(state.OptInt(3, 1))
	// )
	// if from < 0 {
	// 	if from = len(subj) + (from + 1); from < 1 {
	// 		from = 1
	// 	}
	// }
	// switch m, err := strs.Match(subj, patt, from); {
	// 	case err != nil:
	// 		panic(err)
	// 	case m == "":
	// 		state.Push(nil)
	// 	default:
	// 		state.Push(m)
	// }
	// return 1
	// s, p := state.CheckString(1), state.CheckString(2)
	// init := strPos(len(s), int(state.OptInt(3, 1)))
	// switch {
	// 	case init > len(s) + 1:
	// 		state.Push(nil)
	// 		return 1
	// 	case init < 1:
	// 		init = 1
	// }
	// matches := match.Match(s[init-1:], p)
	// if matches == nil {
	// 	state.Push(nil)
	// 	return 1
	// }
	// for _, matched := range matches {
	// 	state.Push(matched[1:])
	// }
	// return len(matches)
	fmt.Println("string.match: TODO")
	state.Debug(true)
	return 0
}

// string.find (s, pattern [, init [, plain]])
//
// Looks for the first match of pattern (see §6.4.1) in the string s.
//
// If it finds a match, then find returns the indices of s where this
// occurrence starts and ends; otherwise, it returns nil.
//
// A third, optional numeric argument init specifies where to start the
// search; its default value is 1 and can be negative. A value of true
// as a fourth, optional argument plain turns off the pattern matching
// facilities, so the function does a plain "find substring" operation,
// with no characters in pattern being considered magic.
//
// Note: If plain is given, then init must be given as well.
//
// If the pattern has captures, then in a successful match the captured
// values are also returned, after the two indices.
//
// https://www.lua.org/manual/5.3/manual.html#pdf-string.find
func strFind(state *lua.State) int {
	// s, p := state.CheckString(1), state.CheckString(2)
	// init := strPos(len(s), int(state.OptInt(3, 1)))
	// switch {
	// 	case init > len(s) + 1:
	// 		state.Push(nil)
	// 		return 1
	// 	case init < 1:
	// 		init = 1
	// }
	// loc := find(s, p, init, state.ToBool(4))
	// if loc == nil {
	// 	state.Push(nil)
	// 	return 1
	// }
	// state.Push(loc[0])
	// state.Push(loc[1])
	// return 2

	fmt.Println("string.find: TODO")
	state.Debug(true)
	return 0

	// var (
	// 	subj = state.CheckString(1)
	// 	patt = state.CheckString(2)
	// 	init = int(state.OptInt(3, int64(len(subj))))
	// )
	// switch init = strPos(len(subj), init); {
	// 	case init > len(subj) + 1:
	// 		state.Push(nil)
	// 		return 1
	// 	case init < 1:
	// 		init = 1
	// }
	// if plain || !strs.IsPattern(patt) {
	// 	if pos := strings.Index(subj[index-1:], patt); pos != -1 {
	// 		state.Push(pos+index)
	// 		state.Push(pos+index+len(patt)-1)
	// 		return 2
	// 	}
	// 	state.Push(nil)
	// 	return 1
	// }
	// switch m, err := strs.Find(subj, patt, index-1); {
	// 	case err != nil:
	// 		panic(err)
	// 	case m == nil:
	// 		state.Push(nil)
	// 		return 1
	// 	default:
	// 		state.Push(index+m[0])
	// 		state.Push(index+m[1]-1)
	// 		return 2
	// }
}

// string.gsub (s, pattern, repl [, n])
//
// Returns a copy of s in which all (or the first n, if given) occurrences of the pattern (see §6.4.1)
// have been replaced by a replacement string specified by repl, which can be a string, a table, or a
// function. gsub also returns, as its second value, the total number of matches that occurred.
//
// The name gsub comes from Global SUBstitution.
//
// If repl is a string, then its value is used for replacement. The character % works as an escape character:
// any sequence in repl of the form %d, with d between 1 and 9, stands for the value of the d-th captured substring.
// The sequence %0 stands for the whole match. The sequence %% stands for a single %.
//
// If repl is a table, then the table is queried for every match, using the first capture as the key.
//
// If repl is a function, then this function is called every time a match occurs, with all captured substrings passed
// as arguments, in order.
//
// In any case, if the pattern specifies no captures, then it behaves as if the whole pattern was inside a capture.
//
// If the value returned by the table query or by the function call is a string or a number, then it is used as the
// replacement string; otherwise, if it is false or nil, then there is no replacement (that is, the original match is
// kept in the string).
//
// Here are some examples:
//
//   x = string.gsub("hello world", "(%w+)", "%1 %1")
//   --> x="hello hello world world"
//
//   x = string.gsub("hello world", "%w+", "%0 %0", 1)
//   --> x="hello hello world"
//     
//   x = string.gsub("hello world from Lua", "(%w+)%s*(%w+)", "%2 %1")
//   --> x="world hello Lua from"
//     
//   x = string.gsub("home = $HOME, user = $USER", "%$(%w+)", os.getenv)
//   --> x="home = /home/roberto, user = roberto"
//     
//   x = string.gsub("4+5 = $return 4+5$", "%$(.-)%$", function (s)
//         return load(s)()
//       end)
//   --> x="4+5 = 9"
//     
//   local t = {name="lua", version="5.3"}
//   x = string.gsub("$name-$version.tar.gz", "%$(%w+)", t)
//   --> x="lua-5.3.tar.gz"
//
// https://www.lua.org/manual/5.3/manual.html#pdf-string.gsub
func strGsub(state *lua.State) int {
	subj := state.CheckString(1)
	patt := state.CheckString(2)
	upto := int(state.OptInt(4, int64(len(subj))))
	var (
		s string
		n int
	)
	switch state.TypeAt(3) {
	case lua.StringType:
		repl := state.CheckString(3)
		s, n = strings.GsubStrAll(
			subj, 
			patt,
			repl,
			upto,
		)
	// case lua.TableType:
	// case lua.FuncType:
	default:
		state.CheckAny(3)
	}
	state.Push(s)
	state.Push(n)
	return 2
}

// string.len (s)
//
// Receives a string and returns its length. The empty string "" has length 0.
// Embedded zeros are counted, so "a\000bc\000" has length 5.
//
// https://www.lua.org/manual/5.3/manual.html#pdf-string.len
func strLen(state *lua.State) int {
	state.Push(len(state.CheckString(1)))
	return 1
}

// string.lower (s)
//
// Receives a string and returns a copy of this string with all uppercase letters changed to lowercase.
// All other characters are left unchanged. The definition of what an uppercase letter is depends on the
// current locale.
//
// https://www.lua.org/manual/5.3/manual.html#pdf-string.lower
func strLower(state *lua.State) int {
	state.Push(gostrs.ToLower(state.CheckString(1)))
	return 1
}

// string.packsize (fmt)
//
// Returns the size of a string resulting from string.pack with the given format.
// The format string cannot have the variable-length options 's' or 'z' (see §6.4.2).
//
// https://www.lua.org/manual/5.3/manual.html#pdf-string.packsize
func strPackSize(state *lua.State) int {
	size, err := packer.Size(state.CheckString(1))
	if err != nil {
		state.Errorf("%v", err)
	}
	state.Push(size)
	return 1
}

// string.unpack (fmt, s [, pos])
//
// Returns the values packed in string s (see string.pack) according to the format string fmt (see §6.4.2).
// An optional pos marks where to start reading in s (default is 1). After the read values, this function
// also returns the index of the first unread byte in s.
//
// https://www.lua.org/manual/5.3/manual.html#pdf-string.unpack
func strUnpack(state *lua.State) int {
	fmt.Println("string.unpack: TODO")
	state.Debug(true)
	return 0
}

// string.pack (fmt, v1, v2, ···)
//
// Returns a binary string containing the values v1, v2, etc. packed (that is, serialized
// in binary form) according to the format string fmt (see §6.4.2).
//
// https://www.lua.org/manual/5.3/manual.html#pdf-string.pack
func strPack(state *lua.State) int {
	// b, err := strs.Pack(state.CheckString(1), state.PopN(state.Top()-1)...)
	// if err != nil {
	// 	panic(err)
	// }
	// state.Push(string(b))
	// return 1
	fmt.Println("string.pack: TODO")
	state.Debug(true)
	return 0
}

// string.rep (s, n [, sep])
//
// Returns a string that is the concatenation of n copies of the string s separated by
// the string sep. The default value for sep is the empty string (that is, no separator).
//
// Returns the empty string if n is not positive.
//
// Note: It is very easy to exhaust the memory of your
// machine with a single call to this function.
//
// https://www.lua.org/manual/5.3/manual.html#pdf-string.rep
func strRep(state *lua.State) int {
	// s, err := repeat(state.CheckString(1), state.OptString(3, ""), state.CheckInt(2))
	// if err != nil {
	// 	panic(err)
	// }
	// state.Push(s)
	// return 1
	fmt.Println("string.rep: TODO")
	state.Debug(true)
	return 0
}

// string.reverse (s)
//
// Returns a string that is the string s reversed.
//
// https://www.lua.org/manual/5.3/manual.html#pdf-string.reverse
func strReverse(state *lua.State) int {
	// arg := state.CheckString(1)
	// str := strs.Reverse(arg)
	// state.Push(str)
	// return 1
	fmt.Println("string.reverse: TODO")
	state.Debug(true)
	return 0
}

// string.sub (s, i [, j])
//
// Returns the substring of s that starts at i and continues until j; i and j can
// be negative. If j is absent, then it is assumed to be equal to -1 (which is the
// same as the string length).
//
// In particular, the call string.sub(s,1,j) returns a prefix of s with length j,
// and string.sub(s,-i) (for i > 0) returns a suffix of s with length i.
//
// If, after the translation of negative indices, i is less than 1, it is corrected to 1.
//
// If j is greater than the string length, it is corrected to that length.
//
// If, after these corrections, i is greater than j, the function returns the empty string.
//
// https://www.lua.org/manual/5.3/manual.html#pdf-string.sub
func strSub(state *lua.State) int {
	// s := state.CheckString(1)
	// i := state.OptInt(2, 1)
	// j := state.OptInt(3, -1)
	// state.Push(strs.SubString(s, int(i), int(j)))
	// return 1
	fmt.Println("string.sub: TODO")
	state.Debug(true)
	return 0
}

// string.upper (s)
//
// Receives a string and returns a copy of this string with all lowercase letters changed
// to uppercase. All other characters are left unchanged. The definition of what a lowercase
// letter is depends on the current locale.
//
// https://www.lua.org/manual/5.3/manual.html#pdf-string.upper
func strUpper(state *lua.State) int {
	state.Push(gostrs.ToUpper(state.CheckString(1)))
	return 1
}