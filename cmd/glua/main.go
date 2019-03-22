package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/fibonacci1729/golua/lua/lua5"
    "github.com/fibonacci1729/golua/lua/luac"
    "github.com/fibonacci1729/golua/lua"
)

var _ = fmt.Println

func must(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "glua: %v\n", err)
        os.Exit(1)
    }
}

func main() {
    flag.Parse()

    var config = &lua.Config{
        Stdlib: lua5.Stdlib,
    }

    fn := luac.Must(luac.Bundle(luac.Defaults, flag.Args()))
    ls := lua.Must(lua.Init(config))
    _, err := ls.Call(ls.Load(fn))
    must(err)
}