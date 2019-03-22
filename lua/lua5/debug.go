package lua5

// import (
//	"github.com/Azure/golua/lua"
// )

// func Debug(ls *lua.Thread) {
// 	stdlib := lua.Library{
// 		Open:  lua.Loader(stdlib۰debug),
// 		Name:  "debug",
// 		Funcs: []*lua.GoFunc{
// 			lua.NewGoFunc("debug", debug.debug),
// 			lua.NewGoFunc("gethook", debug.gethook),
// 			lua.NewGoFunc("getinfo", debug.getinfo),
// 			lua.NewGoFunc("getlocal", debug.getlocal),
// 			lua.NewGoFunc("getmetatable", debug.getmetatable),
// 			lua.NewGoFunc("getregistry", debug.getregistry),
// 			lua.NewGoFunc("getupvalue", debug.getupvalue),
// 			lua.NewGoFunc("getuservalue", debug.getuservalue),
// 			lua.NewGoFunc("sethook", debug.sethook),
// 			lua.NewGoFunc("setlocal", debug.setlocal),
// 			lua.NewGoFunc("setmetatable", debug.setmetatable),
// 			lua.NewGoFunc("setupvalue", debug.setupvalue),
// 			lua.NewGoFunc("setuservalue", debug.setuservalue),
// 			lua.NewGoFunc("traceback", debug.traceback),
// 			lua.NewGoFunc("upvalueid", debug.upvalueid),
// 			lua.NewGoFunc("upvaluejoin", debug.upvaluejoin),
// 		},
// 	}
// 	ls.Require(stdlib, true)
// }

// func stdlib۰debug(ls *lua.Thread) (*lua.Table, error) {
// 	// debug.debug
// 	// debug.gethook
// 	// debug.getinfo
// 	// debug.getlocal
// 	// debug.getmetatable
// 	// debug.getregistry
// 	// debug.getupvalue
// 	// debug.getuservalue
// 	// debug.sethook
// 	// debug.setlocal
// 	// debug.setmetatable
// 	// debug.setupvalue
// 	// debug.setuservalue
// 	// debug.traceback
// 	// debug.upvalueid
// 	// debug.upvaluejoin
// 	return lua.NewTable(), nil
// }
