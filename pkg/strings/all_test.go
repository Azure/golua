package strings

import (
	"testing"
)

func TestMatchAndCaptures(t *testing.T) {
	var tests = []struct {
		pattern string
		subject string
		matches []string
	}{
		{
			pattern: "|.*|",
			subject: "one |two| three |four| five",
			matches: []string{"|two| three |four|"},
		},
		{
			pattern: "|.+|",
			subject: "one |two| three |four| five",
			matches: []string{"|two| three |four|"},
		},
		{
			pattern: "|.-|",
			subject: "one |two| three |four| five",
			matches: []string{"|two|"},
		},
		{
			pattern: "%$(.-)%$",
			subject: "4+5 = $return 4+5$",
			// matches: []string{"return 4 + 5"},
			matches: []string{"$return 4+5$", "return 4+5"},
		},
		{
			pattern: "a*",
			subject: "aaabbb",
			matches: []string{"aaa"},
		},
		{
			pattern: "..z",
			subject: "xyxyz",
			matches: []string{"xyz"},
		},
		{
			pattern: "(a+)bb",
			subject: "aaabbb",
			// matches: []string{"aaa"},
			matches: []string{"aaabb", "aaa"},
		},
		{
			pattern: "a*aaab",
			subject: "aaaaaaaaabcd",
			matches: []string{"aaaaaaaaab"},
		},
		{
			pattern: "(%l+)(%d+)(%l+)",
			subject: "aa22bb",
			// matches: []string{"aa", "22", "bb"},
			matches: []string{"aa22bb", "aa", "22", "bb"},
		},
		{
			pattern: "(%l+)(%d+)%l?",
			subject: "aa22bb",
			// matches: []string{"aa", "22"},
			matches: []string{"aa22b", "aa", "22"},
		},
		{
			pattern: "%$%((%l+)%)",
			subject: "$(var)",
			// matches: []string{"var"},
			matches: []string{"$(var)", "var"},
		},
		{
			pattern: "a(%d+)b",
			subject: "a22b",
			// matches: []string{"22"},
			matches: []string{"a22b", "22"},
		},
		// {
		// 	pattern: "x(%d+(%l+))(zzz)",
		// 	subject: "x123abczzz",
		// },
		// {
		// 	pattern: "^abc",
		// 	subject: "123abc",
		// 	// nil
		// },
		// {
		// 	pattern: "^a-$",
		// 	subject: "aaaa",
		// 	matches: []string{"aaaa"},
		// },
		// {
		// 	pattern: "(..)-%1",
		// 	subject: "xy-yx-xy",
		// },
	}
	for _, test := range tests {
		captures := Match(test.subject, test.pattern)
		if len(captures) != len(test.matches) {
			t.Errorf("%q, expected %d matches, got %d", test.pattern, len(test.matches), len(captures))
		}

		for i := range captures {
			if captures[i] != test.matches[i] {
				t.Errorf("%q, expected %q matches, got %q", test.pattern, test.matches, captures)
				break
			}
		}
	}
}
