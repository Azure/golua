package str

import (
	"strings"
	"fmt"
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