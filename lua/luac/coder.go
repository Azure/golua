package luac

import (
    "fmt"
	"github.com/fibonacci1729/golua/lua/code"
)

var _ = fmt.Println

// instr groups an instruction with its associated line.
type instr struct {
	code *code.Instr
	line int32
}

// Expression and variable descriptor.
//
// Code generation for variables and expressions can be delayed to allow
// optimizations; An 'expression' structure describes a potentially-delayed
// variable/expression. It has a description of its "main" value plus a
// list of conditional jumps that can also produce its value (generated
// by short-circuit operators 'and'/'or').
type expr struct {
    kind exprKind
    info int // generic use
    t, f int // patch list of 'exit when true/false'
    // constant values
    ival int64   // vint
    nval float64 // vfloat
    // indexed variable (vindexed)
    index struct{
        kind exprKind // vlocal or vupval
        t, k int      // table (register or upvalue) / key index (R/K)
    }
}

// kinds of variables/expressions
type exprKind int

const (
    vvoid exprKind = iota // no expression/empty list; when expression describes the last expression in a list
    vnil                  // constant nil
    vtrue                 // constant true
    vfalse                // constant false
    vconst                // constant in 'k'; info = index of constant in 'k'
    vfloat                // floating constant; nval = numerical float value
    vint                  // integer constant; ival = numerical integer value
    vnonreloc             // expression has its value in a fixed register; in
    vlocal                // local variable; info = local register
    vupval                // upvalue variable; info = index of upvalue in 'upvalues'
    vindexed              // indexed variable; kind (vlocal or vupval), t = register or upvalue, k = key's R/K index
    vjump                 // expression is a test/comparison; info = pc of corresponding jump instruction
    vreloc                // expression can put result in any register; info = instruction pc
    vcall                 // expression is a function call; info = instruction pc
    vvararg               // vararg expression; info = instruction pc
)

var exprKinds = [...]string{
    vvoid:     "void",
    vnil:      "nil",
    vtrue:     "true",
    vfalse:    "false",
    vconst:    "constant",
    vfloat:    "float",
    vint:      "int",
    vnonreloc: "nonrelocable",
    vlocal:    "local",
    vupval:    "upval",
    vindexed:  "indexed", 
    vjump:     "jump",
    vreloc:    "relocable",
    vcall:     "call",
    vvararg:   "vararg",
}

func (kind exprKind) String() string { return exprKinds[kind] }

func (e *expr) numeral() bool { return !e.jumps() && (e.kind == vint || e.kind == vfloat) }
func (e *expr) fixed() bool { return e.kind == vnonreloc || e.kind == vlocal }
func (e *expr) isVar() bool { return vlocal <= e.kind && e.kind <= vindexed }
func (e *expr) retsX() bool { return e.kind == vcall || e.kind == vvararg }
func (e *expr) jumps() bool { return e.t != e.f }

// TODO
func (e *expr) String() string { return e.kind.String() }

func (e *expr) init(kind exprKind, info int) *expr {
	e.t, e.f = noJump, noJump
	e.kind   = kind
	e.info   = info
    return e
}

// code generator
type coder struct {}

// code_loadbool
func (c *coder) codeLoadBool(fs *function, a, b, jump int) int {
	fs.pclabel() // those instructions may be jump targets
	return c.codeABC(fs, code.LOADBOOL, a, b, jump)
}

// codeClosure codes instruction to create new closure in parent function.
// The CLOSURE instruction must use the last available register, so that,
// if it invokes the GC, the GC knows which registers are in use at that time.
func (c *coder) codeClosure(fs *function, e *expr) {
	parent := fs.parent
	e.init(vreloc, c.codeABx(parent, code.CLOSURE, 0, len(parent.fn.Protos)-1))
	c.expr2next(parent, e) // fix it at the last register
}

