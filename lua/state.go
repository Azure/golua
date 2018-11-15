package lua

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	// "github.com/Azure/golua/pkg/goutils"
	"github.com/Azure/golua/lua/binary"
	"github.com/Azure/golua/lua/syntax"
)

var _ = os.Exit

type runtimeErr error

// version number for this Lua implementation.
var version = float64(503)

// ThreadStatus is a Lua thread status.
type ThreadStatus int

// thread statuses
const (
	ThreadOK    ThreadStatus = iota // Thread is in success state
	ThreadYield                     // Thread is in suspended state
	ThreadError                     // Thread finished execution with error
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

// Lua execution thread.
type thread struct {
	*State
}

func (x *thread) String() string { return "thread" }
func (x *thread) Type() Type     { return ThreadType }

type (
	// 'per thread' state.
	State struct {
		// shared global state
		global *global
		// execution state
		status ThreadStatus // thread status
		base   Frame        // base call frame
		calls  int          // call count
	}

	// 'global state', shared by all threads of a main state.
	global struct {
		builtins [maxTypeID]*table
		version  *float64
		registry *table
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
		registry = newTable(state, 8, 0)
		globals  = newTable(state, 0, 20)
		thread   = &thread{state}
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

// String returns a printable string of the current executing thread state.
func (ls *State) String() string { return fmt.Sprintf("%p", ls) }

// traceback prints to w a stack trace from the current frame to the base.
func (state *State) traceback(w io.Writer) {
	fmt.Fprintln(w, "#")
	fmt.Fprintf(w, "# traceback (calls = %d)\n", state.calls)
	fmt.Fprintln(w, "#")
	for fr := state.frame(); fr != nil; fr = fr.caller() {
		fmt.Fprintf(w, "function @ %d (%d returns)\n", fr.fnID, fr.rets)

		fmt.Fprintf(w, "    locals (%d)\n", fr.gettop())
		for top := fr.gettop() - 1; top >= 0; top-- {
			indent := "       "
			fmt.Fprintf(w, "%s[%d] %d @ %v (%T)\n",
				indent,
				top,
				top-(fr.gettop()+1),
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
	// goutils.DumpStack(err, state.recover)
	if r := recover(); r != nil {
		state.Debug(false)
		if e, ok := r.(runtimeErr); ok {
			*err = e
		}
		panic(r)
	}
}

// safely executes the function fn returning any errors recovered by the Lua
// runtime
func (ls *State) safely(fn func() error) (err error) {
	defer ls.recover(&err)
	return fn()
}

// value returns the Lua value at the valid index.
func (state *State) value(index int) Value { return state.get(index) }

// errorf reports a formatted error message.
func (state *State) errorf(format string, args ...interface{}) int {
	return state.panic(runtimeErr(fmt.Errorf(format, args...)))
}

// panic generates a Lua error, using the value at the top of the stack at the
// error object. This function panics, and never returns.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_error
func (state *State) panic(err error) int { panic(err) }

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
	fr.next = nil // avoid memory leaks
	fr.prev = nil // avoid memory leaks
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

// globals returns the globals table.
func (state *State) globals() *table {
	return state.global.registry.getInt(GlobalsIndex).(*table)
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

// valueAt returns the Value at the given index if valid; otherwise None.
func (state *State) valueAt(index int) Value {
	if state.isValid(index) {
		return state.get(index)
	}
	return None
}

// isValid reports whether the index points to a valid Value.
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

// Calls a function (Go or Lua). The function to be called is at funcID in the stack.
// The arguments are on the stack in direct order following the function.
//
// On return, all the results are on the stack, starting at the original function position.
func (state *State) call(fr *Frame) {
	// Check that we are below the recursion / call max.
	if state.calls >= MaxCalls {
		state.Errorf("go: call stack overflow")
	}

	// Ensure stack space for new call frame.
	fr.checkstack(InitialStackNew)

	// Push arguments and pop function.
	args := state.frame().popN(state.frame().gettop() - fr.fnID + 1)[1:]

	// Enter and leave frame on return.
	defer state.leave(state.enter(fr))

	fr.pushN(args)

	// Is it a Lua closure?
	if fr.function().isLua() {
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
	} else if fr.function().isGo() {
		// Otherwise Go closure.
		if rets := fr.popN(fr.function().native(state)); fr.rets != 0 {
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
		return
	}
}

func (state *State) init(global *global) {
	state.frame().checkstack(InitialStackNew)
	state.global = global
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
	// fmt.Printf("%v[%v] (%t)\n", obj, key, raw)
	if tbl, ok := obj.(*table); ok {
		if val := tbl.get(key); !IsNone(val) || raw || IsNone(state.metafield(tbl, "__index")) {
			return val
		}
		val, err := tryMetaIndex(state, tbl, key)
		if err != nil {
			state.Errorf("%v", err)
		}
		return val
	}
	if !raw {
		val, err := tryMetaIndex(state, obj, key)
		if err != nil {
			state.Errorf("%v", err)
		}
		return val
	}
	return None
}

func (state *State) settable(obj, key, val Value, raw bool) {
	// fmt.Printf("%v[%v] = %v (%t)\n", obj, key, val, raw)
	if tbl, ok := obj.(*table); ok && (tbl.exists(key) || raw) {
		tbl.set(key, val)
		return
	}
	if !raw {
		if err := tryMetaNewIndex(state, obj, key, val); err != nil {
			state.Errorf("%v", err)
		}
	}
}

func (state *State) metafield(value Value, event string) Value {
	if obj := state.getmetatable(value, true); !IsNone(obj) {
		if tbl, ok := obj.(*table); ok && tbl != nil {
			return tbl.getStr(event)
		}
	}
	return None
}

func (state *State) setmetatable(value, meta Value) {
	mt, ok := meta.(*table)
	if !ok && !IsNone(meta) {
		state.errorf("metatable must be table or nil")
	}
	switch v := value.(type) {
	case *Object:
		v.meta = mt
	case *table:
		v.meta = mt
	default:
		state.global.builtins[v.Type()] = mt
	}
}

func (state *State) getmetatable(value Value, rawget bool) Value {
	var meta Value
	switch value := value.(type) {
	case *Object:
		if !IsNone(value.meta) && value.meta != nil {
			meta = value.meta
		}
	case *table:
		if !IsNone(value.meta) && value.meta != nil {
			meta = value.meta
		}
	default:
		if mt := state.global.builtins[value.Type()]; mt != nil {
			meta = mt
		}
	}
	if !rawget && !IsNone(meta) && meta != nil {
		if mt, ok := meta.(*table); ok {
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

// get resolves the value located at the acceptable index which may be valid and
// point to a stack position or pseudo-index which are used to access the registry
// and upvalues of function.
//
// Any function in the API that receives stack indices works only with valid
// or acceptable indices.
//
// A valid index is an index that refers to a position that stores a modifiable
// lua value (1 <= abs(index) <= top) and pseudo-indices, which represent some
// positions that are accessible to host code but that are not in the stack.
// Pseudo-indices are used to access the registry and the upvalues of a function.
//
// Acceptable indices serve to avoid extra tests against the stack top when querying
// the stack. For instance, a Go function can query its third argument without the
// need to first check wheter there is a third argument, that is, without the need
// to check whether 3 is a valid index.
//
// For functions that can be called with acceptable indices, any non-valid index is
// treated as if it contains a value of a virtual type "none", which behaves like a
// nil value.
func (state *State) get(index int) Value {
	switch frame := state.frame(); {
	//
	// Positive stack index
	//
	case index > 0:
		if index > cap(frame.locals) {
			state.errorf("unacceptable index (%d)", index)
		}
		if index > frame.gettop() {
			return None
		}
		return frame.get(index - 1)
	//
	// Negative stack index
	//
	case !isPseudoIndex(index):
		//state.Logf("get %d (absolute = %d)", index, frame.absindex(index))
		// Debug(state)
		if index = frame.absindex(index); index < 1 || index > frame.gettop() {
			state.errorf("invalid index (%d)", index)
		}
		return frame.get(index - 1)
	//
	// Registry pseudo index
	//
	case index == RegistryIndex:
		return state.global.registry
	//
	// Upvalues pseudo index
	//
	default:
		if index = RegistryIndex - index; index >= MaxUpValues {
			state.errorf("upvalue index too large (%d)", index)
		}
		if nups := len(frame.closure.upvals); nups == 0 || nups > index {
			return None
		}
		return frame.getUp(index - 1).get()
	}
}

func (state *State) set(index int, value Value) {
	switch frame := state.frame(); {
	//
	// Positive stack index
	//
	case index > 0:
		if index > cap(frame.locals) {
			state.errorf("unacceptable index (%d)", index)
		}
		if index > frame.gettop() {
			return
		}
		frame.set(index-1, value)
		return
	//
	// Negative stack index
	//
	case !isPseudoIndex(index):
		if index = frame.absindex(index); index < 1 || index > frame.gettop() {
			state.errorf("invalid index (%d)", index)
		}
		frame.set(index-1, value)
		return
	//
	// Registry pseudo index
	//
	case index == RegistryIndex:
		state.global.registry = value.(*table)
		return

	//
	// Upvalues pseudo index
	//
	default:
		if index = RegistryIndex - index; index >= MaxUpValues {
			state.errorf("upvalue index too large (%d)", index)
		}
		if nups := len(frame.closure.upvals); nups == 0 || nups > index {
			return
		}
		frame.setUp(index-1, value)
		return
	}
}

// converts an integer to a "floating point byte", represented as (eeeeexxx), where the real
// value is (1xxx) * 2^(eeeee - 1) if eeeee != 0 and (xxx) otherwise.
func i2fb(i int) int {
	var (
		u     = uint8(i)
		e int = 0 // exponent
	)
	if u < 8 {
		return i
	}
	for u >= (8 << 4) { // coarse steps
		u = (u + 0xF) >> 4 // x = ceil(x/16)
		e += 4
	}
	for u >= (8 << 1) { // fine steps
		u = (u + 1) >> 1 // x = ceil(x/2)
		e++
	}
	return ((e + 1) << 3) | (int(i) - 8)
}

// converts a "floating point byte" to an integer.
func fb2i(i int) int {
	if i < 8 {
		return i
	}
	return ((i & 7) + 8) << ((uint8(i) >> 3) - 1)
}

// When a function is created, it is possible to associate some values with it, thus creating
// a closure (see PushClosure); these values are called upvalues and are accessible to the
// function whenever it is called.
//
// Whenever a function is called, its upvalues are located at specific pseudo-indices.
// These pseudo-indices are produced by the macro UpValueIndex.
//
// The first upvalue associated with a function is at index UpValueIndex(1), and so on.
// Any access to UpValueIndex(n), where n is greater than the number of upvalues of the
// current function (but not greater than 256, which is one plus the maximum number of
// upvalues in a closure), produces an acceptable but invalid index.
func UpValueIndex(index int) int { return RegistryIndex - index }

// IsUpValueIndex reports true if the index represents an upvalue index.
func isUpValueIndex(index int) bool { return index < RegistryIndex }

// IsStackIndex reports true if the index represents a stack index.
//
// Tests for valid but not pseudo index.
func isStackIndex(index int) bool { return !isPseudoIndex(index) }

// isPseudoIndex reports whether the Index index represents a pseudo-index; that is, an index
// that represents registers that are accessible to host code but that are not in the stack.
//
// Pseudo-indices are used to access the registry and the upvalues of a function.
func isPseudoIndex(index int) bool { return index <= RegistryIndex }
