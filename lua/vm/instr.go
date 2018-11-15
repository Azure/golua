package vm

import (
	"fmt"
)

const (
	MaxArgBX  = 1<<18 - 1
	MaxArgSBX = MaxArgBX >> 1
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
	OpArgU                  // argument is used
	OpArgR                  // argument is a register or a jump offset
	OpArgK                  // argument is a constant or register/constant
)

type Instr uint32

func (instr Instr) Code() Code { return Code(int(instr & 0x3F)) }

func (instr Instr) ABC() (a, b, c int) { return instr.A(), instr.B(), instr.C() }

func (instr Instr) A() (a int) { return int(instr >> 6 & 0xFF) }

func (instr Instr) B() (b int) { return int(instr >> 23 & 0x1FF) }

func (instr Instr) C() (c int) { return int(instr >> 14 & 0x1FF) }

func (instr Instr) AX() (ax int) { return int(instr >> 6) }

func (instr Instr) BX() (bx int) { return int(instr >> 14) }

func (instr Instr) SBX() (sbx int) { return instr.BX() - MaxArgSBX }

func (instr Instr) String() string {
	return fmt.Sprintf("%s %s", instr.Code(), args(instr))
}

func mask1(n, p uint) Instr { return ((^((^Instr(0)) << n)) << p) }
func mask0(n, p uint) Instr { return ^(mask1(n, p)) }

func indexk(x int) int { return x &^ (1 << 8) }
func myk(x int) int    { return -1 - x }
func isk(x int) bool   { return x&(1<<8) != 0 }

func args(instr Instr) string {
	switch code := instr.Code(); code.Mode() {
	case ModeABC:
		return fmt.Sprintf("A=%d B=%d C=%d", instr.A(), instr.B(), instr.C())
	case ModeABx:
		return fmt.Sprintf("A=%d BX=%d", instr.A(), instr.BX())
	case ModeAsBx:
		return fmt.Sprintf("A=%d SBX=%d", instr.A(), instr.SBX())
	case ModeAx:
		return fmt.Sprintf("AX=%d", instr.AX())
	}
	panic(fmt.Sprintf("ir: unknown op mode: %d", instr.Code().Mode()))
}
