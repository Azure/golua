package stdlib

import (
    "fmt"
    "os"
    "github.com/Azure/golua/lua"
)

var _ = fmt.Println
var _ = os.Exit

//
// Lua Standard Library -- Base
//

// collectgarbage
// xpcall
func OpenBase(state *lua.State) int {
    var funcs = map[string]lua.Func{
        "assert":       lua.Func(Assert),
        "dofile":       lua.Func(DoFile),
        "error":        lua.Func(Error),
        "getmetatable": lua.Func(GetMetaTable),
        "ipairs":       lua.Func(IPairs),
        "loadfile":     lua.Func(LoadFile),
        "load":         lua.Func(Load),
        "next":         lua.Func(Next),
        "pairs":        lua.Func(Pairs),
        "pcall":        lua.Func(PCall),
        "print":        lua.Func(Print),
        "rawequal":     lua.Func(RawEqual),
        "rawlen":       lua.Func(RawLen),
        "rawget":       lua.Func(RawGet),
        "rawset":       lua.Func(RawSet),
        "select":       lua.Func(Select),
        "setmetatable": lua.Func(SetMetaTable),
        "tonumber":     lua.Func(ToNumber),
        "tostring":     lua.Func(ToString),
        "type":         lua.Func(Type),
    }

    // Open base library into globals table.
    state.Push(state.Globals())
    state.SetFuncs(funcs, 0)

    // Set global _G.
    state.PushIndex(-1)
    state.SetField(-2, "_G")

    // Set global _VERSION.
    state.Push(lua.Version)
    state.SetField(-2, "_VERSION")

    return 1
}

// assert(v, [, message])
//
// Calls error if the value of its argument v is false (i.e., nil or false);
// otherwise, returns all its arguments. In case of error, message is the error
// object; when absent, it defaults to "assertion failed!"
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-assert
func Assert(state *lua.State) int {
    unimplemented("base: assert")
    return 0
}

// dofile([filename])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-dofile
func DoFile(state *lua.State) int {
    unimplemented("base: assert")
    return 0
}

// error(message [, level])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-error
func Error(state *lua.State) int {
    unimplemented("base: error")
    return 0
}

// getmetatable(object)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-getmetatable
func GetMetaTable(state *lua.State) int {
    unimplemented("base: getmetatable")
    return 0
}

// ipairs(t)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-ipairs
func IPairs(state *lua.State) int {
    unimplemented("base: ipairs")
    return 0
}

// loadfile([filename [, mode [, env]]])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-loadfile
func LoadFile(state *lua.State) int {
    unimplemented("base: loadfile")
    return 0
}

// load(chunk [, chunkname [, mode [, env]]])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-load
func Load(state *lua.State) int {
    unimplemented("base: load")
    return 0
}

// next(table [, index])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-next
func Next(state *lua.State) int {
    unimplemented("base: next")
    return 0
}

// pairs(t)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-pairs
func Pairs(state *lua.State) int {
    unimplemented("base: pairs")
    return 0
}

// pcall(f [, arg1, ...])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-pcall
func PCall(state *lua.State) int {
    unimplemented("base: assert")
    return 0
}

// print(...)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-print
func Print(state *lua.State) int {
    for i := 1; i <= state.Top(); i++ {
        fmt.Printf("%v ", state.Value(i))
    }
    fmt.Println()
    return 0
}

// rawequal(v1, v2)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-rawequal
func RawEqual(state *lua.State) int {
    unimplemented("base: rawequal")
    return 0
}

// rawlen(v)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-rawlen
func RawLen(state *lua.State) int {
    unimplemented("base: rawlen")
    return 0
}

// rawget(table, index)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-rawget
func RawGet(state *lua.State) int {
    unimplemented("base: rawget")
    return 0
}

// rawset(table, index, value)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-rawset
func RawSet(state *lua.State) int {
    unimplemented("base: rawset")
    return 0
}

// select(index, ...)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-select
func Select(state *lua.State) int {
    unimplemented("base: select")
    return 0
}

// setmetatable(table, metatable)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-setmetatable
func SetMetaTable(state *lua.State) int {
    unimplemented("base: setmetatable")
    return 0
}

// tonumber(e [, base])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-tonumber
func ToNumber(state *lua.State) int {
    unimplemented("base: tonumber")
    return 0
}

// tostring(v)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-tostring
func ToString(state *lua.State) int {
    unimplemented("base: tostring")
    return 0
}

// type(v)
// 
// Returns the type of its only argument, coded as a string.
//
// The possible results of this function are "nil" (a string,
// not the value nil), "number", "string", "boolean", "table",
// "function", "thread", and "userdata".
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-type
func Type(state *lua.State) int {
    unimplemented("base: type")
    return 0
}