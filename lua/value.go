package lua

import (
	"fmt"
)

type Type int

const (
	NilType  	 Type = iota
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
	NilType: 	  "nil",
	BoolType: 	  "boolean",
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

func (x *Object) Unwrap() interface{} { return x.data }
func (x *Object) String() string { return "userdata" }
func (x *Object) Type() Type { return UserDataType }

type Table struct { table }
func (x *Table) String() string { return "table" }
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
func (x Float) String() string { return "number" }
func (x Float) Type() Type { return NumberType }
func (Float) number() {}

type Func func(*State)uint32
func (x Func) String() string { return "function" }
func (x Func) Type() Type { return FuncType }

type String string
func (x String) String() string { return string(x) }
func (x String) Type() Type { return StringType }


type Bool bool
const (
	True  = Bool(true)
	False = Bool(false)
)
func (x Bool) String() string { return fmt.Sprintf("%t", x) }
func (x Bool) Type() Type { return BoolType }

type Int int64
func (x Int) String() string { return "number" }
func (x Int) Type() Type { return NumberType }
func (Int) number() {}

type Nil byte
const None = Nil(0)
func (x Nil) String() string { return "nil" }
func (x Nil) Type() Type { return NilType }

//
// Value APIs
//

func ValueOf(value interface{}) Value {
	switch value := value.(type) {
		case func(*State)uint32:
			return Func(value)
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
	udata.meta = metaOf(udata)
	return udata
}

func isNilOrNone(value Value) bool { return value == nil || value == None }