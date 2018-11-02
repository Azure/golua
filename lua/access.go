package lua

import "fmt"

var _ = fmt.Println

// CheckUserData checks whether the function argument arg is a userdata of the type
// metaType (see NewMetaTable) and returns the userdata address (see ToUserData).
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_checkudata
func (state *State) CheckUserData(index int, metaType string) (data interface{}) {
	if data = state.TestUserData(index, metaType); data == nil {
		typeError(state, index, metaType)
	}
	return
}

// TestUserData is equivalent to CheckUserData, except that, when the test fails, it
// returns nil instead of raising an error.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_testudata
func (state *State) TestUserData(index int, metaType string) (data interface{}) {
	if data = state.ToUserData(index); data != nil {
		if state.GetMetaTableAt(index) {
			state.GetMetaTable(metaType)
			if !state.RawEqual(-1, -2) {
				data = nil
			}
			state.PopN(2)
		}
	}
	return
}

// CheckType checks whether the function argument at index has type typ.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_checktype
func (state *State) CheckType(index int, typ Type) {
	if state.TypeAt(index) != typ {
		typeError(state, index, typ.String())
	}
}

// CheckNumber checks whether the function argument arg is a number and
// returns this number.
//
// See http://www.lua.org/manual/5.3/manual.html#luaL_checknumber
func (state *State) CheckNumber(index int) float64 {
	f64, ok := state.TryFloat(index)
	if !ok {
		typeError(state, index, "number")
	}
	return f64
}

// CheckString checks whether the function argument at index is a string and returns
// this string. This function uses ToString to get its result, so all conversions
// and caveats of that function apply here.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_checkstring
func (state *State) CheckString(index int) string {
	str, ok := state.TryString(index)
	if !ok {
		typeError(state, index, "string")
	}
	return str
}

// CheckInt checks whether the function argument at index is an integer (or can be converted to an
// integer) and returns the value as an integer.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_checkinteger
func (state *State) CheckInt(index int) int64 {
	v, ok := toInteger(state.get(index))
	if !ok {
		intError(state, index)
	}
	return int64(v)
}

// CheckAny checks whether the function has an argument of any type (including nil)
// at position index.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_checkany
func (state *State) CheckAny(index int) {
	if state.TypeAt(index) == NoneType {
		argError(state, index, "value expected")
	}
}

// OptString checks if the argument at index is a string and returns this string;
// Otherwise, if absent or nil, returns optStr.
//
// This function invokes CheckString so the caveats of that function apply here.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-optstring
func (state *State) OptString(index int, optStr string) string {
	if state.TypeAt(index) == NoneType {
		return optStr
	}
	return state.CheckString(index)
}

// OptNumber checks if the argument at index is a float and returns the number.
// Otherwise, if absent or nil, returns optNum.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_optnumber
func (state *State) OptNumber(index int, optNum float64) float64 {
	if state.TypeAt(index) == NoneType {
		return optNum
	}
	return state.CheckNumber(index)
}

// OptInt checks if the function argument at index is an integer (or convertible to), and returns
// the integer; Otherwise, if absent or nil, returns optInt.
//
// This function invokes CheckInt so the caveats of that function apply here.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_optinteger
func (state *State) OptInt(index int, optInt int64) int64 {
	if state.TypeAt(index) == NoneType {
		return optInt
	}
	return state.CheckInt(index)
}

// ToUserData returns the value at the given index if type userdata. Otherwise, returns nil.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_touserdata
func (state *State) ToUserData(index int) interface{} {
	if udata, ok := state.get(index).(*Object); ok {
		return udata.data
	}
	return nil
}

// ToStringMeta converts any Lua value at the given index to a C string in a reasonable format. The resulting
// string is pushed onto the stack and also returned by the function. If len is not NULL, the function also
// sets *len with the string length.
//
// If the value has a metatable with a __tostring field, then luaL_tolstring calls the corresponding
// metamethod with the value as argument, and uses the result of the call as its result.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_tolstring
func (state *State) ToStringMeta(index int) string {
	if state.CallMeta(index, "__tostring") {
		s, ok := state.TryString(-1)
		if ok {
			return s
		}
		panic(fmt.Errorf("'__tostring' must return a string"))
	} else {
		switch kind := state.TypeAt(index); kind {
			case NumberType:
				if state.IsInt(index) {
					state.Push(fmt.Sprintf("%d", state.ToInt(index)))
				} else {
					state.Push(fmt.Sprintf("%.14g", state.ToNumber(index)))
				}
			case StringType:
				state.PushIndex(index)
			case BoolType:
				state.Push(state.ToBool(index))
			case NilType:
				state.Push("nil")
			default:
				// TODO: check __name metafield
				state.Push(fmt.Sprintf("%s: %p", kind, state.get(index)))
		}
	}
	return state.ToString(-1)
}

