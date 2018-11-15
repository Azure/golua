package lua

import (
	"fmt"

	"github.com/Azure/golua/lua/binary"
)

type (
	// closure represents a lua or go function closure.
	Closure struct {
		binary *binary.Prototype
		native Func
		upvals []*upValue
	}

	// upValue holds external local variable state.
	upValue struct {
		frame *Frame // frame upValue was opened within.
		index int    // index into stack or enclosing function.
		value Value  // if closed.
	}
)

func newLuaClosure(proto *binary.Prototype) *Closure {
	cls := &Closure{binary: proto}
	if nups := len(proto.UpValues); nups > 0 {
		cls.upvals = make([]*upValue, nups)
	}
	return cls
}

func newGoClosure(native Func, nups int) *Closure {
	cls := &Closure{native: native}
	if nups > 0 {
		cls.upvals = make([]*upValue, nups)
	}
	return cls
}

func (x *Closure) Type() Type {
	if x.isLua() {
		return FuncType
	}
	return x.native.Type()
}

func (x *Closure) String() string {
	return fmt.Sprintf("function: %p", x)
}

func (cls *Closure) isLua() bool { return cls != nil && cls.binary != nil }

func (cls *Closure) isGo() bool { return cls != nil && cls.native != nil }

func (cls *Closure) upvalues() []*upValue {
	if cls != nil {
		return cls.upvals
	}
	return nil
}

func (cls *Closure) getUp(index int) *upValue {
	if cls != nil && index < len(cls.upvals) {
		return cls.upvals[index]
	}
	return nil
}

func (cls *Closure) setUp(index int, value Value) {
	if cls != nil && index < len(cls.upvals) {
		cls.upvals[index].set(value)
	}
}

func (cls *Closure) upName(index int) (name string) {
	if cls.isLua() && index < len(cls.binary.UpNames) {
		name = cls.binary.UpNames[index]
	}
	return
}

func (up *upValue) set(value Value) {
	if up.open() {
		up.frame.set(up.index, value)
		return
	}
	up.value = value
}

func (up *upValue) get() Value {
	if up.open() {
		return up.frame.get(up.index)
	}
	return up.value
}

func (up *upValue) open() bool {
	return up.index != -1
}

func (up *upValue) close() {
	up.value = up.get()
	up.index = -1
}
