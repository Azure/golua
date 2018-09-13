package ir

import (
	"fmt"

	"github.com/Azure/golua/lua/op"
)

const (
	MaxArgBX  = 1<<18-1
	MaxArgSBX = MaxArgBX>>1
)

// Masks for instruction properties. The format is:
// bits 0-1: op mode
// bits 2-3: C arg mode
// bits 4-5: B arg mode
//    bit 6: instruction set register A
//    bit 7: operator is a test (next instruction must be a jump)
type OpArgMask uint8

const (
	OpArgN OpArgMask = iota // argument is not used
	OpArgU 					// argument is used
	OpArgR 					// argument is a register or a jump offset
	OpArgK 					// argument is a constant or register/constant
)

type Instr uint32

func (instr Instr) Code() op.Code {
	//return op.Code((instr >> 0) & mask1(6, 0))
	return op.Code(int(instr&0x3F))
}

func (instr Instr) ABC() (a, b, c int) { return instr.A(), instr.B(), instr.C() }

func (instr Instr) A() (a int) { return int(instr >> 6 & 0xFF) }

func (instr Instr) B() (b int) { return int(instr >> 23 & 0x1FF) }

func (instr Instr) C() (c int) { return int(instr >> 14 & 0x1FF) }

func (instr Instr) AX() (ax int) { return int(instr >> 6) }

func (instr Instr) BX() (bx int) { return int(instr >> 14) }

func (instr Instr) SBX() (sbx int) { return instr.BX() - MaxArgSBX }

func (instr Instr) String() string {
	return fmt.Sprintf("%s %s", instr.Code(), args(instr))
	// switch instr.Code() {
	// 	case op.MOVE: // R(A) := R(B)
	// 	case op.LOADK: // R(A) := Kst(Bx)
	// 	case op.LOADKX: // R(A) := Kst(extra arg)
	// 	case op.LOADBOOL: // R(A) := (Bool)B; if (C) pc++
	// 	case op.LOADNIL: //R(A), R(A+1), ..., R(A+B) := nil
	// 	case op.GETUPVAL: // R(A) := UpValue[B]
	// 	case op.GETTABUP: // R(A) := UpValue[B][RK(C)]
	// 	case op.GETTABLE: // R(A) := R(B)[RK(C)]
	// 	case op.SETTABUP: // UpValue[A][RK(B)] := RK(C) 
	// 	case op.SETTABLE: // R(A)[RK(B)] := RK(C)
	// 	case op.SETUPVAL: // UpValue[B] := R(A)
	// 	case op.NEWTABLE: // R(A) := {} (size = B,C)
	// 	case op.SELF: // R(A+1) := R(B); R(A) := R(B)[RK(C)] 
	// 	case op.ADD: // R(A) := RK(B) + RK(C)               
	// 	case op.SUB: // R(A) := RK(B) - RK(C)
	// 	case op.MUL: // R(A) := RK(B) * RK(C)
	// 	case op.MOD: // R(A) := RK(B) % RK(C)
	// 	case op.POW: // R(A) := RK(B) ^ RK(C)               
	// 	case op.DIV: // R(A) := RK(B) / RK(C)
	// 	case op.UNM: // R(A) := -R(B)
	// 	case op.NOT: // R(A) := not R(B)
	// 	case op.LEN: // R(A) := length of R(B) 
	// 	case op.CONCAT: // R(A) := R(B).. ... ..R(C)
	// 	case op.JMP: // pc+=sBx; if (A) close all upvalues >= R(A-1)
	// 	case op.EQ: // if ((RK(B) == RK(C)) ~= A) then pc++
	// 	case op.LT: // if ((RK(B) <  RK(C)) ~= A) then pc++
	// 	case op.LE: // if ((RK(B) <= RK(C)) ~= A) then pc++
	// 	case op.TEST: // if not (R(A) <=> C) then pc++ 
	// 	case op.TESTSET: // if (R(B) <=> C) then R(A) := R(B) else pc++
	// 	case op.CALL: // R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))
	// 	case op.TAILCALL: // return R(A)(R(A+1), ... ,R(A+B-1))
	// 	case op.RETURN: // return R(A), ... ,R(A+B-2)
	// 	case op.FORLOOP: // R(A)+=R(A+2); if R(A) <?= R(A+1) then { pc+=sBx; R(A+3)=R(A) }                     
	// 	case op.FORPREP: // R(A)-=R(A+2); pc+=sBx
	// 	case op.TFORCALL: // R(A+3), ... ,R(A+2+C) := R(A)(R(A+1), R(A+2))
	// 	case op.TFORLOOP: // if R(A+1) ~= nil then { R(A)=R(A+1); pc += sBx } 
	// 	case op.SETLIST: // R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
	// 	case op.CLOSURE: // R(A) := closure(KPROTO[Bx])
	// 	case op.VARARG: // R(A), R(A+1), ..., R(A+B-2) = vararg
	// 	case op.IDIV: // R(A) := RK(B) // RK(C)
	// 	case op.BAND: // R(A) := RK(B) & RK(C)
	// 	case op.BOR: // R(A) := RK(B) | RK(C) 
	// 	case op.BXOR: // R(A) := RK(B) ~ RK(C)
	// 	case op.SHL: // R(A) := RK(B) << RK(C) 
	// 	case op.SHR: // R(A) := RK(B) >> RK(C)
	// 	case op.BNOT: // R(A) := ~R(B)
	// 	default:
	// 		return fmt.Sprintf("ir: unknown opcode %d", instr.Code())
	// }
}

func mask1(n, p uint) Instr { return ((^((^Instr(0))<<n))<<p) }
func mask0(n, p uint) Instr { return ^(mask1(n, p)) }

func indexk(x int) int { return x&^(1<<8) }
func myk(x int) int { return -1 - x }
func isk(x int) bool { return x&(1<<8)!=0 }

func args(instr Instr) string {
	switch code := instr.Code(); code.Mode() {
		case op.ModeABC:
			return fmt.Sprintf("A=%d B=%d C=%d", instr.A(), instr.B(), instr.C())
		case op.ModeABx:
			return fmt.Sprintf("A=%d BX=%d", instr.A(), instr.BX())
		case op.ModeAsBx:
			return fmt.Sprintf("A=%d SBX=%d", instr.A(), instr.SBX())
		case op.ModeAx:
			return fmt.Sprintf("AX=%d", instr.AX())
	}
	panic(fmt.Sprintf("ir: unknown op mode: %d", instr.Code().Mode()))
}