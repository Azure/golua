package luac

import (
	"fmt"
	"os"
	"github.com/Azure/golua/lua/code"
)

var _ = fmt.Println
var _ = os.Exit

// concatJumps concatenates jump-list 'l2' into jump-list 'l1'
//
// luaK_concat
func concatJumps(fs *function, l1, l2 int) int {
	switch {
		case l2 == noJump:
			// nothing to concatenate
		case l1 == noJump:
			// no original list; 'l1' points to 'l2'
			return l2
		default:
			var (
				list = l1
				next = jumpAddr(fs, list)
			)
			for next != noJump { // find last element
				list, next = next, jumpAddr(fs, next)
			}
			fixJump(fs, list, l2) // last element links to 'l2'
	}
	return l1
}

// goIfTrue emits code to go through if 'e' is true, jump otherwise.
//
// luaK_goiftrue
func goIfTrue(fs *function, e *expr) {
	var pc int // pc of new jump
	switch fs.code.dischargeVars(fs, e); e.kind {
		case vconst, vfloat, vint, vtrue:
			pc = noJump // always true; do nothing
		case vjump: // condition?
			negateCond(fs, e) // jump when it is false
			pc = e.info // save jump position
		default:
			pc = jumpOnCond(fs, e, 0) // jump when false
	}
	e.f = concatJumps(fs, e.f, pc) // insert new jump in false list
	patch2here(fs, e.t) // true list jumps to here (to go through)
	e.t = noJump
}

// goIfFalse emits code to go through if 'e' is false, jump otherwise.
//
// luaK_goiffalse
func goIfFalse(fs *function, e *expr) {
	var pc int // pc of new jump

	switch fs.code.dischargeVars(fs, e); e.kind {
		case vnil, vfalse:
			pc = noJump // always false; do nothing
		case vjump:
			pc = e.info // already jump if true
		default:
			pc = jumpOnCond(fs, e, 1) // jump if true
	}
	e.t = concatJumps(fs, e.t, pc) // insert new jump in 't' list
	patch2here(fs, e.f) // false list jumps to here (to go through)
	e.f = noJump
}

// fixJump fixes the jump instruction at position 'pc' to jump to 'dst'.
//
// Jump addresses are relative in Lua.
//
// fixjump
func fixJump(fs *function, pc, dst int) {
	jmp := fs.instrs[pc].code
	ofs := dst - (pc + 1)
	fs.ls.assert(dst != noJump)
	if abs(ofs) > code.MaxArgSBX {
		fs.ls.syntaxErr("control structure too long")
	}
	jmp.SetSBX(ofs)
}

// jumpAddr returns the destination address of a jump instruction.
//
// Used to traverse a list of jumps.
//
// getjump
func jumpAddr(fs *function, pc int) int {
	// fmt.Printf("jumpAddr(instrs=%d, pc=%d)\n", len(fs.instrs), pc)
	ofs := fs.instrs[pc].code.SBX()
	if ofs == noJump { // point to itself represents end of list
		return noJump // end of list
	}
	return (pc + 1) + ofs // turn offset into absolute position
}

// jumpOnCond emits an instruction to jump if 'e' is 'cond' (that is, if 'cond'
// is true, code will jump if 'e' is true).
//
// Return jump position. Optimize when 'e' is 'not' something, inverting the
// condition and removing the 'not'. 
//
// jumponcond
func jumpOnCond(fs *function, e *expr, cond int) int {
	if e.kind == vreloc {
		inst := fs.instr(e).code
		if inst.Code() == code.NOT {
			fs.instrs = fs.instrs[:fs.pc-1]
			fs.pc-- // remove previous OP_NOT
			return condJump(fs, code.TEST, inst.B(), 0, not(cond))
		}
		// else go through
	}
	fs.code.discharge2any(fs, e)
	fs.code.freeexpr(fs, e)
	return condJump(fs, code.TESTSET, noReg, e.info, cond)
}

// luaK_jumpto
func jumpTo(fs *function, target int) {
	patchList(fs, fs.code.codeJump(fs), target)	
}

// negateCond negates condition 'e' (where 'e' is a comparison).
//
// negatecondition
func negateCond(fs *function, e *expr) {
	inst := jumpCtrl(fs, e.info).code
	fs.ls.assert(inst.Code().Mask().Test() && inst.Code() != code.TESTSET && inst.Code() != code.TEST)
	inst.SetA(uint(not(inst.A())))
}

