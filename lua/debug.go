package lua

import (
	// "strings"
	"fmt"
	// "os"

	// "github.com/fibonacci1729/golua/lua/code"
)

// var _ = fmt.Println
// var _ = os.Exit

type Debug interface {
	// Options(what string) Debug
	CurrentLine() int
	SourceID() string
	ShortSrc() string
	ParamN() int
	UpVarN() int
	Name() string
	Kind() string
	Span() (int, int)
	Vararg() bool
	Tailcall() bool
	Where() string
}

type debug struct {
	ci *call

	source   string // (S)
	short    string // (S)
	event    int
	line     int    // (l) (currentline)
	name     string // (n) (name)
	what     string // (n) 'global', 'local', 'field', 'method' (namewhat)
	kind     string // (S) 'Lua', 'Go', 'main', 'tail' (what)
	span     [2]int // (S) (linedefined/lastlinedefined)
	upvarN   int    // (u) number of upvalues
	paramN   int    // (u) number of parameters
	vararg   bool   // (u) (isvararg)
	tailcall bool   // (t) (istailcall)
}

func (dbg *debug) CurrentLine() int { return dbg.line }
func (dbg *debug) SourceID() string { return dbg.source }
func (dbg *debug) ShortSrc() string { return dbg.short }
func (dbg *debug) ParamN() int { return dbg.paramN }
func (dbg *debug) UpVarN() int { return dbg.upvarN }
func (dbg *debug) Name() string { return dbg.name }
func (dbg *debug) Kind() string { return dbg.kind }
func (dbg *debug) Span() (int, int) { return dbg.span[0], dbg.span[1] }
func (dbg *debug) Vararg() bool { return dbg.vararg }
func (dbg *debug) Tailcall() bool { return dbg.tailcall }

func (dbg *debug) Where() string {
	if dbg != nil {
		if dbg.CurrentLine() > 0 {
			return fmt.Sprintf("%s:%d: ",
				dbg.ShortSrc(),
				dbg.CurrentLine())
		}
	}
	return ""
}

// func funcNameFromGlobals(fr *frame, dbg *debug) (name string) {
// 	var (
// 		ld = fr.ls.global.loaded
// 		found bool
// 	)
// 	if fs := fr.call.fs; fs != nil && fs.closure != nil {
// 		name, found = searchField(ld, fs.closure, 2)
// 	} else {
// 		if fn, ok := fr.call.fn.(Value); ok {
// 			name, found = searchField(ld, fn, 2)
// 		}
// 	}
// 	if found && strings.HasPrefix(name, "_G.") {
// 		name = name[3:]
// 	}
// 	return name
// }

// func searchField(env, fn Value, level int) (name string, found bool) {
// 	if tbl, ok := env.(*Table); ok && level > 0 {
// 		tbl.foreach(func(k, v Value) bool {
// 			if k, isstr := k.(String); isstr {
// 				if found = equals(fn, v); found {
// 					name = string(k)
// 					return false
// 				}
// 			}
// 			var s string
// 			s, found = searchField(v, fn, level-1)
// 			if found {
// 				name = fmt.Sprintf("%s.%s", k, s)
// 				return false
// 			}
// 			return true
// 		})
// 	}
// 	return name, found
// }

