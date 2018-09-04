package stdlib

import (
	"fmt"
	"github.com/Azure/golua/lua"
)

var _ = fmt.Println

//
// Lua Standard Library -- Base
//

// collectgarbage
// xpcall
func OpenBase(state *lua.State) uint32 {
	var funcs = map[string]lua.Func{
		"assert": 		  nil,
		"dofile": 		  nil,
		"error": 		  nil,
		"getmetatable":   nil,
		"ipairs": 		  nil,
		"loadfile": 	  nil,
		"load": 		  nil,
		"next": 		  nil,
		"pairs": 		  nil,
		"pcall": 		  nil,
		"print": 		  nil,
		"rawequal": 	  nil,
		"rawlen": 		  nil,
		"rawget": 	      nil,
		"rawset": 	      nil,
		"select": 	      nil,
		"setmetatable":   nil,
		"tonumber": 	  nil,
		"tostring": 	  nil,
		"type":     	  nil,
		"xpcall": 		  nil,
	}
	// globals := state.Globals()
	// globals.SetFuncs(funcs)
	// globals.SetValue("_G", globals)
	// globals.SetValue("_VERSION", lua.Version)
	
	// table := state.NewTable()
	// table.SetFuncs(funcs)
	// state.SetGlobal("_G", table)
	// state.SetGlobal("_VERSION", lua.Version)

	table := state.SetFuncs(state.Globals(), funcs)
	state.Push(table)
	state.SetField(-3, "_G")
	state.PushString(lua.Version)
	state.SetField(-3, "_VERSION")
	return 1
}

// assert(v, [, message])
//
// Calls error if the value of its argument v is false (i.e., nil or false);
// otherwise, returns all its arguments. In case of error, message is the error
// object; when absent, it defaults to "assertion failed!"
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-assert
// func assert_(state *lua.State) uint32 {
// 	if state.ToBool(1) {
// 		return state.GetTop()
// 	}
// 	state.Errorf("")
// }

// dofile([filename])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-dofile
//func dofile_(state *lua.State) uint32 {}

// error(message [, level])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-error
//func error_(state *lua.State) uint32 {}

// getmetatable(object)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-getmetatable
//func getmetatable_(state *lua.State) uint32 {}

// ipairs(t)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-ipairs
//func ipairs_(state *lua.State) uint32 {}

// loadfile([filename [, mode [, env]]])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-loadfile
//func loadfile_(state *lua.State) uint32 {}

// load(chunk [, chunkname [, mode [, env]]])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-load
//func load_(state *lua.State) uint32 {}

// next(table [, index])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-next
//func next_(state *lua.State) uint32 {}

// pairs(t)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-pairs
//func pairs_(state *lua.State) uint32 {}

// pcall(f [, arg1, ...])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-pcall
//func pcall_(state *lua.State) uint32 {}

// print(...)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-print
// func print_(state *lua.State) uint32 {
// 	top := state.GetTop() // # args
// 	state.GetGlobal("tostring")
// 	for i := 1; i <= top; i++ {
// 		state.PushValue(-1)
// 		state.PushValue(i)
// 		state.Call(1, 1)

// 		str, ok := state.ToString(-1)
// 		if !ok {
// 			panic("'tostring' must return a string to 'print'")
// 		}
// 		fmt.Print(str)
// 		state.Pop(1)
// 	}
// 	fmt.Println()
// 	return 0
// }

// rawequal(v1, v2)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-rawequal
//func rawequal_(state *lua.State) uint32 {}

// rawlen(v)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-rawlen
//func rawlen_(state *lua.State) uint32 {}

// rawget(table, index)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-rawget
//func rawget_(state *lua.State) uint32 {}

// rawset(table, index, value)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-rawset
//func rawset_(state *lua.State) uint32 {}

// select(index, ...)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-select
//func select_(state *lua.State) uint32 {}

// setmetatable(table, metatable)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-setmetatable
//func setmetatable_(state *lua.State) uint32 {}

// tonumber(e [, base])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-tonumber
//func tonumber_(state *lua.State) uint32 {}

// tostring(v)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-tostring
//func tostring_(state *lua.State) uint32 {}

// type(v)
// 
// Returns the type of its only argument, coded as a string.
//
// The possible results of this function are "nil" (a string,
// not the value nil), "number", "string", "boolean", "table",
// "function", "thread", and "userdata".
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-type
//func type_(state *lua.State) uint32 {}