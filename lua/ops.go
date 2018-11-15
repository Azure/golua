package lua

import (
	// "strconv"
	"fmt"
	"math"
	"strings"

	"github.com/Azure/golua/lua/syntax"
)

var _ = fmt.Println

// Op represents a Lua arithmetic, bitwise, or relational operator.
type Op int

const (
	// Arithmetic Operators
	OpAdd   Op = 1 + iota // '+'
	OpSub                 // '-'
	OpMul                 // '*'
	OpMod                 // '%'
	OpPow                 // '^'
	OpDiv                 // '/'
	OpQuo                 // '//'
	OpMinus               // '-' (unary)

	// Bitwise Operators
	OpOr  // '|'
	OpAnd // '&'
	OpXor // '~'
	OpRsh // '>>'
	OpLsh // '<<'
	OpNot // '~' (unary)

	// Relational Operators
	OpLt // '<'
	OpLe // '<='
	OpEq // '=='
	OpGt // '>'
	OpGe // '>='
	OpNe // '~='

	// Miscellaneous Operators
	OpConcat // '..'
	OpLength // '#'
)

// compare applies the relational operator to the values x and y.
//
// These operators always result in true or false.
//
// Equality (==) first compares the type of its operands. If the types are different, then the result
// is false. Otherwise, the values of the operands are compared. Strings are compared in the obvious
// way. Numbers are equal if they denote the same mathematical value.
//
// Tables, userdata, and threads are compared by reference: two objects are considered equal only if
// they are the same object. Every time you create a new object (a table, userdata, or thread), this
// new object is different from any previously existing object. A closure is always equal to itself.
// Closures with any detectable difference (different behavior, different definition) are always different.
// Closures created at different times but with no detectable differences may be classified as equal or
// not (depending on internal caching details).
//
// You can change the way that Lua compares tables and userdata by using the "eq" metamethod (see §2.4).
//
// Equality comparisons do not convert strings to numbers or vice versa. Thus, "0"==0 evaluates to
// false, and t[0] and t["0"] denote different entries in a table.
//
// The operator ~= is exactly the negation of equality (==).
//
// The order operators work as follows. If both arguments are numbers, then they are compared according
// to their mathematical values (regardless of their subtypes). Otherwise, if both arguments are strings,
// then their values are compared according to the current locale. Otherwise, Lua tries to call the "lt"
// or the "le" metamethod (see §2.4). A comparison a > b is translated to b < a and a >= b is translated
// to b <= a.
//
// Following the IEEE 754 standard, NaN is considered neither smaller than, nor equal to, nor greater
// than any value (including itself).
//
// See https://www.lua.org/manual/5.3/manual.html#3.4.5
//
// Returns true if "x op y"; otherwise false.
func (state *State) compare(op Op, x, y Value, raw bool) bool {
	// metamethod event to call if type op type pairs are exhausted.
	var event metaEvent

	switch op {
	case OpEq: // '=='
		if x.Type() != y.Type() {
			if IsNone(x) {
				return IsNone(y)
			}
			return false
		}
		switch x := x.(type) {
		case *Closure:
			if y, ok := y.(*Closure); ok {
				if x.isLua() && y.isLua() {
					return x.binary == y.binary
				}
				return true
			}
		case *Object:
			if y, ok := y.(*Object); ok {
				return x.data == y.data
			}
		case *table:
			if y, ok := y.(*table); ok && (x == y) {
				return true
			}
		case String:
			// x (string) == y (string)
			if y, ok := y.(String); ok {
				return x == y
			}
		case Bool:
			// x (boolean) == y (boolean)
			if y, ok := y.(Bool); ok {
				return x == y
			}
		case Number:
			if x, ok := x.(Int); ok {
				if y, ok := y.(Int); ok {
					return x == y
				}
			}
			if x, ok := x.(Float); ok {
				if y, ok := y.(Float); ok {
					return x == y
				}
			}
			x, ok1 := toInteger(x)
			y, ok2 := toInteger(y)
			return ok1 && ok2 && (x == y)

		// case Float:
		// // x (float) == y (float)
		// if y, ok := y.(Float); ok {
		//     return x == y
		// }
		// // x (float) == y (int)
		// if y, ok := y.(Int); ok {
		//     i, ok := toInteger(x)
		//     if ok {
		//         return i == y
		//     }
		// }
		// case Int:
		//     // x (int) == y (float)
		//     if y, ok := y.(Float); ok {
		//         i, ok := toInteger(y)
		//         if ok {
		//             return x == i
		//         }
		//     }
		//     // x (int) == y (int)
		//     if y, ok := y.(Int); ok {
		//         return x == y
		//     }

		case Nil:
			// x (nil) == y (none)
			// x (nil) == y (nil)
			if IsNone(y) {
				return true
			}
		}

		// try __eq
		event = metaEq

	case OpLe: // '<='
		switch x := x.(type) {
		case String:
			// x (string) <= y (string)
			if y, ok := y.(String); ok {
				return x <= y
			}
		case Float:
			// x (float) <= y (float)
			if y, ok := y.(Float); ok {
				return x <= y
			}
			// x (float) <= y (int)
			if y, ok := y.(Int); ok {
				if x != x { // x is NaN
					return false
				}
				var cmp int
				switch {
				case !math.IsInf(float64(x), 0): // x is finite
					cmp = x.rational().Cmp(y.rational())
				case x > 0: // x is +inf
					cmp = -1
				default: // x is -inf
					cmp = +1
				}
				return threeway(op, cmp)
			}
		case Int:
			// x (int) <= y (float)
			if y, ok := y.(Float); ok {
				if y != y { // y is NaN
					return false
				}
				var cmp int
				switch {
				case !math.IsInf(float64(y), 0): // y is finite
					cmp = x.rational().Cmp(y.rational())
				case y > 0: // y is +inf
					cmp = -1
				default: // y is -inf
					cmp = +1
				}
				return threeway(op, cmp)
			}
			// x (int) <= y (int)
			if y, ok := y.(Int); ok {
				return x <= y
			}
		}
		// try __le
		event = metaLe

	case OpLt: // '<'
		switch x := x.(type) {
		case String:
			// x (string) < y (string)
			if y, ok := y.(String); ok {
				return x < y
			}
		case Float:
			// x (float) < y (float)
			if y, ok := y.(Float); ok {
				return x < y
			}
			// x (float) < y (int)
			if y, ok := y.(Int); ok {
				if x != x { // x is NaN
					return false
				}
				var cmp int
				switch {
				case !math.IsInf(float64(x), 0): // x is finite
					cmp = x.rational().Cmp(y.rational())
				case x > 0: // x is +inf
					cmp = -1
				default: // x is -inf
					cmp = +1
				}
				return threeway(op, cmp)
			}
		case Int:
			// x (int) < y (float)
			if y, ok := y.(Float); ok {
				if y != y { // y is NaN
					return false
				}
				var cmp int
				switch {
				case !math.IsInf(float64(y), 0): // y is finite
					cmp = x.rational().Cmp(y.rational())
				case y > 0: // y is +inf
					cmp = -1
				default: // y is -inf
					cmp = +1
				}
				return threeway(op, cmp)
			}
			// x (int) < y (int)
			if y, ok := y.(Int); ok {
				return x < y
			}
		}
		// try __lt
		event = metaLt
	}
	if !raw {
		// try metamethod event
		val, err := tryMetaCompare(state, x, y, event)
		if err != nil {
			panic(runtimeErr(err))
		}
		return val
	}
	return false
}

