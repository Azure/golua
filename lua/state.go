package lua

import (
    "path/filepath"
    "io/ioutil"
    "os/exec"
	"strings"
	"fmt"
	"io"
	"os"

    "github.com/Azure/golua/lua/binary"
    "github.com/Azure/golua/lua/syntax"
)

var _ = os.Exit

type runtimeErr error

// version number for this Lua implementation.
var version = float64(503)

type (
	// 'per thread' state.
	State struct {
		// shared global state
		global *global
		// execution state
		status ThreadStatus // thread status
		base   Frame 		// base call frame
		calls  int 			// call count
	}

	// 'global state', shared by all threads of a main state.
	global struct {
		builtins [maxTypeID]*Table
		version  *float64
		registry *Table
		thread0  *State
		config   *config
		panicFn  Func
	}
)

// NewState returns a new Lua thread state.
func NewState(opts ...Option) *State {
	// Initialize global configuration from options.
	var cfg config
	for _, opt := range opts {
		opt(&cfg)
	}

	// Lua execution state.
	state := new(State).reset()

	// Set up registry & globals table.
	var (
		registry = &Table{newTable(state, 8, 0)}
		globals  = &Table{newTable(state, 0, 20)}
		thread   = &Thread{state}
	)
	// Initialize registry.
	registry.setInt(MainThreadIndex, thread)
	registry.setInt(GlobalsIndex, globals)

	// Initialize the global state.
	state.enter(new(Frame))
	state.init(&global{
		registry: registry,
		version:  &version,
		thread0:  state,
		config:   &cfg,
	})

	return state
}

// traceback prints to w a stack trace from the current frame to the base.
func (state *State) traceback(w io.Writer) {
	fmt.Fprintln(w, "#")
	fmt.Fprintf(w,  "# traceback (calls = %d)\n", state.calls)
	fmt.Fprintln(w, "#")
	for fr := state.frame(); fr != nil; fr = fr.caller() {
		fmt.Fprintf(w, "function @ %d (%d returns)\n", fr.fnID, fr.rets)

		fmt.Fprintf(w, "    locals (%d)\n", fr.gettop())
		for top := fr.gettop() - 1; top >= 0; top-- {
			indent := "       "
			fmt.Fprintf(w, "%s[%d] %d @ %v (%T)\n",
				indent,
				top,
				top - (fr.gettop() + 1),
				fr.local(top),
				fr.local(top),
			)
		}
		fmt.Fprintln(w, "    end")
	}
	fmt.Fprintln(w)
}

// recover recovers any error thrown by the Lua runtime.
//
// If the error is not a runtimeErr, recover repanics the
// error up the stack.
func (state *State) recover(err *error) {
	if r := recover(); r != nil {
		if e, ok := r.(runtimeErr); ok {
			*err = e
			// if state.global.config.debug {
			// 	Debug(state)
			// }
			//return
		}
		panic(r)
	}
}

// errorf reports a formatted error message.
func (state *State) errorf(format string, args ...interface{}) {
	panic(runtimeErr(fmt.Errorf(format, args...)))
}

// enter enters a new call frame.
func (state *State) enter(fr *Frame) *Frame {
	state.ensure()
	fp := state.base.prev.next
	state.base.prev.next = fr
	fr.prev = state.base.prev
	fr.next = fp
	fp.prev = fr
	fr.state = state
	fr.depth = state.calls
	state.calls++
	return fr
}

// leave leaves the current frame.
func (state *State) leave(fr *Frame) *Frame {
	fr.prev.next = fr.next
	fr.next.prev = fr.prev
	fr.next  = nil // avoid memory leaks
	fr.prev  = nil // avoid memory leaks
	fr.state = nil
	state.calls--
	return fr
}

// reset resets the thread's call stack.
func (state *State) reset() *State {
	state.base.next = &state.base
	state.base.prev = &state.base
	state.calls = 0
	return state
}

// frame returns the current frame or nil.
func (state *State) frame() *Frame {
	if state.depth() == 0 {
		return nil
	}
	return state.base.prev
}

// ensure ensures the call frame stack is initialized.
func (state *State) ensure() {
	if state.base.next == nil {
		state.reset()
	}
}

// depth reports the current call depth.
func (state *State) depth() int { return state.calls }

// Calls a function (Go or Lua). The function to be called is at funcID in the stack.
// The arguments are on the stack in direct order following the function.

