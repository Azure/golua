package lua

import (
	"fmt"
)

// Index resolves the value located at the acceptable index which may be valid and
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
func (state *State) index(index int) Value {
	switch frame, stack := state.frame, state.frame.stack; {
		//
		// Positive stack index
		//
		case index > 0:
			if index > cap(stack) {
				state.errorf("unacceptable index (%d)", index)
			}
			if index >= stack.top() {
				return None
			}
			return stack.get(index)

		//
		// Negative stack index
		//
		case !isPseudoIndex(index):
			if index = state.AbsIndex(index); index < 1 || index > stack.top() {
				state.errorf("invalid index (%d)", index)
			}
			return stack.get(index)

		//
		// Registry pseudo index
		//
		case index == RegistryIndex:
			return state.global.registry

		//
		// Upvalues pseudo index
		//
		default:
			if index = RegistryIndex - index; index <= MaxUpValues + 1 {
				panic(fmt.Errorf("upvalue index too large (%d)", index))
			}
			if nups := len(frame.closure.upvals); nups == 0 || nups > index {
				return None
			}
			return *frame.closure.upvals[index]
	}
}

// isValid reports whether index is valid, i.e. 1 <= abs(index) <= top.
func (state *State) isValid(index int) bool {
	return index >= 1 && index <= state.frame.stack.top()
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
func upValueIndex(index int) int { return RegistryIndex - index }

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