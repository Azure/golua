package lua

import (
    "fmt"
    "os"

    "github.com/Azure/golua/lua/binary"
)

var errNonBinaryChunk = fmt.Errorf("lua: execute: non-binary files not yet supported")

var _ = os.Exit

// Dump dumps a function as a binary chunk. Receives a Lua function on top of the
// stack and produces a binary chunk that, if loaded again, results in a function
// equivalent to the one dumped.
//
// If strip is true, the binary representation may not include all debug information
// about the function, to save space.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_dump
func (state *State) Dump(strip bool) []byte {
    if cls, ok := state.get(-1).(*Closure); ok {
        if cls.isLua() {
            return binary.Dump(cls.binary, strip)
        }
    }
    return nil
}

// SetUpValue sets the value of a closure's upvalue. It assigns the value at the top
// of the stack to the upvalue at index and returns its name. It also pops the value
// from teh stack.
//
// Otherwise returns "" (and pops nothing) when upIndex is greater than the number of
// upvalues.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_setupvalue
func (state *State) SetUpValue(fnIndex, upIndex int) (name string) {
    if cls, ok := state.get(fnIndex).(*Closure); ok {
        if upAt := upIndex - 1; upAt < len(cls.upvals) {
            upvalue := state.frame().pop()
            cls.setUp(upAt, upvalue)
            name = cls.upName(upAt)
        }
    }
    return
}

// Load loads a Lua chunk wihtout running it. If there are no errors, Load pushes
// the compiled chunk as a Lua function on top of the stack, otherwise nothing is
// pushed and the error is returned.
//
// If the resulting function has upvalues, its first upvalue is set to the value of
// the global environment stored at index LUA_RIDX_GLOBALS in the registry (see §4.5).
// When loading main chunks, this upvalue will be the _ENV variable (see §2.2). Other
// upvalues are initialized with nil.
func (state *State) LoadChunk(filename string, source interface{}, mode Mode) error {
    cls, err := state.load(filename, source)
    if err != nil {
        return err
    }
    state.frame().push(cls)
    return nil
}

// Exec loads and runs a Lua chunk returning the result (if any) or an error (if any).
//
// The Lua chunk may be provided via the filename of the source file, or via the
// source parameter.
//
// If source != nil, Do loads the source from source and the filename is only used
// recording position information.
func (state *State) ExecChunk(filename string, source interface{}, mode Mode) error {
    if err := state.LoadChunk(filename, source, mode); err != nil {
        return err
	}
	state.Call(0, -1)
    return nil
}

// Register sets the Go function fn as the new value of global name.
func (state *State) Register(name string, fn Func) {
    state.Push(newGoClosure(fn, 0))
    state.SetGlobal(name)
}

// AtPanic sets a new panic function and returns the old one (see §4.6).
//
// See https://www.lua.org/manual/5.3/manual.html#lua_atpanic
func (state *State) AtPanic(panicFn Func) Func {
    prevFn := state.global.panicFn
    state.global.panicFn = panicFn
    return prevFn
}

// Status returns the status of the thread.
// 
// The status can be 0 (LUA_OK) for a normal thread, an error code if the thread finished
// the execution of a lua_resume with an error, or LUA_YIELD if the thread is suspended.
// 
// You can only call functions in threads with status LUA_OK. You can resume threads with
// status LUA_OK (to start a new coroutine) or LUA_YIELD (to resume a coroutine).
//
// See https://www.lua.org/manual/5.3/manual.html#lua_status
func (state *State) Status() ThreadStatus { unimplemented("Status"); return ThreadError }

// Generates a Lua error, using the value at the top of the stack as the error object.
// 
// This function does a long jump, and therefore never returns (see luaL_error).
//
// See https://www.lua.org/manual/5.3/manual.html#lua_error
func (state *State) Error() int { return state.errorf("%v", state.frame().pop()) }

// Destroys all objects in the given Lua state (calling the corresponding garbage-collection metamethods, if any) and
// frees all dynamic memory used by this state. In several platforms, you may not need to call this function, because
// all resources are naturally released when the host program ends. On the other hand, long-running programs that create
// multiple states, such as daemons or web servers, will probably need to close states as soon as they are not needed.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_close
func (state *State) Close() { /* TODO */ }

// Returns the address of the version number (a C static variable) stored in the Lua core. When called with a valid lua_State,
// returns the address of the version used to create that state. When called with NULL, returns the address of the version
// running the call.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_version
func (state *State) Version() *float64 { unimplemented("Version"); return nil }

