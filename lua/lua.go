package lua

import (
	"fmt"
	"os"

	"github.com/Azure/golua/lua/binary"
	"github.com/Azure/golua/lua/syntax"
)

const Version = "Lua 5.3"

var errNonBinaryChunk = fmt.Errorf("lua: execute: non-binary files not yet supported")

type Config struct {
	Error func(error)
	Trace bool
}

var _ = os.Exit

// CheckStack ensures that the stack has space for at least n extra slots (that is, that you can
// safely push up to n values into it). It returns false if it cannot fulfill the request,
// either because it would cause the stack to be larger than a fixed maximum size (typically
// at least several thousand elements) or because it cannot allocate memory for the extra space.
//
// This function never shrinks the stack; if the stack already has space for the extra slots, it
// is left unchanged.
func (state *State) CheckStack(needed int) bool {
	unimplemented("*State.CheckStack")
	return false
}

// Pop pops n elements from the stack.
func (state *State) Pop(n int) {
	state.stack(0).popN(n)
}

// PushNil pushes a nil value onto the stack.
func (state *State) PushNil() {
	state.stack(0).push(None)
}

// PushInt pushes an integer number onto the stack.
func (state *State) PushInt(v int64) {
	state.stack(0).push(Int(v))
}

// PushBool pushes a boolean value with value b onto the stack.
func (state *State) PushBool(v bool) {
	state.stack(0).push(Bool(v))
}

// PushFloat pushes a float number onto the stack.
func (state *State) PushFloat(v float64) {
	state.stack(0).push(Float(v))
}

// PushString pushes a string value onto the stack.
func (state *State) PushString(v string) {
	state.stack(0).push(String(v))
}

// PushValue pushes a copy of the element at the given index onto the stack.
func (state *State) PushValue(index int) {
	state.stack(0).push(state.index(index))
}

// PushGoFunc pushes a Go function onto the stack. This function receives a pointer
// to a Go function and pushes onto the stack a lua value of type function that,
// when called, invokes the corresponding Go function.
//
// Any function to be callable by lua must follow the correct protocol to receive
// its parameters and return its results.
func (state *State) PushFunc(fn Func) {
	state.PushClosure(fn)
}

// PushGoClosure pushes a new Go closure onto the stack.
//
// When a Go function is created, it is possible to associate some values with it, thus
// creating a Go closure; these values are then accessible to the function whenever it
// is called.
//
// To associate values with a Go function, first these values must be pushed onto the stack
// (when there are multiple values, the first value is pushed first). Then PushGoClosure is
// called to create and push the Go function onto the stack, with the argument n telling how
// many values will be associated with the function.
//
// PushGoClosure also pops these values from the stack.
//
// The maximum value for nups is 255.
//
// When nups is zero, this function creates a light go function, which is just a pointer to the
// go function. In that case, it never raises a memory error. 
func (state *State) PushClosure(fn Func, upvalues ...Value) {
	if len(upvalues) >= MaxUpValues {
		panic(fmt.Errorf("upvalue index too large"))
	}
	state.stack(0).push(newGoClosure(fn, len(upvalues)))
}

// Push pushes the Lua value onto the stack.
func (state *State) Push(value Value) {
	state.stack(0).push(value)
}

// PushGlobals pushes onto the stack the globals table.
func (state *State) Globals() *Table {
	return state.global.registry.getInt(GlobalsIndex).(*Table)
}

// AbsIndex converts the acceptable index idx into an equivalent absolute index;
// that is, one that does not depend on the stack top).
func (state *State) AbsIndex(index int) int {
	// zero, positive, or pseudo index
	if index >= 0 || isPseudoIndex(index) {
		return index
	}
	return state.stack(0).top() + index + 1
}

// GetTop returns the index of the top element in the stack.
//
// Because indices start at 1, this result is equal to the number
// of elements in the stack; in particular, 0 means an empty stack.
func (state *State) GetTop() int {
	return state.stack(0).top()
}

