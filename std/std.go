package std

import (
	"github.com/Azure/golua/std/base"
    "github.com/Azure/golua/std/coro"
    "github.com/Azure/golua/std/debug"
    "github.com/Azure/golua/std/io"
    "github.com/Azure/golua/std/math"
    "github.com/Azure/golua/std/os"
    "github.com/Azure/golua/std/pkg"
    "github.com/Azure/golua/std/str"
    "github.com/Azure/golua/std/table"
    "github.com/Azure/golua/std/utf8"

    "github.com/Azure/golua/lua"
)

func Open(state *lua.State) {
    var libs = []struct{ Name string; Open lua.Func }{
        {"_G",        lua.Func(base.Open)},
        {"package",   lua.Func(pkg.Open)},
        {"coroutine", lua.Func(coro.Open)},
        {"table",     lua.Func(table.Open)},
        {"io",        lua.Func(io.Open)},
        {"os",        lua.Func(os.Open)},
        {"string",    lua.Func(str.Open)},
        {"math",      lua.Func(math.Open)},
        {"utf8",      lua.Func(utf8.Open)},
        {"debug",     lua.Func(debug.Open)},
    }
    for _, lib := range libs {
        state.Logf("open stdlib module %q", lib.Name)
        lua.Require(state, lib.Name, lib.Open, true)
        state.Pop()
    }
}