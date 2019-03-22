package lua

import "fmt"

var _ = fmt.Println

type callstatus int

const (
	allowhooks callstatus = 1 << iota
	luacall
	hooked
	fresh
	ypcall
	tailcall
	hookyield
	lt4le
	finalizer
	mainfunc
)

type (
	frame struct {
		prev *frame
		open *upvar
		call *call
	}

	call struct {
		flag callstatus
		fn   callable
		va   []Value
		fr   *frame
		sp   int
		sb   int
		pc   int
		want int
	}
)

func (ci *call) prepare(ls *thread, fn callable, args *[]Value, want int) *call {
	if fn, ok := fn.(*Func); ok {
		if ci.flag = luacall; fn.proto.Vararg {
			ci.varargs(args, fn.proto.ParamN)
		} else {
			extra := fn.proto.ParamN - len(*args)
			if extra < 0 {
				*args = (*args)[:len(*args)+extra]
			}
			if extra > 0 {
				*args = append(*args, make([]Value, extra)...)
			}
		}
		fn.stack = make([]Value, fn.proto.StackN)
	}
	ci.fr, ls.fr = ls.fr, &frame{prev: ls.fr, call: ci} 
	ci.want = want
	ci.fn = fn
	return ci
}

func (ci *call) varargs(args *[]Value, fixed int) {
	var param int
	for param < fixed && param < len(*args) {
		param++
	}
	for param < fixed {
		param++	
	}
	extra := (*args)[param:]
	*args = (*args)[:param]

	ci.va = make([]Value, len(extra))
	copy(ci.va, extra)
}

func (ci *call) debug(want string) *debug {
	dbg := &debug{ci: ci}
	//ci.sourceInfo(dbg)
	// ci.funcInfo(dbg)
	// ci.funcName(dbg)
	return dbg
}

func (ci *call) call(ls *thread, args []Value) ([]Value, error) {
	rets, err := ci.fn.call(ls, args)
	if err != nil {
		return nil, err
	}
	if ci.want >= 0 {
		switch extra := len(rets) - ci.want; {
			case ci.want == 0:
				rets = rets[:0]
			case extra < 0:
				nils := make([]Value, -extra)
				rets = append(rets, nils...)
			case extra > 0:
				rets = rets[:len(rets)-extra]
		}
	}
	ls.fr = ci.fr
	// ci.fr.leave(ls)
	// fmt.Printf("call: rets = %v (want = %d)\n", rets, ci.want)
	return rets, nil
}