package lua

import (
	"github.com/Azure/golua/lua/ir"
)

func (fr *Frame) instruction(pc int) ir.Instr {
	if fr.closure.proto == nil {
		return ir.Instr(0)
	}
	return ir.Instr(fr.closure.proto.Code[pc])
}

func (fr *Frame) constant(index int) Value {
	if fr.closure.proto == nil {
		return nil
	}
	return ValueOf(fr.closure.proto.Consts[index])
}

func (fr *Frame) recover(err *error) {
	if r := recover(); r != nil {
		if e, ok := r.(error); ok {
			*err = e
		}
	}
}

func (fr *Frame) fetch() (cmd, ir.Instr) {
	instr := fr.instruction(fr.pc)
	fr.pc++
	return ops[instr.Code()], instr
}

func (fr *Frame) exec() {
	for cmd, instr := fr.fetch(); cmd != nil; {
		cmd, instr = cmd(fr, instr)
	}
}