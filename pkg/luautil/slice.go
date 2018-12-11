package luautil

import (
	"reflect"

	"github.com/Azure/golua/lua"
)

func sliceNewIndex(state *lua.State) int {
	var (
		slice = check(state, 1)
		index = int(state.CheckInt(2))
		value = toGoValue(state.CheckAny(3))
	)
	if index < 1 || index > slice.Len() {
		state.ArgError(2, "index out of range")
	}
	slice.Index(index-1).Set(value)
	return 0
}

func sliceIndex(state *lua.State) int {
	var (
		index = int(state.CheckInt(2))
		slice = check(state, 1)
	)
	if index < 1 || index > slice.Len() {
		state.ArgError(2, "index out of range")
	}
	state.Push(valueOf(state,  slice.Index(index-1).Interface()))
	return 1
}

func sliceLength(state *lua.State) int {
	state.Push(check(state, 1).Len())
	return 1
}

func valueFromSlice(state *lua.State, rv reflect.Value) *lua.Object {
	var sliceMetaFuncs = map[string]lua.Func{
		"__newindex":  lua.Func(sliceNewIndex),
		"__index":    lua.Func(sliceIndex),
		"__len":	  lua.Func(sliceLength),
	}
	state.Push(rv.Interface())
	state.NewTableSize(0, len(sliceMetaFuncs))
	state.SetFuncs(sliceMetaFuncs, 0)
	state.SetMetaTableAt(-2)
	return state.ToUserData(-1)
}