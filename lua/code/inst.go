package code

import (
	"fmt"
)

const (
	sizeAx = sizeC + sizeB + sizeA
	sizeBx = sizeC + sizeB
	sizeOp = 6
	sizeA  = 8
	sizeB  = 9
	sizeC  = 9

	posOp  = 0
	posAx  = posA
	posBx  = posC
	posA   = posOp + sizeOp
	posB   = posC + sizeC
	posC   = posA + sizeA

	// bit mask (1 = constant, 0 = register)
	bitRK = 1 << (sizeB - 1)	
)

type Instr uint32

func (instr Instr) Code() Opcode { return Opcode(int(instr&0x3F)) }
func (instr Instr) Mode() Mode { return instr.Code().Mode() }
func (instr Instr) ABC() (a, b, c int) { return instr.A(), instr.B(), instr.C() }
func (instr Instr) A() (a int) { return int(instr >> 6 & 0xFF) }
func (instr Instr) B() (b int) { return int(instr >> 23 & 0x1FF) }
func (instr Instr) C() (c int) { return int(instr >> 14 & 0x1FF) }
func (instr Instr) AX() (ax int) { return int(instr >> 6) }
func (instr Instr) BX() (bx int) { return int(instr >> 14) }
func (instr Instr) SBX() (sbx int) { return instr.BX() - MaxArgSBX }

func (instr *Instr) SetOp(op Opcode) { setarg(instr, uint(op), posOp, sizeOp) }
func (instr *Instr) SetA(a uint) { setarg(instr, a, posA, sizeA) }
func (instr *Instr) SetB(b uint) { setarg(instr, b, posB, sizeB) }
func (instr *Instr) SetC(c uint) { setarg(instr, c, posC, sizeC) }
func (instr *Instr) SetAX(ax uint) { setarg(instr, ax, posAx, sizeAx) }
func (instr *Instr) SetBX(bx uint) { setarg(instr, bx, posBx, sizeBx) }
func (instr *Instr) SetSBX(sbx int) { instr.SetBX(uint(sbx + MaxArgSBX)) }

func (instr Instr) String() string {
	return fmt.Sprintf("%s %s", instr.Code(), args(instr))
}

func MakeAsBx(op Opcode, a, sBx int) Instr {
	return MakeABx(op, a, sBx+MaxArgSBX)
}

func MakeABC(op Opcode, a, b, c int) Instr {
	return Instr(int32(op)<<posOp|int32(a)<<posA|int32(b)<<posB|int32(c)<<posC)
}

func MakeABx(op Opcode, a, bx int) Instr {
	return Instr(int32(op)<<posOp|int32(a)<<posA|int32(bx)<<posBx)
}

func MakeAx(op Opcode, ax int) Instr {
	return Instr(int32(op)<<posOp|int32(ax)<<posAx)
}

func setarg(i *Instr, v, pos, size uint) {
	*i = Instr(uint32(*i) & mask0(size, pos) | uint32(v << pos) & mask1(size, pos))
}

func mask1(n, p uint) uint32{ return ((^((^uint32(0))<<n))<<p) }
func mask0(n, p uint) uint32 { return ^(mask1(n, p)) }