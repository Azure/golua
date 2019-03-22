package luac

import (
	"math"
	"fmt"
	"os"
	"github.com/Azure/golua/lua/code"
)

var _ = fmt.Println
var _ = os.Exit

type (
	assignment struct {
		prev *assignment
		expr expr
	}

	// function state needed to generate code for a given function.
	function struct {
		consts  map[code.Const]int
		parent  *function  // enclosing function
		instrs  []*instr   // instructions
		block   *block     // chain of current block
		code 	*coder 	   // code generation state
		ls 		*lexical   // lexical state
		fn      *code.Proto // current function header
		// transient compile state
		target   int // 'label' of last 'jump label'
		jumppc   int // list of pending jumps to 'pc'
		active   int // number of active local variables
		local0   int // index of first local var (in lexical state)
		nups     int // number of upvalues
		free     int // first free register
		pc 		 int // next position to code (equivalent to 'ncode')
	}

	// variable is a description of an active local variable
	variable struct {
		index int // variable index in stack
	}

	// block is a list of active blocks (Control block).
	block struct {
		parent *block // parent block
		active int    // # of active locals outside the block
		label0 int    // index of first label in thi block
		gotos0 int    // index of first pending goto in this block
		hasup  bool   // true if some variable in the block is an upvalue
		loop   bool   // true if 'block' is a loop
	}

	// label is a description of pending goto statements and label statements
	label struct {
		label string // label identifier
		level int    // local level where it appears in current block
		line  int    // line where it appeared
		pc    int 	 // position in code
	}
)

func (fs *function) numberk(n float64) int {
	if n == 0.0 || math.IsNaN(n) {
		return fs.constant(math.Float64bits(n))
	}
	return fs.constant(n)
}

func (fs *function) constant(value interface{}) int {
	if fs.consts == nil {
		fs.consts = make(map[code.Const]int)
	}
	var k code.Const
	switch v := value.(type) {
		case float64,
			string,
			int64,
			bool,
			nil:
			k = v
		case uint64:
			k = float64(v)
		default:
			fs.ls.syntaxErr(fmt.Sprintf("unexpected constant type %T", v))
	}
	if i, ok := fs.consts[k]; ok {
		return i
	}
	fs.fn.Consts = append(fs.fn.Consts, k)
	fs.consts[k] = len(fs.fn.Consts) - 1
	return fs.consts[k]
}

func (fs *function) checkstack(n int) {
	if stack := fs.free + n; stack > fs.fn.StackN {
		if stack >= maxRegs {
			fs.ls.syntaxErr("function or expression needs too many registers")
		}
		fs.fn.StackN = stack
	}
}

func (fs *function) reserve(n int) {
	fs.checkstack(n)
	fs.free += n
}

func (fs *function) instr(e *expr) *instr {
	return fs.instrs[e.info]
}

func (fs *function) enter(loop bool) *block {
	fs.block = &block{
		label0: len(fs.ls.labels),
		gotos0: len(fs.ls.gotos),
		active: fs.active,
		parent: fs.block,
		loop:   loop,
	}
	fs.ls.assert(fs.free == fs.active)
	return fs.block
}

func (fs *function) leave() {
	var (
		b = fs.block
		ls = fs.ls
	)
	if b.parent != nil && b.hasup {
		// create a 'jump to here' to close upvalues
		jmp := fs.code.codeJump(fs)
		patchClose(fs, jmp, b.active)
		patch2here(fs, jmp)
	}
	if b.loop {
		// create a label named 'break' to resolve
		// break statements, and close pending breaks.
		lbl := fs.label(&ls.labels, "break", 0, fs.pc)
		findgotos(ls, ls.labels[lbl])
	}
	fs.block = b.parent
	removeVars(fs, b.active)
	fs.ls.assert(b.active == fs.active)
	fs.free = fs.active // free registers
	switch {
	 	case b.parent != nil: // inner block?
			// update pending gotos to outer block
			movegotos(fs, b)
		case b.gotos0 < len(ls.gotos): // pending gotos in outer block?
			ls.undefGotoErr(ls.gotos[b.gotos0]) // error
	}
	ls.labels = ls.labels[:b.label0] // remove local labels
}

