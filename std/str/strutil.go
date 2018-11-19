package str

import (
	"fmt"
	"strings"

	strutil "github.com/Azure/golua/pkg/strings"
)

func repeat(str, sep string, count int64) (string, error) {
	switch length := int64(len(str + sep)); {
	case count <= 0:
		return "", nil
	case count == 1:
		return str, nil
	case length*count/count != length:
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
	return str[beg : end+1]
}

func reverse(str string) string {
	var (
		runes = []rune(str)
		count = len(runes)
	)
	for i, j := 0, count-1; i < count/2; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func byteSlice(str string, beg, end int) (bytes []byte) {
	beg, end = bounds(len(str), beg, end)
	if beg > end {
		return
	}
	return []byte(str[beg : end+1])
}

func find(s, p string, init int) (beg, end int, caps []string, ok bool) {
	if loc := strutil.Find(s[init:], p); loc != nil {
		beg = loc[0] + init + 1
		end = loc[1] + init
		ok = true
	}
	return
}