// Performs an arithmetic or bitwise operation over the two values (or one, in the case of negations) at the top of the
// stack, with the value at the top being the second operand, pops these values, and pushes the result of the operation.
// The function follows the semantics of the corresponding Lua operator (that is, it may call metamethods).
// 
// The value of op must be one of the following constants:
//      * LUA_OPADD: performs addition (+)
//      * LUA_OPSUB: performs subtraction (-)
//      * LUA_OPMUL: performs multiplication (*)
//      * LUA_OPDIV: performs float division (/)
//      * LUA_OPIDIV: performs floor division (//)
//      * LUA_OPMOD: performs modulo (%)
//      * LUA_OPPOW: performs exponentiation (^)
//      * LUA_OPUNM: performs mathematical negation (unary -)
//      * LUA_OPBNOT: performs bitwise NOT (~)
//      * LUA_OPBAND: performs bitwise AND (&)
//      * LUA_OPBOR: performs bitwise OR (|)
//      * LUA_OPBXOR: performs bitwise exclusive OR (~)
//      * LUA_OPSHL: performs left shift (<<)
//      * LUA_OPSHR: performs right shift (>>)
//
// See https://www.lua.org/manual/5.3/manual.html#lua_arith
func (state *State) Arith(op Op) {
	y := state.frame().pop()
	x := state.frame().pop()
	state.frame().push(state.arith(op, x, y))
}

// Concatenates the n values at the top of the stack, pops them, and leaves the result at the top. If n is 1, the result
// is the single value on the stack (that is, the function does nothing); if n is 0, the result is the empty string.
// Concatenation is performed following the usual semantics of Lua (see https://www.lua.org/manual/5.3/manual.html#3.4.6).
//
// See https://www.lua.org/manual/5.3/manual.html#lua_concat
func (state *State) Concat(n int) {
    if n > 1 {
        values := state.frame().popN(n)
        result := state.concat(values)
        state.frame().push(result)
    }
}

// Returns the length of the value at the given index. It is equivalent to the ‘\#’ operator in Lua
// (see [Lua 5.3 Reference Manual](https://www.lua.org/manual/5.3/manual.html#3.4.7)) and may trigger
// a metamethod for the “length” event (see [Lua 5.3 Reference Manual](https://www.lua.org/manual/5.3/manual.html#2.4)).
// The result is pushed on the stack.
func (state *State) Length(index int) int { unimplemented("Length"); return 0 }

// Compares two Lua values. Returns 1 if the value at index index1 satisfies op
// when compared with the value at index index2, following the semantics of the
// corresponding Lua operator (that is, it may call metamethods).
//
// Otherwise returns false. Also returns false if any of the indices are invalid.
// 
// The value of op must be one of the following constants:
//      * LUA_OPEQ: compares for equality (==)
//      * LUA_OPLT: compares for less than (<)
//      * LUA_OPLE: compares for less or equal (<=)
func (state *State) Compare(op Op, i1, i2 int) bool {
   return state.compare(op, state.get(i1), state.get(i2), false)
}

// Pushes onto the stack the value of the global name.
// 
// Returns the type of that value.
func (state *State) GetGlobal(name string) Type {
    val := state.gettable(state.globals(), String(name), false)
    state.frame().push(val)
    return val.Type()
}

// Pops a value from the stack and sets it as the new value of global name.
//func (state *State) SetGlobal(name string, value Value) {
func (state *State) SetGlobal(name string) {
    state.settable(state.globals(), String(name), state.Pop(), false)
}

// Creates a new empty table and pushes it onto the stack. Parameter narr is a hint for how
// many elements the table will have as a sequence; parameter nrec is a hint for how many other
// elements the table will have. Lua may use these hints to preallocate memory for the new table.
// This preallocation is useful for performance when you know in advance how many elements the table
// will have. Otherwise you can use the function lua_newtable.
func (state *State) NewTableSize(narr, nrec int) {
    state.frame().push(newTable(state, narr, nrec))
}

// Creates a new empty table and pushes it onto the stack.
// 
// It is equivalent to lua_createtable(L, 0, 0).
func (state *State) NewTable() { state.NewTableSize(0, 0) }

