package strings

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Azure/golua/pkg/pattern"
)

var _ = fmt.Println

// String is a wrapper type that implements various string operations.
type String string

// MatchAll returns the first match found in str. If str is a pattern that
// specifies captures, the captures are returned as well; otherwise ""
// if no match was made and nil if no captures captured.
func (str String) MatchAll(text string, limit int) (captures [][]string) {
	return pattern.MatchAll(text, string(str), limit)
}

// Match returns the first match found in str. If str is a pattern that
// specifies captures, the captures are returned as well; otherwise ""
// if no match was made and nil if no captures captured.
func (str String) Match(text string) (captures []string) {
	return pattern.Match(text, string(str))
}

// FindAll returns the start and end position of the string text found in str.
//
// If str contains a pattern and a match is made with text, the captures
// are returns as a slice of strings; otherwise if no match or no captures
// then captures will be nil.
func (str String) FindAll(text string, limit int) [][]int {
	return pattern.MatchIndexAll(text, string(str), limit)
}

// Find returns the start and end position of the string text found in str.
//
// If str contains a pattern and a match is made with text, the captures
// are returns as a slice of strings; otherwise if no match or no captures
// then captures will be nil.
func (str String) Find(text string) []int {
	return pattern.MatchIndex(text, string(str))
}

// Gsub returns a copy of text in which all (or the upto limit if > 0) occurrences of
// the pattern have been replaced by the specified replacer. The name gsub comes from
// global substitution.
//
// Gsub returns the string with replacements and the number of replacements that occurred.
// If no matches were made, then text is return unmodified with 0 to indicate that no
// replacements were made.
func (str String) Gsub(text string, replacer Replacer) (string, int) {
	return str.GsubAll(text, replacer, 0)
}

// GsubAll returns a copy of text in which all (or the upto limit if > 0) occurrences of
// the pattern have been replaced by the specified replacer. The name gsub comes from
// global substitution.
//
// Gsub returns the string with replacements and the number of replacements that occurred.
// If no matches were made, then text is return unmodified with 0 to indicate that no
// replacements were made.
func (str String) GsubAll(text string, replacer Replacer, limit int) (repl string, count int) {
	// for {
	// 	if captures := str.Find(text); captures !=
	// 	if start == -1 && end == -1 {
	// 		break
	// 	}
	// 	text = text[:start] + replacer.Replace(captures[0]) + text[end:]
	// 	if count++; limit > 0 && count >= limit {
	// 		break
	// 	}
	// }
	// return text, count
	if i, caps := 0, str.FindAll(text, limit); caps != nil {
		var b strings.Builder
		for _, cap := range caps {
			if cap0 := cap[2:]; len(cap) > 0 {
				replace := replacer.Replace(text[cap0[0]:cap0[1]])
				b.WriteString(text[i:cap[0]])
				b.WriteString(replace)
				i = cap[1]
				count++
			}
		}
		b.WriteString(text[i:])
		repl = b.String()
	}
	return repl, count
}

// GsubStr returns a copy of text in which all (or the upto limit if > 0) occurrences
// of the pattern have been replaced by the specified replacement value. The name gsub
// comes from global substitution.
//
// The character '%' works as an escape character: any sequence in replace of the form
// %d, with 1 <= d <= 9, stands for the value of the d-th capture substring. The sequence
// %0 stands for the whole match. The sequence %% stands for a single escaped %.
func (str String) GsubExpr(text, replace string) (repl string, count int) {
	return GsubStrAll(text, string(str), replace, 0)
}

// Gmatch returns an iterator function that, each time it is called, returns the next captures
// from pattern over the string s. If pattern specifies no captures, then the whole match is
// produced in each call.
//
// As an example, the following loop will iterate over all the words from string s, printing
// one per line:
//
//     s = "hello world from Lua"
//     for w in string.gmatch(s, "%a+") do
//       print(w)
//     end
//
// The next example collects all pairs key=value from the given string into a table:
//
//    t = {}
//    s = "from=world, to=Lua"
//    for k, v in string.gmatch(s, "(%w+)=(%w+)") do
//       t[k] = v
//     end
//
// For this function, a caret '^' at the start of a pattern does not work as an anchor, as this
// would prevent the iteration.
func (str String) Gmatch(text string, iter func([]string)) {
	for _, captures := range str.MatchAll(text, 0) {
		iter(captures)
	}
}

// replacerFn implements the Replacer interface.
type replacerFn func(string) string

func (fn replacerFn) Replace(capture string) string { return fn(capture) }

var capRE = regexp.MustCompile("%([0-9])")