// arith performs an arithmetic or bitwise operation.
//
// # Arithmetic Operators (See https://www.lua.org/manual/5.3/manual.html#3.4.1)
//
// With the exception of exponentiation (^) and float division (/), the arithmetic operators
// work as follows: If both operands are integers, the operation is performed over integers
// and the result is an integer. Otherwise, if both operands are numbers or strings that can
// be converted to numbers (see §3.4.3), then they are converted to floats, the operation is
// performed following the usual rules for floating-point arithmetic (usually IEEE 754), and
// the result is a float.
//
// Exponentiation and float division always convert their operands to floats and the result is
// always a float. Exponentiation uses the ISO C function "pow", so it works across non-integer
// exponents too.
//
// Floor division (//) is a division that rounds the quotient towards minus infinity, that is,
// the floor of the division of its operands.
//
// Modulo is defined as the remainder of a division that rounds the quotient towards minus
// infinity (floor division).
//
// In case of overflows in integer arithmetic, all operations wrap around, according to the usual
// rules of two-complement arithmetic. In other words, they return the unique representable integer
// that is equal module 2^64 to the mathematical result.
//
// # Bitwise Operators (See https://www.lua.org/manual/5.3/manual.html#3.4.2)
//
// All bitwise operations convert its operands to integers (see §3.4.3), operate on all bits of those
// integers, and result in an integer.
//
// Both right and left shifts fill the vacant bits with zeros. Negative displacements shift to the
// other direction; displacements with absolute values equal to or higher than the number of bits
// in an integer result in zero (as all bits are shifted out).
//
// Returns the value and nil if successful; otherwise nil and the error.
func (state *State) arith(op Op, x, y Value) Value {
	// metamethod event to call if type op type pairs are exhausted.
	var event metaEvent

	switch op {
	//
	// Arithmetic Operators
	//
	case OpAdd: // '+'
		if isInteger(x) && isInteger(y) {
			return x.(Int) + y.(Int)
		}
		if n1, ok := toFloat(x); ok {
			if n2, ok := toFloat(y); ok {
				return n1 + n2
			}
		}
		// try __add
		event = metaAdd

	case OpSub: // '-'
		if isInteger(x) && isInteger(y) {
			return x.(Int) - y.(Int)
		}
		if n1, ok := toFloat(x); ok {
			if n2, ok := toFloat(y); ok {
				return n1 - n2
			}
		}
		// try __sub
		event = metaSub

	case OpMul: // '*'
		if isInteger(x) && isInteger(y) {
			return x.(Int) * y.(Int)
		}
		if n1, ok := toFloat(x); ok {
			if n2, ok := toFloat(y); ok {
				return n1 * n2
			}
		}
		// try __mul
		event = metaMul

	case OpMod: // '%'
		if isInteger(x) && isInteger(y) {
			m, _ := toInteger(x)
			n, _ := toInteger(y)
			if n == 0 {
				panic(fmt.Errorf("attempt to perform n%%0"))
			}
			if n == -1 {
				return Int(0)
			}
			r := Int(m % n)
			if r != 0 && (m^n) < 0 {
				r += n
			}
			return r
		}
		if n1, ok := toFloat(x); ok {
			if n2, ok := toFloat(y); ok {
				r := math.Mod(float64(n1), float64(n2))
				return Float(r)
			}
		}
		// try __mod
		event = metaMod

	case OpQuo: // '/'
		if isInteger(x) && isInteger(y) {
			m, _ := toInteger(x)
			n, _ := toInteger(y)
			if n == 0 {
				panic(fmt.Errorf("attempt to divide by zero"))
			}
			if n == -1 {
				return Int(0 - m)
			}
			q := Int(m / n)
			if (m^n) < 0 && m%n != 0 {
				q -= 1
			}
			return q
		}
		if n1, ok := toFloat(x); ok {
			if n2, ok := toFloat(y); ok {
				r := math.Floor(float64(n1 / n2))
				return Float(r)
			}
		}
		// try __idiv
		event = metaIdiv

	case OpDiv: // '//'
		if n1, ok := toFloat(x); ok {
			if n2, ok := toFloat(y); ok {
				return n1 / n2
			}
		}
		// try __div
		event = metaDiv

	case OpPow: // '^'
		if n1, ok := toFloat(x); ok {
			if n2, ok := toFloat(y); ok {
				r := math.Pow(float64(n1), float64(n2))
				return Float(r)
			}
		}
		// try __pow
		event = metaPow

	case OpMinus: // '-' (unary)
		if isInteger(x) {
			return -(x.(Int))
		}
		if n, ok := toFloat(x); ok {
			return -n
		}
		// try __unm
		event = metaUnm

	//
	// Bitwise Operators (Integers only)
	//
	case OpOr: // '|'
		if i1, ok := toInteger(x); ok {
			if i2, ok := toInteger(y); ok {
				return i1 | i2
			}
		}
		// try __bor
		event = metaBor

	case OpAnd: // '&'
		if i1, ok := toInteger(x); ok {
			if i2, ok := toInteger(y); ok {
				return i1 & i2
			}
		}
		// try __band
		event = metaBand

	case OpXor: // '~'
		if i1, ok := toInteger(x); ok {
			if i2, ok := toInteger(y); ok {
				return i1 ^ i2
			}
		}
		// try __bxor
		event = metaBxor

	case OpRsh: // '>>'
		if i1, ok := toInteger(x); ok {
			if i2, ok := toInteger(y); ok {
				return shiftRight(i1, i2)
			}
		}
		// try __shr
		event = metaShr

	case OpLsh: // '<<'
		if i1, ok := toInteger(x); ok {
			if i2, ok := toInteger(y); ok {
				return shiftLeft(i1, i2)
			}
		}
		// try __shl
		event = metaShl

	case OpNot: // '~' (unary)
		if i1, ok := toInteger(x); ok {
			return ^i1
		}
		// try __bnot
		event = metaBnot
	}
	// try metamethod event
	val, err := tryMetaBinary(state, x, y, event)
	if err != nil {
		panic(runtimeErr(err))
	}
	return val
}

