package stdlib

import (
    "math/rand"
    "math"
    "fmt"
    "os"
    "github.com/Azure/golua/lua"
)

var _ = fmt.Println
var _ = os.Exit

//
// Lua Standard Library -- math
//

// OpenMath opens the Lua standard math library. This library provides basic mathematical functions.
// It provides all its functions and constants inside the table math. Functions with the annotation
// "integer/float" give integer results for integer arguments and float results for float (or mixed)
// arguments. Rounding functions (math.ceil, math.floor, and math.modf) return an integer when the
// result fits in the range of an integer, or a float otherwise.
//
// See https://www.lua.org/manual/5.3/manual.html#6.7
func OpenMath(state *lua.State) int {
	// Create 'math' table
    var mathFuncs = map[string]lua.Func{
		"abs": 		  lua.Func(mathAbs),
		"acos": 	  lua.Func(mathAcos),
		"asin": 	  lua.Func(mathAsin),
		"atan": 	  lua.Func(mathAtan),
		"ceil": 	  lua.Func(mathCeil),
		"cos": 		  lua.Func(mathCos),
		"deg": 		  lua.Func(mathDeg),
		"exp": 		  lua.Func(mathExp),
		"floor": 	  lua.Func(mathFloor),
		"fmod": 	  lua.Func(mathFmod),
		"log": 		  lua.Func(mathLog),
		"max": 		  lua.Func(mathMax),
		"min": 		  lua.Func(mathMin),
		"modf": 	  lua.Func(mathModf),
		"rad": 		  lua.Func(mathRad),
		"random": 	  lua.Func(mathRand),
		"randomseed": lua.Func(mathRandSeed),
		"sin": 		  lua.Func(mathSin),
		"sqrt": 	  lua.Func(mathSqrt),
		"tan": 		  lua.Func(mathTan),
		"tointeger":  lua.Func(mathToInt),
		"type": 	  lua.Func(mathType),
		"ult": 		  lua.Func(mathUlt),
    }
	state.NewTableSize(0, len(mathFuncs))
    state.SetFuncs(mathFuncs, 0)

    // Set 'pi' field.
    state.Push(math.Pi)
    state.SetField(-2, "pi") // The value of π.

    // Set 'huge' field.
    state.Push(math.MaxFloat64)
    state.SetField(-2, "huge") // A value larger than any other numeric value.

    // Set 'maxinteger' field.
    state.Push(math.MaxInt64)
    state.SetField(-2, "maxinteger") // An integer with the maximum value for an integer.

    // Set 'mininteger' field.
    state.Push(math.MinInt64)
    state.SetField(-2, "mininteger") // An integer with the minimum value for an integer.

	// Return 'math' table
    return 1
}

// math.abs (x)
//
// Returns the absolute value of x. (integer/float)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.abs
func mathAbs(state *lua.State) int {
    unimplemented("math.abs")
    return 0
}

// math.acos (x)
//
// Returns the arc cosine of x (in radians).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.acos
func mathAcos(state *lua.State) int {
    unimplemented("math.acos")
    return 0
}

// math.asin (x)
//
// Returns the arc sine of x (in radians).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.asin
func mathAsin(state *lua.State) int {
    unimplemented("math.asin")
    return 0
}

// math.atan (y [, x])
//
// Returns the arc tangent of y/x (in radians), but uses the signs of both arguments to find the
// quadrant of the result. (It also handles correctly the case of x being zero.)
//
//
// The default value for x is 1, so that the call math.atan(y) returns the arc tangent of y.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.atan
func mathAtan(state *lua.State) int {
    unimplemented("math.atan")
    return 0
}

// math.ceil (x)
//
// Returns the smallest integral value larger than or equal to x.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.ceil
func mathCeil(state *lua.State) int {
    unimplemented("math.ceil")
    return 0
}

// math.cos (x)
//
// Returns the cosine of x (assumed to be in radians).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.cos
func mathCos(state *lua.State) int {
    unimplemented("math.cos")
    return 0
}

