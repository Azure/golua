package lua

import (
	"fmt"
	"io"
)

// CheckStack ensures that the stack has space for at least n extra slots (that is, that you can
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

// DumpStack writes the Lua thread's current stack frame to w.
func (state *State) DumpStack(w io.Writer) {
    if fr := state.frame(); fr != nil {
        for i := fr.gettop() - 1; i >= 0; i-- {
            fmt.Fprintf(w, "[%d] %v\n", i+1, fr.locals[i])
        }
    }
}

// Replace moves the top element into the given valid index without shifting any element
// (therefore replacing the value at that given index), and then pops the top element.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_replace
func (state *State) Replace(index int) { state.frame().replace(index) }

// rotate rotates the stack elements between the valid index and the top of the stack.
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
func (state *State) Rotate(index, n int) { state.frame().rotate(index, n) }

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

// XMove exchanges values between different threads of the same state.
//
// This function pops N values from the state's stack, and pushes them onto
// the state dst's stack.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_xmove
func (state *State) XMove(dst *State, n int) {
    dst.frame().pushN(state.frame().popN(n))
}

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

// PushGlobals pushes the global environment onto the stack.
//
// See https://www.lua.org/manual/5.3/manual.html#lua_pushglobaltable
func (state *State) PushGlobals() {
    state.Push(state.global.registry.getInt(GlobalsIndex))
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