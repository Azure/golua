package lua5

import (
	"fmt"
	"os"
	"github.com/Azure/golua/lua"
)

var _ = fmt.Println
var _ = os.Exit

// setmetatable(table, metatable)
//
// Sets the metatable for the given table. (To change the metatable of other
// types from Lua code, you must use the debug library (§6.10).)
//
// If metatable is nil, removes the metatable of the given table.
// If the original metatable has a __metatable field, raises an error.
//
// This function returns table.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-setmetatable
func base۰setmetatable(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	// args.Arg(0).Type(ls).SetMeta(args.Table(1))
	// return return args[:1], nil
	fmt.Printf("setmetatable(args=%v)\n", args)
	return nil, nil
}

// getmetatable(object)
//
// If object does not have a metatable, returns nil. Otherwise, if the object's
// metatable has a __metatable field, returns the associated value. Otherwise,
// returns the metatable of the given object.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-getmetatable
func base۰getmetatable(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	// fmt.Println(ls.TypeOf(args[0]))
	fmt.Printf("getmetatable(args=%v)\n", args)
    fmt.Println(ls.Caller(1).Where())
	return nil, nil
}

// tostring(v)
//
// Receives a value of any type and converts it to a string in a human-readable format.
//
// For complete control of how numbers are converted, use string.format.
//
// If the metatable of v has a __tostring field, then tostring calls the corresponding
// value with v as argument, and uses the result of the call as its result.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-tostring
func base۰tostring(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	// if fn := args[0].Type(ls).Method("__tostring"); fn != nil {}
	s, ok := lua.ToString(args[0])
	if ok {
		return []lua.Value{s}, nil
	}
	s = lua.String(fmt.Sprintf("%v", args[0]))
	// try "__name"
	return []lua.Value{s}, nil
}

// require(modname)
//
// Loads the given module. The function starts by looking into the package.loaded
// table to determine whether modname is already loaded. If it is, then require
// returns the value stored at package.loaded[modname]. Otherwise, it tries to
// find a loader for the module.
//
// To find a loader, require is guided by the package.searchers sequence. By changing
// this sequence, we can change how require looks for a module. The following explanation
// is based on the default configuration for package.searchers.
//
// First require queries package.preload[modname]. If it has a value, this value (which must
// be a function) is the loader. Otherwise require searches for a Lua loader using the path
// stored in package.path. If that also fails, it searches for a C loader using the path stored
// in package.cpath. If that also fails, it tries an all-in-one loader (see package.searchers).
//
// Once a loader is found, require calls the loader with two arguments: modname and an extra
// value dependent on how it got the loader. (If the loader came from a file, this extra value
// is the file name.) If the loader returns any non-nil value, require assigns the returned value
// to package.loaded[modname]. If the loader does not return a non-nil value and has not assigned
// any value to package.loaded[modname], then require assigns true to this entry. In any case, require
// returns the final value of package.loaded[modname].
//
// If there is any error loading or running the module, or if it cannot find any loader for the module,
// then require raises an error.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-require
func base۰require(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	// name, err := args.String(0)
	// if err != nil {
	// 	return nil, err
	// }
	if name, ok := lua.ToString(args[0]); ok {
		mod, err := lua.Require(ls, string(name))
		return []lua.Value{mod}, err
	}
	err := fmt.Errorf("bad argument #1 to 'require' (string expected, got %s", args[0].String())
	return nil, err
}

// error(message [, level])
//
// Terminates the last protected function called and returns message as
// the error object. Function error never returns.
//
// Usually, error adds some information about the error position at the
// beginning of the message, if the message is a string. The level argument
// specifies how to get the error position. With level 1 (the default), the
// error position is where the error function was called. Level 2 points the
// error to where the function that called error was called; and so on.
// Passing a level 0 avoids the addition of error position information to
// the message.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-error
func base۰error(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
   	// if level := lua.OptInt(args, 2, 1); lua.IsString(args[0])
	// if level := args.OptInt(1, lua.Int(1)); args.IsString(0) && level > 0 {
	// 	ls.Error("%s%s", ls.Where(int(level)), args.String(0))
	// 	return nil
	// }
	// panic("error!")
	// return nil
	fmt.Printf("error(%v)\n", args)
	return nil, nil
}

