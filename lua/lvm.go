package lua

import (
	"fmt"

	//"github.com/Azure/golua/pkg/luautil"
	"github.com/Azure/golua/lua/ir"
)

var _ = fmt.Println

//
// Implementation of Lua v53 Opcodes
//

// MOVE: Copy a value between registers.
//
// @args A B
//
// R(A) := R(B)
func (vm *v53) move(instr ir.Instr) {
	rb := vm.frame().get(instr.B())
	vm.frame().set(instr.A(), rb)
}

// LOADK: Load a constant into a register.
//
// @args A Bx
//
// R(A) := Kst(Bx)
func (vm *v53) loadk(instr ir.Instr) {
	kst := vm.constant(instr.BX())
	vm.frame().set(instr.A(), kst)
}

// LOADKX: Load a constant into a register. The next 'instruction'
// is always EXTRAARG.
//
// @args A
//
// R(A) := Kst(extra arg)
func (vm *v53) loadkx(instr ir.Instr) {
	extra := vm.frame().step(1).AX()
	ra := vm.constant(extra)
	vm.frame().set(instr.A(), ra)
}

// LOADBOOL: Load a boolean into a register.
//
// @args A B C
//
// R(A) := (Bool)B; if (C) pc++
func (vm *v53) loadbool(instr ir.Instr) {
	vm.frame().set(instr.A(), Bool(instr.B() == 1))
	if instr.C() != 0 { vm.frame().step(1) }
}

// LOADNIL: Load nil values into a range of registers.
// 
// @args A B
//
// R(A), R(A+1), ..., R(A+B) := nil
func (vm *v53) loadnil(instr ir.Instr) {
	var (
		a = instr.A()
		b = instr.B()
	)
	for i := a; i <= a + b; i++ {
		vm.frame().set(i, None)
	}
}

// GETUPVAL: Read an upvalue into a register.
//
// @args A B
//       
// R(A) := UpValue[B]
func (vm *v53) getupval(instr ir.Instr) {
	var (
		a = instr.A()
		b = instr.B()
	)
	up := vm.frame().getUp(b)
	vm.frame().set(a, up.get())
}

// SETUPVAL: Write a register value into an upvalue.
//
// @args A B
//
// UpValue[B] := R(A)
func (vm *v53) setupval(instr ir.Instr) {
	var (
		a = instr.A()
		b = instr.B()
	)
	ra := vm.frame().get(a)
	vm.frame().setUp(b, ra)
}

// GETTABLE: Read a table element into a register (locals).
//
// @args A B C
//
// R(A) := R(B)[RK(C)]
func (vm *v53) gettable(instr ir.Instr) {
	var (
		a = instr.A()
		b = instr.B()
		c = instr.C()
	)
	t := vm.frame().get(b)
	k := vm.rk(c)
	v := vm.State.gettable(t, k, false)
	vm.frame().set(a, v)
}

// SETTABLE: Write a register value into a table element (locals).
//
// @args A B C
//
// R(A)[RK(B)] := RK(C)
func (vm *v53) settable(instr ir.Instr) {
	obj := vm.frame().get(instr.A())
	key := vm.rk(instr.B())
	val := vm.rk(instr.C())
	vm.State.settable(obj, key, val, false)
}

// GETTABUP: Read a value from table in up-value into a register (globals).
//
// @args A B C
//   
// R(A) := UpValue[B][RK(C)]
func (vm *v53) gettabup(instr ir.Instr) {
	up := vm.frame().getUp(instr.B()).get()
	rc := vm.rk(instr.C())
	ra := vm.State.gettable(up, rc, false)
	vm.frame().set(instr.A(), ra)
}

