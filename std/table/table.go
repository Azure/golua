package table

import (
	"strings"
	"fmt"
	"github.com/Azure/golua/lua"
)

//
// Lua Standard Library -- table
//

// Open opens the Lua standard Table library. 
//
// This library provides generic functions for table manipulation.
// It provides all its functions inside the table table.
//
// Remember that, whenever an operation needs the length of a table, all
// caveats about the length operator apply (see §3.4.7). All functions
// ignore non-numeric keys in the tables given as arguments.
//
// See https://www.lua.org/manual/5.3/manual.html#6.6
func Open(state *lua.State) int {
	// Create 'table' table.
	var tableFuncs = map[string]lua.Func{
 		"concat": lua.Func(tableConcat),
		"insert": lua.Func(tableInsert),
		"pack":   lua.Func(tablePack),
		"unpack": lua.Func(tableUnpack),
		"remove": lua.Func(tableRemove),
		"move":   lua.Func(tableMove),
		"sort":   lua.Func(tableSort),
	}
	state.NewTableSize(0, 7)
	state.SetFuncs(tableFuncs, 0)

	// Return 'table' table.
	return 1
}

// table.concat (list [, sep [, i [, j]]])
//
// Given a list where all elements are strings or numbers, returns the
// string list[i]..sep..list[i+1] ··· sep..list[j]. The default value
// for sep is the empty string, the default for i is 1, and the default
// for j is #list. If i is greater than j, returns the empty string.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-table.concat
func tableConcat(state *lua.State) int {
	len := length(state, 1, opRead)
	sep := state.OptString(2, "")
	i   := state.OptInt(3, 1)
	j   := state.OptInt(4, len)

	if i > j {
		state.Push("")
		return 1
	}
	buf := make([]string, j-i+1)
	for k := i; k > 0 && k <= j; k++ {
		state.GetI(1, k)
		if !state.IsString(-1) {
			panic(fmt.Errorf("invalid value (%s) at index %d in table for 'concat'",
				state.TypeAt(-1).String(),
				i,
			))
		}
		buf[k-i] = state.ToString(-1)
		state.Pop()
	}
	state.Push(strings.Join(buf, sep))
	return 1
}

// table.insert (list, [pos,] value)
//
// Inserts element value at position pos in list, shifting up the elements
// list[pos], list[pos+1], ···, list[#list]. The default value for pos is
// #list+1, so that a call table.insert(t,x) inserts x at the end of list t.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-table.insert
func tableInsert(state *lua.State) int {
	var (
		len = length(state, 1, opReadWrite) + 1 // first empty element
		pos int64 // where to insert new element
	)
	switch state.Top() {
		case 3:
			if pos = state.CheckInt(2); pos < 1 || pos > len {
				panic(fmt.Errorf("bad argument #2 to 'insert' (position out of bounds)"))
			}
			for i := len; i > pos; i-- { // move up elements
				state.GetI(1, i-1)
				state.SetI(1, i) // t[i] = t[i-1]
			}
		case 2: // called with 2 arguments
			pos = len // insert new element at the end	
		default:
			panic(fmt.Errorf("wrong number of arguments to 'insert'"))
	}
	state.SetI(1, pos) // t[pos] = v
	return 0
}

// table.pack (···)
//
// Returns a new table with all arguments stored into keys 1, 2, etc. and with
// a field "n" with the total number of arguments. Note that the resulting table
// may not be a sequence.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-table.pack
func tablePack(state *lua.State) int {
	elemc := state.Top()
	state.NewTableSize(elemc, 1)
	state.Insert(1)
	for i := int64(elemc); i >= 1; i-- {
		state.SetI(1, i)
	}
	state.Push(elemc)
	state.SetField(1, "n")
	return 1
}