// codeSetList emits a 'SETLIST' instruction.
//
// 'base' is the register that keeps the table; 'elemN' is #table plus those to be stored now;
// 'pending' is the number of values (in registers 'base+1', ...) to add to table (or LUA_MULTRET
// to add up to stack top).
func (c *coder) codeSetList(fs *function, base, elemN, pending int) {
	var (
		n = (elemN - 1)/fieldsPerFlush + 1
		b = pending
	)
	if pending == multRet {
		b = 0
	}
	fs.ls.assert(pending != 0 && pending <= fieldsPerFlush)
	switch {
		case n <= code.MaxArgC:
			c.codeABC(fs, code.SETLIST, base, b, n)
		case n <= code.MaxArgAX:
			c.codeABC(fs, code.SETLIST, base, b, 0)
			c.codeExtra(fs, n)
		default:
			fs.ls.syntaxErr("constructor too long")
	}
	fs.free = base + 1 // free registers with list values
}

// codeBinary emits code for binary expressions that "produce values" (everything but
// logical operators 'and' / 'or' and comparison operators).
//
// Expression to produce final result will be encoded in 'e1'. Because 'expr2rk' can
// free registers, its calls must be in "stack order" (that is, first on 'e2', which
// may have more recent registers to be released).
func (c *coder) codeBinary(fs *function, op code.Opcode, e1, e2 *expr, line int) {
	rk2 := c.expr2rk(fs, e2) // __ both operands are "RK"
	rk1 := c.expr2rk(fs, e1) //
	c.freeexpr2(fs, e1, e2)
	e1.info = c.codeABC(fs, op, 0, rk1, rk2) // generate opcode
	e1.kind = vreloc // result is relocatable
	fs.fixLine(line)
}

// codeUnary emits code for unary expressions that "produce values" (everything but 'not').
//
// Expression to produce the final result will be encoded in 'e'.
func (c *coder) codeUnary(fs *function, op code.Opcode, e *expr, line int) {
	r := c.expr2any(fs, e) // opcodes operate only on registers
	c.freeexpr(fs, e)
	e.info = c.codeABC(fs, op, 0, r, 0) // generate opcode
	e.kind = vreloc // result is relocatable
	fs.fixLine(line)
}

// codeIndex creates the expression 't[k]'.
//
// 't' must have its final result already in a register or upvalue.
func (c *coder) codeIndex(fs *function, t, k *expr) {
	fs.ls.assert(!t.jumps() && (t.fixed() || t.kind == vupval))
	t.index.t = t.info // register or upvalue index
	t.index.k = c.expr2rk(fs, k) // R/K index for key
	if t.kind == vupval {
		t.index.kind = vupval
	} else {
		t.index.kind = vlocal
	}
	t.kind = vindexed
}

func (c *coder) codeReturn(fs *function, first, retN int) {
	c.codeABC(fs, code.RETURN, first, retN+1, 0)
}

// codeSelf emits a SELF instruction (convert expression 'e' into 'e:key(e,').
func (c *coder) codeSelf(fs *function, e, k *expr) {
	c.expr2any(fs, e)
	re := e.info // register where 'e' was placed
	c.freeexpr(fs, e)
	e.info = fs.free // base register for SELF
	e.kind = vnonreloc // self expression has a fixed register
	fs.reserve(2) // function and 'self' produced by SELF
	c.codeABC(fs, code.SELF, e.info, re, c.expr2rk(fs, k))
	c.freeexpr(fs, k)
}

// codeJump emits a jump instruction and returns it position, so its destination can be fixed
// later (with 'fixJump'). If there are jumps to this position (kept in 'jumppc'), link them
// all together so that 'patchlistaux' will fix all them directly to the final destination.
func (c *coder) codeJump(fs *function) int {
	jpc := fs.jumppc // save list of jumps to here
	fs.jumppc = noJump // no more jumps to here
	pc := c.codeAsBx(fs, code.JMP, 0, noJump)
	return concatJumps(fs, pc, jpc)
}

// codeNils creates the LOADNIL instruction, but tries to optimize:
// if the previous instruction is also LOADNIL and ranges are compativle,
// adjust range of previous instruction instead of emitting a new one.
// For instance, 'local a; local b' will generate a single opcode.
func (c *coder) codeNils(fs *function, from, n int) {
	if last := from + n - 1; fs.pc > fs.target { // no jumps to current position?
		previous := fs.instrs[fs.pc-1].code
		if previous.Code() == code.LOADNIL { // previous is LOADNIL?
			// get previous range
			r0 := previous.A()
			r1 := r0 + previous.B()
			if (r0 <= from && from <= r1 + 1) || (from <= r0 && r0 <= last + 1) { // can connect both?
				from = min(from, r0)
				last = max(last, r1)
				previous.SetA(uint(from))
				previous.SetB(uint(last-from))
				return
			}
		}
	}
 	// no optimization
	c.codeABC(fs, code.LOADNIL, from, n - 1, 0)
}