// func funcNameFromCode(ci *call) (name, what string) {
// 	if fn, ok := ci.fn.(*Proto); ok {
// 		if ci.flag & hooked != 0 {
// 			return "?", "hook"
// 		}
// 		var (
// 			inst = fn.Instrs[ci.pc]
// 			meta event
// 		)
// 		switch inst.Code() {
// 			case code.SELF, code.GETTABUP, code.GETTABLE:
// 				meta = _index
// 			case code.SETTABUP, code.SETTABLE:
// 				meta = _newindex
// 			case code.CALL, code.TAILCALL:
// 				return objectName(ci, fn, ci.pc, inst.A())
// 			case code.TFORCALL:
// 				return "for iterator", "for iterator"
// 			case code.IDIV,
// 				code.DIV,
// 				code.ADD,
// 				code.SUB,
// 				code.MUL,
// 				code.MOD,
// 				code.POW,
// 				code.SHL,
// 				code.SHR,
// 				code.BOR,
// 				code.BXOR,
// 				code.BAND:
// 				meta = event(inst.Code()-code.ADD) + _add
// 			case code.CONCAT:
// 				meta = _concat
// 			case code.BNOT:
// 				meta = _bnot
// 			case code.UNM:
// 				meta = _unm
// 			case code.LEN:
// 				meta = _len
// 			case code.EQ:
// 				meta = _eq
// 			case code.LT:
// 				meta = _lt
// 			case code.LE:
// 				meta = _le
// 			default:
// 				return
// 		}
// 		return events[meta], "metamethod"
// 	}
// 	return "", ""
// }

// func objectName(ci *call, fn *Proto, lastpc, reg int) (name, what string) {
// 	if name = localName(fn, lastpc, reg+1); name != "" {
// 		return name, "local"
// 	}
// 	// try symbolic execution
// 	if pc := findSetRegister(ci, fn, lastpc, reg); pc != -1 { // instruction found?
// 		switch inst := fn.Instrs[pc]; inst.Code() {
// 			case code.GETTABUP:
// 				var (
// 					t = inst.B() // table index
// 					k = inst.C() // key index
// 				)
// 				if what = "?"; t < len(fn.UpVars) {
// 					if up := fn.UpVars[t]; up.Name != "" {
// 						if up.Name == EnvID {
// 							what = "global"
// 						} else {
// 							what = "field"
// 						}
// 					}
// 				}
// 				return findNameRK(ci, fn, pc, k), what
// 			case code.GETTABLE:
// 				var (
// 					t = inst.B() // table index
// 					k = inst.C() // key index
// 				)
// 				switch localName(fn, pc, t+1) {
// 					case EnvID:
// 						what = "global"
// 					default:
// 						what = "field"
// 				}
// 				return findNameRK(ci, fn, pc, k), what
// 			case code.GETUPVAL:
// 				if name = "?"; inst.B() < len(fn.UpVars) {
// 					if up := fn.UpVars[inst.B()]; up.Name != "" {
// 						name = up.Name
// 					}
// 				}
// 				return name, "upvalue"
// 			case code.LOADKX:
// 				kst := fn.kst(fn.Instrs[pc+1].AX())
// 				if s, ok := kst.(String); ok {
// 					return string(s), "constant"
// 				}
// 			case code.LOADK:
// 				kst := fn.kst(inst.BX())
// 				if s, ok := kst.(String); ok {
// 					return string(s), "constant"
// 				}
// 			case code.SELF:
// 				return findNameRK(ci, fn, pc, inst.C()), "method"
// 			case code.MOVE:
// 				if inst.B() < inst.A() { // move from 'b' to 'a'
// 					return objectName(ci, fn, pc, inst.B()) // get name for 'b'
// 				}
// 		}
// 	}
// 	return "", ""
// }

// // localName looks for the n-th local variable at line 'line' in function 'func'.
// //
// // Returns the local variable name or "" if not found.
// func localName(fn *Proto, pc, n int) (name string) {
// 	for i := 0; i < len(fn.Locals) && fn.Locals[i].Live <= int32(pc); i++ {
// 		if int32(pc) < fn.Locals[i].Dead { // is variable active?
// 			if n--; n == 0 {
// 				return fn.Locals[i].Name
// 			}
// 		}
// 	}
// 	// not found
// 	return ""
// }

