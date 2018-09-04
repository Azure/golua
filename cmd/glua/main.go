package main

import (
	"path/filepath"
	"io/ioutil"
	"os/exec"
	"strings"
	"flag"
	"fmt"
	"os"
	"io"

	"github.com/Azure/golua/lua/stdlib"
	"github.com/Azure/golua/lua"
)

var trace bool = false

func build(script string) string {
	dir, err := ioutil.TempDir("", "glue")
	must(err)

	tmp := filepath.Join(dir, "glua.bin")
	cmd := exec.Command("luac", "-o", tmp, "-")
	cmd.Stdin = load(script)

	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Println(string(out))
		must(err)
	}

	return tmp
}

func load(path string) io.Reader {
	b, err := ioutil.ReadFile(path)
	must(err)
	return strings.NewReader(string(b))
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	flag.BoolVar(&trace, "trace", trace, "enable tracing")
	flag.Parse()
	if flag.NArg() < 1 {
		must(fmt.Errorf("missing arguments"))
	}

	state := lua.NewState(lua.Config{Trace: trace})
	defer state.Close()
	stdlib.Import(state)

	must(state.Load(build(flag.Arg(0)), nil))
}

func demo(state *lua.State) uint32 {
	fmt.Println("DEMO")
	return 0
}