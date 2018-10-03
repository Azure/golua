// +build ignore

package main

import (
	"fmt"
	"github.com/Azure/golua/pkg/pml"
)


// # Example
//
// pml.Must("(%a+)").Match("hello world")
// Output: "hello"

// StrPos converts a relative string position: negative means back
// from end. The absolute position is returned.
func strPos(len, pos int) int {
	switch {
		case pos >= 0:
			return pos
		case -pos > len:
			return 0
		default:
			return len + pos + 1
	}
}

func main() {
	// assert(string.find("1234567890123456789", ".45", -9) == 13)

	const str = "1234567890123456789"

	fmt.Println(strPos(len(str), -9))

	for _, m := range pml.Must(".45").FindFrom(str, 10) {
		fmt.Println(m)
	}
}