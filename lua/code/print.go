package code

import (
	"fmt"
	"io"
)

func printHeader(w io.Writer, fn *Proto) {
	srcID := fn.Source
	if srcID == "" {
		srcID = "=?"
	}
	switch {
	case srcID[0] == '@' || srcID[0] == '=':
		srcID = srcID[1:]
	case srcID[0] == signature[0]:
		srcID = "(bstring)"
	default:
		srcID = "(string)"
	}
	prefix := "function"
	if fn.SrcPos == 0 {
		prefix = "main"
	}
	fmt.Fprintf(w, "\n%s <%s:%d,%d> (%s at %p)\n",
		prefix,
		srcID,
		fn.SrcPos,
		fn.EndPos,
		plural(len(fn.Instrs), "%d instruction%s"),
		fn,
	)
	if fn.Vararg {
		fmt.Fprintf(w, "%s, ", plural(fn.ParamN, "%d+ param%s"))
	} else {
		fmt.Fprintf(w, "%s, ", plural(fn.ParamN, "%d param%s"))
	}
	fmt.Fprintf(w, "%s, %s, ",
		plural(fn.StackN, "%d slot%s"),
		plural(len(fn.UpVars), "%d upvalue%s"),
	)
	fmt.Fprintf(w, "%s, %s, %s\n",
		plural(len(fn.Locals), "%d local%s"),
		plural(len(fn.Consts), "%d constant%s"),
		plural(len(fn.Protos), "%d function%s"),
	)
}

func printConst(w io.Writer, kst Const) {
	switch kst := kst.(type) {
	case string:
		fmt.Fprintf(w, "%q", kst)
	case float64:
		fmt.Fprintf(w, "%f", kst)
	case bool:
		fmt.Fprintf(w, "%t", kst)
	case int64:
		fmt.Fprintf(w, "%d", kst)
	default: // cannot happen
		if kst == nil {
			fmt.Fprintf(w, "nil")
		} else {
			fmt.Fprintf(w, "? type=%T", kst)
		}
	}
}

func printDebug(w io.Writer, fn *Proto) {
	fmt.Fprintf(w, "constants (%d) for %p:\n", len(fn.Consts), fn)
	for i, kst := range fn.Consts {
		fmt.Fprintf(w, "\t%d\t", i+1)
		printConst(w, kst)
		fmt.Fprintln(w)
	}
	fmt.Fprintf(w, "locals (%d) for %p:\n", len(fn.Locals), fn)
	for i, loc := range fn.Locals {
		fmt.Fprintf(w, "\t%d\t%s\t%d\t%d\n",
			i,
			loc.Name,
			loc.Live+1,
			loc.Dead+1,
		)
	}
	fmt.Fprintf(w, "upvalues (%d) for %p:\n", len(fn.UpVars), fn)
	for i, up := range fn.UpVars {
		fmt.Fprintf(w, "\t%d\t%s\t%d\t%d\n",
			i,
			up.Name,
			b2i(up.Stack),
			up.Index,
		)
	}
}

