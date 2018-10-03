package pml

import "strings"

type (
	Replacer interface {
		// string.gsub (s, pattern, repl [, n])
		ReplaceStringMax(subject, replace string, limit int) (string, int)
		ReplaceStringAll(subject, replace string) (string, int)
		ReplaceString(subject, replace string) (string, int)
	}

	Matcher interface {
		// string.gmatch (s, pattern)
		// string.match (s, pattern [, init])
		MatchFrom(subject string, offset int) string
		MatchIter(subject string) (<-chan string)
		Match(subject string) string
	}

	Finder interface {
		// string.find (s, pattern [, init [, plain]])
		FindFrom(subject string, offset int) []int
		//FindAll(subject string, max int) [][]int
		Find(subject string) []int
	}

	Pattern interface {
		Replacer
		Matcher
		Finder
	}
)

// Must parses a pattern expression and returns, if successful,
// a Pattern object that can be used to match against text;
// otherwise panics on an error.
func Must(pattern string) Pattern{
	p, err := compile(pattern)
	if err != nil {
		panic(err)
	}
	return p
}

// New parses a pattern expression and returns, if successful,
// a Pattern object that can be used to match against text,
// else returns nil and the error.
func New(pattern string) (Pattern, error) {
	return compile(pattern)
}

// MatchFrom is just like Match except that a third numeric argument
// 'offset' specifies where to start the search; its default value
// is 0 and can be negative.
func MatchFrom(subject, pattern string, offset int) {
	Must(pattern).MatchFrom(subject, offset)
}

// Match finds and returns the first match of pattern in the string
// subject. If it finds one, then match returns the captures from the
// pattern.
// 
// Otherwise it returns nil.
//
// If pattern specifies no captures, then the whole match is returned.
func Match(subject, pattern string) (string, error) {
	return Must(pattern).Match(subject), nil
}

// FindFrom is just like Find except that a third numeric argument
// 'offset' indicates where to start the search; its default value
// is 1 and can be negative.
func FindFrom(subject, pattern string, offset int) { /* TODO */ }

// Find looks for and returns the first match of pattern in the string
// subject. If it finds a match, then it returns the indices of subject
// where this occurence starts and ends; otherwise, it returns nil.
//
// If the pattern has captures, then in a successful match the captured
// values are also returned, after the two indices.
func Find(subject, pattern string) { /* TODO */ }

// HasSpecial reports whether str has special characters.
func HasSpecial(str string) bool { return strings.IndexAny(str, special) != -1 }