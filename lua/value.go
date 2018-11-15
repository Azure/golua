package lua

import (
	"fmt"
	"math/big"
	"reflect"
	"runtime"
)

type Type int

const (
	NoneType Type = iota
	NilType
	BoolType
	NumberType
	StringType
	FuncType
	UserDataType
	ThreadType
	TableType
	maxTypeID
)

var types = [...]string{
	NoneType:     "no value",
	NilType:      "nil",
	BoolType:     "boolean",
	NumberType:   "number",
	StringType:   "string",
	FuncType:     "function",
	UserDataType: "userdata",
	ThreadType:   "thread",
	TableType:    "table",
}

func (id Type) String() string { return types[id] }

type (
	Number interface {
		Value
		number()
	}

	Value interface {
		// String returns the canonical string of this value.
		String() string

		// Type returns the type name of this value.
		//
		// Types:
		// nil
		// boolean
		// number
		// string
		// function
		// userdata
		// thread
		// table
		Type() Type
	}
)

// func (x *Object) Format(fmts fmt.State, char rune) {}
// func (x *Table) Format(fmts fmt.State, char rune) {}
// func (x Float) Format(fmts fmt.State, char rune) {}
// func (x String) Format(fmts fmt.State, char rune) {}
// func (x Bool) Format(fmts fmt.State, char rune) {}
// func (x Int) Format(fmts fmt.State, char rune) {}
// func (x Nil) Format(fmts fmt.State, char rune) {}
// func (x Func) Format(fmts fmt.State, char rune) {}
// func (x *Thread) Format(fmts fmt.State, char rune) {}

type Object struct {
	meta *table
	data interface{}
}

func UserData(data interface{}) *Object {
	return &Object{data: data}
}

func (x *Object) Unwrap() interface{} { return x.data }
func (x *Object) String() string      { return fmt.Sprintf("userdata: %p", x) }
func (x *Object) Type() Type          { return UserDataType }

type Float float64

func (x Float) String() string { return fmt.Sprintf("%.14g", float64(x)) }
func (x Float) Type() Type     { return NumberType }
func (Float) number()          {}

type String string

func (x String) String() string { return string(x) }
func (x String) Type() Type     { return StringType }

type Bool bool

const (
	True  = Bool(true)
	False = Bool(false)
)

func (x Bool) String() string { return fmt.Sprintf("%t", bool(x)) }
func (x Bool) Type() Type     { return BoolType }

type Int int64

func (x Int) String() string { return fmt.Sprintf("%v", int64(x)) }
func (x Int) Type() Type     { return NumberType }
func (Int) number()          {}

type Nil byte

const None = Nil(0)
const nilValue = Nil(1)

func (x Nil) String() string { return "nil" }
func (x Nil) Type() Type {
	if x == None {
		return NoneType
	}
	return NilType
}

type Func func(*State) int

// func (x Func) Call(state *State) int {
// 	return x(state)
// }

func (x Func) Type() Type {
	if x == nil {
		return NilType
	}
	return FuncType
}

func (x Func) String() string {
	if x != nil {
		pc := reflect.ValueOf(x).Pointer()
		fn := runtime.FuncForPC(pc)
		_, ln := fn.FileLine(pc)
		return fmt.Sprintf("func@%s:%d", fn.Name(), ln)
	}
	return fmt.Sprintf("func@%s", None)
}

//
// Value APIs
//

func ValueOf(state *State, value interface{}) Value {
	return valueOf(state, value)
}

func valueOf(state *State, value interface{}) Value {
	switch value := value.(type) {
	case func(*State) int:
		return newGoClosure(Func(value), 0)
	case Func:
		return newGoClosure(value, 0)
	case float64:
		return Float(value)
	case float32:
		return Float(float64(value))
	case string:
		return String(value)
	case int64:
		return Int(value)
	case int32:
		return Int(int64(value))
	case int:
		return Int(int64(value))
	case bool:
		return Bool(value)
	case Value:
		return value
	case nil:
		return Nil(1)
	}
	udata := &Object{data: value}
	udata.meta = metaOf(state, udata)
	return udata
}

func IsNumber(value Value) bool { return IsFloat(value) || IsInt(value) }
func IsFloat(value Value) bool  { _, ok := value.(Float); return ok }
func IsInt(value Value) bool    { _, ok := value.(Int); return ok }
func IsNone(value Value) bool {
	return value == nil || value == None || !isNil(value) || value.Type() == NilType
}
func isNil(value Value) bool { return reflect.ValueOf(value).IsValid() } //IsNil() }

func Truth(value Value) bool { return bool(truth(value)) }

func toString(value Value) (string, bool) {
	switch value := value.(type) {
	case String:
		return string(value), true
	case Float:
		s := fmt.Sprintf("%v", float64(value))
		return s, true
	case Int:
		s := fmt.Sprintf("%v", int64(value))
		return s, true
	case Bool:
		s := fmt.Sprintf("%t", bool(value))
		return s, true
	case Nil:
		return value.String(), true
	}
	return "", false
}

func (x Float) rational() *big.Rat { return new(big.Rat).SetFloat64(float64(x)) }
func (x Int) rational() *big.Rat   { return new(big.Rat).SetInt64(int64(x)) }

// truth returns true for any Lua value different from false
// and None (or nil), otherwise returns false.
func truth(value Value) Bool {
	if b, ok := value.(Bool); ok {
		return b
	}
	return Bool(!IsNone(value))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
