package lua

type (
	// __le: the less equal (<=) operation. Unlike other operations, the less-equal
	// operation can use two different events.
	//
	// First, Lua looks for the __le metamethod in both operands, like in the less
	// than operation. If it cannot find such a metamethod, then it will try the __lt
	// metamethod, assuming that a <= b is equivalent to not (b < a). As with the other
	// comparison operators, the result is always a boolean.
	// 
	// This use of the __lt event can be removed in future versions; it is also slower
	// than a real __le metamethod.
	HasLessEqual interface {
		Value		

		LessEqual(Value) (bool, error)
	}

	// __lt: the less than (<) operation. Behavior similar to the addition operation, except that Lua
	// will try a metamethod only when the values being compared are neither both numbers nor both strings.
	//
	// The result of the call is always converted to a boolean.
	HasLessThan interface {
		Value

		LessThan(Value) (bool, error)
	}

	// __newindex: The indexing assignment table[key] = value. Like the index event, this event happens
	// when table is not a table or when key is not present in table. The metamethod is looked up in table.
	// Like with indexing, the metamethod for this event can be either a function or a table.
	//
	// If it is a function, it is called with table, key, and value as arguments. If it is a table, Lua
	// does an indexing assignment to this table with the same key and value. (This assignment is regular,
	// not raw, and therefore can trigger another metamethod.)
	//
	// Whenever there is a __newindex metamethod, Lua does not perform the primitive assignment.
	// If necessary, the metamethod itself can call rawset to do the assignment.
	HasSetIndex interface {
		Value

		Set(index, value Value) error
	}

	// __index: The indexing access operation table[key]. This event happens when table is not a table
	// or when key is not present in table. The metamethod is looked up in table. Despite the name, the
	// metamethod for this event can be either a function or a table. If it is a function, it is called
	// with table and key as arguments, and the result of the call (adjusted to one value) is the result
	// of the operation. If it is a table, the final result is the result of indexing this table with key.
	// (This indexing is regular, not raw, and therefore can trigger another metamethod.)
	HasGetIndex interface {
		Value

		Get(index Value) (Value, error)
	}

	// __len: the length (#) operation. If the object is not a string, Lua will try its metamethod.
	// If there is a metamethod, Lua calls it with the object as argument, and the result of the call
	// (always adjusted to one value) is the result of the operation. If there is no metamethod but
	// the object is a table, then Lua uses the table length operation (see §3.4.7).
	//
	// Otherwise, Lua raises an error.
	HasLength interface {
		Value

		Length() (int, error)
	}

	// __concat: the concatenation (..) operation. Behavior similar to the addition operation, except
	// that Lua will try a metamethod if any operand is neither a string nor a number which is always
	// coercible to a string.
	HasConcat interface {
		Value

		Concat(Value) (Value, error)
	}

	// __call: The call operation func(args). This event happens when Lua tries to call a non-function
	// value (that is, func is not a function). The metamethod is looked up in func. If present, the
	// metamethod is called with func as its first argument, followed by the arguments of the original
	// call (args). All results of the call are the result of the operation.
	//
	// This is the only metamethod that allows multiple results.
	Callable interface {
		Value

		Call(args ...Value) ([]Value, error)
	}

	// __eq: the equal (==) operation. Behavior similar to the addition operation, except that Lua will
	// try a metamethod only when the values being compared are either both tables or both full userdata
	// and they are not primitively equal. The result of the call is always converted to a boolean.
	HasEquals interface {
		Value

		Equals(Value) (Value, error)
	}

	// __unm: the negation (unary -) operation. Behavior similar to the addition operation.
	HasMinus interface {
		Value

		Minus(Value) (Value, error)
	}

	// __add: the addition (+) operation. If any operand for an addition is not a number (nor a string
	// coercible to a number), Lua will try to call a metamethod. First, Lua will check the first operand
	// (even if it is valid). If that operand does not define a metamethod for __add, then Lua will check
	// the second operand. If Lua can find a metamethod, it calls the metamethod with the two operands as
	// arguments, and the result of the call (adjusted to one value) is the result of the operation.
	//
	// Otherwise, it raises an error.
	HasAdd interface {
		Value

		Add(Value) (Value, error)
	}

	// __sub: the subtraction (-) operation. Behavior similar to the addition operation.
	HasSub interface {
		Value

		Sub(Value) (Value, error)
	}
	
	// __mul: the multiplication (*) operation. Behavior similar to the addition operation.
	HasMul interface {
		Value

		Mul(Value) (Value, error)
	}

	//__div: the division (/) operation. Behavior similar to the addition operation.
	HasDiv interface {
		Value

		Div(Value) (Value, error)
	}

	// __mod: the modulo (%) operation. Behavior similar to the addition operation.
	HasMod interface {
		Value

		Mod(Value) (Value, error)
	}
	
	// __pow: the exponentiation (^) operation. Behavior similar to the addition operation.
	HasPow interface {
		Value

		Pow(Value) (Value, error)
	}

	// __band: the bitwise AND (&) operation. Behavior similar to the addition operation, except
	// that Lua will try a metamethod if any operand is neither an integer nor a value coercible
	// to an integer (see §3.4.3).
	HasAnd interface {
		Value

		And(Value) (Value, error)
	}
	
	// __bxor: the bitwise exclusive OR (binary ~) operation. Behavior similar to the bitwise AND operation.
	HasXor interface {
		Value

		Xor(Value) (Value, error)
	}
	
	// __shl: the bitwise left shift (<<) operation. Behavior similar to the bitwise AND operation.
	HasShl interface {
		Value

		Lsh(Value) (Value, error)
	}

	// __shr: the bitwise right shift (>>) operation. Behavior similar to the bitwise AND operation.
	HasShr interface {
		Value

		Rsh(Value) (Value, error)
	}

	// __bnot: the bitwise NOT (unary ~) operation. Behavior similar to the bitwise AND operation.
	HasNot interface {
		Value

		Not() (Value, error)
	}
	
	// __bor: the bitwise OR (|) operation. Behavior similar to the bitwise AND operation.
	HasOr interface {
		Value

		Or(Value) (Value, error)
	}

	// __idiv: the floor division (//) operation. Behavior similar to the addition operation.
)

