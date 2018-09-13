package lua

import (
    "strings"
    "fmt"
    "io"
    "os"
)

var errNonBinaryChunk = fmt.Errorf("lua: execute: non-binary files not yet supported")

var _ = os.Exit

// DumpStack writes the Lua thread's current stack frame to w.
func (state *State) DumpStack(w io.Writer) {
    if fr := state.frame(); fr != nil {
        for i := fr.gettop() - 1; i >= 0; i-- {
            fmt.Fprintf(w, "[%d] %v\n", i+1, fr.locals[i])
        }
    }
}

// EnsureStack ensures that the stack has space for at least n extra slots (that is, that you can
// safely push up to n values into it). It returns false if it cannot fulfill the request,
// either because it would cause the stack to be larger than a fixed maximum size (typically
// at least several thousand elements) or because it cannot allocate memory for the extra space.
//
// This function never shrinks the stack; if the stack already has space for the extra slots, it
// is left unchanged.
func (state *State) EnsureStack(needed int) bool {
    return state.frame().checkstack(needed)
}

// AbsIndex converts the acceptable index idx into an equivalent absolute index;
// that is, one that does not depend on the stack top).
func (state *State) AbsIndex(index int) int {
    return state.frame().absindex(index)
}

// Safely executes the function fn returning any errors recovered by the Lua
// runtime
func (ls *State) Safely(fn func()) (err error) {
    defer ls.recover(&err)
    fn()
    return
}

// String returns a printable string of the current executing thread state.
func (ls *State) String() string {
    var w strings.Builder
    ls.Dump(&w)
    return w.String()
}

// Pushes a new Go closure onto the stack.
//
// When a Go function is created, it is possible to associate some values with it,
// thus creating a Go closure (see §4.4); these values are then accessible to the
// function whenever it is called. To associate values with a Go function, first
// these values must be pushed onto the stack (when there are multiple values, the
// first value is pushed first). Then PushClosure is called to create and push the
// Go function onto the stack, with the argument n telling how many values will be
// associated with the function. PushClosure also pops these values from the stack.
// 
// The maximum value for nups is 255.
//
// When nups is zero, this function creates a light Go function, which is just a pointer
// to the Go function. In that case, it never raises a memory error.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_pushcclosure
func (state *State) PushClosure(fn Func, nups uint8) {
    cls := newGoClosure(fn, int(nups))
    for nups > 0 {
        up := state.Pop()
        cls.upvals[nups-1] = &up
        nups--
    }
    state.Push(cls)
}

// PushIndex pushes a copy of the element at the given index onto the stack.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_pushvalue
func (state *State) PushIndex(index int) {
    state.frame().push(state.get(index))
}

// Removes the element at the given valid index, shifting down the elements
// above this index to fill the gap.
//
// This function cannot be called with a pseudo-index, because a pseudo-index
// is not an actual stack position.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_remove
func (state *State) Remove(index int) { state.frame().remove(index) }

// Insert moves the top element into the given valid index, shifting up the elements
// above this index to open space.
//
// This function cannot be called with a pseudo-index, because a pseudo-index is not
// an actual stack position.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_insert
func (state *State) Insert(index int) { state.frame().rotate(index, 1) }

// Value returns the Lua value at the valid index.
func (state *State) Value(index int) Value { return state.get(index) }

// Top returns the index of the top element in the stack.
//
// Because indices start at 1, this result is equal to the number
// of elements in the stack; in particular, 0 means an empty stack.
func (state *State) Top() int { return state.frame().gettop() }

// Push pushes any value onto the stack first boxing it by the equivalent
// Lua value, returning its position in the frame's local stack (top - 1).
func (state *State) Push(v interface{}) int {
    state.frame().push(ValueOf(v))
    return state.Top() - 1
}

// PopN pops the top n values from the Lua thread's current frame stack.
func (state *State) PopN(n int) []Value { return state.frame().popN(n) }

// Pop pops the top value from the Lua thread's current frame stack.
func (state *State) Pop() Value { return state.frame().pop() }

