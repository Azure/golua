package luac

import (
	"io/ioutil"
	"bytes"
	"math"
	"fmt"
	"os"
	"io"

	"github.com/fibonacci1729/golua/lua/code"
)

var _ = fmt.Println
var _ = os.Exit

func readSource(file string, source interface{}) (buf *bytes.Buffer, err error) {
	var b []byte
	switch src := source.(type) {
	case io.Reader:
		b, err = ioutil.ReadAll(src)
	case string:
		b = []byte(src)
	case []byte:
		b = src
	case nil:
		if file != "" && len(file) > 1 {
			if file[0] == '@' || file[0] == '=' {
				file = file[1:]
			}
		}
		b, err = ioutil.ReadFile(file)
	default:
		return nil, fmt.Errorf("invalid source: %T", src)
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %s", file, err)
	}
	return bytes.NewBuffer(b), nil
}

type parser struct {
	level int
}

func (p *parser) enterLevel(ls *lexical) { p.level++; checkLimit(ls.fs, p.level, maxCalls, "Go levels") }
func (p *parser) leaveLevel(ls *lexical) { p.level-- } // TODO

// mainfunc
func (p *parser) mainFunc(ls *lexical) *code.Proto {
	ls.fs = new(function).open(ls)
	ls.fs.fn.Vararg = true
	e := new(expr).init(vlocal, 0)
	ls.fs.makeup("_ENV", e)
	ls.next()
	// dumptoks(ls)
	p.stmts(ls)
	ls.expect(tEOS)
	fs := ls.fs
	fn := fs.fn
	ls.fs = fs.close()
	ls.assert(fs.parent == nil && fs.nups == 1 && ls.fs == nil)
	// all scopes should be correctly finished
	ls.assert(len(ls.active) == 0 && len(ls.gotos) == 0 && len(ls.labels) == 0)
	return fn
}

//
// Statements
//

// check whether current token is in the follow set of a block.
// 'until' closes syntactical blocks, but do not close scope,
// so it is handled in separate.
//
// block_follow
func (p *parser) follows(ls *lexical, until bool) bool {
	switch ls.token.char {
		case tElse, tElseIf, tEnd, tEOS:
			return true
		case tUntil:
			return until
	}
	return false
}

// statlist
func (p *parser) stmts(ls *lexical) {
	defer un(trace(ls, "Stmts"))
	for !p.follows(ls, true) {
		if ls.token.char == tReturn {
			p.stmt(ls)
			return // 'return' must be last statement
		}
		p.stmt(ls)
	}
}

// statement
func (p *parser) stmt(ls *lexical) {
	defer un(trace(ls, "Stmt"))
	p.enterLevel(ls)
	line := ls.line

	switch ls.token.char {
		case tBreak, tGoto:
			p.gotostmt(ls, ls.fs.code.codeJump(ls.fs))
		case tFunction:
			p.funcstmt(ls, line)
		case tReturn:
			ls.next()
			p.retstmt(ls)
		case tColon2:
			ls.next()
			p.labelstmt(ls, p.ident(ls), line)
		case tRepeat:
			p.repeat(ls, line)
		case tWhile:
			p.while(ls, line)
		case tLocal:
			if ls.next(); ls.test(tFunction) {
				p.localfunc(ls)
			} else {
				p.localstmt(ls)
			}
		case tFor:
			p.forloop(ls, line)
		case tDo:
			ls.next()
			p.block(ls)
			ls.match(tEnd, tDo, line)
		case tIf:
			p.ifstmt(ls, line)
		case ';':
			ls.next()
		default:
			p.exprstmt(ls)
	}
	ls.assert(ls.fs.fn.StackN >= ls.fs.free && ls.fs.free >= ls.fs.active)
	ls.fs.free = ls.fs.active // free registers
	p.leaveLevel(ls)
}

