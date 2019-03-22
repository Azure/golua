package luac

import (
	"fmt"
	"os"
	"github.com/fibonacci1729/golua/lua/code"
)

var _ = fmt.Println
var _ = os.Exit

type constructor struct {
	pending int   // number of array elements pending to be stored
	nh, na  int   // total number of 'record' / 'array' elements
	v 	    expr  // last list item read
	t 	    *expr // table descriptor
}

// constructor -> '{' [ field { sep field } [sep] ] '}' sep -> ',' | ';
func (p *parser) constructor(ls *lexical, e *expr) {
	pc := ls.fs.code.codeABC(ls.fs, code.NEWTABLE, 0, 0, 0)
	cc := &constructor{t: e}
	line := ls.line
	cc.t.init(vreloc, pc) 
	cc.v.init(vvoid, 0) // no value (yet)
	ls.fs.code.expr2next(ls.fs, cc.t) // fix it at stack top
	ls.expect('{')
	ls.next()
	for {
		ls.assert(cc.v.kind == vvoid || cc.pending > 0)
		if ls.token.char == '}' {
			break
		}
		closeListField(ls.fs, cc)
		p.field(ls, cc)
		if !ls.test(',') && !ls.test(';') {
			break
		}
	}
	ls.match('}', '{', line)
	lastListField(ls.fs, cc)
	ls.fs.instrs[pc].code.SetB(uint(int2fb(cc.na))) // set initial array size
	ls.fs.instrs[pc].code.SetC(uint(int2fb(cc.nh))) // set initial table size
}

// field -> listfield | recfield
func (p *parser) field(ls *lexical, cc *constructor) {
	switch ls.token.char {
		case tName: // may be 'listfield' or 'recfield'
			if ls.peek() != '=' { // expression?
				p.listField(ls, cc)
			} else {
				p.hashField(ls, cc)
			}
		case '[':
			p.hashField(ls, cc)
		default:
			p.listField(ls, cc)
	}
}

// listfield -> exp
func (p *parser) listField(ls *lexical, cc *constructor) {
	p.expr(ls, &cc.v)
	checkLimit(ls.fs, cc.na, maxInt, "items in a constructor")
	cc.na++
	cc.pending++
}

// recfield -> (NAME | '['exp1']') = exp1
func (p *parser) hashField(ls *lexical, cc *constructor) {
	register := ls.fs.free
	var ( k, v expr )
	if ls.token.char == tName {
		checkLimit(ls.fs, cc.nh, maxInt, "items in a constructor")
		k.init(vconst, ls.fs.constant(p.ident(ls)))
	} else { // ls.token.char == '['
		p.index(ls, &k)
	}
	cc.nh++
	ls.expect('=')
	ls.next()
	rkkey := ls.fs.code.expr2rk(ls.fs, &k)
	p.expr(ls, &v)
	ls.fs.code.codeABC(ls.fs, code.SETTABLE, cc.t.info, rkkey, ls.fs.code.expr2rk(ls.fs, &v))
	ls.fs.free = register // free registers
}

func closeListField(fs *function, cc *constructor) {
	if cc.v.kind == vvoid { // there is no list item
		return
	}
	fs.code.expr2next(fs, &cc.v)
	cc.v.kind = vvoid
	if cc.pending == fieldsPerFlush {
		fs.code.codeSetList(fs, cc.t.info, cc.na, cc.pending) // flush
		cc.pending = 0
	}
}

// int2fb

func lastListField(fs *function, cc *constructor) {
	if cc.pending == 0 {
		return
	}
	if cc.v.retsX() {
		fs.code.returnX(fs, &cc.v)
		fs.code.codeSetList(fs, cc.t.info, cc.na, multRet)
		cc.na-- // do not count last expression (unknown number of elements)
		return
	}
	if cc.v.kind != vvoid {
		fs.code.expr2next(fs, &cc.v)
	}
	fs.code.codeSetList(fs, cc.t.info, cc.na, cc.pending)
}