package luac

import (
	"math/big"
	"strings"
	"strconv"
	"math"
	"fmt"
	"os"

	"github.com/fibonacci1729/golua/lua/code"
)

const debug = false

func (ls *lexical) printTrace(args ...interface{}) {
	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = len(dots)
	fmt.Printf("%5d: ", ls.line)
	i := 2 * ls.indent
	for i > n {
		fmt.Print(dots)
		i -= n
	}
	// i <= n
	fmt.Print(dots[0:i])
	fmt.Println(args...)
}

func trace(ls *lexical, msg string) *lexical {
	if debug {
		ls.printTrace(msg, "(")
		ls.indent++
	}
	return ls
}

func un(ls *lexical) {
	if debug {
		ls.indent--
		ls.printTrace(")")
	}
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func not(i int) int {
	if i == 0 {
		return 1
	}
	return 0
}

func min(x, y int) int {
	if x <= y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x >= y {
		return x
	}
	return y
}

func dumpcode(fs *function) {
	for pc, inst := range fs.instrs {
		fmt.Printf("[%d] pc@%d = %v\n", inst.line, pc, inst.code)
	}
	os.Exit(1)
}

func dumptoks(ls *lexical) {
	for ls.token.char != tEOS {
		fmt.Println(ls.token)
		ls.next()
	}
	os.Exit(1)
}

// int2fb converts an integer to a "floating point byte", represented as (eeeeexxx),
// where the real value is (1xxx) * 2^(eeeee - 1) if eeeee != 0 and (xxx) otherwise.
func int2fb(x int) int {
	if x < 8 {
		return x
	}
	e := 0
	for ; x >= 0x10; e++ {
		x = (x + 1) >> 1
	}
	return ((e + 1) << 3) | (x - 8)
}

func eval(op code.Op, x, y code.Const) (code.Const, bool) {
	switch op {
		case code.OpDivF, code.OpPow:
			if x, ok := toFloat(x); ok {
				if y, ok := toFloat(y); ok {
					return numop(op, x, y), true
				}
			}

		case code.OpBand,
			code.OpBor,
			code.OpBxor,
			code.OpShl, 
			code.OpShr,
			code.OpBnot:
			
			if x, ok := toInt(x); ok {
				if y, ok := toInt(y); ok {
					return intop(op, x, y), true
				}
			}

		default:
			if x, ok := x.(int64); ok {
				if y, ok := y.(int64); ok {
					return intop(op, x, y), true
				}
			}
			if x, ok := toFloat(x); ok {
				if y, ok := toFloat(y); ok {
					return numop(op, x, y), true
				}
			}
	}
	// TODO: Lua checks metamethods here
	return nil, false
}

func numop(op code.Op, x, y float64) float64 {
	switch op {
		case code.OpMinus:
			return -x
		case code.OpDivF:
			return x / y
		case code.OpDivI:
			return math.Floor(float64(x/y))
		case code.OpAdd:
			return x + y
		case code.OpSub:
			return x - y
		case code.OpMul:
			return x * y
		case code.OpPow:
			return math.Pow(float64(x), float64(y))
		case code.OpMod:
			f64 := math.Mod(float64(x), float64(y))
			if f64 * y < 0 {
				f64 += y
			}
			return f64
	}
	panic(op)
}

func intop(op code.Op, x, y int64) int64 {
	switch op {
		case code.OpMinus:
			return -x
		case code.OpDivI:
			return x / y
		case code.OpBand:
			return x & y
		case code.OpBnot:
			return ^x 
		case code.OpBxor:
			return x ^ y
		case code.OpBor:
			return x | y
		case code.OpAdd:
			return x + y
		case code.OpSub:
			return x - y
		case code.OpMul:
			return x * y
		case code.OpMod:
			if r := (x % y); r != 0 && (x ^ y) < 0 { // 'm/n' would be non-integer negative?
				r += y // correct result for different rounding
				return r
			} else {
				return r
			}
		case code.OpShl:
			return shiftLeft(x, y)
		case code.OpShr:
			return shiftRight(x, y)
	}
	panic(op)
}

// shift left operation
func shiftLeft(x, y int64) int64 {
	if y >= 0 {
		return x << uint64(y)
	}
	return shiftRight(x, -y)
}

// shift right operation
func shiftRight(x, y int64) int64 {
	if y >= 0 {
		return int64(uint64(x) >> uint64(y))
	}
	return shiftLeft(x, -y)
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

func str2num(s string) (code.Const, bool) {
	if i, ok := str2int(s); ok {
		return i, true
	}
	if f, ok := str2float(s); ok {
		return f, true
	}
	return nil, false
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