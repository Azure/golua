package debug

import (
	"strings"
    "fmt"
    "os"
    "github.com/Azure/golua/lua"
)

var _ = fmt.Println
var _ = os.Exit

//
// Lua Standard Library -- debug
//

// Open open the Lua standard debug library. This library provides the functionality of the
// debug interface (§4.9) to Lua programs. You should exert care when using this library.
// 
// Several of its functions violate basic assumptions about Lua code (e.g., that variables local
// to a function cannot be accessed from outside; that userdata metatables cannot be changed by Lua
// code; that Lua programs do not crash) and therefore can compromise otherwise secure code.
// Moreover, some functions in this library may be slow.
//
// All functions in this library are provided inside the debug table. All functions that operate over
// a thread have an optional first argument which is the thread to operate over. The default is always
// the current thread.
func Open(state *lua.State) int {
	// Create 'debug' table.
    var debugFuncs = map[string]lua.Func{
        "debug":        lua.Func(dbgDebug),
		"gethook":      lua.Func(dbgGetHook),
		"getinfo":      lua.Func(dbgGetInfo),
		"getlocal":     lua.Func(dbgGetLocal),
		"getmetatable": lua.Func(dbgGetMetaTable),
		"getregistry":  lua.Func(dbgGetRegistry),
		"getupvalue":   lua.Func(dbgGetUpValue),
		"getuservalue": lua.Func(dbgGetUserValue),
		"sethook":      lua.Func(dbgSetHook),
		"setlocal":     lua.Func(dbgSetLocal),
		"setmetatable": lua.Func(dbgSetMetaTable),
		"setupvalue":   lua.Func(dbgSetUpvalue),
		"setuservalue": lua.Func(dbgSetUserValue),
		"traceback":    lua.Func(dbgTraceback),
		"upvalueid":    lua.Func(dbgUpValueID),
		"upvaluejoin":  lua.Func(dbgUpValueJoin),
    }
	state.NewTableSize(0, len(debugFuncs))
	state.SetFuncs(debugFuncs, 0)

 	// Return 'debug' table.
    return 1
}

// debug.debug ()
//
// Enters an interactive mode with the user, running each string that the user enters. Using simple
// commands and other debug facilities, the user can inspect global and local variables, change their
// values, evaluate expressions, and so on. A line containing only the word cont finishes this function,
// so that the caller continues its execution.
//
// Note that commands for debug.debug are not lexically nested within any function and so
// have no direct access to local variables.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-debug.debug
func dbgDebug(state *lua.State) int {
	unimplemented("debug: debug")
	return 0
}

// debug.getinfo ([thread,] f [, what])
//
// Returns a table with information about a function. You can give the function directly or you can give
// a number as the value of f, which means the function running at level f of the call stack of the given
// thread: level 0 is the current function (getinfo itself); level 1 is the function that called getinfo
// (except for tail calls, which do not count on the stack); and so on. If f is a number larger than the
// number of active functions, then getinfo returns nil.
//
// The returned table can contain all the fields returned by lua_getinfo, with the string what describing
// which fields to fill in. The default for what is to get all information available, except the table of
// valid lines. If present, the option 'f' adds a field named func with the function itself. If present,
// the option 'L' adds a field named activelines with the table of valid lines.
//
// For instance, the expression debug.getinfo(1,"n").name returns a name for the current function, if a
// reasonable name can be found, and the expression debug.getinfo(print) returns a table with all available
// information about the print function.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-debug.getinfo
func dbgGetInfo(state *lua.State) int {
	// (n) => name, kind
	// (S) => span, what, short, source
	// (l) => active
	// (u) => nups, params, vararg
	// (t) => tailcall

	thread, index := getThread(state)
	options := state.OptString(index+2, "flnStu")
	checkstack(state, thread, 3)

	var dbg lua.Debug

	if state.IsFunc(index+1) {
		options = fmt.Sprintf(">%s", options)
		state.PushIndex(index+1)
		state.XMove(thread, 1)
	} else {
		if err := state.GetStack(&dbg, int(state.CheckInt(index+1))); err != nil {
			state.Push(nil)
			return 1
		}
	}
	if err := state.GetInfo(&dbg, options); err != nil {
		panic(fmt.Errorf("bad argument #2 to 'getinfo' %v", err))
	}
	state.NewTable()
	if contains(options, 'S') {
		setFieldStr(state, "source", dbg.Source())
		setFieldStr(state, "short_src", dbg.ShortSrc())
		setFieldInt(state, "linedefined", dbg.LineDefined())
		setFieldInt(state, "lastlinedefined", dbg.LastLineDefined())
		setFieldStr(state, "what", dbg.What())
	}
	if contains(options, 'l') {
		setFieldInt(state, "currentline", dbg.CurrentLine())
	}
	if contains(options, 'u') {
		setFieldInt(state, "nups", dbg.NumUps())
		setFieldInt(state, "nparams", dbg.NumParams())
		setFieldBool(state, "isvararg", dbg.IsVararg())
	}
	if contains(options, 'n') {
		setFieldStr(state, "name", dbg.Name())
		setFieldStr(state, "namewhat", dbg.NameWhat())
	}
	if contains(options, 'T') {
		setFieldBool(state, "istailcall", dbg.IsTailCall())
	}
	if contains(options, 'L') {
		fmt.Println("debug.getinfo: option 'L': TODO")
	}
	if contains(options, 'f') {
		fmt.Println("debug.getinfo: option 'f': TODO")
	}
	fmt.Println()
	return 1
}

