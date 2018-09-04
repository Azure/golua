// TODO: loadkx
// TODO: getupval
// TODO: gettable
// TODO: settabup
// TODO: settable
// TODO: setupval
// TODO: newtable
// TODO: self
// TODO: unm
// TODO: not
// TODO: length
// TODO: concat
// TODO: jmp
// TODO: eq
// TODO: lt
// TODO: le
// TODO: test
// TODO: testset
// TODO: call
// TODO: tailcall
// TODO: returns
// TODO: forloop
// TODO: forprep
// TODO: tforcall
// TODO: tforloop
// TODO: setlist
// TODO: vararg
// TODO: bnot
// TODO: extraarg
package lua

import (
	"fmt"
	"os"

	//"github.com/Azure/golua/pkg/luautil"
	//"github.com/Azure/golua/lua/binary"
	"github.com/Azure/golua/lua/ir"
	"github.com/Azure/golua/lua/op"
)

// TODO: remove
var (
	_ = fmt.Println
	_ = os.Exit
)

// MOVE: Copy a value between registers.
//
// @args A B
//
// R(A) := R(B)
func move(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// LOADK: Load a constant into a register.
//
// @args A Bx
//
// R(A) := Kst(Bx)
func loadk(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// LOADKX: Load a constant into a register. The next 'instruction'
// is always EXTRAARG.
//
// @args A
//
// R(A) := Kst(extra arg)
func loadkx(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// LOADBOOL: Load a boolean into a register.
//
// @args A B C
//
// R(A) := (Bool)B; if (C) pc++
func loadbool(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// LOADNIL: Load nil values into a range of registers.
// 
// @args A B
//
// R(A), R(A+1), ..., R(A+B) := nil
func loadnil(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// GETUPVAL: Read an upvalue into a register.
//
// @args A B
//       
// R(A) := UpValue[B]
func getupval(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}
 
// GETTABUP: Read a value from table in up-value into a register.
//
// @args A B C
//   
// R(A) := UpValue[B][RK(C)]
func gettabup(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// GETTABLE: Read a table element into a register.
//
// @args A B C
//
// R(A) := R(B)[RK(C)]
func gettable(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// SETTABUP: Write a register value into table in up-value
//
// @args A B C
//
// UpValue[A][RK(B)] := RK(C) 
func settabup(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// SETTABLE: Write a register value into a table element.
//
// @args A B C
//
// R(A)[RK(B)] := RK(C)
func settable(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// SETUPVAL: Write a register value into an upvalue.
//
// @args A B
//
// UpValue[B] := R(A)            
func setupval(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
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
func newtable(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// SELF: Prepare an object method for calling.  
//
// @args A B C
//
// R(A+1) := R(B); R(A) := R(B)[RK(C)] 
func self(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// ADD: Addition operator.  
//
// @args A B C
//
// R(A) := RK(B) + RK(C)               
func add(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// SUB: Subtraction operator.
//    
// @args A B C
//
// R(A) := RK(B) - RK(C)
func sub(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// MUL: Multiplication operator.
//     
// @args A B C
//    
// R(A) := RK(B) * RK(C)
func mul(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// MOD: Modulus (remainder) operator.
//
// @args A B C
//
// R(A) := RK(B) % RK(C)                
func mod(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// POW: Exponentation operator.
//
// @args A B C
//
// R(A) := RK(B) ^ RK(C)               
func pow(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// DIV: Division operator.         
//
// @args A B C
//
// R(A) := RK(B) / RK(C)
func div(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// UNM: Unary minus.
//
// @args A B
//
// R(A) := -R(B)
func unm(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// NOT: Logical NOT operator.  
//
// @args A B
//
// R(A) := not R(B)
func not(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}
                  
// LEN: Length operator.
//
// @args A B
//
// R(A) := length of R(B) 
func length(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// CONCAT: Concatenate a range of registers.
//
// @args A B C
//
// R(A) := R(B).. ... ..R(C)
func concat(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// JMP: Unconditional jump.
//
// @args A sBx
//
// pc+=sBx; if (A) close all upvalues >= R(A-1)
func jmp(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// EQ: Equality test, with conditional jump.
//          
// @args A B C
//
// if ((RK(B) == RK(C)) ~= A) then pc++
func eq(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// LT: Less than test, with conditional jump.
//
// @args A B C
//
// if ((RK(B) <  RK(C)) ~= A) then pc++
func lt(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// LE: Less than or equal to test, with conditional jump.
//
// @args A B C
//
// if ((RK(B) <= RK(C)) ~= A) then pc++
func le(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// TEST: Boolean test, with conditional jump.    
//
// @args A C 
//
// if not (R(A) <=> C) then pc++          
func test(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// TESTSET: Boolean test, with conditional jump and assignment.
//
// @args A B C
//
// if (R(B) <=> C) then R(A) := R(B) else pc++
func testset(frame *Frame, instr ir.Instr) {
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
func call(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// TAILCALL: Perform a tail call.   
//
// @args A B C
//
// return R(A)(R(A+1), ... ,R(A+B-1))
func tailcall(frame *Frame, instr ir.Instr) {
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
// @args A B
//
// If (B == 0) then return up to 'top'.
// 
// return R(A), ... ,R(A+B-2)
func returns(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// FORLOOP: Iterate a numeric for loop.
//
// @args A sBx 
//
// R(A)+=R(A+2); if R(A) <?= R(A+1) then { pc+=sBx; R(A+3)=R(A) }                     
func forloop(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// FORPREP: Initialization for a numeric for loop.   
//
// @args A sBx 
//
// R(A)-=R(A+2); pc+=sBx
func forprep(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// TFORCALL: Iterate a generic for loop.
//
// @args A C
//
// R(A+3), ... ,R(A+2+C) := R(A)(R(A+1), R(A+2))
func tforcall(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// TFORLOOP: Initialization for a generic for loop.
//
// @args A sBx 
//
// if R(A+1) ~= nil then { R(A)=R(A+1); pc += sBx } 
func tforloop(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// SETLIST: Set a range of array elements for a table.
//
// @args A B C
//
// R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
func setlist(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// CLOSURE: Create a closure of a function prototype.
//
// @args A Bx
// 
// R(A) := closure(KPROTO[Bx])
func closure(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// VARARG: Assign vararg function arguments to register.  
//
// @args A B
//
// R(A), R(A+1), ..., R(A+B-2) = vararg
func vararg(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// IDIV: Integer division operator.
//      
// @args A B C
//
// R(A) := RK(B) // RK(C)
func idiv(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// BAND: Bit-wise AND operator.        
//
// @args A B C
//
// R(A) := RK(B) & RK(C)
func band(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// BOR: Bit-wise OR operator.
//        
// @args A B C
//
// R(A) := RK(B) | RK(C) 
func bor(frame *Frame, instr ir.Instr) {	
	unimplemented(instr.String())
}

// BXOR: Bit-wise Exclusive OR operator.
//      
// @args A B C
//
// R(A) := RK(B) ~ RK(C)
func bxor(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// SHL: Shift bits left.         
//
// @args A B C
//
// R(A) := RK(B) << RK(C) 
func shl(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// SHR: Shift bits right.         
//
// @args A B C
//
// R(A) := RK(B) >> RK(C)
func shr(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// BNOT: Bit-wise NOT operator.        
//
// @args A B
// 
// R(A) := ~R(B)
func bnot(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// EXTRAARG: Extra (larger) argument for previous opcode.
//
// @args Ax
func extraarg(frame *Frame, instr ir.Instr) {
	unimplemented(instr.String())
}

// cmd is an execution state corresponding to a lua opcode.
type cmd func(*Frame, ir.Instr) (cmd, ir.Instr)

// ops is a table of lua opcode commands.
var ops []cmd

func init() {
	ops = []cmd{
		op.MOVE: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			move(frame, instr)
			return frame.fetch()
		},
		op.LOADK: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			loadk(frame, instr)
			return frame.fetch()
		},
		op.LOADKX: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			loadkx(frame, instr)
			return frame.fetch()
		},
		op.LOADBOOL: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			loadbool(frame, instr)
			return frame.fetch()
		},
		op.LOADNIL: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			loadnil(frame, instr)
			return frame.fetch()
		},
		op.GETUPVAL: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			getupval(frame, instr)
			return frame.fetch()
		},
		op.GETTABUP: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			gettabup(frame, instr)
			return frame.fetch()
		},
		op.GETTABLE: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			gettable(frame, instr)
			return frame.fetch()
		},
		op.SETTABUP: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			settabup(frame, instr)
			return frame.fetch()
		},
		op.SETUPVAL: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			setupval(frame, instr)
			return frame.fetch()
		},
		op.SETTABLE: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			settable(frame, instr)
			return frame.fetch()
		},
		op.NEWTABLE: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			newtable(frame, instr)
			return frame.fetch()
		},
		op.SELF: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			self(frame, instr)
			return frame.fetch()
		},
		op.ADD: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			add(frame, instr)
			return frame.fetch()
		},
		op.SUB: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			sub(frame, instr)
			return frame.fetch()
		},
		op.MUL: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			mul(frame, instr)
			return frame.fetch()
		},
		op.MOD: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			mod(frame, instr)
			return frame.fetch()
		},
		op.POW: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			pow(frame, instr)
			return frame.fetch()
		},
		op.DIV: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			div(frame, instr)
			return frame.fetch()
		},
		op.IDIV: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			idiv(frame, instr)
			return frame.fetch()
		},
		op.BAND: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			band(frame, instr)
			return frame.fetch()
		},
		op.BOR: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			bor(frame, instr)
			return frame.fetch()
		},
		op.BXOR: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			bxor(frame, instr)
			return frame.fetch()
		},
		op.SHL: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			shl(frame, instr)
			return frame.fetch()
		},
		op.SHR: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			shr(frame, instr)
			return frame.fetch()
		},
		op.UNM: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			unm(frame, instr)
			return frame.fetch()
		},
		op.BNOT: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			bnot(frame, instr)
			return frame.fetch()
		},
		op.NOT: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			not(frame, instr)
			return frame.fetch()
		},
		op.LEN: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			length(frame, instr)
			return frame.fetch()
		},
		op.CONCAT: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			concat(frame, instr)
			return frame.fetch()
		},
		op.JMP: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			jmp(frame, instr)
			return frame.fetch()
		},
		op.EQ: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			eq(frame, instr)
			return frame.fetch()
		},
		op.LT: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			lt(frame, instr)
			return frame.fetch()
		},
		op.LE: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			le(frame, instr)
			return frame.fetch()
		},
		op.TEST: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			test(frame, instr)
			return frame.fetch()
		},
		op.TESTSET: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			testset(frame, instr)
			return frame.fetch()
		},
		op.CALL: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			call(frame, instr)
			return frame.fetch()
		},
		op.TAILCALL: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			tailcall(frame, instr)
			return frame.fetch()
		},
		op.RETURN: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			returns(frame, instr)
			return nil, instr
		},
		op.FORLOOP: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			forloop(frame, instr)
			return frame.fetch()
		},
		op.FORPREP: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			forprep(frame, instr)
			return frame.fetch()
		},
		op.TFORCALL: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			tforcall(frame, instr)
			return frame.fetch()
		},
		op.TFORLOOP: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			tforloop(frame, instr)
			return frame.fetch()
		},
		op.SETLIST: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			setlist(frame, instr)
			return frame.fetch()
		},
		op.CLOSURE: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			closure(frame, instr)
			return frame.fetch()
		},
		op.VARARG: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			vararg(frame, instr)
			return frame.fetch()
		},
		op.EXTRAARG: func(frame *Frame, instr ir.Instr) (cmd, ir.Instr) {
			extraarg(frame, instr)
			return frame.fetch()
		},
	}
}