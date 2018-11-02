// TODO: capture index (e.g. %1)
// TODO: balanced captures
// TODO: frontier patterns
// TODO: introduce opSave
// TODO: eliminate opClose
// TODO: pattern sets
// TODO: replace
package pattern

import (
// 	"strings"
	"fmt"
)

func MatchIndexAll(text, expr string, limit int) (captures [][]int) {
	return MustCompile(expr).MatchIndexAll(text, limit)
}

func MatchIndex(text, expr string) (captures []int) {
	return MustCompile(expr).MatchIndex(text)
}

func MatchAll(text, expr string, limit int) (captures [][]string) {
	return MustCompile(expr).MatchAll(text, limit)
}

func Match(text, expr string) (captures []string) {
	return MustCompile(expr).Match(text)
}

func ReplaceAll(text, expr string, repl Replacer, limit int) (string, int) {
	return MustCompile(expr).ReplaceAll(text, repl, limit)
}

func Replace(text, expr string, repl Replacer) (string, int) {
	return MustCompile(expr).Replace(text, repl)
}

const debug = false

func trace(msg string) {
	if debug {
		fmt.Printf("# %s\n", msg)
	}
}

func traceVM(src string, sp, ip int, inst instr) {
	if debug {
		fmt.Printf("[sp=%d|ip=%d] %v (src = %q)\n", sp, ip, inst, src[sp:])
	}
}