// math.deg (x)
//
// Converts the angle x from radians to degrees.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.deg
func mathDeg(state *lua.State) int {
    unimplemented("math.deg")
    return 0
}

// math.exp (x)
//
// Returns the value ex (where e is the base of natural logarithms).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.exp
func mathExp(state *lua.State) int {
    unimplemented("math.exp")
    return 0
}

// math.floor (x)
//
// Returns the largest integral value smaller than or equal to x.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.floor
func mathFloor(state *lua.State) int {
    unimplemented("math.floor")
    return 0
}

// math.fmod (x, y)
//
// Returns the remainder of the division of x by y that rounds the quotient towards zero. (integer/float)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.fmod
func mathFmod(state *lua.State) int {
    unimplemented("math.fmod")
    return 0
}

// math.log (x [, base])
//
// Returns the logarithm of x in the given base. The default for base is e (so that the function returns
// the natural logarithm of x).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.log
func mathLog(state *lua.State) int {
    unimplemented("math.log")
    return 0
}

// math.max (x, ···)
//
// Returns the argument with the maximum value, according to the Lua operator <. (integer/float)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.max
func mathMax(state *lua.State) int {
    unimplemented("math.max")
    return 0
}

// math.min (x, ···)
//
// Returns the argument with the minimum value, according to the Lua operator <. (integer/float)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.min
func mathMin(state *lua.State) int {
    unimplemented("math.min")
    return 0
}

// math.modf (x)
//
// Returns the integral part of x and the fractional part of x. Its second result is always a float.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.modf
func mathModf(state *lua.State) int {
    unimplemented("math.modf")
    return 0
}

// math.rad (x)
//
// Converts the angle x from degrees to radians.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.rad
func mathRad(state *lua.State) int {
    unimplemented("math.rad")
    return 0
}

// math.random ([m [, n]])
// 
// When called without arguments, returns a pseudo-random float with uniform distribution in the range [0,1).
// When called with two integers m and n, math.random returns a pseudo-random integer with uniform distribution
// in the range [m, n]. (The value n-m cannot be negative and must fit in a Lua integer.)
// 
// The call math.random(n) is equivalent to math.random(1,n).
// 
// This function is an interface to the underling pseudo-random generator function provided by Go.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.random
func mathRand(state *lua.State) int {
    unimplemented("math.random")
    return 0
}

// math.randomseed (x)
//
// Sets x as the "seed" for the pseudo-random generator: equal seeds produce equal sequences of numbers.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.randomseed
func mathRandSeed(state *lua.State) int {
    rand.Seed(state.CheckInt(1))
    return 0
}

// math.sin (x)
//
// Returns the sine of x (assumed to be in radians).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.sin
func mathSin(state *lua.State) int {
    unimplemented("math.sin")
    return 0
}

// math.sqrt (x)
//
// Returns the square root of x. (You can also use the expression x^0.5 to compute this value.)
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.sqrt
func mathSqrt(state *lua.State) int {
    unimplemented("math.sqrt")
    return 0
}

// math.tan (x)
//
// Returns the tangent of x (assumed to be in radians).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.tan
func mathTan(state *lua.State) int {
    unimplemented("math.tan")
    return 0
}

// math.tointeger (x)
//
// If the value x is convertible to an integer, returns that integer. Otherwise, returns nil.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.tointeger
func mathToInt(state *lua.State) int {
    unimplemented("math.tointeger")
    return 0
}

// math.type (x)
//
// Returns "integer" if x is an integer, "float" if it is a float, or nil if x is not a number.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.type
func mathType(state *lua.State) int {
    unimplemented("math.type")
    return 0
}

// math.ult (m, n)
//
// Returns a boolean, true if and only if integer m is below integer n when they are compared as unsigned integers.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-math.ult
func mathUlt(state *lua.State) int {
    unimplemented("math.ult")
    return 0
}