// stat -> func | assignment
func (p *parser) exprstmt(ls *lexical) {
	defer un(trace(ls, "ExprStmt"))
	var lhs assignment
	p.suffixed(ls, &lhs.expr)
	if ls.token.char == '=' || ls.token.char == ',' {
		// stat -> assignment
		lhs.prev = nil
		p.assignment(ls, &lhs, 1)
	} else {
		// stat -> func
		if lhs.expr.kind != vcall {
			ls.syntaxErr("syntax error")
		}
		ls.fs.instr(&lhs.expr).code.SetC(1) // call statement uses no results
	}
}

// funcstat -> FUNCTION funcname body
func (p *parser) funcstmt(ls *lexical, line int) {
	var ( v, b expr )
	ls.next() // skip FUNCTION
	p.funcbody(ls, &b, p.funcname(ls, &v), line)
	ls.fs.code.storeVar(ls.fs, &v, &b)
	ls.fs.fixLine(line) // definition "happens" in the first line
}

// funcname -> NAME {fieldsel} [':' NAME]
func (p *parser) funcname(ls *lexical, v *expr) (method bool) {
	for p.variable(ls, v); ls.token.char == '.'; {
		p.selector(ls, v)
	}
	if ls.token.char == ':' {
		method = true
		p.selector(ls, v)
	}
	return method
}

func (p *parser) assignment(ls *lexical, lhss *assignment, varN int) {
	defer un(trace(ls, "Assignment"))
	if !lhss.expr.isVar() {
		ls.syntaxErr("syntax error")
	}
	var rhs expr
	if ls.test(',') {
		// assignment -> ',' suffixedexp assignment
		lhs := &assignment{prev: lhss}
		p.suffixed(ls, &lhs.expr)
		if lhs.expr.kind != vindexed {
			// checkConflict(ls, lhss, &lhs.expr)
		}
		checkLimit(ls.fs, varN + p.level, maxCalls, "Go levels")
		p.assignment(ls, lhs, varN+1)
	} else {
		// assignment -> '=' explist
		ls.expect('=')
		ls.next()
		if n := p.exprs(ls, &rhs); n != varN {
			adjustAssign(ls.fs, varN, n, &rhs)
		} else {
			ls.fs.code.return1(ls.fs, &rhs) // close last expression
			ls.fs.code.storeVar(ls.fs, &lhss.expr, &rhs)
			return
		}
	}
	rhs.init(vnonreloc, ls.fs.free-1)
	ls.fs.code.storeVar(ls.fs, &lhss.expr, &rhs)
}

func (p *parser) gotostmt(ls *lexical, pc int) {
	var (
		line = ls.line
		name string
	)
	if ls.test(tGoto) {
		name = p.ident(ls)
	} else {
		ls.next() // skip 'break'
		name = "break"
	}
	goto_ := ls.fs.label(&ls.gotos, name, line, pc)
	findlabel(ls, goto_) // close it if label already defined
}

// label -> '::' NAME '::'
func (p *parser) labelstmt(ls *lexical, name string, line int) {
	defer un(trace(ls, "LabelStmt"))

	// check for repeated labels in the same block.	
	for i := ls.fs.block.label0; i < len(ls.labels); i++ {
		if lbl := ls.labels[i]; lbl.label == name {
			msg := fmt.Sprintf("label '%s' already defined on line %d", lbl.line)
			ls.semanticErr(msg)
		}
	}
	ls.expect(tColon2)
	ls.next()
	// create new entry for this label
	lbl := ls.fs.label(&ls.labels, name, line, ls.fs.pclabel())
	p.skipnoop(ls) // skip other no-op statements
	if p.follows(ls, false) { // label is last no-op statement in the block?
		// assume that locals are already out of scope
		ls.labels[lbl].level = ls.fs.block.active
	}
	findgotos(ls, ls.labels[lbl])
}

// stat -> LOCAL NAME {',' NAME} ['=' explist]
func (p *parser) localstmt(ls *lexical) {
	defer un(trace(ls, "LocalStmt"))
	var (
		varsN int
		exprN int
		e expr
	)
	for {
		ls.fs.declare(p.ident(ls))
		varsN++

		if !ls.test(',') {
			break
		}
	}
	if ls.test('=') {
		exprN = p.exprs(ls, &e)
	} else {
		e.kind = vvoid
		exprN = 0
	}
	adjustAssign(ls.fs, varsN, exprN, &e)
	adjustLocals(ls.fs, varsN)
}

