package lua

import (
	"strconv"
	"math"
    "fmt"
)

var _ = fmt.Println

// Op represents a Lua arithmetic, bitwise, or relational operator.
type Op int

const (
    // Arithmetic Operators
    OpAdd      Op = 1 + iota  // '+'
    OpSub                     // '-'
    OpMul                     // '*'
    OpMod                     // '%'
    OpPow                     // '^'
    OpDiv                     // '/'
    OpQuo                     // '//'
    OpMinus                   // '-' (unary)

    // Bitwise Operators
    OpOr                      // '|'
    OpAnd                     // '&'
    OpXor                     // '~'
    OpRsh                     // '>>'
    OpLsh                     // '<<'
    OpNot                     // '~' (unary)

    // Relational Operators
    OpLt                     // '<'
    OpLe                     // '<='
    OpEq                     // '=='
    OpGt                     // '>'
    OpGe                     // '>='
    OpNe                     // '~='

    // Miscellaneous Operators
    OpConcat                 // '..'
    OpLength                 // '#'
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
func (state *State) compare(op Op, x, y Value) bool {
    // metamethod event to call if type op type pairs are exhausted.
    var event metaEvent

    switch op {
        //case OpNe: // '~='
        //case OpGe: // '>='
        //case OpGt: // '>'

        case OpEq: // '=='
            if x.Type() != y.Type() {
                return false
            }
            switch x := x.(type) {
                case String:
                    return x == y.(String)
                case Float:
                    return x == y.(Float)
                case Bool:
                    return x == y.(Bool)
                case Int:
                    return x == y.(Int)
                case Nil:
                    return true
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
                        return x <= Float(y)
                    }
                case Int:
                    // x (integer) <= y (float)
                    if y, ok := y.(Float); ok {
                        return Float(x) <= y
                    }
                    // x (integer) <= y (integer)
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
                        return x < Float(y)
                    }
                case Int:
                    // x (integer) < y (float)
                    if y, ok := y.(Float); ok {
                        return Float(x) < y
                    }
                    // x (integer) < y (integer)
                    if y, ok := y.(Int); ok {
                        return x < y
                    }
            }
            // try __lt
            event = metaLt
    }
    // try metamethod
    state.errorf("ops: todo: call relational metamethod: %v(%T,%T)", event.name(), x, y)
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
                return x.(Int) % y.(Int)
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
                return x.(Int) % y.(Int)
            }
            if n1, ok := toFloat(x); ok {
                if n2, ok := toFloat(y); ok {
                    r := math.Pow(float64(n1), float64(n2))
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
                    return i1 >> uint(i2)
                }
            }
            // try __shr
            event = metaShr

        case OpLsh: // '<<'
            if i1, ok := toInteger(x); ok {
                if i2, ok := toInteger(y); ok {
                    return i1 << uint(i2)
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
func (state *State) length(obj Value) int {
    var meta Value
    switch obj := obj.(type) {
        case *Table:
            if meta = state.metafield(obj, "__len"); !IsNone(meta) {
                break
            }
            return len(obj.list)
        case String:
            return len(string(obj))
        default: // try metamethod
            if meta = state.metafield(obj, "__len"); IsNone(meta) {
                state.errorf("attempt to get length of %v value", obj.Type())
            }
    }
    state.errorf("ops: todo: call __len metamethod")
    return 0
}

// toInteger converts a value to an integer.
//
// Returns the integer and true if successful; otherwise 0 and false.
func toInteger(v Value) (Int, bool) {
    switch v := v.(type) {
        case String:
            // i64, err := strconv.ParseInt(string(v), 10, 64)
            // if err != nil {
            //     return Int(0), false
            // }
            // return Int(i64), true
            n, ok := toNumber(v)
            if ok {
                return Int(n.(Float)), true
            }  
        case Float:
           return Int(v), true
        case Int:
            return v, true
    }
    return Int(0), false
}

// tonumber converts a value to a number.
//
// Returns the number and true if successful; otherwise nil and false.
func toNumber(v Value) (Number, bool) {
    switch v := v.(type) {
        case String:
            f64, err := strconv.ParseFloat(string(v), 64)
            if err != nil {
                return nil, false
            }
            return Float(f64), true
        case Float:
            return v, true
        case Int:
            return Float(v), true
    }
    return nil, false
}

// concat returns the concatenation of values.
func (state *State) concat(values []Value) Value {
    lhs := values[0]
    for i := 1; i < len(values); i++ {
        rhs := values[i]
        if s1, ok := toString(lhs); ok {
            if s2, ok := toString(rhs); ok {
                lhs = String(s1 + s2)
                continue
            }
        }
        var err error
        if lhs, err = tryMetaConcat(state, lhs, rhs); err != nil {
            state.errorf("%v", err)
        }
    }
    return lhs
}

// toFloat converts a value to a float.
//
// Returns the float and true if successful; otherwise 0 and false.
func toFloat(v Value) (Float, bool) {
    num, ok := toNumber(v)
    if !ok {
        return Float(0), false
    }
    return num.(Float), true
}

// isInteger returns true if v is an Int.
func isInteger(v Value) bool { _, ok := v.(Int); return ok }

// isNumber returns true if v is a Number.
func isNumber(v Value) bool { _, ok := v.(Number); return ok }

// isNaN returns whether the Lua float number is an IEEE 754 "not-a-number" value.
func isNaN(v Float) bool { return float64(v) == math.NaN() }