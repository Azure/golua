package lua5

import (
	"math"

	"github.com/Azure/golua/lua"
)

// math.abs(x)
//
// Returns the absolute value of x. (integer/float)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.abs
func math۰abs(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	num, err := args.Number(0)
	if err != nil {
		return nil, err
	}
	var abs lua.Value
	switch num := num.(type) {
	case lua.Float:
		abs = lua.Float(math.Abs(float64(num)))
	case lua.Int:
		if abs = num; num < 0 {
			abs = -num
		}
	}
	return []lua.Value{abs}, nil
}

// math.acos(x)
//
// Returns the arc cosine of x (in radians).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.acos
func math۰acos(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.asin(x)
//
// Returns the arc sine of x (in radians).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.asin
func math۰asin(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.atan(y [, x])
//
// Returns the arc tangent of y/x (in radians), but uses the
// signs of both arguments to find the quadrant of the result.
// (It also handles correctly the case of x being zero.)
//
// The default value for x is 1, so that the call math.atan(y)
// returns the arc tangent of y.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.atan
func math۰atan(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.ceil(x)
//
// Returns the smallest integral value larger than or equal to x.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.ceil
func math۰ceil(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.cos(x)
//
// Returns the cosine of x (assumed to be in radians).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.cos
func math۰cos(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.deg(x)
//
// Converts the angle x from radians to degrees.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.deg
func math۰deg(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.exp(x)
//
// Returns the value ex (where e is the base of natural logarithms).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.exp
func math۰exp(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.floor(x)
//
// Returns the largest integral value smaller than or equal to x.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.floor
func math۰floor(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.fmod(x, y)
//
// Returns the remainder of the division of x by y that rounds the
// quotient towards zero. (integer/float)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.fmod
func math۰fmod(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.log(x [, base])
//
// Returns the logarithm of x in the given base. The default for base is e
// (so that the function returns the natural logarithm of x).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.log
func math۰log(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.max(x, ···)
//
// Returns the argument with the maximum value, according to
// the Lua operator <. (integer/float)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.max
func math۰max(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.min(x, ···)
//
// Returns the argument with the minimum value, according to
// the Lua operator <. (integer/float)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.min
func math۰min(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.modf(x)
//
// Returns the integral part of x and the fractional part of x.
// Its second result is always a float.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.modf
func math۰modf(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.rad(x)
//
// Converts the angle x from degrees to radians.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.rad
func math۰rad(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.random([m [, n]])
//
// When called without arguments, returns a pseudo-random float with uniform
// distribution in the range [0,1). When called with two integers m and n,
// math.random returns a pseudo-random integer with uniform distribution in
// the range [m, n]. (The value n-m cannot be negative and must fit in a Lua
// integer.) The call math.random(n) is equivalent to math.random(1,n).
//
// This function is an interface to the underling pseudo-random generator
// function provided by Go.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.random
func math۰random(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.randomseed(x)
//
// Sets x as the "seed" for the pseudo-random generator: equal seeds produce equal
// sequences of numbers.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.randomseed
func math۰randomseed(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.sin(x)
//
// Returns the sine of x (assumed to be in radians).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.sin
func math۰sin(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.sqrt(x)
//
// Returns the square root of x. (You can also use the expression x^0.5
// to compute this value.)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.sqrt
func math۰sqrt(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.tan(x)
//
// Returns the tangent of x (assumed to be in radians).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.tan
func math۰tan(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.tointeger(x)
//
// If the value x is convertible to an integer, returns that integer.
// Otherwise, returns nil.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.tointeger
func math۰tointeger(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.type(x)
//
// Returns "integer" if x is an integer, "float" if it is a float, or
// nil if x is not a number.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.type
func math۰type(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// math.ult (m, n)
//
// Returns a boolean, true if and only if integer m is below integer n
// when they are compared as unsigned integers.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.ult
func math۰ult(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	return nil, nil
}

// This library provides basic mathematical functions. It provides all its functions
// and constants inside the table math. Functions with the annotation "integer/float"
// give integer results for integer arguments and float results for float (or mixed)
// arguments. Rounding functions (math.ceil, math.floor, and math.modf) return an
// integer when the result fits in the range of an integer, or a float otherwise.
//
// See https://www.lua.org/manual/5.3/manual.html#6.7
func stdlib۰math(ls *lua.Thread) (lua.Value, error) {
	return lua.NewTableFromMap(map[string]lua.Value{
		// constants
		"maxinteger": lua.Int(math.MaxInt64),
		"mininteger": lua.Int(math.MinInt64),
		"huge":       lua.Float(math.Inf(1)),
		"pi":         lua.Float(math.Pi),
		// functions
		"abs":        lua.Closure(math۰abs),
		"acos":       lua.Closure(math۰acos),
		"asin":       lua.Closure(math۰asin),
		"atan":       lua.Closure(math۰atan),
		"ceil":       lua.Closure(math۰ceil),
		"cos":        lua.Closure(math۰cos),
		"deg":        lua.Closure(math۰deg),
		"exp":        lua.Closure(math۰exp),
		"floor":      lua.Closure(math۰floor),
		"fmod":       lua.Closure(math۰fmod),
		"log":        lua.Closure(math۰log),
		"max":        lua.Closure(math۰max),
		"min":        lua.Closure(math۰min),
		"modf":       lua.Closure(math۰modf),
		"rad":        lua.Closure(math۰rad),
		"random":     lua.Closure(math۰random),
		"randomseed": lua.Closure(math۰randomseed),
		"sin":        lua.Closure(math۰sin),
		"sqrt":       lua.Closure(math۰sqrt),
		"tan":        lua.Closure(math۰tan),
		"tointeger":  lua.Closure(math۰tointeger),
		"type":       lua.Closure(math۰type),
		"ult":        lua.Closure(math۰ult),
	}), nil
}
