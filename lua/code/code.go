package code

import (
	"fmt"
)

type Error string
func (e Error) Error() string { return string(e) }

const (
	MaxIndexRK = bitRK - 1
	MaxArgAX   = (1 << sizeAx) - 1
	MaxArgBX   = 1<<18-1
	MaxArgSBX  = MaxArgBX>>1
	MaxArgA    = (1 << sizeA) - 1
	MaxArgB    = (1 << sizeB) - 1
	MaxArgC    = (1 << sizeC) - 1
	NoReg 	   = MaxArgA // invalid register that fits in 8-bits
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

type Mode uint8

const (
    ModeABC Mode = iota
    ModeABx
    ModeAsBx
    ModeAx
)

var modes = [...]string{
	ModeABC:  "iABC",
	ModeABx:  "iABx",
	ModeAsBx: "iAsBx",
	ModeAx:   "iAx",
}

func (mode Mode) String() string { return modes[mode] }

type Opcode uint8

const (
	MOVE	 Opcode = iota
	LOADK
	LOADKX
	LOADBOOL
	LOADNIL
	GETUPVAL
	GETTABUP
	GETTABLE
	SETTABUP
	SETUPVAL
	SETTABLE
	NEWTABLE
	SELF
	ADD
	SUB
	MUL
	MOD
	POW
	DIV
	IDIV
	BAND
	BOR
	BXOR
	SHL
	SHR
	UNM
	BNOT
	NOT
	LEN
	CONCAT
	JMP
	EQ
	LT
	LE
	TEST
	TESTSET
	CALL
	TAILCALL
	RETURN
	FORLOOP
	FORPREP
	TFORCALL
	TFORLOOP
	SETLIST
	CLOSURE
	VARARG
	EXTRAARG
)

var names = [...]string{
	MOVE: 	  "MOVE",
	LOADK:    "LOADK",
	LOADKX:   "LOADKX",
	LOADBOOL: "LOADBOOL",
	LOADNIL:  "LOADNIL",
	GETUPVAL: "GETUPVAL",
	GETTABUP: "GETTABUP",
	GETTABLE: "GETTABLE",
	SETTABUP: "SETTABUP",
	SETUPVAL: "SETUPVAL",
	SETTABLE: "SETTABLE",
	NEWTABLE: "NEWTABLE",
	SELF:     "SELF",
	ADD:      "ADD",
	SUB:      "SUB",
	MUL:      "MUL",
	MOD:      "MOD",
	POW:      "POW",
	DIV:      "DIV",
	IDIV:     "IDIV",
	BAND:     "BAND",
	BOR:      "BOR",
	BXOR:     "BXOR",
	SHL:      "SHL",
	SHR:      "SHR",
	UNM:      "UNM",
	BNOT:     "BNOT",
	NOT:      "NOT",
	LEN:      "LEN",
	CONCAT:   "CONCAT",
	JMP:      "JMP",
	EQ: 	  "EQ",
	LT: 	  "LT",
	LE: 	  "LE",
	TEST: 	  "TEST",
	TESTSET:  "TESTSET",
	CALL: 	  "CALL",
	TAILCALL: "TAILCALL",
	RETURN:   "RETURN",
	FORLOOP:  "FORLOOP",
	FORPREP:  "FORPREP",
	TFORCALL: "TFORCALL",
	TFORLOOP: "TFORLOOP",
	SETLIST:  "SETLIST",
	CLOSURE:  "CLOSURE",
	VARARG:   "VARARG",
	EXTRAARG: "EXTRAARG",
}

func (op Opcode) Mask() Mask { return masks[op] }
func (op Opcode) Mode() Mode { return masks[op].Mode() }
func (op Opcode) String() string { return names[op] }

// Masks for instruction arguments.
type ArgMask uint8

const (
	ArgN ArgMask = iota // argument is not used
	ArgU 				// argument is used
	ArgR 				// argument is a register or a jump offset
	ArgK 				// argument is a constant or register/constant
)

// ArgMask represents masks for instruction properties.
//
// The format is:
// bits 0-1: op mode
// bits 2-3: C arg mode
// bits 4-5: B arg mode
// bit 6: instruction set register A
// bit 7: operator is a test (next instruction must be a jump) 
type Mask uint8

var masks = [...]Mask{
    MOVE:     mask(0, 1, ArgR, ArgN, ModeABC),     // MOVE
    LOADK:    mask(0, 1, ArgK, ArgN, ModeABx),     // LOADK
    LOADKX:   mask(0, 1, ArgN, ArgN, ModeABx),     // LOADKX
    LOADBOOL: mask(0, 1, ArgU, ArgU, ModeABC),     // LOADBOOL
    LOADNIL:  mask(0, 1, ArgU, ArgN, ModeABC),     // LOADNIL
    GETUPVAL: mask(0, 1, ArgU, ArgN, ModeABC),     // GETUPVAL
    GETTABUP: mask(0, 1, ArgU, ArgK, ModeABC),     // GETTABUP
    GETTABLE: mask(0, 1, ArgR, ArgK, ModeABC),     // GETTABLE
    SETTABUP: mask(0, 0, ArgK, ArgK, ModeABC),     // SETTABUP
    SETUPVAL: mask(0, 0, ArgU, ArgN, ModeABC),     // SETUPVAL
    SETTABLE: mask(0, 0, ArgK, ArgK, ModeABC),     // SETTABLE
    NEWTABLE: mask(0, 1, ArgU, ArgU, ModeABC),     // NEWTABLE
    SELF:     mask(0, 1, ArgR, ArgK, ModeABC),     // SELF
    ADD:      mask(0, 1, ArgK, ArgK, ModeABC),     // ADD
    SUB:      mask(0, 1, ArgK, ArgK, ModeABC),     // SUB
    MUL:      mask(0, 1, ArgK, ArgK, ModeABC),     // MUL
    MOD:      mask(0, 1, ArgK, ArgK, ModeABC),     // MOD
    POW:      mask(0, 1, ArgK, ArgK, ModeABC),     // POW
    DIV:      mask(0, 1, ArgK, ArgK, ModeABC),     // DIV
    IDIV:     mask(0, 1, ArgK, ArgK, ModeABC),     // IDIV
    BAND:     mask(0, 1, ArgK, ArgK, ModeABC),     // BAND
    BOR:      mask(0, 1, ArgK, ArgK, ModeABC),     // BOR
    BXOR:     mask(0, 1, ArgK, ArgK, ModeABC),     // BXOR
    SHL:      mask(0, 1, ArgK, ArgK, ModeABC),     // SHL
    SHR:      mask(0, 1, ArgK, ArgK, ModeABC),     // SHR
    UNM:      mask(0, 1, ArgR, ArgN, ModeABC),     // UNM
    BNOT:     mask(0, 1, ArgR, ArgN, ModeABC),     // BNOT
    NOT:      mask(0, 1, ArgR, ArgN, ModeABC),     // NOT
    LEN:      mask(0, 1, ArgR, ArgN, ModeABC),     // LEN
    CONCAT:   mask(0, 1, ArgR, ArgR, ModeABC),     // CONCAT
    JMP:      mask(0, 0, ArgR, ArgN, ModeAsBx),    // JMP
    EQ:       mask(1, 0, ArgK, ArgK, ModeABC),     // EQ
    LT:       mask(1, 0, ArgK, ArgK, ModeABC),     // LT
    LE:       mask(1, 0, ArgK, ArgK, ModeABC),     // LE
    TEST:     mask(1, 0, ArgN, ArgU, ModeABC),     // TEST
    TESTSET:  mask(1, 1, ArgR, ArgU, ModeABC),     // TESTSET
    CALL:     mask(0, 1, ArgU, ArgU, ModeABC),     // CALL
    TAILCALL: mask(0, 1, ArgU, ArgU, ModeABC),     // TAILCALL
    RETURN:   mask(0, 0, ArgU, ArgN, ModeABC),     // RETURN
    FORLOOP:  mask(0, 1, ArgR, ArgN, ModeAsBx),    // FORLOOP
    FORPREP:  mask(0, 1, ArgR, ArgN, ModeAsBx),    // FORPREP
    TFORCALL: mask(0, 0, ArgN, ArgU, ModeABC),     // TFORCALL
    TFORLOOP: mask(0, 1, ArgR, ArgN, ModeAsBx),    // TFORLOOP
    SETLIST:  mask(0, 0, ArgU, ArgU, ModeABC),     // SETLIST
    CLOSURE:  mask(0, 1, ArgU, ArgN, ModeABx),     // CLOSURE
    VARARG:   mask(0, 1, ArgU, ArgN, ModeABC),     // VARARG
    EXTRAARG: mask(0, 0, ArgU, ArgU, ModeAx),      // EXTRAARG
}

func (mask Mask) B(barg ArgMask) bool { return barg == ArgMask((mask>>4)&3) }
func (mask Mask) C(carg ArgMask) bool { return carg == ArgMask((mask>>2)&3) }
func (mask Mask) Mode() Mode { return Mode(mask & 3) }
func (mask Mask) SetA() bool { return mask&(1<<6) != 1 }
func (mask Mask) Test() bool { return mask&(1<<7) != 1 }

func mask(t, a uint8, b, c ArgMask, m Mode) Mask {
    return Mask((((t)<<7) | ((a)<<6) | ((uint8(b))<<4) | ((uint8(c))<<2) | (uint8(m))))
}

func Format(instr Instr) string {
	return fmt.Sprintf("%-10s\t%s", instr.Code(), args(instr))
}

func args(instr Instr) (s string) {
	switch code := instr.Code(); code.Mode() {
		case ModeABC:
			var ( a, b, c = instr.A(), instr.B(), instr.C() )

			if !code.Mask().B(ArgN) {
				if IsKst(instr.B()) {
					b = Kst(ToKst(b))
				}
			}
			if !code.Mask().C(ArgN) {
				if IsKst(instr.C()) {
					c = Kst(ToKst(c))
				}
			}
			// return fmt.Sprintf("A=%d B=%d C=%d", a, b, c)
			return fmt.Sprintf("%d %d %d", a, b, c)

		case ModeABx:
			var ( a, bx = instr.A(), instr.BX() )

			if code.Mask().B(ArgK) {
				s += fmt.Sprintf("%d", Kst(bx))
			}
			if code.Mask().B(ArgU) {
				s += fmt.Sprintf(" %d", bx)
			}		
			// return fmt.Sprintf("A=%d %s", a, s)
 			return fmt.Sprintf("%d %s", a, s)

		case ModeAsBx:
			// return fmt.Sprintf("A=%d SBX=%d", instr.A(), instr.SBX())
			return fmt.Sprintf("%d %d", instr.A(), instr.SBX())

		case ModeAx:
			// return fmt.Sprintf("AX=%d", myk(instr.AX()))
			return fmt.Sprintf("%d", Kst(instr.AX()))
	}
	panic(fmt.Sprintf("ir: unknown op mode: %d", instr.Code().Mode()))
}

// IsKst reports whether the register index is a constant.
func IsKst(reg int) bool { return reg&bitRK != 0 }

// ToConstIndex returns the unmasked constant index.
func ToKst(reg int) int { return reg&^bitRK }

// Kst returns the relative constant index of k.
func Kst(k int) int { return -1-k }

// RK returns the constant index reg as an RK value.
func RK(reg int) int { return reg|bitRK }