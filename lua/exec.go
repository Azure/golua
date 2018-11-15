package lua

import (
	"fmt"
	"os"

	"github.com/Azure/golua/lua/binary"
	"github.com/Azure/golua/lua/vm"
)

var _ = fmt.Println
var _ = os.Exit

// v53 is the Lua 5.3 engine.
type v53 struct{ state *State }

// prototype pushes onto the stack a closure for the function prototype
// at index of the binary chunk
func (vm *v53) prototype(index int) *binary.Prototype {
	cls := vm.thread().frame().function()
	return &cls.binary.Protos[index]
}

// constant pushes onto the stack the value of constant at index.
func (vm *v53) constant(index int) Value {
	cls := vm.thread().frame().function()
	return valueOf(vm.thread(), cls.binary.Consts[index])
}

// thread returns the executing thread's state.
func (vm *v53) thread() *State { return vm.state }

// Try to convert a 'for' limit to an integer, preserving the semantics of the loop.
//
// The following explanation assumes a non-negative step; it is valid for negative
// steps mutatis mutandis.
//
// If the limit can be converted to an integer, rounding down, that is it.
//
// Otherwise, check whether the limit can be converted to a number. If the number is
// too large, it is OK to set the limit as LUA_MAXINTEGER, which means no limit.
//
// If the number is too negative, the loop should not run, because any initial integer
// value is larger than the limit. So, it sets the limit to LUA_MININTEGER.
//
// 'stopnow' corrects the extreme case when the initial value is LUA_MININTEGER, in which
// case the LUA_MININTEGER limit would still run the loop once.
// func (vm *v53) forlimit(limit, step Value) (int, bool) {
// 	return 0, false
// }

func (vm *v53) trace(instr vm.Instr) {
	if vm.thread().global.config.debug {
		fmt.Printf("vm @ ip=%02d fp=%02d: %v\n",
			vm.thread().frame().pc,
			vm.thread().frame().depth,
			instr,
		)
	}
}

// fetch returns the next opcode function and instruction to execute
// incrementing the frame's instruction pointer (pc).
func (vm *v53) fetch() (cmd, vm.Instr) {
	i := vm.thread().frame().step(1)
	return ops[i.Code()], i
}

// rk returns the value of the index that is either a register local
// value in the frame's locals stack or the n'th constant in the
// function prototype.
func (vm *v53) rk(index int) Value {
	if index > 0xFF { // Constant value
		return vm.constant(index & 0xFF)
	}
	// Registry value
	return vm.thread().frame().get(index)
}

func execute(vm *v53) {
	for cmd, instr := vm.fetch(); cmd != nil; cmd, instr = cmd(vm, instr) {
		vm.trace(instr)
	}
}

// cmd is an executor for a lua opcode.
type cmd func(*v53, vm.Instr) (cmd, vm.Instr)

// ops is a table of lua v53 opcode commands.
var ops []cmd

func init() {
	ops = []cmd{
		vm.MOVE: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.move(instr)
			return vm.fetch()
		},
		vm.LOADK: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.loadk(instr)
			return vm.fetch()
		},
		vm.LOADKX: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.loadkx(instr)
			return vm.fetch()
		},
		vm.LOADBOOL: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.loadbool(instr)
			return vm.fetch()
		},
		vm.LOADNIL: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.loadnil(instr)
			return vm.fetch()
		},
		vm.GETUPVAL: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.getupval(instr)
			return vm.fetch()
		},
		vm.GETTABUP: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.gettabup(instr)
			return vm.fetch()
		},
		vm.GETTABLE: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.gettable(instr)
			return vm.fetch()
		},
		vm.SETTABUP: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.settabup(instr)
			return vm.fetch()
		},
		vm.SETUPVAL: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.setupval(instr)
			return vm.fetch()
		},
		vm.SETTABLE: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.settable(instr)
			return vm.fetch()
		},
		vm.NEWTABLE: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.newtable(instr)
			return vm.fetch()
		},
		vm.SELF: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.self(instr)
			return vm.fetch()
		},
		vm.ADD: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.add(instr)
			return vm.fetch()
		},
		vm.SUB: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.sub(instr)
			return vm.fetch()
		},
		vm.MUL: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.mul(instr)
			return vm.fetch()
		},
		vm.MOD: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.mod(instr)
			return vm.fetch()
		},
		vm.POW: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.pow(instr)
			return vm.fetch()
		},
		vm.DIV: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.div(instr)
			return vm.fetch()
		},
		vm.IDIV: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.idiv(instr)
			return vm.fetch()
		},
		vm.BAND: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.band(instr)
			return vm.fetch()
		},
		vm.BOR: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.bor(instr)
			return vm.fetch()
		},
		vm.BXOR: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.bxor(instr)
			return vm.fetch()
		},
		vm.SHL: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.shl(instr)
			return vm.fetch()
		},
		vm.SHR: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.shr(instr)
			return vm.fetch()
		},
		vm.UNM: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.unm(instr)
			return vm.fetch()
		},
		vm.BNOT: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.bnot(instr)
			return vm.fetch()
		},
		vm.NOT: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.not(instr)
			return vm.fetch()
		},
		vm.LEN: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.length(instr)
			return vm.fetch()
		},
		vm.CONCAT: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.concat(instr)
			return vm.fetch()
		},
		vm.JMP: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.jmp(instr)
			return vm.fetch()
		},
		vm.EQ: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.eq(instr)
			return vm.fetch()
		},
		vm.LT: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.lt(instr)
			return vm.fetch()
		},
		vm.LE: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.le(instr)
			return vm.fetch()
		},
		vm.TEST: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.test(instr)
			return vm.fetch()
		},
		vm.TESTSET: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.testset(instr)
			return vm.fetch()
		},
		vm.CALL: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.call(instr)
			return vm.fetch()
		},
		vm.TAILCALL: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.tailcall(instr)
			return vm.fetch()
		},
		vm.RETURN: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.returns(instr)
			return nil, instr
		},
		vm.FORLOOP: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.forloop(instr)
			return vm.fetch()
		},
		vm.FORPREP: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.forprep(instr)
			return vm.fetch()
		},
		vm.TFORCALL: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.tforcall(instr)
			return vm.fetch()
		},
		vm.TFORLOOP: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.tforloop(instr)
			return vm.fetch()
		},
		vm.SETLIST: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.setlist(instr)
			return vm.fetch()
		},
		vm.CLOSURE: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.closure(instr)
			return vm.fetch()
		},
		vm.VARARG: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.vararg(instr)
			return vm.fetch()
		},
		vm.EXTRAARG: func(vm *v53, instr vm.Instr) (cmd, vm.Instr) {
			vm.extraarg(instr)
			return vm.fetch()
		},
	}
}
