package lua

// Access functions (stack -> Go)

func (state *State) IsUserData(index int) bool {
	return state.index(index).Type() == UserDataType
}

func (state *State) IsNumber(index int) bool {
	return state.index(index).Type() == NumberType
}

func (state *State) IsString(index int) bool {
	return state.index(index).Type() == StringType
}

func (state *State) IsThread(index int) bool {
	return state.index(index).Type() == ThreadType
}

func (state *State) IsTable(index int) bool {
	return state.index(index).Type() == TableType
}

func (state *State) IsFunc(index int) bool {
	return state.index(index).Type() == FuncType
}

func (state *State) IsBool(index int) bool {
	return state.index(index).Type() == BoolType
}

func (state *State) ToUserData(index int) (*Object, bool) {
	unimplemented("ToUserData")
	return nil, false
}

func (state *State) ToNumber(index int) (float64, bool) {
	unimplemented("ToNumner")
	return 0, false
}

func (state *State) ToString(index int) (string, bool) {
	unimplemented("ToString")
	return "", false
}

func (state *State) ToFunc(index int) (Func, bool) {
	unimplemented("ToFunc")
	return nil, false
}

// ToBool converts the Lua value at the given index to a Go boolean value.
//
// Like all tests in Lua, ToBool returns true for any Lua value different from false and nil;
// otherwise it returns false.
//
// If you want to accept only actual boolean values, use lua_isboolean to test the value's type.
func (state *State) ToBool(index int) bool {
	switch v := state.index(index); v.Type() {
		case BoolType:
			return bool(v.(Bool))
		case NilType:
			return false
	}
	return true
}

func (state *State) ToInt(index int) (int64, bool) {
	unimplemented("ToInt")
	return 0, false
}


func (state *State) ToThread(index int) (*Thread, bool) {
	unimplemented("ToThread")
	return nil, false
}

// Returns the raw "length" of the value at the given index: for strings, this is the string length;
// for tables, this is the result of the length operator (‘\#’) with no metamethods; for userdata,
// this is the size of the block of memory allocated for the userdata; for other values, it is 0.
func (state *State) RawLen(index int) int {
	unimplemented("RawLen")
	return 0
}

func (state *State) TypeAt(index int) Type {
	unimplemented("TypeAt")
	return NilType
}