// table.unpack (list [, i [, j]])
//
// Returns the elements from the given list.
//
// This function is equivalent to
//
//     return list[i], list[i+1], ···, list[j]
//
// By default, i is 1 and j is #list.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-table.unpack
func tableUnpack(state *lua.State) int {
	state.CheckType(1, lua.TableType)
	var (
		i = state.OptInt(2, 1)
		j = state.OptInt(3, int64(state.RawLen(1)))
		n = int(j - i + 1)
	)
	const max = 1000000
	if n <= 0 || n >= max || !state.CheckStack(n) {
		panic(fmt.Errorf("too many results to unpack"))
	}
	for i < j {
		state.GetI(1, i)
		i++
	}
	state.GetI(1, j)
	return n
}

// table.remove (list [, pos])
//
// Removes from list the element at position pos, returning the value of the
// removed element. When pos is an integer between 1 and #list, it shifts down
// the elements list[pos+1], list[pos+2], ···, list[#list] and erases element
// list[#list]; The index pos can also be 0 when #list is 0, or #list + 1; in
// those cases, the function erases the element list[pos].
//
// The default value for pos is #list, so that a call table.remove(l) removes
// the last element of list l.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-table.remove
func tableRemove(state *lua.State) int {
	// var (
	// 	len = length(state, 1, opReadWrite)
	// 	pos = state.OptInt(2, len)
	// )
	// if (pos != len) && (pos < 1 || pos >= len + 1) { // validate pos if given
	// 	panic(fmt.Errorf("bad argument #2 to 'remove' (position out of bounds)"))	
	// }
	// state.GetI(1, pos) // result = t[pos]
	// for ; pos < len; pos++ {
	// 	state.GetI(1, pos + 1)
	// 	state.SetI(1, pos) // t[pos] = t[pos+1]
	// }
	// state.Push(nil)
	// state.SetI(1, pos) // t[pos] = nil
	// return 1
	fmt.Println("table.remove")
	state.Debug(true)
	return 0
}

// table.move (a1, f, e, t [,a2])
//
// Moves elements from table a1 to table a2, performing the equivalent to the
// following multiple assignment: a2[t],··· = a1[f],···,a1[e]. The default for
// a2 is a1. The destination range can overlap with the source range. The number
// of elements to be moved must fit in a Lua integer.
//
// Returns the destination table a2.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-table.move
func tableMove(state *lua.State) int {
	fmt.Println("table.move")
	state.Debug(true)
	return 0
}

// table.sort (list [, comp])
//
// Sorts list elements in a given order, in-place, from list[1] to list[#list].
// If comp is given, then it must be a function that receives two list elements
// and returns true when the first element must come before the second in the
// final order (so that, after the sort, i < j implies not comp(list[j],list[i])).
// If comp is not given, then the standard Lua operator < is used instead.
//
// Note that the comp function must define a strict partial order over the elements
// in the list; that is, it must be asymmetric and transitive. Otherwise, no valid
// sort may be possible.
//
// The sort algorithm is not stable: elements considered equal by the given order
// may have their relative positions changed by the sort.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-table.sort
func tableSort(state *lua.State) int {
	fmt.Println("table.sort")
	state.Debug(true)
	return 0
}

// operations that an object must define to mimic a table (some functions
// only need some of them.)
const (
	opRead      = 1
	opWrite     = 2
	opLength    = 4
	opReadWrite = opRead|opWrite
)

// checkTable checks that 'arg' is either a table or can behave like one (that is,
// it has a metatable with the required metamethods.)
func checkTable(state *lua.State, index, ops int) {
	if state.TypeAt(index) != lua.TableType { // not a table?
		n := 1 // number of elements to pop
		if state.GetMetaTableAt(index) { // must have metatable
			if !((ops&opRead!=0) || checkField(state, "__index", n)) {
				n++	
			}
			if !((ops&opWrite!=0) || checkField(state, "__newindex", n)) {
				n++	
			}
			if !((ops&opRead!=0) || checkField(state, "__len", n)) {
				n++	
			}
			state.PopN(n) // pop metatable and tested metamethods
		}
	} else {
		state.CheckType(index, lua.TableType) // force an error.
	}
}

func checkField(state *lua.State, key string, index int) bool {
	state.Push(key)
	return state.RawGet(-index) != lua.NilType
}

func length(state *lua.State, index, ops int) int64 {
	checkTable(state, index, ops)
	return int64(state.RawLen(index))
}