func (p *parser) localfunc(ls *lexical) {
	defer un(trace(ls, "LocalFunc"))

	ls.fs.declare(p.ident(ls)) // new local variable
	adjustLocals(ls.fs, 1) // enter its scope
	var e expr
	p.funcbody(ls, &e, false, ls.line) // function created in next register
	// debug information will only see the variable after this point!
	ls.fs.local(e.info).Live = int32(ls.fs.pc)
}

// body -> '(' parlist ')' block END
func (p *parser) funcbody(ls *lexical, e *expr, method bool, line int) {
	fs := new(function).open(ls)
	fs.fn.SrcPos = line
	// add new prototype into the current function's list of prototypes.
	ls.fs.fn.Protos = append(ls.fs.fn.Protos, fs.fn)
	ls.fs = fs
	ls.expect('(')
	ls.next()
	if method {
		ls.fs.declare("self") // create 'self' parameter
		adjustLocals(ls.fs, 1)
	}
	p.parameters(ls)
	ls.expect(')')
	ls.next()
	p.stmts(ls)
	fs.fn.EndPos = ls.line
	ls.match(tEnd, tFunction, line)
	fs.code.codeClosure(fs, e)
	ls.fs = fs.close()
}

// parlist -> [ param { ',' param } ]
func (p *parser) parameters(ls *lexical) {
	var params int
	if ls.token.char != ')' { // is parameter list not empty?
		for !ls.fs.fn.Vararg {
			switch ls.token.char {
				case tName: // parameter
					ls.fs.declare(p.ident(ls))
					params++
				case tDots: // vararg
					ls.next()
					ls.fs.fn.Vararg = true
				default:
					ls.syntaxErr("<name> or '...' expected")
			}
			if !ls.test(',') {
				break
			}
		}
	}
	adjustLocals(ls.fs, params)
	ls.fs.fn.ParamN = ls.fs.active
	ls.fs.reserve(ls.fs.active) // reserve registers for parameters
}

// stat -> RETURN [explist] [';']
func (p *parser) retstmt(ls *lexical) {
	defer un(trace(ls, "ReturnStmt"))
	var ( first, retsN int ) // register with returned values
	if p.follows(ls, true) || ls.token.char == ';' {
		first, retsN = 0, 0 // return no values
	} else {
		var e expr
		retsN = p.exprs(ls, &e) // optional return values
		if e.retsX() {
			ls.fs.code.returnX(ls.fs, &e)
			if e.kind == vcall && retsN == 1 { // tail call?
				ls.fs.instr(&e).code.SetOp(code.TAILCALL)
				ls.assert(ls.fs.instr(&e).code.A() == ls.fs.active)
			}
			first = ls.fs.active
			retsN = multRet // return all values
		} else {
			if retsN == 1 { // only one single value?
				first = ls.fs.code.expr2any(ls.fs, &e)
			} else {
				ls.fs.code.expr2next(ls.fs, &e) // values must go to the stack
				first = ls.fs.active // return all active values
				ls.assert(retsN == ls.fs.free - first)
			}
		}
	}
	ls.fs.code.codeReturn(ls.fs, first, retsN)
	ls.test(';') // skip optional semicolon
}


// ifstat -> IF cond THEN block {ELSEIF cond THEN block} [ELSE block] END
func (p *parser) ifstmt(ls *lexical, line int) {
	defer un(trace(ls, "IfStmt"))
	// exit list for finished parts
	escs := p.ifthen(ls, noJump)
	for ls.token.char == tElseIf { // IF cond THEN block
		escs = p.ifthen(ls, escs) // ELSEIF cond THEN block
	}
	if ls.test(tElse) {
		p.block(ls) // 'else' path
	}
	ls.match(tEnd, tIf, line)
	patch2here(ls.fs, escs) // patch escape list to 'if' end
}