// jumpCtrl returns the position of the instruction "controlling" a given jump
// (that is, its condition), or the jump itself if it is unconditional.
//
// getjumpcontrol
func jumpCtrl(fs *function, pc int) *instr {
	if pc >= 1 && fs.instrs[pc-1].code.Code().Mask().Test() {
		return fs.instrs[pc-1]
	}
	return fs.instrs[pc]
}

// condJump codes a 'condition jump', that is, a test or comparison opcode
// followed by a jump.
//
// Returns the jump position.
//
// condjump
func condJump(fs *function, op code.Opcode, a, b, c int) int {
	fs.code.codeABC(fs, op, a, b, c)
	return fs.code.codeJump(fs)
}

// patch2here adds elements in 'list' to list of pending jumps to "here" (current position).
//
// luaK_patchtohere
func patch2here(fs *function, list int) {
	fs.pclabel() // mark "here" as a jump target
	fs.jumppc = concatJumps(fs, fs.jumppc, list)
}

// patchList patches all jumps in 'list' to jump to 'target'.
//
// The assert means that we cannot fix a jump to a forward
// address because we only know addresses once code is
// generated.
//
// luaK_patchlist
func patchList(fs *function, list, target int) {
	if target == fs.pc { // 'target' is current position?
		patch2here(fs, list) // add list to pending jumps
	} else {
		fs.ls.assert(target < fs.pc)
		patchTestList(fs, list, target, noReg, target)
	}
}

// Patch destination register for a TESTSET instruction.
//
// If instruction in position 'node' is not a TESTSET, return 0 ("fails").
// Otherwise, if 'reg' is not 'NO_REG', set it as the destination register.
// Otherwise, change instruction to a simple 'TEST' (produces no register value).
//
// patchtestreg
func patchTest(fs *function, node, reg int) bool {
	inst := jumpCtrl(fs, node)
	if inst.code.Code() != code.TESTSET {
		return false // cannot patch other instructions
	}
	if reg != noReg && reg != inst.code.B() {
		inst.code.SetA(uint(reg))
	} else {
		// no register to put value or register already
		// has the value; change instruction to simple
		// test
		code := code.MakeABC(code.TEST, inst.code.B(), 0, inst.code.C())
		inst.code = &code
	}
	return true
}

// patchClose patches all jumps in 'list' to close upvalues
// up to given 'level' (the assertion checks that jumps either
// were closing nothing or were closing higher levels, from
// inner blocks).
//
// luaK_patchclose
func patchClose(fs *function, list, level int) {
	// argument +1 to reserve 0 as non-op
	for level++; list != noJump; list = jumpAddr(fs, list) {
		inst := fs.instrs[list].code
		fs.ls.assert(inst.Code() == code.JMP && inst.A() == 0 || inst.A() >= level)
		inst.SetA(uint(level))
	}
}

// Traverse a list of tests, patching their destination address and
// registers: tests producing values jump to 'vtarget' (and put their
// values in 'reg'), other tests jump to 'dtarget'.
//
// patchlistaux
func patchTestList(fs *function, list, vtarget, reg, dtarget int) {
	for list != noJump {
		// fmt.Printf("patchTestList(list=%d, vtarget=%d, reg=%d, dtarget=%d)\n", list, vtarget, reg, dtarget)
		next := jumpAddr(fs, list)
		if patchTest(fs, list, reg) {
			fixJump(fs, list, vtarget)
		} else {
			fixJump(fs, list, dtarget) // jump to default target
		}
		list = next
	}
}

// needValues checks whether list has any jumps that do not produce a value
// or produces an inverted value.
//
// need_value
func needValue(fs *function, list int) bool {
	for ; list != noJump; list = jumpAddr(fs, list) {
		if i := jumpCtrl(fs, list).code; i.Code() != code.TESTSET {
			return true
		}
	}
	return false // not found
}

// removeValues traverses a list of tests ensuring no one produces a value.
func removeValues(fs *function, list int) {
	for ; list != noJump; list = jumpAddr(fs, list) {
		patchTest(fs, list, noReg)
	}
}