// length returns the length of the object.
//
// The length operator is denoted by the unary prefix operator #.
//
// The length of a string is its number of bytes (that is, the usual meaning of string
// length when each character is one byte).
//
// The length operator applied on a table returns a border in that table. A border in a
// table t is any natural number that satisfies the following condition:
//
//      (border == 0 or t[border] ~= nil) and t[border + 1] == nil
//
// In words, a border is any (natural) index in a table where a non-nil value is followed
// by a nil value (or zero, when index 1 is nil).
//
// A table with exactly one border is called a sequence. For instance, the table {10, 20, 30, 40, 50}
// is a sequence, as it has only one border (5). The table {10, 20, 30, nil, 50} has two borders
// (3 and 5), and therefore it is not a sequence. The table {nil, 20, 30, nil, nil, 60, nil} has
// three borders (0, 3, and 6), so it is not a sequence, too. The table {} is a sequence with border 0.
// Note that non-natural keys do not interfere with whether a table is a sequence.
//
// When t is a sequence, #t returns its only border, which corresponds to the intuitive notion of the length
// of the sequence. When t is not a sequence, #t can return any of its borders. (The exact one depends on
// details of the internal representation of the table, which in turn can depend on how the table was populated
// and the memory addresses of its non-numeric keys.)
//
// The computation of the length of a table has a guaranteed worst time of O(log n), where n is the largest
// natural key in the table.
//
// A program can modify the behavior of the length operator for any value but strings through the __len
// metamethod (see §2.4).
//
// See https://www.lua.org/manual/5.3/manual.html#3.4.7
func (state *State) length(obj Value) Value {
	if str, ok := obj.(String); ok {
		return Int(len(str))
	}
	val, err := tryMetaLength(state, obj)
	if err == nil {
		return val
	}
	if tbl, ok := obj.(*table); ok {
		return Int(tbl.length())
	}
	panic(runtimeErr(err))
}

