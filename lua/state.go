package lua

import (
	"sync"
	"fmt"
	"github.com/fibonacci1729/golua/lua/code"
)

type (
	packages struct {
		searchers *Table
		loaded    *Table
		preload   *Table
		loader    goloader
	}

	runtime struct {
		packages
		globals *Table
		config  *Config
		thread  *thread
		values  *Table
		wait    sync.WaitGroup
		types   [code.MaxType]*Table
	}
)

func (rt *runtime) init(config *Config) *thread {
	fr := &frame{call: &call{flag: mainfunc}}
	ls := &thread{rt: rt, fr: fr}

	rt.packages = packages{
		searchers: NewTable(),
		preload:   NewTable(),
		loaded:    NewTable(),
	}

	rt.globals = NewTable()
	rt.values  = NewTable()
	rt.thread  = ls
	config.init(rt)

	return ls
}

func (rt *runtime) GoLoader(ls *Thread, file, name string) (Value, error) {
	return rt.loader.load(ls, file, name)
}
// func (rt *runtime) Searcher() Searcher { return nil }
func (rt *runtime) Searchers() *Table { return rt.searchers }
func (rt *runtime) Preload() *Table { return rt.preload }
func (rt *runtime) Globals() *Table { return rt.globals }
func (rt *runtime) Values() *Table { return rt.values }
func (rt *runtime) Loaded() *Table { return rt.loaded }
func (rt *runtime) Config() *Config { return rt.config }

type thread struct {
	// co *coroutine
	rt *runtime
	tt *Thread
	fr *frame
}

func (ls *thread) call(fn Value, args []Value, want int) ([]Value, error) {
	if fn, ok := fn.(callable); ok {
		ci := new(call).prepare(ls, fn, &args, want)
		return ci.call(ls, args)
	}
	method := ls.meta(fn, "__call")
	if method == nil {
		return nil, fmt.Errorf("attempt to call a nil value")
	}
	return ls.call(method, append([]Value{fn}, args...), want)
}

func (ls *thread) load(chunk *code.Chunk) *Func {
	up := make([]*upvar, len(chunk.Main.UpVars))
	up[0] = &upvar{value: ls.rt.globals}

	fn := &Func{proto: chunk.Main}
	fn.closure = closure{fn, up}
	return fn
}

func (ls *thread) meta(value Value, event string) Value {
	if v, ok := value.(HasMeta); ok {
		if meta := v.Meta(); meta != nil {
			return meta.Get(String(event))
		}
	}
	if value != nil {
		meta := ls.rt.types[0x0F&value.kind()]
		if meta != nil {
			return meta.Get(String(event))
		}
	}
	return nil
}

func (ls *thread) typeOf(v Value) *rtype {
	if m, ok := v.(HasMeta); ok {
		return &rtype{m.Meta(), v}
	}
	if v != nil {
		meta := ls.rt.types[0x0F&v.kind()]
		return &rtype{meta, v}
	}
	return nil
}

func (ls *thread) error(err error) error {
	if err == nil {
		return nil
	}
	fmt.Printf("%T\n", err)
	return err
}

func (ls *thread) caller(skip int) *call {
	var ( fr, ci = ls.fr, ls.fr.call )

	for skip > 0 && (ci.flag & mainfunc == 0) {
		fr = fr.prev
		ci = fr.call
		skip--
	}
	if skip == 0 && (ci.flag & mainfunc == 0) {
		return ci
	}
	return nil
}