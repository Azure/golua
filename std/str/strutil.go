package str

import (
	"strings"
	"fmt"

	strutil "github.com/Azure/golua/pkg/strings"
)

const maxInt = int(^uint(0)>>1)

func checkMax(num int) int {
	const max = 16
	if num <= 0 || num > max {
		panic(fmt.Errorf("integral size (%d) out of limits [1,16]", num))
	}
	return num
}

// parseNum converts ASCII to integer.
func parseNum(s string, start, end int) (num int, ok bool, pos int) {
	if start >= end {
		return 0, false, end
	}
	for pos = start; pos < end && isDigit(s[pos]); pos++ {
		num = num * 10 + int(s[pos] - '0')
		ok = true
	}
	return
}

// parseNumOpt converts ASCII to integer and returns the int;
// otherwise if not a number returns the default int opt.
func parseNumOpt(s string, pos, opt int) (int, int) {
	num, isNum, end := parseNum(s, pos, len(s))
	if !isNum {
		return opt, pos
	}
	return num, end
}

func optLimit(format string, opt int) (int, int) {
	if !isDigit(format[0]) { return opt, 0 }

	var ( num, pos int = 0, 0 )
	for ; pos < len(format); pos++ {
		if !isDigit(format[pos]) {
			break
		}
		num = num * 10 + int(format[pos] - '0')
	}
	return num, pos
}

// isUpper reports when b represents a uppercase ascii char.
func isUpper(b byte) bool { return 'A' <= b && b <= 'Z' }

// isDigit reports when b represents a decimal ascii char.
func isDigit(b byte) bool { return '0' <= b && b <= '9' }

func repeat(str, sep string, count int64) (string, error) {
	switch length := int64(len(str + sep)); {
		case count <= 0:
			return "", nil
		case count == 1:
			return str, nil
		case length * count / count != length:
			return "", fmt.Errorf("resulting string too large")
	}
	rep := strings.Repeat(str+sep, int(count))
	return strings.TrimSuffix(rep, sep), nil
}

// strPos converts a relative string position: negative means back
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

func bounds(len, beg, end int) (int, int) {
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

func subStr(str string, beg, end int) (sub string) {
	beg, end = bounds(len(str), beg, end)
	if beg > end {
		return
	}
	return str[beg:end+1]
}

func reverse(str string) string {
	var (
		runes = []rune(str)
		count = len(runes)
	)
	for i,j := 0,count-1; i<count/2; i,j = i+1,j-1{
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func byteSlice(str string, beg, end int) (bytes []byte) {
	beg, end = bounds(len(str), beg, end)
	if beg > end {
		return
	}
	return []byte(str[beg:end+1])
}

func find(s, p string, init int) (beg, end int, caps []string, ok bool) {
	if loc := strutil.Find(s[init:], p); loc != nil {
		beg = loc[0] + init + 1
		end = loc[1] + init
		ok = true
	}
	return
}