// Dump writes a trace of the Lua thread's call stack to w.
func (state *State) Dump(w io.Writer) { state.traceback(w) }

// PushGlobals pushes onto the stack the globals table.
func (state *State) Globals() *Table {
    return state.global.registry.getInt(GlobalsIndex).(*Table)
}

// Load loads a Lua chunk wihtout running it. If there are no errors, Load pushes
// the compiled chunk as a Lua function on top of the stack, otherwise nothing is
// pushed and the error is returned.
//
// If the resulting function has upvalues, its first upvalue is set to the value of
// the global environment stored at index LUA_RIDX_GLOBALS in the registry (see §4.5).
// When loading main chunks, this upvalue will be the _ENV variable (see §2.2). Other
// upvalues are initialized with nil.
func (state *State) Load(filename string, source interface{}) error {
    cls, err := state.load(filename, source)
    if err != nil {
        return err
    }
    state.frame().push(cls)
    return nil
}

// Register sets the Go function fn as the new value of global name.
// state.PushFunc(fn)
// state.SetGlobal(name)
func (state *State) Register(name string, fn Func) { unimplemented("Register") }

// AtPanic sets a new panic function and returns the old one (see §4.6).
func (state *State) AtPanic(fn Func) Func { unimplemented("AtPanic"); return nil }

// Status returns the status of the thread.
// 
// The status can be 0 (LUA_OK) for a normal thread, an error code if the thread finished
// the execution of a lua_resume with an error, or LUA_YIELD if the thread is suspended.
// 
// You can only call functions in threads with status LUA_OK. You can resume threads with
// status LUA_OK (to start a new coroutine) or LUA_YIELD (to resume a coroutine).
func (state *State) Status() ThreadStatus { unimplemented("Status"); return ThreadError }

// Generates a Lua error, using the value at the top of the stack as the error object.
// 
// This function does a long jump, and therefore never returns (see luaL_error).
func (state *State) Error() { unimplemented("Error") }

// Destroys all objects in the given Lua state (calling the corresponding garbage-collection metamethods, if any) and
// frees all dynamic memory used by this state. In several platforms, you may not need to call this function, because
// all resources are naturally released when the host program ends. On the other hand, long-running programs that create
// multiple states, such as daemons or web servers, will probably need to close states as soon as they are not needed.
func (state *State) Close() { unimplemented("Close") }

// Returns the address of the version number (a C static variable) stored in the Lua core. When called with a valid lua_State,
// returns the address of the version used to create that state. When called with NULL, returns the address of the version
// running the call.
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
func (state *State) Arith(op Op) { unimplemented("Arith") }

// Concatenates the n values at the top of the stack, pops them, and leaves the result at the top. If n is 1, the result
// is the single value on the stack (that is, the function does nothing); if n is 0, the result is the empty string.
// Concatenation is performed following the usual semantics of Lua (see https://www.lua.org/manual/5.3/manual.html#3.4.6).
func (state *State) Concat(n int) { unimplemented("Concat") }

// Returns the length of the value at the given index. It is equivalent to the ‘\#’ operator in Lua
// (see [Lua 5.3 Reference Manual](https://www.lua.org/manual/5.3/manual.html#3.4.7)) and may trigger
// a metamethod for the “length” event (see [Lua 5.3 Reference Manual](https://www.lua.org/manual/5.3/manual.html#2.4)).
//The result is pushed on the stack.
func (state *State) Length(index int) int { unimplemented("Length"); return 0 }

// Compares two Lua values. Returns 1 if the value at index index1 satisfies op when compared with the value at index index2, following the semantics of the corresponding Lua operator (that is, it may call metamethods). Otherwise returns 0. Also returns 0 if any of the indices is not valid.
// 
// The value of op must be one of the following constants:
//      * LUA_OPEQ: compares for equality (==)
//      * LUA_OPLT: compares for less than (<)
//      * LUA_OPLE: compares for less or equal (<=)
func (state *State) Compare(op Op, i1, i2 int) bool { unimplemented("Compare"); return false }

