package strs

import (
	"strings"
	"fmt"

	"github.com/Azure/golua/pkg/pml"
)

var _ = fmt.Println

func MatchIter(subject, pattern string) (func()(string, bool), error) {
	patt, err := pml.New(pattern)
	if err != nil {
		return nil, err
	}
	iter := patt.MatchIter(subject)
	return func() (next string, more bool) {
		for next := range iter {
			return next, true
		}
		return "", false
	}, nil
}

func ByteSlice(str string, beg, end int) (bytes []byte) {
	if beg, end = strPos(len(str), beg, end); beg > end {
		return
	}
	return []byte(str[beg:end+1])
}

func SubString(str string, beg, end int) (substr string) {
	if beg, end = strPos(len(str), beg, end); beg > end {
		return
	}
	return str[beg:end+1]
}

func IsPattern(str string) bool {
	return pml.HasSpecial(str)
}

func Reverse(str string) string {
	var (
		runes = []rune(str)
		count = len(runes)
	)
	for i,j := 0,count-1; i<count/2; i,j = i+1,j-1{
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func Repeat(str, sep string, count int) (rep string) {
	// switch {
	// 	case count <= 0:
	// 		return ""
	// 	case count == 1:
	// 		return str
	// }
	if count <= 0 {
		return ""
	}
	rep = strings.Repeat(str+sep, count)
	return strings.TrimSuffix(rep, sep)
}

func Match(subject, pattern string, index int) (string, error) {
	p, err := pml.New(pattern)
	if err != nil {
		return "", err
	}
	return p.Match(subject), nil
}

func Find(subject, pattern string, index int) ([]int, error) {
	p, err := pml.New(pattern)
	if err != nil {
		return nil, err
	}
	return p.FindFrom(subject, index), nil
}

// pml.New(pattern).ReplaceStringMax(subject, replace, limit)
// pml.New(pattern).ReplaceStringAll(subject, replace)
// pml.New(pattern).ReplaceString(subject, replace)

func Gsub(subject, pattern, replacer string, limit int) (string, int, error) {
	p, err := pml.New(pattern)
	if err != nil {
		return "", 0, err
	}
	s, n := p.ReplaceStringMax(subject, replacer, limit)
	return s, n, nil
}

// StrPos converts a relative string position: negative means back
// from end. The absolute position is returned.
func StrPos(len, pos int) int {
	switch {
		case pos >= 0:
			return pos
		case -pos > len:
			return 0
		default:
			return len + pos + 1
	}
}

func strPos(len, beg, end int) (int, int) {
	if beg < 0 {
		beg = len + beg + 1
	}
	if end < 0 {
		end = len + end + 1
	}
	if beg < 1 {
		beg = 1
	}
	if end > len {
		end = len
	}
	beg--
	end--
	return beg, end
}