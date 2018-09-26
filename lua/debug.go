package lua

import (
	"strings"
	"fmt"
	"os"
)

// GetUpValue gets information about the n-th upvalue of the closure at index funcindex.
// It pushes the upvalue's value onto the stack and returns its name. Returns NULL (and
// pushes nothing) when the index n is greater than the number of upvalues.
//
// For C functions, this function uses the empty string "" as a name for all upvalues.
// For Lua functions, upvalues are the external local variables that the function uses,
// and that are consequently included in its closure.
//
// Upvalues have no particular order, as they are active through the whole function. They are numbered in an arbitrary order.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_getupvalue
func (state *State) GetUpValue(function, index int) (name string) {
    if cls, ok := state.get(function).(*Closure); ok {
        if index < len(cls.upvals) {
            up := cls.getUp(index-1)
            state.Push(up.get())
            name = up.ident
        }
    }
    return
}

type Debug interface {
    State() *State
}

type HookFunc func(Debug)

type HookEvent uint

const (
    HookCall     HookEvent = 1 << iota
    HookRets
    HookLine
    HookCount
    HookTailCall
)

func (evt HookEvent) String() string {
    var s string
    if evt & HookCall != 0 {
        s += "call"
        s += ""
    }
    if evt & HookRets != 0 {
        s += "return"
        s += ""
    }
    if evt & HookLine != 0 {
        s += "line"
        s += ""
    }
    if evt & HookCount != 0 {
        s += "count"
        s += ""
    }
    if evt & HookTailCall != 0 {
        s += "tail call"
    }
    return s
}

// The hook table at registry[HookKey] maps threads to their current hook function.
const HookKey = 0

type debug struct {
    state *State
}

func (dbg *debug) State() *State { return dbg.state }

func (state *State) Debug(halt bool) {
    DBG(state.frame(), halt)
}

func DBG(fr *Frame, halt bool) {
    const base = 0

	 var b strings.Builder

	//fr := state.frame()
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
    for i, upval := range fr.closure.upvals {
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