func (c *coder) codeExtra(fs *function, a int) int {
	fs.ls.assert(a <= code.MaxArgAX)
	return c.code(fs, code.MakeAx(code.EXTRAARG, a))
}

func (c *coder) codeKst(fs *function, reg, kst int) int {
	if kst <= code.MaxArgBX {
		return c.codeABx(fs, code.LOADK, reg, kst)
	}
	pc := c.codeABx(fs, code.LOADKX, reg, 0)
	c.codeExtra(fs, kst)
	return pc
}

// codeNot codes a 'not e' instruction, doing constant folding.
func (c *coder) codeNot(fs *function, e *expr) {
	switch c.dischargeVars(fs, e); e.kind {
		case vconst, vfloat, vint, vtrue:
			// false == not "x" == not 0.5 == not 1 == not true 
			e.kind = vfalse
		case vreloc, vnonreloc:
			c.discharge2any(fs, e)
			c.freeexpr(fs, e)
			e.info = c.codeABC(fs, code.NOT, 0, e.info, 0)
			e.kind = vreloc
		case vnil, vfalse:
			// true == not nil == not false
			e.kind = vtrue
		case vjump:
			negateCond(fs, e)
		default:
			panic("unreachable") // cannot happen
	}
	e.t, e.f = e.f, e.t // interchange true and false lists
	removeValues(fs, e.f) // values are useless when negated
	removeValues(fs, e.t)
}

func (c *coder) codeAsBx(fs *function, op code.Opcode, a, sbx int) int {
	return c.codeABx(fs, op, a, sbx + code.MaxArgSBX)
}

func (x *coder) codeABC(fs *function, op code.Opcode, a, b, c int) int {
	fs.ls.assert(op.Mode() == code.ModeABC)
	fs.ls.assert(!op.Mask().B(code.ArgN) || b == 0)
	fs.ls.assert(!op.Mask().C(code.ArgN) || c == 0)
	fs.ls.assert(a <= code.MaxArgA && b <= code.MaxArgB && c <= code.MaxArgC)
	return x.code(fs, code.MakeABC(op, a, b, c))
}

func (c *coder) codeABx(fs *function, op code.Opcode, a, bx int) int {
	fs.ls.assert(op.Mode() == code.ModeABx || op.Mode() == code.ModeAsBx)
	fs.ls.assert(op.Mask().C(code.ArgN))
	fs.ls.assert(a <= code.MaxArgA && bx <= code.MaxArgBX)
	return c.code(fs, code.MakeABx(op, a, bx))
}

func (c *coder) code(fs *function, i code.Instr) int {
	c.dischargeJumps(fs)
	fs.instrs = append(fs.instrs, &instr{&i, int32(fs.ls.last)})
	fs.pc++
	return fs.pc - 1
}

// return1 fixes an expression to return one result.
//
// If expression is not a multi-ret expression (function call or vararg), it already
// returns one result, so nothing needs to be done.
//
// Function calls become VNONRELOC expressions (as its result comes fixed in the base
// register of the call), while vararg expressions become VRELOCABLE (as OP_VARARG
// puts its results where it wants). (Calls are created returning one result, so that
// does not need to be fixed.)
//
// luaK_setoneret
func (c *coder) return1(fs *function, e *expr) {
	switch e.kind {
		case vvararg: // expression is vararg?
			fs.instr(e).code.SetB(2)
			e.kind = vreloc // can relocate its simple result
		case vcall: // expression is an open function call?
			// already returns 1 value
			fs.ls.assert(fs.instr(e).code.C() == 2)
			e.kind = vnonreloc // result has fixed position
			e.info = fs.instr(e).code.A()
	}
}

// Fix an expression to return the number of results 'nresults'. Either 'e' is a multi-ret
// expression (function call or vararg) or 'nresults' is LUA_MULTRET (as any expression can
// satisfy that).
//
// luaK_setreturns
func (c *coder) returnN(fs *function, e *expr, retN int) {
	switch e.kind {
		case vvararg:
			fs.instr(e).code.SetB(uint(retN + 1))
			fs.instr(e).code.SetA(uint(fs.free))
			fs.reserve(1)
		case vcall:
			fs.instr(e).code.SetC(uint(retN+1))
		default:
			fs.ls.assert(retN == multRet)
	}
}

