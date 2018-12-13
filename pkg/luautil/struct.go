package luautil

import (
	"reflect"
	"strings"

	"github.com/Azure/golua/lua"
)

func structIndex(state *lua.State) int {
	var (
		field = state.CheckString(2)
		struc = check(state, 1)
	)
	rv := struc.Elem().FieldByName(strings.Title(field))
	if !rv.IsValid() {
		state.Push(nil)
	} else {
		state.Push(ValueOf(state, rv.Interface()))
	}
	return 1
}

func valueFromStruct(state *lua.State, rv reflect.Value) *lua.Object {
	var structMetaFuncs = map[string]lua.Func{
 		"__index":    lua.Func(structIndex),
	}
	state.Push(rv.Interface())
	state.NewTableSize(0, len(structMetaFuncs))
	state.SetFuncs(structMetaFuncs, 0)
	state.SetMetaTableAt(-2)
	return state.ToUserData(-1)
}