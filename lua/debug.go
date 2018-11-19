package lua

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
)

// debug is a structure used to carry different pieces of information about a function
// or an activation record. state.GetStack fills only the private fields of this struct,
// for later use. To fill the other fields, use state.GetInfo.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_Debug
type Debug struct {
	source   string
	short    string
	name     string
	kind     string
	what     string
	span     [2]int
	nups     int
	active   int
	params   int
	vararg   bool
	tailcall bool
}

func (debug *Debug) Source() string       { return debug.source }
func (debug *Debug) ShortSrc() string     { return debug.short }
func (debug *Debug) CurrentLine() int     { return debug.active }
func (debug *Debug) LineDefined() int     { return debug.span[0] }
func (debug *Debug) LastLineDefined() int { return debug.span[1] }
func (debug *Debug) What() string         { return debug.what }
func (debug *Debug) NumUps() int          { return debug.nups }
func (debug *Debug) NumParams() int       { return debug.params }
func (debug *Debug) IsVararg() bool       { return debug.vararg }
func (debug *Debug) Name() string         { return debug.name }
func (debug *Debug) NameWhat() string     { return debug.kind }
func (debug *Debug) IsTailCall() bool     { return debug.tailcall }

// GetStack returns debug information about the interpreter runtime stack.
//
// This function fills parts of the Debug structure with an identification of
// the activation record of the function executing at a given level. Level 0
// is the current running function, whereas level n+1 is the function that has
// called level n (except for tail calls, which do not count on the stack).
// When there are no errors, StackDepth returns the Debug structure and nil;
// Otherwise returns nil and any error.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_getstack
func (state *State) GetStack(debug *Debug, depth int) error {
	fmt.Println("state.GetStack(): TODO")
	return nil
}

// DebugInfo returns debug information about a specific function or function invocation.
//
// To get information about a function invocation, the parameter ar must be a valid
// activation record that was filled by a previous call to lua_getstack or given as
// argument to a hook (see Hook).
//
// To get information about a function, you push it onto the stack and start the what
// string with the character '>'. (In that case, lua_getinfo pops the function from the
// top of the stack.) For instance, to know in which line a function f was defined, you
// can write the following code:
//
//  lua_Debug ar;
//  lua_getglobal(L, "f");  /* get global 'f' */
//  lua_getinfo(L, ">S", &ar);
//  printf("%d\n", ar.linedefined);
//
// Each character in the string what selects some fields of the structure ar to be filled
// or a value to be pushed on the stack:
//
//  'n': fills in the field name and namewhat;
//  'S': fills in the fields source, short_src, linedefined, lastlinedefined, and what;
//  'l': fills in the field currentline;
//  't': fills in the field istailcall;
//  'u': fills in the fields nups, nparams, and isvararg;
//  'f': pushes onto the stack the function that is running at the given level;
//  'L': pushes onto the stack a table whose indices are the numbers of the lines that are
//       valid on the function. (A valid line is a line with some associated code, that is,
//       a line where you can put a break point. Non-valid lines include empty lines and
//       comments.)
//
// If this option is given together with option 'f', its table is pushed after the function.
//
// This function returns 0 on error (for instance, an invalid option in what).
//
// See https://www.lua.org/manual/5.3/manual.html#lua_getinfo
func (state *State) GetInfo(debug *Debug, options string) error {
	if len(options) > 0 && options[0] == '>' {
		if cls, ok := state.frame().pop().(*Closure); ok {
			return state.getInfo(state.frame(), debug, cls, options[1:])
		}
		return fmt.Errorf("function expected")
	}
	return fmt.Errorf("state.GetInfo(): TODO")
}

func (state *State) getInfo(frame *Frame, debug *Debug, closure *Closure, options string) error {
	for pos := 0; pos < len(options); pos++ {
		switch b := options[pos]; b {
		case 'S':
			funcinfo(frame, debug, closure)
		case 'l':
			if debug.active = -1; closure.isLua() && frame.pc > 0 {
				currentline := int(closure.binary.PcLnTab[frame.pc])
				debug.active = currentline
			}
		case 'u':
			if !closure.isLua() {
				debug.vararg = true
				debug.params = 0
			} else {
				debug.vararg = closure.binary.IsVararg()
				debug.params = closure.binary.NumParams()
			}
		case 't':
			debug.tailcall = frame.status&callStatusTail != 0
		case 'n':
			name, kind := funcname(frame, closure)
			debug.name = name
			debug.kind = kind
		case 'L', 'f':
			// TODO
		default:
			return fmt.Errorf("invalid option: %c", b)
		}
	}
	return nil
}

// GetUpValue gets information about the n-th upvalue of the closure at index funcindex.
// It pushes the upvalue's value onto the stack and returns its name. Returns NULL (and
// pushes nothing) when the index n is greater than the number of upvalues.
//
// For C functions, this function uses the empty string "" as a name for all upvalues.
// For Lua functions, upvalues are the external local variables that the function uses,
// and that are consequently included in its closure.
//
// Upvalues have no particular order, as they are active through the whole function.
// They are numbered in an arbitrary order.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_getupvalue
func (state *State) GetUpValue(function, index int) (name string) {
	if cls, ok := state.get(function).(*Closure); ok {
		if index <= len(cls.upvals) {
			up := cls.getUp(index - 1)
			state.Push(up.get())
			name = cls.upName(index - 1)
		}
	}
	return
}

