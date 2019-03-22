package lua

import (
	"github.com/Azure/golua/lua/code"
)

type Type interface {
	SetMeta(funcs *Table) *Table
	Method(name string) Callable
	Kind() code.Type
	Name() string
}

type rtype struct {
	mt *Table
	tv Value
}

func (t *rtype) SetMeta(funcs *Table) (prev *Table) {
	// if t != nil && t.mt != nil {
	// 	prev = metatypes[t.tv.kind()]
	// 	metatypes[t.tv.kind()] = funcs
	// 	return prev
	// }
	return nil
}

func (t *rtype) Method(name string) Callable {
	if t != nil && t.mt != nil {
		fn, ok := t.mt.Get(String(name)).(Callable)
		if fn != nil && ok {
			return fn
		}
	}
	return nil
}

func (t *rtype) Name() string {
	// TODO: check "__name"
	return t.Kind().String()
}

func (t *rtype) Kind() code.Type {
	if t == nil || t.tv == nil {
		return code.NilType
	}
	return t.tv.kind()
}