// test_then_block -> [IF | ELSEIF] cond THEN block
func (p *parser) ifthen(ls *lexical, escs int) int {
	defer un(trace(ls, "IfThen"))
	var (
		expr expr
		jump int // instruction to skip 'then' code (if condition is false)
	)
	ls.next() // skip IF or ELSEIF
	p.expr(ls, &expr) // read condition
	ls.expect(tThen)
	ls.next()
	if ls.token.char == tGoto || ls.token.char == tBreak {
		goIfFalse(ls.fs, &expr) // will jump to label if condition is true
		ls.fs.enter(false) // must enter block before 'goto'
		p.gotostmt(ls, expr.t)  // handle goto/break
		p.skipnoop(ls) // skip other no-op statements
		if p.follows(ls, false) { // 'goto' is the entire block?
			ls.fs.leave()
			return escs // and that is it
		} else { // must skip over 'then' part if condition is false
			jump = ls.fs.code.codeJump(ls.fs)
		}
	} else { // regular case (not goto/break)
		goIfTrue(ls.fs, &expr)
		ls.fs.enter(false)
		jump = expr.f
	}
	p.stmts(ls) // 'then' part
	ls.fs.leave()
	if ls.token.char == tElse || ls.token.char == tElseIf {
		// 'else'/'elseif' follows 'then'; jump over it
		escs = concatJumps(ls.fs, escs, ls.fs.code.codeJump(ls.fs))
	}
	patch2here(ls.fs, jump)
	return escs
}

// block -> statlist
func (p *parser) block(ls *lexical) {
	ls.fs.enter(false)
	p.stmts(ls)
	ls.fs.leave()
}

func (p *parser) skipnoop(ls *lexical) {
	for ls.token.char == ';' || ls.token.char == tColon2 {
		p.stmt(ls)
	}
}

//
// Expressions
//

// explist -> expr { ',' expr }
func (p *parser) exprs(ls *lexical, e *expr) int {
	defer un(trace(ls, "Exprs"))
	n := 1 // at least one expression
	p.expr(ls, e)
	for ls.test(',') {
		ls.fs.code.expr2next(ls.fs, e)
		p.expr(ls, e)
		n++
	}
	return n
}

func (p *parser) expr(ls *lexical, e *expr) {
	defer un(trace(ls, "Expr"))
	p.subexpr(ls, e, 0)
}

func (p *parser) expr1(ls *lexical) int {
	defer un(trace(ls, "Expr1"))
	var e expr
	p.expr(ls, &e)
	ls.fs.code.expr2next(ls.fs, &e)
	ls.assert(e.kind == vnonreloc)
	return e.info
}

// subexpr -> (simpleexp | unop subexpr) { binop subexpr }
//
// where 'binop' is any binary operator with a priority higher than 'limit'
func (p *parser) subexpr(ls *lexical, e *expr, limit int) (op code.Op) {
	defer un(trace(ls, "SubExpr"))
	p.enterLevel(ls)
	if uop := unaryOp(ls.token.char); uop != code.OpNone {
		line := ls.line
		ls.next()
		p.subexpr(ls, e, unaryPriority)
		p.prefix(ls, uop, e, line)
	} else {
		p.simple(ls, e)
	}
	// expand while operators have prioriteis higher than 'limit'
	op = binaryOp(ls.token.char)
	for op != code.OpNone && priority[op-1].lhs > limit {
		line := ls.line
		var e2 expr
		ls.next()
		p.infix(ls, op, e)
		// read sub-expression with higher priority
		nextop := p.subexpr(ls, &e2, priority[op-1].rhs)
		p.postfix(ls, op, e, &e2, line)
		op = nextop
	}
	p.leaveLevel(ls)
	return op // return first untreated operator
}

// prefix applies the prefix operation 'op' to expression 'e'.
func (p *parser) prefix(ls *lexical, op code.Op, e *expr, line int) {
	y := new(expr).init(vint, 0)
	switch fs := ls.fs; op {
		case code.OpMinus, code.OpBnot: // use 'y' as fake 2nd operand
			if constfold(fs, op, e, y) {
				return
			}
			fallthrough
		case code.OpLen:
			fs.code.codeUnary(fs, code.Opcode(op-code.OpMinus)+code.UNM, e, line)
			return
		case code.OpNot:
			fs.code.codeNot(fs, e)
			return
	}
	panic("unreachable")
}