// UpValueID returns a unique identifier for the upvalue numbered n from the closure at
// index function.
//
// These unique identifiers allow a program to check whether different closures share
// upvalues. Lua closures that share an upvalue (that is, that access a same external
// local variable) will return identical ids for those upvalue indices.
//
// Parameters funcindex and n are as in function lua_getupvalue, but n cannot be greater
// than the number of upvalues.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_upvalueid
func (state *State) UpValueID(function, index int) interface{} {
	if cls, ok := state.get(function).(*Closure); ok {
		if index <= len(cls.upvals) {
			return cls.getUp(index - 1)
		}
	}
	panic(fmt.Errorf("closure expected"))
}

// UpValueJoin makes the n1-th upvalue of the Lua closure at index func1 refer to the n2-th
// upvalue of the Lua closure at index func2.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_upvaluejoin
func (state *State) UpValueJoin(func1, n1, func2, n2 int) {
	if fn1, ok1 := state.get(func1).(*Closure); ok1 {
		if fn2, ok2 := state.get(func2).(*Closure); ok2 {
			if (len(fn1.upvals) >= n1) && (len(fn2.upvals) >= n2) {
				fn1.upvals[n1-1] = fn2.upvals[n2-1]
			}
		}
	}
}

type HookEvent uint

const (
	HookCall HookEvent = 1 << iota
	HookRets
	HookLine
	HookCount
	HookTailCall
)

func (evt HookEvent) String() string {
	var s string
	if evt&HookCall != 0 {
		s += "call"
		s += ""
	}
	if evt&HookRets != 0 {
		s += "return"
		s += ""
	}
	if evt&HookLine != 0 {
		s += "line"
		s += ""
	}
	if evt&HookCount != 0 {
		s += "count"
		s += ""
	}
	if evt&HookTailCall != 0 {
		s += "tail call"
	}
	return s
}

// The hook table at registry[HookKey] maps threads to their current hook function.
const HookKey = 0

func (state *State) Debug(halt bool) {
	DBG(state.frame(), halt)
}

func DBG(fr *Frame, halt bool) {
	const base = 0

	var b strings.Builder

	var pcln string
	if fr.closure.isLua() {
		pcln = fmt.Sprintf("@line = %d", fr.closure.binary.PcLnTab[fr.pc])
	}
	fmt.Fprintf(&b, "\nframe#%d <prev=%p|next=%p> %s\n", fr.depth, fr.prev, fr.next, pcln)
	fmt.Fprintf(&b, "    %s", fr.closure)
	if fr.closure != nil {
		fmt.Fprintf(&b, " (# up = %d)", len(fr.closure.upvals))
	}
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "            * savedpc = %d\n", fr.pc)
	fmt.Fprintf(&b, "            * returns = %d\n\n", fr.rets)

	fmt.Fprintf(&b, "            upvalues\n")
	for i, upval := range fr.function().upvalues() {
		fmt.Fprintf(&b, "                [%d] %v\n", i, *upval)
	}
	fmt.Fprintf(&b, "            end\n\n")

	fmt.Fprintf(&b, "            varargs\n")
	for i, extra := range fr.vararg {
		fmt.Fprintf(&b, "                [%d] %v\n", i, extra)
	}
	fmt.Fprintf(&b, "            end\n")
	fmt.Fprintf(&b, "    end\n\n")
	fmt.Fprintf(&b, "    locals (len=%d, cap=%d, top=%d)\n", len(fr.locals), cap(fr.locals), fr.gettop())
	for i := fr.gettop() - 1; i >= 0; i-- {
		fmt.Fprintf(&b, "        [%d] %v\n", i+base, fr.locals[i])
	}

	fmt.Fprintf(&b, "    end\n")
	fmt.Fprintf(&b, "end\n")

	fmt.Println(b.String())
	if halt {
		os.Exit(1)
	}
}

func funcname(frame *Frame, closure *Closure) (name, what string) {
	if !closure.isLua() {
		pc := reflect.ValueOf(closure.native).Pointer()
		fn := runtime.FuncForPC(pc)
		return fn.Name(), ""
	}
	// TODO
	return "", ""
}

func funcinfo(frame *Frame, debug *Debug, closure *Closure) {
	if closure.isLua() {
		proto := closure.binary
		if debug.source = "=?"; proto.Source != "" {
			debug.source = proto.Source
		}
		debug.span[0] = int(proto.SrcPos)
		debug.span[1] = int(proto.EndPos)
		if debug.span[1] == 0 {
			debug.what = "main"
		} else {
			debug.what = "Lua"
		}
	} else {
		debug.source = "=[Go]"
		debug.span[0] = -1
		debug.span[1] = -1
		debug.what = "Go"
	}
	// TODO: short
	debug.short = chunkID(debug.source)
}

func chunkID(source string) string {
	const maxLen = 60

	if len(source) > 0 {
		switch source[0] {
		case '=':
			if source = source[1:]; len(source) > maxLen {
				source = source[:maxLen-1]
			}
		case '@':
			if source = source[1:]; len(source) > maxLen {
				source = fmt.Sprintf("...%s", source[len(source)-maxLen+4:])
			}
		default:
			if i := strings.IndexByte(source, '\n'); i != -1 {
				source = source[0:i] + "..."
			}
			if max := maxLen - len(`[string " "]`); len(source) > max {
				source = source[0:max-3] + "..."
			}
			source = fmt.Sprintf(`[string "%s"]`, source)
		}
	}
	return source
}
