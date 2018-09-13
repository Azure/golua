package lua

import (
	"fmt"
	"os"

	"github.com/Azure/golua/lua/binary"
	"github.com/Azure/golua/lua/ir"
	"github.com/Azure/golua/lua/op"
)

var _ = fmt.Println
var _ = os.Exit

// v53 is the Lua 5.3 engine.
type v53 struct { *State }

// prototype pushes onto the stack a closure for the function prototype
// at index of the binary chunk
func (vm *v53) prototype(index int) *binary.Prototype {
	cls := vm.frame().closure
	return &cls.proto.Protos[index]
}

// constant pushes onto the stack the value of constant at index.
func (vm *v53) constant(index int) Value {
	cls := vm.frame().closure
	return ValueOf(cls.proto.Consts[index])
}

// fetch returns the next opcode function and instruction to execute
// incrementing the frame's instruction pointer (pc).
func (vm *v53) fetch() (cmd, ir.Instr) {
	i := vm.frame().step(1)
	return ops[i.Code()], i
}

// rk returns the value of the index that is either a register local
// value in the frame's locals stack or the n'th constant in the
// function prototype.
func (vm *v53) rk(index int) Value {
	if index > 0xFF { // Constant value
		return vm.constant(index & 0xFF)
	}
	// Register value
	return vm.frame().get(index + 1)
}

func execute(vm *v53) {
	for cmd, instr := vm.fetch(); cmd != nil; {
		vm.State.Logf("vm @ %02d : %v", vm.frame().pc, instr)
		cmd, instr = cmd(vm, instr)
	}
}

// cmd is an executor for a lua opcode.
type cmd func(*v53, ir.Instr) (cmd, ir.Instr)

// ops is a table of lua v53 opcode commands.
var ops []cmd

func init() {
	ops = []cmd{
		op.MOVE: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.move(instr)
			return vm.fetch()
		},
		op.LOADK: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.loadk(instr)
			return vm.fetch()
		},
		op.LOADKX: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.loadkx(instr)
			return vm.fetch()
		},
		op.LOADBOOL: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.loadbool(instr)
			return vm.fetch()
		},
		op.LOADNIL: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.loadnil(instr)
			return vm.fetch()
		},
		op.GETUPVAL: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.getupval(instr)
			return vm.fetch()
		},
		op.GETTABUP: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.gettabup(instr)
			return vm.fetch()
		},
		op.GETTABLE: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.gettable(instr)
			return vm.fetch()
		},
		op.SETTABUP: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.settabup(instr)
			return vm.fetch()
		},
		op.SETUPVAL: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.setupval(instr)
			return vm.fetch()
		},
		op.SETTABLE: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.settable(instr)
			return vm.fetch()
		},
		op.NEWTABLE: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.newtable(instr)
			return vm.fetch()
		},
		op.SELF: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.self(instr)
			return vm.fetch()
		},
		op.ADD: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.add(instr)
			return vm.fetch()
		},
		op.SUB: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.sub(instr)
			return vm.fetch()
		},
		op.MUL: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.mul(instr)
			return vm.fetch()
		},
		op.MOD: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.mod(instr)
			return vm.fetch()
		},
		op.POW: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.pow(instr)
			return vm.fetch()
		},
		op.DIV: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.div(instr)
			return vm.fetch()
		},
		op.IDIV: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.idiv(instr)
			return vm.fetch()
		},
		op.BAND: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.band(instr)
			return vm.fetch()
		},
		op.BOR: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.bor(instr)
			return vm.fetch()
		},
		op.BXOR: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.bxor(instr)
			return vm.fetch()
		},
		op.SHL: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.shl(instr)
			return vm.fetch()
		},
		op.SHR: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.shr(instr)
			return vm.fetch()
		},
		op.UNM: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.unm(instr)
			return vm.fetch()
		},
		op.BNOT: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.bnot(instr)
			return vm.fetch()
		},
		op.NOT: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.not(instr)
			return vm.fetch()
		},
		op.LEN: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.length(instr)
			return vm.fetch()
		},
		op.CONCAT: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.concat(instr)
			return vm.fetch()
		},
		op.JMP: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.jmp(instr)
			return vm.fetch()
		},
		op.EQ: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.eq(instr)
			return vm.fetch()
		},
		op.LT: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.lt(instr)
			return vm.fetch()
		},
		op.LE: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.le(instr)
			return vm.fetch()
		},
		op.TEST: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.test(instr)
			return vm.fetch()
		},
		op.TESTSET: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.testset(instr)
			return vm.fetch()
		},
		op.CALL: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.call(instr)
			return vm.fetch()
		},
		op.TAILCALL: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.tailcall(instr)
			return vm.fetch()
		},
		op.RETURN: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.returns(instr)
			return nil, instr
		},
		op.FORLOOP: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.forloop(instr)
			return vm.fetch()
		},
		op.FORPREP: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.forprep(instr)
			return vm.fetch()
		},
		op.TFORCALL: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.tforcall(instr)
			return vm.fetch()
		},
		op.TFORLOOP: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.tforloop(instr)
			return vm.fetch()
		},
		op.SETLIST: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.setlist(instr)
			return vm.fetch()
		},
		op.CLOSURE: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.closure(instr)
			return vm.fetch()
		},
		op.VARARG: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.vararg(instr)
			return vm.fetch()
		},
		op.EXTRAARG: func(vm *v53, instr ir.Instr) (cmd, ir.Instr) {
			vm.extraarg(instr)
			return vm.fetch()
		},
	}
}