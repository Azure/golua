package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Azure/golua/lua/stdlib"
	"github.com/Azure/golua/lua"
)

var (
	trace bool = false
	debug bool = false
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
	flag.Parse()

	if flag.NArg() < 1 {
		must(fmt.Errorf("missing arguments"))
	}

	var opts = []lua.Option{lua.WithTrace(trace), lua.WithVerbose(debug)}

	state := lua.NewState(opts...)
	//defer state.Close()
	must(state.Safely(func() {
		stdlib.Import(state)

		must(state.Exec(flag.Arg(0), nil))

		fmt.Println(state.Pop())

		if trace || debug {
			lua.Debug(state)
		}
	}))
}

// func swap(state *lua.State) int {
// 	var (
// 		v1 = state.CheckInt(1)
// 		v2 = state.CheckInt(2)
// 	)
// 	state.Push(v2)
// 	state.Push(v1)
// 	return 2
// }