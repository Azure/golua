package lua

import (
	"fmt"

	"github.com/Azure/golua/lua/vm"
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
func (vm *v53) move(instr vm.Instr) {
	rb := vm.thread().frame().get(instr.B())
	vm.thread().frame().set(instr.A(), rb)
}

// LOADK: Load a constant into a register.
//
// @args A Bx
//
// R(A) := Kst(Bx)
func (vm *v53) loadk(instr vm.Instr) {
	kst := vm.constant(instr.BX())
	vm.thread().frame().set(instr.A(), kst)
}

// LOADKX: Load a constant into a register. The next 'instruction'
// is always EXTRAARG.
//
// @args A
//
// R(A) := Kst(extra arg)
func (vm *v53) loadkx(instr vm.Instr) {
	extra := vm.thread().frame().step(1).AX()
	ra := vm.constant(extra)
	vm.thread().frame().set(instr.A(), ra)
}

// LOADBOOL: Load a boolean into a register.
//
// @args A B C
//
// R(A) := (Bool)B; if (C) pc++
func (vm *v53) loadbool(instr vm.Instr) {
	vm.thread().frame().set(instr.A(), Bool(instr.B() == 1))
	if instr.C() != 0 { vm.thread().frame().step(1) }
}

// LOADNIL: Load nil values into a range of registers.
// 
// @args A B
//
// R(A), R(A+1), ..., R(A+B) := nil
func (vm *v53) loadnil(instr vm.Instr) {
	var (
		a = instr.A()
		b = instr.B()
	)
	for i := a; i <= a + b; i++ {
		vm.thread().frame().set(i, Nil(1))
	}
}

// GETUPVAL: Read an upvalue into a register.
//
// @args A B
//       
// R(A) := UpValue[B]
func (vm *v53) getupval(instr vm.Instr) {
	var (
		a = instr.A()
		b = instr.B()
	)
	up := vm.thread().frame().getUp(b)
	vm.thread().frame().set(a, up.get())
}

// SETUPVAL: Write a register value into an upvalue.
//
// @args A B
//
// UpValue[B] := R(A)
func (vm *v53) setupval(instr vm.Instr) {
	var (
		a = instr.A()
		b = instr.B()
	)
	ra := vm.thread().frame().get(a)
	vm.thread().frame().setUp(b, ra)
}

// GETTABLE: Read a table element into a register (locals).
//
// @args A B C
//
// R(A) := R(B)[RK(C)]
func (vm *v53) gettable(instr vm.Instr) {
	var (
		a = instr.A()
		b = instr.B()
		c = instr.C()
	)
	t := vm.thread().frame().get(b)
	k := vm.rk(c)
	v := vm.thread().gettable(t, k, false)
	vm.thread().frame().set(a, v)
}

// SETTABLE: Write a register value into a table element (locals).
//
// @args A B C
//
// R(A)[RK(B)] := RK(C)
func (vm *v53) settable(instr vm.Instr) {
	obj := vm.thread().frame().get(instr.A())
	key := vm.rk(instr.B())
	val := vm.rk(instr.C())
	vm.thread().settable(obj, key, val, false)
}

// GETTABUP: Read a value from table in up-value into a register (globals).
//
// @args A B C
//   
// R(A) := UpValue[B][RK(C)]
func (vm *v53) gettabup(instr vm.Instr) {
	up := vm.thread().frame().getUp(instr.B()).get()
	rc := vm.rk(instr.C())
	ra := vm.thread().gettable(up, rc, false)
	vm.thread().frame().set(instr.A(), ra)
}

// SETTABUP: Write a register value into table in up-value (globals).
//
// @args A B C
//
// UpValue[A][RK(B)] := RK(C) 
func (vm *v53) settabup(instr vm.Instr) {
	up := vm.thread().frame().getUp(instr.A()).get()
	rb := vm.rk(instr.B())
	rc := vm.rk(instr.C())
	vm.thread().settable(up, rb, rc, false)
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
func (vm *v53) newtable(instr vm.Instr) {
	var (
		a = instr.A()
		b = instr.B()
		c = instr.C()
	)
	t := newTable(vm.thread(), fb2i(b), fb2i(c))
	vm.thread().frame().set(a, t)
}

// SELF: Prepare an object method for calling.  
//
// @args A B C
//
// R(A+1) := R(B); R(A) := R(B)[RK(C)] 
func (vm *v53) self(instr vm.Instr) {
	var (
		obj = vm.thread().frame().get(instr.B())
		key = vm.rk(instr.C())
		fn  = vm.thread().gettable(obj, key, false)
	)
	vm.thread().frame().set(instr.A(), fn)
	vm.thread().frame().set(instr.A()+1, obj)
}

// ADD: Addition operator.  
//
// @args A B C
//
// R(A) := RK(B) + RK(C)               
func (vm *v53) add(instr vm.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.thread().arith(OpAdd, rb, rc)
	)
	vm.thread().frame().set(instr.A(), ra)
}

// SUB: Subtraction operator.
//    
// @args A B C
//
// R(A) := RK(B) - RK(C)
func (vm *v53) sub(instr vm.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.thread().arith(OpSub, rb, rc)
	)
	vm.thread().frame().set(instr.A(), ra)
}

