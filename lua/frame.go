package lua

import (
    "fmt"
    "os"
    
    "github.com/Azure/golua/lua/ir"
)

var _ = fmt.Println
var _ = os.Exit

type (
    // CallInfo holds information about a call.
    CallInfo struct {
        prev, next *CallInfo // dynamic call link to caller and callee
        savedpc    int       // saved pc; return address (lua only)
        status     int       // call status indicating success or failure
        funcID     int       // function index in the stack
        nrets      int       // expected number of results
        top        int       // stack top for this function
        frame      *Frame    // call frame
    }

    // Frame is the context to execute a function closure.
    Frame struct {
        prev, next *Frame // dynamic link caller and callee frame
        closure *Closure  // frame closure
        vararg  []Value   // variable arguments
        locals  []Value   // frame stack locals 
        upvals  []int     // list of open upvalue indices
        state   *State    // thread state
        depth   int       // call frame ID
        fnID    int       // function index
        rets    int       // # expected returns
        //top     int       // next free slot in frame stack
        pc      int       // last executed instruction pc
    }
)

// checkstack checks that there atleast needed slots available.
func (fr *Frame) checkstack(needed int) bool {
    if space := cap(fr.locals) - fr.gettop(); space < needed {
        fr.extend(needed - space)
    }
    return true
}

// extend grows the frame's locals stack by grow.
func (fr *Frame) extend(grow int) {
    fr.locals = append(
        fr.locals[:len(fr.locals)],
        make([]Value, 0, grow)...,
    )
}

// absindex converts the acceptable index into an equivalent
// absolute index; that is, one that does not depend on the
// stack top).
func (fr *Frame) absindex(index int) int {
    // zero, positive, or pseudo index
    if index > 0 || isPseudoIndex(index) {
        return index
    }
    // negative
    return fr.gettop() + index + 1
}

// Reverse reverses the frame's locals stack starting from the src to dst indices.
func (fr *Frame) reverse(src, dst int) {
    for locals := fr.locals; src < dst; {
        locals[src], locals[dst] = locals[dst], locals[src]
        src++
        dst--
    }
}

// Replace moves the top element into the given valid index without shifting
// any element (therefore replacing the value at that given index), and then
// pops the top element.
func (fr *Frame) replace(index int) {
    if v := fr.pop(); fr.gettop() == 0 {
        fr.push(v)
    } else {
        fr.set(index, v)
    }
}

// Rotate rotates the stack elements between the valid index and the top of the stack.
//
// The elements are rotated n positions in the direction of the top, if positive;
// otherwise -n positions in the direction of the bottom, if negative.
//
// The absolute value of n must not be greater than the size of the slice being rotated.
//
// This function cannot be called with a pseudo-index, because a pseudo-index is not an
// actual stack position.
//
// Let x = AB, where A is a prefix of length 'n'.
// Then, rotate x n == BA. But BA == (A^r . B^r)^r
//
// See https://www.lua.org/manual/5.3/manual.html#lua_rotate
func (fr *Frame) rotate(index, n int) {
    var (
        abs = fr.absindex(index) - 1 // start of stack segment
        top = fr.gettop() - 1        // top of stack segment
        mid int
    )
    if n >= 0 { mid = top - n } else { mid = abs - n - 1 }
    fr.reverse(abs, mid)     // reverse prefix with length n
    fr.reverse(mid+1, top)   // reverse suffix
    fr.reverse(abs, top)     // reverse segment
}

// Remove removes the element at the given valid index, shifting down the elements
// above this index to fill the gap.
//
// This function cannot be called with a pseudo-index, because a pseudo-index
// is not an actual stack position.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_remove
func (fr *Frame) remove(index int) {
    fr.rotate(fr.absindex(index), -1)
    fr.pop()
}

// Insert moves the top element into the given valid index, shifting up the
// elements above this index to open space.
//
// This function cannot be called with a pseudo-index, because a pseudo-index
// is not an actual stack position.
//func (fr *Frame) insert(index int) { fr.rotate(index, 1) }

// caller returns the frame's caller frame.
func (fr *Frame) caller() *Frame {
    if fp := fr.prev; fr.state != nil && fp != &fr.state.base {
        return fp
    }
    return nil
}

// callee returns the frame's callee frame (if any).
func (fr *Frame) callee() *Frame {
    if fp := fr.next; fr.state != nil && fp != &fr.state.base {
        return fp
    }
    return nil
}

// upvalue returns the upvalue at index.
func (fr *Frame) upvalue(index int) Value {
    return *fr.closure.upvals[index]
}

// gettop returns the index of the top element in the stack.
//
// Because indices start at 1, this result is equal to the
// number of elements in the stack; in particular, 0 means
// an empty stack.
func (fr *Frame) gettop() int { return len(fr.locals) }

// local returns the n'th local in the frame's stack.
//
// TODO: bounds check
func (fr *Frame) local(index int) Value {
     if index = fr.absindex(index); fr.instack(index) {
        return fr.locals[index-1]
    }
    return None
}

// pushN pushes N values onto the frame's stack.
//
// TODO: ensure stack
func (fr *Frame) pushN(vs []Value) {
    for _, v := range vs {
        fr.push(v)
    }
}

// Copy copies the element at index src into the valid index dst,
// replacing the value at that position.
//
// Values at other positions are not affected.
//
// TODO: bounds check
func (fr *Frame) copy(src, dst int) { fr.set(src, fr.get(dst)) }

// push pushes 1 values onto the frame's stack.
//
// TODO: ensure stack
func (fr *Frame) push(v Value) {
    fr.locals = append(fr.locals, v)
}

// pop pops 1 value from the frame's stack.
//
// TODO: ensure stack
func (fr *Frame) pop() Value {
    if fr.gettop() == 0 {
        return None
    }
    top := fr.gettop() - 1
    val := fr.locals[top]
    fr.locals = fr.locals[:top]
    return val
}

// popN pops N values from the frame's stack.
//
// TODO: ensure stack
func (fr *Frame) popN(n int) (vs []Value) {
    vs = make([]Value, n)
    for i := n-1; i >= 0; i-- {
        vs[i] = fr.pop()
    }
    return vs
}

// step returns the next instruction for the frame's
// current instruction pointer and increments by n.
//
// TODO: bounds check
func (fr *Frame) step(n int) ir.Instr {
    i := ir.Instr(fr.closure.proto.Code[fr.pc])
    fr.pc += n
    return i
}

// set sets the frame local value at index to value.
//
// TODO: pseudo & upvalue indices.
// TODO: bounds and stack check.
func (fr *Frame) set(index int, value Value) {
    // if index == fr.gettop() {
    //     fr.push(value)
    // } else {
    //     fr.locals[index-1] = value
    // }
    if fr.gettop() == index {
        fr.push(value)
        return
    }
    if index = fr.absindex(index); fr.instack(index) {
        fr.locals[index-1] = value
    }
}

// get returns the value located in the frame's locals
// stack, or at pseudo- & upvalue- index.
//
// TODO: pseudo & upvalue indices.
// TODO: bounds and stack check.
func (fr *Frame) get(index int) Value {
    if index = fr.absindex(index); fr.instack(index) {
        return fr.locals[index-1]
    }
    return None
}

// instack reports whether index is in the frame's locals stack.
func (fr *Frame) instack(index int) bool {
    return index > 0 && index <= fr.gettop()
}