// concat returns the concatenation of values.
func (state *State) concat(values []Value) Value {
	rhs := values[len(values)-1]
	for i := len(values) - 2; i >= 0; i-- {
		lhs := values[i]
		switch lhs.(type) {
		case String, Number:
			if _, ok := rhs.(String); ok {
				s1, _ := toString(lhs)
				s2, _ := toString(rhs)
				rhs = String(s1 + s2)
				continue
			}
			if _, ok := rhs.(Number); ok {
				s1, _ := toString(lhs)
				s2, _ := toString(rhs)
				rhs = String(s1 + s2)
				continue
			}
		}
		var err error
		if rhs, err = tryMetaConcat(state, lhs, rhs); err != nil {
			state.Errorf("%v", err)
		}
	}
	return rhs
}

// tonumber converts a value to a number.
//
// Returns the number and true if successful; otherwise nil and false.
func toNumber(v Value) (Number, bool) {
	switch v := v.(type) {
	case String:
		return str2num(strings.ToLower(strings.TrimSpace(string(v))))
	case Float:
		return v, true
	case Int:
		return v, true
	}
	return nil, false
}

// toInteger converts a value to an integer.
//
// Returns the integer and true if successful; otherwise 0 and false.
func toInteger(v Value) (Int, bool) {
	if num, ok := toNumber(v); ok {
		switch num := num.(type) {
		case Float:
			if float64(int64(num)) == float64(num) {
				return Int(num), true
			}
		case Int:
			return num, true
		}
	}
	return Int(0), false
}

