package lua

import (
	"runtime"
	"reflect"
	"fmt"
)

type Type int

const (
	NilType  	 Type = iota
	IntType
	BoolType
	FloatType
	NumberType
	StringType
	FuncType
	UserDataType
	ThreadType
	TableType
	maxTypeID
)

var types = [...]string{
	NilType: 	  "nil",
	IntType: 	  "int",
	BoolType: 	  "boolean",
	FloatType:    "float",
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

type Object struct {
	meta *Table
	data interface{}
}

func UserData(data interface{}) *Object {
	return &Object{data: data}
}

func (x *Object) Unwrap() interface{} { return x.data }
func (x *Object) String() string { return fmt.Sprintf("userdata: %p", x) }
func (x *Object) Type() Type { return UserDataType }

type Table struct { table }
func (x *Table) String() string { return fmt.Sprintf("table: %p", x) }
func (x *Table) Type() Type { return TableType }

//func (x *Table) Next(key Value) (Value, Value) {}
//func (x *Table) Length() int {}
//func (x *Table) Append(value Value) {}
//func (x *Table) Insert(index int, value Value) {}
//func (x *Table) Remove(index int) (value Value) {}

//func (x *Table) RawGet(key Value) (value Value) {}
//func (x *Table) RawGetInt(key int) (value Value) {}
//func (x *Table) RawGetKey(key Value) (value Value) {}
//func (x *Table) RawGetStr(key string) (value Value) {}

//func (x *Table) RawSet(key, value Value) {}
//func (x *Table) RawSetKey(key, value Value) {}
//func (x *Table) RawSetInt(key int, value Value) {}
//func (x *Table) RawSetStr(key string, value Value) {}

type Float float64
func (x Float) String() string { return fmt.Sprintf("%v", float64(x)) }
func (x Float) Type() Type { return FloatType }
func (Float) number() {}

type String string
func (x String) String() string { return string(x) }
func (x String) Type() Type { return StringType }

type Bool bool
const (
	True  = Bool(true)
	False = Bool(false)
)
func (x Bool) String() string { return fmt.Sprintf("%t", bool(x)) }
func (x Bool) Type() Type { return BoolType }

type Int int64
func (x Int) String() string { return fmt.Sprintf("%v", int64(x)) }
func (x Int) Type() Type { return IntType }
func (Int) number() {}

type Nil byte
const None = Nil(0)
func (x Nil) String() string { return "nil" }
func (x Nil) Type() Type { return NilType }

type Func func(*State)int

func (x Func) Call(state *State) int {
	return x(state)
}

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

func (x *Table) ForEach(fn func(k, v Value)) {
	for i, v := range x.table.list {
		fn(Int(i), v)
	}
	for k, v := range x.table.hash {
		fn(k, v)
	}
}

func valueOf(state *State, value interface{}) Value {
	switch value := value.(type) {
		case func(*State)int:
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
			return None
	}
	udata := &Object{data: value}
	udata.meta = metaOf(state, udata)
	return udata
}

func IsNumber(value Value) bool { return IsFloat(value) || IsInt(value) }
func IsFloat(value Value) bool { return value.Type() == FloatType }
func IsNone(value Value) bool { return value == nil || value == None || !isNil(value) }
func IsInt(value Value) bool { return value.Type() == IntType }
func isNil(value Value) bool { return reflect.ValueOf(value).IsValid() }//IsNil() }

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
	}
	return "", false
}

// truth returns true for any Lua value different from false
// and None (or nil), otherwise returns false.
func truth(value Value) Bool {
	if b, ok := value.(Bool); ok {
		return b
	}
	return Bool(!IsNone(value))
}

func max(a, b int) int { if a > b { return a }; return b }
func min(a, b int) int { if a < b { return a }; return b }