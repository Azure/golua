package stdlib

import (
    "fmt"

	"github.com/Azure/golua/lua"
)

func Load(state *lua.State) {
    var stdlibs = []struct{ Name string; Load lua.Func }{
        {"_G",        lua.Func(OpenBase)},
        {"package",   lua.Func(OpenPackage)},
        {"coroutine", lua.Func(OpenCoroutine)},
        {"table",     lua.Func(OpenTable)},
        {"io",        lua.Func(OpenIO)},
        {"os",        lua.Func(OpenOS)},
        {"string",    lua.Func(OpenString)},
        {"math",      lua.Func(OpenMath)},
        {"utf8",      lua.Func(OpenUTF8)},
        {"debug",     lua.Func(OpenDebug)},
    }
    for _, stdlib := range stdlibs {
        state.Logf("open stdlib module %q", stdlib.Name)
        lua.Require(state, stdlib.Name, stdlib.Load, true)
        state.Pop()
    }
}

func unimplemented(msg string) { panic(fmt.Errorf(msg)) }