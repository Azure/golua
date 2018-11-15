package strings

import (
	"strings"

	"github.com/Azure/golua/pkg/pattern"
)

type Replacer interface {
	Replace(string) string
}

// MatchAll returns the first match found in str. If str is a pattern that
// specifies captures, the captures are returned as well; otherwise ""
// if no match was made and nil if no captures captured.
func MatchAll(text, expr string, limit int) (captures [][]string) {
	return String(expr).MatchAll(text, limit)
}

// Match returns the first match found in str. If str is a pattern that
// specifies captures, the captures are returned as well; otherwise ""
// if no match was made and nil if no captures captured.
func Match(text, expr string) (captures []string) {
	return String(expr).Match(text)
}

// FindAll finds the first instance of pattern in string and returns its start
// and end position. If the pattern specified captures, they are returned
// as well as a slice of strings. If no captures were specified the entire
// match is returned as a single capture in the slice. If no matches were
// made, start and end are returned as -1, -1, and captures is nil.
func FindAll(text, expr string, limit int) (captures [][]int) {
	return String(expr).FindAll(text, limit)
}

// Find finds the first instance of pattern in string and returns its start
// and end position. If the pattern specified captures, they are returned
// as well as a slice of strings. If no captures were specified the entire
// match is returned as a single capture in the slice. If no matches were
// made, start and end are returned as -1, -1, and captures is nil.
func Find(text, expr string) (captures []int) {
	return String(expr).Find(text)
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
func Gmatch(text, expr string, iter func([]string)) {
	String(expr).Gmatch(text, iter)
}

// Gsub returns a copy of text in which all (or the upto limit if > 0) occurrences of
// the pattern have been replaced by the specified replacer. The name gsub comes from
// global substitution.
//
// Gsub returns the string with replacements and the number of replacements that occurred.
// If no matches were made, then text is return unmodified with 0 to indicate that no
// replacements were made.
func Gsub(text, expr string, replacer Replacer) (repl string, count int) {
	return GsubAll(text, expr, replacer, 1)
}

// GsubAll returns a copy of text in which all (or the upto limit if > 0) occurrences of
// the pattern have been replaced by the specified replacer. The name gsub comes from
// global substitution.
//
// Gsub returns the string with replacements and the number of replacements that occurred.
// If no matches were made, then text is return unmodified with 0 to indicate that no
// replacements were made.
func GsubAll(text, expr string, replacer Replacer, limit int) (repl string, count int) {
	return String(expr).GsubAll(text, replacer, limit)
}

// GsubStr returns a copy of text in which all (or the upto limit if > 0) occurrences
// of the pattern have been replaced by the specified replacement value. The name gsub
// comes from global substitution.
//
// The character '%' works as an escape character: any sequence in replace of the form
// %d, with 1 <= d <= 9, stands for the value of the d-th capture substring. The sequence
// %0 stands for the whole match. The sequence %% stands for a single escaped %.
func GsubStr(text, expr, replace string) (repl string, count int) {
	return GsubStrAll(text, expr, replace, 1)
}

// GsubStrAll returns a copy of text in which all (or the upto limit if > 0) occurrences
// of the pattern have been replaced by the specified replacement value. The name gsub
// comes from global substitution.
//
// The character '%' works as an escape character: any sequence in replace of the form
// %d, with 1 <= d <= 9, stands for the value of the d-th capture substring. The sequence
// %0 stands for the whole match. The sequence %% stands for a single escaped %.
func GsubStrAll(text, expr, replace string, limit int) (repl string, count int) {
	var (
		b strings.Builder
		i = 0
	)
	for _, caps := range pattern.MatchIndexAll(text, expr, limit) {
		gsub := capRE.ReplaceAllStringFunc(replace, func(k string) string {
			if i := k[1] - '0'; 0 <= i && i <= 9 {
				// TODO: check that i is valid capture index
				var (
					from, to = 2 * i, 2*i + 1
				)
				return text[caps[from]:caps[to]]
			}
			return k
		})

		b.WriteString(text[i:caps[0]])
		b.WriteString(gsub)
		i = caps[1]
		count++
	}
	b.WriteString(text[i:])
	return b.String(), count
}

// GsubFunc is just like Gsub except that replace is called to retrieve replacement values.
// The function replace is called every time a match occurs, with all captured substrings
// passed as arguments, in order.
func GsubFunc(text, expr string, replace func(string) string) (string, int) {
	return GsubFuncAll(text, expr, replacerFn(replace), 1)
}

// GsubFuncAll is just like GsubFunc except that upto limit matches are replaced.
func GsubFuncAll(text, expr string, replace func(string) string, limit int) (string, int) {
	return GsubAll(text, expr, replacerFn(replace), limit)
}

// GsubMap is just like Gsub except that values are retrieved from the vars map.
// For every match, the table is queried using the first capture as the key.
func GsubMap(text, expr string, vars map[string]string) (string, int) {
	return GsubMapAll(text, expr, vars, 1)
}

// GsubMapAll is just like GsubMax upto limit matches are replaced.
func GsubMapAll(text, expr string, vars map[string]string, limit int) (string, int) {
	replacer := replacerFn(func(capture string) string {
		v, ok := vars[capture]
		if ok {
			return v
		}
		return capture
	})
	return GsubAll(text, expr, replacer, limit)
}
