package lua

import (
	"strings"
	"fmt"
	"os"
)

//
// Configuration for paths
//
const (
	// GLUA_DEFAULT
	// GLUA_V53
	// GLUA
	//
	// GLUAGO_DEFAULT
	// GLUAGO_v53
	// GLUAGO

	// Lua search paths.
	LUAPATH_DEFAULT = LUA_DIR + "?.lua;" + LUA_DIR + "?/init.lua;" + GO_DIR + "?.lua;" + GO_DIR + "?/init.lua;" + "./?.lua;" + "./?/init.lua"
	LUA_DIR 	    = ROOT_DEFAULT + "share/lua/5.3/"
	LUAPATH 	    = "GLUA_PATH"
	// Go search paths.
	GOPATH_DEFAULT = GO_DIR + "?.so;" + GO_DIR + "loadall.so;" + "./?.so"
	GO_DIR 		   = ROOT_DEFAULT + "lib/lua/5.3/"
	GOPATH 		   = "GLUA_GOPATH"
	// System root paths.
	ROOT_DEFAULT   = "/usr/local/"
	ROOT   	       = "LUA_ROOT"

	// PATH_MARK is the string that marks the substitution points in a template.
	PATH_MARK = "?"

	// PATH_SEP is the character that separates templates in a path.
	PATH_SEP = ";"

	// EXEC_DIR in a Windows path is replaced by the executable's directory.
	EXEC_DIR = "!"

	// LUA_DIRSEP is the directory separator (for submodules).
	//
	// CHANGE it if your machine does not use "/" as the directory
	// separator and is not Windows.
	//
	// Windows Lua automatically uses "\".
	DIR_SEP = "/"

	// MOD_SEP (LUA_CSUBSEP/LUA_LSUBSEP) is the character that replaces dots
	// in submodule names when searching for a Go/Lua loader.
	MOD_SEP = DIR_SEP

	// LUA_IGMARK is a mark to ignore all before it when building
	// the luaopen_ function name.
	IGNORE_MARK = "-"
)

// @path templates
//
// A typical path is a list of directories wherein to search for a given file.
// The path used by require is a list of "templates", each of them specifying
// an alternative way to transform a module name (the argument to require) into
// a file name. Each template in the path is a file name containing optional
// question marks. For each template, require substitutes the module name for
// each question mark and checks whether there is a file with the resulting name;
// if not, it goes to the next template. The templates in a path are separated by
// semicolons, a character seldom used for file names in most operating systems.
//
// For instance, consider the path: "?;?.lua;c:\windows\?;/usr/local/lua/?/?.lua"
// 
// With this path, the call require"sql" will try to open the following Lua files:
//		sql
//		sql.lua
//		c:\windows\sql
//		/usr/local/lua/sql/sql.lua
//
// The path require uses to search for Lua files is aways the current value of the
// variable "package.path". When the module package is initialized, it sets the
// variable with the value of the environment variable LUA_PATH_5_3; if this
// environment variable is undefined, Lua tries the environment variable LUA_PATH.
// If both are undefined, Lua uses a compiled-defined default path (-E prevents
// the use of these environment variables and forces the default).
//
// When using the value of an environment variable, Lua substitues the default path
// for any substring ";;". For instance, if we set LUA_PATH_5_3 to "mydir/?.lua;;",
// the final path will be the template "mydir/?.lua" followed by the default path.
//
// Path is a template (or list of templates separated by ';') used to search for
// Lua and Go packages.
type Path string

