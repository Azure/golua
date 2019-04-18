package luac

import (
	"fmt"
	"os"

	"github.com/Azure/golua/lua/code"
)

var _ = fmt.Println
var _ = os.Exit

// forbody -> DO block
func (p *parser) forbody(ls *lexical, base, line, varN int, numeric bool) {
	adjustLocals(ls.fs, 3) // control variables
	ls.expect(tDo)
	ls.next()

	var prep int
	if numeric {
		prep = ls.fs.code.codeAsBx(ls.fs, code.FORPREP, base, noJump)
	} else {
		prep = ls.fs.code.codeJump(ls.fs)
	}
	ls.fs.enter(false) // scope for declared variables
	adjustLocals(ls.fs, varN)
	ls.fs.reserve(varN)
	p.block(ls)
	ls.fs.leave() // end of scope for declared variables
	patch2here(ls.fs, prep)
	var end int
	if numeric { // numeric for?
		end = ls.fs.code.codeAsBx(ls.fs, code.FORLOOP, base, noJump)
	} else { // generic for?
		ls.fs.code.codeABC(ls.fs, code.TFORCALL, base, 0, varN)
		ls.fs.fixLine(line)
		end = ls.fs.code.codeAsBx(ls.fs, code.TFORLOOP, base+2, noJump)
	}
	patchList(ls.fs, end, prep+1)
	ls.fs.fixLine(line)
}

// fornum -> NAME = exp1,exp1[,exp1] forbody
func (p *parser) forlist(ls *lexical, name string, line int) {
	base := ls.fs.free
	ls.fs.declare("(for index)")
	ls.fs.declare("(for limit)")
	ls.fs.declare("(for step)")
	ls.fs.declare(name)
	// initial value
	ls.expect('=')
	ls.next()
	p.expr1(ls)
	// limit value
	ls.expect(',')
	ls.next()
	p.expr1(ls)
	// optional step
	if ls.test(',') {
		p.expr1(ls)
	} else {
		// default step = 1
		ls.fs.code.codeKst(ls.fs, ls.fs.free, ls.fs.constant(int64(1)))
		ls.fs.reserve(1)
	}
	p.forbody(ls, base, line, 1, true)
}

// forlist -> NAME {,NAME} IN explist forbody
func (p *parser) foriter(ls *lexical, name string) {
	var (
		base = ls.fs.free
		varN = 4 // generator, state, control, plus at least one declared variable
	)
	// create control variables
	ls.fs.declare("(for generator)")
	ls.fs.declare("(for state)")
	ls.fs.declare("(for control)")
	// create declared variables
	ls.fs.declare(name)

	for ls.test(',') {
		ls.fs.declare(p.ident(ls))
		varN++
	}
	ls.expect(tIn)
	ls.next()
	line := ls.line
	var e expr
	adjustAssign(ls.fs, 3, p.exprs(ls, &e), &e)
	ls.fs.checkstack(3) // extra space to call generator
	p.forbody(ls, base, line, varN-3, false)
}

// forstat -> FOR (fornum | forlist) END
func (p *parser) forloop(ls *lexical, line int) {
	ls.fs.enter(true) // scope for loop and control variables
	ls.next()         // skip 'for'
	switch name := p.ident(ls); ls.token.char {
	case ',', tIn:
		p.foriter(ls, name)
	case '=':
		p.forlist(ls, name, line)
	default:
		ls.syntaxErr("'=' or 'in' expected")
	}
	ls.match(tEnd, tFor, line)
	ls.fs.leave() // loop scope ('break' jumps to this point)
}

// repeatstat -> REPEAT block UNTIL cond
func (p *parser) repeat(ls *lexical, line int) {
	init := ls.fs.pclabel()
	ls.fs.enter(true)  // loop block
	ls.fs.enter(false) // scope block
	ls.next()          // skip REPEAT
	p.stmts(ls)
	ls.match(tUntil, tRepeat, line)
	exit := p.condition(ls) // read condition (inside scope block)
	if ls.fs.block.hasup {  // upvalues?
		patchClose(ls.fs, exit, ls.fs.block.active)
	}
	ls.fs.leave()                // finish scope
	patchList(ls.fs, exit, init) // close the loop
	ls.fs.leave()                // finish loop
}

// whilestat -> WHILE cond DO block END
func (p *parser) while(ls *lexical, line int) {
	ls.next() // skip WHILE
	var (
		init = ls.fs.pclabel()
		exit = p.condition(ls)
	)
	ls.fs.enter(true)
	ls.expect(tDo)
	ls.next()
	p.block(ls)
	jumpTo(ls.fs, init)
	ls.match(tEnd, tWhile, line)
	ls.fs.leave()
	patch2here(ls.fs, exit) // false conditions finish the loop
}
