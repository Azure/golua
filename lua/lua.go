package lua

import (
    "fmt"
    "io"
    "os"

    "github.com/Azure/golua/lua/binary"
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
func (state *State) CheckStack(needed int) bool {
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
    return fmt.Sprintf("%p", ls)
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
        cls.upvals[nups-1] = &upValue{
            index: -1,
            value: state.Pop(),
        }
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

// SetTop accepts any index, or 0, and sets the stack top to this index. If the new top
// is larger than the old one, then the new elements are filled with nil. If index is 0,
// then all stack elements are removed.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_settop
func (state *State) SetTop(top int) {
    if top = state.frame().absindex(top); top < 0 {
        panic(runtimeErr(fmt.Errorf("stack underflow!")))
    }
    state.frame().settop(top)
}

// Top returns the index of the top element in the stack.
//
// Because indices start at 1, this result is equal to the number
// of elements in the stack; in particular, 0 means an empty stack.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_gettop
func (state *State) Top() int { return state.frame().gettop() }

// Push pushes any value onto the stack first boxing it by the equivalent
// Lua value, returning its position in the frame's local stack (top - 1).
func (state *State) Push(any interface{}) int {
    state.frame().push(valueOf(state, any))
    return state.Top() - 1
}

// PopN pops the top n values from the Lua thread's current frame stack.
func (state *State) PopN(n int) []Value { return state.frame().popN(n) }

// Pop pops the top value from the Lua thread's current frame stack.
func (state *State) Pop() Value { return state.frame().pop() }

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

// PushGlobals pushes onto the stack the globals table.
func (state *State) Globals() *Table {
    return state.global.registry.getInt(GlobalsIndex).(*Table)
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
func (state *State) Load(filename string, source interface{}, mode Mode) error {
    cls, err := state.load(filename, source)
    if err != nil {
        return err
    }
    state.frame().push(cls)
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
func (state *State) Error() int {
    return state.errorf("%v", state.frame().pop())
}

// Destroys all objects in the given Lua state (calling the corresponding garbage-collection metamethods, if any) and
// frees all dynamic memory used by this state. In several platforms, you may not need to call this function, because
// all resources are naturally released when the host program ends. On the other hand, long-running programs that create
// multiple states, such as daemons or web servers, will probably need to close states as soon as they are not needed.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_close
func (state *State) Close() { fmt.Println("lua: state.Close(): TODO") }

// Returns the address of the version number (a C static variable) stored in the Lua core. When called with a valid lua_State,
// returns the address of the version used to create that state. When called with NULL, returns the address of the version
// running the call.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_version
func (state *State) Version() *float64 { unimplemented("Version"); return nil }

// XMove exchanges values between different threads of the same state.
//
// This function pops N values from the state's stack, and pushes them onto
// the state dst's stack.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_xmove
func (state *State) XMove(dst *State, n int) {
    dst.frame().pushN(state.frame().popN(n))
}

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
    obj := state.global.registry.getInt(GlobalsIndex)
    val := state.gettable(obj, String(name), false)
    state.frame().push(val)
    return val.Type()
}

// Pops a value from the stack and sets it as the new value of global name.
//func (state *State) SetGlobal(name string, value Value) {
func (state *State) SetGlobal(name string) {
    key := String(name)
    val := state.Pop()
    state.settable(state.Globals(), key, val, false)
    //state.settable(state.Globals(), String(name), value, 1)
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
    tbl, ok := state.get(index).(*Table)
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
func (state *State) SetMetaTableAt(index int) {
    state.setmetatable(state.get(index), state.frame().pop())
}

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
func (state *State) GetI(index int, entry int64) Type {
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
func (state *State) SetI(index int, entry int64) {
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
        case *Table:
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
func (state *State) RawGetI(index int, entry int) Type {
    var (
        obj = state.get(index)
        key = Int(entry)
        val = state.gettable(obj, key, true)
    )
    state.frame().push(val)
    return val.Type()
}

// Pushes onto the stack the value t[k], where t is the table at the given index and
// k is the pointer p represented as a light userdata.
//
// The access is raw; that is, it does not invoke the __index metamethod.
// 
// Returns the type of the pushed value.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_rawgetp
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
//
// See https://www.lua.org/manual/5.3/manual.html#lua_rawsetp
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

    state.Logf("call (func @ %d) %v (# args = %d, # rets = %d)\n", funcID, value, args, rets)
    
    if !ok && !tryMetaCall(state, value, funcID, args, rets) {
        state.errorf("attempt to call a %s value", value.Type())
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

// Exec loads and runs a Lua chunk returning the result (if any) or an error (if any).
//
// The Lua chunk may be provided via the filename of the source file, or via the
// source parameter.
//
// If source != nil, Do loads the source from source and the filename is only used
// recording position information.
func (state *State) Exec(filename string, source interface{}, mode Mode) error {
    if err := state.Load(filename, source, mode); err != nil {
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
        //state.SetGlobal(module, state.frame().pop())
    }
}

//
// TODO: remove me
//

func unimplemented(msg string) { panic(fmt.Errorf(msg)) }