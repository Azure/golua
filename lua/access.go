package lua

// CheckType checks whether the function argument at index has type typ.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_checktype
func (state *State) CheckType(index int, typ Type) {
	if state.TypeAt(index) != typ {
		typeError(state, index, typ)
	}
}

// CheckNumber checks whether the function argument arg is a number and
// returns this number.
//
// See http://www.lua.org/manual/5.3/manual.html#luaL_checknumber
func (state *State) CheckNumber(index int) float64 {
	f64, ok := state.TryFloat(index)
	if !ok {
		typeError(state, index, NumberType)
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
		typeError(state, index, StringType)
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
	if state.IsNoneOrNil(index) {
		return optStr
	}
	return state.CheckString(index)
}

// OptNumber checks if the argument at index is a float and returns the number.
// Otherwise, if absent or nil, returns optNum.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_optnumber
func (state *State) OptNumber(index int, optNum float64) float64 {
	if state.IsNoneOrNil(index) {
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
	if state.TypeAt(index) != IntType {
		return optInt
	}
	return state.CheckInt(index)
}

// ToString converts the Lua value at the given index to a Go string. The Lua value must be a string
// or a number; otherwise, the function returns ("", false). If the value is a number, then ToString
// also changes the actual value in the stack to a string. (This change confuses Next(...) when ToString
// is applied to keys during a table traversal.)
//
// ToString returns a copy of the string inside the Lua state.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_tolstring
// func (state *State) ToString(index int) (string, bool) {
func (state *State) ToString(index int) string {
	val := state.get(index)
 	if meta := state.metafield(val, "__tostring"); !IsNone(meta) {
        if cls, ok := meta.(*Closure); ok {
            state.frame().push(cls)
            state.frame().push(val)
            state.Call(1, 1)
            val = state.frame().pop()
        }
    }
	s, ok := toString(val)
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
	if ok {
		state.set(index, Int(i64))
	}
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
func (state *State) ToThread(index int) *Thread {
	if v, ok := state.get(index).(*Thread); ok {
		return v
	}
	return nil
}

// IsThread returns true if the value at the given index is a thread; otherwise false.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_isthread
func (state *State) IsThread(index int) bool {
	return state.TypeAt(index) == ThreadType
}

// IsNumber returns true if the value at the given index is a number or a string
// convertible to a number, and false otherwise.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_isnumber
func (state *State) IsNumber(index int) bool { return state.TypeAt(index) == FloatType || state.IsInt(index) }

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

// IsInt returns true if the value at the given index is an integer
// (that is, the value is a number and is represented as an integer),
// and false otherwise.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_isinteger
func (state *State) IsInt(index int) bool { return state.TypeAt(index) == IntType }

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