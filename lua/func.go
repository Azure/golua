package lua

import (
	"fmt"
	"github.com/fibonacci1729/golua/lua/code"
)

type (
	Function interface {
		Callable
		// Info() Debug
		// Up(int) Value
	}

	// callable is implemented by all values that are callable: *GoFunc / *Func.
	callable interface {
		call(*thread, []Value) ([]Value, error)
	}

	// closure represents a Lua/Go closure.
	closure struct {
		fn callable
		up []*upvar
	}

	// upvar represents a Lua upvalue.
	upvar struct {
		value Value
		index int
		open  bool
		next  *upvar
	}
)

// set the upvalue's inner value.
func (up *upvar) set(v Value) {
	if up.open {
		// (*up.stack)[up.level] = v
		// return
		panic("upvar.set: open!")
	}
	up.value = v
}

// get the upvalue's inner value.
func (up *upvar) get() Value {
	if up.open {
		// return (*up.stack)[up.level]
		panic("upvar.get: open!")
	}
	return up.value
}

// closure represents a closure value embedded into a callable values.
func (cls *closure) String() string {
	return fmt.Sprintf("function: %p", cls.fn)
}

// GoFunc represents a Go builtin function value.
type GoFunc struct {
	closure
	name string
	args argsCheck
	impl func(*Thread, Tuple) ([]Value, error)
}

// NewGoFunc creates and returns a new *GoFunc.
func NewGoFunc(name string, impl func(*Thread, Tuple) ([]Value, error), vars ...Value) *GoFunc {
	fn := Closure(impl, vars...)
	fn.name = name
	return fn
}

func Closure(impl func(*Thread, Tuple) ([]Value, error), vars ...Value) *GoFunc {
	fn := &GoFunc{impl: impl}
	up := make([]*upvar, len(vars))
	for i, v := range vars {
		up[i] = &upvar{value: v}
	}
	fn.closure = closure{fn, up}
	return fn
}

func (fn *GoFunc) check(ls *thread, argv []Value) (args Tuple, err error) {
	if args = Tuple(argv); fn.args != nil {
		if err = args.Check(ls.tt, fn.args); err != nil {
			return nil, err
		}
	}
	return args, nil
}

// call implements the callable interface for Go funcs. 
func (fn *GoFunc) call(ls *thread, argv []Value) ([]Value, error) {
	args, err := fn.check(ls, argv)
	if err != nil {
		return nil, err
	}
	return fn.impl(ls.tt, args)
}

// A Func represents a Lua function value.
type Func struct {
	closure
	stack []Value
	proto *code.Proto
}

// call implements the callable interface for Lua funcs. 
func (fn *Func) call(ls *thread, args []Value) ([]Value, error) {
	// fmt.Printf("call: args=%v, varg=%v\n", args, ls.fr.call.va)
	// fn.stack = append(fn.stack, args)
	for i, arg := range args {
		fn.stack[i] = arg
	}
	return ls.exec(fn)
}

// kst returns the function i'th constant.
func (fn *Func) kst(i int) (c Constant) {
	switch kst := fn.proto.Consts[i].(type) {
		case float64:
			return Float(kst)
		case string:
			return String(kst)
		case int64:
			return Int(kst)
		case bool:
			if kst {
				return True
			}
			return False
	}
	return c
}

// rk returns the i'th stack value or the i'th constant if
// 'i' is a constant index.
func (fn *Func) rk(i int) Value {
	if code.IsKst(i) {
		return fn.kst(code.ToKst(i))
	}
	return fn.stack[i]
}

func (fn *Func) close(level int) {
	fmt.Println("close!")
}

func (fn *Func) open(stack []Value, encup ...*upvar) {
	cls := closure{fn: fn, up: make([]*upvar, len(fn.proto.UpVars))}
	fn.closure = cls
	for i, up := range fn.proto.UpVars {
		if up.Stack {
			// upvalue refers to local variable
			// cls.up[i] = stack[up.Index]
			panic("open: up in stack!")
		} else {
			// upvalue is in enclosing function
			cls.up[i] = encup[up.Index]
		}
	}
}

func (fn *Func) checkstack(top, n int) {
	if room := len(fn.stack) - top; room < n {
		space := make([]Value, n - room)
		fn.stack = append(fn.stack, space...)
	}
}