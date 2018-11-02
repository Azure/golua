package base

import (
    "runtime"
    "strings"
    "strconv"
    "fmt"
    "os"

    "github.com/Azure/golua/lua"
)

var _ = fmt.Println
var _ = os.Exit

//
// Lua Standard Library -- basic
//

// Open opens the Lua standard basic library. The basic library provides core
// functions to Lua.
// 
// See https://www.lua.org/manual/5.3/manual.html#6.1
func Open(state *lua.State) int {
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
// Terminates the last protected function called and returns message
// as the error object. Function error never returns. Usually, error
// adds some information about the error position at the beginning of
// the message, if the message is a string. The level argument specifies
// how to get the error position. With level 1 (the default), the error
// position is where the error function was called. Level 2 points the
// error to where the function that called error was called; and so on.
// Passing a level 0 avoids the addition of error position information
// to the message.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-error
func baseError(state *lua.State) int {
    level := int(state.OptInt(2, 1))
    state.SetTop(1)
    if state.TypeAt(1) == lua.StringType && level > 0 {
        state.Where(level)
        state.PushIndex(1)
        state.Concat(2)
    }
    return state.Error()
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
        if t := state.GetI(1, i); t == lua.NilType || t == lua.NoneType {
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
        state.Push(nil)
        state.Push(err.Error())
        // state.Insert(-2) // put before error message
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
    if !state.IsNoneOrNil(4) {
        env = 4
    }
    chunk, ok := state.TryString(1)
    if ok && chunk != "" { // loading a string?
        name = state.OptString(2, "")    
    } else {
        // otherwise loading from a reader
        name = state.OptString(3, "=(load)")
        state.CheckType(1, lua.FuncType)
        // create reserved slot
        state.Debug(true)
    }
    if err := state.Load(name, chunk, mode); err != nil {
        // error message is on top of the stack
        state.Push(err.Error())
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
// Allows a program to traverse all fields of a table. Its first argument is a table and its second argument
// is an index in this table. next returns the next index of the table and its associated value. When called
// with nil as its second argument, next returns an initial index and its associated value. When called with
// the last index, or with nil in an empty table, next returns nil. If the second argument is absent, then it
// is interpreted as nil. In particular, you can use next(t) to check whether a table is empty.
// 
// The order in which the indices are enumerated is not specified, even for numeric indices. (To traverse a
// table in numerical order, use a numerical for.)
//
// The behavior of next is undefined if, during the traversal, you assign any value to a non-existent field
// in the table. You may however modify existing fields. In particular, you may clear existing fields.
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
    if state.CheckAny(1); state.GetMetaField(1, "__pairs") == lua.NoneType { // no metamethod?
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
// Calls function f with the given arguments in protected mode.
// This means that any error inside f is not propagated; instead,
// pcall catches the error and returns a status code. Its first
// result is the status code (a boolean), which is true if the
// call succeeds without errors. In such case, pcall also returns
// all results from the call, after this first result. In case of
// any error, pcall returns false plus the error message.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-pcall
func basePCall(state *lua.State) int {
    if err := state.PCall(state.Top()-1, -1, 0); err != nil {
        state.Push(false)
        state.Push(err.Error())
        return 2
    }
    state.Push(true)
    state.Insert(1)
    return state.Top()
}

// print(...)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-print
func basePrint(state *lua.State) int {
    var (
        n = state.Top()
        i = 1
    )
    state.GetGlobal("tostring")
    for ; i <= n; i++ {
        state.PushIndex(-1)
        state.PushIndex(i)
        state.Call(1, 1)
        str, ok := state.TryString(-1)
        if !ok {
            panic(fmt.Errorf("'tostring' must return a string to 'print'"))
        }
        if i > 1 {
            fmt.Print("\t")
        }
        fmt.Print(str)
        state.Pop()
    }
    fmt.Println()
    return 0
}

// rawequal(v1, v2)
//
// Checks whether v1 is equal to v2, without invoking the __eq metamethod.
// Returns a boolean.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-rawequal
func baseRawEqual(state *lua.State) int {
    state.CheckAny(1)
    state.CheckAny(2)
    state.Push(state.RawEqual(1, 2))
    return 1
}

// rawlen(v)
//
// Returns the length of the object v, which must be a table or a string,
// without invoking the __len metamethod. Returns an integer.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-rawlen
func baseRawLen(state *lua.State) int {
    if t := state.TypeAt(1); t != lua.StringType && t != lua.TableType {
        panic(fmt.Errorf("table or string expected"))
    }
    state.Push(state.RawLen(1))
    return 1
}

// rawget(table, index)
//
// Gets the real value of table[index], without invoking the __index metamethod.
// table must be a table; index may be any value.
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
// Sets the real value of table[index] to value, without invoking
// the __newindex metamethod. table must be a table, index any
// value different from nil and NaN, and value any Lua value.
//
// This function returns table.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-rawset
func baseRawSet(state *lua.State) int {
    state.CheckType(1, lua.TableType)
    state.CheckAny(2)
    state.CheckAny(3)
    state.SetTop(3)
    state.RawSet(1)
    return 1  
}

// select(index, ...)
//
// If index is a number, returns all arguments after argument number index;
// a negative number indexes from the end (-1 is the last argument).
//
// Otherwise, index must be the string "#", and select returns the total
// number of extra arguments it received.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-select
func baseSelect(state *lua.State) int {
   if state.TypeAt(1) == lua.StringType && state.CheckString(1) == "#" {
        state.Push(state.Top() - 1)
        return 1
    }
    var (
        sel = state.CheckInt(1)
        top = int64(state.Top())
    )
    switch {
        case sel > top:
            sel = top
        case sel < 0:
            sel = top + sel
    }
    if sel < 1 {
        panic(fmt.Errorf("bad argument to 'select' (index out of range)"))
    }
    return int(top - sel)
}

// setmetatable(table, metatable)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-setmetatable
func baseSetMetaTable(state *lua.State) int {
    if t := state.TypeAt(2); t != lua.NilType && t != lua.TableType {
        state.Errorf("nil or table expected")
    }
    state.CheckType(1, lua.TableType)
    if state.GetMetaField(1, "__metatable") != lua.NoneType {
        state.Errorf("cannot change a protected metatable")
    }
    state.SetTop(2)
    state.SetMetaTableAt(1)
    return 1
}

// tonumber(e [, base])
//
// When called with no base, tonumber tries to convert its argument to a number.
// If the argument is already a number or a string convertible to a number, then
// tonumber returns this number; otherwise, it returns nil.
//
// The conversion of strings can result in integers or floats, according to the
// lexical conventions of Lua (see §3.1). (The string may have leading and trailing
// spaces and a sign.)
//
// When called with base, then e must be a string to be interpreted as an integer
// numeral in that base. The base may be any integer between 2 and 36, inclusive.
// In bases above 10, the letter 'A' (in either upper or lower case) represents 10,
// 'B' represents 11, and so forth, with 'Z' representing 35. If the string e is not
// a valid numeral in the given base, the function returns nil.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-tonumber
func baseToNumber(state *lua.State) int {
    if state.IsNoneOrNil(2) {
        state.CheckAny(1)
        if state.TypeAt(1) == lua.NumberType {
            state.SetTop(1)
            return 1
        }
        n, ok := state.TryNumber(1)
        if ok {
            state.Push(n)
        } else {
            state.Push(nil)
        }
        return 1
    }
    str := strings.ToLower(strings.TrimSpace(state.CheckString(1)))
    base := state.CheckInt(2)
    if base < 2 || base > 36 {
        panic(fmt.Errorf("bad argument to 'tonumber' (base out of range)"))
    }
    num, err := strconv.ParseInt(str, int(base), 64)
    if err != nil {
        state.Push(nil)
        return 1
    }
    state.Push(num)
    return 1
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
    state.Push(state.ToStringMeta(1))
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
    state.CheckAny(1)
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

func unimplemented(msg string) { panic(fmt.Errorf(msg)) }