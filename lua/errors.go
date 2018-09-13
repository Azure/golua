package lua

import (
	"fmt"
)

func argErrorf(state *State, argAt int, format string, args ...interface{}) {
	// TODO: stack analysis and debugging info if available.
	msg := fmt.Sprintf(format, args...)
	argError(state, argAt, msg)
}

func argError(state *State, argAt int, msg string) {
	// TODO: stack analysis and debugging info if available.
	state.Errorf("bad argument #%d (%s)", argAt, msg)
}

func typeError(state *State, argAt int, want Type) {
	// TODO: stack analysis and debugging info if available.
	state.Errorf("%s expected @ %d, got %s", want, argAt, state.Value(argAt).Type())
}