// ToNumber converts the Lua value at the given index to a Number. The Lua value must be a number or a string
// convertible to a number (see §3.4.3); otherwise, otherwise, returns 0.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_tonumber
func (state *State) ToNumber(index int) float64 {
	f64, ok := state.TryFloat(index)
	if !ok {
		return 0
	}
	state.set(index, Float(f64))
	return f64
}

// ToString converts the Lua value at the given index to a Go string. The Lua value must be a string
// or a number; otherwise, the function returns ("", false). If the value is a number, then ToString
// also changes the actual value in the stack to a string. (This change confuses Next(...) when ToString
// is applied to keys during a table traversal.)
//
// ToString returns a copy of the string inside the Lua state.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_tolstring
func (state *State) ToString(index int) string {
	s, ok := toString(state.get(index))
	if ok {
		state.set(index, String(s))
	}
	return s
}

// ToInt converts the Lua value at the given index to an int64. The Lua value must be an integer,
// or a number or string convertible to an integer (see §3.4.3); otherwise, returns 0.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_tointeger
func (state *State) ToInt(index int) int64 {
	i64, ok := state.TryInt(index)
	if !ok {
		return 0
	}
	state.set(index, Int(i64))
	return i64
}

// ToBool converts the Lua value at the given index to a Go boolean value. Like all tests in Lua, ToBool
// returns true for any Lua value different from false and nil; otherwise it returns false.
//
// If you want to accept only actual boolean values, use IsBool to test the value's type.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_toboolean
func (state *State) ToBool(index int) bool {
	if state.isValid(index) {
		return Truth(state.get(index))
	}
	return false
}

// ToThread converts the value at the given index to a Lua thread (*State). This value must be a thread;
// otherwise, the function returns nil.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_tothread
func (state *State) ToThread(index int) *State {
	if v, ok := state.get(index).(*Thread); ok {
		return v.State
	}
	return nil
}

// IsThread returns true if the value at the given index is a thread; otherwise false.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_isthread
func (state *State) IsThread(index int) bool {
	return state.TypeAt(index) == ThreadType
}

// IsString returns true if the value at the given index is a string or a number (which is
// always convertible to a string); otherwise, returns false.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_isstring
func (state *State) IsString(index int) bool {
	t := state.TypeAt(index)
	return t == StringType || t == NumberType
}

// IsNumber returns true if the value at the given index is a number or a string
// convertible to a number, and false otherwise.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_isnumber
func (state *State) IsNumber(index int) bool { _, ok := state.TryNumber(index); return ok }

// IsNoneOrNil returns true if the given index is not valid or if the value at this index is nil and
// false otherwise.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_isnoneornil
func (state *State) IsNoneOrNil(index int) bool { return state.IsNone(index) || state.IsNil(index) }

// IsNone returns true if the given index is not valid, and false otherwise.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_isnone
func (state *State) IsNone(index int) bool { return state.TypeAt(index) == NoneType }

// IsNil returns true if the value at the given index is nil, and false otherwise.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_isnil
func (state *State) IsNil(index int) bool { return state.TypeAt(index) == NilType }

// Returns true if the value at the given index is a boolean; otherwise falsse.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_isboolean
func (state *State) IsBool(index int) bool { return state.TypeAt(index) == BoolType }

// IsFunc returns true if the value at the given index is a function (either Go or Lua);
// otherwise returns false.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_isfunction
func (state *State) IsFunc(index int) bool { return state.TypeAt(index) == FuncType }

// IsGoFunc returns true if the value at the given index is a Go function;
// otherwise returns false.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_iscfunction
func (state *State) IsGoFunc(index int) bool {
	cls, ok := state.get(index).(*Closure)
	return ok && !cls.isLua()
}

// IsFloat returns true if the value at the given index is an float
// (that is, the value is a number and is represented as a float),
// and false otherwise.
func (state *State) IsFloat(index int) bool { return IsFloat(state.get(index)) }

// IsInt returns true if the value at the given index is an integer
// (that is, the value is a number and is represented as an integer),
// and false otherwise.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_isinteger
func (state *State) IsInt(index int) bool { return IsInt(state.get(index)) }

func (state *State) TryNumber(index int) (Number, bool) {
	if v, ok := toInteger(state.get(index)); ok {
		return v, ok
	}
	if v, ok := toFloat(state.get(index)); ok {
		return v, ok
	}
	return nil, false
}

func (state *State) TryString(index int) (string, bool) {
	return toString(state.get(index))
}

func (state *State) TryFloat(index int) (float64, bool) {
	f64, ok := toFloat(state.get(index))
	return float64(f64), ok
}

func (state *State) TryInt(index int) (int64, bool) {
	i64, ok := toInteger(state.get(index))
	return int64(i64), ok
}