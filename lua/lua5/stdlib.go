package lua5

import (
	"github.com/fibonacci1729/golua/lua"
)

func CoroutineLib(ls *lua.Thread) error {
	// coroutine.create
	// coroutine.isyieldable
	// coroutine.resume
	// coroutine.running
	// coroutine.status
	// coroutine.wrap
	// coroutine.yield
	lib := lua.Library{Name: "coroutine", Open: stdlib۰coroutine}
	_, err := ls.Require(lib, true)
	return err
}

func PackageLib(ls *lua.Thread) error {
	// package.config
	// package.cpath
	// package.loaded
	// package.loadlib
	// package.path
	// package.preload
	// package.searchers
	// package.searchpath
	lib := lua.Library{Name: "package", Open: stdlib۰package}
	_, err := ls.Require(lib, true)
	return err
}

func MathLib(ls *lua.Thread) error {
	// math.abs
	// math.acos
	// math.asin
	// math.atan
	// math.ceil
	// math.cos
	// math.deg
	// math.exp
	// math.floor
	// math.fmod
	// math.huge
	// math.log
	// math.max
	// math.maxinteger
	// math.min
	// math.mininteger
	// math.modf
	// math.pi
	// math.rad
	// math.random
	// math.randomseed
	// math.sin
	// math.sqrt
	// math.tan
	// math.tointeger
	// math.type
	// math.ult
	lib := lua.Library{Name: "math", Open: stdlib۰math}
	_, err := ls.Require(lib, true)
	return err
}

func BaseLib(ls *lua.Thread) error {
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
	lib := lua.Library{Name: "_G", Open: stdlib۰base}
	_, err := ls.Require(lib, true)
	return err
}

// func TableLib(ls *lua.Thread) {}
// func IOLib(ls *lua.Thread) {}
// func OSLib(ls *lua.Thread) {}
// func StringLib(ls *lua.Thread) {}
// func UTF8Lib(ls *lua.Thread) {}
// func DebugLib(ls *lua.Thread) {}

// Stdlib requires in the Lua standard libraries.
func Stdlib(ls *lua.Thread) error {
	// var stdlibs = []lua.Library{
	// 	// lua.Library{Name: "_G", Open: lua.Loader(stdlib۰base)},
	// 	// lua.Library{Name: "package", Open: lua.Loader(stdlib۰package)},
	// 	// lua.Library{Name: "coroutine", Open: lua.Loader(stdlib۰coroutine)},
	// 	// lua.Library{Name: "table", Open: lua.Loader(stdlib۰table)},
	// 	// lua.Library{Name: "io", Open: lua.Loader(stdlib۰io)},
	// 	// lua.Library{Name: "os", Open: lua.Loader(stdlib۰os)},
	// 	// lua.Library{Name: "string", Open: lua.Loader(stdlib۰string)},
	// 	// lua.Library{Name: "math", Open: lua.Loader(stdlib۰math)},
	// 	// lua.Library{Name: "utf8", Open: lua.Loader(stdlib۰utf8)},
	// 	// lua.Library{Name: "debug", Open: lua.Loader(stdlib۰debug)},
	// }
	// for _, stdlib := range stdlibs {
	// 	if err := ls.Require(stdlib, true); err != nil {
	// 		return err
	// 	}
	// }
	if err := BaseLib(ls); err != nil {
		return err
	}
	if err := PackageLib(ls); err != nil {
		return err
	}
	if err := CoroutineLib(ls); err != nil {
		return err
	}
	// table
	// io
	// os
	// string
	if err := MathLib(ls); err != nil {
		return err
	}
	return nil
}