// MUL: Multiplication operator.
//     
// @args A B C
//    
// R(A) := RK(B) * RK(C)
func (vm *v53) mul(instr vm.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.thread().arith(OpMul, rb, rc)
	)
	vm.thread().frame().set(instr.A(), ra)
}

// MOD: Modulus (remainder) operator.
//
// @args A B C
//
// R(A) := RK(B) % RK(C)                
func (vm *v53) mod(instr vm.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.thread().arith(OpMod, rb, rc)
	)
	vm.thread().frame().set(instr.A(), ra)
}

// POW: Exponentation operator.
//
// @args A B C
//
// R(A) := RK(B) ^ RK(C)               
func (vm *v53) pow(instr vm.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.thread().arith(OpPow, rb, rc)
	)
	vm.thread().frame().set(instr.A(), ra)
}

// DIV: Division operator.         
//
// @args A B C
//
// R(A) := RK(B) / RK(C)
func (vm *v53) div(instr vm.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.thread().arith(OpDiv, rb, rc)
	)
	vm.thread().frame().set(instr.A(), ra)
}

// UNM: Unary minus.
//
// @args A B
//
// R(A) := -R(B)
func (vm *v53) unm(instr vm.Instr) {
	var (
		rb = vm.rk(instr.B())
		ra = vm.thread().arith(OpMinus, rb, None)
	)
	vm.thread().frame().set(instr.A(), ra)
}

// NOT: Logical NOT operator.  
//
// @args A B
//
// R(A) := not R(B)
func (vm *v53) not(instr vm.Instr) {
	rb := vm.thread().frame().get(instr.B())
	vm.thread().frame().set(instr.A(), !truth(rb))
}
                  
// LEN: Length operator.
//
// @args A B
//
// R(A) := length of R(B) 
func (vm *v53) length(instr vm.Instr) {
	rb := vm.thread().frame().get(instr.B())
	vm.thread().frame().set(instr.A(), vm.thread().length(rb))
}

// CONCAT: Concatenate a range of registers.
//
// @args A B C
//
// R(A) := R(B).. ... ..R(C)
func (vm *v53) concat(instr vm.Instr) {
	var (
		a = instr.A()
		b = instr.B()
		c = instr.C()
	)
	vm.thread().Concat(c-b+1)
	vm.thread().frame().replace(a)
}

// JMP: Unconditional jump.
//
// @args A sBx
//
// pc+=sBx; if (A) close all upvalues >= R(A-1)
func (vm *v53) jmp(instr vm.Instr) {
	vm.thread().frame().step(instr.SBX())
	if a := instr.A(); a != 0 {
		vm.thread().frame().closeUp(a-1)
	}
}