func printCode(w io.Writer, fn *Proto) {
	for pc := 0; pc < len(fn.Instrs); pc++ {
		inst := fn.Instrs[pc]
		mode := inst.Code().Mode()
		mask := inst.Code().Mask()

		fmt.Fprintf(w, "\t%d\t", pc+1)
		if line := funcLine(fn, pc); line > 0 {
			fmt.Fprintf(w, "[%d]\t", line)
		} else {
			fmt.Fprint(w, "[-]\t")
		}
		fmt.Fprintf(w, "%-9s\t", inst.Code())

		switch mode {
		case ModeAsBx:
			fmt.Fprintf(w, "%d %d", inst.A(), inst.SBX())
		case ModeABC:
			fmt.Fprintf(w, "%d", inst.A())
			if !mask.B(ArgN) {
				if IsKst(inst.B()) {
					fmt.Fprintf(w, " %d", Kst(ToKst(inst.B())))
				} else {
					fmt.Fprintf(w, " %d", inst.B())
				}
			}
			if !mask.C(ArgN) {
				if IsKst(inst.C()) {
					fmt.Fprintf(w, " %d", Kst(ToKst(inst.C())))
				} else {
					fmt.Fprintf(w, " %d", inst.C())
				}
			}
		case ModeABx:
			fmt.Fprintf(w, "%d", inst.A())
			if mask.B(ArgK) {
				fmt.Fprintf(w, " %d", Kst(inst.BX()))
			}
			if mask.B(ArgU) {
				fmt.Fprintf(w, " %d", inst.BX())
			}
		case ModeAx:
			fmt.Fprintf(w, "%d", Kst(inst.AX()))
		}

		switch inst.Code() {
		case GETUPVAL, SETUPVAL:
			fmt.Fprintf(w, "\t; %s", upVarID(fn.UpVars[inst.B()]))
		case GETTABLE, SELF:
			if IsKst(inst.C()) {
				fmt.Fprint(w, "\t; ")
				printConst(w, fn.Consts[ToKst(inst.C())])
			}
		case SETTABUP:
			fmt.Fprintf(w, "\t; %s", upVarID(fn.UpVars[inst.A()]))
			if IsKst(inst.B()) {
				fmt.Fprint(w, " ")
				printConst(w, fn.Consts[ToKst(inst.B())])
			}
			if IsKst(inst.C()) {
				fmt.Fprint(w, " ")
				printConst(w, fn.Consts[ToKst(inst.C())])
			}
		case GETTABUP:
			fmt.Fprintf(w, "\t; %s", upVarID(fn.UpVars[inst.B()]))
			if IsKst(inst.C()) {
				fmt.Fprint(w, " ")
				printConst(w, fn.Consts[ToKst(inst.C())])
			}
		case SETTABLE,
			ADD,
			SUB,
			MUL,
			POW,
			DIV,
			IDIV,
			BAND,
			BOR,
			BXOR,
			SHL,
			SHR,
			EQ,
			LT,
			LE:
			if IsKst(inst.B()) || IsKst(inst.C()) {
				fmt.Fprint(w, "\t; ")
				if IsKst(inst.B()) {
					printConst(w, fn.Consts[ToKst(inst.B())])
				} else {
					fmt.Fprint(w, "-")
				}
				fmt.Fprint(w, " ")
				if IsKst(inst.C()) {
					printConst(w, fn.Consts[ToKst(inst.C())])
				} else {
					fmt.Fprint(w, "-")
				}
			}
		case TFORLOOP,
			FORPREP,
			FORLOOP,
			JMP:
			fmt.Fprintf(w, "\t; to %d", inst.SBX()+pc+2)
		case CLOSURE:
			fmt.Fprintf(w, "\t; %p", fn.Protos[inst.BX()])
		case SETLIST:
			if inst.C() == 0 {
				fmt.Fprintf(w, "\t; %d", fn.Instrs[pc+1])
			} else {
				fmt.Fprintf(w, "\t; %d", inst.C())
			}
		case EXTRAARG:
			fmt.Fprint(w, "\t; ")
			printConst(w, fn.Consts[inst.AX()])
		case LOADK:
			fmt.Fprint(w, "\t; ")
			printConst(w, fn.Consts[inst.BX()])
		}
		fmt.Fprintln(w)
	}
}

func printFunc(w io.Writer, fn *Proto, full bool) {
	printHeader(w, fn)
	printCode(w, fn)
	if full {
		printDebug(w, fn)
	}
	for _, fn := range fn.Protos {
		printFunc(w, fn, full)
	}
}

func upVarID(up *UpVar) string {
	if up.Name != "" {
		return up.Name
	}
	return "-"
}

func funcLine(fn *Proto, pc int) int32 {
	if len(fn.PcLine) > 0 {
		return fn.PcLine[pc]
	}
	return -1
}

func plural(n int, format string) string {
	if n == 1 {
		return fmt.Sprintf(format, n, "")
	}
	return fmt.Sprintf(format, n, "s")
}