// // findNameRK finds a name for the RK value 'rk'.
// func findNameRK(ci *call, fn *Proto, pc, rk int) (name string) {
// 	if code.IsKst(rk) { // is 'rk' a constant?
// 		if s, ok := fn.kst(code.ToKst(rk)).(String); ok { // literal constant?
// 			return string(s) // it is its own name
// 		}
// 		// else no reasonable name found
// 		return ""
// 	}
// 	// else 'rk' is a register
// 	name, what := objectName(ci, fn, pc, rk)
// 	if what == "constant" {
// 		return name
// 	}
// 	return "?"
// }

// // findSetRegister tries to find the last instruction before 'lastpc' that modified register 'reg'.
// func findSetRegister(ci *call, fn *Proto, lastpc, reg int) (pc int) {
// 	var (
// 		set = -1 // keep last instruction that changed 'reg'
// 		jmp = 0  // any code before this address is conditional
// 	)
// 	filterPC := func(pc, jmp int) int {
// 		if pc < jmp { // is code conditional (inside a jump)?
// 			return -1 // cannt know who sets that register
// 		}
// 		return pc // current position sets that register
// 	}
// 	for pc := 0; pc < lastpc; pc++ {
// 		switch inst := fn.Instrs[pc]; inst.Code() {
// 			case code.CALL, code.TAILCALL:
// 				if reg >= inst.A() { // affect all registers above base
// 					set = filterPC(pc, jmp)
// 				}
// 			case code.TFORCALL:
// 				if reg >= inst.A() + 2 { // affect all registers above its base
// 					set = filterPC(pc, jmp)
// 				}
// 			case code.LOADNIL:
// 				if a, b := inst.A(), inst.B(); reg >= a && reg <= a + b {
// 					// set register from 'a' to 'a+b'
// 					set = filterPC(pc, jmp)
// 				}
// 			case code.JMP:
// 				// jump is forward and do not skip 'lastpc'?
// 				if dst := pc + 1 + inst.SBX(); dst > pc && dst <= lastpc {
// 					if dst > jmp {
// 						jmp = dst // update jump target
// 					}
// 				}
// 			default:
// 				if inst.Code().Mask().SetA() && (reg == inst.A()) {
// 					// any instruction that set A
// 					set = filterPC(pc, jmp)
// 				}
// 		}
// 	}
// 	return set
// }

// func sourceInfo(dbg *debug) *debug {
// 	if fn, ok := dbg.ci.fn.(*Func); ok {
// 		if dbg.source = fn.Source; fn.Source == "" {
// 			dbg.source = "=?"
// 		}
// 		if dbg.kind = "Lua"; fn.SrcPos == 0 {
// 			dbg.kind = "main"
// 		}
// 		dbg.span[0] = fn.SrcPos
// 		dbg.span[1] = fn.EndPos
// 		dbg.line = -1
// 		if len(fn.PcLine) > 0 {
// 			dbg.line = int(fn.PcLine[ci.pc-1])
// 		}
// 	} else {
// 		dbg.source  = "=[Go]"
// 		dbg.kind    = "Go"
// 		dbg.line    = -1
// 		dbg.span[0] = -1
// 		dbg.span[1] = -1
// 	}
// 	dbg.short = chunkID(dbg.source)
// 	return dbg
// }

// func funcInfo(dbg *debug) *debug {
// 	if fn, ok := ci.fn.(*Proto); ok {
// 		dbg.tailcall = ci.flag & tailcall != 0
// 		dbg.upvarN   = len(fn.UpVars)
// 		dbg.paramN   = fn.ParamN
// 		dbg.vararg   = fn.Vararg
// 	} else {
// 		dbg.vararg = true
// 	}
// 	return dbg
// }

// func funcName(dbg *debug) *debug {
// 	if caller := ci.fr.call; isLua(caller) {
// 		if (ci.flag & tailcall == 0) {
// 			dbg.name, dbg.what = funcNameFromCode(caller)
// 		}
// 	}
// 	if dbg.what == "" {
// 		dbg.name = "<unset>"
// 	}
// 	return dbg
// }