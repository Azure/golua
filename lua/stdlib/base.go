package stdlib

import (
    "runtime"
    "fmt"
    "os"

    "github.com/Azure/golua/lua"
)

var _ = fmt.Println
var _ = os.Exit

//
// Lua Standard Library -- basic
//

// OpenBase opens the Lua standard basic library. The basic library provides core
// functions to Lua.
// 
// See https://www.lua.org/manual/5.3/manual.html#6.1
func OpenBase(state *lua.State) int {
    var baseFuncs = map[string]lua.Func{
        "assert":         lua.Func(baseAssert),
        "dofile":         lua.Func(baseDoFile),
        "error":          lua.Func(baseError),
        "getmetatable":   lua.Func(baseGetMetaTable),
        "ipairs":         lua.Func(baseIPairs),
        "loadfile":       lua.Func(baseLoadFile),
        "load":           lua.Func(baseLoad),
        "next":           lua.Func(baseNext),
        "pairs":          lua.Func(basePairs),
        "pcall":          lua.Func(basePCall),
        "print":          lua.Func(basePrint),
        "rawequal":       lua.Func(baseRawEqual),
        "rawlen":         lua.Func(baseRawLen),
        "rawget":         lua.Func(baseRawGet),
        "rawset":         lua.Func(baseRawSet),
        "select":         lua.Func(baseSelect),
        "setmetatable":   lua.Func(baseSetMetaTable),
        "tonumber":       lua.Func(baseToNumber),
        "tostring":       lua.Func(baseToString),
        "type":           lua.Func(baseType),
        "xpcall":         lua.Func(baseXpcall),
        "collectgarbage": lua.Func(baseGC),
    }

    // Open base library into globals table.
    state.Push(state.Globals())
    state.SetFuncs(baseFuncs, 0)

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
func baseAssert(state *lua.State) int {
    if state.ToBool(1) {
        return state.Top()
    }
    state.CheckAny(1)
    state.Remove(1)
    state.Push("assert failed!")
    state.SetTop(1)
    return baseError(state)
}

// dofile([filename])
//
// Opens the named file and executes its contents as a Lua chunk.
// When called without arguments, dofile executes the contents of
// the standard input (stdin). Returns all values returned by the
// chunk. In case of errors, dofile propagates the error to its
// caller (that is, dofile does not run in protected mode).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-dofile
func baseDoFile(state *lua.State) int {
    var (
        file = state.OptString(1, "")
        src interface{} = nil
    )
    state.SetTop(1)
    if file == "" {
        file = "stdin"
        src = os.Stdin
    }
    if err := state.Load(file, src, 0); err != nil {
        panic(err)
    }
    state.Call(0, lua.MultRets)
    return state.Top() - 1
}

// error(message [, level])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-error
func baseError(state *lua.State) int {
    state.Debug(true)
    unimplemented("base: error")
    return 0
}

// getmetatable(object)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-getmetatable
func baseGetMetaTable(state *lua.State) int {
    state.CheckAny(1)
    if !state.GetMetaTableAt(1) {
        state.Push(nil)
        return 1 // no metatable
    }
    state.GetMetaField(1, "__metatable")
    return 1 // return either __metatable field (if present) or metatable.
}

// ipairs(t)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-ipairs
func baseIPairs(state *lua.State) int {
    ipairs := func(state *lua.State) int { // iterator function
        i := state.CheckInt(2) + 1
        state.Push(i)
        if state.GetI(1, i) == lua.NilType {
            return 1
        }
        return 2
    }
    state.CheckAny(1)
    state.PushClosure(ipairs, 0) // return iterator,
    state.PushIndex(1) // state,
    state.Push(0) // initial value
    return 3
}

// loadfile([filename [, mode [, env]]])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-loadfile
func baseLoadFile(state *lua.State) int {
    var (
        name string = state.OptString(1, "")
        env  int = 0   // env index (0 if no env)
    )
    var mode lua.Mode // mode "b", "b", or "bt"
    switch state.OptString(2, "bt") {
        case "b":
            mode |= lua.BinaryMode
        case "t":
            mode |= lua.TextMode
    }
    if !state.IsNone(3) {
        env = 3
    }
    if err := state.Load(name, nil, mode); err != nil {
        // error message is on top of the stack
        state.Push(err.Error())
        state.Insert(-2) // put before error message
        return 2 // return nil plus error message
    } else {
        if env != 0 { // 'env' parameter?
            state.PushIndex(env) // push environment for loaded function
            if state.SetUpValue(-2, 1) == "" { // set is as 1st upvalue
                state.Pop() // remove 'env' if not used by previous call
            }
        }
        return 1
    }
}

// load(chunk [, chunkname [, mode [, env]]])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-load
func baseLoad(state *lua.State) int {
    var (
        mode lua.Mode  // mode "b", "b", or "bt"
        name string
        env  int = 0   // env index (0 if no env)
    )
    switch state.OptString(3, "bt") {
        case "b":
            mode |= lua.BinaryMode
        case "t":
            mode |= lua.TextMode
    }
    if state.TypeAt(4) != lua.NilType {
        env = 4
    }
    chunk, ok := state.TryString(1)
    if ok && chunk != "" { // loading a string?
        name = state.OptString(2, chunk)    
    } else {
        // otherwise loading from a reader
        name = state.OptString(3, "=(load)")
        state.CheckType(1, lua.FuncType)
        // create reserved slot
        state.Debug(true)
    }
    if err := state.Load(name, chunk, mode); err != nil {
        // error message is on top of the stack
        state.Push(nil)
        state.Insert(-2) // put before error message
        return 2 // return nil plus error message
    } else {
        if env != 0 { // 'env' parameter?
            state.PushIndex(env) // push environment for loaded function
            if state.SetUpValue(-2, 1) == "" { // set is as 1st upvalue
                state.Pop() // remove 'env' if not used by previous call
            }
        }
        return 1
    }
}

// next(table [, index])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-next
func baseNext(state *lua.State) int {
    state.CheckType(1, lua.TableType)
    state.SetTop(2) // create a 2nd argument if there isn't one
    if state.Next(1) {
        return 2
    }
    state.Push(nil)
    return 1
}

// pairs(t)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-pairs
func basePairs(state *lua.State) int {
    fmt.Println(state.GetMetaField(1, "__pairs"))
    state.Debug(true)
    if state.CheckAny(1); state.GetMetaField(1, "__pairs") == lua.NilType { // no metamethod?
        state.PushClosure(baseNext, 0) // will return generator
        state.PushIndex(1)             // state,
        state.Push(nil)                // and initial value
    } else {
        state.PushIndex(1) // argument 'self' to metamethod
        state.Call(1, 3)   // get 3 values from metamethod
    }
    return 3
}

// pcall(f [, arg1, ...])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-pcall
func basePCall(state *lua.State) int {
    unimplemented("base: assert")
    return 0
}

// print(...)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-print
func basePrint(state *lua.State) int {
    for i := 1; i <= state.Top(); i++ {
        fmt.Printf("%v\t", state.Value(i))
    }
    fmt.Println()
    return 0
}

// rawequal(v1, v2)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-rawequal
func baseRawEqual(state *lua.State) int {
    unimplemented("base: rawequal")
    return 0
}

// rawlen(v)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-rawlen
func baseRawLen(state *lua.State) int {
    unimplemented("base: rawlen")
    return 0
}

// rawget(table, index)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-rawget
func baseRawGet(state *lua.State) int {
    state.CheckType(1, lua.TableType)
    state.CheckAny(2)
    state.SetTop(2)
    state.RawGet(1)
    return 1
}

// rawset(table, index, value)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-rawset
func baseRawSet(state *lua.State) int {
    unimplemented("base: rawset")
    return 0
}

// select(index, ...)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-select
func baseSelect(state *lua.State) int {
    unimplemented("base: select")
    return 0
}

// setmetatable(table, metatable)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-setmetatable
func baseSetMetaTable(state *lua.State) int {
    if t := state.TypeAt(2); t != lua.NilType && t != lua.TableType {
        state.Errorf("nil or table expected")
    }
    state.CheckType(1, lua.TableType)
    if state.GetMetaField(1, "__metatable") != lua.NilType {
        state.Errorf("cannot change a protected metatable")
    }
    state.SetTop(2)
    state.SetMetaTableAt(1)
    return 1
}

// tonumber(e [, base])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-tonumber
func baseToNumber(state *lua.State) int {
    unimplemented("base: tonumber")
    return 0
}

// tostring(v)
//
// Receives a value of any type and converts it to a string in a human-readable
// format. For complete control of how numbers are converted, use string.format.
//
// If the metatable of v has a __tostring field, then tostring calls the corresponding
// value with v as argument, and uses the result of the call as its result.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-tostring
func baseToString(state *lua.State) int {
    state.CheckAny(1)
    state.Push(state.ToString(1))
    return 1
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
func baseType(state *lua.State) int {
    state.Push(state.TypeAt(1).String())
    return 1
}

// collectgarbage ([opt [, arg]])
//
// This function is a generic interface to the garbage collector. It performs different functions
// according to its first argument, opt:
//
//    "collect": performs a full garbage-collection cycle. This is the default option.
//  
//       "stop": stops automatic execution of the garbage collector. The collector will
//               run only when explicitly invoked, until a call to restart it.
//  
//    "restart": restarts automatic execution of the garbage collector.
//  
//      "count": returns the total memory in use by Lua in Kbytes. The value has a fractional
//               part, so that it multiplied by 1024 gives the exact number of bytes in use by
//               Lua (except for overflows).
//
//       "step": performs a garbage-collection step. The step "size" is controlled by arg.
//               With a zero value, the collector will perform one basic (indivisible) step.
//               For non-zero values, the collector will perform as if that amount of memory
//               (in KBytes) had been allocated by Lua. Returns true if the step finished a
//               collection cycle.
//
//   "setpause": sets arg as the new value for the pause of the collector (see §2.5).
//               Returns the previous value for pause.
//
// "setstepmul": sets arg as the new value for the step multiplier of the collector (see §2.5).
//               Returns the previous value for step.
//
//  "isrunning": returns a boolean that tells whether the collector is running (i.e., not stopped).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-collectgarbage
func baseGC(state *lua.State) int {
    // TODO: finish him
    switch opt := state.OptString(1, "collect"); opt {
        // case "count":
        // case "stop":
        // case "restart":
        // case "setpause":
        // case "setstepmul":
        // case "isrunning":
        case "collect":
            runtime.GC()
            state.Push(0)
        case "step":
            runtime.GC()
            state.Push(true)
        default:
            state.Push(-1)
    }
    return 1
}

// See https://www.lua.org/manual/5.3/manual.html#pdf-xpcall
func baseXpcall(state *lua.State) int {
    unimplemented("base: xpcall")
    return 0
}