// SetTop accepts any index, or 0, and sets the stack top to this index.
//
// If the new top is larger than the old one, then the new elements are
// filled with nil. If index is 0, then all stack elements are removed.
func (state *State) SetTop(int) {
	unimplemented("*State.SetTop")
}

// Copy copies the element at index fromidx into the valid index toidx, replacing the value at that position.
//
// Values at other positions are not affected.
func (state *State) Copy(from, to int) {
	unimplemented("*State.Copy")
}

// Insert moves the top element into the given valid index, shifting up the elements above this
// index to open space.
//
// This function cannot be called with a pseudo-index, because a pseudo-index is not an actual
// stack position.
func (state *State) Insert(index int) {
	unimplemented("*State.Insert")
}

// Replace moves the top element into the given valid index without shifting any element (therefore
// replacing the value at that given index), and then pops the top element.
func (state *State) Replace(index int) {
	unimplemented("*State.Replace")
}

// Remove removes the element at the given valid index, shifting down the elements above
// this index to fill the gap. This function cannot be called with a pseudo-index, because
// a pseudo-index is not an actual stack position.
func (state *State) Remove(index int) {
	unimplemented("*State.Remove")
}

// Rotate rotates the stack elements between the valid index idx and the top of the stack.
//
// The elements are rotated n positions in the direction of the top, for a positive n, or -n
// positions in the direction of the bottom, for a negative n. The absolute value of n must
// not be greater than the size of the slice being rotated. This function cannot be called
// with a pseudo-index, because a pseudo-index is not an actual stack position.
func (state *State) Rotate(index int, n int) {
	unimplemented("*State.Rotate")
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
	b, err := syntax.Source(filename, source)
	if err != nil {
		return err
	}
	if !binary.IsChunk(b) {
		return errNonBinaryChunk
	}
	c, err := binary.Unpack(b)
	if err != nil {
		return err
	}

	if state.global.config.Trace {
		fmt.Println(&c)
	}

	// type NativeFn func(*State)int
	// type Function struct {
	//		proto  *binary.Prototype
	//		native NativeFn
	// }

	cls := newLuaClosure(&c.Entry)
	if len(cls.upvals) > 0 {
		globals := state.global.registry.getInt(GlobalsIndex)
		cls.upvals[0] = &globals		
	}
	state.stack(0).push(cls)

	// if err := lvm.chunk(&c); err != nil {
	// 	return nil, err
	// }
	// return lvm.thread.Stack().Index(-1), nil
	return nil
}

// If modname is not already present in package.loaded, calls function loader with string modname
// as an argument and sets the call result in package.loaded[modname], as if that function has been
// called through require.
//
// If global is true, also stores the module into global modname.
//
// Leaves a copy of the module on the stack.
func (state *State) Require(module string, loader Func, global bool) {
	state.GetSubTable(RegistryIndex, "_LOADED")
	state.GetField(-1, module)
	if !state.ToBool(-1) {
		state.Pop(1)
		state.PushFunc(loader)
		state.PushString(module)
		state.Call(1, 1)
		state.PushValue(-1)
		state.SetField(-3, module)
	}
	state.Remove(-2)

	if global {
		state.PushValue(-1)
		state.SetGlobal(module)
	}
}

// Register sets the Go function fn as the new value of global name.
func (state *State) Register(name string, fn Func) {
	state.PushFunc(fn)
	state.SetGlobal(name)
}

// AtPanic sets a new panic function and returns the old one (see §4.6).
func (state *State) AtPanic(fn Func) Func {
	unimplemented("*State.AtPanic")
	return nil
}

// Status returns the status of the thread.
// 
// The status can be 0 (LUA_OK) for a normal thread, an error code if the thread finished
// the execution of a lua_resume with an error, or LUA_YIELD if the thread is suspended.
// 
// You can only call functions in threads with status LUA_OK. You can resume threads with
// status LUA_OK (to start a new coroutine) or LUA_YIELD (to resume a coroutine).
func (state *State) Status() ThreadStatus {
	unimplemented("*State.Status")
	return ThreadError
}