// infix processes the 1st operand 'v' of a binary operation 'op' before
// reading the 2nd operand.
func (p *parser) infix(ls *lexical, op code.Op, e *expr) {
	defer un(trace(ls, "Infix"))
	switch op {
		case code.OpConcat:
			// operand must be on the 'stack'
			ls.fs.code.expr2next(ls.fs, e)
		case code.OpAnd:
			// go ahead only if 'e' is true
			goIfTrue(ls.fs, e)
		case code.OpOr:
			// go ahead only if 'e' is false
			goIfFalse(ls.fs, e)
		case code.OpAdd,
			code.OpSub,
			code.OpMul,
			code.OpDivF,
			code.OpDivI,
			code.OpMod,
			code.OpPow,
			code.OpBand,
			code.OpBor,
			code.OpBxor,
			code.OpShl,
			code.OpShr:
			if !e.numeral() {
				ls.fs.code.expr2rk(ls.fs, e)
			}
		default:
			ls.fs.code.expr2rk(ls.fs, e)
	}
}

// postfix finalizes code for a binary operation, after reading the 2nd operand.
// For '(a .. b .. c)' which is '(a .. (b .. c))', because concatenation is
// right associative, merge the second concat into first one.
func (p *parser) postfix(ls *lexical, op code.Op, e1, e2 *expr, line int) {
	defer un(trace(ls, "Postfix"))
	switch fs := ls.fs; op {
		case code.OpConcat:
			if fs.code.expr2val(fs, e2); e2.kind == vreloc {
				if inst := fs.instr(e2); inst.code.Code() == code.CONCAT {
					ls.assert(e1.info == inst.code.B() - 1)
					fs.code.freeexpr(fs, e1)
					inst.code.SetB(uint(e1.info))
					e1.info = e2.info
					e1.kind = vreloc
					return
				}
			}
			fs.code.expr2next(fs, e2) // operand must be on the 'stack'
			fs.code.codeBinary(fs, code.CONCAT, e1, e2, line)
			return

		case code.OpAnd:
			ls.assert(e1.t == noJump) // list closed by 'infix'
			fs.code.dischargeVars(fs, e2)
			e2.f = concatJumps(fs, e2.f, e1.f)
			*e1 = *e2
			return

		case code.OpOr:
			ls.assert(e1.f == noJump) // list closed by 'infix'
			fs.code.dischargeVars(fs, e2)
			e2.t = concatJumps(fs, e2.t, e1.t)
			*e1 = *e2
			return

		case code.OpAdd,
			code.OpSub,
			code.OpMul,
			code.OpDivF,
			code.OpDivI,
			code.OpMod,
			code.OpPow,
			code.OpBand,
			code.OpBor,
			code.OpBxor,
			code.OpShl,
			code.OpShr:

			if !constfold(fs, op, e1, e2) {
				fs.code.codeBinary(fs, code.Opcode(op-1)+code.ADD, e1, e2, line)
			}
			return

		case code.OpEq,
			code.OpLt,
			code.OpLe,
			code.OpNe,
			code.OpGt,
			code.OpGe:

			p.compare(ls, op, e1, e2)
			return
		
	}
	panic("unreachable")
}

// suffixedexp -> primaryexp { '.' NAME | '[' exp ']' | ':' NAME funcargs | funcargs }
func (p *parser) suffixed(ls *lexical, e *expr) {
	defer un(trace(ls, "SuffixedExpr"))
	line := ls.line
	p.primary(ls, e)
	fs := ls.fs
	for {
		switch ls.token.char {
			case '(', tString, '{':
				fs.code.expr2next(fs, e)
				p.arguments(ls, e, line)
			case '[':
				fs.code.expr2ru(fs, e)
				var k expr
				p.index(ls, &k)
				fs.code.codeIndex(fs, e, &k)
			case ':':
				ls.next()
				var k expr
				k.init(vconst, fs.constant(p.ident(ls)))
				fs.code.codeSelf(fs, e, &k)
				p.arguments(ls, e, line)
			case '.':
				p.selector(ls, e)
			default:
				return
		}
	}
}