// toFloat converts a value to a float.
//
// Returns the float and true if successful; otherwise 0 and false.
func toFloat(v Value) (Float, bool) {
	if num, ok := toNumber(v); ok {
		switch num := num.(type) {
		case Float:
			return num, true
		case Int:
			return Float(num), true
		}
	}
	return Float(0), false
}

// isInteger returns true if v is an Int.
func isInteger(v Value) bool { _, ok := v.(Int); return ok }

// isNumber returns true if v is a Number.
func isNumber(v Value) bool { _, ok := v.(Number); return ok }

// isNaN returns whether the Lua float number is an IEEE 754 "not-a-number" value.
func isNaN(v Float) bool { return float64(v) == math.NaN() }

// threeway interprets a three-way comparsion value cmp (-1, 0, +1) as a boolean
// comparison between two values (e.g. x < y).
func threeway(op Op, cmp int) bool {
	switch op {
	case OpLt:
		return cmp < 0
	case OpLe:
		return cmp <= 0
	case OpEq:
		return cmp == 0
	case OpGt:
		return cmp > 0
	case OpGe:
		return cmp >= 0
	case OpNe:
		return cmp != 0
	}
	panic(op)
}

func str2num(str string) (num Number, ok bool) {
	if num, ok = str2int(str); ok {
		return num, ok
	}
	if num, ok = str2float(str); ok {
		return num, ok
	}
	return nil, false
}

func str2int(str string) (Int, bool) {
	// num, err := strconv.ParseInt(str, 0, 64)
	// fmt.Println(err)
	// return Int(num), err == nil
	i64, ok := syntax.StrToI64(str)
	return Int(i64), ok
}

func str2float(str string) (Float, bool) {
	// if strings.Contains(str, "nan") || strings.Contains(str, "inf") {
	//     return 0, false
	// }
	// // num, err := strconv.ParseUint(str, 0, 64)
	// // if err != nil {
	// //     return Float(0), false
	// // }
	// // return Float(math.Float64frombits(num)), true
	// num, err := strconv.ParseFloat(str, 64)
	// return Float(num), err == nil
	f64, ok := syntax.StrToF64(str)
	return Float(f64), ok
}

func shiftLeft(x, y Int) Int {
	if y >= 0 {
		return x << uint64(y)
	}
	return shiftRight(x, -y)
}

func shiftRight(x, y Int) Int {
	if y >= 0 {
		return Int(uint64(x) >> uint64(y))
	}
	return shiftLeft(x, -y)
}
