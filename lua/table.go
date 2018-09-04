package lua

import (
	"strings"
	"math"
	"fmt"
)

// Implementation of tables (aka arrays, objects, or hash tables). Tables keep
// its elements in two parts: an array part and a hash part. Non-negative integer
// keys are all candidates to be kept in the array part. The actual size of the
// array is the largest 'n' such that more than half the slots between 1 and n
// are in use. Hash uses a mix of chained scatter table with Brent's variation.
// A main invariant of these tables is that, if an element is not in its main
// position (i.e. the 'original' position that its hash gives to it), then the
// colliding element is in its own main position. Hence even when the load factor
// reaches 100%, performance remains good.
type table struct {
	// thread state
	state *State

	// table state
	hash map[Value]Value
	list []Value
	meta *Table
	size int

	// iterator state
	iter []Value
	keys map[Value]int
}

// newtable returns a new table initialized using the provided sizes
// arrayN and hashN to create the underlying hash and array part.
func newTable(state *State, arrayN, hashN int) table {
	t := table{state: state}
	if arrayN > 0 {
		t.list = make([]Value, arrayN)
	}
	if hashN > 0 {
		t.hash = make(map[Value]Value, hashN)
	} else {
		t.hash = make(map[Value]Value)
	}
	return t
}

func (t *table) Set(k, v Value) {
	if isNilOrNone(k) {
		return
	}
	if n, ok := k.(Number); ok {
		i := arrayIndex(n) - 1
		if i >= 0 && i < len(t.list) {
			t.list[i] = v
			return
		}
		// TODO: resize & rehash
	}
	if isNilOrNone(v) {
		delete(t.hash, k)
		return
	}
	t.hash[k] = v
}

func (t *table) Get(k Value) Value {
	if isNilOrNone(k) {
		return None
	}
	if n, ok := k.(Number); ok {
		i := arrayIndex(n) - 1
		if i >= 0 && i < len(t.list) {
			return t.list[i]
		}
	}
	if v, ok := t.hash[k]; ok {
		return v
	}
	return None
}

func (t *table) String() string {
	var b strings.Builder
	for i, v := range t.list {
		fmt.Fprintf(&b, "<%d,%v>", i, v)
	}
	for k, v := range t.hash {
		fmt.Fprintf(&b, "<%v,%v>", k, v)
	}
	return b.String()
}

func (t *table) Next(value Value) (Value, Value) {
	fmt.Println("Next: TODO")
	return None, None
}

func (t *table) getStr(key string) Value {
	return t.Get(String(key)) 
}

func (t *table) getInt(key int64) Value {
	return t.Get(Int(key))
}

func (t *table) setStr(key string, value Value) {
	t.Set(String(key), value)
}

func (t *table) setInt(key int64, value Value) {
	t.Set(Int(key), value)
}

func (t *table) exists(key Value) bool {
	return !isNilOrNone(t.Get(key))
}

const maxInt = int(^uint(0) >> 1)

func float2int(f64 float64) (int, bool) {
	if math.IsInf(f64, 0) || math.IsNaN(f64) {
		return 0, false
	} else {
		if i64 := int64(f64); float64(i64) == f64 {
			return int(i64), true
		}
		return 0, false
	}
}

func arrayIndex(num Number) int {
	switch num := num.(type) {
		case Float:
			if x, ok := float2int(float64(num)); ok {
				return x
			}
		case Int:
			if x := int(num); x > 0 && x < maxInt {
				return x
			}
	}
	return 0
}