// EQ: Equality test, with conditional jump.
//          
// @args A B C
//
// if ((RK(B) == RK(C)) ~= A) then pc++
func (vm *v53) eq(instr vm.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		aa = (instr.A() != 0)
	)
	if vm.thread().compare(OpEq, rb, rc, false) != aa {
		vm.thread().frame().step(1)
	}
}

// LT: Less than test, with conditional jump.
//
// @args A B C
//
// if ((RK(B) <  RK(C)) ~= A) then pc++
func (vm *v53) lt(instr vm.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		aa = (instr.A() == 1)
	)
	if vm.thread().compare(OpLt, rb, rc, false) != aa {
		vm.thread().frame().step(1)
	}
}

// LE: Less than or equal to test, with conditional jump.
//
// @args A B C
//
// if ((RK(B) <= RK(C)) ~= A) then pc++
func (vm *v53) le(instr vm.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		aa = (instr.A() == 1)
	)
	if vm.thread().compare(OpLe, rb, rc, false) != aa {
		vm.thread().frame().step(1)
	}
}

// TEST: Boolean test, with conditional jump.    
//
// @args A C 
//
// if not (R(A) <=> C) then pc++          
func (vm *v53) test(instr vm.Instr) {
	var (
		ra = vm.thread().frame().get(instr.A())
		cc = (instr.C() == 1)
	)
	if Truth(ra) != cc {
		vm.thread().frame().step(1)
	}
}