// Search searches for the given name in the given path.
//
// A path is a string containing a sequence of templates separated by semicolons.
// For each template, the function replaces each interrogation mark (if any) in
// the template with a copy of name wherein all occurences of sep (a dot, by default)
// were replaced by rep (the system's directory separator, by default), and then tries
// to ope nthe resulting file name.
//
// For instance, if the path is the string:
//
//		"./?.lua;./.lc;/usr/local/?/init.lua"
//
// The search for the name "foo.bar" will try to open the files (in order):
//
//		"./foo/bar.lua"
//		"./foo/bar.lc"
//		"/usr/local/foo/bar/init.lua"
//
// Returns the resulting name of the first file that it can open in read mode
// (after closing the file), or "" and the error if none succeeds (this error
// message lists all the file names it tried to open).
// func SearchPath(ls *Thread, name, path, sep, dirsep string) (string, error) {
func (tpl Path) Search(name, sep, dirsep string) (string, error) {
	path := string(tpl)

	if path == "" || name == "" {
		return "", nil
	}
	if sep != "" {
		// non-empty separator then replace it by 'dirsep'
		name = strings.Replace(name, sep, dirsep, -1)
	}

	path = strings.Replace(path, PATH_MARK, name, -1)
	path = strings.TrimSuffix(path, PATH_SEP)
	var b strings.Builder

	for _, file := range strings.Split(path, PATH_SEP) {
		f, err := os.OpenFile(file, os.O_RDONLY, 0666)
		if f.Close(); err != nil {
			switch {
				case os.IsPermission(err):
					// file is not readable
					fmt.Fprintf(&b, "\n\tno file '%s'", file)
					continue
				case os.IsNotExist(err):
					// file does not exist
					fmt.Fprintf(&b, "\n\tno file '%s'", file)
					continue
			}
			// uh-oh
			panic(err)
		}
		return file, nil
	}
	return "", fmt.Errorf("%s", b.String())
}

// @require(module)
//
// 1: Check in package.loaded whether "module" is already loaded.
// 	  If so, return corresponding value. Therefore, once a module
//	  is loaded, other calls requiring the same module simply return
//	  the same value, without running any code again.
//
// 2: Otherwise, searches for a Lua file with the module name (this
//	  search is guided by "package.path"). If it finds such a file,
//	  it loads it with LoadFile. The result is a function that we
//	  call a "Loader".
//
//	  Loader: A function that, when called, loads the module.
//
// 3: Otherwise, searches for a Go library with that name (this search is
//    guided by "package.gopath"). If it finds a Go library, it loads it
//	  with the low-level function "package.loadlib", looking for a function
//	  called "luaopen_modname". The loader in this case is the result of
// 	  "loadlib", which is the Go function "luaopen_modname" represented as
//	  a Lua function. 
//
// 4: No matter whether the module was found in a Lua file or a Go library,
//	  require now has a loader for it. To finally load the module, require
//	  calls the loader with two arguments: the module name and the name
//	  of the file where it got the loader (most modules just ignore these
// 	  arguments).
//
//	  If the loader returns any value, require returns this value and stores
//	  it in the "package.loaded" table, to return the same value in future
//	  calls for this same module. IF the loader returns no value, and the
//	  table entry "package.loaded[module]" is still empty, require behaves
//	  as if the module returned true. Without this correction, a subsequent
//	  call to require would run the module again.
func Require(ls *Thread, module string) (Value, error) {
	// Check in package.loaded whether "module" is already loaded.
	// If so, return corresponding value. Therefore, once a module
	// is loaded, other calls requiring the same module simply
	// return the same value, without running any code again.
	pkgs := ls.Context().Loaded()
	name := String(module)

	if mod := pkgs.Get(name); mod != nil {
		return mod, nil
	}
	loader, err := Search(ls, module)
	if err != nil {
		return nil, err
	}
	mod, err := loader(ls)
	if err != nil {
		return nil, err
	}
	if mod != nil {
		pkgs.Set(name, mod)
	}
	if pkgs.Get(name) == nil {
		pkgs.Set(name, True)
	}
	return pkgs.Get(name), nil
}

func Search(ls *Thread, module string) (Loader, error) {
	var (
		searchers = ls.Context().Searchers().Slice()
		b strings.Builder
	)
	for _, searcher := range searchers {
		rets, err := ls.CallN(searcher, []Value{String(module)}, 2)
		if err != nil {
			fmt.Fprintf(&b, "\t%v", err)
			continue
		}
		if _, ok := rets[0].(Callable); ok {
			return Loader(func(ls *Thread) (Value, error) {
				rets, err := ls.CallN(rets[0], []Value{rets[1]}, 1)
				if err != nil {
					return nil, err
				}
				return rets[0], nil
			}), nil
		}
	}
	err := fmt.Errorf("module '%s' not found:\n%s", module, b.String())
	return nil, err
}