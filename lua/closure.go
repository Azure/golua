package lua

import (
	"github.com/Azure/golua/lua/binary"
)

// closure represents a lua or go function closure.
type Closure struct {
	proto  *binary.Prototype
	native Func
	upvals []*Value
}

func newLuaClosure(proto *binary.Prototype) *Closure {
	cls := &Closure{proto: proto}
	if nups := len(proto.UpValues); nups > 0 {
		cls.upvals = make([]*Value, nups)
	}
	return cls
}

func newGoClosure(native Func, nups int) *Closure {
	cls := &Closure{native: native}
	if nups > 0 {
		cls.upvals = make([]*Value, nups)
	}
	return cls
}

func (x *Closure) String() string { return "closure" }
func (x *Closure) Type() Type { return FuncType }