// debug.getmetatable (value)
//
// Returns the metatable of the given value or nil if it does not have a metatable.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-debug.getmetatable
func dbgGetMetaTable(state *lua.State) int {
	state.CheckAny(1)
	if !state.GetMetaTableAt(1) {
		state.Push(nil) // no metatable
	}
	return 1
}

// debug.setmetatable (value, table)
//
// Sets the metatable for the given value to the given table (which can be nil).
//
// Returns value.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-debug.setmetatable
func dbgSetMetaTable(state *lua.State) int {
	switch state.TypeAt(2) {
		case lua.NilType, lua.TableType:
			state.SetTop(2)
			state.SetMetaTableAt(1)
			return 1
		default:
			panic(fmt.Errorf("nil or table expected"))
	}
}

// debug.getregistry ()
//
// Returns the registry table (see §4.5).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-debug.getregistry
func dbgGetRegistry(state *lua.State) int {
	state.PushIndex(lua.RegistryIndex)
	return 1
}

// debug.getupvalue (f, up)
//
// This function returns the name and the value of the upvalue with index up of the
// function f. The function returns nil if there is no upvalue with the given index.
//
// Variable names starting with '(' (open parenthesis) represent variables with no
// known names (variables from chunks saved without debug information).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-debug.getupvalue
func dbgGetUpValue(state *lua.State) int {
	index := int(state.CheckInt(2))
	state.CheckType(1, lua.FuncType)
	ident := state.GetUpValue(1, index)
	if ident == "" {
		return 0
	}
	state.Push(ident)
	state.Insert(-2)
	return 2
}

// debug.setupvalue (f, up, value)
//
// This function assigns the value value to the upvalue with index up of the function f. The function returns
// nil if there is no upvalue with the given index. Otherwise, it returns the name of the upvalue.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-debug.setupvalue
func dbgSetUpvalue(state *lua.State) int {
	unimplemented("debug: setupvalue")
	return 0
}

// debug.gethook ([thread])
//
// Returns the current hook settings of the thread, as three values: the current hook function, the current
// hook mask, and the current hook count (as set by the debug.sethook function).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-debug.gethook
func dbgGetHook(state *lua.State) int {
	unimplemented("debug: gethook")
	return 0
}

// debug.sethook([thread,] hook, mask [, count])
//
// Sets the given function as a hook. The string mask and the number count describe
// when the hook will be called. The string mask may have any combination of the
// following characters, with the given meaning:
//
// 		'c': the hook is called every time Lua calls a function;
// 		'r': the hook is called every time Lua returns from a function;
// 		'l': the hook is called every time Lua enters a new line of code.
//
// Moreover, with a count different from zero, the hook is called also after every
// count instructions.
//
// When called without arguments, debug.sethook turns off the hook.
//
// When the hook is called, its first argument is a string describing the event
// that has triggered its call: "call" (or "tail call"), "return", "line", and
// "count". For line events, the hook also gets the new line number as its second
// parameter. Inside a hook, you can call getinfo with level 2 to get more information
// about the running function (level 0 is the getinfo function, and level 1 is the hook
// function).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-debug.sethook
func dbgSetHook(state *lua.State) int {
	unimplemented("debug: sethook")
	return 0
}

// debug.getlocal ([thread,] f, local)
//
// This function returns the name and the value of the local variable with index local of the function
// at level f of the stack. This function accesses not only explicit local variables, but also parameters,
// temporaries, etc.
//
// The first parameter or local variable has index 1, and so on, following the order that they are declared
// in the code, counting only the variables that are active in the current scope of the function. Negative
// indices refer to vararg arguments; -1 is the first vararg argument. The function returns nil if there is
// no variable with the given index, and raises an error when called with a level out of range. (You can call
// debug.getinfo to check whether the level is valid.)
//
// Variable names starting with '(' (open parenthesis) represent variables with no known names (internal variables
// such as loop control variables, and variables from chunks saved without debug information).
//
// The parameter f may also be a function. In that case, getlocal returns only the name of function parameters.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-debug.getlocal
func dbgGetLocal(state *lua.State) int {
	unimplemented("debug: getlocal")
	return 0
}

