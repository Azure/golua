package pml

import (
	"strings"
	"regexp"
	"fmt"
)

// Special characters escaped by '%'.
//const special = `^$()%.[]*+-?`
const maxMatch = 200
const special  = "^$*+?.([%-"
const escape   = "%"

var classes = map[byte]string{
	'a': "[[:alpha:]]",  'A': "[[:^alpha:]]",  // class of all letters
	'c': "[[:cntrl:]]",  'C': "[[:^cntrl:]]",  // class of all control characters
	'd': "[[:digit:]]",  'D': "[[:^digit:]]",  // class of all digits
	'g': "[[:graph:]]",  'G': "[[:^graph:]]",  // class of all printable characters except space
	'l': "[[:lower:]]",  'L': "[[:^lower:]]",  // class of all lowercase characters
	'p': "[[:punct:]]",  'P': "[[:^punct:]]",  // class of all punctuation characters
	's': "[[:space:]]",  'S': "[[:^space:]]",  // class of all space characters
	'u': "[[:upper:]]",  'U': "[[:^upper:]]",  // class of all uppercase characters
	'w': "[[:word:]]",   'W': "[[:^word:]]",   // class of all alphanumeric characters
	'x': "[[:xdigit:]]", 'X': "[[:^xdigit:]]", // class of all hexadecimal digits
}

type pattern struct {
	re2 *regexp.Regexp
	src string
}

// func gsub(subject, pattern string, replacer string, limit int) (string, int) {
// 	var (
// 		regexpr = regexp.MustCompile(pattern)
// 		matches = regexpr.FindAllStringIndex(subject, limit)
// 		replace = strings.Replace(replacer, "%", "$", -1)
// 		total   = len(matches)
// 	)
// 	last := matches[total-1][1]
// 	head := subject[:last]
// 	tail := subject[last:]
// 	return regexpr.ReplaceAllString(head, replace) + tail, total
// }

func (patt *pattern) ReplaceStringMax(subject, replace string, limit int) (string, int) {
	return patt.replace(subject, replace, limit)
}

func (patt *pattern) ReplaceStringAll(subject, replace string) (string, int) {
	return patt.replace(subject, replace, len(subject))
}

func (patt *pattern) ReplaceString(subject, replace string) (string, int) {
	return patt.replace(subject, replace, 1)
}

func (patt *pattern) MatchFrom(str string, ofs int) string {
	if captures := patt.matches(str, ofs); captures != nil {
		start, end := captures[0], captures[1]
		return str[start:end]
	}
	return ""
}

func (patt *pattern) MatchIter(str string) (<-chan string) {
	iter := make(chan string, 1)
	go func() {
		for _, match := range patt.findMax(str, len(str)) {
			iter <- str[match[0]:match[1]]
		}
		close(iter)
	}()
	return iter
}

func (patt *pattern) Match(str string) string {
	return patt.MatchFrom(str, 1)
}

func (patt *pattern) FindFrom(str string, ofs int) []int {
	if captures := patt.matches(str, ofs); captures != nil {
		return captures[:2]
	}
	return nil
}

func (patt *pattern) Find(str string) []int {
	return patt.FindFrom(str, 1)
}

func compile(src string) (*pattern, error) {
	var (
		out strings.Builder
		str = src
	)
	brackets := false
	for len(str) > 0 {
		var ( b0, b1 byte )

		switch b0, str = str[0], str[1:]; b0 {
			case '{', '}', '\\': // escapes '{}\' chars
				out.WriteByte('\\')
				out.WriteByte(b0)
			case '[':
				brackets = true
				out.WriteByte(b0)
			case ']':
				brackets = false
				out.WriteByte(b0)
			case '-':
				if brackets {
					out.WriteByte(b0)
					continue
				}
				out.WriteString("*?")
			case '^':
				if !brackets && out.Len() > 0 {
					out.WriteString("\\")
				}
				out.WriteByte(b0)
			case '$':
				if len(str) > 0 {
					out.WriteString("\\")
				}
				out.WriteByte(b0)
			case '%':
				if len(str) == 0 {
					return nil, fmt.Errorf("malformed pattern (ends with '%')")
				}
				switch b1, str = str[0], str[1:]; b1 {
					case '(', ')', '[', ']', '.', '?', '*', '+', '-', '^', '$':
						out.WriteByte('\\')
						out.WriteByte(b1)
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						return nil, fmt.Errorf("backreference is unimplemented")
					case 'b':
						return nil, fmt.Errorf("'%b' is unimplemented")
					case 'f':
						return nil, fmt.Errorf("'%f' is unimplemented")
					case '%':
						out.WriteByte('%')
					default:
						if cc, ok := classes[b1]; ok {
							if brackets {
								out.WriteString(cc[1:len(cc)-1])
							} else {
								out.WriteString(cc)
							}
						} else {
							out.WriteByte(b1)
						}
				}
			default:
				out.WriteByte(b0)
		}
	}
	//fmt.Printf("re2: %s\n", out.String())
	re2, err := regexp.Compile(out.String())
	if err != nil {
		return nil, err
	}
	return &pattern{re2: re2, src: src}, nil
}

func (patt *pattern) replace(subject, replace string, limit int) (string, int) {
	if matches := patt.findMax(subject, limit); matches != nil {
		var (
			repl = strings.Replace(replace, "%", "$", -1)
			last = matches[len(matches)-1][1]
			head = subject[:last]
			tail = subject[last:]
		)
		return patt.re2.ReplaceAllString(head, repl) + tail, len(matches)	
	}
	return subject, 0
}

func (patt *pattern) matches(str string, ofs int) []int {
	return patt.re2.FindStringSubmatchIndex(str[ofs:])
}

func (patt *pattern) findMax(str string, max int) [][]int {
	return patt.re2.FindAllStringIndex(str, max)
}