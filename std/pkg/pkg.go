package pkg

import (
	"fmt"
	"os"
	"plugin"
	"strings"

	"github.com/Azure/golua/lua"
)

//
// Lua Standard Library -- load/package
//

func Open(state *lua.State) int {
	// Create 'package' table.
	var packageFuncs = map[string]lua.Func{
		"searchpath": lua.Func(pkgSearchPath),
		"loadlib":    lua.Func(pkgLoadLibrary),
	}
	state.NewTableSize(0, 8)
	state.SetFuncs(packageFuncs, 0)

	// Create 'searchers' table.
	createSearchersTable(state)

	// Set 'path' field.
	state.Push(lua.EnvPath)
	state.SetField(-2, "path")

	// Set 'cpath' field.
	state.Push(lua.EnvHome)
	state.SetField(-2, "gopath")

	// Set 'config' field.
	state.Push(lua.Config)
	state.SetField(-2, "config")

	// Set 'loaded' field.
	state.GetSubTable(lua.RegistryIndex, lua.LoadedKey)
	state.SetField(-2, "loaded")

	// Set 'preload' field.
	state.GetSubTable(lua.RegistryIndex, lua.PreloadKey)
	state.SetField(-2, "preload")

	// Set global 'require' function with 'package' table as an upvalue.
	var loadFuncs = map[string]lua.Func{
		"require": lua.Func(require),
	}
	state.PushGlobals()
	state.PushIndex(-2)
	state.SetFuncs(loadFuncs, 1)

	// Pop 'globals' table.
	state.Pop()

	// Return 'package' table
	return 1
}

// require(modname)
//
// Loads the given module. The function starts by looking into the package.loaded
// table to determine whether module is already loaded. If it is, then require
// returns the value stored at package.loaded[module]. Otherwise, it tries to find
// a loader for the module.
//
// To find a loader, require is guided by the package.loaders array. By changing
// this array, we can change how require looks for a module.
//
// First require queries package.preload[module]. If it has avalue, this value
// (which should be a function) is the loader. Otherwise require searches for a
// Lua loader using the path stored in package.path.
//
// Once a loader is found, require calls the loader with a single argument, module.
//
// If the loader returns any value, require assigns the returned value to package.loaded[module].
//
// If the loader returns no value and has not assigned any value to package.loaded[module],
// then require assigns true to this entry. In any case, require returns the final value of
// package.loaded[module].
//
// If there is any error loading or running the module, or if it cannot find any loader for
// the module, then require signals an error.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-require
func require(state *lua.State) int {
	modname := state.CheckString(1) // name of module to require

	// Push the LOADED table and check if module was already loaded.
	state.GetField(lua.RegistryIndex, lua.LoadedKey)
	state.GetField(-1, modname)
	if state.ToBool(-1) {
		return 1 // Module is already loaded
	}

	// Otherwise we must load the module.

	// First remove the result of GetField.
	state.Pop()

	// Find a loader for the module.
	searchLoader(state, modname)

	// Pass name as argument to module loader.
	state.Push(modname)

	// Name is 1st argument (before search data).
	state.Insert(-2)

	// Run the loader to load the module.
	state.Call(2, 1)

	// If non-nil return, then LOADED[modname] = returned value
	// if !lua.IsNone(state.Value(-1)) {
	// 	state.SetField(2, modname)
	// }

	// If module set no value, use true as result (LOADED[modname] = true).
	// if state.GetField(2, modname) == lua.NilType {
	// 	state.Push(true) // Value stored
	// 	state.Push(true) // Value returned
	// 	state.SetField(2, modname)
	// }

	return 1
}

// package.searchpath(name, path [, sep [, rep]])
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-package.searchpath
func pkgSearchPath(state *lua.State) int {
	unimplemented("package: searchpath")
	return 0
}

// package.loadlib(libname, funcname)
//
// Dynamically links the host program with the C library libname.
//
// If funcname is "*", then it only links with the library, making the symbols exported
// by the library available to other dynamically linked libraries. Otherwise, it looks
// for a function funcname inside the library and returns this function as a C function.
// So, funcname must follow the lua_CFunction prototype (see lua_CFunction).

// This is a low-level function. It completely bypasses the package and module system.
// Unlike require, it does not perform any path searching and does not automatically
// adds extensions. libname must be the complete file name of the C library, including
// if necessary a path and an extension. funcname must be the exact name exported by the
// C library (which may depend on the C compiler and linker used).
//
// This function is not supported by Standard C. As such, it is only available on
// some platforms (Windows, Linux, Mac OS X, Solaris, BSD, plus other Unix systems
// that support the dlfcn standard).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-package.loadlib
func pkgLoadLibrary(state *lua.State) int {
	var (
		path = state.CheckString(1)
		init = state.CheckString(2)
	)
	return lookForFunc(state, path, init)
}

func createSearchersTable(state *lua.State) {
	var searchers = []lua.Func{
		// preload searcher
		lua.Func(searchPreload),
		// lua searcher
		lua.Func(searchLua),
		// go searcher
		lua.Func(searchGo),
		// all-in-one loader (root)
		//lua.Func(searchRoot),
	}
	// Create 'searchers' table.
	state.NewTableSize(len(searchers), 0)

	// Fill it with the predefined searchers.
	for i := 0; i < len(searchers); i++ {
		state.PushIndex(-2)
		state.PushClosure(searchers[i], 1)
		state.RawSetIndex(-2, i+1)
	}

	// Put it in field 'searchers'.
	state.SetField(-2, "searchers")
}