// SETTABUP: Write a register value into table in up-value (globals).
//
// @args A B C
//
// UpValue[A][RK(B)] := RK(C) 
func (vm *v53) settabup(instr ir.Instr) {
	up := vm.frame().getUp(instr.A()).get()
	rb := vm.rk(instr.B())
	rc := vm.rk(instr.C())
	vm.State.settable(up, rb, rc, false)
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
func (vm *v53) newtable(instr ir.Instr) {
	var (
		a = instr.A()
		b = instr.B()
		c = instr.C()
	)
	t := newTable(vm.State, fb2i(b), fb2i(c))
	vm.frame().set(a, &Table{t})
}

// SELF: Prepare an object method for calling.  
//
// @args A B C
//
// R(A+1) := R(B); R(A) := R(B)[RK(C)] 
func (vm *v53) self(instr ir.Instr) {
	var (
		obj = vm.frame().get(instr.B())
		key = vm.rk(instr.C())
		fn  = vm.State.gettable(obj, key, false)
	)
	vm.frame().set(instr.A(), fn)
	vm.frame().set(instr.A()+1, obj)
}

// ADD: Addition operator.  
//
// @args A B C
//
// R(A) := RK(B) + RK(C)               
func (vm *v53) add(instr ir.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.arith(OpAdd, rb, rc)
	)
	vm.frame().set(instr.A(), ra)
}

// SUB: Subtraction operator.
//    
// @args A B C
//
// R(A) := RK(B) - RK(C)
func (vm *v53) sub(instr ir.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.arith(OpSub, rb, rc)
	)
	vm.frame().set(instr.A(), ra)
}

// MUL: Multiplication operator.
//     
// @args A B C
//    
// R(A) := RK(B) * RK(C)
func (vm *v53) mul(instr ir.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.arith(OpMul, rb, rc)
	)
	vm.frame().set(instr.A(), ra)
}

// MOD: Modulus (remainder) operator.
//
// @args A B C
//
// R(A) := RK(B) % RK(C)                
func (vm *v53) mod(instr ir.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.arith(OpMod, rb, rc)
	)
	vm.frame().set(instr.A(), ra)
}

// POW: Exponentation operator.
//
// @args A B C
//
// R(A) := RK(B) ^ RK(C)               
func (vm *v53) pow(instr ir.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.arith(OpPow, rb, rc)
	)
	vm.frame().set(instr.A(), ra)
}

// DIV: Division operator.         
//
// @args A B C
//
// R(A) := RK(B) / RK(C)
func (vm *v53) div(instr ir.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.arith(OpDiv, rb, rc)
	)
	vm.frame().set(instr.A(), ra)
}

// UNM: Unary minus.
//
// @args A B
//
// R(A) := -R(B)
func (vm *v53) unm(instr ir.Instr) {
	var (
		rb = vm.rk(instr.B())
		ra = vm.arith(OpMinus, rb, None)
	)
	vm.frame().set(instr.A(), ra)
}

// NOT: Logical NOT operator.  
//
// @args A B
//
// R(A) := not R(B)
func (vm *v53) not(instr ir.Instr) {
	rb := vm.frame().get(instr.B())
	vm.frame().set(instr.A(), truth(rb))
}
                  
// LEN: Length operator.
//
// @args A B
//
// R(A) := length of R(B) 
func (vm *v53) length(instr ir.Instr) {
	rb := vm.frame().get(instr.B())
	vm.frame().set(instr.A(), Int(vm.State.length(rb)))
}

// CONCAT: Concatenate a range of registers.
//
// @args A B C
//
// R(A) := R(B).. ... ..R(C)
func (vm *v53) concat(instr ir.Instr) {
	var (
		a = instr.A()
		b = instr.B()
		c = instr.C()
	)
	vm.State.Concat(c-b+1)
	vm.frame().replace(a)
}

// JMP: Unconditional jump.
//
// @args A sBx
//
// pc+=sBx; if (A) close all upvalues >= R(A-1)
func (vm *v53) jmp(instr ir.Instr) {
	vm.frame().step(instr.SBX())
	if a := instr.A(); a != 0 {
		vm.frame().closeUp(a-1)
	}
}

// EQ: Equality test, with conditional jump.
//          
// @args A B C
//
// if ((RK(B) == RK(C)) ~= A) then pc++
func (vm *v53) eq(instr ir.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		aa = (instr.A() != 0)
	)
	if vm.compare(OpEq, rb, rc) != aa {
		vm.frame().step(1)
	}
}

