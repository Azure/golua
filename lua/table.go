package lua

import (
	"fmt"
)

type border int

const (
	borderOk border = iota
	borderUp
	borderDown
)

type Table struct {
	kvs  map[Value]entry
	seqN Int
	key0 Value
	flag border
	meta *Table
}

type entry struct {
	value Value
	next  Value
}

func NewTableFromMap(kvs map[string]Value) *Table {
	t := NewTableSize(0, len(kvs))
	for k, v := range kvs {
		t.Set(String(k), v)
	}
	return t
}

func NewTableSize(arrN, kvsN int) *Table {
	return &Table{kvs: make(map[Value]entry, kvsN)}
}

func NewTable() *Table { return NewTableSize(0, 0) }

func (t *Table) String() string      { return fmt.Sprintf("table: %p", t) }
func (t *Table) SetMeta(meta *Table) { t.meta = meta }
func (t *Table) Meta() *Table        { return t.meta }

func (t *Table) Slice() (slice []Value) {
	var (
		i = Int(1)
		v = t.Get(i)
	)
	for v != nil {
		slice = append(slice, v)
		i++
		v = t.Get(i)
	}
	return slice
}

func (t *Table) Next(key Value) (k, v Value, ok bool) {
	var e entry
	if key == nil {
		k, ok = t.key0, true
	} else {
		if e, ok = t.kvs[key]; !ok {
			return
		}
		k = e.next
	}
	for k != nil {
		e = t.kvs[k]
		if v = e.value; v != nil {
			return
		}
		k = e.next
	}
	return
}

func (t *Table) Length() Int {
	switch t.flag {
	case borderDown:
		for t.seqN > 0 && t.kvs[t.seqN].value == nil {
			t.seqN--
		}
	case borderUp:
		for t.kvs[t.seqN+1].value != nil {
			t.seqN++
		}
	}
	t.flag = borderOk
	return t.seqN
}

func (t *Table) Get(k Value) Value {
	if n, ok := k.(Float); ok {
		if i := Int(n); Float(i) == n {
			k = i
		}
	}
	return t.kvs[k].value
}

func (t *Table) Set(k, v Value) {
	switch k := k.(type) {
	case Float:
		if i := Int(k); Float(i) == k {
			t.setInt(i, v)
			return
		}
		t.set(k, v)
	case Int:
		t.setInt(k, v)
	default:
		t.set(k, v)
	}
}

func (t *Table) setInt(k Int, v Value) {
	switch {
	case v == nil && t.seqN > 0 && k == t.seqN:
		t.flag = borderDown
		t.seqN--
	case k > t.seqN && v != nil:
		t.flag = borderUp
		t.seqN = k
	}
	t.set(k, v)
}

func (t *Table) set(k, v Value) {
	entry, ok := t.kvs[k]
	if v == nil && !ok {
		return
	}
	entry.value = v
	if !ok {
		entry.next = t.key0
		t.key0 = k
	}
	t.kvs[k] = entry
}

func (t *Table) foreach(fn func(k, v Value) bool) {
	for k, v, ok := t.Next(nil); k != nil && ok; {
		if !fn(k, v) {
			return
		}
		k, v, ok = t.Next(k)
	}
}

// fb2int converts a "floating point byte" back to an integer.
func fb2int(x int) int {
	e := (x >> 3) & 0x1F
	if e == 0 {
		return x
	}
	return ((x & 7) + 8) << uint(e-1)
}