// fieldsel -> ['.' | ':'] NAME
func (p *parser) selector(ls *lexical, e *expr) {
	ls.fs.code.expr2ru(ls.fs, e)
	ls.next() // skip '.' or ':'
	var k expr
	k.init(vconst, ls.fs.constant(p.ident(ls)))
	ls.fs.code.codeIndex(ls.fs, e, &k)
}

// compare emits code for comparisons.
//
// 'e1' was already put in R/K form by 'infix'.
func (p *parser) compare(ls *lexical, op code.Op, e1, e2 *expr) {
	var ( rk1, rk2 int )
	if e1.kind == vconst {
		rk1 = code.RK(e1.info)
	} else {
		ls.assert(e1.kind == vnonreloc)
		rk1 = e1.info
	}
	rk2 = ls.fs.code.expr2rk(ls.fs, e2)
	ls.fs.code.freeexpr2(ls.fs, e1, e2)
	switch fs := ls.fs; op {
		case code.OpGt, code.OpGe:
			// '(a > b)' ==> '(b < a)'; '(a >= b)' ==> '(b <= a)'
			e1.info = condJump(fs, code.Opcode(op-code.OpNe)+code.EQ, 1, rk2, rk1)
		case code.OpNe:
			// '(a ~= b)' ==> 'not (a == b)'
			e1.info = condJump(fs, code.EQ, 0, rk1, rk2)
		default:
			// '==', '<', '<=' use their own opcodes
			e1.info = condJump(fs, code.Opcode(op-code.OpEq)+code.EQ, 1, rk1, rk2)
	}
	e1.kind = vjump
}

// primaryexp -> NAME | '(' expr ')'
func (p *parser) primary(ls *lexical, e *expr) {
	defer un(trace(ls, "PrimaryExpr"))
	switch ls.token.char {
		case tName:
			p.variable(ls, e)
		case '(':
			line := ls.line
			ls.next()
			p.expr(ls, e)
			ls.match(')', '(', line)
			ls.fs.code.dischargeVars(ls.fs, e)
		default:
			ls.syntaxErr("unexpected symbol")
	}
}

// simpleexp -> FLT | INT | STRING | NIL | TRUE | FALSE | ... | constructor | FUNCTION body | suffixedexp
func (p *parser) simple(ls *lexical, e *expr) {
	switch ls.token.char {
		case tFunction:
			ls.next()
			p.funcbody(ls, e, false, ls.line)
			return
		case tString:
			e.init(vconst, ls.fs.constant(ls.token.sval))
		case tFloat:
			e.init(vfloat, 0)
			e.nval = ls.token.nval
		case tFalse:
			e.init(vfalse, 0)
		case tTrue:
			e.init(vtrue, 0)
		case tDots:
			if !ls.fs.fn.Vararg {
				ls.syntaxErr("cannot use '...' outside a vararg function")
			}
			e.init(vvararg, ls.fs.code.codeABC(ls.fs, code.VARARG, 0, 1, 0))
		case tInt:
			e.init(vint, 0)
			e.ival = ls.token.ival
		case tNil:
			e.init(vnil, 0)
		case '{':
			p.constructor(ls, e)
			return
		default:
			p.suffixed(ls, e)
			return
	}
	ls.next()
}

// cond -> exp
func (p *parser) condition(ls *lexical) int {
	var e expr
	if p.expr(ls, &e); e.kind == vnil {
		e.kind = vfalse
	}
	goIfTrue(ls.fs, &e)
	return e.f
}