// print(...)
//
// Receives any number of arguments and prints their values to stdout,
// using the tostring function to convert each argument to a string.
// print is not intended for formatted output, but only as a quick
// way to show a value, for instance for debugging. For complete control
// over the output, use string.format and io.write.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-print
func base۰print(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	tostring := ls.Global("tostring")
	for i, arg := range args {
		v, err := ls.CallN(tostring, []lua.Value{arg}, 1)
		if err != nil {
			return nil, err
		}
		s, ok := lua.ToString(v[0])
		if !ok {
			return nil, fmt.Errorf("'tostring' must return a string to 'print'")
		}
		if i > 0 {
			fmt.Print("\t")
		}
		fmt.Print(s)
	}
	fmt.Println()
    return nil, nil
}

// ipairs(t)
//
// Returns three values (an iterator function, the table t, and 0) so
// that the construction
//
//     for i,v in ipairs(t) do body end
//
// will iterate over the key–value pairs (1,t[1]), (2,t[2]), ..., up to
// the first nil value.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-ipairs
func base۰ipairs(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	ipairs := func(ls *lua.Thread, args lua.Tuple) (rets []lua.Value, err error) {
		var (
			i = args[1].(lua.Int) + 1
			t = args[0]
			v lua.Value
		)
		if v, err = ls.Index(t, i); err != nil {
			return nil, err
		}
		if v != nil {
			rets = []lua.Value{i, v}
		}
		return rets, err
	}
	return []lua.Value{lua.Closure(ipairs), args[0], lua.Int(0)}, nil
}

// pairs(t)
//
// If t has a metamethod __pairs, calls it with t as argument and
// returns the first three results from the call.
//
// Otherwise, returns three values: the next function, the table t,
// and nil, so that the construction
// 
//     for k,v in pairs(t) do body end
//
// will iterate over all key–value pairs of table t.
// 
// See function next for the caveats of modifying the table during
// its traversal.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-pairs
func base۰pairs(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	// args[0].Type(ls).Method("__pairs")
	// rets, ok, err := ls.TypeOf(args[0]).CallMetaN(ls, "__pairs", args, 3)
	// if err != nil {
	// 	ls.Error(err.Error())
	// }
	// if ok {
	// 	return rets
	// }
    return []lua.Value{ls.Global("next"), args[0], nil}, nil
}

// next(table [, index])
//
// Allows a program to traverse all fields of a table. Its first argument
// is a table and its second argument is an index in this table. next returns
// the next index of the table and its associated value. When called with nil
// as its second argument, next returns an initial index and its associated value.
// When called with the last index, or with nil in an empty table, next returns nil.
// If the second argument is absent, then it is interpreted as nil. In particular,
// you can use next(t) to check whether a table is empty.
//
// The order in which the indices are enumerated is not specified, even for numeric
// indices. (To traverse a table in numerical order, use a numerical for.)
// 
// The behavior of next is undefined if, during the traversal, you assign any value to
// a non-existent field in the table. You may however modify existing fields.
// In particular, you may clear existing fields.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-next
func base۰next(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	var (
		tbl = args[0].(*lua.Table)
		key lua.Value
	)
	if len(args) > 1 {
		key = args[1]
	}
	k, v, ok := tbl.Next(key)
	if k != nil && ok {
		return []lua.Value{k, v}, nil
	}
	return nil, nil
}

func stdlib۰base(ls *lua.Thread) (lua.Value, error) {
	// _G
	// _VERSION
	// assert
	// collectgarbage
	// dofile
	// error
	// getmetatable
	// ipairs
	// load
	// loadfile
	// next
	// pairs
	// pcall
	// print
	// rawequal
	// rawget
	// rawlen
	// rawset
	// require
	// select
	// setmetatable
	// tonumber
	// tostring
	// type
	// xpcall
	ls.SetGlobal("setmetatable", lua.NewGoFunc("setmetatable", base۰setmetatable))
	ls.SetGlobal("getmetatable", lua.NewGoFunc("getmetatable", base۰getmetatable))
	ls.SetGlobal("tostring", lua.NewGoFunc("tostring", base۰tostring))
	ls.SetGlobal("error", lua.NewGoFunc("error", base۰error))
	ls.SetGlobal("print", lua.NewGoFunc("print", base۰print))
	ls.SetGlobal("ipairs", lua.NewGoFunc("ipairs", base۰ipairs))
	ls.SetGlobal("pairs", lua.NewGoFunc("pairs", base۰pairs))
	ls.SetGlobal("next", lua.NewGoFunc("next", base۰next))
	return ls.Globals(), nil
}