// debug.setlocal ([thread,] level, local, value)
//
// This function assigns the value value to the local variable with index local of the function at level level
// of the stack. The function returns nil if there is no local variable with the given index, and raises an error
// when called with a level out of range. (You can call getinfo to check whether the level is valid.) Otherwise,
// it returns the name of the local variable.
//
// See debug.getlocal for more information about variable indices and names.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-debug.setlocal
func dbgSetLocal(state *lua.State) int {
	unimplemented("debug: setlocal")
	return 0
}

// debug.getuservalue (u)
//
// Returns the Lua value associated to u. If u is not a full userdata, returns nil.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-debug.getuservalue
func dbgGetUserValue(state *lua.State) int {
	fmt.Println("debug.getuservalue")
	state.Debug(true)
	// if state.TypeAt(1) != lua.UserDataType {
	// 	state.Push(nil)
	// } else {
	// 	state.GetUserValue(1)
	// }
	return 1
}

//debug.setuservalue (udata, value)
//
// Sets the given value as the Lua value associated to the given udata. udata must be a full userdata.
//
// Returns udata.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-debug.setuservalue
func dbgSetUserValue(state *lua.State) int {
	fmt.Println("debug.setuservalue")
	state.Debug(true)
	// state.CheckType(1, lua.UserDataType)
	// state.CheckAny(2)
	// state.SetTop(2)
	// state.SetUserValue(1)
	return 1
}

// debug.traceback ([thread,] [message [, level]])
//
// If message is present but is neither a string nor nil, this function returns message without further processing.
// Otherwise, it returns a string with a traceback of the call stack. The optional message string is appended at the
// beginning of the traceback. An optional level number tells at which level to start the traceback (default is 1, the
// function calling traceback).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-debug.traceback
func dbgTraceback(state *lua.State) int {
	unimplemented("debug: traceback")
	return 0
}

// debug.upvalueid (f, n)
//
// Returns a unique identifier (as a light userdata) for the upvalue numbered n from the given function.
//
// These unique identifiers allow a program to check whether different closures share upvalues.
// Lua closures that share an upvalue (that is, that access a same external local variable) will
// return identical ids for those upvalue indices.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-debug.upvalueid
func dbgUpValueID(state *lua.State) int {
	upID := state.UpValueID(1, checkUpValue(state, 1, 2))
	state.Push(upID)
	return 1
}

// debug.upvaluejoin (f1, n1, f2, n2)
//
// Make the n1-th upvalue of the Lua closure f1 refer to the n2-th upvalue of the Lua closure f2.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-debug.upvaluejoin
func dbgUpValueJoin(state *lua.State) int {
	var (
		n1 = checkUpValue(state, 1, 2)
		n2 = checkUpValue(state, 3, 4)
	)
	// state.ArgCheck(!state.IsGoFunc(1), 1, "Lua function expected")
	// state.ArgCheck(!state.IsGoFunc(3), 3, "Lua function expected")
	if state.IsGoFunc(1) {
		panic(fmt.Errorf("bad argument #1 to 'upvaluejoin' (Lua function expected)"))
	}
	if state.IsGoFunc(3) {
		panic(fmt.Errorf("bad argument #3 to 'upvaluejoin' (Lua function expected)"))
	}
	state.UpValueJoin(1, n1, 3, n2)
	return 0
}

// hookMask converts a string mask (for 'sethook') into a bit mask.
func hookEvent(name string, count int64) (mask lua.HookEvent) {
	if strings.Index(name, "c") != -1 {
		mask |= lua.HookCall
	}
	if strings.Index(name, "r") != -1 {
		mask |= lua.HookRets
	}
	if strings.Index(name, "l") != -1 {
		mask |= lua.HookLine
	}
	if count > 0 {
		mask |= lua.HookCount
	}
	return mask
}

func unimplemented(msg string) { panic(fmt.Errorf(msg)) }

func contains(options string, option byte) bool {
	return strings.IndexByte(options, option) != -1
}

// Auxiliary function used by several library functions: check for
// an optional thread as function's first argument and set 'arg'
// with 1 if this argument is present (so that functions can skip
// it to access their other arguments).
func getThread(state *lua.State) (*lua.State, int) {
	if state.IsThread(1) {
		return state.ToThread(1), 1
	}
	return state, 0
}

// convenience functions.

func checkstack(l1, l2 *lua.State, n int) {
	fmt.Println("debug: checkstack: TODO")
}

func setFieldStr(state *lua.State, key, value string) {
	state.Push(value)
	state.SetField(-2, key)
}

func setFieldInt(state *lua.State, key string, value int) {
	state.Push(value)
	state.SetField(-2, key)
}

func setFieldBool(state *lua.State, key string, value bool) {
	state.Push(value)
	state.SetField(-2, key)
}

// checkUpValue checks whether a given upvalue from a given closure exists and returns its index.
func checkUpValue(state *lua.State, function, index int) (up int) {
    state.CheckType(1, lua.FuncType)
    up = int(state.CheckInt(index))
	if state.GetUpValue(function, up) == "" {
		panic(fmt.Errorf("invalid upvalue index"))
	}
	return up
}