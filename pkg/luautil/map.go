package luautil

import (
	"reflect"

	"github.com/Azure/golua/lua"
)

func mapNewIndex(state *lua.State) int {
	if state.IsNoneOrNil(2) {
		return 0
	}

	check(state, 1).SetMapIndex(
		reflect.ValueOf(state.CheckString(2)),
		toGoValue(state.CheckAny(3)),
	)
	return 0
}

func mapIndex(state *lua.State) int {
	var (
		key = state.CheckString(2)
		obj = check(state, 1)
	)
	if val := obj.MapIndex(reflect.ValueOf(key)); val.IsValid() {
		state.Push(valueOf(state, val.Interface()))
	} else {
		state.Push(nil)
	}
	return 1
}

func mapLength(state *lua.State) int {
	state.Push(check(state, 1).Len())
	return 1
}

func valueFromMap(state *lua.State, rv reflect.Value) *lua.Object {
	var mapMetaFuncs = map[string]lua.Func{
		"__newindex": lua.Func(mapNewIndex),
		"__index":    lua.Func(mapIndex),
		"__len":      lua.Func(mapLength),
	}
	state.Push(rv.Interface())
	state.NewTableSize(0, len(mapMetaFuncs))
	state.SetFuncs(mapMetaFuncs, 0)
	state.SetMetaTableAt(-2)
	return state.ToUserData(-1)
}
