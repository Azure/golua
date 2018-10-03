package lua

import (
	//"strings"
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
		ident string // name if debugging enabled.
		local bool   // true if in frame local's stack.
		index int    // index into stack or enclosing function.
		value Value  // if closed.
	}
)

func newLuaClosure(proto *binary.Prototype) *Closure {
	cls := &Closure{binary: proto}
	if nups := len(proto.UpValues); nups > 0 {
		cls.upvals = make([]*upValue, nups, nups)
	}
	return cls
}

func newGoClosure(native Func, nups int) *Closure {
	cls := &Closure{native: native}
	if nups > 0 {
		cls.upvals = make([]*upValue, nups, nups)
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
	// if x == nil {
	// 	return "closure(none)"
	// }
	// var b strings.Builder
	// if x.isLua() {
	// 	fmt.Fprintf(&b, "closure(lua:func@%s:%d)", x.binary.Source, x.binary.SrcPos)
	// } else {
	// 	fmt.Fprintf(&b, "closure(go:%s)", x.native.String())
	// }
	// return b.String()
	return fmt.Sprintf("function: %p", x)
}

func (x *Closure) isLua() bool { return x.binary != nil }

func (cls *Closure) numParams() int {
	if !cls.isLua() {
		return 0
	}
	return cls.binary.NumParams()
}

func (cls *Closure) getUp(index int) *upValue {
	if index < len(cls.upvals) {
		return cls.upvals[index]
	}
	return nil
}

func (cls *Closure) setUp(index int, value Value) {
	if index < len(cls.upvals) {
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