// luaK_setmultret
func (c *coder) returnX(fs *function, e *expr) {
	c.returnN(fs, e, multRet)
}

// storeVar generates code to store the result of expression 'e'
// into variable 'v'.
//
// luaK_storevar
func (c *coder) storeVar(fs *function, v, e *expr) {
	switch v.kind {
		case vindexed:
			var op code.Opcode
			if v.index.kind == vlocal {
				op = code.SETTABLE
			} else {
				op = code.SETTABUP
			}
			c.codeABC(fs, op, v.index.t, v.index.k, c.expr2rk(fs, e))
			c.freeexpr(fs, e)
			return
		case vupval:
			c.codeABC(fs, code.SETUPVAL, c.expr2any(fs, e), v.info, 0)
			c.freeexpr(fs, e)
			return
		case vlocal:
			c.freeexpr(fs, e)
			c.expr2reg(fs, e, v.info) // compute 'e' into proper place
			return
	}
	panic("unreachable")
}

// dischargeVars ensures that expression 'e' is not a variable.
func (c *coder) dischargeVars(fs *function, e *expr) {
	switch e.kind {
		case vvararg, vcall:
			c.return1(fs, e)
		case vindexed:
			var op code.Opcode
			c.freereg(fs, e.index.k)
			if e.index.kind == vlocal {
				// 't' is in a register
				c.freereg(fs, e.index.t)
				op = code.GETTABLE
			} else {
				// 't' is in an upvalue
				fs.ls.assert(e.index.kind == vupval)
				op = code.GETTABUP
			}
			e.info = c.codeABC(fs, op, 0, e.index.t, e.index.k)
			e.kind = vreloc
		case vlocal:
			// already in a register; becomes a non-relocatable value
			e.kind = vnonreloc
		case vupval:
			// move value to some (pending) register
			e.info = c.codeABC(fs, code.GETUPVAL, 0, e.info, 0)
			e.kind = vreloc
	}
}

// dischargeToAny ensures the expression 'e' is in any register.
//
// discharge2anyreg
func (c *coder) discharge2any(fs *function, e *expr) {
	if e.kind != vnonreloc { // no fixed register yet?
		fs.reserve(1) // get a register
		c.discharge2reg(fs, e, fs.free-1) // put 'e' there
	}
}

// dischargeTo ensures the expression 'e' value is in register 'reg' (and therefore
// 'e' will become a non-relocatable expression).
//
// discharge2reg
func (c *coder) discharge2reg(fs *function, e *expr, reg int) {
	switch c.dischargeVars(fs, e); e.kind {
		case vtrue, vfalse:
			fs.code.codeABC(fs, code.LOADBOOL, reg, b2i(e.kind == vtrue), 0)
		case vnonreloc:
			if reg != e.info {
				fs.code.codeABC(fs, code.MOVE, reg, e.info, 0)
			}
		case vreloc:
			// instruction will put result in 'reg'
			fs.instr(e).code.SetA(uint(reg))
		case vconst:
			c.codeKst(fs, reg, e.info)
		case vfloat:
			c.codeKst(fs, reg, fs.numberk(e.nval))
		case vint:
			c.codeKst(fs, reg, fs.constant(e.ival))
		case vnil:
			c.codeNils(fs, reg, 1)
		default:
			fs.ls.assert(e.kind == vjump)
			return
	}
	e.kind = vnonreloc
	e.info = reg
}

// dischargeJumps ensures all pending jumps to current position are fixed
// (jumping to current position with no values) and reset list of pending
// jumps.
func (c *coder) dischargeJumps(fs *function) {
	patchTestList(fs, fs.jumppc, fs.pc, noReg, fs.pc)
	fs.jumppc = noJump
}

// freeExprs frees registers used by expressions 'e1' and 'e2' (if any)
// in proper order.
func (c *coder) freeexpr2(fs *function, e1, e2 *expr) {
	var ( r1, r2 = -1, -1 )

	if e1.kind == vnonreloc {
		r1 = e1.info
	}
	if e2.kind == vnonreloc {
		r2 = e2.info
	}
	if r1 > r2 {
		c.freereg(fs, r1)
		c.freereg(fs, r2)
	} else {
		c.freereg(fs, r2)
		c.freereg(fs, r1)
	}
}

