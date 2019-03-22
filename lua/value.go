package lua

import (
	"fmt"
	"github.com/Azure/golua/lua/code"
)

type (
	Callable interface {
		callable
	}

	HasMeta interface {
		SetMeta(*Table)
		Meta() *Table
	}

	Value interface {
		String() string
		Type(*Thread) Type
		kind() code.Type
	}
)

const (
	NilType     = code.NilType    
	BoolType    = code.BoolType
	NumberType  = code.NumberType 
	StringType  = code.StringType 
	TableType   = code.TableType  
	FuncType    = code.FuncType   
	GoType      = code.GoType     
	ThreadType  = code.ThreadType
	IntType     = code.IntType
	FloatType   = code.FloatType
)

type GoValue struct {
	Value interface{}
	funcs *Table
}

func (v *GoValue) String() string { return fmt.Sprintf("userdata: %p", v) }

func (v *closure) Type(t *Thread) Type { return t.ls.typeOf(v) }
func (v *GoValue) Type(t *Thread) Type { return t.ls.typeOf(v) }
func (v *Thread) Type(t *Thread) Type { return t.ls.typeOf(v) }
func (v *Table) Type(t *Thread) Type { return t.ls.typeOf(v) }
func (v String) Type(t *Thread) Type { return t.ls.typeOf(v) }
func (v Float) Type(t *Thread) Type { return t.ls.typeOf(v) }
func (v Int) Type(t *Thread) Type { return t.ls.typeOf(v) }
func (v Bool) Type(t *Thread) Type { return t.ls.typeOf(v) }

func (*closure) kind() code.Type { return FuncType }
func (*GoValue) kind() code.Type { return GoType }
func (*Thread) kind() code.Type { return ThreadType }
func (*Table) kind() code.Type { return TableType }
func (String) kind() code.Type { return StringType }
func (Float) kind() code.Type { return FloatType }
func (Int) kind() code.Type {  return IntType }
func (Bool) kind() code.Type { return BoolType }