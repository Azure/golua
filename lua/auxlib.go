package lua

import (
	"fmt"

	"github.com/Azure/golua/lua/code"
)

var _ = fmt.Println

type Tuple []Value

//
// Accessors
//

func (args *Tuple) Thread(arg int) (*Thread, error) {
	if ls := ToThread(args.Arg(arg)); ls != nil {
		return ls, nil
	}
	typ := TypeName(args.Arg(arg))
	return nil, TypeErr(arg, typ, "thread")
}

func (args *Tuple) Table(arg int) (*Table, error) {
	if tbl := ToTable(args.Arg(arg)); tbl != nil {
		return tbl, nil
	}
	typ := TypeName(args.Arg(arg))
	return nil, TypeErr(arg, typ, "table")
}

func (args *Tuple) GoFunc(arg int) (*GoFunc, error) {
	if fn := ToGoFunc(args.Arg(arg)); fn != nil {
		return fn, nil
	}
	typ := TypeName(args.Arg(arg))
	return nil, TypeErr(arg, typ, "function")
}

func (args *Tuple) Number(arg int) (Number, error) {
	if num := ToNumber(args.Arg(arg)); num != nil {
		return num, nil
	}
	typ := TypeName(args.Arg(arg))
	return nil, TypeErr(arg, typ, "number")
}

func (args *Tuple) String(arg int) (String, error) {
	if str, ok := ToString(args.Arg(arg)); ok {
		return str, nil
	}
	typ := TypeName(args.Arg(arg))
	return "", TypeErr(arg, typ, "string")
}

func (args *Tuple) Float(arg int) (Float, error) {
	if f64, ok := ToFloat(args.Arg(arg)); ok {
		return f64, nil
	}
	typ := TypeName(args.Arg(arg))
	return 0, TypeErr(arg, typ, "float")
}

func (args *Tuple) Int(arg int) (Int, error) {
	if i64, ok := ToInt(args.Arg(arg)); ok {
		return i64, nil
	}
	typ := TypeName(args.Arg(arg))
	return 0, TypeErr(arg, typ, "integer")
}

func (args *Tuple) Bool(arg int) (Bool, error) {
	if !IsBool(args.Arg(arg)) {
		typ := TypeName(args.Arg(arg))
		return false, TypeErr(arg, typ, "boolean")
	}
	return Bool(Truth(args.Arg(arg))), nil
}

//
// Accessors with options
//

func (args *Tuple) StringOpt(arg int, opt String) String {
	if v, err := args.String(arg); err == nil {
		return v
	}
	return opt
}

func (args *Tuple) NumberOpt(arg int, opt Number) Number {
	if v, err := args.Number(arg); err == nil {
		return v
	}
	return opt
}

func (args *Tuple) FloatOpt(arg int, opt Float) Float {
	if v, err := args.Float(arg); err == nil {
		return v
	}
	return opt
}

func (args *Tuple) IntOpt(arg int, opt Int) Int {
	if v, err := args.Int(arg); err == nil {
		return v
	}
	return opt
}

func (args *Tuple) BoolOpt(arg int, opt Bool) Bool {
	if v, err := args.Bool(arg); err == nil {
		return v
	}
	return opt
}

//
// Accessor helpers
//

func (args *Tuple) Arg(arg int) Value {
	if arg < 0 || arg >= len(*args) {
		return nil
	}
	return (*args)[arg]
}

func (args *Tuple) Check(ls *Thread, check argsCheck) error {
	if err := check(ls, *args); err != nil {
		return err
	}
	return nil
}

//
// Argument Checks
//

type argsCheck func(*Thread, Tuple) error

func ValidArgs(kinds ...code.Type) argsCheck {
	return func(ls *Thread, args Tuple) error {
		if err := MinArgs(len(kinds))(ls, args); err != nil {
			return err
		}
		for arg, val := range args {
			if val == nil {
				return ArgErr(arg, fmt.Errorf("value expected"))
			}
			if arg < len(kinds) {
				if val.kind() != kinds[arg] {
					return TypeErr(arg, val.kind().String(), kinds[arg].String())
				}
			}
		}
		return nil
	}
}

func RangeArgs(min, max int) argsCheck {
	return func(ls *Thread, args Tuple) error {
		if len(args) < min || len(args) > max {
			return ArgErr(-1, fmt.Errorf("expects between %d and %d arguments", min, max))
		}
		return nil
	}
}

func ExactArgs(want int) argsCheck {
	return func(ls *Thread, args Tuple) error {
		if len(args) != want {
			return ArgErr(-1, fmt.Errorf("expects exactly %d arguments", want))
		}
		return nil
	}
}

func MaxArgs(max int) argsCheck {
	return func(ls *Thread, args Tuple) error {
		if len(args) > max {
			return ArgErr(-1, fmt.Errorf("expects at most %d arguments", max))
		}
		return nil
	}
}

func MinArgs(min int) argsCheck {
	return func(ls *Thread, args Tuple) error {
		if len(args) < min {
			return ArgErr(-1, fmt.Errorf("expects at least %d arguments", min))
		}
		return nil
	}
}

func WithArgs(check argsCheck, fn *GoFunc) *GoFunc {
	fn.args = check
	return fn
}

//
// Helpers
//

func TypeName(v Value) string {
	if v == nil {
		return "nil"
	}
	return v.kind().String()
}
