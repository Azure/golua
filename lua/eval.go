package lua

import (
	"fmt"
	"math"
	"reflect"
)

func Compare(t *Thread, op Op, x, y Value) (bool, error) {
	return compare(t.ls, op, x, y)
}

func Equals(t *Thread, x, y Value) (bool, error) {
	return compare(t.ls, OpEq, x, y)
}

func Length(t *Thread, x Value) (Int, error) {
	v, err := Unary(t, OpLen, x)
	if err != nil {
		return 0, err
	}
	n, ok := ToInt(v)
	if !ok {
		// What is the error?
		return 0, fmt.Errorf("length!")
	}
	return n, nil
}

func Binary(t *Thread, op Op, x, y Value) (Value, error) {
	return binary(t.ls, op, x, y)
}

func Unary(t *Thread, op Op, x Value) (Value, error) {
	return unary(t.ls, op, x)
}

// func compare(ls *thread, op Op, x, y Value) (bool, error) {
// 	switch op {
// 		case OpEq:
// 			return Equals(x, y), nil
// 		case OpLe:
// 			// if x, ok := x.(Number); ok {
// 			// 	if y, ok := y.(Number); ok {
// 			// 		return lesseq(x, y), nil
// 			// 	}
// 			// }
// 			// if x, ok := x.(String); ok {
// 			// 	if y, ok := y.(String); ok {
// 			// 		return x <= y, nil
// 			// 	}
// 			// }
// 			panic("Compare: OpLe")

// 		case OpLt:
// 			// if x, ok := x.(Number); ok {
// 			// 	if y, ok := y.(Number); ok {
// 			// 		return less(x, y), nil
// 			// 	}
// 			// }
// 			// if x, ok := x.(String); ok {
// 			// 	if y, ok := y.(String); ok {
// 			// 		return x < y, nil
// 			// 	}
// 			// }
// 			panic("Compare: OpLt")

// 		// case OpNe:
// 		// case OpGt:
// 		// case OpGe:
// 		default:
// 			panic(op)
// 	}
// 	// return false, &orderErr{op, x, y}
// 	return false, nil
// }

// func equals(ls *thread, x, y Value) (bool, error) {
// 	if !sameType(x, y) {
// 		switch x := x.(type) {
// 			case Float:
// 				if y, ok := y.(Int); ok {
// 					return x == Float(y)
// 				}
// 				y, ok := y.(Float)
// 				return ok && (x == y)
// 			case Int:
// 				if y, ok := y.(Float); ok {
// 					return Float(x) == y
// 				}
// 				y, ok := y.(Int)
// 				return ok && (x == y)
// 			default:
// 				return false
// 		}
// 	}
// 	switch x := x.(type) {
// 		case *GoFunc:
// 			y, ok := y.(*GoFunc)
// 			return ok && (x == y)
// 		case *Func:
// 			y, ok := y.(*Func)
// 			return ok && (x == y)
// 		case *Table:
// 			y, ok := y.(*Table)
// 			return ok && (x == y)
// 		case String:
// 			y, ok := y.(String)
// 			return ok && (x == y)
// 		case Float:
// 			y, ok := y.(Float)
// 			return ok && (x == y)
// 		case Bool:
// 			y, ok := y.(Bool)
// 			return ok && (x == y)
// 		case Int:
// 			y, ok := y.(Int)
// 			return ok && (x == y)
// 		case nil:
// 			return (y == nil)
// 		default:
// 			panic(fmt.Errorf("equals(%T, %T)", x, y))
// 	}
// }

// func lesseq(ls *thread, x, y Value) (bool, error) {
// 	switch x := x.(type) {
// 		case Float:
// 			if y, ok := y.(Int); ok {
// 				return x <= Float(y)
// 			}
// 			return x < y.(Float)
// 		case Int:
// 			if y, ok := y.(Float); ok {
// 				return Float(x) < y
// 			}
// 			return x <= y.(Int)
// 	}
// 	panic("unreachable")
// }

