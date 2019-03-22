package lua

import (
	"fmt"
	"os"
	"github.com/fibonacci1729/golua/lua/code"
)

var _ = fmt.Println
var _ = os.Exit

func gettable(ls *thread, t, k Value) (Value, error) {
	// - If 't' is a table and 't[k]' is not nil, return value.
	// - Otherwise check 't' for '__index' metamethod.
	// - If metamethod is nil, return nil.
	// - If metamethod exists and table, repeat lookuped with t = m.
	// - If metamethod exists and function, call 't.__index(t, k)'.
	for loop := 0; loop < maxMetaLoop; loop++ {
		if t, ok := t.(*Table); ok {
			if v := t.Get(k); v != nil {
				return v, nil
			}
		}
		switch m := ls.meta(t, "__index").(type) {
			case callable:
				panic("gettable: callable")
			case *Table:
				t = m
			default:
				return nil, nil
		}
	}
	return nil, fmt.Errorf("'__index' chain too long; possible loop")
}

func settable(ls *thread, t, k, v Value) error {
	// - If 't' is a table and 't[k]' is not nil, then 't[k]=v' and return nil.
	// - Otherwise check 't' for '__newindex' metamethod.
	// - If metamethod is nil, return nil.
	// - If metamethod exists and table, repeat lookup with t = m.
	// - If metamethod exists and function, call 't.__index(t, k)'.
	for loop := 0; loop < maxMetaLoop; loop++ {
		if t, ok := t.(*Table); ok {
			if v := t.Get(k); v != nil {
				t.Set(k, v)
				return nil
			}
		}
		if m := ls.meta(t, "__newindex"); m != nil {
			switch m := m.(type) {
				case callable:
					panic("set: callable!")
					// lua۰push(t, m, tbl, key, value)
					// lua۰call(t, 3, 0)
					// 	panic(runtimeErr("settable: call __newindex: todo"))
				case *Table:
					t = m
					continue
			}
		}
		if t, ok := t.(*Table); ok {
			t.Set(k, v)
		}
		return nil
	}
	return fmt.Errorf("'__newindex' chain too long; possible loop")
}

func compare(ls *thread, op Op, x, y Value) (bool, error) {
	switch op {
		case OpNe, OpEq:
			switch eq, err := equals(ls, x, y); {
				case err != nil:
					return false, err
				case op == OpNe:
					return !eq, nil
				default:
					return eq, nil
			}
		case OpGt:
			lt, err := less(ls, y, x)
			if err != nil {
				return false, err
			}
			return lt, nil
		case OpGe:
			le, err := lesseq(ls, y, x)
			if err != nil {
				return false, err
			}
			return le, nil
		case OpLt:
			return less(ls, x, y)
		case OpLe:
			return lesseq(ls, x, y)
	}
	panic(fmt.Errorf("unexpected comparison operator '%v'", op))
}

func equals(ls *thread, x, y Value) (bool, error) {
	return false, fmt.Errorf("equals: todo!")
}

func lesseq(ls *thread, x, y Value) (bool, error) {
	switch x := x.(type) {
		case String:
			if y, ok := y.(String); ok {
				return x <= y, nil
			}
		case Float:
			switch y := y.(type) {
				case Float:
					return x <= y, nil
				case Int:
					return x <= Float(y), nil
			}
		case Int:
			switch y := y.(type) {
				case Float:
					return Float(x) <= y, nil
				case Int:
					return x <= y, nil
			}
	}

	return false, fmt.Errorf("lesseq: meta: todo!")
}

func less(ls *thread, x, y Value) (bool, error) {
	switch x := x.(type) {
		case String:
			if y, ok := y.(String); ok {
				return x < y, nil
			}
		case Float:
			switch y := y.(type) {
				case Float:
					return x < y, nil
				case Int:
					return x < Float(y), nil
			}
		case Int:
			switch y := y.(type) {
				case Float:
					return Float(x) < y, nil
				case Int:
					return x < y, nil
			}
	}
	return false, fmt.Errorf("less: meta: todo!")
}

func length(ls *thread, x Value) (Value, error) {
	switch x := x.(type) {
		case String:
			return Int(len(x)), nil
		case *Table:
			if ls == nil || x.meta == nil {
				return x.Length(), nil
			}
	}
	return nil, fmt.Errorf("length: meta: todo!")
}

// UNM, BNOT, NOT, LEN
func unary(ls *thread, op Op, x Value) (Value, error) {
	switch op {
		case OpMinus:
			return binary(ls, OpMinus, x, Int(0))
		case OpBnot:
			return binary(ls, OpBnot, x, Int(0))
		case OpNot:
			return Bool(!Truth(x)), nil
		case OpLen:
			return length(ls, x)
	}
	panic(fmt.Errorf("unexpected unary operator '%v'", op))
}

