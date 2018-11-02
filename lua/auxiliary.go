package lua

import (
	"syscall"
	"fmt"
)

var _ = fmt.Println

// ArgCheck checks whether cond is true. If it is not, raises an error with a
// standard message (see ArgError).
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_argcheck
func (state *State) ArgCheck(cond bool, arg int, msg string) {
	if !cond {
		state.ArgError(arg, msg)	
	}
}

// ArgError raises an error reporting a problem with argument arg of the Go
// function that called it, using a standard message that includes msg as
// a comment: ```bad argument #arg to 'funcname' (msg)```.
//
// This function never returns.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_argerror
func (state *State) ArgError(arg int, msg string) int {
	// TODO: to 'funcname'
	return state.Errorf("bad argument #%d (%s)", arg, msg)
}

// FileResult procudes the return values for file-related function in the standard library
// (io.open, os.rename, file:seek, etc.).
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_fileresult
func (state *State) FileResult(err error, filename string) int {
	if err == nil {
		state.Push(true)
		return 1
	}
	state.Push(nil)
	if errStr := err.Error(); filename != "" {
		state.Push(fmt.Sprintf("%s: %s", filename, errStr))
	} else {
		state.Push(errStr)
	}
	if errno, ok := err.(syscall.Errno); ok {
		state.Push(int(errno))
	} else {
		state.Push(0)
	}
	return 3
}

// If the registry already has the key name, return false. Otherwise, creates a new table to
// be used as a metatable for userdata, adds to this new table the pair __name = name, adds
// to the registry the pair [name] = new table, and returns true. The entry __name is used
// by some error-reporting functions.
//
// In both cases pushes onto the stack the final value associated with name in the registry.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_newmetatable
func (state *State) NewMetaTable(name string) bool {
	if mt := state.GetMetaTable(name); mt != NilType && mt != NoneType {
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
    return state.GetField(RegistryIndex, name)
}

// SetMetaTable sets the metatable of the object at the top of the stack as the metatable
// associated with name tname in the registry (see NewMetaTable).
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_setmetatable
func (state *State) SetMetaTable(name string) {
	state.GetMetaTable(name)
	state.SetMetaTableAt(-2)
}

// CallMeta calls a metamethod.
//
// If the object at index obj has a metatable and this metatable has a field event,
// this function calls this field passing the object as its only argument. In this
// case this function returns true and pushes onto the stack the value returned by
// the call. If there is no metatable or no metamethod, this function returns false
// (without pushing any value on the stack).
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_callmeta
func (state *State) CallMeta(index int, event string) bool {
	val := state.get(index)
	if meta := state.metafield(val, event); !IsNone(meta) {
        if cls, ok := meta.(*Closure); ok {
            state.frame().push(cls)
            state.frame().push(val)
            state.Call(1, 1)
			return true
        }
    }
	return false
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

// Pushes onto the stack a string identifying the current position
// of the control at level in the call stack. Typically this string
// has the following format:
//
// 		chunkname:currentline:
//
// Level 0 is the running function,
// level 1 is the function that called the running function, etc.
//
// This function is used to build a prefix for error messages.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_where
func (state *State) Where(level int) {
	// var debug Debug
	// if state.GetStack(level, &debug) {
	//		state.GetInfo("Sl", &debug)
	//		if debug.Line > 0 {
	//			state.Push(fmt.Sprintf("%s:%d: ", debug.Source, debug.Line))
	//			return
	//		}
	// }
	state.Push("")
}

// Errorf raises an error.
//
// The error message format is given by fmt plus any extra arguments, following
// the same rules of lua_pushfstring. It also adds at the beginning of the message
// the file name and the line number where the error occurred, if this information
// is available.
//
// This function never returns, but it is an idiom to use it in Go
// functions as return state.Errorf(fmt, args).
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_error
func (state *State) Errorf(format string, args ...interface{}) int {
	return state.errorf(format, args...)
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