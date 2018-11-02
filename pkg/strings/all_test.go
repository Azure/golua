package strings

import (
	"testing"
	"fmt"
)

func TestMatchAndCaptures(t *testing.T) {
	var tests = []struct{
		pattern string
		subject string
		matches bool
		invalid bool
	}{
		{
			pattern: "|.*|",
			subject: "one |two| three |four| five",
			// [|two| three |four|]
		},
		{
			pattern: "|.+|",
			subject: "one |two| three |four| five",
			// [|two| three |four|]
		},
		{
			pattern: "|.-|",
			subject: "one |two| three |four| five",
			// [|two|]
		},
		{
			pattern: "%$(.-)%$",
			subject: "4+5 = $return 4+5$",
			// [return 4 + 5]
		},
		{
			pattern: "a*",
			subject: "aaabbb",
			// [aaa]
		},
		{
			pattern: "..z",
			subject: "xyxyz",
			// [xyz]
		},
		{
			pattern: "(a+)bb",
			subject: "aaabbb",
			// [aaa]
		},
		{
			pattern: "a*aaab",
			subject: "aaaaaaaaabcd",
			// [aaaaaaaaab]
		},
		{
			pattern: "(%l+)(%d+)(%l+)",
			subject: "aa22bb",
			// [aa 22 bb]
		},
		{
			pattern: "(%l+)(%d+)%l?",
			subject: "aa22bb",
			// [aa 22]
		},
		{
			pattern: "%$%((%l+)%)",
			subject: "$(var)",
			// [var]
		},
		{
			pattern: "a(%d+)b",
			subject: "a22b",
			// [22]
		},
		{
			pattern: "a(%d+)b",
			subject: "a22b",
			// [22]
		},
		// {
		// 	pattern: "x(%d+(%l+))(zzz)",
		// 	subject: "x123abczzz",
		// 	// []
		// },
		// {
		// 	pattern: "^abc",
		// 	subject: "123abc",
		// 	// nil
		// },
		// {
		// 	pattern: "^a-$",
		// 	subject: "aaaa",
		// 	// [aaaa]
		// },
		// {
		// 	pattern: "(..)-%1",
		// 	subject: "xy-yx-xy",
		// 	// []
		// },
	}
	for _, test := range tests {
		captures, matched := Match(test.subject, test.pattern)
		fmt.Printf("%v (match = %t)\n", captures, matched)
	}
}