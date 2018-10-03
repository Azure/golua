package lua

import "fmt"

var _ = fmt.Println

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
	return NoneType
}

func (state *State) valueAt(index int) Value {
	if state.isValid(index) {
		return state.get(index)
	}
	return None
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