// On return, all the results are on the stack, starting at the original function position.
func (state *State) call(fr *Frame) {
	// Check that we are below the recursion / call max.
	if state.calls >= MaxCalls {
		panic(runtimeErr(fmt.Errorf("go: call stack overflow")))
	}

	// Ensure stack space for new call frame.
	fr.checkstack(InitialStackNew)

	// Push arguments and pop function.
	args := state.frame().popN(state.frame().gettop()-fr.fnID+1)[1:]
	fr.pushN(args)

	// Enter frame and leave on return.
	defer state.leave(state.enter(fr))
	
 	// Is it a Lua closure?
	if fr.closure.isLua() {
		// Ensure stack has space.
		fr.checkstack(fr.closure.binary.StackSize())

 		// Adjust the stack; params is the # of fixed real
		// parameters from the prototype, and fr.top holds
		// the # of passed arguments currently on the frame
		// local stack.
		switch params := fr.closure.binary.NumParams(); {
			case fr.gettop() < params: // # arguments < # parameters
				for fr.gettop() < params {
					fr.push(None) // nil to top
				}
			case fr.gettop() > params: // # arguments > # parameters
				extras := fr.popN(fr.gettop() - params)
				if fr.closure.binary.IsVararg() {
				    fr.vararg = extras
				}
		}

		// Execute the closure.
		execute(&v53{state})
		return
	}

	// Otherwise Go closure.
	if rets := fr.popN(fr.closure.native(state)); fr.rets != 0 {
		switch retc := len(rets); {
			case retc < fr.rets:
				for retc < fr.rets {
					rets = append(rets, None)
					retc++
				}
			case retc > fr.rets:
				if fr.rets != MultRets {
					rets = rets[:fr.rets]
				}
		}
		fr.caller().pushN(rets)
	}
}

func (state *State) init(g *global) {
	state.frame().checkstack(InitialStackNew)
	state.global = g
}

// EmitIR compiles the Lua script similarily to Compile but instead
// uses the long listing options to capture the compiled IR dump.
// EmitIR returns the Lua bytecode encoded as a string or any error.
func (state *State) emit(script string) {
    src, err := ioutil.ReadFile(script)
    if err != nil {
       panic(err)
    }
    cmd := exec.Command("luac", "-l", "-")
    cmd.Stdin = strings.NewReader(string(src))
    out, err := cmd.CombinedOutput()
    if err != nil {
        panic(fmt.Errorf("%v: %s", err, string(out)))
    }
    fmt.Fprintln(os.Stdout, string(out))
}

func (state *State) load(filename string, source interface{}) (*Closure, error) {
	var (
		src []byte
		err error
	)
	if src, err = syntax.Source(filename, source); err != nil {
		return nil, err
	}

	if !binary.IsChunk(src) {
	    dir, err := ioutil.TempDir("", "glua")
	    if err != nil {
	        return nil, err
	    }
	    tmp := filepath.Join(dir, "glua.bin")
	    cmd := exec.Command("luac", "-o", tmp, "-")
	    cmd.Stdin = strings.NewReader(string(src))

	    out, err := cmd.CombinedOutput()
	    if err != nil {
	        return nil, fmt.Errorf("%v: %s", err, string(out))
	    }

	    if src, err = ioutil.ReadFile(tmp); err != nil {
	        return nil, err
	    }
	}

	if state.global.config.debug {
		state.emit(filename)
	}

	chunk, err := binary.Load(src)
	if err != nil {
		return nil, err
	}

	cls := newLuaClosure(&chunk.Entry)
    if len(cls.upvals) > 0 {
        globals := state.global.registry.getInt(GlobalsIndex)
        cls.upvals[0] = &upValue{index: -1, value: globals}
    }
	return cls, nil
}

func (state *State) gettable(obj, key Value, raw bool) Value {
	if tbl, ok := obj.(*Table); ok {
		if val := tbl.Get(key); !IsNone(val) || raw || IsNone(state.metafield(tbl, "__index")) {
			return val
		}
		val, err := tryMetaIndex(state, tbl, key)
		if err != nil {
			panic(err)
		}
		return val
	}
	fmt.Printf("gettable(%v, %v, %t)\n", obj, key, raw)
	state.Debug(true)
	panic(runtimeErr(fmt.Errorf("table expected, got %v", obj)))
}

func (state *State) settable(obj, key, val Value, raw bool) {
	if tbl, ok := obj.(*Table); ok && (!tbl.exists(key) || raw) {
		tbl.Set(key, val)
		return
	}
	if err := tryMetaNewIndex(state, obj, key, val); err != nil {
		panic(err)
	}
}

func (state *State) metafield(value Value, event string) Value {
	if obj := state.getmetatable(value, true); !IsNone(obj) {
		if tbl, ok := obj.(*Table); ok && tbl != nil {
			return tbl.getStr(event)
		}
	}
	return nilValue
}

func (state *State) setmetatable(value, meta Value) {
	mt, ok := meta.(*Table)
	if !ok && !IsNone(meta) {
		state.errorf("metatable must be table or nil")
	}
	switch v := value.(type) {
		case *Object:
			v.meta = mt
		case *Table:
			v.meta = mt
		default:
			state.global.builtins[v.Type()] = mt
	}
}

func (state *State) getmetatable(value Value, rawget bool) Value {
	var meta Value = None
	switch value := value.(type) {
		case *Object:
			if !IsNone(value.meta) {
				meta = value.meta
			}
		case *Table:
			if !IsNone(value.meta) {
				meta = value.meta
			}
		default:
			meta = state.global.builtins[value.Type()] 
	}
	if !rawget && !IsNone(meta) {
		if mt, ok := meta.(*Table); ok {
			if mm := mt.getStr("__metatable"); !IsNone(mm) {
				meta = mm
			}
		}
	}
	return meta
}

func (state *State) Logf(format string, args ...interface{}) {
	state.Log(fmt.Sprintf(format, args...))
}

func (state *State) Log(args ...interface{}) {
	if state.global.config.trace {
		fmt.Fprintf(os.Stdout, "lua: %v\n", fmt.Sprint(args...))
	}
}