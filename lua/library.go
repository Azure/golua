package lua

import (
	"strings"
	"plugin"
	"fmt"
	"os"
)

var _ = fmt.Println
var _ = os.Exit

type goloader func(ls *Thread, file, name string) (Value, error)

// Try to find a load function for module 'modname' at file 'filename'.
// First, change '.' to '_' in 'modname'; then, if 'modname' has
// the form X-Y (that is, it has an "ignore mark"), build a function
// name "luaopen_X" and look for it. (For compatibility, if that
// fails, it also tries "luaopen_Y".) If there is no ignore mark,
// look for a function named "luaopen_modname".
func (fn goloader) load(ls *Thread, file, name string) (Value, error) {
	// name = "Load" + strings.Title(strings.Replace(name, ".", "_", -1))
	var _ = strings.Title
	if fn == nil { // default implementation
		p, err := plugin.Open(file)
		if err != nil {
			return nil, err
		}
		s, err := p.Lookup("Loader")
		if err != nil {
			return nil, err
		}
		if fn, ok := s.(*Loader); ok {
			return Closure(func(ls *Thread, args Tuple) ([]Value, error) {
				lib, err := (*fn)(ls)
				if err != nil {
					return nil, err
				}
				return []Value{lib}, nil
			}), nil
		}
		return nil, fmt.Errorf("plugin %q symbol %q has incorrect type", file, name)
	}
	return nil, fmt.Errorf("go library loading not supported")
}

// searcher_preload
// searcher_lua
// searcher_go
// searcher_root

// const (
// 	GoPathEnvVar = "GLUA_GOPATH"
// 	PathEnvVar = "GLUA_PATH"
// )

// type Searcher interface {
// 	Search(string) Loader
// }

// func Require(ls *Thread, name string) (Value, error) {
// 	if mod := ls.Context().Loaded().Get(String(name)); Truth(mod) {
// 		return mod, nil
// 	}
// 	for _, searcher := range ls.Context().Searchers() {
// 		if load := searcher.Search(name); loader != nil {
// 			t, err := load(ls, name)
// 			if err != nil {
// 				//
// 			}
// 		}
// 	}
// }