func (fs *function) open(ls *lexical) *function {
	fs.parent = ls.fs
	fs.local0 = len(ls.active)
	fs.jumppc = noJump
	fs.target = 0
	fs.nups   = 0
	fs.free   = 0
	fs.pc     =	0
	fs.ls     = ls
	fs.fn = &code.Proto{
		Source: ls.name,
		StackN: 2, // registers 0/1 are always valid
	}
	fs.enter(false)
	return fs
}

func (fs *function) close() *function {
	fs.code.codeReturn(fs, 0, 0)
	fs.leave()
	fs.ls.assert(fs.block == nil)

	fs.fn.Instrs = make([]code.Instr, len(fs.instrs))
	fs.fn.PcLine = make([]int32, len(fs.instrs))

	for i := 0; i < len(fs.instrs); i++ {
		fs.fn.Instrs[i] = *fs.instrs[i].code
		fs.fn.PcLine[i] = int32(fs.instrs[i].line)
	}

	return fs.parent
}

func (fs *function) label(list *[]*label, name string, line, pc int) int {
	*list = append(*list, &label{
		level: fs.active,
		label: name,
		line:  line,
		pc:    pc,
	})
	return len(*list)-1
}

// searchvar
func (fs *function) search(name string) int {
	for i := fs.active - 1; i >= 0; i-- {
		if fs.local(i).Name == name {
			return i
		}
	}
	return -1
}

// new_localvar
func (fs *function) declare(name string) {
	fs.fn.Locals = append(fs.fn.Locals, &code.Local{Name: name})
	checkLimit(fs, len(fs.ls.active) + 1 - fs.local0, maxVars, "local variables")
	fs.ls.active = append(fs.ls.active, &variable{len(fs.fn.Locals)-1})
}

// getlocvar
func (fs *function) local(i int) *code.Local {
	// fmt.Printf("local(index=%d, first=%d, total=%d)\n", i, fs.local0, len(fs.ls.active))
	index := fs.ls.active[fs.local0 + i].index
	fs.ls.assert(index < len(fs.fn.Locals))
	return fs.fn.Locals[index]
}

func search(fs *function, name string, e *expr, base bool) *expr {
	if fs == nil {
		return e.init(vvoid, 0)
	}
	if i := fs.search(name); i >= 0 {
		e.init(vlocal, i)
		if !base {
			fs.markup(i)
		}
		return e
	}
	i := fs.findup(name)
	if i < 0 {
		if e = search(fs.parent, name, e, false); e.kind == vvoid {
			return e
		}
		i = fs.makeup(name, e)
	}
	return e.init(vupval, i)
}

// searchupvalue
func (fs *function) findup(name string) int {
	for i, up := range fs.fn.UpVars {
		if up.Name == name {
			return i
		}
	}
	return -1
}

// newupvalue
func (fs *function) makeup(name string, e *expr) int {
	checkLimit(fs, len(fs.fn.UpVars) + 1, maxUp, "upvalues")
	fs.fn.UpVars = append(fs.fn.UpVars, &code.UpVar{
		Stack: (e.kind == vlocal),
		Index: e.info,
		Name:  name,
	})
	fs.nups++
	return len(fs.fn.UpVars) - 1
}

// markup marks the block where variable at the given level was defined
// (to emit close instructions later).
//
// markupval
func (fs *function) markup(level int) {
	b := fs.block
	for b.active > level {
		b = b.parent 
	}
	b.hasup = true
}

func (fs *function) fixLine(line int) {
	fs.instrs[fs.pc-1].line = int32(line)
}

// pclabel returns the current 'pc' and marks it as a jump target (to avoid
// optimizations with consecutive instructions not in the same basic block).
func (fs *function) pclabel() int {
	fs.target = fs.pc
	return fs.pc
}