package lua

import (
	"fmt"
)

type Error interface {
	error
	format(*Thread) string
	traceback(*Thread) string
}

// "'typename' expected, got 'typename'"
func TypeErr(arg int, typ, want string) error {
	return ArgErr(arg, &typeErr{typ, want})
}

// "bad argument #arg to 'funcname' (extramsg)"
func ArgErr(arg int, err error) error {
	return &argErr{arg, err}
}

type (
	runtimeErr struct {
		fr  *frame
		err error
	}

	evalErr struct {
		ctx string
		err error
	}

	typeErr struct {
		typ  string
		want string
	}

	argErr struct {
		arg int
		err error
	}
)

func (*runtimeErr) Error() string { return "runtime error!" }

func (e *evalErr) Error() string {
	return "eval error!"
}

func (e *typeErr) Error() string {
	return fmt.Sprintf("%s expected, got %s", e.want, e.typ)
}

func (e *argErr) Error() string {
	// "bad argument #arg to 'funcname' (extramsg)"
	return fmt.Sprintf("bad argument #%d (%v)", e.arg, e.err)
}
