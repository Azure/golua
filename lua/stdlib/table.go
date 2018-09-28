package stdlib

import (
	"github.com/Azure/golua/lua"
)

//
// Lua Standard Library -- table
//

func OpenTable(state *lua.State) int {
	// Create 'table' table.
	var tableFuncs = map[string]lua.Func{
 		"concat": lua.Func(tableConcat),
		"insert": lua.Func(tableInsert),
		"pack":   lua.Func(tablePack),
		"unpack": lua.Func(tableUnpack),
		"remove": lua.Func(tableRemove),
		"move":   lua.Func(tableMove),
		"sort":   lua.Func(tableSort),
	}
	state.NewTableSize(0, 7)
	state.SetFuncs(tableFuncs, 0)

	// Return 'table' table.
	return 1
}

func tableConcat(state *lua.State) int {
	unimplemented("tableConcat")
	return 0
}

func tableInsert(state *lua.State) int {
	unimplemented("tableInsert")
	return 0
}

func tablePack(state *lua.State) int {
	unimplemented("tablePack")
	return 0
}

func tableUnpack(state *lua.State) int {
	unimplemented("tableUnpack")
	return 0
}

func tableRemove(state *lua.State) int {
	unimplemented("tableRemove")
	return 0
}

func tableMove(state *lua.State) int {
	unimplemented("tableMove")
	return 0
}

func tableSort(state *lua.State) int {
	unimplemented("tableSort")
	return 0
}