// TESTSET: Boolean test, with conditional jump and assignment.
//
// @args A B C
//
// if (R(B) <=> C) then R(A) := R(B) else pc++
func (vm *v53) testset(instr vm.Instr) {
	var (
		rb = vm.thread().frame().get(instr.B())
		cc = (instr.C() == 1)
	)
	if Truth(rb) == cc {
		vm.thread().frame().set(instr.A(), rb)
	} else {
		vm.thread().frame().step(1)
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
func (vm *v53) call(instr vm.Instr) {
	var (
		a = instr.A()
		b = instr.B()
		c = instr.C()
	)
	// arguments
	if b != 0 {
		vm.thread().frame().settop(a+b)
		vm.thread().Call(b-1, c-1)
	} else {
		vm.thread().Call(vm.thread().frame().gettop()-a-1, c-1)
	}
	// returns
	if c--; c > 0 {
		for i, v := range vm.thread().frame().popN(c) {
			vm.thread().frame().set(a+i, v)
		}
	} 
	// C=0 so return values indicated by 'top'
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
func (vm *v53) tailcall(instr vm.Instr) {
	// TODO: proper tail call elimination
	var (
		a = instr.A()
		b = instr.B()
		c = instr.C()
	)
	// arguments
	if b != 0 {
		vm.thread().frame().settop(a+b)
		vm.thread().Call(b-1, c-1)
	} else {
		vm.thread().Call(vm.thread().frame().gettop()-a-1, c-1)
	}
	// returns
	if c--; c > 0 {
		for i, v := range vm.thread().frame().popN(c) {
			vm.thread().frame().set(a+i, v)
		}
	} 
	// C=0 so return values indicated by 'top'
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
func (vm *v53) returns(instr vm.Instr) {
	var (
		a = instr.A()
		b = instr.B()
	)
	if want := vm.thread().frame().rets; want != 0 {
		b--
		var (
			retc int = b
			rets []Value
		)
		if b == -1 {
			retc = vm.thread().frame().gettop()-a
		}
		switch {
			case want > retc: // # wanted > # returned
				for i := a; i < a + retc; i++ {
					rets = append(rets, vm.thread().frame().get(i))
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
					rets = append(rets, vm.thread().frame().get(i))
				}
		}
		vm.thread().frame().caller().pushN(rets)
	}
}

// FORLOOP: Iterate a numeric for loop.
//
// @args A sBx 
//
// R(A)+=R(A+2); if R(A) <?= R(A+1) then { pc+=sBx; R(A+3)=R(A) }                     
func (vm *v53) forloop(instr vm.Instr) {
	var (
		item = vm.thread().frame().get(instr.A())
		upto = vm.thread().frame().get(instr.A()+1)
		step = vm.thread().frame().get(instr.A()+2)
	)
	if isInteger(item) { // integer loop?
		i1 := item.(Int) 
		i2 := upto.(Int)
		i3 := step.(Int)
		i1 += i3 // increment index
		if (i3 > 0 && (i1 <= i2)) || (i3 < 0 && (i1 > i2)) {
			vm.thread().frame().set(instr.A(), i1)   // update internal index...
			vm.thread().frame().set(instr.A()+3, i1) // ... and external index
			vm.thread().frame().step(instr.SBX())    // jump back
		}
	} else { // floating loop
		f1 := item.(Float)
		f2 := upto.(Float)
		f3 := step.(Float)
		f1 += f3
		if (f3 > 0 && (f1 <= f2)) || (f3 < 0 && (f1 > f2)) {
			vm.thread().frame().set(instr.A(), f1)   // update internal index...
			vm.thread().frame().set(instr.A()+3, f1) // ... and external index
			vm.thread().frame().step(instr.SBX())    // jump back
		}
	}
}

// FORPREP: Initialization for a numeric for loop.   
//
// @args A sBx 
//
// R(A)-=R(A+2); pc+=sBx
func (vm *v53) forprep(instr vm.Instr) {
	var (
		init = vm.thread().frame().get(instr.A())
		upto = vm.thread().frame().get(instr.A()+1)
		step = vm.thread().frame().get(instr.A()+2)
	)
	// Try for values as integers.
	var (
		i1, ok1 = toInteger(init)
		i2, ok2 = toInteger(upto)
		i3, ok3 = toInteger(step)
	)
	if ok1 && ok2 && ok3 {
		// TODO: Try converting forlimit to an integer rounding if possible.
		vm.thread().frame().set(instr.A(), i1-i3)
		vm.thread().frame().set(instr.A()+1, i2)
		vm.thread().frame().set(instr.A()+2, i3)
		vm.thread().frame().step(instr.SBX())
		return
	}
	// Try for values as numbers.
	var f1, f2, f3 Float
	if f1, ok1 = toFloat(init); !ok1 {
		vm.thread().errorf("'for' init must be a number")
	}
	if f2, ok2 = toFloat(upto); !ok2 {
		vm.thread().errorf("'for' limit must be a number")
	}
	if f3, ok3 = toFloat(step); !ok3 {
		vm.thread().errorf("'for' step must be a number")
	}
	
	vm.thread().frame().set(instr.A(), f1-f3)
	vm.thread().frame().set(instr.A()+1, f2)
	vm.thread().frame().set(instr.A()+2, f3)
	vm.thread().frame().step(instr.SBX())
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
func (vm *v53) tforcall(instr vm.Instr) {
	// tforcall expects the for variables below to be at a fixed
	// position in the stack for every iteration, so we need to
	// adjust the stack to ensure this to avoid side effects.

	var (
		a = instr.A()
		c = instr.C()
	)

	var (
		iter = vm.thread().frame().get(a)   // iterator function
		data = vm.thread().frame().get(a+1) // state
		ctrl = vm.thread().frame().get(a+2) // control variable / initial value
		base = instr.A()+3
	)

	vm.thread().frame().push(iter)
	vm.thread().frame().push(data)
	vm.thread().frame().push(ctrl)

	vm.thread().Call(2, c)


	rets := vm.thread().frame().popN(c)
	for i, v := range rets {
		vm.thread().frame().set(base+i, v)
	}
}

// TFORLOOP: Initialization for a generic for loop.
//
// @args A sBx 
//
// if R(A+1) ~= nil then { R(A)=R(A+1); pc += sBx } 
func (vm *v53) tforloop(instr vm.Instr) {
	if ctrl := vm.thread().frame().get(instr.A()+1); !IsNone(ctrl) {
		vm.thread().frame().set(instr.A(), ctrl)
		vm.thread().frame().step(instr.SBX())
		return
	}
	// loop done, reset top
	vm.thread().frame().settop(instr.A())
}

// SETLIST: Set a range of array elements for a table.
//
// @args A B C
//
// R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
func (vm *v53) setlist(instr vm.Instr) {
	var (
		a = instr.A()
		b = instr.B()
		c = instr.C()
	)
	if b == 0 { b = vm.thread().frame().gettop() - a - 1}
	o := (c-1) * FieldsPerFlush
	t := vm.thread().frame().get(a).(*table)
	for i := 1; i <= b; i++ {
		t.setInt(int64(o+i), vm.thread().frame().get(a+i))
	}
	vm.thread().frame().popN(b)
}

// CLOSURE: Create a closure of a function prototype.
//
// @args A Bx
// 
// R(A) := closure(KPROTO[Bx])
func (vm *v53) closure(instr vm.Instr) {
	cls := newLuaClosure(vm.prototype(instr.BX()))
	vm.thread().frame().openUp(cls)
	vm.thread().frame().push(cls)
	vm.thread().frame().replace(instr.A())
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
func (vm *v53) vararg(instr vm.Instr) {
	var (
		a = instr.A()
		b = instr.B()
	)
	for i, v := range vm.thread().frame().varargs(b-1) {
		if v == nil { v = None }
		vm.thread().frame().set(a+i, v)	
	}
}

// IDIV: Integer division operator.
//      
// @args A B C
//
// R(A) := RK(B) // RK(C)
func (vm *v53) idiv(instr vm.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.thread().arith(OpQuo, rb, rc)
	)
	vm.thread().frame().set(instr.A(), ra)
}

// BAND: Bit-wise AND operator.        
//
// @args A B C
//
// R(A) := RK(B) & RK(C)
func (vm *v53) band(instr vm.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.thread().arith(OpAnd, rb, rc)
	)
	vm.thread().frame().set(instr.A(), ra)
}

// BOR: Bit-wise OR operator.
//        
// @args A B C
//
// R(A) := RK(B) | RK(C) 
func (vm *v53) bor(instr vm.Instr) {	
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.thread().arith(OpOr, rb, rc)
	)
	vm.thread().frame().set(instr.A(), ra)
}

// BXOR: Bit-wise Exclusive OR operator.
//      
// @args A B C
//
// R(A) := RK(B) ~ RK(C)
func (vm *v53) bxor(instr vm.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.thread().arith(OpXor, rb, rc)
	)
	vm.thread().frame().set(instr.A(), ra)
}

// SHL: Shift bits left.         
//
// @args A B C
//
// R(A) := RK(B) << RK(C) 
func (vm *v53) shl(instr vm.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.thread().arith(OpLsh, rb, rc)
	)
	vm.thread().frame().set(instr.A(), ra)
}

// SHR: Shift bits right.         
//
// @args A B C
//
// R(A) := RK(B) >> RK(C)
func (vm *v53) shr(instr vm.Instr) {
	var (
		rb = vm.rk(instr.B())
		rc = vm.rk(instr.C())
		ra = vm.thread().arith(OpRsh, rb, rc)
	)
	vm.thread().frame().set(instr.A(), ra)
}

// BNOT: Bit-wise NOT operator.        
//
// @args A B
// 
// R(A) := ~R(B)
func (vm *v53) bnot(instr vm.Instr) {
	var (
		rb = vm.rk(instr.B())
		ra = vm.thread().arith(OpNot, rb, None)
	)
	vm.thread().frame().set(instr.A(), ra)
}

// EXTRAARG: Extra (larger) argument for previous opcode.
//
// @args Ax
func (vm *v53) extraarg(instr vm.Instr) {
	// This op func should never execute directly.
	unimplemented(instr.String())
}