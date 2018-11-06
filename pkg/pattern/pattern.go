// TODO: capture index (e.g. %1)
// TODO: balanced captures
// TODO: frontier patterns
// TODO: introduce opSave
// TODO: eliminate opClose
// TODO: pattern sets
// TODO: replace
// TODO: anchor tail $
package pattern

import (
	"fmt"
)

// MatchIndexAll matches all items in text that match the pattern expr upto limit
// (or all if limit is -1). The matches are returned as a two-dimensional slice
// of integer offsets that points to the boundaries of any captures found in the
// match. 
func MatchIndexAll(text, expr string, limit int) (captures [][]int) {
	return MustCompile(expr).MatchIndexAll(text, limit)
}

// MatchIndex matches the first item in text to match the pattern expr. The matche
// is returned as a slice of integer offsets that points to the boundaries of any
// captures found in the match. 
func MatchIndex(text, expr string) (captures []int) {
	return MustCompile(expr).MatchIndex(text)
}

// MatchAll is like MatchIndexAll excepts that the capture string are returned
// instead of the positions.
func MatchAll(text, expr string, limit int) (captures [][]string) {
	return MustCompile(expr).MatchAll(text, limit)
}

// Match is like MatchIndex except that the capture are returned instead of the
// positions.
func Match(text, expr string) (captures []string) {
	return MustCompile(expr).Match(text)
}

// ReplaceAll replaces all (or upto limit) matches of text in expr using the
// provided Replacer repl. The new string and number of replacements made is
// returned. If no match was made, "", 0 is returned.
func ReplaceAll(text, expr string, repl Replacer, limit int) (string, int) {
	return MustCompile(expr).ReplaceAll(text, repl, limit)
}

// Replace replaces the first match of text in expr using the provided Replacer
// repl. On success, the replaced string and number of replacements is returned;
// Otherwise "", 0 is returned.
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