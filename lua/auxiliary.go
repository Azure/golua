package lua

import (
	"fmt"
)

// CheckString checks whether the function argument at index is a string and returns
// this string. This function uses ToString to get its result, so all conversions
// and caveats of that function apply here.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_checkstring
func (state *State) CheckString(index int) string {
	v, ok := state.ToString(index)
	if !ok {
		typeError(state, index, StringType)
	}
	return v
}

// ToString converts the Lua value at the given index to a Go string. IThe Lua value must be a string
// or a number; otherwise, the function returns ("", false). If the value is a number, then ToString
// also changes the actual value in the stack to a string. (This change confuses Next(...) when ToString
// is applied to keys during a table traversal.)
//
// ToString returns a copy of the string inside the Lua state.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_tolstring
func (state *State) ToString(index int) (string, bool) {
	switch v := state.get(index).(type) {
		case String:
			return string(v), true
		case Float:
			s := fmt.Sprintf("%v", float64(v))
			state.set(index, String(s))
			return s, true
		case Int:
			s := fmt.Sprintf("%v", int64(v))
			state.set(index, String(s))
			return s, true
	}
	return "", false
}

// ToBool converts the Lua value at the given index to a Go boolean value. Like all tests in Lua, ToBool
// returns true for any Lua value different from false and nil; otherwise it returns false.
//
// If you want to accept only actual boolean values, use IsBool to test the value's type.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_toboolean
func (state *State) ToBool(index int) bool { return Truth(state.Value(index)) }

// TypeAt returns the type of the value in the given valid index.
//
// TypeAt returns NilType for a non-valid (but acceptable) index.
//
// Otherwise, TypeAt returns one of:
//	LUA_TNUMBER
//	LUA_TBOOLEAN
//	LUA_TSTRING
//	LUA_TTABLE
//	LUA_TFUNCTION
//	LUA_TUSERDATA
//	LUA_TTHREAD
//	LUA_TLIGHTUSERDATA
func (state *State) TypeAt(index int) Type {
	return state.get(index).Type()
}