package lua

import (
	"math/big"
	"strings"
	"strconv"
	"math"
)

// Convert and returns float f as an int64, true.
//
// Otherwise returns 0, false on failure.
func float2int(f float64) (int64, bool) {
	if isNaN(f) {
		return 0, false
	}
	if i := int64(f); float64(i) == f {
		return i, true
	}
	return 0, false
}

// Convert string 's' to a Lua number (put in 'result'). Return 0, false
// on fail or the address of the ending '\0' on success. 'pmode' points to
// (and 'mode' contains) special things in the string:
// 	- 'x'/'X' means an hexadecimal numeral
// 	- 'n'/'N' means 'inf' or 'nan' (which should be rejected)
// 	- '.' just optimizes the search for the common case (nothing special)
//
// This function accepts both the current locale or a dot as the radix mark.
// If the convertion fails, it may mean number has a dot but locale accepts
// something else. In that case, the code copies 's' to a buffer (because 's'
// is read-only), changes the dot to the current locale radix mark, and tries
// to convert again.
func str2float(s string) (float64, bool) {
	if s = strings.TrimSpace(s); strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		num, _, err := new(big.Float).Parse(strings.ToLower(s), 0)
		if err != nil {
			return 0, false
		}
		f64, acc := num.Float64()
		return f64, (acc == big.Exact)
	}
	f64, err := strconv.ParseFloat(s, 10)
	return f64, (err == nil)
}

func str2int(s string) (int64, bool) {
	var (
		sign int64 = 1
		acc uint64
		pos int
	)
	if s = strings.TrimSpace(s); s[pos] == '-' || s[pos] == '+' {
		if pos++; s[pos-1] == '-' {
			sign = -1
		}
	}
	if s[pos] == '0' && pos+1 < len(s) && (s[pos+1] == 'x' || s[pos+1] == 'X') {
		for pos += 2; pos < len(s) && isHexDigit(rune(s[pos])); pos++ {
			acc = acc * 16 + uint64(hex2int(s[pos]))
		}
		return sign*int64(acc), (pos == len(s))
	}
	const (
		maxBy10 = uint64(maxInt/10)
		maxLast = int64(maxInt%10)
	)
	for pos < len(s) && isDigit(rune(s[pos])) {
		dig := uint64(s[pos] - '0')
		if acc >= maxBy10 && (acc > maxBy10 || dig > uint64(maxLast + sign)) {
			return 0, false
		}
		acc = acc * 10 + dig
		pos++
	}
	return sign*int64(acc), (pos == len(s))
}

func hex2int(r byte) int {
	switch {
	 	case '0' <= r && r <= '9':
			r = r - '0'
		case 'a' <= r && r <= 'f':
			r = r - 'a' + 10
		case 'A' <= r && r <= 'F':
			r = r - 'A' + 10
	}
	return int(r)
}

func isHexDigit(r rune) bool { return isDigit(r) || 'a' <= r && r <= 'f' || 'A' <= r && r <= 'F' }
func isDigit(r rune) bool { return '0' <= r && r <= '9' }
func isNaN(f64 float64) bool { return math.IsNaN(f64) }
const maxInt = int(^uint(0)>>1)