type metaEvent int

const (
	metaAdd metaEvent = iota + 1
	metaSub
	metaMul
	metaDiv
	metaMod
	metaPow
	metaUnm
	metaIdiv
	metaBand
	metaBor
	metaBxor
	metaBnot
	metaShl
	metaShr
	metaConcat
	metaLen
	metaEq
	metaLt
	metaLe
	metaIndex
	metaNewIndex
	metaCall
	metaMode
	metaTagN
)

var metaFields = [...]string{
	metaAdd:      "add",
	metaSub:      "sub",
	metaMul:      "mul",
	metaDiv:      "div",
	metaMod:      "mod",
	metaPow:      "pow",
	metaUnm:      "unm",
	metaIdiv:     "idiv",
	metaBand:     "band",
	metaBor:      "bor",
	metaBxor:     "bxor",
	metaBnot:     "bnot",
	metaShl:      "shl",
	metaShr:      "shr",
	metaConcat:   "concat",
	metaLen:      "len",
	metaEq:       "eq",
	metaLt: 	  "lt",
	metaLe: 	  "le",
	metaIndex: 	  "index",
	metaNewIndex: "newindex",
	metaCall: 	  "call",
	metaMode: 	  "mode",
}

func (evt metaEvent) toName() string { return "__" + metaFields[evt] }

// TODO: idiv
func metaOf(v Value) *Table {
	var events Table
	switch v := v.(type) {
		case *Object:
			var u interface{}
			if u = v.Unwrap(); u == nil {
				break
			}
			if o, ok := u.(HasSetIndex); ok { // __newindex
				events.setStr(metaNewIndex.toName(), o)
			}
			if o, ok := u.(HasGetIndex); ok { // __index
				events.setStr(metaIndex.toName(), o)
			}
			if o, ok := u.(HasLength); ok { // __len
				events.setStr(metaLen.toName(), o)
			}
			if o, ok := u.(Callable); ok { // __call
				events.setStr(metaCall.toName(), o)
			}
			if o, ok := u.(HasConcat); ok { // __concat
				events.setStr(metaConcat.toName(), o)
			}
			if o, ok := u.(HasMinus); ok { // __unm
				events.setStr(metaUnm.toName(), o)
			}
			if o, ok := u.(HasAdd); ok { // __add
				events.setStr(metaAdd.toName(), o)
			}
			if o, ok := u.(HasSub); ok { // __sub
				events.setStr(metaSub.toName(), o)
			}
			if o, ok := u.(HasMul); ok { // __mul
				events.setStr(metaMul.toName(), o)
			}
			if o, ok := u.(HasDiv); ok { // __div, __idiv
				events.setStr(metaDiv.toName(), o)
				events.setStr(metaIdiv.toName(), o)
			}
			if o, ok := u.(HasMod); ok { // __mod
				events.setStr(metaMod.toName(), o)
			}
			if o, ok := u.(HasPow); ok { // __pow
				events.setStr(metaPow.toName(), o)
			}
			if o, ok := u.(HasEquals); ok { // __eq
				events.setStr(metaEq.toName(), o)
			}
			if o, ok := u.(HasLessThan); ok { // __lt
				events.setStr(metaLt.toName(), o)
			}
			if o, ok := u.(HasLessEqual); ok { // __le
				events.setStr(metaLe.toName(), o)
			}
			if o, ok := u.(HasAnd); ok { // __band
				events.setStr(metaBand.toName(), o)
			}
			if o, ok := u.(HasOr); ok { // __bor
				events.setStr(metaBor.toName(), o)
			}
			if o, ok := u.(HasXor); ok { // __bxor
				events.setStr(metaBxor.toName(), o)
			}
			if o, ok := u.(HasNot); ok { // __bnot
				events.setStr(metaBnot.toName(), o)
			}
			if o, ok := u.(HasShl); ok { // __shl
				events.setStr(metaShl.toName(), o)
			}
			if o, ok := u.(HasShr); ok { // __shr
				events.setStr(metaShr.toName(), o)
			}
	}
	return &events
}