package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Azure/golua/lua"
	"github.com/Azure/golua/std"
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

func init() {
	flag.BoolVar(&debug, "debug", debug, "enable verbose logging")
	flag.BoolVar(&trace, "trace", trace, "enable tracing")
	flag.BoolVar(&tests, "tests", trace, "execute tests")
	flag.Parse()
}

func main() {
	if flag.NArg() < 1 {
		must(fmt.Errorf("missing arguments"))
	}

	var opts = []lua.Option{lua.WithTrace(trace), lua.WithVerbose(debug)}
	state := lua.NewState(opts...)
	defer state.Close()
	std.Open(state)

	if tests {
		state.Push(true)
		state.SetGlobal("_U")
	}

	must(state.Main(flag.Args()...))
}