// Next pops a key from the stack, and pushes a key-value pair from the table at the given index
// (the "next" pair after the given key). If there are no more elements in the table, then Next
// returns false (and pushes nothing).
//
// While traversing a table, do not call lua.ToString directly on a key, unless you know that the
// key is actually a string. Recall that lua.ToString may change the value at the given index;
// this confuses the next call to lua.Next.
//
// See function "next" for the caveats of modifying the table during its traversal.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_next
func (state *State) Next(index int) (more bool) {
    tbl, ok := state.get(index).(*table)
    if !ok {
        state.errorf("table expected")
    }
    var k, v Value
    if k, v, more = tbl.next(state.frame().pop()); more {
        // fmt.Printf("next: key=%v, value=%v (more = %t)\n", k, v, more)
        state.frame().push(k)
        state.frame().push(v)
    }
    return
}

// If the value at the given index has a metatable, the function pushes that metatable onto
// the stack and returns 1. Otherwise, the function returns 0 and pushes nothing on the stack.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_getmetatable
func (state *State) GetMetaTableAt(index int) bool {
    if meta := state.getmetatable(state.get(index), true); !IsNone(meta) {
        state.Push(meta)
        return true
    }
    return false
}

// Pops a table from the stack and sets it as the new metatable for the value at the given index.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_setmetatable
func (state *State) SetMetaTableAt(index int) { state.setmetatable(state.get(index), state.frame().pop()) }

// GetTable pushes onto the stack the value t[k], where t is the value at the given index
// and k is the value at the top of the stack.
//
// This function pops the key from the stack, pushing the resulting value in its place.
// As in Lua, this function may trigger a metamethod for the "index" event (see §2.4).
//
// See https://www.lua.org/manual/5.3/manual.html#lua_gettable
func (state *State) GetTable(index int) Type {
    var (
        key = state.frame().pop()
        obj = state.get(index)
    )
    val := state.gettable(obj, key, false)
    state.frame().push(val)
    return val.Type()
}

// SetTable does the equivalent to t[k] = v, where t is the value at the given index, v is the
// value at the top of the stack, and k is the value just below the top.
//
// This function pops both the key and the value from the stack. As in Lua, this function may
// trigger a metamethod for the "newindex" event (see §2.4).
//
// See https://www.lua.org/manual/5.3/manual.html#lua_settable
func (state *State) SetTable(index int) {
    var (
        val = state.frame().pop()
        key = state.frame().pop()
        obj = state.get(index)
    )
    state.settable(obj, key, val, false)
}

// Pushes onto the stack the value t[k], where t is the value at the given index.
//
// As in Lua, this function may trigger a metamethod for the "index" event (see §2.4).
//
// Returns the type of the pushed value.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_getfield
func (state *State) GetField(index int, field string) Type {
    v := state.gettable(state.get(index), String(field), false)
    state.frame().push(v)
    return v.Type()
}

// Does the equivalent to t[k] = v, where t is the value at the given index and v is the
// value at the top of the stack.
// 
// This function pops the value from the stack.
//
// As in Lua, this function may trigger a metamethod for the "newindex" event (see §2.4).
//
// See https://www.lua.org/manual/5.3/manual.html#lua_setfield
func (state *State) SetField(index int, field string) {
    obj := state.get(index)
    key := String(field)
    val := state.frame().pop()
    state.settable(obj, key, val, false)
}

// Pushes onto the stack the value t[i], where t is the value at the given index.
//
// As in Lua, this function may trigger a metamethod for the "index" event (see §2.4).
// 
// Returns the type of the pushed value.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_geti
func (state *State) GetIndex(index int, entry int64) Type {
    obj := state.get(index)
    key := Int(entry)
    val := state.gettable(obj, key, false)
    state.frame().push(val)
    return val.Type()
}

// Does the equivalent to t[n] = v, where t is the value at the given index and v is the
// value at the top of the stack.
// 
// This function pops the value from the stack.
//
// As in Lua, this function may trigger a metamethod for the "newindex" event (see §2.4).
//
// See https://www.lua.org/manual/5.3/manual.html#lua_seti
func (state *State) SetIndex(index int, entry int64) {
    tbl := state.get(index)
    key := Int(entry)
    val := state.frame().pop()
    state.settable(tbl, key, val, false)
}

// Similar to lua_gettable, but does a raw access (i.e., without metamethods).
//
// See https://www.lua.org/manual/5.3/manual.html#lua_rawget
func (state *State) RawGet(index int) Type {
    var (
        key = state.frame().pop()
        obj = state.get(index)
	)
    val := state.gettable(obj, key, true)
    state.frame().push(val)
    return val.Type()
}

