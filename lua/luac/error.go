package luac

import (
	"fmt"
	"strings"

	"github.com/Azure/golua/lua/code"
)

// undefGotoErr generates an error for an undefined 'goto'; choose appropriate
// message when label name is a reserved word (which can only be 'break').
func (ls *lexical) undefGotoErr(g *label) {
	msgFmt := "no visible label '%s' for <goto> at line %d"
	if _, reserved := keywords[g.label]; reserved {
		msgFmt = "<%s> at line %d not inside a loop"
	}
	ls.semanticErr(fmt.Sprintf(msgFmt, g.label, g.line))
}

func (ls *lexical) semanticErr(msg string) {
	ls.token.char = 0 // remove "near <token>" from final message
	ls.syntaxErr(msg)
}

func (ls *lexical) syntaxErr(msg string) {
	ls.scanErr(msg, ls.token.char)
}

func (ls *lexical) expectErr(tok rune) {
	ls.syntaxErr(fmt.Sprintf("%s expected", ls.tok2str(tok)))
}

func (s *lexical) escapeErr(chars []rune, message string) {
	s.buffer.Reset()
	s.save('\\')
	for _, r := range chars {
		if r == eof {
			break
		}
		s.save(r)
	}
	s.scanErr(message, tString)
}

func (ls *lexical) numberErr() {
	ls.scanErr("malformed number", tFloat)
}

func (ls *lexical) scanErr(msg string, tok rune) {
	if msg = addErrInfo(msg, ls.name, ls.line); tok != 0 {
		msg = fmt.Sprintf("%s near %s", msg, ls.tok2str(tok))
	}
	panic(code.Error(msg))
}

func (ls *lexical) checkEsc(cond bool, msg string) {
	if !cond {
		if ls.char != eof {
			ls.consume() // add current to buffer for error message
		}
		ls.scanErr(msg, tString)
	}
}

func (ls *lexical) limitErr(fs *function, limit int, what string) {
	var (
		where = "main function"
		line0 = fs.fn.SrcPos
	)
	if line0 != 0 {
		where = fmt.Sprintf("function at line %d", line0)
	}
	msg := fmt.Sprintf("too many %s (limit is %d) in %s", what, limit, where)
	ls.syntaxErr(msg)
}

// add src:line information to message
func addErrInfo(msg, src string, line int) string {
	chunkID := "?"
	if len(src) > 0 {
		switch src[0] {
		case '=':
			if len(src) <= maxID {
				chunkID = src[1:]
			} else {
				chunkID = src[1:maxID]
			}
		case '@':
			if len(src) <= maxID {
				chunkID = src[1:]
			} else {
				chunkID = "..." + src[1:maxID-3]
			}
		default:
			src = strings.Split(src, "\n")[0]
			tpl := "[string \"...\"]"
			if len(src) > maxID-len(tpl) {
				chunkID = "[string \"" + src + "...\"]"
			} else {
				chunkID = "[string \"" + src + "\"]"
			}
		}
	}
	return fmt.Sprintf("%s:%d: %s", chunkID, line, msg)
}

func checkLimit(fs *function, value, limit int, what string) {
	if value > limit {
		fs.ls.limitErr(fs, limit, what)
	}
}

func (ls *lexical) assert(cond bool) {
	if !cond {
		panic(addErrInfo("failed assertion!", ls.name, ls.line))
	}
}