// Generates a Lua error, using the value at the top of the stack as the error object.
// 
// This function does a long jump, and therefore never returns (see luaL_error).
func (state *State) Error() {
	unimplemented("*State.Error")
}

// Destroys all objects in the given Lua state (calling the corresponding garbage-collection metamethods, if any) and
// frees all dynamic memory used by this state. In several platforms, you may not need to call this function, because
// all resources are naturally released when the host program ends. On the other hand, long-running programs that create
// multiple states, such as daemons or web servers, will probably need to close states as soon as they are not needed.
func (state *State) Close() {
	//unimplemented("*State.Close")
}

// Returns the address of the version number (a C static variable) stored in the Lua core. When called with a valid lua_State,
// returns the address of the version used to create that state. When called with NULL, returns the address of the version
// running the call.
func (state *State) Version() *float64 {
	unimplemented("*State.Version")
	return nil
}

// Performs an arithmetic or bitwise operation over the two values (or one, in the case of negations) at the top of the
// stack, with the value at the top being the second operand, pops these values, and pushes the result of the operation.
// The function follows the semantics of the corresponding Lua operator (that is, it may call metamethods).
// 
// The value of op must be one of the following constants:
// 		* LUA_OPADD: performs addition (+)
// 		* LUA_OPSUB: performs subtraction (-)
// 		* LUA_OPMUL: performs multiplication (*)
// 		* LUA_OPDIV: performs float division (/)
// 		* LUA_OPIDIV: performs floor division (//)
// 		* LUA_OPMOD: performs modulo (%)
// 		* LUA_OPPOW: performs exponentiation (^)
// 		* LUA_OPUNM: performs mathematical negation (unary -)
// 		* LUA_OPBNOT: performs bitwise NOT (~)
// 		* LUA_OPBAND: performs bitwise AND (&)
// 		* LUA_OPBOR: performs bitwise OR (|)
// 		* LUA_OPBXOR: performs bitwise exclusive OR (~)
// 		* LUA_OPSHL: performs left shift (<<)
// 		* LUA_OPSHR: performs right shift (>>)
func (state *State) Arith(op Op) {
	unimplemented("*State.Arith")
}

// Concatenates the n values at the top of the stack, pops them, and leaves the result at the top. If n is 1, the result
// is the single value on the stack (that is, the function does nothing); if n is 0, the result is the empty string.
// Concatenation is performed following the usual semantics of Lua (see https://www.lua.org/manual/5.3/manual.html#3.4.6).
func (state *State) Concat(n int) {
	unimplemented("*State.Concat")
}

// Returns the length of the value at the given index. It is equivalent to the ‘\#’ operator in Lua
// (see [Lua 5.3 Reference Manual](https://www.lua.org/manual/5.3/manual.html#3.4.7)) and may trigger
// a metamethod for the “length” event (see [Lua 5.3 Reference Manual](https://www.lua.org/manual/5.3/manual.html#2.4)).
//The result is pushed on the stack.
func (state *State) Length(index int) int {
	unimplemented("*State.Length")
	return 0
}

// Compares two Lua values. Returns 1 if the value at index index1 satisfies op when compared with the value at index index2, following the semantics of the corresponding Lua operator (that is, it may call metamethods). Otherwise returns 0. Also returns 0 if any of the indices is not valid.
// 
// The value of op must be one of the following constants:
// 		* LUA_OPEQ: compares for equality (==)
// 		* LUA_OPLT: compares for less than (<)
// 		* LUA_OPLE: compares for less or equal (<=)
func (state *State) Compare(op Op, i1, i2 int) bool {
	unimplemented("*State.Compare")
	return false
}

// Pushes onto the stack the value of the global name.
// 
// Returns the type of that value.
func (state *State) GetGlobal(name string) Type {
	tbl := state.global.registry.getInt(GlobalsIndex)
	return state.gettable(tbl, String(name), 1).Type()
}

