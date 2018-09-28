package lua

import (
	"fmt"
	"os"
)

var (
	_ = fmt.Println
	_ = os.Exit
)

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
			return frame.get(index-1)
		//
		// Negative stack index
		//
		case !isPseudoIndex(index):
			//state.Logf("get %d (absolute = %d)", index, frame.absindex(index))
			// Debug(state)
			if index = frame.absindex(index); index < 1 || index > frame.gettop() {
				state.errorf("invalid index (%d)", index)
			}
			return frame.get(index-1)
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
			return frame.getUp(index-1).get()
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
			state.global.registry = value.(*Table)
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
		u = uint8(i)
		e int = 0 // exponent
	)
	if u < 8 {
		return i
	}
	for u >= (8 << 4) {    // coarse steps
		u = (u + 0xF) >> 4 // x = ceil(x/16)
		e += 4
	}
	for u >= (8 << 1) {  // fine steps
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