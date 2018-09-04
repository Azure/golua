package lua

import (
	"fmt"
)

const (
	// RegistryIndex is the pseudo-index for the registry table.
	RegistryIndex = -DefaultStackMax - 1000

	// MainThreadIndex is the registry index for the main thread of the main state.
	MainThreadIndex = 1

	// GlobalsIndex is the registry index for the global environment.
	GlobalsIndex = 2
)

const (
	// Maximum valid index and maximum size of stack.
	DefaultStackMax = 1000000

	// Minimum Lua stack available to a function.
	DefaultStackMin = 20

	// Size allocated for new stacks.
	InitialStackNew = 2 * DefaultStackMin

	// Initial space allocate for UpValues.
	InitialFreeMax = 5

	// Maximum number of upvalues in a closure (both lua and go).
	// Value must fit in a VM register.
	MaxUpValues = 255

	// Limit for table tag-method chains (to avoid loops).
	MethodChainMax = 2000
)

// version number for this Lua implementation.
var version = float64(503)

// ThreadStatus is a Lua thread status.
type ThreadStatus int

// thread statuses
const (
	ThreadOK ThreadStatus = iota // Thread is in success state 
	ThreadYield 				 // Thread is in suspended state
	ThreadError 				 // Thread finished execution with error
)

// String returns the canoncial string of the thread status.
func (status ThreadStatus) String() string {
	switch status {
		case ThreadOK:
			return "OK"
		case ThreadYield:
			return "YIELD"
		case ThreadError:
			return "ERROR"
	}
	return fmt.Sprintf("unknown thread status %d", status)
}

// Op represents a Lua arithmetic or comparison operator.
type Op int

type (
	// 'per thread' state.
	State struct {
		status ThreadStatus
		global *global
		frame  *Frame
	}

	// call frame for function invocations.
	Frame struct {
		closure *Closure // closure for this frame
		varargs []Value  // variable arguments
		thread  *State   // executing thread
		caller  *Frame   // parent call frame
		stack   stack    // call frame stack
		nrets   int      // expected number of returns
		status  int      // call status filled in on return
		pc 		int 	 // instruction program counter
	}

	// 'global state', shared by all threads of a main state.
	global struct {
		builtins [maxTypeID]*Table
		version  *float64
		registry *Table
		thread0  *State
		config   Config
		panicFn  Func
	}
)

// stack_init
// init_registry
// luaS_init
// luaT_init
// set version
func NewState(config Config) *State {
	// Lua execution state.
	var state State

	// Set up registry & globals table.
	var (
		registry = &Table{newTable(&state, 8, 0)}
		globals  = &Table{newTable(&state, 0, 20)}
		thread   = &Thread{&state}
	)
	registry.setInt(MainThreadIndex, thread)
	registry.setInt(GlobalsIndex, globals)

	// Initialize global state.
	global := &global{
		registry: registry,
		version:  &version,
		thread0:  &state,
		config:   config,
	}

	// Setup initial call frame.
	var (
		stack = make(stack, 0, InitialStackNew)
		frame = &Frame{thread: &state, stack: stack}
	)

	// Initialize state.
	state = State{global: global, frame: frame}
	return &state
}

// pops a table from the stack and sets it as the new metatable
// for the value at the given index.
func (state *State) setmetatable(index int) {
	unimplemented("setmetatable")
}

func (state *State) getmetatable(value Value, rawget bool) Value {
	if isNilOrNone(value) {
		return None
	}
	var meta Value = None
	switch value := value.(type) {
		case *Object:
			if !isNilOrNone(value.meta) {
				meta = value.meta
			}
		case *Table:
			if !isNilOrNone(value.meta) {
				meta = value.meta
			}
		default:
			if mt := state.global.builtins[value.Type()]; !isNilOrNone(mt) {
				meta = mt
			}
	}
	if !rawget && !isNilOrNone(value) {
		if mt, ok := meta.(*Table); ok {
			if mm := mt.getStr("__metatable"); !isNilOrNone(mm) {
				meta = mm
			}
		}
	}
	return meta
}

func (state *State) metafield(value Value, event metaEvent) Value {
	if meta := state.getmetatable(value, true); meta != nil {
		if tbl, ok := meta.(*Table); ok && tbl != nil {
			return tbl.getStr(event.toName())
		}
	}
	return None
}

func (state *State) gettable(obj, key Value, chain int) Value {
	if chain > MethodChainMax {
		state.errorf("'__index' chain too long; possible loop")
	}
	if tbl, ok := obj.(*Table); ok && tbl.exists(key) {
		return tbl.Get(key)
	}
	switch meta := state.metafield(obj, metaIndex).(type) {
		case *Closure:
			state.stack(0).push(meta)
			state.stack(0).push(obj)
			state.stack(0).push(key)
			state.Call(2, 1)
			return state.stack(0).pop()
		case *Table:
			return state.gettable(meta, key, chain+1)
		default:
			if _, ok := obj.(*Table); !ok {
				state.errorf("attempt to index a %s value (%v)", obj.Type(), key)
			}
	}
	return None
}

func (state *State) settable(obj, key, val Value, chain int) {
	if chain > MethodChainMax {
		state.errorf("'__newindex' chain too long; possible loop")
	}
	if tbl, ok := obj.(*Table); ok && tbl.exists(key) {
		tbl.Set(key, val)
		return
	}
	switch meta := state.metafield(obj, metaNewIndex).(type) {
		case *Closure:
			state.stack(0).push(meta)
			state.stack(0).push(obj)
			state.stack(0).push(key)
			state.stack(0).push(val)
			state.Call(3, 0)
		case *Table:
			state.settable(meta, key, val, chain+1)
		default:
			if _, ok := obj.(*Table); !ok {
				state.errorf("attempt to index a %s value (%v)", obj.Type(), key)
			}
	}
}

func (state *State) errorf(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}

func (state *State) stack(depth int) *stack { return &state.frame.stack }

func (state *State) call(cls *Closure, argN, retN int) {
	switch {
		case cls.native != nil: // call go
			fmt.Printf("TODO: CALL GO (args=%d, rets=%d)\n", argN, retN)
			state.Dump()
			fmt.Println(cls.native(state))
		case cls.proto != nil: // call lua
			fmt.Println("TODO: CALL LUA")
	}
}

//
// get functions (lua -> stack)
//

// auxgetStr
// getglobal
// gettable
// getfield
// geti
// rawget
// rawgeti
// rawgetp
// createtable
// getmetatable
// getuservalue

//
// set functions (stack -> lua)
//

// auxsetstr
// setglobal
// settable
// setfield
// seti
// rawset
// rawseti
// rawsetp
// setmetatable
// setuservalue