// freeExpr frees the register used by expression 'e' (if any)
func (c *coder) freeexpr(fs *function, e *expr) {
	if e.kind == vnonreloc {
		c.freereg(fs, e.info)
	}
}

// freereg frees register 'reg', if it is neither a constant index nor
// a local variable.
func (c *coder) freereg(fs *function, reg int) {
	if !code.IsKst(reg) && reg >= fs.active {
		fs.free--
		fs.ls.assert(reg == fs.free)
	}
}

func (c *coder) expr2next(fs *function, e *expr) {
	c.dischargeVars(fs, e)
	c.freeexpr(fs, e)
	fs.reserve(1)
	c.expr2reg(fs, e, fs.free-1)
}

// expr2reg ensures final expression result (including results from its jump lists)
// is in register 'reg'. If expression has jumps, need to patch these jumps either
// to its final position or to "load" instructions (for those tests that do not
// produce values).
//
// exp2reg
func (c *coder) expr2reg(fs *function, e *expr, reg int) {
	if c.discharge2reg(fs, e, reg); e.kind == vjump { // expression itself is a test?
		e.t = concatJumps(fs, e.t, e.info) // put this jump in 't' list
	}
	if e.jumps() {
		var (
			pt, pf = noJump, noJump // position of eventual LOAD true/false
			final int   			// position after whole expression
		)
		if needValue(fs, e.t) || needValue(fs, e.f) {
			var jp int
			if jp = noJump; e.kind != vjump {
				jp = fs.code.codeJump(fs)
			} 
			pf = fs.code.codeLoadBool(fs, reg, 0, 1)
			pt = fs.code.codeLoadBool(fs, reg, 1, 0)
			patch2here(fs, jp)
		}
		final = fs.pclabel()
		patchTestList(fs, e.f, final, reg, pf)
		patchTestList(fs, e.t, final, reg, pt)
	}
	e.t, e.f = noJump, noJump
	e.kind = vnonreloc
	e.info = reg
}

// expr2any ensures the final expression result (including results from
// its jump lists) is in some (any) register and return that register.
func (c *coder) expr2any(fs *function, e *expr) int {
	if c.dischargeVars(fs, e); e.kind == vnonreloc {
		if !e.jumps() {
			return e.info
		}
		if e.info >= fs.active {
			c.expr2reg(fs, e, e.info)
			return e.info
		}
	}
	c.expr2next(fs, e)
	return e.info
}

// expr2val ensures the final expression result is either in a register
// or it is a constant.
func (c *coder) expr2val(fs *function, e *expr) {
	if e.jumps() {
		c.expr2any(fs, e)
	} else {
		c.dischargeVars(fs, e)
	}
}

// expr2ru ensures the final expression result is either in a register or
// in an upvalue.
func (c *coder) expr2ru(fs *function, e *expr) {
	if e.kind != vupval || e.jumps() {
		c.expr2any(fs, e)
	}
}

// expr2rk ensures the final expression result is in a valid R/K index;
// that is, it is either in a register or in 'k' with an index in the
// range of R/K indices.
//
// Returns the R/K index.
func (c *coder) expr2rk(fs *function, e *expr) int {
	switch c.expr2val(fs, e); e.kind {
		case vfloat:
			e.info = fs.numberk(e.nval)
			e.kind = vconst
		case vfalse:
			e.info = fs.constant(false)
			e.kind = vconst
		case vtrue:
			e.info = fs.constant(true)
			e.kind = vconst
		case vnil:
			e.info = fs.constant(nil)
			e.kind = vconst
		case vint:
			e.info = fs.constant(e.ival)
			e.kind = vconst
	}
	if e.kind == vconst {
		if e.info <= code.MaxIndexRK { // constant fits in 'argC'?
			return code.RK(e.info)
		}
	}
	return c.expr2any(fs, e)
}

func expr2str(e *expr) string {
	if e.kind == vindexed {
		return fmt.Sprintf("indexed.%s %d:%d", e.index.kind, e.index.t, e.index.k)
	}
	return e.kind.String()
}