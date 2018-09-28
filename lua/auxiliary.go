package lua

// If the registry already has the key name, return false. Otherwise, creates a new table to
// be used as a metatable for userdata, adds to this new table the pair __name = name, adds
// to the registry the pair [name] = new table, and returns true. The entry __name is used
// by some error-reporting functions.
//
// In both cases pushes onto the stack the final value associated with name in the registry.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_newmetatable
func (state *State) NewMetaTable(name string) bool {
	if state.GetMetaTable(name) != NilType {
		return false
	}
	state.Pop()
	state.NewTableSize(0, 2) 
	state.Push(name)
	state.SetField(-2, "__name") // metatable.__name = name
	state.PushIndex(-1)
	state.SetField(RegistryIndex, name) // registry.name = metatable
	return true
}

// GetSubTable ensures that stack[index][field] has a table and pushes
// that table onto the stack.
//
// Returns true if the table already exists at index; otherwise false
// if the table didn't exist but was created.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_getsubtable
func (state *State) GetSubTable(index int, field string) bool {
    if state.GetField(index, field) == TableType {
        return true               // table already exists
    }
    state.Pop()                   // remove previous result
    index = state.AbsIndex(index) // in frame stack
    state.NewTable()              // create table
    state.PushIndex(-1)           // copy to be left at top
    state.SetField(index, field)  // assign new table to field
    return false
}

// GetMetaField pushes onto the stack the field event from the metatable of the object
// at index and returns the type of the pushed value. If the object does not have a 
// metatable, or if the metatable does not have this field, pushes nothing and returns
// NilType.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_getmetafield
func (state *State) GetMetaField(index int, event string) Type {
	meta := state.metafield(state.get(index), event)
	if !IsNone(meta) {
		state.Push(meta)
	}
	return meta.Type()
}

// GetMetaTable pushes onto the stack the metatable associated with name in the registry
// (see luaL_newtable) (nil if there is no metatable associated with that name). Returns
// the type of the pushed value.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_getmetatable
func (state *State) GetMetaTable(name string) Type {
    unimplemented("GetMetaTable")
    return NilType
}

// CheckType checks whether the function argument at index has type typ.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_checktype
func (state *State) CheckType(index int, typ Type) {
	if state.TypeAt(index) != typ {
		typeError(state, index, typ)
	}
}

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
	if state.TypeAt(index) == NilType {
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
	if state.TypeAt(index) == NilType {
		return optStr
	}
	return state.CheckString(index)
}

// OptInt checks if the function argument at index is an integer (or convertible to), and returns
// the integer; Otherwise, if absent or nil, returns optInt.
//
// This function invokes CheckInt so the caveats of that function apply here.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_optinteger
func (state *State) OptInt(index int, optInt int64) int64 {
	if state.TypeAt(index) == IntType {
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
func (state *State) ToString(index int) (string, bool) {
	// TODO: check __tostring ??
	s, ok := toString(state.get(index))
	if ok {
		state.set(index, String(s))
	}
	return s, ok
}

// ToBool converts the Lua value at the given index to a Go boolean value. Like all tests in Lua, ToBool
// returns true for any Lua value different from false and nil; otherwise it returns false.
//
// If you want to accept only actual boolean values, use IsBool to test the value's type.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_toboolean
func (state *State) ToBool(index int) bool {
	return Truth(state.Value(index))
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

// IsNoneOrNil returns true if the given index is not valid or if the value at this index is nil and
// false otherwise.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_isnoneornil
func (state *State) IsNoneOrNil(index int) bool {
	return state.IsNone(index) || state.IsNil(index)
}

// IsNone returns true if the given index is not valid, and false otherwise.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_isnone
func (state *State) IsNone(index int) bool { return state.TypeAt(index) == NilType }

// IsNil returns true if the value at the given index is nil, and false otherwise.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_isnil
func (state *State) IsNil(index int) bool { return IsNone(state.get(index)) }

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
	if state.isValid(index) {
		return state.get(index).Type()
	}
	return NilType
}

func (state *State) isValid(index int) bool {
	switch {
		case index == RegistryIndex: // registry
			index = UpValueIndex(index) - 1
			cls := state.frame().closure
			return cls != nil && index < len(cls.upvals)
		case index < RegistryIndex: // upvalues
			return true
	}
	var (
		abs = state.frame().absindex(index)
		top = state.frame().gettop()
	)
	return abs > 0 && abs <= top
}