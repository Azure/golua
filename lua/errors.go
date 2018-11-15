package lua

import (
	"fmt"
)

func argError(state *State, argAt int, msg string) {
	// TODO: stack analysis and debugging info if available.
	panic(fmt.Errorf("bad argument #%d (%s)", argAt, msg))
}

func intError(state *State, argAt int) {
	if isNumber(state.get(argAt)) {
		argError(state, argAt, "number has not ineger representation")
	}
	typeError(state, argAt, "number")
}

func typeError(state *State, argAt int, want string) {
	// TODO: stack analysis and debugging info if available.
	panic(fmt.Errorf("%s expected @ %d, got %s", want, argAt, state.valueAt(argAt).Type()))
}

// luaG_typerror 		"attempt to %s a %s value%s"
// luaG_concaterror 	typeerror()
// luaG_opinterror
// luaG_tointerror
// luaG_ordererror
// luaG_runerror

// https://www.lua.org/manual/5.3/manual.html#lua_error
// https://www.lua.org/manual/5.3/manual.html#luaL_error
// https://www.lua.org/manual/5.3/manual.html#luaL_argcheck
// https://www.lua.org/manual/5.3/manual.html#luaL_argerror
