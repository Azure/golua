//func (x *Table) Next(key Value) (Value, Value) {}
//func (x *Table) Length() int {}
//func (x *Table) Append(value Value) {}
//func (x *Table) Insert(index int, value Value) {}
//func (x *Table) Remove(index int) (value Value) {}
//func (x *Table) RawGet(key Value) (value Value) {}
//func (x *Table) RawGetInt(key int) (value Value) {}
//func (x *Table) RawGetKey(key Value) (value Value) {}
//func (x *Table) RawGetStr(key string) (value Value) {}
//func (x *Table) RawSet(key, value Value) {}
//func (x *Table) RawSetKey(key, value Value) {}
//func (x *Table) RawSetInt(key int, value Value) {}
//func (x *Table) RawSetStr(key string, value Value) {}
package lua

import (
	"fmt"
	"math"
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
	meta *table

	// iterator state
	iter []Value
	keys map[Value]int
}

func (x *table) String() string { return fmt.Sprintf("table: %p", x) }
func (x *table) Type() Type     { return TableType }

// newtable returns a new table initialized using the provided sizes
// arrayN and hashN to create the underlying hash and array part.
func newTable(state *State, arrayN, hashN int) *table {
	t := table{state: state}
	if arrayN > 0 {
		t.list = make([]Value, arrayN)
	}
	if hashN > 0 {
		t.hash = make(map[Value]Value, hashN)
	} else {
		t.hash = make(map[Value]Value)
	}
	return &t
}

func (t *table) set(k, v Value) {
	if IsNone(k) {
		return
	}
	if n, ok := k.(Number); ok {
		i := arrayIndex(n) - 1
		if i >= 0 && i < len(t.list) {
			t.list[i] = v
			return
		}
		if i == len(t.list) {
			t.list = append(t.list, v)
			return
		}
		// TODO: resize & rehash
	}
	if IsNone(v) {
		delete(t.hash, k)
		return
	}
	t.hash[k] = v
}

func (t *table) get(k Value) Value {
	if IsNone(k) {
		return None
	}
	if n, ok := k.(Number); ok {
		i := arrayIndex(n) - 1
		// fmt.Printf("table[%v] (%T) @ %d\n", k, k, i)
		if i >= 0 && i < len(t.list) {
			return t.list[i]
		}
	}
	if v, ok := t.hash[k]; ok {
		return v
	}
	return None
}

func (t *table) getStr(key string) Value {
	return t.get(String(key))
}

func (t *table) getInt(key int64) Value {
	return t.get(Int(key))
}

func (t *table) setStr(key string, value Value) {
	t.set(String(key), value)
}

func (t *table) setInt(key int64, value Value) {
	t.set(Int(key), value)
}

func (t *table) exists(key Value) bool {
	return !IsNone(t.get(key))
}

func (t *table) length() int {
	return len(t.list)
}

func (t *table) next(key Value) (k, v Value, more bool) {
	if IsNone(key) || t.keys == nil { // first iteration?
		t.keys = make(map[Value]int, len(t.hash))
		t.iter = make([]Value, 0, len(t.hash))
		for k := range t.hash {
			t.iter = append(t.iter, k)
			t.keys[k] = len(t.iter) - 1
		}
	}
	if index := t.iterKey(key); index < len(t.list) {
		for index++; index <= len(t.list); index++ {
			if !IsNone(t.list[index-1]) {
				k = Int(index)
				v = t.list[index-1]
				return k, v, true
			}
		}
	} else {
		if index = index - len(t.list); index < len(t.iter) {
			k := t.iter[index]
			v := t.hash[k]
			return k, v, true
		}
	}

	// Key did not exist or iteration ended.
	t.iter = nil
	t.keys = nil

	return None, None, false
}

// iterKey returns the index of a 'key' for table traversals. First goes
// all elements in the array part, then elements in the hash part. The
// beginning of a traversal is signaled by 0.
func (t *table) iterKey(key Value) (index int) {
	if IsNone(key) {
		return 0
	} // first iteration?
	index = arrayIndex(key)
	if index != 0 && index <= len(t.list) { // key in array?
		return index // found index
	}
	// otherwise key is in hash part.
	var found bool
	if index, found = t.keys[key]; !found {
		panic(runtimeErr(fmt.Errorf("invalid key (%v) to 'next'", key)))
	}
	// hash elements are numbered after array ones.
	return index + 1 + len(t.list)
}

const maxInt = int(^uint(0) >> 1)

func arrayIndex(val Value) int {
	switch val := val.(type) {
	case Float:
		if x, ok := float2int(float64(val)); ok {
			return x
		}
	case Int:
		if x := int(val); x > 0 && x < maxInt {
			return x
		}
	}
	return 0
}

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
