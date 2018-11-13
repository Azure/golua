package lua

import (
    "fmt"
)

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
        //Value       

        LessEqual(Value) (bool, error)
    }

    // __lt: the less than (<) operation. Behavior similar to the addition operation, except that Lua
    // will try a metamethod only when the values being compared are neither both numbers nor both strings.
    //
    // The result of the call is always converted to a boolean.
    HasLessThan interface {
        //Value

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
    HasNewIndex interface {
        //Value

        SetIndex(key, value Value) error
    }

    // __index: The indexing access operation table[key]. This event happens when table is not a table
    // or when key is not present in table. The metamethod is looked up in table. Despite the name, the
    // metamethod for this event can be either a function or a table. If it is a function, it is called
    // with table and key as arguments, and the result of the call (adjusted to one value) is the result
    // of the operation. If it is a table, the final result is the result of indexing this table with key.
    // (This indexing is regular, not raw, and therefore can trigger another metamethod.)
    HasIndex interface {
        //Value

        //Get(index Value) (Value, error)
        Index(key Value) (Value, error)
    }

    // __len: the length (#) operation. If the object is not a string, Lua will try its metamethod.
    // If there is a metamethod, Lua calls it with the object as argument, and the result of the call
    // (always adjusted to one value) is the result of the operation. If there is no metamethod but
    // the object is a table, then Lua uses the table length operation (see §3.4.7).
    //
    // Otherwise, Lua raises an error.
    HasLength interface {
        //Value

        Length() (int, error)
    }

    // __concat: the concatenation (..) operation. Behavior similar to the addition operation, except
    // that Lua will try a metamethod if any operand is neither a string nor a number which is always
    // coercible to a string.
    HasConcat interface {
        //Value

        Concat(Value) (Value, error)
    }

    // __call: The call operation func(args). This event happens when Lua tries to call a non-function
    // value (that is, func is not a function). The metamethod is looked up in func. If present, the
    // metamethod is called with func as its first argument, followed by the arguments of the original
    // call (args). All results of the call are the result of the operation.
    //
    // This is the only metamethod that allows multiple results.
    Callable interface {
        //Value

        Call(args ...Value) ([]Value, error)
    }

    // __eq: the equal (==) operation. Behavior similar to the addition operation, except that Lua will
    // try a metamethod only when the values being compared are either both tables or both full userdata
    // and they are not primitively equal. The result of the call is always converted to a boolean.
    HasEquals interface {
        //Value

        Equals(Value) (Value, error)
    }

    // __unm: the negation (unary -) operation. Behavior similar to the addition operation.
    HasMinus interface {
        //Value

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
        //Value

        Add(Value) (Value, error)
    }

    // __sub: the subtraction (-) operation. Behavior similar to the addition operation.
    HasSub interface {
        //Value

        Sub(Value) (Value, error)
    }
    
    // __mul: the multiplication (*) operation. Behavior similar to the addition operation.
    HasMul interface {
        //Value

        Mul(Value) (Value, error)
    }

    //__div: the division (/) operation. Behavior similar to the addition operation.
    HasDiv interface {
        //Value

        Div(Value) (Value, error)
    }

    // __mod: the modulo (%) operation. Behavior similar to the addition operation.
    HasMod interface {
        //Value

        Mod(Value) (Value, error)
    }
    
    // __pow: the exponentiation (^) operation. Behavior similar to the addition operation.
    HasPow interface {
        //Value

        Pow(Value) (Value, error)
    }

    // __band: the bitwise AND (&) operation. Behavior similar to the addition operation, except
    // that Lua will try a metamethod if any operand is neither an integer nor a value coercible
    // to an integer (see §3.4.3).
    HasAnd interface {
        //Value

        And(Value) (Value, error)
    }
    
    // __bxor: the bitwise exclusive OR (binary ~) operation. Behavior similar to the bitwise AND operation.
    HasXor interface {
        //Value

        Xor(Value) (Value, error)
    }
    
    // __shl: the bitwise left shift (<<) operation. Behavior similar to the bitwise AND operation.
    HasShl interface {
        //Value

        Lsh(Value) (Value, error)
    }

    // __shr: the bitwise right shift (>>) operation. Behavior similar to the bitwise AND operation.
    HasShr interface {
        //Value

        Rsh(Value) (Value, error)
    }

    // __bnot: the bitwise NOT (unary ~) operation. Behavior similar to the bitwise AND operation.
    HasNot interface {
        //Value

        Not() (Value, error)
    }
    
    // __bor: the bitwise OR (|) operation. Behavior similar to the bitwise AND operation.
    HasOr interface {
        //Value

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

var (
    name2event = map[string]metaEvent{
        "add":      metaAdd,
        "sub":      metaSub,
        "mul":      metaMul,
        "div":      metaDiv,
        "mod":      metaMod,
        "pow":      metaPow,
        "unm":      metaUnm,
        "idiv":     metaIdiv,
        "band":     metaBand,
        "bor":      metaBor,
        "bxor":     metaBxor,
        "bnot":     metaBnot,
        "shl":      metaShl,
        "shr":      metaShr,
        "concat":   metaConcat,
        "len":      metaLen,
        "eq":       metaEq,
        "lt":       metaLt,
        "le":       metaLe,
        "index":    metaIndex,
        "newindex": metaNewIndex,
        "call":     metaCall,
        "mode":     metaMode,
    }

    event2name = [...]string{
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
        metaLt:       "lt",
        metaLe:       "le",
        metaIndex:    "index",
        metaNewIndex: "newindex",
        metaCall:     "call",
        metaMode:     "mode",
    }
)

func (evt metaEvent) ID() string { return "__" + event2name[evt] }

func (evt metaEvent) name() string { return event2name[evt] }

func metaOf(state *State, v Value) *table {
    events := newTable(state, 0, 0)
    switch v := v.(type) {
        case *Object:
            var u interface{}
            if u = v.Unwrap(); u == nil {
                break
            }
            if o, ok := u.(HasNewIndex); ok { // __newindex
                method := Func(func(state *State) int {
                    var (
                        key = state.frame().pop()
                        val = state.frame().pop()
                    )
                    if err := o.SetIndex(key, val); err != nil {
                        state.errorf("%v", err)
                    }
                    return 0
                })
                events.setStr(metaNewIndex.ID(), newGoClosure(method, 0))
            }
            if o, ok := u.(HasIndex); ok { // __index
                method := Func(func(state *State) int {
                    v, err := o.Index(state.frame().pop())
                    if err != nil {
                        state.errorf("%v", err)
                    }
                    state.Push(v)
                    return 1
                })
                events.setStr(metaIndex.ID(), newGoClosure(method, 0))
            }
            if o, ok := u.(Callable); ok { // __call
                method := Func(func(state *State) int {
                    args := state.frame().popN(state.frame().gettop())
                    vs, err := o.Call(args...)
                    if err != nil {
                        state.errorf("%v", err)
                    }
                    if vs == nil || len(vs) == 0 {
                        return 0
                    }
                    for _, v := range vs {
                        state.Push(v)
                    }
                    return len(vs)
                })
                events.setStr(metaCall.ID(), newGoClosure(method, 0))
            }
            if o, ok := u.(HasConcat); ok { // __concat
                method := Func(func(state *State) int {
                    v, err := o.Concat(state.frame().pop())
                    if err != nil {
                        state.errorf("%v", err)
                    }
                    state.Push(v)
                    return 1
                })
                events.setStr(metaConcat.ID(), newGoClosure(method, 0))
            }

            if o, ok := u.(HasAdd); ok { // __add
                method := Func(func(state *State) int {
                    v, err := o.Add(state.frame().pop())
                    if err != nil {
                        state.errorf("%v", err)
                    }
                    state.Push(v)
                    return 1
                })
                events.setStr(metaAdd.ID(), newGoClosure(method, 0))
			}
    }
    return events
}

// tryMetaNewIndex performs the indexing assignment table[key] = value. Like the
// index event, this event happens when object is not a table or when key is not
// present in table. The metamethod is looked up in object.
//
// Like with indexing, the metamethod for this event can be either a function or
// a table. If it is a function, it is called with arguments object, key, and value.
// This assignment is regular, not raw, and therefore can trigger another metamethod.
//
// Whenever there is a __newindex metamethod, Lua does not perform the primitive
// assignment. If necessary, the metamethod itself can call rawset to perform the
// assignment directly.
func tryMetaNewIndex(state *State, object, key, value Value) error {
    const event = metaNewIndex

    for loop, meta, object := 0, Value(None), object; loop < metaLoopMax; loop++ {
        if meta = state.metafield(object, event.ID()); !IsNone(meta) {
            switch meta := meta.(type) {
                case *Closure:
                    state.frame().push(meta)
                    state.frame().push(object)
                    state.frame().push(key)
                    state.frame().push(value)
                    state.Call(3, 0)
                    return nil
                case *table:
                    object = meta
                    continue
                default:
                    return fmt.Errorf("meta method '__newindex' not a table or function")
            }
        }
        if table, ok := object.(*table); ok {
            table.set(key, value)
            return nil
        }
    }
    return fmt.Errorf("'__newindex' chain too long; possible loop")
}

// tryMetaIndex performs the indexing access operation table[key]. This event
// happens when object is not a table or when key is not present in table.
// The metamethod is looked up in table.
//
// Despite the name, the metamethod for this event can be either a function or
// a table. If it is a function, it is called with object and key as arguments,
// and the result of the call (adjusted to one value) is the result of the
// operation. This indexing is regular, not raw, and therefore can trigger
// another metamethod.
func tryMetaIndex(state *State, object, key Value) (Value, error) {
    const event = metaIndex

    for loop := 0; loop < metaLoopMax; loop++ {
        if table, ok := object.(*table); ok {
            if table.exists(key) {
                return table.get(key), nil
            }
        }
        switch meta := state.metafield(object, event.ID()).(type) {
            case *Closure:
                state.frame().push(meta)
                state.frame().push(object)
                state.frame().push(key)
                state.Call(2, 1)
                return state.frame().pop(), nil
            case *table:
                object = meta
                continue
            default:
                return None, nil
        }
    }
    return None, fmt.Errorf("'__index' chain too long; possible loop")
}

// tryMetaBinary performs one of the following binary operations:
//
// __add: the addition (+) operation. First, Lua checks the lhs operand. If that
// operand does not define a metamethod for __add, then Lua will check the rhs
// operand. If Lua can find a metamethod, it calls the metamethod with the two
// operands as arguments, and the result of the call (adjusted to one value) is
// the result of the operation. Otherwise, it raises an error.
func tryMetaBinary(state *State, lhs, rhs Value, event metaEvent) (Value, error) {
    if meta := state.metafield(lhs, event.ID()); !IsNone(meta) { // try lhs operand
        if cls, ok := meta.(*Closure); ok {
            state.frame().push(cls)
            state.frame().push(lhs)
            state.frame().push(rhs)
            state.Call(2, 1)
            return state.frame().pop(), nil
        }
    }
    if meta := state.metafield(rhs, event.ID()); !IsNone(meta) { // try rhs operand
        if cls, ok := meta.(*Closure); ok {
            state.frame().push(cls)
            state.frame().push(lhs)
            state.frame().push(rhs)
            state.Call(2, 1)
            return state.frame().pop(), nil
        }
    }
    return None, fmt.Errorf("attempt to apply %s on %v %v value", event.ID(), lhs.Type(), rhs.Type())
}

// tryMetaCompare performs one of the follow Lua comparison metamethods: __lt, __le, __eq
//
// __lt: the less than (<) operation. Behavior similar to the addition operation, except that
// Lua will try a metamethod only when the values being compared are neither both numbers nor
// both strings. The result of the call is always converted to a boolean.
//
// __le: the less equal (<=) operation. Unlike other operations, the less-equal operation can
// use two different events. First, Lua looks for the __le metamethod in both operands, like in
// the less than operation. If it cannot find such a metamethod, then it will try the __lt metamethod,
// assuming that a <= b is equivalent to (not b < a). As with the other comparison operators, the result
// is always a boolean. (This use of the __lt event can be removed in future versions; it is also slower
// than a real __le metamethod.)
//
// __eq: the equal (==) operation. Behavior similar to the addition operation, except that Lua will try a
// metamethod only when the values being compared are either both tables or both full userdata and they are
// not primitively equal. The result of the call is always converted to a boolean.
//
// See https://www.lua.org/manual/5.3/manual.html#2.4
func tryMetaCompare(state *State, lhs, rhs Value, event metaEvent) (cmp bool, err error) {
    if meta := state.metafield(lhs, event.ID()); !IsNone(meta) { // try lhs operand
        if cls, ok := meta.(*Closure); ok {
            state.frame().push(cls)
            state.frame().push(lhs)
            state.frame().push(rhs)
            state.Call(2, 1)
            cmp = state.ToBool(-1)
            state.Pop()
            return cmp, nil
        }
    }
    if meta := state.metafield(rhs, event.ID()); !IsNone(meta) { // try rhs operand
        if cls, ok := meta.(*Closure); ok {
            state.frame().push(cls)
            state.frame().push(lhs)
            state.frame().push(rhs)
            state.Call(2, 1)
            cmp = state.ToBool(-1)
            state.Pop()
            return cmp, nil
        }
    }
    if event == metaLe {
        cmp, err = tryMetaCompare(state, rhs, lhs, metaLt)
        return !cmp, err
    }
    return false, fmt.Errorf("attempt to compare %s with %s", lhs.Type(), rhs.Type())
}

// tryMetaConcat (__concat) performs the concatenation (..) operation. Behavior similar
// to the addition operation, except that Lua will try a metamethod if any operand is
// neither a string nor a number (which is always coercible to a string).
func tryMetaConcat(state *State, lhs, rhs Value) (Value, error) {
    const event = metaConcat

    if meta := state.metafield(lhs, event.ID()); !IsNone(meta) { // try lhs operand
        if cls, ok := meta.(*Closure); ok {
            state.frame().push(cls)
            state.frame().push(lhs)
            state.frame().push(rhs)
            state.Call(2, 1)
            return state.frame().pop(), nil
        }
    }
    if meta := state.metafield(rhs, event.ID()); !IsNone(meta) { // try rhs operand
        if cls, ok := meta.(*Closure); ok {
            state.frame().push(cls)
            state.frame().push(lhs)
            state.frame().push(rhs)
            state.Call(2, 1)
            return state.frame().pop(), nil
        }
    }
    return None, fmt.Errorf("attempt to apply '%s' on %v and %v values", event.ID(), lhs.Type(), rhs.Type())
}

// tryMetaLength (__len) performs the length (#) operation. If the object is not a string, Lua
// will try its metamethod. If there is a metamethod, Lua calls it with the object as argument,
// and the result of the call (always adjusted to one value) is the result of the operation.
// If there is no metamethod but the object is a table, then Lua uses the table length operation
// (see §3.4.7). Otherwise, Lua raises an error.
func tryMetaLength(state *State, obj Value) (Value, error) {
    if meta := state.metafield(obj, "__len"); !IsNone(meta) {
        if cls, ok := meta.(*Closure); ok {
            state.frame().push(cls)
            state.frame().push(obj)
            state.Call(1, 1)
            return state.frame().pop(), nil
        }
    }
    return nil, fmt.Errorf("attempt to get length of %v value", obj.Type())
}

// tryMetaCall performs the call operation func(args). This event happens when
// Lua tries to call a non-function value (that is, func is not a function).
// The metamethod is looked is looked up in func. If present, the metamethod
// is called with func as its first argument, followed by the arguments of the
// origin call (args). All results of the call are the result of the operation.
// This is the only metamethod that allows multiple results.
func tryMetaCall(state *State, value Value, fnID, args, rets int) bool {
    const event = metaCall

	if meta := state.metafield(value, event.ID()); !IsNone(meta) {
        if cls, ok := meta.(*Closure); ok {
			state.Push(cls)
			state.Insert(-(args+2))
			args++
			state.Call(args, rets)
            return true
        }
    }
    return false
}