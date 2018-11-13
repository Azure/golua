package main

import (
	"flag"
	"fmt"
	"os"

	stdlib "github.com/Azure/golua/std"
	"github.com/Azure/golua/lua"
)

var (
	trace bool = false
	debug bool = false
	tests bool = false
)

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	flag.BoolVar(&debug, "debug", debug, "display Lua IR")
	flag.BoolVar(&trace, "trace", trace, "enable tracing")
	flag.BoolVar(&tests, "tests", trace, "execute tests")
	flag.Parse()

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		if err, ok := r.(error); ok {
	// 			must(err)
	// 		}
	// 	}
	// }()

	var opts = []lua.Option{lua.WithTrace(trace), lua.WithVerbose(debug)}

	state := lua.NewState(opts...)
	//defer state.Close()

	if flag.NArg() < 1 {
		must(fmt.Errorf("missing arguments"))
	}

	if tests {
		state.Push(true)
		state.SetGlobal("_U")	
	}
	stdlib.Open(state)

	mode := lua.BinaryMode|lua.TextMode
	must(state.Exec(flag.Arg(0), nil, mode))

	if result := state.Pop(); result != lua.None {
		fmt.Println(result)
	}
}

// func Main(state *lua.State) {
// 	//argc := state.ToInt(1)
// 	//argv := state.ToUserData(2)
// 	//args := collectArgs(argv)
// }