func lookForFunc(state *lua.State, path, init string) int {
	p, err := plugin.Open(path)
	if err != nil {
		state.Push(err.Error())
		return 1
	}
	s, err := p.Lookup(init)
	if err != nil {
		state.Push(err.Error())
		return 1
	}
	l, ok := s.(func(*lua.State) int)
	if !ok {
		err := fmt.Errorf("plugin %q symbol %q has incorrect type", path, init)
		state.Push(err.Error())
		return 1
	}
	state.Push(l)
	return 1
}

func searchLoader(state *lua.State, modname string) {
	// Push 'package.searchers' onto the stack.
	if state.GetField(lua.UpValueIndex(1), "searchers") != lua.TableType {
		state.Errorf("'package.searchers' must be a table")
	}

	// Builds the error if module is not found.
	var errs strings.Builder

	// Iterate over available searchers to find a loader.
	for arr, at := state.Top(), 1; state.RawGetIndex(arr, at) != lua.NilType; at++ {
		// Push modname argument and call searcher.
		state.Push(modname)
		state.Call(1, 2)

		switch state.TypeAt(-2) {
		case lua.StringType:
			// The loader returned an error message
			// pop it from the stack and throw it.
			state.Pop() // remove extra argument.
			fmt.Fprintf(&errs, "%v", state.Pop())
		case lua.FuncType:
			// A loader was found, simply return.
			return
		default: // Remove both returns.
			state.PopN(2)
		}
	}

	// No loader found, pop last result and throw error.
	state.Pop() // nil from final query.

	state.Errorf("module '%s' not found:%s", modname, errs.String())
}

func searchPreload(state *lua.State) int {
	modname := state.CheckString(1)
	state.GetField(lua.RegistryIndex, lua.PreloadKey)
	if state.GetField(-1, modname) == lua.NilType {
		state.Push(fmt.Sprintf("\n\tno field package.preload['%s']", modname))
	}
	return 1
}

func searchLua(state *lua.State) int {
	var (
		modname  = state.CheckString(1)
		filename string
	)
	if filename = findFile(state, modname, "path"); filename == "" {
		// Module not found in this path.
		return 1
	}
	if err := state.LoadChunk(filename, nil, 0); err != nil {
		// Module didn't load successfully.
		state.Push(fmt.Sprintf("error loading module '%s' from file '%s':\n\t%v",
			modname,
			filename,
			err,
		))
		return 1
	}
	// Module loaded successfully. Push the script path
	// as 2nd argument to module invocator. Return 2
	// because after the above Load, the module Loader
	// function closure is on top of the stack.
	state.Push(filename)
	return 2
}

func searchGo(state *lua.State) int {
	var (
		name = strings.Replace(state.CheckString(1), ".", "_", -1)
		open = "Open"
		file string
	)
	if file = findFile(state, name, "gopath"); file == "" {
		// Module not found on this path.
		return 1
	}
	if dash := strings.Index(name, "-"); dash != -1 {
		open += strings.Title(name[dash+1:])
	}
	return lookForFunc(state, file, open)
}

// func searchRoot(state *lua.State) int {
// 	unimplemented("searcher: all-in-one")
// 	return 0
// }

func findFile(state *lua.State, name, pathkey string) string {
	state.GetField(lua.UpValueIndex(1), pathkey)
	path, ok := state.TryString(-1)
	if !ok {
		state.Errorf("'package.%s' must be a string", pathkey)
	}
	return searchPath(state, name, path, ".")
}

func searchPath(state *lua.State, name, path, sep string) string {
	if sep != "" { // non-empty separator?
		name = strings.Replace(name, sep, string(os.PathSeparator), -1)
	}
	var errMsg string
	for _, file := range strings.Split(path, ";") {
		file = strings.Replace(file, "?", name, -1)
		if _, err := os.Stat(file); err == nil {
			return file
		}
		errMsg = fmt.Sprintf("%s\n\tno file '%s'", errMsg, file)
	}
	state.Push(errMsg)
	return ""
}

//	stdin:1: module 'mymodule' not found:
//		no field package.preload['mymodule']
//		no file '/usr/local/share/lua/5.3/mymodule.lua'
//		no file '/usr/local/share/lua/5.3/mymodule/init.lua'
//		no file '/usr/local/lib/lua/5.3/mymodule.lua'
//		no file '/usr/local/lib/lua/5.3/mymodule/init.lua'
//		no file './mymodule.lua'
//		no file './mymodule/init.lua'
//		no file '/usr/local/lib/lua/5.3/mymodule.so'
//		no file '/usr/local/lib/lua/5.3/loadall.so'
//		no file './mymodule.so'
//	stack traceback:
//		[C]: in function 'require'
//		stdin:1: in main chunk
//		[C]: in ?

func unimplemented(msg string) { panic(fmt.Errorf(msg)) }