// LT: Less than test, with conditional jump.
//
// @args A B C
//
// if ((RK(B) <  RK(C)) ~= A) then pc++
func (vm *v53) lt(instr ir.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		aa = (instr.A() == 1)
	)
	if vm.compare(OpLt, rb, rc) != aa {
		vm.frame().step(1)
	}
}

// LE: Less than or equal to test, with conditional jump.
//
// @args A B C
//
// if ((RK(B) <= RK(C)) ~= A) then pc++
func (vm *v53) le(instr ir.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		aa = (instr.A() == 1)
	)
	if vm.compare(OpLe, rb, rc) != aa {
		vm.frame().step(1)
	}
}

// TEST: Boolean test, with conditional jump.    
//
// @args A C 
//
// if not (R(A) <=> C) then pc++          
func (vm *v53) test(instr ir.Instr) {
	var (
		ra = vm.frame().get(instr.A())
		cc = (instr.C() == 1)
	)
	if Truth(ra) != cc {
		vm.frame().step(1)
	}
}

// TESTSET: Boolean test, with conditional jump and assignment.
//
// @args A B C
//
// if (R(B) <=> C) then R(A) := R(B) else pc++
func (vm *v53) testset(instr ir.Instr) {
	var (
		rb = vm.frame().get(instr.B())
		cc = (instr.C() == 1)
	)
	if Truth(rb) == cc {
		vm.frame().set(instr.A(), rb)
	} else {
		vm.frame().step(1)
	} 
}

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
func (vm *v53) call(instr ir.Instr) {
	var (
		a = instr.A()
		b = instr.B()
		c = instr.C()
	)
	// arguments
	if b != 0 {
		vm.frame().settop(a+b)
		vm.State.Call(b-1, c-1)
	} else {
		vm.State.Call(vm.frame().gettop()-a-1, c-1)
	}
	// returns
	if c != 0 {
		rets := vm.frame().popN(c-1)
		for i, v := range rets {
			vm.frame().set(a+i, v)
		}
	} else {
		// C=0 so return values indicated by 'top'
		// TODO ??
	}
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
func (vm *v53) tailcall(instr ir.Instr) {
	var (
		a = instr.A()
		b = instr.B()
		c = instr.C()
	)
	// TODO: optimize tailcalls (reuse frame)
	// TODO: vm.frame().closeups()
	if b != 0 {
		vm.frame().settop(a+b)
		vm.State.Call(b-1, c-1)
	} else {
		vm.State.Call(vm.frame().gettop()-a-1, c-1)	
	}
}

// RETURN: Returns from function call.
//
// Returns to the calling function, with optional return values. First, op.RETURN closes any
// open upvalues by calling frame.closeup().
//
// If B == 0, the set of values ranges from R(A) to the top of the stack.
// If B == 1, there are no return values.
// If B >= 2, there are (B-1) return values, located in consecutive
// register from R(A) ... R(A+B-1).
//
// It is assumed that if the VM is returning to a lua function, then it is within the
// same invocation of 'exec'. Otherwise, it is assumed that 'exec' is being invoked
// from a Go function.
//
// If B == 0, then the previous instruction (which must be either op.CALL or op.VARARG)
// would have set state top to indicate how many values to return. The number of values
// to be returned in this case is R(A) to state.GetTop().
//
// If B > 0, then the number of values to be returned is simply B-1.
// 
// If (B == 0) then return up to 'top'.
//
// @args A B
// 
// return R(A), ... ,R(A+B-2)
func (vm *v53) returns(instr ir.Instr) {
	var (
		a = instr.A()
		b = instr.B()
	)
	//fmt.Printf("RETURN A=%d B=%d (want = %d)\n", a, b, vm.frame().rets)
	//vm.State.Debug(false)
	if want := vm.frame().rets; want != 0 {
		var (
			retc int = b - 1
			rets []Value
		)
		if b == 0 {
			retc = vm.frame().gettop()-a
		}
		switch {
			case want > retc: // # wanted > # returned
				for i := a; i < a + retc; i++ {
					rets = append(rets, vm.frame().get(i))
				}
				for retc < want {
					rets = append(rets, None)
					retc++
				}
				
			case want <= retc: // # wanted <= # returned
				if want == MultRets {
					want = retc
				}
				for i := a; i < a + want; i++ {
					rets = append(rets, vm.frame().get(i))
				}
		}
		vm.frame().caller().pushN(rets)
	}
}

// FORLOOP: Iterate a numeric for loop.
//
// @args A sBx 
//
// R(A)+=R(A+2); if R(A) <?= R(A+1) then { pc+=sBx; R(A+3)=R(A) }                     
func (vm *v53) forloop(instr ir.Instr) {
	var (
		item = vm.frame().get(instr.A())
		upto = vm.frame().get(instr.A()+1)
		step = vm.frame().get(instr.A()+2)
	)
	if isInteger(item) { // integer loop?
		i1 := item.(Int) 
		i2 := upto.(Int)
		i3 := step.(Int)
		i1 += i3 // increment index
		if (i3 > 0 && (i1 <= i2)) || (i3 < 0 && (i1 > i2)) {
			vm.frame().set(instr.A(), i1)   // update internal index...
			vm.frame().set(instr.A()+3, i1) // ... and external index
			vm.frame().step(instr.SBX())    // jump back
		}
	} else { // floating loop
		f1 := item.(Float)
		f2 := upto.(Float)
		f3 := step.(Float)
		f1 += f3
		if (f3 > 0 && (f1 <= f2)) || (f3 < 0 && (f1 > f2)) {
			vm.frame().set(instr.A(), f1)   // update internal index...
			vm.frame().set(instr.A()+3, f1) // ... and external index
			vm.frame().step(instr.SBX())    // jump back
		}
	}
}

// FORPREP: Initialization for a numeric for loop.   
//
// @args A sBx 
//
// R(A)-=R(A+2); pc+=sBx
func (vm *v53) forprep(instr ir.Instr) {
	var (
		init = vm.frame().get(instr.A())
		upto = vm.frame().get(instr.A()+1)
		step = vm.frame().get(instr.A()+2)
	)
	// Try for values as integers.
	var (
		i1, ok1 = toInteger(init)
		i2, ok2 = toInteger(upto)
		i3, ok3 = toInteger(step)
	)
	if ok1 && ok2 && ok3 {
		vm.frame().set(instr.A(), i1-i3)
		vm.frame().set(instr.A()+1, i2)
		vm.frame().set(instr.A()+2, i3)
		vm.frame().step(instr.SBX())
		return
	}
	// Try for values as numbers.
	var f1, f2, f3 Float
	if f1, ok1 = toFloat(init); !ok1 {
		vm.State.errorf("'for' init must be a number")
	}
	if f2, ok2 = toFloat(upto); !ok2 {
		vm.State.errorf("'for' limit must be a number")
	}
	if f3, ok3 = toFloat(step); !ok3 {
		vm.State.errorf("'for' step must be a number")
	}
	
	vm.frame().set(instr.A(), f1-f3)
	vm.frame().set(instr.A()+1, f2)
	vm.frame().set(instr.A()+2, f3)
	vm.frame().step(instr.SBX())
}

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
func (vm *v53) tforcall(instr ir.Instr) {
	// tforcall expects the for variables below to be at a fixed
	// position in the stack for every iteration, so we need to
	// adjust the stack to ensure this to avoid side effects.

	var (
		a = instr.A()
		c = instr.C()
	)

	var (
		iter = vm.frame().get(a)   // iterator function
		data = vm.frame().get(a+1) // state
		ctrl = vm.frame().get(a+2) // control variable / initial value
		base = instr.A()+3
	)

	vm.frame().push(iter)
	vm.frame().push(data)
	vm.frame().push(ctrl)

	vm.State.Call(2, c)


	rets := vm.frame().popN(c)
	for i, v := range rets {
		vm.frame().set(base+i, v)
	}
}

// TFORLOOP: Initialization for a generic for loop.
//
// @args A sBx 
//
// if R(A+1) ~= nil then { R(A)=R(A+1); pc += sBx } 
func (vm *v53) tforloop(instr ir.Instr) {
	if ctrl := vm.frame().get(instr.A()+1); !IsNone(ctrl) {
		vm.frame().set(instr.A(), ctrl)
		vm.frame().step(instr.SBX())
		return
	}
	// loop done, reset top
	vm.frame().settop(instr.A())
}

// SETLIST: Set a range of array elements for a table.
//
// @args A B C
//
// R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
func (vm *v53) setlist(instr ir.Instr) {
	var (
		a = instr.A()
		b = instr.B()
		c = instr.C()
	)
	if b == 0 { b = vm.frame().gettop() - a - 1}
	o := (c-1) * FieldsPerFlush
	t := vm.frame().get(a).(*Table)
	for i := 1; i <= b; i++ {
		t.setInt(int64(o+i), vm.frame().get(a+i))
	}
	vm.frame().popN(b)
}

// CLOSURE: Create a closure of a function prototype.
//
// @args A Bx
// 
// R(A) := closure(KPROTO[Bx])
func (vm *v53) closure(instr ir.Instr) {
	cls := newLuaClosure(vm.prototype(instr.BX()))
	vm.frame().openUp(cls)
	vm.frame().push(cls)
	vm.frame().replace(instr.A())
	// TODO: caching?
}

// VARARG: Assign vararg function arguments to registers.
//
// If B == 0, load all varargs.
// If B >= 1, load B-1 varargs.
//
// @args A B
//
// R(A), R(A+1), ..., R(A+B-2) = vararg
func (vm *v53) vararg(instr ir.Instr) {
	var (
		a = instr.A()
		b = instr.B()
	)
	for i, v := range vm.frame().varargs(b-1) {
		if v == nil { v = None }
		vm.frame().set(a+i, v)	
	}
}

// IDIV: Integer division operator.
//      
// @args A B C
//
// R(A) := RK(B) // RK(C)
func (vm *v53) idiv(instr ir.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.arith(OpQuo, rb, rc)
	)
	vm.frame().set(instr.A(), ra)
}