// Similar to lua_settable, but does a raw assignment (i.e., without metamethods).
//
// See https://www.lua.org/manual/5.3/manual.html#lua_rawset
func (state *State) RawSet(index int) {
    var (
        val = state.frame().pop()
        key = state.frame().pop()
        obj = state.get(index)
    )
    state.settable(obj, key, val, true)
}

// RawLen returns the raw "length" of the value at the given index: for strings, this
// is the string length; for tables, this is the result of the length operator ('#')
// with no metamethods; for userdata, this is the size of the block of memory allocated
// for the userdata; for other values, it is 0.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_rawlen
func (state *State) RawLen(index int) int {
    switch v := state.get(index).(type) {
        case String:
            return len(v)
        case *table:
            return len(v.list)
    }
    return 0
}

// Pushes onto the stack the value t[n], where t is the table at the given index.
//
// The access is raw, that is, it does not invoke the __index metamethod.
// 
// Returns the type of the pushed value.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_rawgeti
func (state *State) RawGetIndex(index, entry int) Type {
    var (
        obj = state.get(index)
        key = Int(entry)
        val = state.gettable(obj, key, true)
    )
    state.frame().push(val)
    return val.Type()
}

// Does the equivalent of t[i] = v, where t is the table at the given index and
// v is the value at the top of the stack.
// 
// This function pops the value from the stack.
//
// The assignment is raw, that is, it does not invoke the __newindex metamethod.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_rawseti
func (state *State) RawSetIndex(index, entry int) {
    tbl, ok := state.get(index).(*table)
    if !ok {
        state.errorf("table expected")
        return
    }
    tbl.setInt(int64(entry), state.Pop())
}

// Pushes onto the stack the value t[k], where t is the table at the given index and
// k is the pointer p represented as a light userdata.
//
// The access is raw; that is, it does not invoke the __index metamethod.
// 
// Returns the type of the pushed value.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_rawgetp
func (state *State) RawGetPtr(index int, udata *Object) Type {
    unimplemented("RawGetPtr")
    return NilType
}

// Does the equivalent of t[p] = v, where t is the table at the given index, p is
// encoded as a light userdata, and v is the value at the top of the stack.
// 
// This function pops the value from the stack.
//
// The assignment is raw, that is, it does not invoke __newindex metamethod.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_rawsetp
func (state *State) RawSetPtr(index int, udata *Object) {
    unimplemented("RawSetPtr")
}

// Returns true if the two values in indices index1 and index2 are primitively equal
// (that is, without calling the __eq metamethod).
//
// Otherwise returns false. or if any of the indices are not valid.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_rawequal
func (state *State) RawEqual(i1, i2 int) bool {
    if !state.isValid(i1) || !state.isValid(i2) {
        return false
    }
    return state.compare(OpEq, state.get(i1), state.get(i2), true)
}

// PCall calls a function in protected mode.
//
// Both nargs and nresults have the same meaning as in lua_call. If there are no errors during the call,
// lua_pcall behaves exactly like lua_call. However, if there is any error, lua_pcall catches it, pushes
// a single value on the stack (the error object), and returns an error code. Like lua_call, lua_pcall
// always removes the function and its arguments from the stack.
//
// If msgh is 0, then the error object returned on the stack is exactly the original error object.
// Otherwise, msgh is the stack index of a message handler. (This index cannot be a pseudo-index.)
// In case of runtime errors, this function will be called with the error object and its return value
// will be the object returned on the stack by lua_pcall.
//
// Typically, the message handler is used to add more debug information to the error object, such as a
// stack traceback. Such information cannot be gathered after the return of lua_pcall, since by then the
// stack has unwound.
//
// The lua_pcall function returns one of the following constants (defined in lua.h):
//
//  LUA_OK (0): success.
//  LUA_ERRRUN: a runtime error.
//  LUA_ERRMEM: memory allocation error. For such errors, Lua does not call the message handler.
//  LUA_ERRERR: error while running the message handler.
//  LUA_ERRGCMM: error while running a __gc metamethod. For such errors, Lua does not call
//               the message handler (as this kind of error typically has no relation with
//               the function being called).
//
// See https://www.lua.org/manual/5.3/manual.html#lua_pcall
func (state *State) PCall(args, rets, msgh int) (err error) {
    defer func(err *error) {
        if r := recover(); r != nil {
            if e, ok := r.(error); ok {
                *err = e
            }
        }
    }(&err)
    state.Call(args, rets)
    return
}

