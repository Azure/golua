// +build ignore

package main

import (
	"io/ioutil"
	"fmt"
	"os"

	"github.com/Azure/golua/lua/binary"
)

func read(filename string) []byte {
	b, err := ioutil.ReadFile(filename)
	must(err)
	return b
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	c1, err := binary.Load(read("luac.out"))
	must(err)

	c2, err := binary.Load(binary.Dump(&c1.Entry))
	must(err)

	fmt.Printf("%+v\n", c1.Entry)
	fmt.Printf("%+v\n", c2.Entry)
}

// {
//		Source:=stdin 
//		SrcPos:0 
//		EndPos:0 
//		Params:0 
//		Vararg:1 
//		Stack:2 
//		Code:[8388646] 
//		Consts:[] 
//		UpValues:[{InStack:1 Index:0}] 
//		Protos:[] 
//		PcLnTab:[1] 
//		Locals:[] 
//		UpNames:[_ENV]
// }
//
// {
//		Source:stdin&�_ENV 
//		SrcPos:0 
//		EndPos:0 
//		Params:0 
//		Vararg:0 
//		Stack:0 
//		Code:[] 
//		Consts:[] 
//		UpValues:[] 
//		Protos:[] 
//		PcLnTab:[] 
//		Locals:[] 
//		UpNames:[]
// }