// Pops a value from the stack and sets it as the new value of global name.
func (state *State) SetGlobal(name string) {
	var (
		tbl = state.global.registry.getInt(GlobalsIndex)
		val = state.stack(0).pop()
	)
	state.settable(tbl, String(name), val, 1)
}

// Creates a new empty table and pushes it onto the stack. Parameter narr is a hint for how
// many elements the table will have as a sequence; parameter nrec is a hint for how many other
// elements the table will have. Lua may use these hints to preallocate memory for the new table.
// This preallocation is useful for performance when you know in advance how many elements the table
// will have. Otherwise you can use the function lua_newtable.
func (state *State) CreateTable(narr, nrec int) {
	state.stack(0).push(&Table{newTable(state, narr, nrec)})
}

// Creates a new empty table and pushes it onto the stack.
// 
// It is equivalent to lua_createtable(L, 0, 0).
func (state *State) NewTable() {
	state.CreateTable(0, 0)
}

func (state *State) GetSubTable(index int, field string) bool {
	if state.GetField(index, field) == TableType {
		return true
	}
	state.Pop(1)
	index = state.AbsIndex(index)
	state.NewTable()
	state.PushValue(-1)
	state.SetField(index, field)
	return false
}

// If the value at the given index has a metatable, the function pushes that metatable onto the stack
// and returns 1. Otherwise, the function returns 0 and pushes nothing on the stack.
func (state *State) GetMetaTable(index int) bool {
	unimplemented("*State.GetMetaTable")
	return false
}

// Pops a table from the stack and sets it as the new metatable for the value at the given index.
func (state *State) SetMetaTable(index int) {
	unimplemented("*State.SetMetaTable")
}

// GetTable pushes onto the stack the value t[k], where t is the value at the given index
// and k is the value at the top of the stack.
//
// This function pops the key from the stack, pushing the resulting value in its place.
// As in Lua, this function may trigger a metamethod for the "index" event (see §2.4).
func (state *State) GetTable(index int) Type {
	var (
		key = state.stack(0).pop()
		tbl = state.index(index)
	)
	v := state.gettable(tbl, key, 1)
	state.stack(0).push(v)
	return v.Type()
}

// SetTable does the equivalent to t[k] = v, where t is the value at the given index, v is the
// value at the top of the stack, and k is the value just below the top.
//
// This function pops both the key and the value from the stack. As in Lua, this function may
// trigger a metamethod for the "newindex" event (see §2.4).
func (state *State) SetTable(index int) {
	unimplemented("*State.SetTable")
}

// Pushes onto the stack the value t[k], where t is the value at the given index.
//
// As in Lua, this function may trigger a metamethod for the "index" event (see §2.4).
//
// Returns the type of the pushed value.
func (state *State) GetField(index int, field string) Type {
	v := state.gettable(state.index(index), String(field), 1)
	state.stack(0).push(v)
	return v.Type()
}

// Does the equivalent to t[k] = v, where t is the value at the given index and v is the
// value at the top of the stack.
// 
// This function pops the value from the stack.
//
// As in Lua, this function may trigger a metamethod for the "newindex" event (see §2.4).
func (state *State) SetField(index int, field string) {
	var (
		val = state.stack(0).pop()
		tbl = state.index(index)
	)
	state.settable(tbl, String(field), val, 1)
}

// Pushes onto the stack the value t[i], where t is the value at the given index.
//
// As in Lua, this function may trigger a metamethod for the "index" event (see §2.4).
// 
// Returns the type of the pushed value.
func (state *State) GetI(index, entry int) Type {
	unimplemented("*State.GetI")
	return NilType
}

// Does the equivalent to t[n] = v, where t is the value at the given index and v is the
// value at the top of the stack.
// 
// This function pops the value from the stack.
//
// As in Lua, this function may trigger a metamethod for the "newindex" event (see §2.4).
func (state *State) SetI(index, entry int) {
	unimplemented("*State.SetI")
}

