package lua5

import (
	"fmt"

	"github.com/Azure/golua/lua"
)

var _ = fmt.Println

// coroutine.isyieldable()
//
// Returns true when the running coroutine can yield.
//
// A running coroutine is yieldable if it is not the main thread and it is not inside a non-yieldable
// Go function.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-coroutine.isyieldable
func coroutine۰isyieldable(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	fmt.Printf("coroutine۰isyieldable(%v)\n", args)
	return nil, nil
}

// coroutine.create(f)
//
// Creates a new coroutine, with body f. f must be a function. Returns this new coroutine, an object
// with type "thread".
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-coroutine.create
func coroutine۰create(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	fmt.Printf("coroutine۰create(%v)\n", args)
	return nil, nil
}

// coroutine.resume(co [, val1, ···])
//
// Starts or continues the execution of coroutine co. The first time you resume a coroutine, it starts running
// its body. The values val1, ... are passed as the arguments to the body function. If the coroutine has yielded,
// resume restarts it; the values val1, ... are passed as the results from the yield.
//
// If the coroutine runs without any errors, resume returns true plus any values passed to yield (when the
// coroutine yields) or any values returned by the body function (when the coroutine terminates). If there
// is any error, resume returns false plus the error message.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-coroutine.resume
func coroutine۰resume(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	fmt.Printf("coroutine۰resume(%v)\n", args)
	return nil, nil
}

// coroutine.running()
//
// Returns the running coroutine plus a boolean, true when the running coroutine is the main one.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-coroutine.running
func coroutine۰running(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	fmt.Printf("coroutine۰running(%v)\n", args)
	return nil, nil
}

// coroutine.status(co)
//
// Returns the status of coroutine co, as a string: "running", if the coroutine is running (that is,
// it called status); "suspended", if the coroutine is suspended in a call to yield, or if it has not
// started running yet; "normal" if the coroutine is active but not running (that is, it has resumed
// another coroutine); and "dead" if the coroutine has finished its body function, or if it has stopped
// with an error.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-coroutine.status
func coroutine۰status(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	fmt.Printf("coroutine۰status(%v)\n", args)
	return nil, nil
}

// coroutine.wrap(f)
//
// Creates a new coroutine, with body f. f must be a function. Returns a function that resumes the coroutine
// each time it is called. Any arguments passed to the function behave as the extra arguments to resume.
// Returns the same values returned by resume, except the first boolean. In case of error, propagates the error.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-coroutine.wrap
func coroutine۰wrap(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	fmt.Printf("coroutine۰wrap(%v)\n", args)
	return nil, nil
}

// coroutine.yield(···)
//
// Suspends the execution of the calling coroutine. Any arguments to yield are passed as extra results to resume.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-coroutine.yield
func coroutine۰yield(ls *lua.Thread, args lua.Tuple) ([]lua.Value, error) {
	fmt.Printf("coroutine۰yield(%v)\n", args)
	return nil, nil
}

func stdlib۰coroutine(ls *lua.Thread) (lua.Value, error) {
	// coroutine.create
	// coroutine.isyieldable
	// coroutine.resume
	// coroutine.running
	// coroutine.status
	// coroutine.wrap
	// coroutine.yield
	return lua.NewTableFromMap(map[string]lua.Value{
		"isyieldable": lua.Closure(coroutine۰isyieldable),
		"running":     lua.Closure(coroutine۰running),
		"status":      lua.Closure(coroutine۰status),
		"create":      lua.Closure(coroutine۰create),
		"resume":      lua.Closure(coroutine۰resume),
		"yield":       lua.Closure(coroutine۰yield),
		"wrap":        lua.Closure(coroutine۰wrap),
	}), nil
}