// Pushes onto the stack the value of the global name.
// 
// Returns the type of that value.
func (state *State) GetGlobal(name string) Type {
    obj := state.global.registry.getInt(GlobalsIndex)
    val := state.gettable(obj, String(name), 1)
    return val.Type()
}

// Pops a value from the stack and sets it as the new value of global name.
func (state *State) SetGlobal(name string) {
    key := String(name)
    val := state.Pop()
    state.settable(state.Globals(), key, val, 1)
}

// Creates a new empty table and pushes it onto the stack. Parameter narr is a hint for how
// many elements the table will have as a sequence; parameter nrec is a hint for how many other
// elements the table will have. Lua may use these hints to preallocate memory for the new table.
// This preallocation is useful for performance when you know in advance how many elements the table
// will have. Otherwise you can use the function lua_newtable.
func (state *State) NewTableSize(narr, nrec int) {
    tbl := &Table{newTable(state, narr, nrec)}
    state.frame().push(tbl)
}

// Creates a new empty table and pushes it onto the stack.
// 
// It is equivalent to lua_createtable(L, 0, 0).
func (state *State) NewTable() { state.NewTableSize(0, 0) }

// GetSubTable ensures that stack[index][field] has a table and pushes
// that table onto the stack.
//
// Returns true if the table already exists at index; otherwise false
// if the table didn't exist but was created.
//
// See: https://www.lua.org/manual/5.3/manual.html#luaL_getsubtable
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

// If the value at the given index has a metatable, the function pushes that metatable onto the stack
// and returns 1. Otherwise, the function returns 0 and pushes nothing on the stack.
func (state *State) GetMetaTable(index int) bool {
    unimplemented("GetMetaTable")
    return false
}

// Pops a table from the stack and sets it as the new metatable for the value at the given index.
func (state *State) SetMetaTable(index int) {
    unimplemented("SetMetaTable")
}

// GetTable pushes onto the stack the value t[k], where t is the value at the given index
// and k is the value at the top of the stack.
//
// This function pops the key from the stack, pushing the resulting value in its place.
// As in Lua, this function may trigger a metamethod for the "index" event (see §2.4).
func (state *State) GetTable(index int) Type {
    var (
        key = state.frame().pop()
        obj = state.get(index)
    )
    val := state.gettable(obj, key, 1)
    state.frame().push(val)
    return val.Type()
}

// SetTable does the equivalent to t[k] = v, where t is the value at the given index, v is the
// value at the top of the stack, and k is the value just below the top.
//
// This function pops both the key and the value from the stack. As in Lua, this function may
// trigger a metamethod for the "newindex" event (see §2.4).
func (state *State) SetTable(index int) {
    var (
        val = state.frame().pop()
        key = state.frame().pop()
        obj = state.get(index)
    )
    state.settable(obj, key, val, 1)
}

// Pushes onto the stack the value t[k], where t is the value at the given index.
//
// As in Lua, this function may trigger a metamethod for the "index" event (see §2.4).
//
// Returns the type of the pushed value.
func (state *State) GetField(index int, field string) Type {
    v := state.gettable(state.get(index), String(field), 1)
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
    state.settable(obj, key, val, 1)
}

// Pushes onto the stack the value t[i], where t is the value at the given index.
//
// As in Lua, this function may trigger a metamethod for the "index" event (see §2.4).
// 
// Returns the type of the pushed value.
func (state *State) GetI(index, entry int) Type {
    unimplemented("GetI")
    return NilType
}

// Does the equivalent to t[n] = v, where t is the value at the given index and v is the
// value at the top of the stack.
// 
// This function pops the value from the stack.
//
// As in Lua, this function may trigger a metamethod for the "newindex" event (see §2.4).
func (state *State) SetI(index, entry int) {
    unimplemented("SetI")
}