// Similar to lua_gettable, but does a raw access (i.e., without metamethods).
func (state *State) RawGet(index int) Type {
	unimplemented("*State.RawGet")
	return NilType
}

// Similar to lua_settable, but does a raw assignment (i.e., without metamethods).
func (state *State) RawSet(index int) {
	unimplemented("*State.RawSet")
}

// Pushes onto the stack the value t[n], where t is the table at the given index.
//
// The access is raw, that is, it does not invoke the __index metamethod.
// 
// Returns the type of the pushed value.
func (state *State) RawGetI(index int, entry int) Type {
	tbl, ok := state.index(index).(*Table)
	if !ok {
		state.errorf("expected table, found %s", tbl.Type())
	}
	val := tbl.Get(Int(entry))
	state.stack(0).push(val)
	return val.Type()
}

// Pushes onto the stack the value t[k], where t is the table at the given index and
// k is the pointer p represented as a light userdata.
//
// The access is raw; that is, it does not invoke the __index metamethod.
// 
// Returns the type of the pushed value.
func (state *State) RawGetP(index int, udata *Object) Type {
	unimplemented("*State.RawGetP")
	return NilType
}

// Does the equivalent of t[p] = v, where t is the table at the given index, p is
// encoded as a light userdata, and v is the value at the top of the stack.
// 
// This function pops the value from the stack.
//
// The assignment is raw, that is, it does not invoke __newindex metamethod.
func (state *State) RawSetP(index int, udata *Object) {
	unimplemented("*State.RawSetP")
}

// Does the equivalent of t[i] = v, where t is the table at the given index and
// v is the value at the top of the stack.
// 
// This function pops the value from the stack.
//
// The assignment is raw, that is, it does not invoke the __newindex metamethod.
func (state *State) RawSetI(index int, entry int) {
	unimplemented("*State.RawSetI")
}

// Returns true if the two values in indices index1 and index2 are primitively equal
// (that is, without calling the __eq metamethod).
//
// Otherwise returns false. or if any of the indices are not valid.
func (state *State) RawEqual(i1, i2 int) bool {
	unimplemented("*State.RawEquals")
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
	unimplemented("*State.CallMeta")
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
func (state *State) Call(argN, retN int) {
	var (
		value = state.stack(0).get(state.AbsIndex(-(argN+1)))
		c, ok = value.(*Closure)
	)
	if !ok {
		// TODO: __call
	}
	state.call(c, argN, retN)
	os.Exit(1)
}

// Registers all functions in the array l (see luaL_Reg) into the table on the top of the stack (below optional
// upvalues, see next).
//
// When nup is not zero, all functions are created sharing nup upvalues, which must be previously pushed on the
// stack on top of the library table. These values are popped from the stack after the registration.
func (state *State) SetFuncs(table *Table, funcs map[string]Func, upvalues ...Value) *Table {
	for name, fn := range funcs {
		cls := newGoClosure(fn, len(upvalues))
		for i, upvalue := range upvalues {
			cls.upvals[i] = &upvalue
		}
		table.setStr(name, cls)
	}
	return table
}

// Errorf raises an error with message formatted using the format string and provided
// arguments. The error message is prefixed with the file name and line number where
// the error occurred, if this information is available. This function never returns.
//
// [-0, +0, v]
func (state *State) Errorf(format string, args ...interface{}) {
	unimplemented("*State.Errorf")
}

//
// TODO: remove me
//

func (state *State) Safely(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				state.Dump()
				err = e
				panic(err)
			}
		}
	}()
	fn()
	// dump registry
	// dump globals
	// dump stack
	return
}

func (state *State) Dump() {
	// fmt.Println(state.global.registry.table)
	fmt.Println(state.stack(0))
}

func unimplemented(msg string) {
	panic(fmt.Errorf(msg))
}

func assert(cond bool, ctxt string) {
	if cond {
		panic(fmt.Errorf(ctxt))
	}
}