// BAND: Bit-wise AND operator.        
//
// @args A B C
//
// R(A) := RK(B) & RK(C)
func (vm *v53) band(instr ir.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.arith(OpAnd, rb, rc)
	)
	vm.frame().set(instr.A(), ra)
}

// BOR: Bit-wise OR operator.
//        
// @args A B C
//
// R(A) := RK(B) | RK(C) 
func (vm *v53) bor(instr ir.Instr) {	
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.arith(OpOr, rb, rc)
	)
	vm.frame().set(instr.A(), ra)
}

// BXOR: Bit-wise Exclusive OR operator.
//      
// @args A B C
//
// R(A) := RK(B) ~ RK(C)
func (vm *v53) bxor(instr ir.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.arith(OpXor, rb, rc)
	)
	vm.frame().set(instr.A(), ra)
}

// SHL: Shift bits left.         
//
// @args A B C
//
// R(A) := RK(B) << RK(C) 
func (vm *v53) shl(instr ir.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.arith(OpLsh, rb, rc)
	)
	vm.frame().set(instr.A(), ra)
}

// SHR: Shift bits right.         
//
// @args A B C
//
// R(A) := RK(B) >> RK(C)
func (vm *v53) shr(instr ir.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.arith(OpRsh, rb, rc)
	)
	vm.frame().set(instr.A(), ra)
}

// BNOT: Bit-wise NOT operator.        
//
// @args A B
// 
// R(A) := ~R(B)
func (vm *v53) bnot(instr ir.Instr) {
	var (
		rb = vm.rk(instr.B())
		ra = vm.arith(OpNot, rb, None)
	)
	vm.frame().set(instr.A(), ra)
}

// EXTRAARG: Extra (larger) argument for previous opcode.
//
// @args Ax
func (vm *v53) extraarg(instr ir.Instr) {
	// This op func should never execute directly.
	unimplemented(instr.String())
}