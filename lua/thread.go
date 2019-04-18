package lua

import (
	"fmt"

	"github.com/Azure/golua/lua/code"
)

type Context interface {
	GoLoader(*Thread, string, string) (Value, error)
	Searchers() *Table
	Preload() *Table
	Globals() *Table
	Values() *Table
	Loaded() *Table
	Config() *Config
}

type Thread struct {
	ls *thread
}

func (t *Thread) IsMainThread() bool { return t.ls == t.ls.rt.thread }

func (t *Thread) SetGlobal(name string, global Value) *Thread {
	if env := t.Globals(); env != nil {
		env.Set(String(name), global)
	}
	return t
}

func (t *Thread) Global(name string) Value {
	if env := t.Globals(); env != nil {
		return env.Get(String(name))
	}
	return nil
}

func (t *Thread) Globals() *Table {
	return t.Context().Globals()
}

func (t *Thread) Context() Context {
	return t.ls.rt
}

func (t *Thread) Require(lib Library, global bool) (Value, error) {
	loaded := t.ls.rt.loaded
	tbl, err := lib.Open(t)
	if err != nil {
		return nil, err
	}
	loaded.Set(String(lib.Name), tbl)
	if global {
		t.SetGlobal(lib.Name, tbl)
	}
	return tbl, nil
}

func (t *Thread) Index(tbl, key Value) (Value, error) {
	return gettable(t.ls, tbl, key)
}

func (t *Thread) TypeOf(obj Value) Type {
	return t.ls.typeOf(obj)
}

func (t *Thread) ExecN(chunk *code.Chunk, args []Value, want int) ([]Value, error) {
	return t.ls.call(t.Load(chunk), args, want)
}

func (t *Thread) Exec(chunk *code.Chunk, args ...Value) ([]Value, error) {
	return t.ls.call(t.Load(chunk), args, -1)
}

func (t *Thread) CallN(fv Value, args []Value, want int) ([]Value, error) {
	return t.ls.call(fv, args, want)
}

func (t *Thread) Call(fv Value, args ...Value) ([]Value, error) {
	return t.ls.call(fv, args, -1)
}

func (t *Thread) Load(chunk *code.Chunk) *Func {
	return t.ls.load(chunk)
}

func (t *Thread) String() string {
	return fmt.Sprintf("thread: %p", t)
}

func (t *Thread) Caller(level int) Debug {
	if ci := t.ls.caller(level); ci != nil {
		return ci.debug("Sl")
	}
	return nil
}
