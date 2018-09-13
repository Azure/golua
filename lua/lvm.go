// # TODO
// not
// loadbool
// lt
// le
// eq
// jmp

// concat
// test
// testset
// loadkx
// loadnil
// gettable
// settable
// newtable

// getupval
// setupval
// self
// tailcall
// forloop
// forprep
// tforcall
// tforloop
// setlist
// vararg
// extraarg
package lua

import (
	//"github.com/Azure/golua/pkg/luautil"
	"github.com/Azure/golua/lua/ir"
)

//
// Implementation of Lua v53 Opcodes
//

// MOVE: Copy a value between registers.
//
// @args A B
//
// R(A) := R(B)
func (vm *v53) move(instr ir.Instr) {
	rb := vm.frame().get(instr.B()+1)
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
	unimplemented(instr.String())
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
	unimplemented(instr.String())
}

// GETUPVAL: Read an upvalue into a register.
//
// @args A B
//       
// R(A) := UpValue[B]
func (vm *v53) getupval(instr ir.Instr) {
	unimplemented(instr.String())
}

// SETUPVAL: Write a register value into an upvalue.
//
// @args A B
//
// UpValue[B] := R(A)
func (vm *v53) setupval(instr ir.Instr) {
	unimplemented(instr.String())
}

// GETTABLE: Read a table element into a register (locals).
//
// @args A B C
//
// R(A) := R(B)[RK(C)]
func (vm *v53) gettable(instr ir.Instr) {
	unimplemented(instr.String())
}

// SETTABLE: Write a register value into a table element (locals).
//
// @args A B C
//
// R(A)[RK(B)] := RK(C)
func (vm *v53) settable(instr ir.Instr) {
	unimplemented(instr.String())
}

// GETTABUP: Read a value from table in up-value into a register (globals).
//
// @args A B C
//   
// R(A) := UpValue[B][RK(C)]
func (vm *v53) gettabup(instr ir.Instr) {
	up := vm.frame().upvalue(instr.B())
	rc := vm.rk(instr.C())
	ra := vm.State.gettable(up, rc, 0)
	vm.frame().set(instr.A(), ra)
}

// SETTABUP: Write a register value into table in up-value (globals).
//
// @args A B C
//
// UpValue[A][RK(B)] := RK(C) 
func (vm *v53) settabup(instr ir.Instr) {
	up := vm.frame().upvalue(instr.A())
	rb := vm.rk(instr.B())
	rc := vm.rk(instr.C())
	vm.State.settable(up, rb, rc, 0)
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
	unimplemented(instr.String())
}

// SELF: Prepare an object method for calling.  
//
// @args A B C
//
// R(A+1) := R(B); R(A) := R(B)[RK(C)] 
func (vm *v53) self(instr ir.Instr) {
	unimplemented(instr.String())
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
	unimplemented(instr.String())
}

// CONCAT: Concatenate a range of registers.
//
// @args A B C
//
// R(A) := R(B).. ... ..R(C)
func (vm *v53) concat(instr ir.Instr) {
	unimplemented(instr.String())
}

// JMP: Unconditional jump.
//
// @args A sBx
//
// pc+=sBx; if (A) close all upvalues >= R(A-1)
func (vm *v53) jmp(instr ir.Instr) {
	vm.frame().step(instr.SBX())
	if a := instr.A(); a != 0 {
		var (
			ra = vm.frame().get(a-1)
			ls = vm.State
		)
		vm.frame().closure.closeUpValues(ls, ra)
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
		aa = instr.A() == 1
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
		aa = instr.A() == 1
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
		aa = instr.A() == 1
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
	unimplemented(instr.String())
}

// TESTSET: Boolean test, with conditional jump and assignment.
//
// @args A B C
//
// if (R(B) <=> C) then R(A) := R(B) else pc++
func (vm *v53) testset(instr ir.Instr) {
	unimplemented(instr.String())
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
// @args A B C
//
// R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))  
func (vm *v53) call(instr ir.Instr) {
	var (
		//a = instr.A()
		b = instr.B()
		c = instr.C()
	)

	// OPEN THE UPVALUES

	if b != 0 {
		vm.State.Call(b-1, c-1)
		return
	}
	unimplemented("call: b == 0")
}

// TAILCALL: Perform a tail call.   
//
// @args A B C
//
// return R(A)(R(A+1), ... ,R(A+B-1))
func (vm *v53) tailcall(instr ir.Instr) {
	unimplemented(instr.String())
}

// RETURN: Returns from function call.
//
// Returns to the calling function, with optional return values. First, op.RETURN closes any
// open upvalues by calling state.Close().
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
	if b > 0 {
		rets := vm.frame().popN(a+(b-1))
		vm.frame().caller().pushN(rets)
		return
	}
	unimplemented("returns: b == 0")
}

// FORLOOP: Iterate a numeric for loop.
//
// @args A sBx 
//
// R(A)+=R(A+2); if R(A) <?= R(A+1) then { pc+=sBx; R(A+3)=R(A) }                     
func (vm *v53) forloop(instr ir.Instr) {
	unimplemented(instr.String())
}

// FORPREP: Initialization for a numeric for loop.   
//
// @args A sBx 
//
// R(A)-=R(A+2); pc+=sBx
func (vm *v53) forprep(instr ir.Instr) {
	unimplemented(instr.String())
}

// TFORCALL: Iterate a generic for loop.
//
// @args A C
//
// R(A+3), ... ,R(A+2+C) := R(A)(R(A+1), R(A+2))
func (vm *v53) tforcall(instr ir.Instr) {
	unimplemented(instr.String())
}

// TFORLOOP: Initialization for a generic for loop.
//
// @args A sBx 
//
// if R(A+1) ~= nil then { R(A)=R(A+1); pc += sBx } 
func (vm *v53) tforloop(instr ir.Instr) {
	unimplemented(instr.String())
}

// SETLIST: Set a range of array elements for a table.
//
// @args A B C
//
// R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
func (vm *v53) setlist(instr ir.Instr) {
	unimplemented(instr.String())
}

// CLOSURE: Create a closure of a function prototype.
//
// @args A Bx
// 
// R(A) := closure(KPROTO[Bx])
func (vm *v53) closure(instr ir.Instr) {
	cls := newLuaClosure(vm.prototype(instr.BX()))
	cls.openUpValues(vm.State)
	vm.frame().push(cls)
	vm.frame().replace(instr.A()+1)
}

// VARARG: Assign vararg function arguments to register.  
//
// @args A B
//
// R(A), R(A+1), ..., R(A+B-2) = vararg
func (vm *v53) vararg(instr ir.Instr) {
	unimplemented(instr.String())
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
	unimplemented(instr.String())
}