// Similar to lua_gettable, but does a raw access (i.e., without metamethods).
func (state *State) RawGet(index int) Type {
    unimplemented("RawGet")
    return NilType
}

// Similar to lua_settable, but does a raw assignment (i.e., without metamethods).
func (state *State) RawSet(index int) {
    unimplemented("RawSet")
}

// Pushes onto the stack the value t[n], where t is the table at the given index.
//
// The access is raw, that is, it does not invoke the __index metamethod.
// 
// Returns the type of the pushed value.
func (state *State) RawGetI(index int, entry int) Type {
    tbl, ok := state.get(index).(*Table)
    if !ok {
        state.errorf("expected table, found %s", tbl.Type())
    }
    val := tbl.getInt(int64(entry))
    state.frame().push(val)
    return val.Type()
}

// Pushes onto the stack the value t[k], where t is the table at the given index and
// k is the pointer p represented as a light userdata.
//
// The access is raw; that is, it does not invoke the __index metamethod.
// 
// Returns the type of the pushed value.
func (state *State) RawGetP(index int, udata *Object) Type {
    unimplemented("RawGetP")
    return NilType
}

// Does the equivalent of t[p] = v, where t is the table at the given index, p is
// encoded as a light userdata, and v is the value at the top of the stack.
// 
// This function pops the value from the stack.
//
// The assignment is raw, that is, it does not invoke __newindex metamethod.
func (state *State) RawSetP(index int, udata *Object) {
    unimplemented("RawSetP")
}

// Does the equivalent of t[i] = v, where t is the table at the given index and
// v is the value at the top of the stack.
// 
// This function pops the value from the stack.
//
// The assignment is raw, that is, it does not invoke the __newindex metamethod.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_rawseti
func (state *State) RawSetI(index int, entry int) {
    tbl, ok := state.get(index).(*Table)
    if !ok {
        state.errorf("table expected")
        return
    }
    tbl.setInt(int64(entry), state.Pop())
}

// Returns true if the two values in indices index1 and index2 are primitively equal
// (that is, without calling the __eq metamethod).
//
// Otherwise returns false. or if any of the indices are not valid.
func (state *State) RawEqual(i1, i2 int) bool {
    unimplemented("RawEquals")
    return false
}

// CallMeta calls a metamethod.
//
// If the object at index obj has a metatable and this metatable has a field event,
// this function calls this field passing the object as its only argument. In this
// case this function returns true and pushes onto the stack the value returned by
// the call. If there is no metatable or no metamethod, this function returns false
// (without pushing any value on the stack).
//
// [-0, +(0|1), e]
func (state *State) CallMeta(index int, event string) bool {
    unimplemented("CallMeta")
    return false
}

// Calls a function.
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
        // value  = state.frame().get(-(args+1))
        value  = state.frame().get(funcID)
        c, ok  = value.(*Closure)
    )

    state.Logf("call (func @ %d) %v (# args = %d, # rets = %d)", funcID, value, args, rets)

    if !ok {
        Debug(state)
        unimplemented(fmt.Sprintf("Call: __call: %v", value))
        return
    }
    state.call(&Frame{closure: c, fnID: funcID, rets: rets})
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

// Errorf raises an error with message formatted using the format string and provided
// arguments. The error message is prefixed with the file name and line number where
// the error occurred, if this information is available. This function never returns.
//
// [-0, +0, v]
func (state *State) Errorf(format string, args ...interface{}) { state.errorf(format, args...) }

// Exec loads and runs a Lua chunk returning the result (if any) or an error (if any).
//
// The Lua chunk may be provided via the filename of the source file, or via the
// source parameter.
//
// If source != nil, Do loads the source from source and the filename is only used
// recording position information.
func (state *State) Exec(filename string, source interface{}) error {
    if err := state.Load(filename, source); err != nil {
        return err
    }
    state.Call(0, 0)
    return nil
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
func Require(state *State, module string, loader Func, global bool) {
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

//
// TODO: remove me
//

func unimplemented(msg string) { panic(fmt.Errorf(msg)) }