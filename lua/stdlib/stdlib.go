package stdlib

import (
	"github.com/Azure/golua/lua"
)

func Import(state *lua.State) {
    var stdlibs = map[string]lua.Func{
        "_G":        lua.Func(OpenBase),
        // "package":   lua.Func(Package),
        // "coroutine": lua.Func(Coroutine),
        // "table":     lua.Func(Table),
        // "io":        lua.Func(IO),
        // "os":        lua.Func(OS),
        // "string":    lua.Func(String),
        // "math":      lua.Func(Math),
        // "utf8":      lua.Func(UTF8),
        // "debug":     lua.Func(Debug),
    }
    for name, load := range stdlibs {
        state.Require(name, load, true)
        state.Pop(1)
    }
}