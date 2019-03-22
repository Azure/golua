package lua5

import (
	"fmt"
	"os"
	"github.com/Azure/golua/lua"
)

var _ = fmt.Println
var _ = os.Exit

func searcher۰preload(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	var (
		libs = ls.Context().Preload()
		name = args[0].(lua.String)
	)
	if mod := libs.Get(name); mod != nil {
		return []lua.Value{mod}, nil
	}
	return nil, fmt.Errorf("no field package.preload['%s']", name)
}

func searcher۰lua(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	var (
		path = ls.Context().Config().Path
		name = args[0].(lua.String)
	)
	file, err := path.Search(string(name), ".", lua.MOD_SEP)
	if err != nil {
		return nil, err
	}
	if file != "" {
		fn, err := lua.LoadFile(ls, file)
		if err != nil {
			return nil, err
		}
		return []lua.Value{fn, name}, nil
	}
	return nil, nil
}

func searcher۰go(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	var (
		path = ls.Context().Config().GoPath
		name = args[0].(lua.String)
	)
	file, err := path.Search(string(name), ".", lua.MOD_SEP)
	if err != nil {
		return nil, err
	}
	if file != "" {
		fn, err2 := ls.Context().GoLoader(ls, file, string(name))
		if err2 != nil {
			return nil, fmt.Errorf("%v\n\tinvalid '%s' (%v)", err, file, err2)
		}
		return []lua.Value{fn, name}, nil
	}
	return nil, nil
}

func searcher۰root(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	panic("searcher۰root") // TODO
}

func stdlib۰package(ls *lua.Thread) (lua.Value, error) {
	ls.SetGlobal("require", lua.NewGoFunc("require", base۰require))
	pkg := lua.NewTableFromMap(map[string]lua.Value{
		// package.path
		// package.cpath
		// package.config
		// package.loaded
		// package.loadlib
		// package.preload
		// package.searchers
		// package.searchpath
		"config": lua.String(fmt.Sprintf(
			"%s\n%s\n%s\n%s\n%s\n",
	 			lua.DIR_SEP,
				lua.PATH_SEP,
				lua.PATH_MARK,
				lua.EXEC_DIR,
				lua.IGNORE_MARK,
		)),
		"gopath":     lua.String(ls.Context().Config().GoPath),
		"path":       lua.String(ls.Context().Config().Path),
		"loaded":     ls.Context().Loaded(),
		"preload":    ls.Context().Preload(),
		"searchers":  nil,
		"searchpath": nil,
	})
	searchers := ls.Context().Searchers()
	searchers.Set(lua.Int(1), lua.Closure(searcher۰preload, pkg))
	searchers.Set(lua.Int(2), lua.Closure(searcher۰lua, pkg))
	searchers.Set(lua.Int(3), lua.Closure(searcher۰go, pkg))
	// searchers.Set(lua.Int(4), lua.Closure(searcher۰root, pkg))
	pkg.Set(lua.String("searchers"), searchers)
	return pkg, nil
}