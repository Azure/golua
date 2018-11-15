package coro

import (
	"fmt"

	"github.com/Azure/golua/lua"
)

//
// Lua Standard Library -- coroutine
//

func Open(state *lua.State) int {
	// Create 'coroutine' table.
	var coroutineFuncs = map[string]lua.Func{
		"create":      lua.Func(coroutineCreate),
		"resume":      lua.Func(coroutineResume),
		"running":     lua.Func(coroutineRunning),
		"status":      lua.Func(coroutineStatus),
		"wrap":        lua.Func(coroutineWrap),
		"yield":       lua.Func(coroutineYield),
		"isyieldable": lua.Func(coroutineIsYieldable),
	}
	state.NewTableSize(0, 7)
	state.SetFuncs(coroutineFuncs, 0)

	// Return 'coroutine' table.
	return 1
}

func coroutineCreate(state *lua.State) int {
	unimplemented("coroutine")
	return 0
}

func coroutineResume(state *lua.State) int {
	unimplemented("coroutine")
	return 0
}

func coroutineRunning(state *lua.State) int {
	unimplemented("coroutine")
	return 0
}

func coroutineStatus(state *lua.State) int {
	unimplemented("coroutine")
	return 0
}

func coroutineWrap(state *lua.State) int {
	unimplemented("coroutine")
	return 0
}

func coroutineYield(state *lua.State) int {
	unimplemented("coroutine")
	return 0
}

func coroutineIsYieldable(state *lua.State) int {
	unimplemented("coroutine")
	return 0
}

func unimplemented(msg string) { panic(fmt.Errorf(msg)) }