// func less(ls *thread, x, y Number) (bool, error) {
// 	switch x := x.(type) {
// 		case Float:
// 			if y, ok := y.(Int); ok {
// 				return x < Float(y)
// 			}
// 			return x < y.(Float)
// 		case Int:
// 			if y, ok := y.(Float); ok {
// 				return Float(x) < y
// 			}
// 			return x < y.(Int)
// 	}
// 	panic("unreachable")
// }

func binary(ls *thread, op Op, x, y Value) (Value, error) {
	switch op {
	case OpDivF, OpPow:
		if x, ok := ToFloat(x); ok {
			if y, ok := ToFloat(y); ok {
				return numop(op, x, y), nil
			}
		}

	case OpBand,
		OpBor,
		OpBxor,
		OpShl,
		OpShr,
		OpBnot:

		if x, ok := ToInt(x); ok {
			if y, ok := ToInt(y); ok {
				return intop(op, x, y), nil
			}
		}

	default:
		if x, ok := x.(Int); ok {
			if y, ok := y.(Int); ok {
				return intop(op, x, y), nil
			}
		}
		if x, ok := ToFloat(x); ok {
			if y, ok := ToFloat(y); ok {
				return numop(op, x, y), nil
			}
		}
	}
	e := event(op-OpEq) + _add
	return e.binary(ls, x, y)
}

func numop(op Op, x, y Float) Float {
	switch op {
	case OpMinus:
		return -x
	case OpDivF:
		return x / y
	case OpDivI:
		return Float(math.Floor(float64(x / y)))
	case OpAdd:
		return x + y
	case OpSub:
		return x - y
	case OpMul:
		return x * y
	case OpPow:
		f64 := math.Pow(float64(x), float64(y))
		return Float(f64)
	case OpMod:
		f64 := Float(math.Mod(float64(x), float64(y)))
		if f64*y < 0 {
			f64 += y
		}
		return f64
	}
	panic(fmt.Errorf("unexpected binary operator '%v'", op))
}

func intop(op Op, x, y Int) Int {
	switch op {
	case OpMinus:
		return -x
	case OpDivI:
		return x / y
	case OpBand:
		return x & y
	case OpBnot:
		return ^x
	case OpBxor:
		return x ^ y
	case OpBor:
		return x | y
	case OpAdd:
		return x + y
	case OpSub:
		return x - y
	case OpMul:
		return x * y
	case OpMod:
		if r := (x % y); r != 0 && (x^y) < 0 { // 'm/n' would be non-integer negative?
			r += y // correct result for different rounding
			return Int(r)
		} else {
			return Int(r)
		}
	case OpShl:
		return shiftLeft(x, y)
	case OpShr:
		return shiftRight(x, y)
	}
	panic(fmt.Errorf("unexpected binary operator '%v'", op))
}

// shift left operation
func shiftLeft(x, y Int) Int {
	if y >= 0 {
		return x << uint64(y)
	}
	return shiftRight(x, -y)
}

// shift right operation
func shiftRight(x, y Int) Int {
	if y >= 0 {
		return Int(uint64(x) >> uint64(y))
	}
	return shiftLeft(x, -y)
}

// // threeway interprets a three-way comparison value cmp (-1, 0, +1)
// // as a boolean comparison (e.g. x < y).
// func threeway(op Op, cmp int) bool {
// 	switch op {
// 		case OpEq:
// 			return cmp == 0
// 		case OpNe:
// 			return cmp != 0
// 		case OpLt:
// 			return cmp < 0
// 		case OpLe:
// 			return cmp <= 0
// 		case OpGt:
// 			return cmp > 0
// 		case OpGe:
// 			return cmp >= 0
// 	}
// 	panic(op)
// }

func sameType(x, y Value) bool {
	return reflect.TypeOf(x) == reflect.TypeOf(y) || x.kind() == y.kind()
}
