package lua

import "fmt"

var _ = fmt.Println

type event int

const (
	_index event = iota
	_newindex
	_gc
	_mode
	_len
	_add
	_sub
	_mul
	_mod
	_pow
	_div
	_idiv
	_band
	_bor
	_bxor
	_shl
	_shr
	_unm
	_bnot
	_eq
	_lt
	_le
	_concat
	_call
	maxEvent
)

var events = [...]string{
	_index:    "index",
	_newindex: "newindex",
	_gc:       "gc",
	_mode:     "mode",
	_len:      "len",
	_add:      "add",
	_sub:      "sub",
	_mul:      "mul",
	_mod:      "mod",
	_pow:      "pow",
	_div:      "div",
	_idiv:     "idiv",
	_band:     "band",
	_bor:      "bor",
	_bxor:     "bxor",
	_shl:      "shl",
	_shr:      "shr",
	_unm:      "unm",
	_bnot:     "bnot",
	_eq:       "eq",
	_lt:       "lt",
	_le:       "le",
	_concat:   "concat",
	_call:     "call",
}

func (evt event) String() string { return "__" + events[evt] }

func (evt event) compare(ls *thread, x, y Value) (v Value, err error) {
	if fn, ok := ls.meta(x, evt.String()).(callable); ok {
		rets, err := ls.call(fn.(Value), []Value{x, y}, 1)
		if err != nil {
			return nil, err
		}
		return rets[0], nil
	}
	if fn, ok := ls.meta(y, evt.String()).(callable); ok {
		rets, err := ls.call(fn.(Value), []Value{x, y}, 1)
		if err != nil {
			return nil, err
		}
		return rets[0], nil
	}
	return nil, fmt.Errorf("compare meta-event: todo!")
}

func (evt event) binary(ls *thread, x, y Value) (v Value, err error) {
	// if fn, ok := ls.meta(x, evt.String()).(callable); ok {
	// 	rets, err := ls.pcall(fn, stack{x, y}, 1)
	// 	if err == nil {
	// 		return rets[0]
	// 	}
	// }
	// if fn, ok := ls.meta(y, evt.String()).(callable); ok {
	// 	rets, err := ls.pcall(fn, stack{x, y}, 1)
	// 	if err == nil {
	// 		return rets[0]
	// 	}
	// }
	// ls.error(&evalErr{evt, x, y})
	// panic("unreachable")

	// if method := ls.meta(x, evt.String()); method != nil { // try 1st operand
	// 	return evt.call(ls, method, x, y, true)
	// }
	// if method := ls.meta(y, evt.String()); method != nil { // try 2nd operand
	// 	return evt.call(ls, method, x, y, true)
	// }
	// func (evt event) call(ls *State, fn, arg1, arg2 lua.Value, hasResult bool) (lua.Value, error) {
	// 	return nil, fmt.Errorf("%s.call(%v, %v, %v) (result = %t)\n", fn, arg1, arg2, hasResult)
	// }
	return nil, fmt.Errorf("binary meta-event: todo!")
}