// Call calls a function.
//
// To call a function you must use the following protocol: first, the function to be called is pushed onto the stack;
// then, the arguments to the function are pushed in direct order; that is, the first argument is pushed first.
// Finally you call lua_call; nargs is the number of arguments that you pushed onto the stack. All arguments and the
// function value are popped from the stack when the function is called. The function results are pushed onto the stack
// when the function returns. The number of results is adjusted to nresults, unless nresults is LUA_MULTRET.
// In this case, all results from the function are pushed; Lua takes care that the returned values fit into the stack
// space, but it does not ensure any extra space in the stack. The function results are pushed onto the stack in direct
// order (the first result is pushed first), so that after the call the last result is on the top of the stack.
//
// Note that the code above is balanced: at its end, the stack is back to its original configuration.
// This is considered good programming practice.
func (state *State) Call(args, rets int) {
    //checkNumStack(state, argN + 1)
    //checkResults(state, argN, retN)
    var (
        funcID = state.frame().absindex(-(args+1))
        value  = state.frame().get(funcID-1)
        c, ok  = value.(*Closure)
	)

	// state.Logf("call (func @ %d) %v (# args = %d, # rets = %d)\n", funcID, value, args, rets)
    
    if !ok {
		if !tryMetaCall(state, value, funcID, args, rets) {
			state.errorf("attempt to call a %s value @ %d (%T)\n%v\n", value.Type(), funcID, value, state.frame().locals)
		}
    } else {
		state.call(&Frame{closure: c, fnID: funcID, rets: rets})
	}
}

// Registers all functions in the array l (see luaL_Reg) into the table on the top of the stack (below optional
// upvalues, see next).
//
// When nup is not zero, all functions are created sharing nup upvalues, which must be previously pushed on the
// stack on top of the library table. These values are popped from the stack after the registration.
func (state *State) SetFuncs(funcs map[string]Func, nups uint8) {
    upvalues := state.frame().popN(int(nups))
    for name, fn := range funcs {
        state.frame().pushN(upvalues)
        state.PushClosure(fn, nups)
        state.SetField(-2, name)
    }
}

// Require loads the module.
//
// If modname is not already present in package.loaded, calls function loader with string modname
// as an argument and sets the call result in package.loaded[modname], as if that function has been
// called through require.
//
// If global is true, also stores the module into global modname.
//
// Leaves a copy of the module on the stack.
//
// See https://www.lua.org/manual/5.3/manual.html#luaL_requiref
func (state *State) Require(module string, loader Func, global bool) {
    state.Logf("require %q (global = %t)", module, global)
    state.GetSubTable(RegistryIndex, LoadedKey)
    state.GetField(-1, module)        // LOADED[module]

    if !Truth(state.get(-1)) {        // package not already loaded?
        state.Pop()                   // remove field
        state.PushClosure(loader, 0)  // push loader
        state.Push(module)            // argument to open function
        state.Call(1, 1)              // call 'loader' to open module
        state.PushIndex(-1)           // make copy of module (call result)
        state.SetField(-3, module)    // LOADED[modname] = module
    }

    state.Remove(-2) // remove LOADED table

    // If global, load copy of module and set _G[modname] = module
    if global {
        state.PushIndex(-1)
        state.SetGlobal(module)
    }
}

// Preload preloads a Lua module into the package.preload table.
func (state *State) Preload(module string, loader Func) {
	state.Logf("preload %q", module)

	state.GetSubTable(RegistryIndex, PreloadKey)
	state.GetField(-1, module) // PRELOAD

	if !Truth(state.get(-1)) { 		 // package not already preloaded?
		state.Pop() 		  		// remove field
		state.PushClosure(loader, 0) // push loader
		state.SetField(-2, module)   // PRELOAD[module] = loader
	}
	state.Remove(-2) // remove PRELOAD table
}

func (state *State) Main(args ...string) error {
	return state.safely(func() error { // pmain
		defer state.Close()
		// TODO: open stdlib
		return state.ExecFile(args[0])
	})
}

//
// TODO: remove me
//
func unimplemented(msg string) { panic(fmt.Errorf(msg)) }