func (ls *thread) exec(fn *Func) (rets []Value, err error) {
	const trace = false
	var (
		fp = fn.proto
		fr = ls.fr
	)
	frame:
		for ci := fr.call; ci.pc < len(fp.Instrs); ci.pc++ {
			if trace {
				fmt.Printf("[%d] %v\n", ci.pc, fp.Instrs[ci.pc])
			}
			switch inst := fp.Instrs[ci.pc]; inst.Code() {
				// Unary operators.
				//
				// @args A B
				//
				// R(A) := OP RK(B)
		        case code.BNOT, code.UNM, code.NOT, code.LEN:
					var (
						op = Op(inst.Code()-code.UNM)+OpMinus
						x  = fn.rk(inst.B())
						v Value
					)
					if v, err = unary(ls, op, x); err != nil {
						break frame
					}
					fn.stack[inst.A()] = v
		        
				// Comparison operators with conditional jump.
				//
				// @args A B C
				//
 				// if ((RK(B) OP RK(C)) ~= A) then pc++
				case code.EQ, code.LT, code.LE:
				    var (
						op = Op(inst.Code()-code.EQ)+OpEq
				        x  = fn.rk(inst.B())
				        y  = fn.rk(inst.C())
						v bool
				    )
				    if v, err = compare(ls, op, x, y); err != nil {
						break frame
					}
					if v != (inst.A() == 1) {
				    	ci.pc++
				    }

				// Binary operators.
				//
				// @args A B C
				//
				// R(A) := RK(B) OP RK(C) 
		        case code.ADD,
		        	code.SUB,
		        	code.MUL,
		        	code.MOD,
		        	code.POW,
		        	code.DIV,
		        	code.IDIV,
		        	code.BAND,
		        	code.BOR,
		        	code.BXOR,
		        	code.SHL,
		        	code.SHR:
					var (
						op = Op(inst.Code()-code.ADD)+OpAdd
						x  = fn.rk(inst.B())
						y  = fn.rk(inst.C())
						v Value
					)
					if v, err = binary(ls, op, x, y); err != nil {
						break frame
					}
				    fn.stack[inst.A()] = v

				// CONCAT: Concatenate a range of registers.
				//
				// @args A B C
				//
				// R(A) := R(B).. ... ..R(C)
				case code.CONCAT:
					// if xs := fn.stack[inst.B():inst.C()+1]; len(xs) > 1 {
					//     fn.stack[inst.A()] = ls.concat(xs)
					//     fn.sp = inst.A()    
					// }
					panic("code.CONCAT")

				// FORPREP: Initialization for a numeric for loop.   
				//
				// @args A sBx 
				//
				// R(A) -= R(A+2); pc+=sBx
				case code.FORPREP:
					panic("code.FORPREP")

				// FORLOOP: Iterate a numeric for loop.
				//
				// @args A sBx 
				//
				// R(A) += R(A+2); if R(A) <?= R(A+1) then { pc+=sBx; R(A+3)=R(A) }
				case code.FORLOOP:
					panic("code.FORLOOP")

				// TFORCALL: Iterate a generic for loop.
				//
				// R(A) is the iterator function, R(A+1) is the state, R(A+2) is the
				// control variable. At the start, R(A+2) has an initial value.
				//
				// Loop variables reside at locations R(A+3) and up, and their count
				// is specified in operand C. Operand C must be at least 1.
				//
				// Each time tforcall executes, the iterator function referenced by
				// R(A) is called with two arguments, the state R(A+1) and control
				// variable R(A+2). The results are returned in the local loop
				// variables, from R(A+3) up to R(A+2+C).
				//
				// @args A C
				//
				// R(A+3), ... ,R(A+2+C) := R(A)(R(A+1), R(A+2))
				case code.TFORCALL:
					// tforcall expects the for variables below to be at a fixed
					// position in the stack for every iteration, so we need to
					// adjust the stack to ensure this to avoid side effects.
					var (
			    		ctrl = fn.stack[inst.A()+2]
						data = fn.stack[inst.A()+1]
					    iter = fn.stack[inst.A()]
						base = inst.A() + 3
						rvs []Value
					)
					rvs, err = ls.call(iter, []Value{data, ctrl}, inst.C())
					if err != nil {
						break frame
					}
					for i, ret := range rvs {
						fn.stack[base+i] = ret
					}

				// TFORLOOP: Initialization for a generic for loop.
				//
				// @args A sBx 
				//
				// if R(A+1) ~= nil then { R(A)=R(A+1); pc += sBx } 
				case code.TFORLOOP:
					if ctrl := fn.stack[inst.A()+1]; ctrl != nil { // continue loop?
					    fn.stack[inst.A()] = ctrl // save control variable
					    ci.pc += inst.SBX() // jump back
					}

				// NEWTABLE: Create a new table.
				//
				// Creates a new empty table at register R(A). B and C are the encoded size information
				// for the array part and the hash part of the table, respectively. Appropriate values
				// for B and C are set in order to avoid rehashing when initially populating the table
				// with array values or hash key-value pairs.
				//
				// Operand B and C are both encoded as a "floating point byte" (see lobject.c)
				// which is eeeeexxx in binary, where x is the mantissa and e is the exponent.
				// The actual value is calculated as 1xxx*2^(eeeee-1) if eeeee is greater than
				// 0 (a range of 8 to 15*2^30).
				//
				// If eeeee is 0, the actual value is xxx (a range of 0 to 7.)
				//
				// If an empty table is created, both sizes are zero. If a table is created with a number
				// of objects, the code generator counts the number of array elements and the number of
				// hash elements. Then, each size value is rounded up and encoded in B and C using the
				// floating point byte format.
				//
				// @args A B C
				//
				// R(A) := {} (size = B,C)
				case code.NEWTABLE:
				    var (
				        arrN = fb2int(inst.B())
				        kvsN = fb2int(inst.C())
				    )
				    fn.stack[inst.A()] = NewTableSize(arrN, kvsN)

				// GETTABLE: Read a table element into a register (locals).
				//
				// @args A B C
				//
				// R(A) := R(B)[RK(C)]
				case code.GETTABLE:
				    var (
				        t = fn.stack[inst.B()]
				        k = fn.rk(inst.C())
						v Value
				    )
					if v, err = gettable(ls, t, k); err != nil {
						break frame
					}
					fn.stack[inst.A()] = v

				// SETTABLE: Write a register value into a table element (locals).
				//
				// @args A B C
				//
				// R(A)[RK(B)] := RK(C)
				case code.SETTABLE:
				    var (
				        t = fn.stack[inst.A()]
				        k = fn.rk(inst.B())
				        v = fn.rk(inst.C())
				    )
					if err := settable(ls, t, k, v); err != nil {
						break frame
					}

				// SETLIST: Set a range of array elements for a table.
				//
				// Sets the values for a range of array elements in a table referenced by R(A). Field B is the number
				// of elements to set. Field C encodes the block number of the table to be initialized.
				//
				// The values used to initialize the table are located in registers R(A+1), R(A+2), and so on.
				//
				// The block size is denoted by FPF. FPF is ‘fields per flush’, defined as LFIELDS_PER_FLUSH in the source
				// file lopcodes.h, with a value of 50. For example, for array locations 1 to 20, C will be 1 and B will
				// be 20.
				// 
				// If B is 0, the table is set with a variable number of array elements, from register R(A+1) up to the top
				// of the stack. This happens when the last element in the table constructor is a function call or a vararg
				// operator.
				// 
				// If C is 0, the next instruction is cast as an integer, and used as the C value. This happens only when
				// operand C is unable to encode the block number, i.e. when C > 511, equivalent to an array index greater
				// than 25550.
				//
				// @args A B C
				//
				// R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
				case code.SETLIST:
				    var (
				        a = inst.A()
				        b = inst.B()
				        c = inst.C()
				    )
				    if b == 0 {
				        b = (ci.sp - a) - 1
				    }
				    if c == 0 {
				        // ASSERT: fn.Instrs[fn.pc+1] == EXTRAARG)
				        c = fp.Instrs[ci.pc+1].AX()
				        ci.pc++
				    }
				    o := (c - 1) * fieldsPerFlush + b
				    t := fn.stack[a].(*Table)
				    for b > 0 {
				        t.Set(Int(o), fn.stack[a + b])
				        o--
				        b--
				    }
				
				// SELF: Prepare an object method for calling.  
				//
				// @args A B C
				//
				// R(A+1) := R(B); R(A) := R(B)[RK(C)] 
				case code.SELF:
				    var (
				        self = fn.stack[inst.B()]
				        k    = fn.rk(inst.C())
				        v Value   
				    )
					v, err = gettable(ls, self, k)
					if err != nil {
						break frame
					}
				    fn.stack[inst.A()+1] = self
				    fn.stack[inst.A()] = v

				// GETTABUP: Read a value from table in
				// up-value into a register (globals).
				//
				// @args A B C
				//   
				// R(A) := UpValue[B][RK(C)]
		        case code.GETTABUP:
				    var (
				        t = fn.up[inst.B()].get()
				        k = fn.rk(inst.C())
						v Value
				    )
					if v, err = gettable(ls, t, k); err != nil {
						break frame
					}
					fn.stack[inst.A()] = v

				// SETTABUP: Write a register value into table in up-value (globals).
				//
				// @args A B C
				//
				// UpValue[A][RK(B)] := RK(C)
				case code.SETTABUP:
				    var (
				        t = fn.up[inst.A()].get()
				        k = fn.rk(inst.B())
				        v = fn.rk(inst.C())
				    )
				    if err := settable(ls, t, k, v); err != nil {
						break frame
					}

				// GETUPVAL: Read an upvalue into a register.
				//
				// @args A B
				//       
				// R(A) := UpValue[B]
				case code.GETUPVAL:
					// fn.stack[inst.A()] = fn.up[inst.B()].get()
					panic("code.GETUPVAL")

				// SETUPVAL: Write a register value into an upvalue.
				//
				// @args A B
				//
				// UpValue[B] := R(A)
				case code.SETUPVAL:
				    // fn.up[inst.B()].set(fn.stack[inst.A()])
					panic("code.SETUPVAL")
				
				// TESTSET: Boolean test, with conditional jump and assignment.
				//
				// @args A B C
				//
				// if (R(B) <=> C) then R(A) := R(B) else pc++
				case code.TESTSET:
				    if Truth(fn.stack[inst.B()]) != (inst.C() == 1) {
				    	ci.pc++
				    }
				    fn.stack[inst.A()] = fn.stack[inst.B()]

				// TEST: Boolean test, with conditional jump.    
				//
				// @args A C 
				//
				// if not (R(A) <=> C) then pc++  
				case code.TEST:
				    if Truth(fn.stack[inst.A()]) != (inst.C() == 1) {
				        ci.pc++
				    }

				// LOADNIL: Load nil values into a range of registers.
				// 
				// @args A B
				//
				// R(A), R(A+1), ..., R(A+B) := nil
				case code.LOADNIL:
					for i := inst.A(); i <= inst.A() + inst.B(); i++ {
					    fn.stack[i] = nil
					}

				// LOADBOOL: Load a boolean into a register.
				//
				// @args A B C
				//
				// R(A) := (Bool)B; if (C) pc++
				case code.LOADBOOL:
					truth := (Bool(inst.B() == 1))
    				fn.stack[inst.A()] = truth
					if inst.C() != 0 {
        				ci.pc++
					}

				// LOADKX: Load a constant into a register.
				// The next 'instruction' is always EXTRAARG.
				//
				// @args A
				//
				// R(A) := Kst(extra arg)
				case code.LOADKX:
					fn.stack[inst.A()] = fn.kst(fp.Instrs[ci.pc+1].AX())
					ci.pc++

				// LOADK: Load a constant into a register.
				//
				// @args A Bx
				//
				// R(A) := Kst(Bx)
				case code.LOADK:
					fn.stack[inst.A()] = fn.kst(inst.BX())

				// MOVE: Copy a value between registers.
				//
				// @args A B
				//
				// R(A) := R(B)
				case code.MOVE:
    				fn.stack[inst.A()] = fn.stack[inst.B()]

				// JMP: Unconditional jump.
				//
				// @args A sBx
				//
				// pc+=sBx; if (A) close all upvalues >= R(A-1)
				case code.JMP:
    				// if (A) close all upvalues >= R(A-1)
    				if inst.A() != 0 {
						fn.close(inst.A()-1)
					}
			    	ci.pc += inst.SBX()

				// CLOSURE: Create a closure of a function prototype.
				//
				// @args A Bx
				// 
				// R(A) := closure(KPROTO[Bx])
				case code.CLOSURE:
					cls := &Func{proto: fp.Protos[inst.BX()]}
					fn.stack[inst.A()] = cls
					cls.open(fn.stack, fn.up...)

				// VARARG: Assign vararg function arguments to registers.
 				//
				// VARARG copies B-1 parameters into a number of registers starting from R(A),
				// padding with nils if there aren’t enough values. If B is 0, VARARG copies as
				// many values as it can based on the number of parameters passed.
				//
				// If a fixed number of values is required, B is a value greater than 1.
				// If any number of values is required, B is 0.
				//
				// If B == 0, load all varargs.
				// If B >= 1, load B-1 varargs.
				//
				// @args A B
				//
				// R(A), R(A+1), ..., R(A+B-2) = vararg
				case code.VARARG:
					var (
					    b = inst.B() - 1
						a = inst.A()
					    n = len(ci.va)
					    i int
					)
					if b < 0 {
					    // fn.check(inst.A() + n)
					    ci.sp = inst.A() + n
					    b = n
					}
					for i < b && i < n {
					    fn.stack[a+i] = ci.va[i]
					    i++
					}
					for i < b {
					    fn.stack[a+i] = nil
					    i++
					}

				// TAILCALL: Perform a tail call.
				//
				// TAILCALL performs a tail call, which happens when a return statement has a single
				// function call as the expression, e.g. return foo(bar). A tail call results in the
				// function being interpreted within the same call frame as the caller -- the stack
				// is replaced and then a 'goto' executed to start at the entry point in the VM. Only
				// Lua functions can be tailcalled. Tail calls allow infinite recursion without growing
				// the stack.
				//
				// Like OP_CALL, registry R(A) holds the reference to the function object to be called.
				// B encodes the number of parameters in the same way as in OP_CALL.
				//
				// C isn't used by TAILCALL, since all return results are used. In any case, Lua always
				// generates a 0 for C denoting multiple return results. 
				//
				// @args A B C
				//
				// return R(A)(R(A+1), ... ,R(A+B-1))
				case code.TAILCALL:
					panic("tailcall!")

				// CALL: Calls a function.
				//
				// CALL performs a function call, with register R(A) holding the reference to
				// the function object to be called. Parameters to the function are placed in
				// the registers following R(A).
				//
				// If B is 1, the function has no parameters.
				//
				// If B is 2 or more, there are (B-1) parameters, and upon entry entry to the
				// called function, R(A+1) will become the base.
				//
				// If B is 0, then B = 'top', i.e., the function parameters range from R(A+1) to
				// the top of the stack. This form is used when the number of parameters to pass
				// is set by the previous VM instruction, which has to be one of OP_CALL or OP_VARARG.
				//
				// If C is 1, no return results are saved. If C is 2 or more, (C-1) return values are
				// saved.
				//
				// If C == 0, then 'top' is set to last_result+1, so that the next open instruction
				// (i.e. OP_CALL, OP_RETURN, OP_SETLIST) can use 'top'.
				//
				// If C > 1, results returned by the function call are placed in registers ranging
				// from R(A) to R(A+C-1).
				//
				// @args A B C
				//
				// R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))
				case code.CALL:
					var (
						base = inst.A() + 1
						argc = inst.B() - 1
						want = inst.C() - 1
						args []Value
						rvs  []Value
					)
					if argc < 0 {
						args = fn.stack[base:(ci.sp-inst.A())+base]
						// panic("call: argc < 0")
					} else {
						args = fn.stack[base:base+argc]
					}
					rvs, err = ls.call(fn.stack[base-1], args, want)
					if err != nil {
						break frame
					}
					fn.checkstack(base, len(rvs))
					for i, ret := range rvs {
						fn.stack[base+i-1] = ret
					}
					if want < 0 {
						ci.sp = inst.A() + len(rvs) - 1
						// panic("call: want < 0")
					} else {
						ci.sp = base + want
					}

				// RETURN: Returns from function call.
				//
				// Returns to the calling function, with optional return values.
				// First, op.RETURN closes any open upvalues by calling fn.close().
				//
				// If B == 0, the set of values ranges from R(A) to the top of the stack.
				// If B == 1, there are no return values.
				// If B >= 2, there are (B-1) return values, located in consecutive
				// register from R(A) ... R(A+B-1).
				//
				// If B == 0, then the previous instruction (which must be either op.CALL or
				// op.VARARG) would have set state top to indicate how many values to return.
				// The number of values to be returned in this case is R(A) to ci.top.
				//
				// If B > 0, then the number of values to be returned is simply B-1.
				// 
				// If (B == 0) then return up to 'top'.
				//
				// @args A B
				// 
				// return R(A), ... ,R(A+B-2)
		        case code.RETURN:
					var (
						b = inst.B() - 1
						a = inst.A()
					)
					if b < 0 {
						b = ci.sp - a
					}
					rets = fn.stack[a:a+b]
					break frame

				// EXTRAARG: Extra (larger) argument for previous opcode.
				//
				// @args Ax
				case code.EXTRAARG:
    				// This op func should never execute directly.
			    	panic("unreachable")

				default:
					panic(fmt.Errorf("unhandled instruction: %v", inst)) // Fatal
			}
		}

	return rets, ls.error(err)
}