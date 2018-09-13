package stdlib

import (
    "fmt"

	"github.com/Azure/golua/lua"
)

// TODO: LibTable  (Table, OpenTable)
// TODO: LibCoro   (Coro/Coroutine, OpenCoro)
// TODO: LibIO     (IO, OpenIO)
// TODO: LibOS     (OS, OpenOS)
// TODO: LibLoad   (Package, OpenPackage)
// TODO: LibMath   (Math, OpenMath)
// TODO: LibUTF8   (UTF8, OpenUTF8)
// TODO: LibDebug  (Debug, OpenDebug)
// TODO: LibString (String, OpenString)
func Import(state *lua.State) {
    var stdlibs = []struct{ Name string; Load lua.Func }{
        {"_G",      lua.Func(OpenBase)},
        {"package", lua.Func(OpenLoad)},
    }
    for _, stdlib := range stdlibs {
        state.Logf("open stdlib module %q", stdlib.Name)
        lua.Require(state, stdlib.Name, stdlib.Load, true)
        state.Pop()
    }
}

func unimplemented(msg string) { panic(fmt.Errorf(msg)) }