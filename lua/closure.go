package lua

import (
	"strings"
	"fmt"

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

func (x *Closure) Type() Type {
	if x.isLua() {
		return FuncType
	}
	return x.native.Type()
}

func (x *Closure) String() string {
	if x == nil {
		return "closure(none)"
	}
	var b strings.Builder
	if x.isLua() {
		fmt.Fprintf(&b, "closure(lua:func@%s:%d)", x.proto.Source, x.proto.SrcPos)
	} else {
		fmt.Fprintf(&b, "closure(go:%s)", x.native.String())
	}
	return b.String()
}

func (x *Closure) isLua() bool { return x.proto != nil } 

func (cls *Closure) openUpValues(state *State) {
	state.Log("openUpValues: TODO")
}

func (cls *Closure) closeUpValues(state *State, upto Value) {
	state.Log("closeUpValues: TODO")
}