func (p *parser) variable(ls *lexical, e *expr) {
	defer un(trace(ls, "Variable"))
	name := p.ident(ls)
	search(ls.fs, name, e, true)
	if e.kind == vvoid {
		search(ls.fs, "_ENV", e, true)
		ls.assert(e.kind != vvoid)
		var k expr
		k.init(vconst, ls.fs.constant(name))
		ls.fs.code.codeIndex(ls.fs, e, &k)
	}
}

func (p *parser) arguments(ls *lexical, e *expr, line int) {
	defer un(trace(ls, "FuncArgs"))
	var (
		fs = ls.fs
		args expr
		argc int
		base int
	)
	switch ls.token.char {
		case tString:
			args.init(vconst, fs.constant(ls.token.sval))
			ls.next()
		case '(':
			if ls.next(); ls.token.char == ')' {
				args.kind = vvoid
			} else {
				p.exprs(ls, &args)
				fs.code.returnX(fs, &args)
			}
			ls.match(')', '(', line)
		case '{':
			p.constructor(ls, &args)
		default:
			ls.syntaxErr("function arguments expected")
	}
	ls.assert(e.kind == vnonreloc)
	if base = e.info; args.retsX() {
		argc = multRet
	} else {
		if args.kind != vvoid {
			ls.fs.code.expr2next(fs, &args)
		}
		argc = fs.free - (base + 1)
	}
	e.init(vcall, fs.code.codeABC(fs, code.CALL, base, argc+1, 2))
	fs.fixLine(line)
	fs.free = base+1
}

// index -> '[' expr ']'
func (p *parser) index(ls *lexical, e *expr) {
	ls.next() // skip the '['
	p.expr(ls, e)
	ls.fs.code.expr2val(ls.fs, e)
	ls.expect(']')
	ls.next()
}

func (p *parser) ident(ls *lexical) string {
	defer un(trace(ls, "Ident"))
	ls.expect(tName)
	s := ls.token.sval
	ls.next()
	return s
}

// constfold tries to "constant-fold" an operation; return true
// iff successful. In this case, 'e1' holds the final result.
func constfold(fs *function, op code.Op, e1, e2 *expr) bool {
	if n1, ok := toNumeral(e1); ok {
		if n2, ok := toNumeral(e2); ok {
			if checkOp(op, n1, n2) {
				if v, ok := eval(op, n1, n2); ok {
					if i, ok := v.(int64); ok {
						e1.kind = vint
						e1.ival = i
						return true
					}
					// folds neither NaN nor 0.0 (to avoid problems with -0.0)
					if n, ok := toFloat(v); ok && (n != 0) && !isNaN(n) {
						e1.kind = vfloat
						e1.nval = n
						return true
					}
				}
			}
		}
	}
	// non-numeric operands are not safe to fold.
	return false
}

// checkOp returns false if folding can raise an error. Bitwise
// operations need oparands that are convertible to integers;
// division operations cannot have 0 as a divisor.
func checkOp(op code.Op, v1, v2 code.Const) bool {
	switch op {
		case code.OpBand,
			code.OpBor,
			code.OpBxor,
			code.OpShl,
			code.OpShr,
			code.OpBnot:
			// conversion errors
			_, ok1 := toInt(v1)
			_, ok2 := toInt(v2)
			return ok1 && ok2	
		case code.OpDivI,
			code.OpDivF,
			code.OpMod:
			n, ok := toFloat(v2)
			return ok && (n != 0)
	}
	return true
}

func toNumeral(e *expr) (code.Const, bool) {
	if !e.numeral() {
		return nil, false
	}
	switch e.kind {
		case vfloat:
			return e.nval, true
		case vint:
			return e.ival, true
	}
	panic("unreachable")
}

func toFloat(v code.Const) (float64, bool) {
	switch v := v.(type) {
		case float64:
			return v, true
		case int64:
			return float64(v), true
	}
	return 0, false
}

func toInt(v code.Const) (int64, bool) {
	switch v := v.(type) {
		case string:
			panic("str2int: TODO!")
		case float64:
			if i := int64(v); float64(i) == float64(v) {
				return i, true
			}
		case int64:
			return v, true
	}
	return 0, false
}

func isNaN(f64 float64) bool {
	return math.IsNaN(f64)
}